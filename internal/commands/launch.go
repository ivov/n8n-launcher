package commands

import (
	"fmt"
	"log"
	"n8n-launcher/internal/auth"
	"n8n-launcher/internal/config"
	"os"
	"os/exec"
	"strings"
)

type LaunchCommand struct {
	RunnerType string
}

func (l *LaunchCommand) Execute() error {
	log.Println("Starting `launch` command execution...")

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("Error reading config: %v", err)
		return err
	}

	cfgNum := len(cfg.TaskRunners)

	if cfgNum == 0 {
		return fmt.Errorf("found no task runner configs in launcher config")
	} else if cfgNum == 1 {
		log.Print("Config file loaded. Found a single runner config")
	} else {
		log.Printf("Config file loaded. Found %d runner configurations", cfgNum)
	}

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

	if err := os.Chdir(runnerConfig.WorkDir); err != nil {
		log.Printf("Failed to change dir to %s: %v", runnerConfig.WorkDir, err)
		return fmt.Errorf("failed to chdir into configured dir (%s): %w", runnerConfig.WorkDir, err)
	}

	log.Printf("Changed working directory")

	defaultEnvs := []string{"LANG", "PATH", "TZ", "TERM"}
	allowedEnvs := append(defaultEnvs, runnerConfig.AllowedEnv...)
	env := filterEnvToAllowedOnly(allowedEnvs)

	log.Printf("Filtered environment variables")

	if token := os.Getenv("N8N_RUNNERS_AUTH_TOKEN"); token != "" {
		n8nUri := os.Getenv("N8N_RUNNERS_N8N_URI")
		if n8nUri == "" {
			return fmt.Errorf("N8N_RUNNERS_N8N_URI is required when N8N_RUNNERS_AUTH_TOKEN is set")
		}

		log.Printf("Found auth token and n8n URI, attempting to fetch grant token")
		grantToken, err := auth.FetchGrantToken(n8nUri, token)
		if err != nil {
			log.Printf("Failed to fetch grant token: %v", err)
			return fmt.Errorf("failed to fetch grant token: %w", err)
		}

		log.Printf("Obtained grant token from n8n main instance")
		env = append(env, fmt.Sprintf("N8N_RUNNERS_GRANT_TOKEN=%s", grantToken))
	}

	log.Printf("Launching runner with command: %s", runnerConfig.Command)
	log.Printf("Arguments: %v", runnerConfig.Args)

	cmd := exec.Command(runnerConfig.Command, runnerConfig.Args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Starting task runner...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Task runner failed to start: %v", err)
		return err
	}

	log.Printf("Successfully started task runner")
	return nil
}

func filterEnvToAllowedOnly(allowed []string) []string {
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
