package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var CONFIG_BROKER_URL string = "tcp://52.66.70.168:1883"
var c mqtt.Client = nil

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// ALL MQTT Topics uplink/downlink are described here
func Mqtt_set_routing() {
	// ***** UPLINK ***** both gw and sensor add messages to be sent every 15 min by gw
	// This acts as a heartbeat as well as to rebuild the cloud DB in the evernt of a
	// restart.
	// Update GW Info - Control packets receive handler setup
	if token := c.Subscribe("gurupada/gw/add", 0, gwMqttRcv); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	// Update SENSORS Info - Control packets receive handler setup
	if token := c.Subscribe("gurupada/sensor/add", 0, sensorMqttRcv); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	// Telemetry data - Data Path - All GWs will send Telemetry data to gurupada/data/<custid>
	if token := c.Subscribe("gurupada/data/#", 0, telemetryDataRecv); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// ***** DOWNLINK *****
	// The custid would be sent to the GWs using MQTT send channel - gurupada/<gwid>/downlink topic
	// All GWs need to listen on to this channel to get info
	// Data JSON Format:
	// { "custid":100, "cmdType": "http|mqtt", "cmdDst": "<url>|<mqttTopic>, "cmdVal": <data> }
}

func Mqtt_disconnect() {
	fmt.Printf("\n Mqtt_disconnect called")
	//time.Sleep(6 * time.Second)
	if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}

// https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#section-readme
func Mqtt_init() {
	fmt.Printf("\n Mqtt_init called")
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(CONFIG_BROKER_URL).SetClientID("gp_client")
	opts.SetUsername("draadmin")
	opts.SetPassword("DRAAdmin@123")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(30 * time.Second)

	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	// 	Subscribe(topic string, qos byte, callback MessageHandler) Token
	if token := c.Subscribe("gurupada/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Publish a test message
	// 	Publish(topic string, qos byte, retained bool, payload interface{}) Token
	token := c.Publish("gurupada/1", 0, false, "Gurupada IOT starting")
	token.Wait()
	fmt.Println("MQTT init completed...")
}
