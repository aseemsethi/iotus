package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
)

type Customer struct {
	Cid      int       `json:"cid"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Address  string    `json:"address"`
	Gw       []Gateway `json:"gateway"`
}

type Gateway struct {
	GwId     string   `json:"gwid"`
	TypeGw   string   `json:"type"`
	Location string   `json:"location"`
	IP       string   `json:"ip"`
	Sensors  []Sensor `json:"sensor"`
}

type Sensor struct {
	GwId     string `json:"gwid"`
	SensorId string `json:"sensorid"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	RW       string `json:"rw"`
}

type Customers struct {
	Customers []Customer `json:"customers"`
}

var C Customers

var dbg *sql.DB

func Db_init() {
	fmt.Println("Initialize DB")
	readCustomerFile()
	fmt.Println("...Done")
}

func readCustomerFile() {
	jsonFile, err := os.Open("cfg/Customer.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\nSuccessfully Opened Customer.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//fmt.Println("JSON File: ", byteValue)
	json.Unmarshal(byteValue, &C)
	for i, _ := range C.Customers {
		fmt.Printf("\n%+v", C.Customers[i])
	}
}

func Db_gw_add(gwid string, typegw string, location string, ip string) {
	fmt.Println("Updating gw row..")
	for i, v := range C.Customers {
		for i1, v1 := range v.Gw {
			if v1.GwId == gwid {
				fmt.Printf("\n GW %d updated in customer %d", gwid, v.Cid)
				C.Customers[i].Gw[i1].TypeGw = typegw
				C.Customers[i].Gw[i1].Location = location
				C.Customers[i].Gw[i1].IP = ip
				for j, _ := range C.Customers {
					fmt.Printf("\n%+v", C.Customers[j])
				}
				return
			}
		}
	}
	fmt.Printf("\n GW %d not updated in any customer row", gwid)
}

func Db_sensor_add(gwid string, sensorid string, typeSensor string, protocol string, rw string) {
	fmt.Println("Updating gw row..")
	for i, v := range C.Customers {
		for i1, v1 := range v.Gw {
			for i2, _ := range v1.Sensors {
				if v1.GwId == gwid {
					fmt.Printf("\n Sensor %d under GW %d updated in customer %d",
						sensorid, gwid, v.Cid)
					C.Customers[i].Gw[i1].Sensors[i2].Type = typeSensor
					C.Customers[i].Gw[i1].Sensors[i2].Protocol = protocol
					C.Customers[i].Gw[i1].Sensors[i2].RW = rw
					for j, _ := range C.Customers {
						fmt.Printf("\n%+v", C.Customers[j])
					}
					return
				}
			}
		}
	}
	fmt.Printf("\n GW %d not updated in any customer row", gwid)
}
