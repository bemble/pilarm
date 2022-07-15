package devices

import (
	"pilarm/core"
	"pilarm/hardware"
)

type Audio struct {
	dev            hardware.PwnAudio
	player         core.AlsaPlayer
	alarmSoundFile string
}

func NewAudio(pinNumber int, alarmSoundFile string) *Audio {
	dev := hardware.NewPwnAudio(pinNumber)
	player := core.NewAlsaPlayer()
	return &Audio{dev: dev, player: player, alarmSoundFile: alarmSoundFile}
}

func (a *Audio) playAlarm(duration int) error {
	defer a.dev.TurnOff()

	e := a.dev.TurnOn()
	if e == nil {
		e = a.player.Play(a.alarmSoundFile, duration)
	}
	return e
}

func (a *Audio) PlayAlarm() error {
	return a.playAlarm(0)
}

func (a *Audio) StopAlarm() error {
	a.player.Stop()
	return a.dev.TurnOff()
}

func (a *Audio) PlayAlarmFor(duration int) error {
	return a.playAlarm(duration)
}

func (a *Audio) Stop() error {
	return a.dev.ShutDown()
}
