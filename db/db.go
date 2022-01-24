package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Customer struct {
	Cid      int    `json:"cid"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"address"`
}

/*
mysql> CREATE USER 'root'@'%' IDENTIFIED BY 'PASSWORD';
mysql> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;
mysql> FLUSH PRIVILEGES;
*/

func Db_init() {
	fmt.Println("Connect to MySQL IOT DB")
	fmt.Println("Reading Environment Variable")
	var databasePass string
	databasePass = os.Getenv("DATABASE_PASS")
	fmt.Printf("Database Password: %s\n", databasePass)
	connString := "root:" + databasePass + "@tcp(15.206.73.249:3306)/IOT"
	//fmt.Println("Conn string: ", connString)

	// Open up our database connection.
	db, err := sql.Open("mysql", connString)
	// defer the close till after the main function has finished executing
	defer db.Close()

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Println("Connection succeeded...")

	// Execute the query
	insert, err := db.Query("INSERT INTO customer VALUES ( 99, 'PMOA', 'Bangalore', 'Varthur')")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer insert.Close()

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
