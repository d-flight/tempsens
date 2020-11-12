package sensor

import "fmt"

// Reading ...
type Reading struct {
	Temperature float64
	Humidity    float64
}

// Avg ...
func Avg(readings []Reading) Reading {
	var reading = Reading{Temperature: 0., Humidity: 0.}
	var total float64 = 0.
	for _, r := range readings {
		reading.Temperature += r.Temperature
		reading.Humidity += r.Humidity
		total++
	}

	fmt.Printf("Calculated from %v readings\n", total)

	return Reading{
		Temperature: reading.Temperature / total,
		Humidity:    reading.Humidity / total,
	}
}
