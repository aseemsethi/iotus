package utils

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type gw struct {
	CustomerId int    `json:"custid"`
	GwId       int    `json:"gwid"`
	Type       string `json:"type"`
	Location   string `json:"location"`
	IP         string `json:"ip"`
}

var gw1 gw

/*
 * The GW should send a publish messahe to gurupada/gw/add with the following body
 * {
	"custid"     :  100,
	"gwid"     : 100,
	"type"     : "esp32",
	"location" : "bangalore",
	"ip"       : "1.1.1.1"
 * }
*/
// Sample output of program -
// GW JSON recvd:::: {100 100 esp32 bangalore 1.1.1.1}
var gwMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd GW Control msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &gw1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n GW JSON recvd:::: %v", gw1)
	// Save this GW in a DB now
}
