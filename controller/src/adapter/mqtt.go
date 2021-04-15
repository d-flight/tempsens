package adapter

import (
	"fmt"
	"tempsens/data"
	"tempsens/mqtt"
	"time"

	mqtt_ "github.com/eclipse/paho.mqtt.golang"
)

const connectionRetryAttempts = 10
const connectionRetryAfter = 1 * time.Second

const qosAtMostOnce = 0
const qosAtLeastOnce = 1
const qosExactlyOnce = 2

const reportChannel = "tempsens/report"
const controlChannel = "tempsens/control"

type Mqtt struct {
	client mqtt_.Client

	// latest reading, desired temperature, current heating state
	OnNewReport func(*data.Reading, data.Temperature, data.HeatingState)
}

func NewMqttAdapter(mqttConfig *mqtt_.ClientOptions) *Mqtt {
	return &Mqtt{client: mqtt_.NewClient(mqttConfig)}
}

func (m *Mqtt) Boot() {
	// setup connection
	e := retry(m.client.Connect(), connectionRetryAttempts, connectionRetryAfter, nil)

	if e != nil {
		panic(e)
	} else {
		fmt.Println("connected")
	}

	// subscribe to the tempsens topic
	e = retry(
		m.client.Subscribe(reportChannel, qosAtLeastOnce, m.onReportMessage),
		3,
		1*time.Second,
		func(e error) {
			fmt.Println(e)
		},
	)

	if e != nil {
		panic(e)
	} else {
		fmt.Printf("subscribed")
	}
}

func (m *Mqtt) SetDesiredTemperature(t data.Temperature) {
	message, e := mqtt.NewChangeDesiredTemperatureMessage(t).Serialize()
	if e != nil {
		panic(e)
	}

	m.client.Publish(
		controlChannel,
		qosAtLeastOnce,
		false,
		message,
	)
}

func (m *Mqtt) ToggleActive(t bool) {
	message, e := mqtt.NewToggleActiveMessage(t).Serialize()
	if e != nil {
		panic(e)
	}

	m.client.Publish(
		controlChannel,
		qosAtLeastOnce,
		false,
		message,
	)
}

func (m *Mqtt) onReportMessage(client mqtt_.Client, msg mqtt_.Message) {
	report, e := mqtt.DeserializeReportMessage(msg.Payload())

	if e != nil {
		fmt.Printf("Received invalid report message (%v): %v\n", e, msg.Payload())
		return
	}

	fmt.Printf("Recieved new report message: %v\n", report)

	m.OnNewReport(&report.Reading, report.Desired, report.HeatingState)
}

func wait(t mqtt_.Token) (e error) {
	t.Wait()
	return t.Error()
}

func retry(t mqtt_.Token, attempts int, delay time.Duration, callback func(error)) (e error) {
	for i := 0; i < attempts; i++ {
		if wait(t) == nil {
			break
		} else {
			if callback != nil {
				callback(t.Error())
			}
			time.Sleep(delay)
		}
	}

	return t.Error()
}
