package credentials

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Credentials struct {
	URL   string `json:"url"`
	Token string `json:"token"`
	Email string `json:"email"`
}

// GetCredentialsPath returns the path to the credentials file
func GetCredentialsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	kasokuDir := filepath.Join(homeDir, ".kasoku")
	if err := os.MkdirAll(kasokuDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create .kasoku directory: %w", err)
	}

	return filepath.Join(kasokuDir, "credentials.json"), nil
}

// Load reads credentials from disk
func Load() (*Credentials, error) {
	path, err := GetCredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not logged in: run 'kasoku login' first")
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// Save writes credentials to disk
func Save(creds *Credentials) error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// Clear removes stored credentials
func Clear() error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	return nil
}

// Exists checks if credentials are stored
func Exists() bool {
	path, err := GetCredentialsPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}
