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
}
