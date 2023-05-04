/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	Before        []Hook   `yaml:"before"`
	After         []Hook   `yaml:"after"`
	RequireBranch string   `yaml:"requireBranch"`
	Ignore        []string `yaml:"ignore"`
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
