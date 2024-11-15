// Package config provides functions to use the launcher configuration file.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type TaskRunnerConfig struct {
	// Type of runner, currently only "javascript" supported
	RunnerType string `json:"runner-type"`

	// Path to directory containing launcher (Go binary)
	WorkDir string `json:"workdir"`

	// Command to execute to start runner
	Command string `json:"command"`

	// Arguments for command to execute, currently path to task runner entrypoint
	Args []string `json:"args"`

	// Env vars allowed to be passed by launcher to task runner
	AllowedEnv []string `json:"allowed-env"`
}

type LauncherConfig struct {
	TaskRunners []TaskRunnerConfig `json:"task-runners"`
}

func getConfigPath() string {
	if os.Getenv("SECURE_MODE") == "true" {
		return "/etc/n8n-task-runners.json"
	}
	return "./config.json"
}

func ReadConfig() (*LauncherConfig, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file at %s: %w", configPath, err)
	}

	var config LauncherConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file at %s: %w", configPath, err)
	}

	return &config, nil
}
