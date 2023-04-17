package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mielleman/nightscout-shell/internal/app/prompt"
	"github.com/mielleman/nightscout-shell/internal/app/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Default values shared between sub commands
	defaultConfigFilename := os.Getenv("HOME") + "/.config/nightscout-shell/config.json"

	// The service runs and updates the cache file
	serviceCmd := flag.NewFlagSet("service", flag.ExitOnError)
	serviceConfig := serviceCmd.String("config", defaultConfigFilename, "Location of the configuration file")

	// The prompt is a single-shot and just reads the cache file
	promptCmd := flag.NewFlagSet("prompt", flag.ExitOnError)
	promptConfig := promptCmd.String("config", defaultConfigFilename, "Location of the configuration file")

	// Set the logger
	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true,
	})
	log.Info("Main application start")

	// No subcommand, then we run as prompt
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "prompt")
	}

	// Act per sub command
	switch os.Args[1] {
	case "service":
		serviceCmd.Parse(os.Args[2:])
		s := service.New(*serviceConfig)
		s.Start()

	case "prompt":
		promptCmd.Parse(os.Args[2:])
		p := prompt.New(*promptConfig)
		p.Main()

	default:
		fmt.Println("ERROR: Expected either 'service' or 'prompt' subcommands")
		os.Exit(1)
	}
}
