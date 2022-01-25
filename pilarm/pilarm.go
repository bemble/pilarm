package pilarm

import (
	"pilarm/core"
	"pilarm/core/config"
	"pilarm/devices"
	"time"

	log "github.com/sirupsen/logrus"
)

type Pilarm struct {
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

func NewPilarm() (*Pilarm, error) {
	pilarm := Pilarm{
		canWakeUpLed:       devices.NewLed(27), //rpi.P1_13
		stayInBedLed:       devices.NewLed(17), //rpi.P1_11
		sonar:              devices.NewSonar(6, 13),
		screen:             nil,
		currentLed:         nil,
		scheduler:          core.NewScheduler(),
		wasOn:              false,
		sonarOnSince:       0,
		canWakeUpAnimation: nil,
	}

	screen, errorScreen := devices.NewScreen()
	if errorScreen != nil {
		log.WithError(errorScreen).Warn("No screen found")
	}

	if screen != nil {
		errorAnimation := error(nil)
		animation, errorAnimation := core.Gif2Animation(screen.Bounds().Dx(), screen.Bounds().Dy(), config.GetRessourcePath("pikachu.gif"), 1500*time.Millisecond)
		if errorAnimation != nil {
			log.WithError(errorAnimation).Warn("Animation not found")
		} else {
			pilarm.screen = screen
			pilarm.canWakeUpAnimation = animation
		}
	}

	pilarm.sonar.AddCallback(pilarm.sonarCallback)

	return &pilarm, nil
}

func (m *Pilarm) sonarCallback(d float32) {
	if d <= 0.6 {
		if m.sonarOnSince == 0 {
			m.sonarOnSince = makeTimestamp()
			return
		}
		if makeTimestamp()-m.sonarOnSince > 150 {
			if !m.wasOn {
				log.WithFields(log.Fields{"component": "hardware", "category": "sonar"}).Debug("Sonar signal interuption for ", makeTimestamp()-m.sonarOnSince, "ms")
				canWakeUpStr := "can wake up"
				if !m.scheduler.CanWakeUp {
					canWakeUpStr = "shoud stay in bed"
				}
				log.WithFields(log.Fields{"component": "hardware", "category": "sonar"}).Info("Triggered, ", canWakeUpStr)
			}
			m.wasOn = true
			led := m.stayInBedLed
			if m.scheduler.CanWakeUp {
				led = m.canWakeUpLed
				if m.screen != nil {
					go m.screen.PlayAnimation(m.canWakeUpAnimation)
				}
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

func (m *Pilarm) Start() {
	go m.scheduler.Start()
	go m.sonar.Start()
}

func (m *Pilarm) Stop() {
	m.canWakeUpLed.Stop()
	m.stayInBedLed.Stop()
	m.sonar.Stop()
	if m.screen != nil {
		m.screen.Stop()
	}
}
