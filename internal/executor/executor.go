package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/thebushidocollective/kasoku/internal/analytics"
	"github.com/thebushidocollective/kasoku/internal/cache/local"
	"github.com/thebushidocollective/kasoku/internal/config"
	"github.com/thebushidocollective/kasoku/internal/hash"
	"github.com/thebushidocollective/kasoku/internal/storage"
)

type Executor struct {
	config        *config.Config
	localCache    *local.Cache
	remoteStorage storage.Backend
	analytics     *analytics.Analytics
}

type ExecutionResult struct {
	CacheHit  bool
	Hash      string
	Stdout    string
	Stderr    string
	ExitCode  int
	Duration  time.Duration
	TimeSaved time.Duration
}

func NewExecutor(cfg *config.Config) (*Executor, error) {
	var localCache *local.Cache
	var remoteStorage storage.Backend
	var err error

	if cfg.Cache.Local.Enabled {
		maxSize, err := local.ParseSize(cfg.Cache.Local.MaxSize)
		if err != nil {
			return nil, fmt.Errorf("failed to parse max cache size: %w", err)
		}

		localCache, err = local.NewCacheWithSize(cfg.Cache.Local.Path, maxSize)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize local cache: %w", err)
		}
	}

	if cfg.Cache.Remote.Enabled {
		storageCfg := storage.Config{
			Type:           cfg.Cache.Remote.Type,
			CustomURL:      cfg.Cache.Remote.URL,
			CustomToken:    cfg.Cache.Remote.Token,
			FilesystemPath: cfg.Cache.Remote.FilesystemPath,
			S3Bucket:       cfg.Cache.Remote.S3Bucket,
			S3Region:       cfg.Cache.Remote.S3Region,
			S3Endpoint:     cfg.Cache.Remote.S3Endpoint,
			GCSBucket:      cfg.Cache.Remote.GCSBucket,
			GCSCredentials: cfg.Cache.Remote.GCSCredentials,
			AzureAccount:   cfg.Cache.Remote.AzureAccount,
			AzureContainer: cfg.Cache.Remote.AzureContainer,
			AzureKey:       cfg.Cache.Remote.AzureKey,
		}

		remoteStorage, err = storage.NewBackend(storageCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize remote storage: %w", err)
		}
	}

	analytics, err := analytics.New("")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize analytics: %w", err)
	}

	return &Executor{
		config:        cfg,
		localCache:    localCache,
		remoteStorage: remoteStorage,
		analytics:     analytics,
	}, nil
}

func (e *Executor) Execute(commandName string) (*ExecutionResult, error) {
	cmd, ok := e.config.Commands[commandName]
	if !ok {
		return nil, fmt.Errorf("command %q not found in config", commandName)
	}

	workingDir := cmd.WorkingDir
	if workingDir == "" {
		workingDir = "."
	}

	hasher := hash.NewHasher(workingDir)
	hashResult, err := hasher.ComputeHash(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hash: %w", err)
	}

	fmt.Printf("🔍 Hash: %s\n", hashResult.Hash)

	if e.localCache != nil {
		if entry := e.checkCache(hashResult.Hash); entry != nil {
			return e.restoreFromCache(entry, cmd, workingDir)
		}
	}

	return e.executeCommand(cmd, hashResult, workingDir)
}

func (e *Executor) ExecuteRaw(rawCommand string) (*ExecutionResult, error) {
	cmdConfig, actualCommand, ok := e.config.MatchCommand(rawCommand)
	if !ok {
		return e.executePassthrough(rawCommand)
	}

	cmdWithActual := *cmdConfig
	cmdWithActual.Command = actualCommand

	workingDir := cmdConfig.WorkingDir
	if workingDir == "" {
		workingDir = "."
	}

	hasher := hash.NewHasher(workingDir)
	hashResult, err := hasher.ComputeHash(cmdWithActual)
	if err != nil {
		return nil, fmt.Errorf("failed to compute hash: %w", err)
	}

	fmt.Printf("🔍 Hash: %s\n", hashResult.Hash)

	if e.localCache != nil {
		if entry := e.checkCache(hashResult.Hash); entry != nil {
			return e.restoreFromCache(entry, cmdWithActual, workingDir)
		}
	}

	return e.executeCommand(cmdWithActual, hashResult, workingDir)
}

