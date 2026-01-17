package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Server struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		ReadinessPath string `yaml:"readinessPath"`
		LivenessPath  string `yaml:"livenessPath"`
	} `yaml:"server"`
	Resources struct {
		Memory  string `yaml:"memory"`
		CPU     string `yaml:"cpu"`
		Storage string `yaml:"storage"`
		Threads int    `yaml:"threads"`
	} `yaml:"resources"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
}

func LoadConfig(cfg *Config, path string) error {
	configFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if err := configFile.Close(); err != nil {
			log.Printf("Error closing config file: %v", err)
		}
	}()

	content, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	expanded := os.ExpandEnv(string(content))

	err = yaml.Unmarshal([]byte(expanded), cfg)
	if err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	return nil
}
