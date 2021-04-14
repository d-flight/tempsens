package application

import (
	"tempsens/adapter"
	"tempsens/data"
)

var deltaCap = &data.DeltaCap{Humidity: 200, Temperature: 50}

// Controller ...
type Controller struct {
	view        *View
	state       *data.State
	mqttAdapter *adapter.Mqtt
}

// NewController ...
func NewController(view *View, state *data.State, mqttAdapter *adapter.Mqtt) *Controller {
	return &Controller{view, state, mqttAdapter}
}

// SetDesiredTemperature ...
func (c *Controller) SetDesiredTemperature(t data.Temperature) {
	if t == c.state.GetDesiredTemperature() {
		return
	}

	c.state.SetDesiredTemperature(t)

	// propagate
	c.mqttAdapter.SetDesiredTemperature(t)

	// update view
	c.updateView()
}

func (c *Controller) HandleNewReport(reading *data.Reading, desired data.Temperature, actualHeatingState data.HeatingState) {
	// update latest reading
	c.SetLatestReading(reading)

	// report mismatches state, re-propagate
	stateHeatingState := c.state.GetHeatingState()
	if actualHeatingState != stateHeatingState {
		if data.HEATING_STATE_OFF == stateHeatingState { // should be off
			c.mqttAdapter.ToggleActive(false)
		} else { // state needs an update
			c.state.SetHeatingState(actualHeatingState)
		}
	} else {
		temperatureDesired := c.state.GetDesiredTemperature()
		if temperatureDesired.IsValid() && temperatureDesired != desired {
			c.mqttAdapter.SetDesiredTemperature(temperatureDesired)
		}
	}

	c.updateView()
}

func (c *Controller) SetLatestReading(reading *data.Reading) {
	if reading == nil {
		return
	}

	latestReading := c.state.GetLatestReading()

	if latestReading != nil && latestReading.Equals(reading, deltaCap) {
		// if the reading didn't change we exit here
		return
	}

	// update state
	c.state.SetLatestReading(reading)

	// update view
	c.updateView()
}

func (c *Controller) SetHeatingState(newState data.HeatingState) {
	// update state
	c.state.SetHeatingState(newState)

	// propagate
	if data.HEATING_STATE_OFF == newState {
		c.mqttAdapter.ToggleActive(false)
	} else {
		c.mqttAdapter.ToggleActive(true)
	}

	// update view
	c.updateView()
}

func (c *Controller) updateView() {
	c.view.ViewState(*c.state)
}
