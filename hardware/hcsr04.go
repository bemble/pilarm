package hardware

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"

	log "github.com/sirupsen/logrus"
)

const (
	soundSpeed       float32       = 343.0
	measurementCycle time.Duration = 60 // 60ms between two measurements
	// MonitorUpdate is the time between each monitor update
	HCSR04MonitorUpdate time.Duration = 100 * time.Millisecond
)

// HCSR04 instance
type HCSR04 struct {
	triggerPinID           int
	echoPinID              int
	triggerPin             gpio.PinOut
	echoPin                gpio.PinIn
	mux                    sync.Mutex
	Measure                float32 // The last measure
	distanceMonitorControl chan int
	distanceMonitorStarted bool
}

// NewHCSR04 creates a new HCSR04 instance
func NewHCSR04(triggerPinID int, echoPinID int) *HCSR04 {
	hcsr04 := HCSR04{
		triggerPinID: triggerPinID,
		echoPinID:    echoPinID,
	}
	hcsr04.triggerPin = gpioreg.ByName(strconv.Itoa(hcsr04.triggerPinID))
	hcsr04.triggerPin.Out(gpio.Low)

	hcsr04.echoPin = gpioreg.ByName(strconv.Itoa(hcsr04.echoPinID))
	hcsr04.echoPin.In(gpio.PullNoChange, gpio.RisingEdge)
	return &hcsr04
}

// MeasureDistance measure the distance in front of sensor in meters
// and returns the measure
// MeasureDistance triggers a distance measurement by the sensor
//
// ! MeasureDistance is not design to work in a fast loop
// For this specific usage, use StartDistanceMonitor associated with GetDistance Instead
func (hcsr04 *HCSR04) MeasureDistance() (float32, error) {
	hcsr04.mux.Lock()
	defer hcsr04.mux.Unlock()
	pulseDuration, err := hcsr04.measurePulse()
	if err != nil {
		return 0, err
	}
	hcsr04.Measure = pulseToDistance(pulseDuration)
	return hcsr04.Measure, nil
}

func (hcsr04 *HCSR04) emitTrigger() {
	hcsr04.triggerPin.Out(gpio.High)
	time.Sleep(10 * time.Microsecond)
	hcsr04.triggerPin.Out(gpio.Low)
}

func (hcsr04 *HCSR04) measurePulse() (int64, error) {
	startChan := make(chan int64)
	stopChan := make(chan int64)
	startQuit := false
	stopQuit := false
	var startTime int64
	var stopTime int64
	go getPinStateChangeTime(hcsr04.echoPin, gpio.High, startChan, &startQuit)
	hcsr04.emitTrigger()
	if hcsr04.echoPin.Read() == gpio.High {
		return 0, errors.New("already receiving echo")
	}
	select {
	case t := <-startChan:
		startTime = t
	case <-time.After(measurementCycle * time.Millisecond):
		startQuit = true
		return 0, fmt.Errorf("echo not received after %d milliseconds", measurementCycle)
	}
	go getPinStateChangeTime(hcsr04.echoPin, gpio.Low, stopChan, &stopQuit)
	select {
	case t := <-stopChan:
		stopTime = t
	case <-time.After(measurementCycle * time.Millisecond):
		stopQuit = true
		return 0, fmt.Errorf("echo received for more than %d milliseconds", measurementCycle)
	}
	return stopTime - startTime, nil
}

func getPinStateChangeTime(pin gpio.PinIn, state gpio.Level, outChan chan int64, quit *bool) {
	for pin.Read() != state && !*quit {
	}
	if pin.Read() == state && !*quit {
		outChan <- time.Now().UnixNano()
	}
}

func pulseToDistance(pulseDuration int64) float32 {
	return float32(pulseDuration) / 1000000000.0 * soundSpeed / 2
}

// GetDistance returns the last distance measured
// Contrary to MeasureDistance, GetDistance does not trigger a distance measurement
func (hcsr04 *HCSR04) GetDistance() float32 {
	return hcsr04.Measure
}

// StartDistanceMonitor starts a process which will keep Measure updated
func (hcsr04 *HCSR04) StartDistanceMonitor() error {
	hcsr04.distanceMonitorControl = make(chan int)
	if hcsr04.distanceMonitorStarted {
		return errors.New("monitor already started")
	}
	go hcsr04.distanceMonitor()
	return nil
}

// StopDistanceMonitor stop the monitor process
func (hcsr04 *HCSR04) StopDistanceMonitor() {
	if hcsr04.distanceMonitorStarted {
		hcsr04.distanceMonitorControl <- 1
	}
}

func (hcsr04 *HCSR04) distanceMonitor() {
	for {
		select {
		case <-hcsr04.distanceMonitorControl:
			hcsr04.distanceMonitorStarted = false
			return
		default:
			if _, err := hcsr04.MeasureDistance(); err != nil {
				log.WithField("error", err).Trace("impossible to measure distance")
			}
		}
		time.Sleep(HCSR04MonitorUpdate)
	}
}
