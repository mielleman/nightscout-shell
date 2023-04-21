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

	// Set the logger format
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		DisableQuote:           true,
		DisableLevelTruncation: true,
	})

	// No sub-command, then we explain the usage
	if len(os.Args) < 2 {
		fmt.Printf("%s (service|prompt)\n", os.Args[0])
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  prompt    Get the value to be used in your prompt")
		fmt.Println("  service   Run the service to retrieve and update the latest value used by prompt")
		fmt.Println()
		os.Exit(1)
	}

	// Act per sub command
	switch os.Args[1] {
	case "service":
		log.Info("Service started")
		serviceCmd.Parse(os.Args[2:])
		s := service.New(*serviceConfig)
		s.Start()

	case "prompt":
		promptCmd.Parse(os.Args[2:])
		p := prompt.New(*promptConfig)
		p.Main()

	default:
		log.Error("Expected either the 'service' or 'prompt' sub-command.")
		os.Exit(1)
	}
}
