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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootArgs struct {
	logLevel string
	quiet    bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "since",
	Short: "Parse git history and changelog files",
	Long:  `Parses git logs and changelog files and lists changes for a given version.`,
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVarP(&rootArgs.logLevel, "log-level", "l", "debug", "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().BoolVarP(&rootArgs.quiet, "quiet", "q", false, "Disable logging (useful for scripting)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initLogging() {
	if rootArgs.logLevel != "" {
		var logLevel logrus.Level
		if rootArgs.quiet {
			logLevel = logrus.PanicLevel
		} else {
			ll, err := logrus.ParseLevel(rootArgs.logLevel)
			if err != nil {
				ll = logrus.DebugLevel
			}
			logLevel = ll
		}
		logrus.SetLevel(logLevel)
	}
}
