package devices

import (
	"time"

	"miveil/hardware"
)

type Led struct {
	dev hardware.Led
}

func NewLed(portNumber int) Led {
	dev := hardware.NewLed(portNumber)
	return Led{dev: dev}
}

func (l *Led) TurnOn() error {
	return l.dev.TurnOn()
}

func (l *Led) TurnOff() error {
	return l.dev.TurnOff()
}

func (l *Led) TurnOnFor(duration time.Duration) {
	defer l.TurnOff()
	l.TurnOn()
	time.Sleep(duration)
}

func (l *Led) Toggle() error {
	if l.dev.IsOn {
		return l.dev.TurnOff()
	} else {
		return l.dev.TurnOn()
	}
}
