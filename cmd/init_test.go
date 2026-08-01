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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/release-tools/since/cfg"
)

func Test_runInit(t *testing.T) {
	t.Run("creates config in current directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		defer os.Chdir(originalWd)
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		err = runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		configPath := filepath.Join(tmpDir, "since.yaml")
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if len(content) == 0 {
			t.Error("config file is empty")
		}
	})

	t.Run("creates config in specified output directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = tmpDir
		defer func() { initSubCmd.outputFile = "" }()

		err := runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		configPath := filepath.Join(tmpDir, "since.yaml")
		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if len(content) == 0 {
			t.Error("config file is empty")
		}
	})

	t.Run("embeds requireBranch example", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		defer os.Chdir(originalWd)
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		err = runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		content, err := os.ReadFile("since.yaml")
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if !strings.Contains(string(content), "requireBranch:") {
			t.Error("config file does not contain requireBranch example")
		}
	})

	t.Run("embeds hook examples", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		defer os.Chdir(originalWd)
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		err = runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		content, err := os.ReadFile("since.yaml")
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if !strings.Contains(string(content), "command:") {
			t.Error("config file does not contain command hook example")
		}
	})

	t.Run("embeds ignore examples", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		defer os.Chdir(originalWd)
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		err = runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		content, err := os.ReadFile("since.yaml")
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if !strings.Contains(string(content), "ignore:") {
			t.Error("config file does not contain ignore example")
		}
	})

	t.Run("embeds script-based hook example", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("failed to get working directory: %v", err)
		}
		defer os.Chdir(originalWd)
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatalf("failed to chdir: %v", err)
		}

		err = runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		content, err := os.ReadFile("since.yaml")
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if !strings.Contains(string(content), "script:") {
			t.Error("config file does not contain script-based hook example")
		}
	})

	t.Run("refuses to overwrite existing config file without force", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = tmpDir
		defer func() { initSubCmd.outputFile = "" }()

		configPath := filepath.Join(tmpDir, "since.yaml")
		if err := os.WriteFile(configPath, []byte("stale: content"), 0644); err != nil {
			t.Fatalf("failed to seed existing config file: %v", err)
		}

		err := runInit()
		if err == nil {
			t.Fatal("runInit() expected an error when config file exists, got nil")
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}
		if !strings.Contains(string(content), "stale: content") {
			t.Error("existing config file was overwritten despite no force flag")
		}
	})

	t.Run("refuses to overwrite existing config file with alternate extension", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = tmpDir
		defer func() { initSubCmd.outputFile = "" }()

		ymlPath := filepath.Join(tmpDir, "since.yml")
		if err := os.WriteFile(ymlPath, []byte("stale: content"), 0644); err != nil {
			t.Fatalf("failed to seed existing config file: %v", err)
		}

		err := runInit()
		if err == nil {
			t.Fatal("runInit() expected an error when since.yml exists, got nil")
		}

		if _, err := os.Stat(filepath.Join(tmpDir, "since.yaml")); !os.IsNotExist(err) {
			t.Error("since.yaml was created despite existing since.yml and no force flag")
		}
		content, err := os.ReadFile(ymlPath)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}
		if !strings.Contains(string(content), "stale: content") {
			t.Error("existing since.yml was modified despite no force flag")
		}
	})

	t.Run("overwrites existing config file with force", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = tmpDir
		initSubCmd.force = true
		defer func() {
			initSubCmd.outputFile = ""
			initSubCmd.force = false
		}()

		configPath := filepath.Join(tmpDir, "since.yaml")
		if err := os.WriteFile(configPath, []byte("stale: content"), 0644); err != nil {
			t.Fatalf("failed to seed existing config file: %v", err)
		}

		err := runInit()
		if err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read created config file: %v", err)
		}

		if strings.Contains(string(content), "stale: content") {
			t.Error("config file was not overwritten")
		}
		if !strings.Contains(string(content), "requireBranch:") {
			t.Error("overwritten config file does not contain template content")
		}
	})

	t.Run("returns error when config check cannot stat the path", func(t *testing.T) {
		tmpDir := t.TempDir()

		// point the output at a path whose parent is a regular file, so statting
		// the config file yields a non-IsNotExist error (ENOTDIR).
		notADir := filepath.Join(tmpDir, "regular-file")
		if err := os.WriteFile(notADir, []byte("not a directory"), 0644); err != nil {
			t.Fatalf("failed to seed file: %v", err)
		}

		initSubCmd.outputFile = notADir
		defer func() { initSubCmd.outputFile = "" }()

		err := runInit()
		if err == nil {
			t.Fatal("runInit() expected an error when the config path cannot be checked, got nil")
		}
	})

	t.Run("returns error when output directory does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = filepath.Join(tmpDir, "does-not-exist")
		defer func() { initSubCmd.outputFile = "" }()

		err := runInit()
		if err == nil {
			t.Fatal("runInit() expected an error when writing to a missing directory, got nil")
		}
	})

	t.Run("generated config loads through the real loader", func(t *testing.T) {
		tmpDir := t.TempDir()

		initSubCmd.outputFile = tmpDir
		defer func() { initSubCmd.outputFile = "" }()

		if err := runInit(); err != nil {
			t.Fatalf("runInit() unexpected error: %v", err)
		}

		// The template ships with every example commented out, so the loader
		// should parse it without error and yield an empty (default) config.
		// This guards against invalid YAML or accidentally uncommented lines
		// slipping into the template.
		config, err := cfg.LoadConfig(tmpDir)
		if err != nil {
			t.Fatalf("generated config failed to load: %v", err)
		}

		if config.RequireBranch != "" {
			t.Errorf("expected empty RequireBranch, got %q", config.RequireBranch)
		}
		if len(config.Before) != 0 || len(config.After) != 0 || len(config.Ignore) != 0 {
			t.Error("expected no active hooks or ignore patterns in default template")
		}
	})
}
