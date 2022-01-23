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

	"miveil/miveil"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
	})

	if _, err := host.Init(); err != nil {
		log.WithField("error", err).Fatal("Could not load drivers")
	}
}

func main() {
	miveil := miveil.NewMiveil()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("")
		miveil.Stop()
		log.WithField("category", "general").Info("Stopped")
		os.Exit(0)
	}()

	miveil.Start()
	log.WithField("category", "general").Info("Started")

	for {
		time.Sleep(5 * time.Minute)
	}
}
