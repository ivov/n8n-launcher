package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"task-runner-launcher/internal/auth"
	"task-runner-launcher/internal/config"
)

type LaunchCommand struct {
	RunnerType string
}

func (l *LaunchCommand) Execute() error {
	log.Println("Starting launch command execution...")

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("Error reading config: %v", err)
		return err
	}

	log.Printf("Config loaded successfully. Found %d runner configurations", len(cfg.TaskRunners))

	var runnerConfig config.TaskRunnerConfig
	found := false
	for _, r := range cfg.TaskRunners {
		log.Printf("Checking runner type: %s", r.RunnerType)
		if r.RunnerType == l.RunnerType {
			runnerConfig = r
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("unknown runner type: %s", l.RunnerType)
	}

	log.Printf("Found matching runner config. WorkDir: %s", runnerConfig.WorkDir)

	if err := os.Chdir(runnerConfig.WorkDir); err != nil {
		log.Printf("Failed to change directory to %s: %v", runnerConfig.WorkDir, err)
		return fmt.Errorf("failed to chdir into configured directory (%s): %w", runnerConfig.WorkDir, err)
	}

	log.Printf("Changed working directory successfully")

	defaultEnvs := []string{"LANG", "PATH", "TZ", "TERM"}
	allowedEnvs := append(defaultEnvs, runnerConfig.AllowedEnv...)
	env := filterEnv(allowedEnvs)

	log.Printf("Filtered environment variables. Count: %d", len(env))

	if token := os.Getenv("N8N_RUNNERS_AUTH_TOKEN"); token != "" {
		log.Printf("Found auth token, attempting to fetch grant token")
		n8nUri := os.Getenv("N8N_RUNNERS_N8N_URI")
		if n8nUri == "" {
			return fmt.Errorf("N8N_RUNNERS_N8N_URI is required when N8N_RUNNERS_AUTH_TOKEN is set")
		}

		grantToken, err := auth.FetchGrantToken(n8nUri, token)
		if err != nil {
			log.Printf("Failed to fetch grant token: %v", err)
			return fmt.Errorf("failed to fetch grant token: %w", err)
		}

		log.Printf("Successfully obtained grant token")
		env = append(env, fmt.Sprintf("N8N_RUNNERS_GRANT_TOKEN=%s", grantToken))
	}

	log.Printf("Launching runner with command: %s", runnerConfig.Command)
	log.Printf("Arguments: %v", runnerConfig.Args)

	cmd := exec.Command(runnerConfig.Command, runnerConfig.Args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting process...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Process failed to start: %v", err)
		return err
	}

	log.Printf("Process started successfully")
	return nil
}

func filterEnv(allowed []string) []string {
	var filtered []string

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		for _, allowedKey := range allowed {
			if key == allowedKey {
				filtered = append(filtered, env)
				break
			}
		}
	}

	return filtered
}
