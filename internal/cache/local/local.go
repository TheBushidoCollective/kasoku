package local

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/thebushidocollective/kasoku/internal/hash"
)

type Cache struct {
	basePath string
	maxSize  int64
}

type CacheEntry struct {
	Hash      string            `json:"hash"`
	Timestamp time.Time         `json:"timestamp"`
	Command   string            `json:"command"`
	Outputs   []string          `json:"outputs"`
	Stdout    string            `json:"stdout"`
	Stderr    string            `json:"stderr"`
	ExitCode  int               `json:"exit_code"`
	Duration  time.Duration     `json:"duration"`
	Details   hash.HashDetails  `json:"details"`
	Metadata  map[string]string `json:"metadata"`
}

func NewCache(basePath string) (*Cache, error) {
	return NewCacheWithSize(basePath, 10*1024*1024*1024)
}

func ParseSize(sizeStr string) (int64, error) {
	original := sizeStr
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))

	units := []struct {
		suffix     string
		multiplier int64
	}{
		{"TB", 1024 * 1024 * 1024 * 1024},
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
		{"B", 1},
	}

	for _, unit := range units {
		if strings.HasSuffix(sizeStr, unit.suffix) {
			numStr := strings.TrimSpace(strings.TrimSuffix(sizeStr, unit.suffix))
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid size format: %s", original)
			}
			return int64(num * float64(unit.multiplier)), nil
		}
	}

	num, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size format: %s", original)
	}
	return num, nil
}

func NewCacheWithSize(basePath string, maxSize int64) (*Cache, error) {
	if basePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(homeDir, ".brisk", "cache")
	}

	basePath = os.ExpandEnv(basePath)

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		basePath: basePath,
		maxSize:  maxSize,
	}, nil
}

func (c *Cache) Get(hash string) (*CacheEntry, error) {
	entryPath := c.entryPath(hash)

	data, err := os.ReadFile(filepath.Join(entryPath, "metadata.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read cache entry: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to parse cache entry: %w", err)
	}

	return &entry, nil
}

func (c *Cache) Put(hash string, entry *CacheEntry) error {
	entryPath := c.entryPath(hash)

	if err := os.MkdirAll(entryPath, 0755); err != nil {
		return fmt.Errorf("failed to create entry directory: %w", err)
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize entry: %w", err)
	}

	metadataPath := filepath.Join(entryPath, "metadata.json")
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

func (c *Cache) StoreOutputs(hash string, outputPaths []string, workingDir string) error {
	entryPath := c.entryPath(hash)
	artifactsPath := filepath.Join(entryPath, "artifacts.tar.gz")

	f, err := os.Create(artifactsPath)
	if err != nil {
		return fmt.Errorf("failed to create artifacts archive: %w", err)
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	for _, outputPath := range outputPaths {
		fullPath := filepath.Join(workingDir, outputPath)
		if err := c.addToTar(tw, fullPath, outputPath); err != nil {
			return fmt.Errorf("failed to archive %s: %w", outputPath, err)
		}
	}

	return nil
}

func (c *Cache) RestoreOutputs(hash string, workingDir string) error {
	entryPath := c.entryPath(hash)
	artifactsPath := filepath.Join(entryPath, "artifacts.tar.gz")

	f, err := os.Open(artifactsPath)
	if err != nil {
		return fmt.Errorf("failed to open artifacts archive: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}

		target := filepath.Join(workingDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		}
	}

	return nil
}

func (c *Cache) addToTar(tw *tar.Writer, sourcePath, tarPath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if info.IsDir() {
		entries, err := os.ReadDir(sourcePath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			srcPath := filepath.Join(sourcePath, entry.Name())
			tPath := filepath.Join(tarPath, entry.Name())
			if err := c.addToTar(tw, srcPath, tPath); err != nil {
				return err
			}
		}
		return nil
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = tarPath

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	f, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(tw, f)
	return err
}

func (c *Cache) entryPath(hash string) string {
	return filepath.Join(c.basePath, hash[:2], hash[2:4], hash)
}

func (c *Cache) Exists(hash string) bool {
	entry, err := c.Get(hash)
	return err == nil && entry != nil
}

func (c *Cache) List() ([]*CacheEntry, error) {
	var entries []*CacheEntry

	err := filepath.Walk(c.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "metadata.json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			var entry CacheEntry
			if err := json.Unmarshal(data, &entry); err != nil {
				return nil
			}

			entries = append(entries, &entry)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *Cache) Delete(hash string) error {
	entryPath := c.entryPath(hash)

	if err := os.RemoveAll(entryPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}

	return nil
}

func (c *Cache) GetSize() (int64, error) {
	var totalSize int64

	err := filepath.Walk(c.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalSize, nil
}

func (c *Cache) EnforceMaxSize() error {
	currentSize, err := c.GetSize()
	if err != nil {
		return fmt.Errorf("failed to get cache size: %w", err)
	}

	if currentSize <= c.maxSize {
		return nil
	}

	entries, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list cache: %w", err)
	}

	type entryWithSize struct {
		entry *CacheEntry
		size  int64
	}

	entriesWithSize := make([]entryWithSize, 0, len(entries))
	for _, entry := range entries {
		entryPath := c.entryPath(entry.Hash)
		size, err := c.getDirSize(entryPath)
		if err != nil {
			continue
		}
		entriesWithSize = append(entriesWithSize, entryWithSize{entry: entry, size: size})
	}

	for i := 0; i < len(entriesWithSize)-1; i++ {
		for j := 0; j < len(entriesWithSize)-i-1; j++ {
			if entriesWithSize[j].entry.Timestamp.After(entriesWithSize[j+1].entry.Timestamp) {
				entriesWithSize[j], entriesWithSize[j+1] = entriesWithSize[j+1], entriesWithSize[j]
			}
		}
	}

	for _, ews := range entriesWithSize {
		if currentSize <= c.maxSize {
			break
		}

		if err := c.Delete(ews.entry.Hash); err != nil {
			continue
		}

		currentSize -= ews.size
	}

	return nil
}

func (c *Cache) getDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func (c *Cache) GetBasePath() string {
	return c.basePath
}
