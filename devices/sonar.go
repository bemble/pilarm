package devices

import (
	"time"

	log "github.com/sirupsen/logrus"

	"miveil/hardware"
)

type Sonar struct {
	h                  *hardware.HCSR04
	onMeasureCallbacks []func(float32)
}

func NewSonar(trigger int, echo int) Sonar {
	h := hardware.NewHCSR04(trigger, echo)
	return Sonar{h: h, onMeasureCallbacks: []func(float32){}}
}

func (r *Sonar) AddCallback(f func(float32)) {
	r.onMeasureCallbacks = append(r.onMeasureCallbacks, f)
}

func (r *Sonar) Start() {
	go func() {
		if err := r.h.StartDistanceMonitor(); err != nil {
			log.Fatal("impossible to start distance monitor")
		} else {
			defer r.h.StopDistanceMonitor()
		}
		for {
			distance := r.h.GetDistance()
			if distance > 0 {
				for i := 0; i < len(r.onMeasureCallbacks); i++ {
					r.onMeasureCallbacks[i](distance)
				}
				time.Sleep(hardware.HCSR04MonitorUpdate)
			}
		}
	}()
}
