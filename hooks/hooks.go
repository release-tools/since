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

package hooks

import (
	"fmt"
	"github.com/outofcoffee/since/cfg"
	"github.com/outofcoffee/since/vcs"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type HookType string

const (
	Before HookType = "before"
	After  HookType = "after"
)

// ExecuteHooks executes all hooks of the given type
func ExecuteHooks(config cfg.SinceConfig, hookType HookType, metadata vcs.ReleaseMetadata) error {
	var hooks []cfg.Hook
	switch hookType {
	case Before:
		hooks = config.Before
	case After:
		hooks = config.After
	default:
		return fmt.Errorf("invalid hook type: %s", hookType)
	}

	logrus.Tracef("%d %v hooks found", len(hooks), hookType)
	for _, hook := range hooks {
		err := executeHook(hook, metadata)
		if err != nil {
			return fmt.Errorf("error executing hook '%s %s': %v", hook.Command, strings.Join(hook.Args, " "), err)
		}
	}

	return nil
}

// executeHook executes a hook command with the given arguments
func executeHook(hook cfg.Hook, metadata vcs.ReleaseMetadata) error {
	logrus.Debugf("executing hook '%s %s'", hook.Command, strings.Join(hook.Args, " "))

	cmd := exec.Command(hook.Command, hook.Args...)
	cmd.Dir = metadata.RepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), []string{
		"SINCE_NEW_VERSION=" + metadata.NewVersion,
		"SINCE_OLD_VERSION=" + metadata.OldVersion,
		"SINCE_SHA=" + metadata.Sha,
	}...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing hook '%s %s': %v", hook.Command, strings.Join(hook.Args, " "), err)
	}
	if cmd.ProcessState.Success() {
		logrus.Debugf("hook '%s %s' executed successfully", hook.Command, strings.Join(hook.Args, " "))
	} else {
		logrus.Warnf("hook '%s %s' executed with errors", hook.Command, strings.Join(hook.Args, " "))
	}
	return nil
}
