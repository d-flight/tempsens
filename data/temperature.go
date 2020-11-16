package data

// Temperature ...
type Temperature int64

const (
	TEMPERATURE_MIN = -4000
	TEMPERATURE_MAX = +8000
)

// InCelsius ...
func (t Temperature) InCelsius() float64 {
	return float64(t) / 100.
}

// IsValid ...
func (t Temperature) IsValid() bool {
	return TEMPERATURE_MIN <= t && TEMPERATURE_MAX >= t
}

func InvalidTemperature() Temperature {
	return Temperature(TEMPERATURE_MIN - 1)
}

// FromCelsius ...
func FromCelsius(celsius float64) Temperature {
	return Temperature(celsius * 100)
}
