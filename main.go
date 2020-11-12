package main

import (
	"fmt"
	"time"

	"tempsens/sensor"
)

func main() {
	dht1Pin := "GPIO4"
	dht1, e := sensor.NewDHT22(dht1Pin)

	if e != nil {
		fmt.Println(e)
		return
	}

	ms := sensor.NewMultiSensor([]sensor.Sensor{&dht1})

	ms.Start()
	defer ms.Stop()

	for range time.Tick(2 * time.Second) {
		reading, err := ms.Read()

		if err != nil {
			fmt.Println("Read error", err)
		} else {
			fmt.Printf("temp: %v, humidity: %v \n", reading.Temperature, reading.Humidity)
		}
	}

}
