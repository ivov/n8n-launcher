package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"task-runner-launcher/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument. Expected 'launch' or 'kill' subcommand")
		os.Exit(1)
	}

	var cmd commands.Command

	switch os.Args[1] {
	case "launch":
		launchCmd := flag.NewFlagSet("launch", flag.ExitOnError)
		runnerType := launchCmd.String("type", "", "Runner type to launch")
		launchCmd.Parse(os.Args[2:])

		if *runnerType == "" {
			launchCmd.PrintDefaults()
			os.Exit(1)
		}

		cmd = &commands.LaunchCommand{RunnerType: *runnerType}

	case "kill":
		killCmd := flag.NewFlagSet("kill", flag.ExitOnError)
		runnerType := killCmd.String("type", "", "Runner type to kill")
		pid := killCmd.Int("pid", 0, "Process ID to kill")
		killCmd.Parse(os.Args[2:])

		if *runnerType == "" || *pid == 0 {
			killCmd.PrintDefaults()
			os.Exit(1)
		}

		cmd = &commands.KillCommand{
			RunnerType: *runnerType,
			PID:        *pid,
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Expected 'launch' or 'kill' subcommand")
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
