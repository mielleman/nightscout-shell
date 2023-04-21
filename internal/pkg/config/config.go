package config

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	NightscoutUrl   string `json:"nightscout_url"`
	NightscoutToken string `json:"nightscout_token"`
	CacheFile       string `json:"cache_file"`
	ServiceInterval int    `json:"service_interval"`
}

func new() *Configuration {
	// Set some defaults
	return &Configuration{
		CacheFile:       os.Getenv("HOME") + "/.cache/nightscout-shell/prompt.dat",
		ServiceInterval: 5,
	}
}

func ParseConfigFile(filename string) *Configuration {
	// Open the configuration file
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		log.Panicf("Failed to load the given configuration file: %s", filename)
		os.Exit(1)
	}

	// Defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// Load the config into the struct
	conf := new()
	err = json.Unmarshal(byteValue, conf)
	if err != nil {
		log.Error(err)
		log.Panicf("Failed to load the given configuration file: %s", filename)
		os.Exit(1)
	}

	// Make sure the URL and the Token are set
	if conf.NightscoutUrl == "" {
		log.Panicf("Missing Nightscout URL, this value must be set in your configuration file!")
		os.Exit(1)
	}
	if conf.NightscoutToken == "" {
		log.Panicf("Missing Nightscout Token, this value must be set in your configuration file!")
		os.Exit(1)
	}

	return conf
}
