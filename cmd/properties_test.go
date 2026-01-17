package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tempFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(tempFile.Name()) })

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name()
}

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name          string
		setup         func(t *testing.T)
		configFile    string
		configContent string
		expectErr     bool
		validate      func(t *testing.T, cfg *Config)
	}{
		{
			name:       "loads config from file",
			configFile: filepath.Join("..", "resources", "test.yaml"),
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Database.Host != "localhost" {
					t.Errorf("Expected host 'localhost', got %s", cfg.Database.Host)
				}
				if cfg.Database.Port != 5432 {
					t.Errorf("Expected port 5432, got %d", cfg.Database.Port)
				}
				if cfg.Server.Port != 8081 {
					t.Errorf("Expected server port 8081, got %d", cfg.Server.Port)
				}
				if cfg.Logging.Level != "debug" {
					t.Errorf("Expected log level 'debug', got %s", cfg.Logging.Level)
				}
			},
		},
		{
			name: "loads config with environment variable substitution",
			setup: func(t *testing.T) {
				t.Setenv("TEST_HOST", "envhost")
				t.Setenv("TEST_PORT", "5678")
			},
			configContent: `database:
  host: ${TEST_HOST}
  port: ${TEST_PORT}
`,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Database.Host != "envhost" {
					t.Errorf("Expected host 'envhost', got %s", cfg.Database.Host)
				}
				if cfg.Database.Port != 5678 {
					t.Errorf("Expected port 5678, got %d", cfg.Database.Port)
				}
			},
		},
		{
			name:       "returns error when config file not found",
			configFile: "non-existent-file.yaml",
			expectErr:  true,
		},
		{
			name:          "returns error for invalid YAML",
			configContent: `database: { port: "not-a-number" }`,
			expectErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}

			configFile := tc.configFile
			if tc.configContent != "" {
				configFile = createTempConfigFile(t, tc.configContent)
			}

			var cfg Config
			err := LoadConfig(&cfg, configFile)

			if tc.expectErr {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if tc.validate != nil {
				tc.validate(t, &cfg)
			}
		})
	}
}
