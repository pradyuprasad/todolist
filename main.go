package main

import (
	"fmt"

	"todolist/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	db, err := utils.DBopen()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", utils.Home).Methods("GET")
	router.HandleFunc("/createuser", utils.CreateUserGET).Methods("GET")
	router.HandleFunc("/createuser", utils.CreateUserPOST).Methods("POST")
	//router.HandleFunc("/login", )
	utils.Serverrun(router)

}
