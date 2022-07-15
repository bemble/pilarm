package hardware

import (
	"strconv"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type Led struct {
	pin  gpio.PinIO
	IsOn bool
}

func NewLed(portNumber int) Led {
	l := Led{pin: gpioreg.ByName(strconv.Itoa(portNumber))}
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
	err := l.pin.Out(newValue)
	if err == nil {
		l.IsOn = on
	}
	return err
}

func (l *Led) TurnOn() error {
	return l.turn(true, false)
}

func (l *Led) TurnOff() error {
	return l.turn(false, false)
}

func (l *Led) ShutDown() error {
	return l.turn(false, true)
}
