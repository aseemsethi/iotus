package utils

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Telemerty struct {
	GwId     string `json:"gwid"`
	SensorId string `json:"sensorid"`
	Data     string `json:"data"`
}

var t1 Telemerty

// Telemetry data - Data Path - All GWs will send Telemetry data to gurupada/data/<custid>
var telemetryDataRecv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd Add GW Control Telemerty msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &t1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n GW JSON recvd:::: %v", t1)
	// Save to customer specific file
}
