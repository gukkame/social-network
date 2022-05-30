package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

//MAIN DATA BASE CONNECTION FOR GLOBAL USE
var DBC *sql.DB

func Database() {
	DBC = OpenDB()

	Initalize(DBC, "../database/dbtemplate.sql")

	fmt.Println("Connection to database establised")
}

//READS `dbtemplate.sql` AND CREATES TABLES IF THEY DONT EXIST
func Initalize(db *sql.DB, path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("(db.go) Unable to os.ReadFile due to:")
		log.Fatal(err)
	}
	//LOOPS THROUGH ALL THE SQL STATEMENTS IN `dbtemplate.sql`
	requests := strings.SplitAfter(string(file), ";")
	for _, request := range requests {
		db.Exec(request)
	}
}

//OPENING DATABASE CONNECTION
func OpenDB() *sql.DB {
	db, err := sql.Open("sqlite3", "../database/database.db")
	if err != nil {
		fmt.Println("(db.go) Unable to open db due to:")
		log.Fatal(err)
	}

	//USING PING TO ESTABLISH A NEW CONNECTION TO THE DATABASE
	err = db.Ping()
	if err != nil {
		fmt.Println("(db.go) Unable to ping due to:")
		log.Fatal(err)
	}

	// RETURNS DB CONNECTION POOl
	return db
}