func (e *Executor) executePassthrough(rawCommand string) (*ExecutionResult, error) {
	startTime := time.Now()

	shellCmd := exec.Command("bash", "-c", rawCommand)
	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr

	err := shellCmd.Run()
	duration := time.Since(startTime)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	return &ExecutionResult{
		CacheHit:  false,
		Hash:      "",
		Stdout:    "",
		Stderr:    "",
		ExitCode:  exitCode,
		Duration:  duration,
		TimeSaved: 0,
	}, nil
}

func (e *Executor) checkCache(hash string) *local.CacheEntry {
	if e.localCache == nil {
		return nil
	}

	entry, err := e.localCache.Get(hash)
	if err != nil {
		fmt.Printf("⚠️  Failed to check cache: %v\n", err)
		return nil
	}

	if entry != nil {
		return entry
	}

	if e.remoteStorage != nil {
		fmt.Println("🔍 Checking remote cache...")
		if remoteEntry := e.pullFromRemote(hash); remoteEntry != nil {
			return remoteEntry
		}
	}

	return nil
}

func (e *Executor) pullFromRemote(hash string) *local.CacheEntry {
	ctx := context.Background()

	metadataKey := fmt.Sprintf("%s/metadata.json", hash)
	exists, err := e.remoteStorage.Exists(ctx, metadataKey)
	if err != nil || !exists {
		return nil
	}

	reader, err := e.remoteStorage.Get(ctx, metadataKey)
	if err != nil {
		fmt.Printf("⚠️  Failed to get remote metadata: %v\n", err)
		return nil
	}
	defer reader.Close()

	var entry local.CacheEntry
	if err := json.NewDecoder(reader).Decode(&entry); err != nil {
		fmt.Printf("⚠️  Failed to decode remote metadata: %v\n", err)
		return nil
	}

	artifactsKey := fmt.Sprintf("%s/artifacts.tar.gz", hash)
	artifactsReader, err := e.remoteStorage.Get(ctx, artifactsKey)
	if err != nil {
		fmt.Printf("⚠️  Failed to get remote artifacts: %v\n", err)
		return nil
	}
	defer artifactsReader.Close()

	if e.localCache != nil {
		entryPath := filepath.Join(e.localCache.GetBasePath(), hash[:2], hash[2:4], hash)
		if err := os.MkdirAll(entryPath, 0755); err != nil {
			fmt.Printf("⚠️  Failed to create local cache directory: %v\n", err)
			return &entry
		}

		localArtifactsPath := filepath.Join(entryPath, "artifacts.tar.gz")
		f, err := os.Create(localArtifactsPath)
		if err != nil {
			fmt.Printf("⚠️  Failed to create local artifacts file: %v\n", err)
			return &entry
		}
		defer f.Close()

		if _, err := io.Copy(f, artifactsReader); err != nil {
			fmt.Printf("⚠️  Failed to copy artifacts to local cache: %v\n", err)
			return &entry
		}

		if err := e.localCache.Put(hash, &entry); err != nil {
			fmt.Printf("⚠️  Failed to store metadata locally: %v\n", err)
		} else {
			fmt.Println("📥 Pulled from remote cache")
		}
	}

	return &entry
}

