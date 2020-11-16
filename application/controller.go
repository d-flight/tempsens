package application

import (
	"tempsens/data"
	"tempsens/sensor"
	"time"
)

// Buffers for temperature, when heating is canceled/triggered
// In Celsius / 100
const (
	UPPER_TEMP_BUFFER = 30
	LOWER_TEMP_BUFFER = 50
)

const (
	STATE_OFF     = 0
	STATE_HEATING = 1
	STATE_IDLE    = 2
)

// Controller ...
type Controller struct {
	LatestReading      *sensor.Reading
	DesiredTemperature *data.Temperature
	State              int
	Schedule           *data.HeatingSchedule
}

// NewController ...
func NewController(schedule *data.HeatingSchedule, state int) *Controller {
	return &Controller{
		Schedule: schedule,
		State:    state,
	}
}

// SetDesiredTemperature ...
func (c *Controller) SetDesiredTemperature(t *data.Temperature) {
	c.DesiredTemperature = t
	c.UpdateState()
}

// GetDesiredTemperature ...
func (c *Controller) GetDesiredTemperature() data.Temperature {
	if c.DesiredTemperature != nil {
		return *c.DesiredTemperature
	}

	return c.Schedule.GetTemperature(time.Now())
}

// UpdateState ...
func (c *Controller) UpdateState() int {
	// we can't update state unless we have a reading
	// also we won't update the state, if it is set to STATE_OFF
	if STATE_OFF == c.State || c.LatestReading == nil {
		return c.State
	}

	target := c.GetDesiredTemperature()

	// if we reached a tempearture that is higher than the target, switch to idle
	if UPPER_TEMP_BUFFER < c.LatestReading.Temperature-target {
		c.State = STATE_IDLE
	}

	// if we reached a temperature that is lower than the target, switch to heating
	if LOWER_TEMP_BUFFER < target-c.LatestReading.Temperature {
		c.State = STATE_HEATING
	}

	return c.State
}
