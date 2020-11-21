package data

import (
	"time"
)

// HeatingSchedule contains schedule at which hour which temperature is desired
type HeatingSchedule [24]Temperature

type Setting struct {
	Hour        int
	Temperature Temperature
}

func NewHeatingSchedule(settings ...*Setting) *HeatingSchedule {
	s := &HeatingSchedule{}

	// fill with invalid temperatures (defaults)
	for i := 0; i < 24; i++ {
		s[i] = InvalidTemperature()
	}

	// then apply settings
	for _, setting := range settings {
		s[setting.Hour] = setting.Temperature
	}

	return s
}

// GetTemperature get the temperature for the given time
func (s *HeatingSchedule) GetTemperature(t time.Time) (temperature Temperature) {
	temperature = InvalidTemperature()
	hour := t.Hour()

	// first we try to find the first setting that applies
	// starting from the given hour until hour 0
	for i := hour; 0 <= i; i-- {
		temp := s[i]
		if temp.IsValid() {
			return temp
		}
	}
	// if we couldn't find a setting, we have to take the setting
	// that is last in the schedule
	for i := len(s) - 1; hour < i; i-- {
		temp := s[i]
		if temp.IsValid() {
			return temp
		}
	}

	return temperature
}
