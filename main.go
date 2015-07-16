// Linda - little princess, programmed to serve in mens-only online conference rooms.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/kpashka/linda/config"
	"github.com/kpashka/linda/kernel"
)

func main() {
	// Parse command-line flags
	location := flag.String("c", "config.toml", "Config file location (URL or filesystem)")
	flag.Parse()

	// Load configuration
	config := config.New()
	err := config.Load(*location)
	if err != nil {
		log.Printf("Can't load configuration from location: %s", *location)
		log.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	// Create bot instance and run
	d := kernel.NewLinda(config)
	d.Start()
}
