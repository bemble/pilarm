package core

import (
	"time"
)

type Scheduler struct {
	CanWakeUp bool
}

func canWakeUp() bool {
	now := time.Now()
	maxSleeping := time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, now.Location())
	minSleeping := time.Date(now.Year(), now.Month(), now.Day(), 19, 30, 0, 0, now.Location())
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
