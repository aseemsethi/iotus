package sched

import (
	"fmt"
	"time"
	//"github.com/aseemsethi/iotus/db"
)

/*
 This file needs to read in Trigger-Action rules of the following type
 This is only for "write" type of sensors - as noted in JSON file
 This needs to be read from mySQL DB or something similar. Can't be in JSON,
 since it could change based on user inputs very frequently.

Case:1 - one time tasks, based on HTTP API
 "custid": custid, "gwid": gwid, "sensorid": sensorid",
 "Trigger": Change State "on"|"off"
 "Action" : Send MQTT Msg to MQTT based Sensor

Case:2 - one time tasks, based on Trigger/Action Alarm
 "custid": custid, "gwid": gwid, "sensorid": sensorid",
 "Trigger": value > 50 (temp), value < 10% (water level), etc
 "Action" : Alarm

Case:3 - scheduled tasks
 "custid": custid, "gwid": gwid, "sensorid": sensorid",
 "Trigger": start, "TimeStart": time, "TimeEnd": time,
 "Action" : on|off
*/

func SchedInit() {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	tm := time.Now().In(loc)
	fmt.Printf("\nSchedInit called...%s", tm.Format("2006-01-02 15:04:05"))
}
