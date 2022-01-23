package miveil

import (
	"miveil/core"
	"miveil/devices"
	"time"

	log "github.com/sirupsen/logrus"
)

type Miveil struct {
	canWakeUpLed       devices.Led
	stayInBedLed       devices.Led
	sonar              devices.Sonar
	screen             *devices.Screen
	currentLed         *devices.Led
	scheduler          core.Scheduler
	wasOn              bool
	sonarOnSince       int64
	canWakeUpAnimation *core.Animation
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func NewMiveil() Miveil {
	screen, _ := devices.NewScreen()
	animation, _ := core.Gif2Animation(screen.Bounds().Dx(), screen.Bounds().Dy(), "./ressources/pikachu.gif", 1500*time.Millisecond)

	miveil := Miveil{
		canWakeUpLed:       devices.NewLed(27), //rpi.P1_13
		stayInBedLed:       devices.NewLed(17), //rpi.P1_11
		sonar:              devices.NewSonar(6, 13),
		screen:             screen,
		currentLed:         nil,
		scheduler:          core.NewScheduler(),
		wasOn:              false,
		sonarOnSince:       0,
		canWakeUpAnimation: animation,
	}

	miveil.sonar.AddCallback(miveil.sonarCallback)

	return miveil
}

func (m *Miveil) sonarCallback(d float32) {
	if d <= 0.6 {
		if m.sonarOnSince == 0 {
			m.sonarOnSince = makeTimestamp()
			return
		}
		if makeTimestamp()-m.sonarOnSince > 150 {
			log.WithFields(log.Fields{"component": "hardware", "category": "sonar"}).Debug("triggered while ", makeTimestamp()-m.sonarOnSince, "ms")
			m.wasOn = true
			led := m.stayInBedLed
			if m.scheduler.CanWakeUp {
				led = m.canWakeUpLed
				go m.screen.PlayAnimation(m.canWakeUpAnimation)
			}
			if m.currentLed != nil && led != *m.currentLed {
				m.currentLed.TurnOff()
			}
			m.currentLed = &led
			m.currentLed.TurnOn()
		}
	} else {
		m.sonarOnSince = 0
		if m.wasOn {
			log.WithFields(log.Fields{"component": "hardware", "category": "sonar"}).Info("Triggered")
			m.wasOn = false
			if m.currentLed != nil {
				go func() {
					ledDuration := 1 * time.Second
					if m.scheduler.CanWakeUp {
						ledDuration = 2 * time.Second
					}

					m.currentLed.TurnOnFor(ledDuration)
					m.currentLed = nil
				}()
			}
		}
	}
}

func (m *Miveil) Start() {
	go m.scheduler.Start()
	go m.sonar.Start()
}

func (m *Miveil) Stop() {
	m.canWakeUpLed.Stop()
	m.stayInBedLed.Stop()
	m.sonar.Stop()
	m.screen.Stop()
}
