package core

import (
	"image"
	"miveil/devices"
	"time"
)

type Miveil struct {
	canWakeUpLed       devices.Led
	stayInBedLed       devices.Led
	sonar              devices.Sonar
	screen             *devices.Screen
	currentLed         *devices.Led
	scheduler          Scheduler
	wasOn              bool
	canWakeUpAnimation []*image.Gray
}

func NewMiveil() Miveil {
	animation, _ := Gif2Animation(128, 64, "./ressources/pikachu.gif")
	screen, _ := devices.NewScreen()

	miveil := Miveil{
		canWakeUpLed:       devices.NewLed(27), //rpi.P1_13
		stayInBedLed:       devices.NewLed(17), //rpi.P1_11
		sonar:              devices.NewSonar(6, 13),
		screen:             screen,
		currentLed:         nil,
		scheduler:          NewScheduler(),
		wasOn:              false,
		canWakeUpAnimation: animation,
	}

	miveil.sonar.AddCallback(miveil.sonarCallback)

	return miveil
}

func (m *Miveil) sonarCallback(d float32) {
	if d <= 0.6 {
		m.wasOn = true
		led := m.canWakeUpLed
		if !m.scheduler.CanWakeUp {
			led = m.stayInBedLed
		}
		if m.currentLed != nil && led != *m.currentLed {
			m.currentLed.TurnOff()
		}
		m.currentLed = &led
		m.currentLed.TurnOn()
	} else {
		if m.wasOn {
			m.wasOn = false
			ledDuration := 1 * time.Second
			if m.scheduler.CanWakeUp {
				ledDuration = 2 * time.Second
				go func() {
					m.screen.PlayAnimation(m.canWakeUpAnimation)
				}()
			}

			if m.currentLed != nil {
				go func() {
					m.currentLed.TurnOnFor(ledDuration)
					m.currentLed = nil
				}()
			}
		}
	}
}

func (m *Miveil) Start() {
	m.scheduler.Start()
	m.sonar.Start()

	for {
	}
}

func (m *Miveil) Stop() {
	m.canWakeUpLed.TurnOff()
	m.stayInBedLed.TurnOff()
}
