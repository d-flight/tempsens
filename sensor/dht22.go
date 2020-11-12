package sensor

import (
	"errors"
	"fmt"
	"time"

	"github.com/MichaelS11/go-dht"
)

// keeps 10 readings which are averaged to a reading when Read() is called
const readingCount = 10

// DHT22 ...
type DHT22 struct {
	dht      *dht.DHT
	readings []Reading
	stop     chan struct{}
	stopped  chan struct{}
}

func (d DHT22) Read() (result Reading, err error) {
	if 1 > len(d.readings) {
		err = errors.New("No readings for sensor")
	} else {
		result = Avg(d.readings)
	}

	return result, err
}

// Start ...
func (d *DHT22) Start() error {
	d.stop = make(chan struct{})
	d.stopped = make(chan struct{})

	go d.poll()

	return nil
}

func (d *DHT22) poll() {
	const interval = 2 * time.Second

Loop:

	for range time.Tick(interval) {
		startTime := time.Now()
		h, t, e := d.dht.Read()

		if e == nil {
			reading := Reading{Temperature: t, Humidity: h}
			d.readings = append(d.sliceReadings(), reading)
		}

		select {
		case <-time.After(interval - time.Since(startTime)):
		case <-d.stop:
			break Loop
		}
	}

	close(d.stopped)
}

func (d DHT22) sliceReadings() []Reading {
	const max = readingCount - 1
	if max < len(d.readings) {
		return d.readings[:max]
	}

	return d.readings
}

// Stop ...
func (d DHT22) Stop() error {
	close(d.stop)

	return nil
}

// NewDHT22 ...
func NewDHT22(pin string) (d DHT22, err error) {
	err = dht.HostInit()
	if err != nil {
		fmt.Println("Hostinit error", err)
		return
	}

	dhtInstance, err := dht.NewDHT(pin, dht.Celsius, "")
	if err != nil {
		fmt.Println("NewDHT error", err)
		return
	}

	return DHT22{dht: dhtInstance}, nil
}
