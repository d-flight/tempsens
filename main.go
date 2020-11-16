package main

import (
	"fmt"
	"os"
	"os/signal"
	"tempsens/application"
	"tempsens/data"
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
	raspi := raspi.NewAdaptor()
	rgb := tbot.NewRgbLedDriver(raspi, "11", "12", "13")
	dht22 := tbot.NewDHT22Driver(raspi, "7")
	relay := gpio.NewRelayDriver(raspi, "8")

	// setup homekit
	info := accessory.Info{Name: "Thermostat"}
	// limitations as defined by the DHT22 Datasheet
	hcThermostat := accessory.NewThermostat(info, 22, -40., 80., 0.1)
	// DHT22 supports humidity, which we also expose to homekit
	hcHumidity := characteristic.NewCurrentRelativeHumidity()
	hcThermostat.Thermostat.AddCharacteristic(hcHumidity.Characteristic)

	// setup application
	view := application.NewView(rgb, hcThermostat, hcHumidity)
	controller := application.NewController(view, application.NewState(), relay)

	// connect gobot
	work := func() {
		dht22.On(gpio.Error, func(data interface{}) {
			fmt.Println(data)
		})
		dht22.On(tbot.TemperatureUpdated, func(data interface{}) {
			reading := data.(*sensor.Reading)
			fmt.Printf("New Reading: T %v*C, H %v%%\n", reading.Temperature, reading.Humidity)

			controller.SetLatestReading(reading)
		})
	}

	// connect homekit
	hcThermostat.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(val float64) {
		controller.SetDesiredTemperature(data.FromCelsius(val), true)
	})
	hcThermostat.Thermostat.CurrentHeatingCoolingState.OnValueRemoteUpdate(func(val int) {
		if characteristic.CurrentHeatingCoolingStateOff == val {
			controller.SetHeatingState(application.HEATING_STATE_OFF)
		} else if characteristic.CurrentHeatingCoolingStateHeat == val {
			controller.SetHeatingState(application.HEATING_STATE_ON)
		}
	})

	// start gobot
	go gobot.NewRobot("tempsens bot",
		[]gobot.Connection{raspi},
		[]gobot.Device{rgb, dht22, relay},
		work,
	).Start()

	// start homekit
	transportConfig := hc.Config{Pin: "00102003"}
	hcTransport, err := hc.NewIPTransport(transportConfig, hcThermostat.Accessory)
	if err != nil {
		panic(err)
	}

	hc.OnTermination(func() {
		<-hcTransport.Stop()
	})

	fmt.Println("starting homekit service..")
	go hcTransport.Start()

	// default desired temperature
	controller.SetDesiredTemperature(data.FromCelsius(23.5), false)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	if sig, ok := <-c; ok {
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		os.Exit(1)
	}
}
