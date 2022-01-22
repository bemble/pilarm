package devices

import (
	"strconv"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Led struct {
	port gpio.PinIO
	IsOn bool
}

func NewLed(portNumber int) Led {
	l := Led{port: gpioreg.ByName(strconv.Itoa(portNumber))}
	l.turn(false, true)
	return l
}

func (l *Led) turn(on bool, force bool) error {
	if on == l.IsOn && !force {
		return nil
	}

	newValue := gpio.High
	if !on {
		newValue = gpio.Low
	}
	err := l.port.Out(newValue)
	if err == nil {
		l.IsOn = on
	}
	return err
}

func (l *Led) TurnOnFor(duration time.Duration) {
	defer l.TurnOff()
	l.TurnOn()
	time.Sleep(duration)
}

func (l *Led) TurnOn() error {
	return l.turn(true, false)
}

func (l *Led) TurnOff() error {
	return l.turn(false, false)
}

func (l *Led) Toggle() error {
	if l.IsOn {
		return l.TurnOff()
	} else {
		return l.TurnOn()
	}
}
