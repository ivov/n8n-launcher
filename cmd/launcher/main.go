package main

import (
	"flag"
	"os"

	"n8n-launcher/internal/commands"
	"n8n-launcher/internal/logs"
)

func main() {
	if len(os.Args) < 2 {
		logs.Logger.Fatal("Missing argument. Expected `launch` or `kill` subcommand")
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

	default:
		logs.Logger.Printf("Unknown command: %s\nExpected `launch` or `kill` subcommand", os.Args[1])
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		logs.Logger.Printf("Failed to execute command: %s", err)
	}
}