func (e *Executor) pushToRemote(hash string, entry *local.CacheEntry) error {
	if e.remoteStorage == nil {
		return nil
	}

	ctx := context.Background()

	metadataBytes, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataKey := fmt.Sprintf("%s/metadata.json", hash)
	if err := e.remoteStorage.Put(ctx, metadataKey, bytes.NewReader(metadataBytes)); err != nil {
		return fmt.Errorf("failed to push metadata: %w", err)
	}

	if e.localCache != nil {
		entryPath := filepath.Join(e.localCache.GetBasePath(), hash[:2], hash[2:4], hash)
		artifactsPath := filepath.Join(entryPath, "artifacts.tar.gz")

		f, err := os.Open(artifactsPath)
		if err != nil {
			return fmt.Errorf("failed to open artifacts: %w", err)
		}
		defer f.Close()

		artifactsKey := fmt.Sprintf("%s/artifacts.tar.gz", hash)
		if err := e.remoteStorage.Put(ctx, artifactsKey, f); err != nil {
			return fmt.Errorf("failed to push artifacts: %w", err)
		}
	}

	fmt.Println("📤 Pushed to remote cache")
	return nil
}

func (e *Executor) restoreFromCache(entry *local.CacheEntry, cmd config.Command, workingDir string) (*ExecutionResult, error) {
	fmt.Println("✨ Cache hit! Restoring outputs...")

	startTime := time.Now()

	if err := e.localCache.RestoreOutputs(entry.Hash, workingDir); err != nil {
		return nil, fmt.Errorf("failed to restore outputs: %w", err)
	}

	restoreDuration := time.Since(startTime)

	fmt.Println("\n" + entry.Stdout)
	if entry.Stderr != "" {
		fmt.Fprintln(os.Stderr, entry.Stderr)
	}

	timeSaved := entry.Duration - restoreDuration
	if timeSaved > 0 {
		fmt.Printf("\n⚡ Saved %v (original: %v, restore: %v)\n", timeSaved, entry.Duration, restoreDuration)
	}

	if e.analytics != nil {
		e.analytics.RecordExecution(analytics.ExecutionEvent{
			Command:   cmd.Name,
			CacheHit:  true,
			Duration:  restoreDuration,
			TimeSaved: timeSaved,
			Timestamp: time.Now(),
		})
	}

	return &ExecutionResult{
		CacheHit:  true,
		Hash:      entry.Hash,
		Stdout:    entry.Stdout,
		Stderr:    entry.Stderr,
		ExitCode:  entry.ExitCode,
		Duration:  restoreDuration,
		TimeSaved: timeSaved,
	}, nil
}

func (e *Executor) executeCommand(cmd config.Command, hashResult *hash.HashResult, workingDir string) (*ExecutionResult, error) {
	fmt.Println("🚀 Executing command...")

	startTime := time.Now()

	var stdoutBuf, stderrBuf bytes.Buffer

	shellCmd := exec.Command("bash", "-c", cmd.Command)
	shellCmd.Dir = workingDir
	shellCmd.Env = append(os.Environ(), cmd.Environment...)

	stdoutWriter := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderrWriter := io.MultiWriter(os.Stderr, &stderrBuf)
	shellCmd.Stdout = stdoutWriter
	shellCmd.Stderr = stderrWriter

	err := shellCmd.Run()
	duration := time.Since(startTime)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}
	}

	outputPaths := make([]string, 0, len(cmd.Outputs))
	for _, output := range cmd.Outputs {
		fullPath := output.Path
		if _, err := os.Stat(fullPath); err == nil || !output.Optional {
			outputPaths = append(outputPaths, output.Path)
		}
	}

	if e.localCache != nil && exitCode == 0 {
		entry := &local.CacheEntry{
			Hash:      hashResult.Hash,
			Timestamp: time.Now(),
			Command:   cmd.Command,
			Outputs:   outputPaths,
			Stdout:    stdoutBuf.String(),
			Stderr:    stderrBuf.String(),
			ExitCode:  exitCode,
			Duration:  duration,
			Details:   hashResult.Details,
		}

		if err := e.localCache.Put(hashResult.Hash, entry); err != nil {
			fmt.Printf("⚠️  Failed to cache metadata: %v\n", err)
		}

		if err := e.localCache.StoreOutputs(hashResult.Hash, outputPaths, workingDir); err != nil {
			fmt.Printf("⚠️  Failed to cache outputs: %v\n", err)
		} else {
			fmt.Println("💾 Cached outputs for future use")

			if err := e.localCache.EnforceMaxSize(); err != nil {
				fmt.Printf("⚠️  Failed to enforce cache size limit: %v\n", err)
			}

			if e.remoteStorage != nil {
				if err := e.pushToRemote(hashResult.Hash, entry); err != nil {
					fmt.Printf("⚠️  Failed to push to remote cache: %v\n", err)
				}
			}
		}
	}

	if e.analytics != nil {
		e.analytics.RecordExecution(analytics.ExecutionEvent{
			Command:   cmd.Name,
			CacheHit:  false,
			Duration:  duration,
			TimeSaved: 0,
			Timestamp: time.Now(),
		})
	}

	return &ExecutionResult{
		CacheHit:  false,
		Hash:      hashResult.Hash,
		Stdout:    stdoutBuf.String(),
		Stderr:    stderrBuf.String(),
		ExitCode:  exitCode,
		Duration:  duration,
		TimeSaved: 0,
	}, nil
}

