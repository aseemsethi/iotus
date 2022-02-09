package utils

import (
	"encoding/json"
	"fmt"
	db "github.com/aseemsethi/iotus/db"
	//sched "github.com/aseemsethi/iotus/sched"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
	"time"
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
	cid, sensorType, sensorName := db.Db_telemetry_update(t1)
	if cid == 0 {
		fmt.Printf("\n No customer fond for this MQTT telemetry update")
		return
	}

	// Send this msg to the Android App waiting on
	// gurupada/<custid>/<sensorType>/<sensorName>
	sendTopic := fmt.Sprintf("gurupada/%s/%s/%s",
		strconv.Itoa(cid), sensorType, sensorName)
	fmt.Printf("\nMQTT Assist: Send to %s, msg:%s", sendTopic, msg.Payload())
	token := c.Publish(sendTopic, 0, false, msg.Payload())
	token.Wait()

	//Also check for Alarm condition here, send mqtt alarm to Android app if yes
	checkAlarm(cid, t1)
}

func checkTempAlarm(t1 db.Telemerty, v2 db.SensorT) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	tempValue := strings.Split(t1.Data, ":")[1]
	a, _ := strconv.ParseFloat(tempValue, 64)
	b, _ := strconv.Atoi(v2.Compare)
	fmt.Printf("\n Triggers: Temp value: %f %f", a, float64(b))
	currentTime := time.Now().In(loc)
	tm := currentTime.Format("03:04 PM")
	tmNow, _ := time.ParseInLocation("03:04 PM", tm, loc)
	//fmt.Printf("\n Time Now: %v", tmNow) // 2022-02-06 21:00
	sensorStartTime, _ := time.ParseInLocation("03:04 PM", v2.TimeStart, loc)
	sensorEndTime, _ := time.ParseInLocation("03:04 PM", v2.TimeEnd, loc)
	//fmt.Printf("\nTriggers: SensorTimes: %v : %v", sensorStartTime, sensorEndTime)
	if tmNow.After(sensorStartTime) &&
		tmNow.Before(sensorEndTime) {
		fmt.Printf("\n Triggers: Time Alarm !!")
		if a > float64(b) {
			fmt.Printf("\n Triggers: Val Alarm !!")
		} else {
			fmt.Printf("\n Triggers: No Val Alarm !!")
		}
	}
}

func sendAlarm(cid int, msg string) {
	// Send this msg to the Android App waiting on gurupada/<custid>
	sendTopic := fmt.Sprintf("gurupada/%s/alarm", strconv.Itoa(cid))
	fmt.Printf("\nMQTT Assist Alarm: Send to %s, msg:%s", sendTopic, msg)
	token := c.Publish(sendTopic, 0, false, msg)
	token.Wait()
}

func checkAlarm(cid int, msg db.Telemerty) {
	for _, v := range db.T.Triggers {
		if v.Cid == cid {
			fmt.Printf("\n Triggers: Customer found...")
			for _, v1 := range v.Gw {
				if v1.GwId == msg.GwId {
					fmt.Printf("\n Triggers: GW found..." + v1.GwId)
					for _, v2 := range v1.Sensors {
						if v2.SensorId == msg.SensorId {
							fmt.Printf("\n Triggers: Sensor found..." + v2.SensorId)
							if v2.Type == "temperature" && v2.Trigger == ">" {
								if strings.Contains(strings.Split(msg.Data, ":")[0], "hex") {
									fmt.Printf("\niSensor HEX values recvd")
									// TBD - call checkTempAlarmHex()
								} else {
									checkTempAlarm(msg, v2)
								}
							} else if v2.Type == "door" && v2.Trigger == "=" {
								if strings.Contains(msg.Data, "Open") {
									fmt.Printf("\nDoor Opened %s", v2.SensorId)
									alarmMsg := fmt.Sprintf("%s:Open", v2.SensorId)
									sendAlarm(cid, alarmMsg)
								}
							}
						}
					}
				}
			}
		}
	}
}
