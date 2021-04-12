package adapter

import (
	"tempsens/data"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
)

// thermostat configuration
const (
	thermostat_default_temp = 22
	thermostat_min_temp     = 0
	thermostat_max_temp     = 35
	thermostat_temp_step    = 0.1
)

type Homekit struct {
	pin        string
	thermostat *accessory.Thermostat
	humidity   *characteristic.CurrentRelativeHumidity

	OnDesiredTemperatureChanged func(data.Temperature)
	OnHeatingStateChanged       func(data.HeatingState)
}

func NewHomekitAdapter(info accessory.Info, pin string) *Homekit {
	humidity := characteristic.NewCurrentRelativeHumidity()
	thermostat := accessory.NewThermostat(
		info,
		thermostat_default_temp,
		thermostat_min_temp,
		thermostat_max_temp,
		thermostat_temp_step,
	)
	thermostat.Thermostat.AddCharacteristic(humidity.Characteristic)

	return &Homekit{
		pin,
		thermostat,
		humidity,
		func(t data.Temperature) {},
		func(s data.HeatingState) {},
	}
}

func (d *Homekit) Boot() {
	// add listeners
	d.thermostat.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(temperature float64) {
		d.OnDesiredTemperatureChanged(data.FromCelsius(temperature))
	})
	d.thermostat.Thermostat.TargetHeatingCoolingState.OnValueRemoteUpdate(func(state int) {
		var applicationState data.HeatingState
		if characteristic.TargetHeatingCoolingStateOff == state {
			applicationState = data.HEATING_STATE_OFF
		} else {
			applicationState = data.HEATING_STATE_ON
		}

		d.OnHeatingStateChanged(applicationState)
	})

	// configure & start transport
	transport, err := hc.NewIPTransport(
		hc.Config{Pin: d.pin},
		d.thermostat.Accessory,
	)
	if err != nil {
		panic(err)
	}

	hc.OnTermination(func() { <-transport.Stop() })

	transport.Start()
}

func (d *Homekit) SetHeatingState(state data.HeatingState) {
	hcState := characteristic.CurrentHeatingCoolingStateOff

	switch state {
	case data.HEATING_STATE_IDLE:
		hcState = characteristic.CurrentHeatingCoolingStateCool
	case data.HEATING_STATE_ON:
		hcState = characteristic.CurrentHeatingCoolingStateHeat
	case data.HEATING_STATE_OFF:
	default:
	}

	d.thermostat.Thermostat.CurrentHeatingCoolingState.SetValue(hcState)
}

func (d *Homekit) SetDesiredTemperature(temp data.Temperature) {
	d.thermostat.Thermostat.TargetTemperature.SetValue(temp.InCelsius())
}

func (d *Homekit) SetLatestReading(reading *data.Reading) {
	d.thermostat.Thermostat.CurrentTemperature.SetValue(reading.Temperature.InCelsius())
	d.humidity.SetValue(reading.Humidity)
}
