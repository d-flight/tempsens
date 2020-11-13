package main

import (
	"fmt"
	tbot "tempsens/gobot"
	"tempsens/sensor"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
)

func main() {
	// setup gobot
	r := raspi.NewAdaptor()
	dht22 := tbot.NewDHT22Driver(r, "7")
	button := gpio.NewButtonDriver(r, "11")

	// setup homekit
	info := accessory.Info{Name: "Thermostat"}
	// limitations as defined by the DHT22 Datasheet
	hcThermostat := accessory.NewThermostat(info, 0., -40., 80., 0.1)
	// DHT22 supports humidity, which we also expose to homekit
	hcHumidity := characteristic.NewCurrentRelativeHumidity()
	hcThermostat.Thermostat.AddCharacteristic(hcHumidity.Characteristic)

	transportConfig := hc.Config{Pin: "00102003"}
	hcTransport, err := hc.NewIPTransport(transportConfig, hcThermostat.Accessory)
	if err != nil {
		panic(err)
	}

	// connect them
	work := func() {
		dht22.On(gpio.Error, func(data interface{}) {
			fmt.Println(data)
		})
		dht22.On(tbot.TemperatureUpdated, func(data interface{}) {
			reading := data.(*sensor.Reading)
			fmt.Printf("New Reading: T %v*C, H %v%%\n", reading.Temperature, reading.Humidity)

			// new reading received, update homekit state
			hcThermostat.Thermostat.CurrentTemperature.SetValue(reading.Temperature)
			hcHumidity.SetValue(reading.Humidity)
		})

		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("Button: On")
		})
		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("Button: Off")
		})
	}

	// start gobot
	go gobot.NewRobot("tempsens bot",
		[]gobot.Connection{r},
		[]gobot.Device{dht22, button},
		work,
	).Start()

	// start homekit
	go func() {
		hcThermostat.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(val float64) {
			fmt.Printf("Remote update: %v", val)
		})
	}()

	hc.OnTermination(func() {
		<-hcTransport.Stop()
	})

	fmt.Println("starting homekit service..")
	hcTransport.Start()
}
