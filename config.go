package main

import (
	"os"
	"log"
	"github.com/BurntSushi/toml"
)

// Info from config file
type Config struct {
	Email     string
	ApiKey    string
	Domain 	  string
	DNSRecord string
}

// Reads info from config file
func ReadConfig(configFile string) Config {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("Config file is missing: ", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}