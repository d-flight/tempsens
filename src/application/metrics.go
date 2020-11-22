package application

import (
	"tempsens/data"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricService struct {
	actualTemperature  prometheus.Gauge
	actualHumidity     prometheus.Gauge
	desiredTemperature prometheus.Gauge
	isHeating          prometheus.Gauge
	isUserControlled   prometheus.Gauge
}

func NewMetricService() *MetricService {
	actualTemperature := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_actual_temperature",
		Help: "Actual temperature as percieved by the BME280 Sensor",
	})
	actualHumidity := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_actual_humidity",
		Help: "Actual humidity as percieved by the BME280 Sensor",
	})
	desiredTemperature := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_desired_temperature",
		Help: "Temperature desired either by schedule or the user",
	})
	isHeating := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_is_heating",
		Help: "Indicator whether the heating is enabled or not",
	})
	isUserControlled := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_is_user_controlled",
		Help: "Indicator whether the desired temperature is controlled by the user or the schedule",
	})

	return &MetricService{
		actualTemperature,
		actualHumidity,
		desiredTemperature,
		isHeating,
		isUserControlled,
	}
}

func (m *MetricService) RecordLatestReading(reading *data.Reading) {
	m.actualHumidity.Set(reading.Humidity)
	m.actualTemperature.Set(reading.Temperature.InCelsius())
}

func (m *MetricService) RecordHeatingState(state data.HeatingState) {
	if data.HEATING_STATE_ON == state {
		m.isHeating.Set(1)
	} else {
		m.isHeating.Set(0)
	}
}

func (m *MetricService) RecordDesiredTemperature(temp data.Temperature) {
	m.desiredTemperature.Set(temp.InCelsius())
}

func (m *MetricService) RecordUserControlled(userControlled bool) {
	if userControlled {
		m.isUserControlled.Set(1)
	} else {
		m.isUserControlled.Set(0)
	}
}
