package sensor

// Sensor ...
type Sensor interface {
	Read() (Reading, error)
	Start() error
	Stop() error
}
