package hardware

import (
	"strconv"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

type PwnAudio struct {
	powerPin gpio.PinIO
	IsOn     bool
}

func NewPwnAudio(portNumber int) PwnAudio {
	d := PwnAudio{powerPin: gpioreg.ByName(strconv.Itoa(portNumber))}
	d.turn(false, true)
	return d
}

func (d *PwnAudio) turn(on bool, force bool) error {
	if on == d.IsOn && !force {
		return nil
	}

	newValue := gpio.High
	if !on {
		newValue = gpio.Low
	}
	err := d.powerPin.Out(newValue)
	if err == nil {
		d.IsOn = on
	}
	return err
}

func (d *PwnAudio) TurnOn() error {
	return d.turn(true, false)
}

func (d *PwnAudio) TurnOff() error {
	return d.turn(false, false)
}

func (d *PwnAudio) ShutDown() error {
	return d.turn(false, true)
}
