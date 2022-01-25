package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

/* MY SQL Tables
mysql> CREATE TABLE customer (cid INT UNSIGNED NOT NULL,
name VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, address VARCHAR(20) NOT NULL);
mysql> CREATE TABLE gw (cid INT UNSIGNED NOT NULL, gwid INT UNSIGNED NOT NULL,
typegw VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, ip VARCHAR(20) NOT NULL);
mysql> CREATE TABLE sensor (cid INT UNSIGNED NOT NULL, gwid INT UNSIGNED NOT NULL,
sensorid INT UNSIGNED NOT NULL, type VARCHAR(20) NOT NULL);
*/

type Customer struct {
	Cid      int    `json:"cid"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"adress"`
}

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

	// read in the JSON file for customer tables- this is static file
	Db_customer_add(100, "PMOA", "Bangalore", "Varthur Rd")

	results, err := db.Query("SELECT * FROM customer")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var customer Customer
		err = results.Scan(&customer.Cid, &customer.Name, &customer.Location, &customer.Address)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Printf("%v", customer)
	}
}

func Db_customer_add(cid int, name string, location string, address string) {
	fmt.Println("Adding customer row...", cid, name, location, address)
	insert, err := dbg.Query("INSERT INTO customer VALUES ( ?,?,?,?)", cid, name, location, address)
	//insert, err := dbg.Query("INSERT INTO customer VALUES ( 100, 'Aseem', 'Blr', 'Varthur')")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer insert.Close()
}

func Db_gw_add(cid int, gwid int, typegw string, location string, ip string) {
	fmt.Println("Adding gw row...", cid, gwid, typegw, location, ip)
	insert, err := dbg.Query("INSERT INTO gw VALUES ( ?,?,?,?, ?)", cid, gwid, typegw, location, ip)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer insert.Close()
}
