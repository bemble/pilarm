package devices

import (
	"image"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"
)

type Screen struct {
	dev       *ssd1306.Dev
	isPlaying bool
}

func NewScreen() (*Screen, error) {
	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	// Open a handle to a ssd1306 connected on the I²C bus:
	dev, err := ssd1306.NewI2C(bus, &ssd1306.Opts{W: 128, H: 64})
	if err != nil {
		return nil, err
	}

	dev.SetContrast(0x00)
	dev.Halt()

	s := Screen{dev: dev}

	return &s, nil
}

func (s *Screen) PlayAnimation(animation []*image.Gray) {
	if !s.isPlaying {
		s.isPlaying = true
		for i := 0; i < len(animation); i++ {
			c := time.After(10 * time.Millisecond)
			img := animation[i]
			s.dev.Draw(img.Bounds(), img, image.Point{})
			<-c
		}
		s.dev.Halt()
		s.isPlaying = false
	}
}
