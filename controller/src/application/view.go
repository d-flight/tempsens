package application

import (
	"tempsens/adapter"
	"tempsens/data"
	"time"
)

type View struct {
	gobotAdapter   *adapter.Gobot
	homekitAdapter *adapter.Homekit
	promAdapter    *adapter.PrometheusAdapter
	booting        bool
}

func NewView(
	gobotAdapter *adapter.Gobot,
	homekitAdapter *adapter.Homekit,
	prometheusAdapter *adapter.PrometheusAdapter,
) *View {
	v := &View{gobotAdapter, homekitAdapter, prometheusAdapter, true}

	// booting LED blinks blue
	v.gobotAdapter.StatusLed.SetColor(data.Blue())
	v.gobotAdapter.StatusLed.Blink(2*time.Second, 2*time.Second)

	return v
}

func (v *View) finishBooting() { v.booting = false }

func (v *View) ViewState(state data.State) {
	if v.booting {
		return
	}

	v.ViewHeatingState(state.GetHeatingState())
	v.ViewDesiredTemperature(state.GetDesiredTemperature())
	v.ViewLatestReading(state.GetLatestReading())
}

func (v *View) ViewHeatingState(state data.HeatingState) {
	switch state {
	case data.HEATING_STATE_IDLE:
		// update relay
		v.gobotAdapter.HeatingRelay.Off()

		// update LED
		v.gobotAdapter.StatusLed.SetColor(data.Green())
		v.gobotAdapter.StatusLed.On()
	case data.HEATING_STATE_OFF:
		// update relay
		v.gobotAdapter.HeatingRelay.Off()

		// update LED
		v.gobotAdapter.StatusLed.SetColor(data.None())
		v.gobotAdapter.StatusLed.Off()
	case data.HEATING_STATE_ON:
		// update relay
		v.gobotAdapter.HeatingRelay.On()

		// update LED
		v.gobotAdapter.StatusLed.SetColor(data.Red())
		v.gobotAdapter.StatusLed.On()
	}

	// update HomeKit
	v.homekitAdapter.SetHeatingState(state)

	// update metrics
	v.promAdapter.RecordHeatingState(state)
}

func (v *View) ViewDesiredTemperature(t data.Temperature) {
	// update HomeKit
	v.homekitAdapter.SetDesiredTemperature(t)

	// update metrics
	v.promAdapter.RecordDesiredTemperature(t)
}

func (v *View) ViewLatestReading(reading *data.Reading) {
	if reading == nil {
		return
	}

	// update HomeKit
	v.homekitAdapter.SetLatestReading(reading)

	// update metrics
	v.promAdapter.RecordLatestReading(reading)
}
