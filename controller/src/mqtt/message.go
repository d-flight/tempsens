package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"

	"tempsens/data"
)

type ReportMessage struct {
	Desired data.Temperature
	Reading data.Reading
	State   data.HeatingState
}

func (m ReportMessage) Validate() (errs []error) {
	if !m.Desired.IsValid() {
		errs = append(errs, errors.New("Invalid desired temperature"))
	}

	if !m.Reading.Temperature.IsValid() {
		errs = append(errs, errors.New("Invalid reading temperature"))
	}

	if m.Reading.Humidity > 100 || 0 == m.Reading.Humidity {
		errs = append(errs, errors.New("Invalid reading humidity"))
	}

	return
}

// DeserializeReportMessage ...
func DeserializeReportMessage(rawMessage []byte) (*ReportMessage, error) {
	var message ReportMessage
	if e := json.Unmarshal(rawMessage, &message); e != nil {
		return nil, fmt.Errorf("error with unmarshal: %v", e)
	}

	return &message, nil
}

const CONTROL_MESSAGE_DESIRED_TEMPERATURE = 1
const CONTROL_MESSAGE_HEATING_STATE = 2

type ControlMessage struct {
	Type byte
}

// ChangeDesiredTemperatureMessage ...
type ChangeDesiredTemperatureMessage struct {
	Desired data.Temperature

	*ControlMessage
}

func NewChangeDesiredTemperatureMessage(desired data.Temperature) *ChangeDesiredTemperatureMessage {
	return &ChangeDesiredTemperatureMessage{
		desired,
		&ControlMessage{CONTROL_MESSAGE_DESIRED_TEMPERATURE},
	}
}

func (m *ChangeDesiredTemperatureMessage) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

type ToggleActiveMessage struct {
	Active bool

	*ControlMessage
}

func NewToggleActiveMessage(active bool) *ToggleActiveMessage {
	return &ToggleActiveMessage{
		active,
		&ControlMessage{CONTROL_MESSAGE_HEATING_STATE},
	}
}

func (m *ToggleActiveMessage) Serialize() ([]byte, error) {
	return json.Marshal(m)
}
