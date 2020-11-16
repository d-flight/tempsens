package sensor

import (
	"fmt"
	"tempsens/data"
)

// Reading ...
type Reading struct {
	Temperature data.Temperature
	Humidity    float64
}

/*
	Delta when two readings are not equal anymore.
	Due to problems when comparing floating point numbers,
	these are 1/100 of the actual value. e.g. 1.2% Humidity equals 120 delta.
*/
type DeltaCap struct {
	Temperature int
	Humidity    int
}

const (
	// according to the datasheet the DHT is accurate up to 0.5 celsius and 2% humidity
	deltaTemperature = 50
	deltaHumidity    = 200
)

// Avg ...
func Avg(readings []Reading) *Reading {
	var temperature, humidity, total float64
	for _, r := range readings {
		temperature += r.Temperature.InCelsius()
		humidity += r.Humidity
		total++
	}

	fmt.Printf("Calculated from %v readings\n", total)

	return &Reading{
		Temperature: data.FromCelsius(temperature / total),
		Humidity:    humidity / total,
	}
}

// Equals compares two Readings for equality
func (r *Reading) Equals(other *Reading, cap *DeltaCap) bool {
	if nil == other {
		return false
	}
	if r == other {
		return true
	}

	deltaTemperature := data.Abs(int(r.Temperature*100) - int(other.Temperature*100))
	deltaHumidity := data.Abs(int(r.Humidity*100) - int(other.Humidity*100))

	return deltaTemperature < cap.Temperature && deltaHumidity < cap.Humidity
}
