package main

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"periph.io/x/host/v3"

	"miveil/miveil"
)

// Add process env + config

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	miveil := miveil.NewMiveil()

	// does not work: defer miveil.Stop()
	miveil.Start()
}
