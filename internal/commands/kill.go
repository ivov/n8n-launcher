package commands

import (
	"fmt"
	"os"
	"syscall"
	"task-runner-launcher/internal/config"
)

type KillCommand struct {
	RunnerType string
	PID        int
}

func (k *KillCommand) Execute() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	found := false
	for _, r := range cfg.TaskRunners {
		if r.RunnerType == k.RunnerType {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("unknown runner type: %s", k.RunnerType)
	}

	process, err := os.FindProcess(k.PID)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	return process.Signal(syscall.SIGTERM)
}
