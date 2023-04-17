package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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
		CacheFile:       os.Getenv("HOME") + "/.config/nightscout-shell/cache.dat",
		ServiceInterval: 5,
	}
}

func ParseConfigFile(filename string) *Configuration {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR: Failed to load the given configuration file: " + filename)
		fmt.Println(err)
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
