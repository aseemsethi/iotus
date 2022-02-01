package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Customer struct {
	Cid         int       `json:"cid"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Address     string    `json:"address"`
	Gw          []Gateway `json:"gateway"`
	LastUpdated time.Time
}

type Gateway struct {
	GwId     string `json:"gwid"`
	TypeGw   string `json:"type"`
	Location string
	IP       string   `json:"ip"`
	Sensors  []Sensor `json:"sensor"`
}

type Sensor struct {
	GwId     string `json:"gwid"`
	SensorId string `json:"sensorid"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	RW       string `json:"rw"`
	State    string // "unknown", "open", "close", etc
}

type Customers struct {
	Customers []Customer `json:"customers"`
}

var C Customers

type Telemerty struct {
	GwId     string `json:"gwid"`
	SensorId string `json:"sensorid"`
	Data     string `json:"data"`
}

var dbg *sql.DB

func Db_init() {
	fmt.Printf("\nInitialize DB")
	readCustomerFile()
	fmt.Println("...Done")
}

func readCustomerFile() {
	jsonFile, err := os.Open("cfg/Customer.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\nOpened Customer.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//fmt.Println("JSON File: ", byteValue)
	json.Unmarshal(byteValue, &C)
	for i, _ := range C.Customers {
		fmt.Printf("\n%+v", C.Customers[i])
	}
}

func Db_gw_add(gwid string, typegw string, ip string) int {
	fmt.Println("Updating gw row..")
	loc, _ := time.LoadLocation("Asia/Kolkata")

	for i, v := range C.Customers {
		for i1, v1 := range v.Gw {
			if v1.GwId == gwid {
				fmt.Printf("\n GW %d updated in customer %d", gwid, v.Cid)
				C.Customers[i].Gw[i1].TypeGw = typegw
				C.Customers[i].Gw[i1].IP = ip
				C.Customers[i].LastUpdated = time.Now().In(loc)
				for j, _ := range C.Customers {
					fmt.Printf("\n%+v", C.Customers[j])
				}
				return v.Cid
			}
		}
	}
	fmt.Printf("\n GW %d not updated in any customer row", gwid)
	return 0
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

func Db_telemetry_update(t Telemerty) (int, string) {
	fmt.Printf("\nUpdating telemerty data..")
	for i, v := range C.Customers {
		for i1, v1 := range v.Gw {
			if v1.GwId == t.GwId {
				for i2, v2 := range v1.Sensors {
					if v2.SensorId == t.SensorId {
						fmt.Printf("\n Sensor %d under GW %d updated in customer %d to - %s",
							v2.SensorId, v1.GwId, v.Cid, t.Data)
						C.Customers[i].Gw[i1].Sensors[i2].State = t.Data
						// Add Telemetry data line to customer file - customer-<cid>.stats
						// Line Format - cid, gwid, sensorid, data
						Db_telemetry_save(
							C.Customers[i].Cid, C.Customers[i].Gw[i1].GwId,
							C.Customers[i].Gw[i1].Sensors[i2].SensorId, t.Data)
						return C.Customers[i].Cid, C.Customers[i].Gw[i1].Sensors[i2].Type
					}
				}
			}
		}
	}
	fmt.Printf("\n GW %d not updated in any customer row", t.GwId)
	return 0, ""
}

func Db_telemetry_save(cid int, gwid string, sensorid string, data string) {
	fmt.Printf("\n DB telemerty save to file..")
	currentTime := time.Now()
	tm := currentTime.Format("2006-01-02 15:04:05")
	// If the file doesn't exist, create it, or append to the file
	filename := fmt.Sprintf("stats/customer-%s", strconv.Itoa(cid))
	fmt.Sprintf("\n Writing stats to file %s", filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.Write([]byte(tm + ":" + gwid + ":" +
		sensorid + ":" + data + "\n")); err != nil {
		fmt.Println(err)
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("..done")
}
