package application

import (
	"tempsens/data"
	"tempsens/gobot"
	"tempsens/sensor"
	"time"

	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
)

type View struct {
	statusLed              *gobot.RgbLedDriver
	thermostat             *accessory.Thermostat
	humidityCharacteristic *characteristic.CurrentRelativeHumidity
	booting                bool
}

func NewView(statusLed *gobot.RgbLedDriver, thermostat *accessory.Thermostat, humidityCharacteristic *characteristic.CurrentRelativeHumidity) *View {
	thermostat.Thermostat.TargetTemperature.Unit = characteristic.UnitCelsius
	thermostat.Thermostat.CurrentTemperature.Unit = characteristic.UnitCelsius

	booting := true

	v := &View{
		statusLed,
		thermostat,
		humidityCharacteristic,
		booting,
	}

	// booting LED blinks blue
	v.statusLed.SetColor(data.Blue())
	v.statusLed.Blink(2*time.Second, 2*time.Second)

	return v
}

func (v *View) finishBooting() { v.booting = false }

func (v *View) ViewState(state State) {
	if v.booting {
		return
	}

	v.ViewHeatingState(state.GetHeatingState())
	v.ViewDesiredTemperature(state.desiredTemperature)
	v.ViewLatestReading(state.latestReading)
}

func (v *View) ViewHeatingState(state int) {
	switch state {
	case HEATING_STATE_IDLE:
		// update LED
		v.statusLed.SetColor(data.Orange())
		v.statusLed.On()

		// update HomeKit
		v.thermostat.Thermostat.CurrentHeatingCoolingState.SetValue(characteristic.CurrentHeatingCoolingStateHeat)
	case HEATING_STATE_OFF:
		// update LED
		v.statusLed.SetColor(data.None())
		v.statusLed.Off()

		// update HomeKit
		v.thermostat.Thermostat.CurrentHeatingCoolingState.SetValue(characteristic.CurrentHeatingCoolingStateOff)
	case HEATING_STATE_ON:
		// update LED
		v.statusLed.SetColor(data.Red())
		v.statusLed.Blink(5*time.Second, 2*time.Second)

		// update HomeKit
		v.thermostat.Thermostat.CurrentHeatingCoolingState.SetValue(characteristic.CurrentHeatingCoolingStateHeat)
	}
}

func (v *View) ViewDesiredTemperature(t data.Temperature) {
	// update HomeKit
	v.thermostat.Thermostat.TargetTemperature.SetValue(t.InCelsius())
}

func (v *View) ViewLatestReading(reading *sensor.Reading) {
	if reading == nil {
		return
	}

	// update HomeKit
	v.thermostat.Thermostat.CurrentTemperature.SetValue(reading.Temperature.InCelsius())
	v.humidityCharacteristic.SetValue(reading.Humidity)
}
