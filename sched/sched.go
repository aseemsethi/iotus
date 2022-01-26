package sched

import (
	"fmt"
	//"github.com/aseemsethi/iotus/db"
)

/*
 This file needs to read in Trigger-Action rules of the following type
 This is only for "write" type of sensors - as noted in JSON file
 This needs to be read from mySQL DB or something similar. Can't be in JSON,
 since it could change based on user inputs very frequently.

 "custid": custid, "gwid": gwid, "sensorid": sensorid",
 "Trigger": start, "TimeStart": time, "TimeEnd": time,
 "Action" : on|off

*/

func SchedInit() {
	fmt.Printf("\nSchedInit called...")
}
