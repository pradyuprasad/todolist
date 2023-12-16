package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	db, err := DBopen()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/login", LoginGET).Methods("GET")
	serverrun(router)

}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	db, err := DBopen()

	if err != nil {
		http.Error(w, "Error opening Database", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	var username = r.FormValue("username")
	var password = r.FormValue("password")

	_, err = db.Exec("INSERT into users (username, password) VALUES (?, ?)", username, password)

	if err != nil {
		http.Error(w, "Error creating User", http.StatusInternalServerError)
		return
	}

}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello, Go!")

}

func serverrun(router *mux.Router) {
	// always keep these two together v
	/*fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)*/
	// always keep these 2 together ^ idk why time to find out

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

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
