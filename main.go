package main

import (
	"fmt"
	tbot "tempsens/gobot"
	"tempsens/sensor"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	// setup gobot
	r := raspi.NewAdaptor()
	dht22 := tbot.NewDHT22Driver(r, "7")
	button := gpio.NewButtonDriver(r, "11")

	work := func() {
		dht22.On(gpio.Error, func(data interface{}) {
			fmt.Println(data)
		})
		dht22.On(tbot.TemperatureUpdated, func(data interface{}) {
			reading := data.(*sensor.Reading)
			fmt.Printf("New Reading: T %v*C, H %v%%\n", reading.Temperature, reading.Humidity)
		})

		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("Button: On")
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Button: Off")
		})
	}

	// start gobot
	robot := gobot.NewRobot("tempsens bot",
		[]gobot.Connection{r},
		[]gobot.Device{dht22, button},
		work,
	)

	robot.Start()

}
