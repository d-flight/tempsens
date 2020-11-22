package data

// Reading ...
type Reading struct {
	Temperature Temperature
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

// Equals compares two Readings for equality
func (r *Reading) Equals(other *Reading, cap *DeltaCap) bool {
	if nil == other {
		return false
	}
	if r == other {
		return true
	}

	deltaTemperature := Abs(int(r.Temperature*100) - int(other.Temperature*100))
	deltaHumidity := Abs(int(r.Humidity*100) - int(other.Humidity*100))

	return deltaTemperature < cap.Temperature && deltaHumidity < cap.Humidity
}
