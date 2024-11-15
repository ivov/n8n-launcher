package commands

import (
	"fmt"
	"n8n-launcher/internal/config"
	"os"
	"syscall"
)

type KillCommand struct {
	// Type of runner to kill, currently only "javascript" supported
	RunnerType string

	// Process ID of runner to kill
	PID int
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
		return fmt.Errorf("failed to find requested runner type in config: %s", k.RunnerType)
	}

	process, err := os.FindProcess(k.PID)
	if err != nil {
		return fmt.Errorf("failed to find requested process ID: %w", err)
	}

	return process.Signal(syscall.SIGKILL)
}
