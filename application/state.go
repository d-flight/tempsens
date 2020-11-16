package application

import (
	"tempsens/data"
	"tempsens/sensor"
)

const (
	HEATING_STATE_OFF  = 0
	HEATING_STATE_ON   = 1
	HEATING_STATE_IDLE = 2
)

type State struct {
	heatingState       int
	desiredTemperature data.Temperature
	latestReading      *sensor.Reading
	isUserControlled   bool
}

func NewState() *State {
	return &State{
		heatingState:       HEATING_STATE_IDLE,
		desiredTemperature: data.InvalidTemperature(),
		latestReading:      nil,
		isUserControlled:   false,
	}
}

func (s *State) GetHeatingState() int { return s.heatingState }

func (s *State) GetLatestReading() *sensor.Reading { return s.latestReading }

func (s *State) SetLatestReading(reading *sensor.Reading) { s.latestReading = reading }

func (s *State) GetDesiredTemperature() data.Temperature { return s.desiredTemperature }

func (s *State) IsUserControlled() bool { return s.isUserControlled }

func (s *State) SetDesiredTemperature(tempearture data.Temperature, userControlled bool) {
	s.desiredTemperature = tempearture
	s.isUserControlled = userControlled
}
