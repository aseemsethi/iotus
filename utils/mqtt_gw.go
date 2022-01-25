package utils

import (
	"encoding/json"
	"fmt"
	db "github.com/aseemsethi/iotus/db"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Gw struct {
	GwId     int    `json:"gwid"`
	Type     string `json:"type"`
	Location string `json:"location"`
	IP       string `json:"ip"`
}

var gw1 Gw

/*
 * The GW should send a publish message to gurupada/gw/add with the following body
 * for it to add itself to the DB tree
 * {
	"gwid"     : 10001,
	"type"     : "esp32",
	"location" : "bangalore",
	"ip"       : "1.1.1.1"
 * }
*/
// Sample output of program -
// GW JSON recvd:::: {100 100 esp32 bangalore 1.1.1.1}
var gwMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd Add GW Control msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &gw1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n GW JSON recvd:::: %v", gw1)
	// Save this GW in the DB now - 1st param cid is null, since GW does not send
	// us the customer id. It does know it.
	db.Db_gw_add(0, gw1.GwId, gw1.Type, gw1.Location, gw1.IP)
}

type sensor struct {
	CustomerId int    `json:"custid"`
	GwId       int    `json:"gwid"`
	SensorId   int    `json:"sensorid"`
	Type       string `json:"type"`
}

var sensor1 sensor

/*
 * The GW should send a publish message to gurupada/sensor/add with the following body
 * for it to add sensors to the DB tree under its GW struct
 * {
	"custid"   : 100,
	"gwid"     : 10001,
	"sensorid" : 1000101,
	"type"     : "temp"
* }
*/
// Sample output of program -
// GW JSON recvd:::: {100 10001 1000101 temp}
var sensorMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd Add Sensor Control msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &sensor1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n SENSOR JSON recvd:::: %v", sensor1)
	// Save this Sensor under the GW in the DB now
}
