package pilarm

import (
	"pilarm/core"
	"pilarm/core/config"
	"pilarm/devices"
	"time"

	log "github.com/sirupsen/logrus"
)

type Pilarm struct {
	canWakeUpLed       *devices.Led
	stayInBedLed       *devices.Led
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
	Config := config.Get()
	log.Println(Config.Leds)
	pilarm := Pilarm{
		canWakeUpLed:       nil,
		stayInBedLed:       nil,
		sonar:              devices.NewSonar(Config.Sonar.TriggerPin, Config.Sonar.EchoPin),
		screen:             nil,
		currentLed:         nil,
		scheduler:          core.NewScheduler(),
		wasOn:              false,
		sonarOnSince:       0,
		canWakeUpAnimation: nil,
	}

	if Config.Leds.ArePresent && Config.Leds.StayInBedPin > 0 && Config.Leds.CanWakeUpPin > 0 {
		pilarm.canWakeUpLed = devices.NewLed(Config.Leds.CanWakeUpPin)
		pilarm.stayInBedLed = devices.NewLed(Config.Leds.StayInBedPin)

	}

	if Config.Screen.IsPresent {
		screen, errorScreen := devices.NewScreen()
		if errorScreen != nil {
			log.WithError(errorScreen).Warn("No screen found")
		}

		if screen != nil {
			pilarm.screen = screen
			Config := config.Get()

			if Config.Screen.CanWakeUpAnimationFile != "" && Config.Screen.CanWakeUpAnimationDuration > 0 {
				errorAnimation := error(nil)
				animation, errorAnimation := core.Gif2Animation(screen.Bounds().Dx(), screen.Bounds().Dy(), config.GetRessourcePath(Config.Screen.CanWakeUpAnimationFile), time.Duration(Config.Screen.CanWakeUpAnimationDuration)*time.Second)
				if errorAnimation != nil {
					log.WithError(errorAnimation).Warn("Animation not found")
				} else {
					pilarm.canWakeUpAnimation = animation
				}
			}
		}
	}

	pilarm.sonar.AddCallback(pilarm.sonarCallback)

	return &pilarm, nil
}

func (m *Pilarm) sonarCallback(d float32) {
	Config := config.Get()
	if d <= Config.Sonar.MaxDistance && d >= Config.Sonar.MinDistance {
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
			}
			if m.currentLed != nil && led != m.currentLed {
				m.currentLed.TurnOff()
			}
			if m.screen != nil {
				go func() {
					if m.scheduler.CanWakeUp && m.canWakeUpAnimation != nil {
						m.screen.PlayAnimation(m.canWakeUpAnimation)
					}
					screenDuration := time.Duration(Config.Screen.StayInBedDisplayTimeDuration) * time.Second
					if m.scheduler.CanWakeUp {
						screenDuration = time.Duration(Config.Screen.CanWakeUpDisplayTimeDuration) * time.Second
					}
					if screenDuration > 0 {
						m.screen.DisplayTimeFor(screenDuration)
					}
				}()
			}
			m.currentLed = led
			if m.currentLed != nil {
				m.currentLed.TurnOn()
			}
		}
	} else {
		m.sonarOnSince = 0
		if m.wasOn {
			m.wasOn = false
			if m.currentLed != nil {
				go func() {
					ledDuration := time.Duration(Config.Leds.StayInBedDisplayDuration) * time.Second
					if m.scheduler.CanWakeUp {
						ledDuration = time.Duration(Config.Leds.CanWakeUpDisplayDuration) * time.Second
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
	if m.canWakeUpLed != nil {
		m.canWakeUpLed.Stop()
	}
	if m.stayInBedLed != nil {
		m.stayInBedLed.Stop()
	}
	m.sonar.Stop()
	if m.screen != nil {
		m.screen.Stop()
	}
}
