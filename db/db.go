package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
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
	Cid      int    `json:"cid"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"adress"`
	Gw       Gateway
}

type Gateway struct {
	GwId     int    `json:"gwid"`
	TypeGw   string `json:"type"`
	Location string `json:"location"`
	IP       string `json:"ip"`
}

var customers []Customer

var dbg *sql.DB

/*
mysql> CREATE USER 'root'@'%' IDENTIFIED BY 'PASSWORD';
mysql> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
mysql> FLUSH PRIVILEGES;
*/

func Db_init() {
	fmt.Println("Connect to MySQL IOT DB")
	fmt.Println("Reading Environment Variable")
	// Set via bash shell - export DATABASE_PASS=xxxx
	var databasePass string
	databasePass = os.Getenv("DATABASE_PASS")
	fmt.Printf("Database Password: %s\n", databasePass)
	connString := "root:" + databasePass + "@tcp(15.206.73.249:3306)/IOT"
	//fmt.Println("Conn string: ", connString)

	db, err := sql.Open("mysql", connString)
	dbg = db
	//defer db.Close()

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connection succeeded...")

	// read in the JSON file per customer - this is static file
	// for now just give the commands here....but this should be a JSON
	// file read eventually
	Db_customer_add(100, "PMOA", "Bangalore", "Varthur Rd")
	Db_gw_add(100, 10010, "Temperature", "Bangalore", "0.0.0.0")
	viewCustomers()
}

func viewCustomers() {
	fmt.Printf("\nCustomers InMem DB: \n")
	for _, v := range customers {
		fmt.Printf("%v", v)
	}
}

// Not used as of now....
// This function basically ensres that all data structures are in order for all
// customers and their gw/sensors after every update
func updateCustomersInMem() {
	fmt.Println("updateCustomersInMem called..")
	results, err := dbg.Query("SELECT * FROM customer")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var customer Customer
		err = results.Scan(&customer.Cid, &customer.Name, &customer.Location, &customer.Address)
		if err != nil {
			panic(err.Error())
		}
		// and then print out the tag's Name attribute
		fmt.Printf("%v", customer)
		customers = append(customers, customer)
	}
	fmt.Printf("\n%d Customers DataBase: \n%v", len(customers), customers)
}

func Db_customer_add(cid int, name string, location string, address string) {
	fmt.Println("Adding customer row...", cid, name, location, address)
	insert, err := dbg.Query("INSERT INTO customer VALUES ( ?,?,?,?)", cid, name, location, address)
	if err != nil {
		fmt.Println("DB_customer_add: ", err.Error())
		if strings.Contains(err.Error(), "Duplicate") {
			fmt.Println("Duplicate customer entry - ignore")
		}
	} else {
		defer insert.Close()
	}
	customers = append(customers, Customer{
		Cid: cid, Name: name, Location: location, Address: address,
		Gw: Gateway{GwId: 0, TypeGw: "", Location: "", IP: "0.0.0.0"}})
}

// cid will be non null when called at bootup time reading JSON file
// cid will be null when called from MQTT
func Db_gw_add(cid int, gwid int, typegw string, location string, ip string) {
	fmt.Println("Adding gw row..")
	insert, err := dbg.Query("INSERT INTO gw VALUES (?,?,?, ?)",
		gwid, typegw, location, ip)
	if err != nil {
		fmt.Println("DB_gw_add: ", err.Error())
		if strings.Contains(err.Error(), "Duplicate") {
			fmt.Println("Duplicate gw entry - ignore")
		}
	} else {
		defer insert.Close()
	}

	// Now update the InMem DB
	for i, v := range customers {
		if v.Cid == cid || v.Gw.GwId == gwid {
			fmt.Printf("\n GW %d added to customer %d", gwid, v.Cid)
			customers[i].Gw.GwId = gwid
			customers[i].Gw.TypeGw = typegw
			customers[i].Gw.Location = location
			customers[i].Gw.IP = ip
			viewCustomers()
			return
		}
	}
	fmt.Printf("\n GW %d not attached to any customer", gwid)
}