func (e *Executor) ShowHashDetails(commandName string) error {
	cmd, ok := e.config.Commands[commandName]
	if !ok {
		return fmt.Errorf("command %q not found in config", commandName)
	}

	workingDir := cmd.WorkingDir
	if workingDir == "" {
		workingDir = "."
	}

	hasher := hash.NewHasher(workingDir)
	hashResult, err := hasher.ComputeHash(cmd)
	if err != nil {
		return fmt.Errorf("failed to compute hash: %w", err)
	}

	fmt.Printf("Hash: %s\n\n", hashResult.Hash)
	fmt.Println("Command:")
	fmt.Printf("  %s\n\n", hashResult.Details.Command)

	if len(hashResult.Details.Files) > 0 {
		fmt.Println("Files:")
		for file, fileHash := range hashResult.Details.Files {
			fmt.Printf("  %s: %s\n", file, fileHash[:16])
		}
		fmt.Println()
	}

	if len(hashResult.Details.Environment) > 0 {
		fmt.Println("Environment:")
		for key, value := range hashResult.Details.Environment {
			fmt.Printf("  %s=%s\n", key, value)
		}
		fmt.Println()
	}

	if len(hashResult.Details.Commands) > 0 {
		fmt.Println("Command Inputs:")
		for cmd, output := range hashResult.Details.Commands {
			fmt.Printf("  %s → %s\n", cmd, output)
		}
		fmt.Println()
	}

	return nil
}

func (e *Executor) ListCache() error {
	if e.localCache == nil {
		fmt.Println("⚠️  Local cache is not enabled")
		return nil
	}

	entries, err := e.localCache.List()
	if err != nil {
		return fmt.Errorf("failed to list cache: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("📦 Cache is empty")
		return nil
	}

	fmt.Printf("\n📦 Cached Entries (%d total)\n", len(entries))
	fmt.Println("═══════════════════════════════════════════════════════════════")

	for _, entry := range entries {
		fmt.Printf("\nHash:      %s\n", entry.Hash[:16]+"...")
		fmt.Printf("Command:   %s\n", entry.Command)
		fmt.Printf("Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("Duration:  %v\n", entry.Duration)
		fmt.Printf("ExitCode:  %d\n", entry.ExitCode)
		if len(entry.Outputs) > 0 {
			fmt.Printf("Outputs:   %v\n", entry.Outputs)
		}
		fmt.Println("───────────────────────────────────────────────────────────────")
	}

	return nil
}

func (e *Executor) ClearCache() error {
	if e.localCache == nil {
		fmt.Println("⚠️  Local cache is not enabled")
		return nil
	}

	entries, err := e.localCache.List()
	if err != nil {
		return fmt.Errorf("failed to list cache: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("📦 Cache is already empty")
		return nil
	}

	fmt.Printf("🗑️  Clearing %d cache entries...\n", len(entries))

	for _, entry := range entries {
		if err := e.localCache.Delete(entry.Hash); err != nil {
			fmt.Printf("⚠️  Failed to delete %s: %v\n", entry.Hash[:16], err)
		}
	}

	fmt.Println("✅ Cache cleared successfully")
	return nil
}
