package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// Parse command-line flags
	cfgFilePath := flag.String("c", "config.json", "Config file location")
	flag.Parse()

	// Read config file
	bytes, err := ioutil.ReadFile(*cfgFilePath)
	if err != nil {
		log.WithFields(log.Fields{"file": *cfgFilePath}).Fatalf("Couldn't open config file: %s", err.Error())
	}

	// Create config instance
	config := Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.WithFields(log.Fields{"file": *cfgFilePath}).Fatalf("Error in config file: %s", err.Error())
	}

	// Create bot instance and run
	d := NewDumbSlut(&config)
	d.Start()
}
