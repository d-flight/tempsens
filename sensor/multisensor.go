package sensor

// MultiSensor ...
type MultiSensor struct {
	sensors []Sensor
}

func (m MultiSensor) Read() (reading Reading, err error) {
	var readings []Reading
	for _, s := range m.sensors {
		r, e := s.Read()
		if e != nil {
			err = e
			return
		}

		readings = append(readings, r)
	}

	return Avg(readings), nil
}

// Start ...
func (m *MultiSensor) Start() (err error) {
	for _, sensor := range m.sensors {
		err = sensor.Start()
		if err != nil {
			break
		}
	}

	return err
}

// Stop ...
func (m *MultiSensor) Stop() (err error) {
	for _, sensor := range m.sensors {
		err = sensor.Stop()
		if err != nil {
			break
		}
	}

	return err
}

// NewMultiSensor ...
func NewMultiSensor(sensors []Sensor) MultiSensor {
	return MultiSensor{sensors: sensors}
}
