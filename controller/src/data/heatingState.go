package data

var (
	HEATING_STATE_OFF  = HeatingState(0)
	HEATING_STATE_ON   = HeatingState(1)
	HEATING_STATE_IDLE = HeatingState(2)
)

type HeatingState int
