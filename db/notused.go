package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

/*
mysql> CREATE TABLE customer (cid INT UNSIGNED NOT NULL PRIMARY KEY,
name VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, address VARCHAR(20) NOT NULL);
mysql> CREATE TABLE gw (gwid INT UNSIGNED NOT NULL PRIMARY KEY,
typegw VARCHAR(20) NOT NULL, location VARCHAR(20) NOT NULL, ip VARCHAR(20) NOT NULL);
mysql> CREATE TABLE sensor (cid INT UNSIGNED NOT NULL, gwid INT UNSIGNED NOT NULL,
sensorid INT UNSIGNED NOT NULL, type VARCHAR(20) NOT NULL);

mysql> CREATE USER 'root'@'%' IDENTIFIED BY 'PASSWORD';
mysql> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
mysql> FLUSH PRIVILEGES;
*/

func Db_init_notused() {
	fmt.Println("Connect to MySQL IOT DB")
	fmt.Println("Reading Environment Variable")
	// Set via bash shell - export DATABASE_PASS=xxxx
	var databasePass string
	databasePass = os.Getenv("DATABASE_PASS")
	fmt.Printf("Database Password: %s\n", databasePass)
	connString := "root:" + databasePass + "@tcp(15.206.73.249:3306)/IOT"

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
}

func Db_gw_add_notused(gwid int, typegw string, location string, ip string) {
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
}

func viewCustomersDB() {
	fmt.Println("viewCustomersDB called..")
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
		fmt.Printf("%v", customer)
	}
}
