package gobot

/*
	This module reads the DHT22 sensor trough a Gobot Connection
	and emits events when the reading changes.

	All the communication and conversion of bytes to floats comes
	from github.com/MichaelS11/go-dht

*/
import (
	"fmt"
	"runtime/debug"
	"time"

	"tempsens/sensor"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	// TemperatureUpdated event
	TemperatureUpdated = "event_temp_updated"
)

const (
	high = 1
	low  = 0
)

// DHT22Driver a gobot driver for the DHT22
type DHT22Driver struct {
	name            string
	connection      gobot.Connection
	pin             string
	lastRead        time.Time
	lastReading     *sensor.Reading
	readingInterval time.Duration
	errors          int
	halt            chan bool
	readingDelta    *sensor.DeltaCap
	gobot.Eventer
}

// NewDHT22Driver ...
func NewDHT22Driver(c gobot.Connection, pin string) *DHT22Driver {
	d := &DHT22Driver{
		name:            gobot.DefaultName("DHT22"),
		connection:      c,
		pin:             pin,
		Eventer:         gobot.NewEventer(),
		halt:            make(chan bool),
		readingInterval: 2 * time.Second,
		// according to the DHT22 datasheet, temperature might be off by 0.5C, Humidity by 2%
		readingDelta: &sensor.DeltaCap{Temperature: 50, Humidity: 200},
	}

	d.AddEvent(gpio.Error)
	d.AddEvent(TemperatureUpdated)

	return d
}

// Name ...
func (d *DHT22Driver) Name() string {
	return d.name
}

// SetName ...
func (d *DHT22Driver) SetName(s string) {
	d.name = s
}

// Start ...
func (d *DHT22Driver) Start() (err error) {
	// set pin to high so ready for first read
	if err = d.resetSensor(); err != nil {
		return
	}

	// set lastRead a second before to give the pin a second to warm up
	d.lastRead = time.Now().Add(-1 * time.Second)

	// start polling
	go func() {
		for {
			startTime := time.Now()

			// set sleepTime
			var sleepTime time.Duration
			if d.errors < 57 {
				sleepTime = (2 * time.Second) + (time.Duration(d.errors) * 500 * time.Millisecond)
			} else {
				// sleep max of 30 seconds
				sleepTime = 30 * time.Second
			}
			sleepTime -= time.Since(d.lastRead)

			// sleep between 2 and 30 seconds
			time.Sleep(sleepTime)

			// read bits from sensor
			var bits []int
			if bits, err = d.readBits(); err != nil {
				d.Publish(gpio.Error, err)
			}

			// covert bits to humidity and temperature
			var newReading *sensor.Reading
			if err == nil {
				if newReading, err = createReadingFromBits(bits); err != nil {
					d.Publish(gpio.Error, err)
				}
			}

			// emit event if the reading changed
			if err == nil && !newReading.Equals(d.lastReading, d.readingDelta) {
				d.lastReading = newReading
				d.Publish(TemperatureUpdated, newReading)
			}

			select {
			// sleep before next read
			case <-time.After(d.readingInterval - time.Since(startTime)):
			// stop when halt was sent
			case <-d.halt:
				break
			}
		}
	}()

	return
}

