package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type TaskRunnerConfig struct {
	RunnerType string   `json:"runner-type"`
	WorkDir    string   `json:"workdir"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
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
		return nil, fmt.Errorf("failed to open config file: %s: %w", configPath, err)
	}

	var config LauncherConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse launcher config file: %w", err)
	}

	return &config, nil
}
