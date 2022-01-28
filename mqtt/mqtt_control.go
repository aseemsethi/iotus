package utils

import (
	"encoding/json"
	"fmt"
	db "github.com/aseemsethi/iotus/db"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var gw1 db.Gateway

/*
 * The GW should send a publish message every 15 min to gurupada/gw/add
 * with the following body for it to update data about itself in the DB tree
 {
	"gwid"     : "10010",
	"type"     : "esp32",
	"location" : "bangalore",
	"ip"       : "1.1.1.1"
 }
*/
// GW JSON recvd:::: {10010 esp32 bangalore 1.1.1.1}
var gwMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd Add GW Control msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &gw1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n GW JSON recvd:::: %v", gw1)
	// Update additional info like IP etc, info recvd for this GW in the DB
	db.Db_gw_add(gw1.GwId, gw1.TypeGw, gw1.Location, gw1.IP)
}

var sensor1 db.Sensor

/*
 * The GW should send a publish message every 15 min to gurupada/sensor/add
 * with the following body to update sensors data in the DB tree under GW struct
 {
	"gwid"     : "10010",
	"sensorid" : "1001001",
	"type"     : "sonoff",
	"protocol" : "ble",
	"rw"       : "write"
}
*/
// GW JSON recvd:::: {10001 1000101 temp ble write}
var sensorMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd Add Sensor Control msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &sensor1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n SENSOR JSON recvd:::: %v", sensor1)
	// Update Sensor under the GW in the DB now
	db.Db_sensor_add(sensor1.GwId, sensor1.SensorId, sensor1.Type, sensor1.Protocol, sensor1.RW)
}
