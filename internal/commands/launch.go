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
	log.Println("Starting to execute `launch` command...")

	token := os.Getenv("N8N_RUNNERS_AUTH_TOKEN")
	n8nUri := os.Getenv("N8N_RUNNERS_N8N_URI")

	if token == "" || n8nUri == "" {
		return fmt.Errorf("both N8N_RUNNERS_AUTH_TOKEN and N8N_RUNNERS_N8N_URI are required")
	}

	// 1. read configuration

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("Error reading config: %v", err)
		return err
	}

	var runnerConfig config.TaskRunnerConfig
	found := false
	for _, r := range cfg.TaskRunners {
		if r.RunnerType == l.RunnerType {
			runnerConfig = r
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("config file does not contain requested runner type : %s", l.RunnerType)
	}

	cfgNum := len(cfg.TaskRunners)

	if cfgNum == 1 {
		log.Print("Loaded config file loaded with a single runner config")
	} else {
		log.Printf("Loaded config file with %d runner configs", cfgNum)
	}

	// 2. change into working directory

	if err := os.Chdir(runnerConfig.WorkDir); err != nil {
		log.Printf("Failed to change dir to %s: %v", runnerConfig.WorkDir, err)
		return fmt.Errorf("failed to chdir into configured dir (%s): %w", runnerConfig.WorkDir, err)
	}

	log.Printf("Changed into working directory: %s", runnerConfig.WorkDir)

	// 3. filter environment variables

	defaultEnvs := []string{"LANG", "PATH", "TZ", "TERM"}
	allowedEnvs := append(defaultEnvs, runnerConfig.AllowedEnv...)
	env := filterEnvToAllowedOnly(allowedEnvs)

	log.Printf("Filtered environment variables")

	// 4. authenticate with n8n main instance

	log.Printf("Attempting to authenticate with n8n main instance...")
	grantToken, err := auth.FetchGrantToken(n8nUri, token)
	if err != nil {
		log.Printf("Failed to fetch grant token from n8n main instance: %v", err)
		return fmt.Errorf("failed to fetch grant token from n8n main instance: %w", err)
	}

	env = append(env, fmt.Sprintf("N8N_RUNNERS_GRANT_TOKEN=%s", grantToken))

	log.Printf("Authenticated with n8n main instance")

	// 5. launch runner

	log.Printf("Launching runner...")
	log.Printf("Command: %s", runnerConfig.Command)
	log.Printf("Args: %v", runnerConfig.Args)
	log.Printf("Env vars: %v", envKeys(env))

	cmd := exec.Command(runnerConfig.Command, runnerConfig.Args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("Launching task runner...")
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to launch task runner: %v", err)
		return err
	}

	log.Printf("Successfully launched task runner")

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

func envKeys(envVars []string) []string {
	keys := make([]string, len(envVars))
	for i, env := range envVars {
		keys[i] = strings.SplitN(env, "=", 2)[0]
	}
	return keys
}
