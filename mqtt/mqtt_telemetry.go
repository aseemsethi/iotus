package utils

import (
	"encoding/json"
	"fmt"
	db "github.com/aseemsethi/iotus/db"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
)

var t1 db.Telemerty

// Telemetry data - Data Path - All GWs will send Telemetry data to gurupada/data/<custid>
var telemetryDataRecv mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("\n Recvd Telemerty msg..")
	fmt.Printf("\nTOPIC: %s", msg.Topic())
	fmt.Printf("\nMSG: %s", msg.Payload())

	err := json.Unmarshal([]byte(msg.Payload()), &t1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n GW JSON recvd:::: %v", t1)
	// Save to customer specific file
	cid, sensorType := db.Db_telemetry_update(t1)

	// Send this msg to the Android App waiting on gurupada/<custid>
	sendTopic := fmt.Sprintf("gurupada/%s/%s", strconv.Itoa(cid), sensorType)
	fmt.Printf("\nMQTT Assist: Send to %s, msg:%s", sendTopic, msg.Payload())
	token := c.Publish(sendTopic, 0, false, msg.Payload())
	token.Wait()
}
