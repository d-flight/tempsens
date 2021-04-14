package application

import (
	"tempsens/adapter"
	"tempsens/data"
)

type View struct {
	homekitAdapter *adapter.Homekit
	promAdapter    *adapter.PrometheusAdapter
}

func NewView(
	homekitAdapter *adapter.Homekit,
	prometheusAdapter *adapter.PrometheusAdapter,
) *View {
	return &View{homekitAdapter, prometheusAdapter}
}

func (v *View) ViewState(state data.State) {
	v.ViewHeatingState(state.GetHeatingState())
	v.ViewDesiredTemperature(state.GetDesiredTemperature())
	v.ViewLatestReading(state.GetLatestReading())
}

func (v *View) ViewHeatingState(state data.HeatingState) {
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
