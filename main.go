package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gobot.io/x/gobot/platforms/raspi"

	"tempsens/adapter"
	"tempsens/application"
	"tempsens/data"

	"github.com/brutella/hc/accessory"
)

// configuration for the tempsens application
const (
	homekit_name = "tempsens"
	homekit_id   = "diy-42069"
	homekit_pin  = "00102003"
)

func main() {
	// setup prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

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
	view := application.NewView(gobotAdapter, homekitAdapter)
	controller := application.NewController(view, data.NewState())
	schedule := application.NewSchedule(*data.NewHeatingSchedule(
		&data.Setting{Hour: 6, Temperature: data.FromCelsius(23)},
		&data.Setting{Hour: 22, Temperature: data.FromCelsius(20)},
	), controller)

	// connect gobot
	gobotAdapter.OnNewReading = controller.SetLatestReading
	gobotAdapter.OnScheduleButton = func() {
		controller.SetUserControlled(false)
		schedule.Trigger(time.Now())
	}

	// connect homekit
	homekitAdapter.OnDesiredTemperatureChanged = func(t data.Temperature) {
		controller.SetUserControlled(true)
		controller.SetDesiredTemperature(t, true)
	}
	homekitAdapter.OnHeatingStateChanged = controller.SetHeatingState

	// start homekit, gobot, schedule
	go homekitAdapter.Boot()
	go gobotAdapter.Boot()
	go schedule.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	if sig, ok := <-c; ok {
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		schedule.Stop()
		os.Exit(1)
	}
}
