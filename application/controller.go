package application

import (
	"tempsens/data"
	"tempsens/sensor"

	"gobot.io/x/gobot/drivers/gpio"
)

// Buffers for temperature, when heating is canceled/triggered
// In Celsius / 100
const (
	UPPER_TEMP_BUFFER = 30
	LOWER_TEMP_BUFFER = 50
)

const (
	// according to the datasheet the DHT is accurate up to 0.5 celsius and 2% humidity
	deltaTemperature = 50
	deltaHumidity    = 200
)

var deltaCap = &sensor.DeltaCap{Humidity: deltaHumidity, Temperature: deltaTemperature}

// Controller ...
type Controller struct {
	view  *View
	state *State
	relay *gpio.RelayDriver
}

// NewController ...
func NewController(view *View, state *State, relay *gpio.RelayDriver) *Controller {
	return &Controller{view, state, relay}
}

// SetDesiredTemperature ...
func (c *Controller) SetDesiredTemperature(t data.Temperature, userControlled bool) {
	// if the heating is user controlled right now we ignore the temperature
	if !userControlled && c.state.IsUserControlled() {
		return
	}

	// update desired temperature in state
	c.state.SetDesiredTemperature(t, userControlled)

	// trigger heating state update
	c.updateHeatingState()

	// update view
	c.updateView()
}

func (c *Controller) SetLatestReading(reading *sensor.Reading) {
	if reading == nil {
		return
	}

	if c.state.latestReading == nil {
		// first reading, booting is done
		c.view.finishBooting()
	} else if c.state.latestReading.Equals(reading, deltaCap) {
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
	if HEATING_STATE_OFF == c.state.heatingState {
		c.relay.Off()
		return
	}

	lastReading := c.state.latestReading
	desiredTemperature := c.state.GetDesiredTemperature()

	if nil == lastReading || lastReading.Temperature == desiredTemperature {
		return
	}

	// save old heating state and calculate new one
	actualTemperature := lastReading.Temperature
	newHeatingState := c.state.GetHeatingState()

	// if we reached a tempearture that is higher than the desiredTemperature, switch to idle
	if UPPER_TEMP_BUFFER < actualTemperature-desiredTemperature {
		newHeatingState = HEATING_STATE_IDLE
	}

	// if we reached a temperature that is lower than the desiredTemperature, switch to heating
	if LOWER_TEMP_BUFFER < desiredTemperature-actualTemperature {
		newHeatingState = HEATING_STATE_ON
	}

	// update state
	c.state.heatingState = newHeatingState

	// update relay
	if HEATING_STATE_ON == newHeatingState {
		c.relay.On()
	} else {
		c.relay.Off()
	}
}

func (c *Controller) SetHeatingState(newState int) {
	// update state
	c.state.heatingState = newState

	// apply
	c.updateHeatingState()

	// update view
	c.updateView()
}

func (c *Controller) updateView() {
	c.view.ViewState(*c.state)
}
