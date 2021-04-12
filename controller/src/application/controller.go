package application

import (
	"tempsens/adapter"
	"tempsens/data"
)

// TODO: Move to slave
// Buffers for temperature, when heating is canceled/triggered
// In Celsius / 100
const (
	upper_temp_buffer = 25
	lower_temp_buffer = 15
)

// TODO: Move to slave
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
	if stateDesired := c.state.GetDesiredTemperature(); stateDesired.IsValid() && stateDesired != desired {
		c.mqttAdapter.SetDesiredTemperature(stateDesired)
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

// TODO: Move to slave
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
