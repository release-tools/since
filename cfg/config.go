package cfg

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type Hook struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type SinceConfig struct {
	Before        []Hook `json:"before"`
	After         []Hook `json:"after"`
	RequireBranch string `json:"requireBranch"`
}

const defaultConfigFile = "since.yaml"

// LoadConfig loads the YAML config file from the given directory.
func LoadConfig(dir string) (SinceConfig, error) {
	return loadConfig(path.Join(dir, defaultConfigFile))
}

// loadConfig loads the YAML config file from the given path.
// If the file does not exist, an empty config is returned.
func loadConfig(configPath string) (SinceConfig, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Tracef("config file '%s' not found", configPath)
		return SinceConfig{}, nil
	}

	var config SinceConfig
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return SinceConfig{}, fmt.Errorf("error: %v", err)
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return SinceConfig{}, fmt.Errorf("error: %v", err)
	}
	return config, nil
}
