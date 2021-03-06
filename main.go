package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"periph.io/x/host/v3"

	"pilarm/core/config"
	"pilarm/pilarm"
)

func init() {
	Config := config.Get()

	logLevel := log.InfoLevel
	if Config.Debug {
		logLevel = log.DebugLevel
	}

	log.SetLevel(logLevel)
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	if _, err := host.Init(); err != nil {
		log.WithField("error", err).Fatal("Could not load drivers")
	}
}

func main() {
	log.Debug("Starting Pilarm in debug mode")
	pilarm, _ := pilarm.NewPilarm()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("")
		pilarm.Stop()
		log.WithField("category", "general").Info("Stopped")
		os.Exit(0)
	}()

	pilarm.Start()
	log.WithField("category", "general").Info("Started")

	for {
		time.Sleep(5 * time.Minute)
	}
}
