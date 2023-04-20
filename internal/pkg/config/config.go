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
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Error(err)
		log.Panic("Failed to load the given configuration file: " + filename)
		os.Exit(1)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	// Load the config into the struct
	conf := new()
	json.Unmarshal(byteValue, conf)

	return conf
}
