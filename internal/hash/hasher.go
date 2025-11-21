package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/thebushidocollective/kasoku/internal/config"
)

type Hasher struct {
	workingDir string
}

type HashResult struct {
	Hash    string
	Details HashDetails
}

type HashDetails struct {
	Files       map[string]string
	Environment map[string]string
	Commands    map[string]string
	Command     string
	Timestamp   string
}

func NewHasher(workingDir string) *Hasher {
	return &Hasher{workingDir: workingDir}
}

func (h *Hasher) ComputeHash(cmd config.Command) (*HashResult, error) {
	hasher := sha256.New()
	details := HashDetails{
		Files:       make(map[string]string),
		Environment: make(map[string]string),
		Commands:    make(map[string]string),
	}

	if err := h.hashCommand(hasher, cmd.Command, &details); err != nil {
		return nil, err
	}

	if err := h.hashFiles(hasher, cmd.Inputs.Files, &details); err != nil {
		return nil, err
	}

	if err := h.hashGlobs(hasher, cmd.Inputs.Globs, &details); err != nil {
		return nil, err
	}

	if err := h.hashEnvironment(hasher, cmd.Inputs.Environment, &details); err != nil {
		return nil, err
	}

	if err := h.hashCommandInputs(hasher, cmd.Inputs.Commands, &details); err != nil {
		return nil, err
	}

	hashBytes := hasher.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)

	return &HashResult{
		Hash:    hashStr,
		Details: details,
	}, nil
}

func (h *Hasher) hashCommand(hasher io.Writer, command string, details *HashDetails) error {
	details.Command = command
	_, err := hasher.Write([]byte(fmt.Sprintf("command:%s\n", command)))
	return err
}

func (h *Hasher) hashFiles(hasher io.Writer, files []string, details *HashDetails) error {
	sort.Strings(files)

	for _, file := range files {
		path := filepath.Join(h.workingDir, file)
		fileHash, err := h.hashFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to hash file %s: %w", file, err)
		}

		details.Files[file] = fileHash
		_, err = hasher.Write([]byte(fmt.Sprintf("file:%s:%s\n", file, fileHash)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hasher) hashGlobs(hasher io.Writer, globs []string, details *HashDetails) error {
	allFiles := make(map[string]bool)

	for _, pattern := range globs {
		negate := strings.HasPrefix(pattern, "!")
		if negate {
			pattern = pattern[1:]
		}

		matches, err := filepath.Glob(filepath.Join(h.workingDir, pattern))
		if err != nil {
			return fmt.Errorf("failed to expand glob %s: %w", pattern, err)
		}

		for _, match := range matches {
			relPath, err := filepath.Rel(h.workingDir, match)
			if err != nil {
				continue
			}

			if negate {
				delete(allFiles, relPath)
			} else {
				info, err := os.Stat(match)
				if err != nil || info.IsDir() {
					continue
				}
				allFiles[relPath] = true
			}
		}
	}

	sortedFiles := make([]string, 0, len(allFiles))
	for file := range allFiles {
		sortedFiles = append(sortedFiles, file)
	}
	sort.Strings(sortedFiles)

	for _, file := range sortedFiles {
		path := filepath.Join(h.workingDir, file)
		fileHash, err := h.hashFile(path)
		if err != nil {
			return fmt.Errorf("failed to hash file %s: %w", file, err)
		}

		details.Files[file] = fileHash
		_, err = hasher.Write([]byte(fmt.Sprintf("glob_file:%s:%s\n", file, fileHash)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hasher) hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (h *Hasher) hashEnvironment(hasher io.Writer, envVars []string, details *HashDetails) error {
	sort.Strings(envVars)

	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		details.Environment[envVar] = value
		_, err := hasher.Write([]byte(fmt.Sprintf("env:%s:%s\n", envVar, value)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hasher) hashCommandInputs(hasher io.Writer, commands []config.CommandInput, details *HashDetails) error {
	for _, cmdInput := range commands {
		output, err := h.executeCommand(cmdInput)
		if err != nil {
			return fmt.Errorf("failed to execute input command %s: %w", cmdInput.Command, err)
		}

		output = strings.TrimSpace(output)
		details.Commands[cmdInput.Command] = output
		_, err = hasher.Write([]byte(fmt.Sprintf("cmd_input:%s:%s\n", cmdInput.Command, output)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Hasher) executeCommand(cmdInput config.CommandInput) (string, error) {
	shell := cmdInput.Shell
	if shell == "" {
		shell = "bash"
	}

	cmd := exec.Command(shell, "-c", cmdInput.Command)
	cmd.Dir = h.workingDir

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
