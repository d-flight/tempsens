package application

import (
	"tempsens/data"
)

// Buffers for temperature, when heating is canceled/triggered
// In Celsius / 100
const (
	upper_temp_buffer = 25
	lower_temp_buffer = 15
)

var deltaCap = &data.DeltaCap{Humidity: 200, Temperature: 50}

// Controller ...
type Controller struct {
	view  *View
	state *data.State
}

// NewController ...
func NewController(view *View, state *data.State) *Controller {
	return &Controller{view, state}
}

// SetDesiredTemperature ...
func (c *Controller) SetDesiredTemperature(t data.Temperature) {
	if t == c.state.GetDesiredTemperature() {
		return
	}

	c.state.SetDesiredTemperature(t)

	// trigger heating state update
	c.updateHeatingState()

	// update view
	c.updateView()
}

func (c *Controller) SetLatestReading(reading *data.Reading) {
	if reading == nil {
		return
	}

	if c.state.GetLatestReading() == nil {
		// first reading, booting is done
		c.view.finishBooting()
	} else if c.state.GetLatestReading().Equals(reading, deltaCap) {
		// if the reading didn't change we exit here
		return
	}

	// update state
	c.state.SetLatestReading(reading)

	// update heating state
	c.updateHeatingState()

	// update view
	c.updateView()
}

func (c *Controller) updateHeatingState() {
	newHeatingState := c.state.GetHeatingState()

	if data.HEATING_STATE_OFF != newHeatingState {
		lastReading := c.state.GetLatestReading()
		desiredTemperature := c.state.GetDesiredTemperature()

		if nil == lastReading || lastReading.Temperature == desiredTemperature {
			return
		}

		// save old heating state and calculate new one
		actualTemperature := lastReading.Temperature
		newHeatingState = c.state.GetHeatingState()

		// if we reached a tempearture that is higher than the desiredTemperature, switch to idle
		if upper_temp_buffer < actualTemperature-desiredTemperature {
			newHeatingState = data.HEATING_STATE_IDLE
		}

		// if we reached a temperature that is lower than the desiredTemperature, switch to heating
		if lower_temp_buffer < desiredTemperature-actualTemperature {
			newHeatingState = data.HEATING_STATE_ON
		}
	}

	// update state
	c.state.SetHeatingState(newHeatingState)
}

func (c *Controller) SetHeatingState(newState data.HeatingState) {
	// update state
	c.state.SetHeatingState(newState)

	// apply
	c.updateHeatingState()

	// update view
	c.updateView()
}

func (c *Controller) updateView() {
	c.view.ViewState(*c.state)
}
