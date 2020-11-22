package application

import (
	"tempsens/data"
	"time"
)

const (
	// schedule triggers every 10 seconds
	schedule_update_interval = 10 * time.Second
)

type Schedule struct {
	heatingSchedule *data.HeatingSchedule
	halt            chan bool

	OnTick func(data.Temperature)
}

func NewSchedule(heatingSchedule data.HeatingSchedule, controller *Controller) *Schedule {
	return &Schedule{
		heatingSchedule: &heatingSchedule,
		halt:            make(chan bool),

		OnTick: func(d data.Temperature) {
			controller.SetDesiredTemperature(d, false)
		},
	}
}

func (s *Schedule) Start() {
	for {
		startTime := time.Now()

		s.Trigger(startTime)

		select {
		case <-time.After(schedule_update_interval - time.Since(startTime)):
		case <-s.halt:
			return
		}
	}
}

func (s *Schedule) Trigger(time time.Time) {
	currentTemperature := s.heatingSchedule.GetTemperature(time)

	s.OnTick(currentTemperature)
}

func (s *Schedule) Stop() {
	s.halt <- true
}
