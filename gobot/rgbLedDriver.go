package gobot

import (
	"fmt"
	"tempsens/data"
	"time"

	_g "gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// RgbLedDriver represents a digital RGB Led
type RgbLedDriver struct {
	pinRed     string
	redColor   byte
	pinGreen   string
	greenColor byte
	pinBlue    string
	blueColor  byte
	name       string
	connection gpio.DigitalWriter
	high       bool
	blink      chan bool
}

// NewRgbLedDriver return a new RgbLedDriver given a DigitalWriter and
// 3 pins: redPin, greenPin, and bluePin
//
// Adds the following API Commands:
//	"SetRGB" - See RgbLedDriver.SetRGB
//	"Toggle" - See RgbLedDriver.Toggle
//	"On" - See RgbLedDriver.On
//	"Off" - See RgbLedDriver.Off
func NewRgbLedDriver(a gpio.DigitalWriter, redPin string, greenPin string, bluePin string) *RgbLedDriver {
	return &RgbLedDriver{
		name:       _g.DefaultName("RGBLED"),
		pinRed:     redPin,
		pinGreen:   greenPin,
		pinBlue:    bluePin,
		connection: a,
		high:       false,
	}
}

// Start implements the Driver interface
func (l *RgbLedDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *RgbLedDriver) Halt() (err error) {
	l.cancelBlink()
	return
}

// Name returns the RgbLedDrivers name
func (l *RgbLedDriver) Name() string { return l.name }

// SetName sets the RgbLedDrivers name
func (l *RgbLedDriver) SetName(n string) { l.name = n }

// Pin returns the RgbLedDrivers pins
func (l *RgbLedDriver) Pin() string {
	return "r=" + l.pinRed + ", g=" + l.pinGreen + ", b=" + l.pinBlue
}

// RedPin returns the RgbLedDrivers redPin
func (l *RgbLedDriver) RedPin() string { return l.pinRed }

// GreenPin returns the RgbLedDrivers redPin
func (l *RgbLedDriver) GreenPin() string { return l.pinGreen }

// BluePin returns the RgbLedDrivers bluePin
func (l *RgbLedDriver) BluePin() string { return l.pinBlue }

// Connection returns the RgbLedDriver Connection
func (l *RgbLedDriver) Connection() _g.Connection {
	return l.connection.(_g.Connection)
}

// State return true if the led is On and false if the led is Off
func (l *RgbLedDriver) State() bool {
	return l.high
}

// On sets the led's pins to their various states
func (l *RgbLedDriver) On() (err error) {
	l.cancelBlink()
	return l.turnOn()
}

func (l *RgbLedDriver) turnOn() (err error) {
	if err = l.SetLevel(l.pinRed, l.redColor); err != nil {
		return
	}

	if err = l.SetLevel(l.pinGreen, l.greenColor); err != nil {
		return
	}

	if err = l.SetLevel(l.pinBlue, l.blueColor); err != nil {
		return
	}

	l.high = true
	return
}

// Off sets the led to black.
func (l *RgbLedDriver) Off() (err error) {
	l.cancelBlink()
	return l.turnOff()
}

func (l *RgbLedDriver) turnOff() (err error) {
	if err = l.SetLevel(l.pinRed, 0); err != nil {
		return
	}

	if err = l.SetLevel(l.pinGreen, 0); err != nil {
		return
	}

	if err = l.SetLevel(l.pinBlue, 0); err != nil {
		return
	}

	l.high = false
	return
}

// Toggle sets the led to the opposite of it's current state
func (l *RgbLedDriver) Toggle() (err error) {
	if l.State() {
		err = l.Off()
	} else {
		err = l.On()
	}
	return
}

// SetLevel sets the led to the specified color level
func (l *RgbLedDriver) SetLevel(pin string, level byte) (err error) {
	if 0xff == level || 0 == level {
		if writer, ok := l.connection.(gpio.DigitalWriter); ok {
			return writer.DigitalWrite(pin, level)
		}
	}

	if writer, ok := l.connection.(gpio.PwmWriter); ok {
		return writer.PwmWrite(pin, level)
	}
	return gpio.ErrPwmWriteUnsupported
}

// SetColor sets the Red Green Blue value of the LED.
func (l *RgbLedDriver) SetColor(c *data.Color) error {
	l.redColor = c.Red()
	l.greenColor = c.Green()
	l.blueColor = c.Blue()

	return l.On()
}

func (l *RgbLedDriver) Blink(i time.Duration) {
	l.blink = make(chan bool)

	go func() {
	Loop:
		for range time.Tick(i) {
			if l.State() {
				l.turnOff()
			} else {
				l.turnOn()
			}

			select {
			case <-l.blink:
				fmt.Println("received halt")
				close(l.blink)
				l.blink = nil
				break Loop
			default:
			}
		}
	}()
}

func (l *RgbLedDriver) cancelBlink() {
	if l.blink != nil {
		l.blink <- true
	}
}
