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
	"gobot.io/x/gobot/platforms/raspi"
)

// configuration for the tempsens application
const (
	homekit_name = "tempsens"
	homekit_id   = "diy-42069"
	homekit_pin  = "00102003"
)

func main() {
	// setup prometheus
	promAdapter := adapter.NewPrometheusAdapter()

	// setup gobot
	gobotAdapter := adapter.NewGobotAdapter(
		raspi.NewAdaptor(), "tempsens",
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
	view := application.NewView(gobotAdapter, homekitAdapter, promAdapter)
	controller := application.NewController(view, data.NewState())

	// connect gobot
	gobotAdapter.OnNewReading = controller.SetLatestReading

	// connect homekit
	homekitAdapter.OnDesiredTemperatureChanged = func(t data.Temperature) {
		controller.SetDesiredTemperature(t)
	}
	homekitAdapter.OnHeatingStateChanged = controller.SetHeatingState

	// start homekit, gobot, prometheus
	go homekitAdapter.Boot()
	go gobotAdapter.Boot()
	go promAdapter.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	if sig, ok := <-c; ok {
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		os.Exit(1)
	}
}
