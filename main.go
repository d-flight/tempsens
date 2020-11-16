package main

import (
	"tempsens/application"
	"tempsens/data"
	tbot "tempsens/gobot"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	booting := true

	// setup a new controller
	c := application.NewController(&data.HeatingSchedule{
		// from 6AM until 10PM 22*C
		6: data.FromCelsius(22),
		// from 10PM until 6AM 18*C
		22: data.FromCelsius(18),
	}, application.STATE_IDLE)

	// setup gobot
	r := raspi.NewAdaptor()
	rgb := tbot.NewRgbLedDriver(r, "11", "12", "13")
	rgb.SetColor(data.Blue())
	rgb.Blink(1500 * time.Millisecond)
	// dht22 := tbot.NewDHT22Driver(r, "7")

	// // setup homekit
	// info := accessory.Info{Name: "Thermostat"}
	// // limitations as defined by the DHT22 Datasheet
	// hcThermostat := accessory.NewThermostat(info, 0., -40., 80., 0.1)
	// // DHT22 supports humidity, which we also expose to homekit
	// hcHumidity := characteristic.NewCurrentRelativeHumidity()
	// hcThermostat.Thermostat.AddCharacteristic(hcHumidity.Characteristic)

	// transportConfig := hc.Config{Pin: "00102003"}
	// hcTransport, err := hc.NewIPTransport(transportConfig, hcThermostat.Accessory)
	// if err != nil {
	// 	panic(err)
	// }

	// connect them
	work := func() {
		// dht22.On(gpio.Error, func(data interface{}) {
		// 	fmt.Println(data)
		// })
		// dht22.On(tbot.TemperatureUpdated, func(data interface{}) {
		// 	reading := data.(*sensor.Reading)
		// 	fmt.Printf("New Reading: T %v*C, H %v%%\n", reading.Temperature, reading.Humidity)

		// 	// new reading received, update homekit state
		// 	hcThermostat.Thermostat.CurrentTemperature.SetValue(reading.Temperature)
		// 	hcHumidity.SetValue(reading.Humidity)
		// })
	}

	// start gobot
	bot := gobot.NewRobot("tempsens bot",
		[]gobot.Connection{r},
		[]gobot.Device{},
		work,
	)

	bot.Start()

	// start homekit
	// go func() {
	// 	hcThermostat.Thermostat.TargetTemperature.OnValueRemoteUpdate(func(val float64) {
	// 		fmt.Printf("Remote update: %v", val)
	// 	})
	// }()

	// hc.OnTermination(func() {
	// 	<-hcTransport.Stop()
	// })

	// fmt.Println("starting homekit service..")
	// go hcTransport.Start()
}
