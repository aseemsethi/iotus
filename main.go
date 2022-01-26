package main

import (
	"fmt"
	"github.com/aseemsethi/iotus/db"
	"github.com/aseemsethi/iotus/httpG"
	"github.com/aseemsethi/iotus/mqtt"
	"net/http"
)

func main() {
	fmt.Printf("\nIOTUS Tool Starting..")

	utils.Mqtt_init()
	utils.Mqtt_set_routing()

	db.Db_init()

	http.HandleFunc("/api/customers", httpG.ApiCustomers)
	http.ListenAndServe(":8090", nil)
}

// Example o/p
//aseemsethi@yahoo.com:~/environment $ curl localhost:8090/api/customers
//User-Agent: curl/7.79.1
//Accept: */*
/*
[]db.Customer{
	db.Customer{Cid:100, Name:"PMOA", Location:"DG Room", Address:"Varthur Rd",
		Gw:[]
		db.Gateway{db.Gateway{GwId:10010, TypeGw:"ESP32", Location:"Bangalore", IP:"0.0.0.0",
			Sensors:[]db.Sensor{db.Sensor{SensorId:1001001, Type:"temperature", Protocol:""},
					db.Sensor{SensorId:1001002, Type:"onoff", Protocol:""}}},
		db.Gateway{GwId:10020, TypeGw:"ESP32", Location:"Bangalore", IP:"0.0.0.0",
			Sensors:[]db.Sensor(nil)}}},
	db.Customer{Cid:200, Name:"Prestige", Location:"DG Room", Address:"Varthur Rd",
		Gw:[2]db.Gateway{db.Gateway{GwId:20010, TypeGw:"ESP32", Location:"Bangalore", IP:"0.0.0.0",
			Sensors:[]db.Sensor(nil)}, db.Gateway{GwId:20020, TypeGw:"ESP32", Location:"Bangalore", IP:"0.0.0.0",
				Sensors:[]db.Sensor(nil)}}}}aseemsethi@yahoo.com:~/environment $
*/
