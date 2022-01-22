package devices

import (
	"log"
	"time"

	"github.com/raspberrypi-go-drivers/hcsr04"
)

type Sonar struct {
	h                  *hcsr04.HCSR04
	onMeasureCallbacks []func(float32)
}

func NewSonar(trigger uint8, echo uint8) Sonar {
	h := hcsr04.NewHCSR04(trigger, echo)
	return Sonar{h: h, onMeasureCallbacks: []func(float32){}}
}

func (r *Sonar) AddCallback(f func(float32)) {
	r.onMeasureCallbacks = append(r.onMeasureCallbacks, f)
}

func (r *Sonar) Start() {
	go func() {
		if err := r.h.StartDistanceMonitor(); err != nil {
			log.Panic("impossible to start distance monitor")
		} else {
			defer r.h.StopDistanceMonitor()
		}
		for {
			distance := r.h.GetDistance()
			if distance > 0 {
				for i := 0; i < len(r.onMeasureCallbacks); i++ {
					r.onMeasureCallbacks[i](distance)
				}
				time.Sleep(hcsr04.MonitorUpdate)
			}
		}
	}()
}
