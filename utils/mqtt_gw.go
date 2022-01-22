package utils

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var gwMqttRcv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("\n Recvd GW Control msg..")
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
