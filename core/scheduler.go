package core

import (
	"pilarm/core/config"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Scheduler struct {
	CanWakeUp bool
}

func canWakeUp() bool {
	now := time.Now()
	Config := config.Get()
	weekday := time.Now().Weekday()
	timesRv := reflect.ValueOf(Config.Times.WakeUp)
	fv := timesRv.FieldByName(weekday.String())

	wakeUpTime := strings.Split(fv.String(), ":")
	wakeUpHours, _ := strconv.Atoi(wakeUpTime[0])
	wakeUpMinutes, _ := strconv.Atoi(wakeUpTime[1])

	toBedTime := strings.Split(Config.Times.ToBed, ":")
	toBedHours, _ := strconv.Atoi(toBedTime[0])
	toBedMinutes, _ := strconv.Atoi(toBedTime[1])

	maxSleeping := time.Date(now.Year(), now.Month(), now.Day(), wakeUpHours, wakeUpMinutes, 0, 0, now.Location())
	minSleeping := time.Date(now.Year(), now.Month(), now.Day(), toBedHours, toBedMinutes, 0, 0, now.Location())
	return !(now.Before(maxSleeping) || now.After(minSleeping))
}

func NewScheduler() Scheduler {
	s := Scheduler{CanWakeUp: canWakeUp()}
	return s
}

func (s *Scheduler) Start() {
	for {
		s.CanWakeUp = canWakeUp()
		time.Sleep(25 * time.Second)
	}
}
