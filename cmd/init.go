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

package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/release-tools/since/cfg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//go:embed templates/since.yaml
var defaultConfig string

var initSubCmd struct {
	outputFile string
	force      bool
}

// initSubCmd represents the init subcommand
var initSubCmdCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new since.yaml config file",
	Long: `Creates a new since.yaml config file with example configuration,
including branch requirements, pre/post hook scripts, and commit exclusions.
If a config file already exists, the command fails unless --force is set.`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	configDir := initSubCmd.outputFile
	if configDir == "" {
		var err error
		configDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	if !initSubCmd.force {
		for _, name := range cfg.SupportedConfigFiles {
			existing := filepath.Join(configDir, name)
			if _, err := os.Stat(existing); err == nil {
				return fmt.Errorf("config file '%s' already exists; use --force to overwrite", existing)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to check config file: %w", err)
			}
		}
	}

	configPath := filepath.Join(configDir, cfg.DefaultConfigFile)
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	logrus.Infof("created config file '%s'", configPath)
	return nil
}

func init() {
	rootCmd.AddCommand(initSubCmdCmd)

	initSubCmdCmd.Flags().StringVarP(&initSubCmd.outputFile, "output", "o", "", "Directory to write the config file to (default: current directory)")
	initSubCmdCmd.Flags().BoolVarP(&initSubCmd.force, "force", "f", false, "Overwrite an existing config file")
}
