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

	// To the original message recvd from GW, add in the
	// sensorName and sensorType too.
	var m map[string]interface{}
	err = json.Unmarshal(msg.Payload(), &m)
	m["name"] = sensorName
	m["type"] = sensorType
	msgNew, err := json.Marshal(m)

	// Send this msg to the Android App waiting on
	// gurupada/<custid>/<sensorType>/
	sendTopic := fmt.Sprintf("gurupada/%s/%s",
		strconv.Itoa(cid), sensorType)
	fmt.Printf("\nMQTT Assist: Send to %s, msg:%s", sendTopic, msgNew) // msg.Payload())
	token := c.Publish(sendTopic, 0, false, msgNew)                    // msg.Payload())
	token.Wait()

	//Also check for Alarm condition here, send mqtt alarm to Android app if yes
	checkAlarm(cid, t1)
}

func checkTempAlarmHex(t1 db.Telemerty, v2 db.SensorT) bool {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	tempValue := strings.Split(t1.Data, ":")[1]
	a, _ := strconv.ParseUint(tempValue, 16, 32)
	a1 := uint32(a)
	fmt.Printf("\na: %d", a1)
	c := float32(a1)
	fmt.Printf("\nc: %f", c)
	b, _ := strconv.Atoi(v2.Compare)
	fmt.Printf("\n Triggers: Temp value: %f %f", c, float32(b))
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
		if c > float32(b) {
			fmt.Printf("\n Triggers: Val Alarm !!")
			return true
		} else {
			fmt.Printf("\n Triggers: No Val Alarm !!")
			return false
		}
	}
	return false
}

func checkTempAlarm(t1 db.Telemerty, v2 db.SensorT) bool {
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
			return true
		} else {
			fmt.Printf("\n Triggers: No Val Alarm !!")
			return false
		}
	}
	return false
}

func sendAlarm(cid int, msg string) {
	// Send this msg to the Android App waiting on gurupada/<custid>

	sendTopic := fmt.Sprintf("gurupada/%s/alarm", strconv.Itoa(cid))
	fmt.Printf("\nMQTT Assist Alarm: Send to %s, msg:%s", sendTopic, msg)
	token := c.Publish(sendTopic, 0, false, msg)
	token.Wait()
}

func checkAlarm(cid int, msg db.Telemerty) {
	found := false

	loc, _ := time.LoadLocation("Asia/Kolkata")
	currentTime := time.Now().In(loc)
	tm := currentTime.Format("03:04 PM")
	tmNow, _ := time.ParseInLocation("03:04 PM", tm, loc)
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
									fmt.Printf("\nSensor HEX values recvd")
									if checkTempAlarmHex(msg, v2) == true {
										alarmMsg := fmt.Sprintf("%s:High Temp", v2.Name)
										sendAlarm(cid, alarmMsg)
									}
								} else {
									if checkTempAlarm(msg, v2) == true {
										alarmMsg := fmt.Sprintf("%s:High Temp", v2.Name)
										sendAlarm(cid, alarmMsg)
									}
								}
							} else if v2.Type == "door" && v2.Trigger == "=" {
								// We are triggering on "open" only for now...disregarding json file
								if strings.Contains(msg.Data, "Open") {
									fmt.Printf("\nDoor Opened %s", v2.SensorId)
									alarmMsg := fmt.Sprintf("%s:Open", v2.Name)
									sensorStartTime, _ := time.ParseInLocation("03:04 PM", v2.TimeStart, loc)
									sensorEndTime, _ := time.ParseInLocation("03:04 PM", v2.TimeEnd, loc)
									if tmNow.After(sensorStartTime) &&
										tmNow.Before(sensorEndTime) {
										fmt.Printf("\nDoor Opened - send alarm - %s", v2.Name)
										sendAlarm(cid, alarmMsg)
									} else {
										fmt.Printf("\nDoor Opened - no alarm")
									}
								}
							}
							found = true
						} // end if sensor found
					} // end for loop
					if found == false {
						fmt.Printf("\nTriggers: Sensor not found...")
					}
				}
			}
		}
	}
}
