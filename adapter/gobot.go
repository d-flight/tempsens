package adapter

import (
	"fmt"
	"tempsens/data"
	tbot "tempsens/gobot"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
)

// pin configuration
const (
	// status rgb
	pin_status_rgb_red   = "11"
	pin_status_rgb_green = "12"
	pin_status_rgb_blue  = "13"

	// heating relay
	pin_heating_relay = "8"

	// user control led
	pin_user_control_led = "15"

	// button to switch to schedule
	pin_schedule_button = "16"
)

// i2c address configuration
const (
	i2c_bme280_bus     = 1
	i2c_bme280_address = 0x76
)

// adapter configuration
const (
	reading_polling_rate = 10 * time.Second
)

type Gobot struct {
	Name    string
	adaptor *gobot.Adaptor

	// TODO: abstract output components as well
	StatusLed      *tbot.RgbLedDriver
	HeatingRelay   *gpio.RelayDriver
	UserControlLed *gpio.LedDriver

	bme280         *i2c.BME280Driver
	scheduleButton *gpio.ButtonDriver

	OnNewReading     func(*data.Reading)
	OnScheduleButton func()
}

// NewGobotAdapter creates a new adaptor for tempsens using gobot.
// The adapter must implement the gpio.DigitalWriter, gpio.DigitalReader and i2c.Connector interface from gobot
func NewGobotAdapter(adaptor gobot.Adaptor, name string) *Gobot {
	return &Gobot{
		Name:             name,
		adaptor:          &adaptor,
		StatusLed:        tbot.NewRgbLedDriver(adaptor.(gpio.DigitalWriter), pin_status_rgb_red, pin_status_rgb_green, pin_status_rgb_blue),
		HeatingRelay:     gpio.NewRelayDriver(adaptor.(gpio.DigitalWriter), pin_heating_relay),
		UserControlLed:   gpio.NewLedDriver(adaptor.(gpio.DigitalWriter), pin_user_control_led),
		bme280:           i2c.NewBME280Driver(adaptor.(i2c.Connector), i2c.WithBus(i2c_bme280_bus), i2c.WithAddress(i2c_bme280_address)),
		scheduleButton:   gpio.NewButtonDriver(adaptor.(gpio.DigitalReader), pin_schedule_button),
		OnNewReading:     func(*data.Reading) {},
		OnScheduleButton: func() {},
	}
}

func (d *Gobot) Boot() error {
	return gobot.NewRobot(
		d.Name,
		[]gobot.Connection{*d.adaptor},
		[]gobot.Device{d.StatusLed, d.HeatingRelay, d.bme280, d.UserControlLed, d.scheduleButton},
		func() {
			// poll reading once and then in the configured polling rate
			d.pollReading()

			gobot.Every(reading_polling_rate, d.pollReading)

			// handle button release
			err := d.scheduleButton.On(gpio.ButtonRelease, func(s interface{}) {
				d.OnScheduleButton()
			})

			if err != nil {
				panic(err)
			}
		},
	).Start()
}

func (d *Gobot) pollReading() {
	if reading, err := d.getLatestReading(); err == nil {
		d.OnNewReading(reading)
	}
}

func (d *Gobot) getLatestReading() (reading *data.Reading, err error) {
	temp, err := d.bme280.Temperature()
	if err != nil {
		fmt.Println("Error reading temperature", err)
		return
	}
	humidity, err := d.bme280.Humidity()
	if err != nil {
		fmt.Println("Error reading humidity", err)
		return
	}

	reading = &data.Reading{
		Temperature: data.FromCelsius(float64(temp)),
		Humidity:    float64(humidity),
	}

	fmt.Printf("New Reading: T %v*C, H %v%%\n", reading.Temperature, reading.Humidity)

	return reading, nil
}
