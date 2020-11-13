package sensor

import (
	"fmt"
)

// Reading ...
type Reading struct {
	Temperature float64
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
	var reading = Reading{Temperature: 0., Humidity: 0.}
	var total float64 = 0.
	for _, r := range readings {
		reading.Temperature += r.Temperature
		reading.Humidity += r.Humidity
		total++
	}

	fmt.Printf("Calculated from %v readings\n", total)

	return &Reading{
		Temperature: reading.Temperature / total,
		Humidity:    reading.Humidity / total,
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

	deltaTemperature := abs(int(r.Temperature*100) - int(other.Temperature*100))
	deltaHumidity := abs(int(r.Humidity*100) - int(other.Humidity*100))

	return deltaTemperature < cap.Temperature && deltaHumidity < cap.Humidity
}

func abs(v int) int {
	if 0 > v {
		return -v
	}
	return v
}
