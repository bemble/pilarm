package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"periph.io/x/host/v3"

	"miveil/core"
)

// Add process env + config

func init() {
	log.SetLevel(log.FatalLevel)

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	defer rpio.Close()
	miveil := core.NewMiveil()

	// does not work: defer miveil.Stop()
	miveil.Start()
}
