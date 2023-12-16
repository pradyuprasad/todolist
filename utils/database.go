package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func DBopen() (*sql.DB, error) {
	// the output type is a pointer because sql.Open returns a pointer to the object

	// loads the goenv package
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil, err
	}

	//loads the username and password from the .env file
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := "todolist"

	// formates the dbURI string
	dbURI := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", username, password, dbName)

	// opens the SQL server
	// note that sql.Open returns a pointer to the actual db object
	db, err := sql.Open("mysql", dbURI)
	if err != nil {

		fmt.Println("Error:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println(err)
		db.Close()
		return nil, err

	}

	// returns the db object
	return db, nil

}

func CreateUsers(db *sql.DB) error {
	var CreateUsersQuery = `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) UNIQUE,
		password VARCHAR(255)
	);`

	_, err := db.Exec(CreateUsersQuery)

	if err != nil {
		return err
	}

	return nil

}

func ShowTables(db *sql.DB) error {

	rows, err := db.Query("SHOW TABLES;")

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var tablename string

		if err := rows.Scan(&tablename); err != nil {
			return err
		}

		fmt.Println(tablename)

	}

	return nil

}

func DeleteUser(db *sql.DB) error {

	var DropUsers = "DROP TABLE IF EXISTS users"

	_, err := db.Exec(DropUsers)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
