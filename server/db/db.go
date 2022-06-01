package database

import (
	"database/sql"
	"fmt"
	"log"

	migrate "github.com/rubenv/sql-migrate"
)

//MAIN DATA BASE CONNECTION FOR GLOBAL USE
var DBC *sql.DB

func Database() {
	DBC = OpenDB()

	Migrations(DBC)

	fmt.Println("Connection to database establised")
}

func Migrations(db *sql.DB) {

	//TAKES `../database/migrations` DIR AS A SOURCE OF MIGRATIONS
	migrations := &migrate.FileMigrationSource{
		Dir: "../database/migrations",
	}

	//EXECUTES THE MIGRATION (db conn pool, name, source, action)
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		fmt.Println("Error applying migrations: db.go")
	}
	fmt.Printf("Applied %d migrations to database.db!", n)
	fmt.Println()
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
