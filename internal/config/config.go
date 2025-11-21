package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version  string             `yaml:"version"`
	Commands map[string]Command `yaml:"commands"`
	Cache    CacheConfig        `yaml:"cache"`
}

type Command struct {
	Name        string         `yaml:"name"`
	Command     string         `yaml:"command"`
	Patterns    []string       `yaml:"patterns"`
	Inputs      InputsConfig   `yaml:"inputs"`
	Outputs     []OutputConfig `yaml:"outputs"`
	Environment []string       `yaml:"environment"`
	WorkingDir  string         `yaml:"working_dir"`
}

type InputsConfig struct {
	Files        []string       `yaml:"files"`
	Globs        []string       `yaml:"globs"`
	Environment  []string       `yaml:"environment"`
	Commands     []CommandInput `yaml:"commands"`
	Dependencies []string       `yaml:"dependencies"`
}

type CommandInput struct {
	Command string `yaml:"command"`
	Shell   string `yaml:"shell"`
}

type OutputConfig struct {
	Path     string `yaml:"path"`
	Optional bool   `yaml:"optional"`
}

type CacheConfig struct {
	Local  LocalCacheConfig  `yaml:"local"`
	Remote RemoteCacheConfig `yaml:"remote"`
}

type LocalCacheConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	MaxSize string `yaml:"max_size"`
}

type RemoteCacheConfig struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`

	// For custom kasoku server
	URL   string `yaml:"url"`
	Token string `yaml:"token"`

	// Legacy/deprecated fields
	Endpoint string `yaml:"endpoint"`
	Secure   bool   `yaml:"secure"`

	// S3 backend
	S3Bucket   string `yaml:"s3_bucket"`
	S3Region   string `yaml:"s3_region"`
	S3Endpoint string `yaml:"s3_endpoint"`

	// GCS backend
	GCSBucket      string `yaml:"gcs_bucket"`
	GCSCredentials string `yaml:"gcs_credentials"`

	// Azure backend
	AzureAccount   string `yaml:"azure_account"`
	AzureContainer string `yaml:"azure_container"`
	AzureKey       string `yaml:"azure_key"`

	// Filesystem backend
	FilesystemPath string `yaml:"filesystem_path"`
}

func Load(path string) (*Config, error) {
	configPath, err := findConfig(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	cfg.ApplyDefaults()
	return &cfg, nil
}

func findConfig(startPath string) (string, error) {
	if startPath != "kasoku.yaml" {
		if _, err := os.Stat(startPath); err == nil {
			return startPath, nil
		}
		return "", fmt.Errorf("config file not found: %s", startPath)
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		configPath := filepath.Join(dir, "kasoku.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("kasoku.yaml not found in current directory or any parent directory")
}

func (c *Config) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}

	for name, cmd := range c.Commands {
		if len(cmd.Outputs) == 0 {
			return fmt.Errorf("command %q: at least one output must be defined", name)
		}
	}

	return nil
}

func (c *Config) ApplyDefaults() {
	if c.Cache.Local.Path == "" {
		homeDir, _ := os.UserHomeDir()
		c.Cache.Local.Path = homeDir + "/.kasoku/cache"
	}

	if c.Cache.Local.MaxSize == "" {
		c.Cache.Local.MaxSize = "10GB"
	}

	for name, cmd := range c.Commands {
		if cmd.Name == "" {
			cmd.Name = name
			c.Commands[name] = cmd
		}
		if cmd.WorkingDir == "" {
			cmd.WorkingDir = "."
			c.Commands[name] = cmd
		}
	}
}

func (c *Config) MatchCommand(actualCommand string) (*Command, string, bool) {
	for _, cmd := range c.Commands {
		if len(cmd.Patterns) == 0 {
			return &cmd, actualCommand, true
		}

		for _, pattern := range cmd.Patterns {
			if matchPattern(pattern, actualCommand) {
				return &cmd, actualCommand, true
			}
		}
	}

	return nil, "", false
}

func matchPattern(pattern, command string) bool {
	if pattern == "" {
		return false
	}

	if pattern == "*" {
		return true
	}

	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(command) >= len(prefix) && command[:len(prefix)] == prefix
	}

	return pattern == command
}
