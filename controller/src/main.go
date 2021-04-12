package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tempsens/adapter"
	"tempsens/application"
	"tempsens/data"

	"github.com/brutella/hc/accessory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// configuration for the tempsens application
const (
	homekit_name = "tempsens-v2"
	homekit_id   = "diy-42069-v2"
	homekit_pin  = "00102003"
)

func main() {
	// setup prometheus
	promAdapter := adapter.NewPrometheusAdapter()

	// setup mqtt
	mqttAdapter := adapter.NewMqttAdapter(
		mqtt.NewClientOptions().
			SetClientID("tempsens controller").
			AddBroker("tcp://localhost:1883"),
	)

	// setup homekit
	homekitAdapter := adapter.NewHomekitAdapter(
		accessory.Info{
			Name:         homekit_name,
			SerialNumber: homekit_id,
			Model:        "DIY",
			Manufacturer: "github.com/d-flight",
		},
		homekit_pin,
	)

	// setup application
	view := application.NewView(homekitAdapter, promAdapter)
	controller := application.NewController(view, data.NewState(), mqttAdapter)

	// connect mqtt
	mqttAdapter.OnNewReport = controller.HandleNewReport

	// connect homekit
	homekitAdapter.OnDesiredTemperatureChanged = func(t data.Temperature) {
		controller.SetDesiredTemperature(t)
	}
	homekitAdapter.OnHeatingStateChanged = controller.SetHeatingState

	// start homekit, mqtt, prometheus
	go homekitAdapter.Boot()
	go mqttAdapter.Boot()
	go promAdapter.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	if sig, ok := <-c; ok {
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		os.Exit(1)
	}
}
