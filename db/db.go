package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
)

/* MY SQL Tables
mysql> CREATE TABLE customer (cid INT UNSIGNED NOT NULL PRIMARY KEY,
name VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, address VARCHAR(20) NOT NULL);
mysql> CREATE TABLE gw (gwid INT UNSIGNED NOT NULL PRIMARY KEY,
typegw VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, ip VARCHAR(20) NOT NULL);
mysql> CREATE TABLE sensor (cid INT UNSIGNED NOT NULL, gwid INT UNSIGNED NOT NULL,
sensorid INT UNSIGNED NOT NULL, type VARCHAR(20) NOT NULL);
*/

type Customer struct {
	Cid      int        `json:"cid"`
	Name     string     `json:"name"`
	Location string     `json:"location"`
	Address  string     `json:"address"`
	Gw       [2]Gateway `json:"gateway"`
}

type Gateway struct {
	GwId     int    `json:"gwid"`
	TypeGw   string `json:"type"`
	Location string `json:"location"`
	IP       string `json:"ip"`
}

type Customers struct {
	Customers []Customer `json:"customers"`
}

var c Customers

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
	json.Unmarshal(byteValue, &c)
	for i, _ := range c.Customers {
		fmt.Printf("\n%+v", c.Customers[i])
	}
}

func Db_gw_add(gwid int, typegw string, location string, ip string) {
	fmt.Println("Updating gw row..")
	for i, v := range c.Customers {
		for i1, v1 := range v.Gw {
			if v1.GwId == gwid {
				fmt.Printf("\n GW %d updated in customer %d", gwid, v.Cid)
				c.Customers[i].Gw[i1].TypeGw = typegw
				c.Customers[i].Gw[i1].Location = location
				c.Customers[i].Gw[i1].IP = ip
				for j, _ := range c.Customers {
					fmt.Printf("\n%+v", c.Customers[j])
				}
				return
			}
		}
	}
	fmt.Printf("\n GW %d not updated in any customer row", gwid)
}
