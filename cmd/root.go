/*
Copyright © 2023 Pete Cornish <outofcoffee@gmail.com>
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootArgs struct {
	logLevel string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "changelog-parser",
	Short: "Parse changelog files",
	Long:  `Parse changelog file and return changes for a given version.`,
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVarP(&rootArgs.logLevel, "log-level", "l", "debug", "Log level (debug, info, warn, error, fatal, panic)")
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
		ll, err := logrus.ParseLevel(rootArgs.logLevel)
		if err != nil {
			ll = logrus.DebugLevel
		}
		logrus.SetLevel(ll)
	}
}
