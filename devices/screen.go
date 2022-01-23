package devices

import (
	"image"
	"miveil/core"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ssd1306"

	log "github.com/sirupsen/logrus"
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

	log.Debug(dev)

	s := Screen{dev: dev}

	return &s, nil
}

func (s *Screen) Bounds() image.Rectangle {
	return s.dev.Bounds()
}

func (s *Screen) PlayAnimation(animation *core.Animation) {
	if !s.isPlaying {
		s.isPlaying = true
		currentFrame := 0

		for i := 0; i < len(animation.Sequence); i++ {
			c := time.After(animation.FrameDuration[currentFrame])
			img := animation.Frames[currentFrame]
			//log.Debug(img.Bounds(), animation.Sequence[i])

			currentImage := core.CreateFrame(s.dev.Bounds().Dx(), s.dev.Bounds().Dy(), img, *animation.Sequence[i])

			s.dev.Draw(currentImage.Bounds(), currentImage, image.Point{})
			currentFrame = (currentFrame + 1) % len(animation.Frames)
			<-c
		}
		s.dev.Halt()
		s.isPlaying = false
	}
}

func (s *Screen) Stop() {
	s.dev.Halt()
}