func (d *DHT22Driver) readBits() (bits []int, err error) {
	// create variables ahead of time before critical timing part
	var i int
	var startTime time.Time
	var levelPrevious int
	var level int
	levels := make([]int, 0, 84)
	durations := make([]time.Duration, 0, 84)

	// set lastRead so do not read more than once every 2 seconds
	d.lastRead = time.Now()

	// disable garbage collection during critical timing part
	gcPercent := debug.SetGCPercent(-1)

	// send start low
	if err := d.writer().DigitalWrite(d.pin, low); err != nil {
		d.resetSensor()
		return nil, fmt.Errorf("pin out low error: %v", err)
	}
	time.Sleep(time.Millisecond)

	// switch to read
	if levelPrevious, err = d.reader().DigitalRead(d.pin); err != nil {
		return nil, fmt.Errorf("pin in error: %v", err)
	}

	// read levels and durations with busy read
	// hope there is a better way in the future
	// tried to use WaitForEdge but seems to miss edges and/or take too long to detect them
	// note that pin read takes around .2 microsecond (us) on Raspberry PI 3
	// note that 1000 microsecond (us) = 1 millisecond (ms)
	level = levelPrevious
	for i = 0; i < 84; i++ {
		startTime = time.Now()
		for levelPrevious == level && time.Since(startTime) < time.Millisecond {
			if level, err = d.reader().DigitalRead(d.pin); err != nil {
				return nil, fmt.Errorf("pin in error when reading sensor: %v", err)
			}
		}
		durations = append(durations, time.Since(startTime))
		levels = append(levels, levelPrevious)
		levelPrevious = level
	}

	// enable garbage collection, done with critical part
	debug.SetGCPercent(gcPercent)

	// set pin to high so ready for next time
	d.resetSensor()

	// get last low reading so know start of data
	var endNumber int
	for i = len(levels) - 1; ; i-- {
		if levels[i] == low {
			endNumber = i
			break
		}
		if i < 80 {
			// not enough readings, i = 79 means endNumber is 78 or less
			return nil, fmt.Errorf("missing some readings - low level not found")
		}
	}
	startNumber := endNumber - 79

	// covert pulses into bits and check high levels
	bits = make([]int, 40)
	index := 0
	for i = startNumber; i < endNumber; i += 2 {
		// check high levels
		if levels[i] != high {
			return nil, fmt.Errorf("missing some readings - level not high")
		}
		// high should not be longer then 90 microseconds
		if durations[i] > 90*time.Microsecond {
			return nil, fmt.Errorf("missing some readings - high level duration too long: %v", durations[i])
		}
		// bit is 0 if less than or equal to 30 microseconds
		if durations[i] > 30*time.Microsecond {
			// bit is 1 if more than 30 microseconds
			bits[index] = 1
		}
		index++
	}

	// check low levels
	for i = startNumber + 1; i < endNumber+1; i += 2 {
		// check low levels
		if levels[i] != low {
			return nil, fmt.Errorf("missing some readings - level not low")
		}
		// low should not be longer then 70 microseconds
		if durations[i] > 70*time.Microsecond {
			return nil, fmt.Errorf("missing some readings - low level duration too long: %v", durations[i])
		}
		// low should not be shorter then 35 microseconds
		if durations[i] < 35*time.Microsecond {
			return nil, fmt.Errorf("missing some readings - low level duration too short: %v", durations[i])
		}
	}

	return bits, nil
}

func (d *DHT22Driver) resetSensor() error {
	return d.writer().DigitalWrite(d.pin, high)
}

func (d *DHT22Driver) writer() gpio.DigitalWriter {
	return d.connection.(gpio.DigitalWriter)
}

func (d *DHT22Driver) reader() gpio.DigitalReader {
	return d.connection.(gpio.DigitalReader)
}

// Halt ...
func (d *DHT22Driver) Halt() (err error) {
	d.halt <- true
	return
}

// Connection ...
func (d *DHT22Driver) Connection() gobot.Connection {
	return d.connection
}

// GetReading ...
func (d *DHT22Driver) GetReading() *sensor.Reading {
	return d.lastReading
}

func createReadingFromBits(bits []int) (reading *sensor.Reading, err error) {
	var sum8 uint8
	var sumTotal uint8
	var checkSum uint8
	var i int
	var humidityInt int
	var temperatureInt int

	// get humidityInt value
	for i = 0; i < 16; i++ {
		humidityInt = humidityInt << 1
		humidityInt += bits[i]
		// sum 8 bits for checkSum
		sum8 = sum8 << 1
		sum8 += uint8(bits[i])
		if i == 7 || i == 15 {
			// got 8 bits, add to sumTotal for checkSum
			sumTotal += sum8
			sum8 = 0
		}
	}

	// get temperatureInt value
	for i = 16; i < 32; i++ {
		temperatureInt = temperatureInt << 1
		temperatureInt += bits[i]
		// sum 8 bits for checkSum
		sum8 = sum8 << 1
		sum8 += uint8(bits[i])
		if i == 23 || i == 31 {
			// got 8 bits, add to sumTotal for checkSum
			sumTotal += sum8
			sum8 = 0
		}
	}
	// if high 16 bit is set, value is negtive
	// 1000000000000000 = 0x8000
	if (temperatureInt & 0x8000) > 0 {
		// flip bits 16 and lower to get negtive number for int
		// 1111111111111111 = 0xffff
		temperatureInt |= ^0xffff
	}

	// get checkSum value
	for i = 32; i < 40; i++ {
		checkSum = checkSum << 1
		checkSum += uint8(bits[i])
	}

	// humidity is between 0 % to 100 %
	if humidityInt < 0 || humidityInt > 1000 {
		err = fmt.Errorf("bad data - humidity: %v", humidityInt)
		return
	}
	// temperature between -40 C to 80 C
	if temperatureInt < -400 || temperatureInt > 800 {
		err = fmt.Errorf("bad data - temperature: %v", temperatureInt)
		return
	}
	// check checkSum
	if checkSum != sumTotal {
		err = fmt.Errorf("bad data - check sum fail")
	}

	return &sensor.Reading{Temperature: float64(temperatureInt) / 10.0, Humidity: float64(humidityInt) / 10.0}, nil
}
