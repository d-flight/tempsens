package adapter

import (
	// "net/http"
	"tempsens/data"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusAdapter struct {
	registry *prometheus.Registry

	actualTemperature  prometheus.Gauge
	actualHumidity     prometheus.Gauge
	desiredTemperature prometheus.Gauge
	isHeating          prometheus.Gauge
}

func NewPrometheusAdapter() *PrometheusAdapter {
	registry := prometheus.NewRegistry()
	metricFactory := promauto.With(registry)

	actualTemperature := metricFactory.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_actual_temperature",
		Help: "Actual temperature as percieved by the BME280 Sensor",
	})
	actualHumidity := metricFactory.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_actual_humidity",
		Help: "Actual humidity as percieved by the BME280 Sensor",
	})
	desiredTemperature := metricFactory.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_desired_temperature",
		Help: "Temperature desired either by the user",
	})
	isHeating := metricFactory.NewGauge(prometheus.GaugeOpts{
		Name: "tempsens_is_heating",
		Help: "Indicator whether the heating is enabled or not",
	})

	return &PrometheusAdapter{
		registry,
		actualTemperature,
		actualHumidity,
		desiredTemperature,
		isHeating,
	}
}

func (a *PrometheusAdapter) Start() error {
	// // application metrics
	// http.Handle("/metrics/app", promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{}))
	// // process metrics
	// http.Handle("/metrics/process", promhttp.Handler())

	// return http.ListenAndServe(":2112", nil)

	// disabled for testing
	return nil
}

func (a *PrometheusAdapter) RecordLatestReading(reading *data.Reading) {
	a.actualHumidity.Set(reading.Humidity)
	a.actualTemperature.Set(reading.Temperature.InCelsius())
}

func (a *PrometheusAdapter) RecordHeatingState(state data.HeatingState) {
	if data.HEATING_STATE_ON == state {
		a.isHeating.Set(1)
	} else {
		a.isHeating.Set(0)
	}
}

func (a *PrometheusAdapter) RecordDesiredTemperature(temp data.Temperature) {
	a.desiredTemperature.Set(temp.InCelsius())
}
