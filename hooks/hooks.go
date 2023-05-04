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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"SINCE_NEW_VERSION=" + metadata.NewVersion,
		"SINCE_OLD_VERSION=" + metadata.OldVersion,
		"SINCE_SHA=" + metadata.Sha,
	}
	cmd.Dir = metadata.RepoPath

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
