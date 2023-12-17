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
	router.HandleFunc("/login", utils.LoginGET).Methods("GET")
	router.HandleFunc("/login", utils.LoginPOST).Methods("POST")
	router.HandleFunc("/loggedin", utils.LoginPOST).Methods("GET")
	router.HandleFunc("/logout", utils.LogoutHandle).Methods("GET")

	SubRouter := router.PathPrefix("/").Subrouter()
	SubRouter.HandleFunc("/newtodo", utils.NewTodoGET)
	SubRouter.Use(utils.AuthRequired)
	utils.Serverrun(router) // ALWAYS RUN THIS AS THE LAST THING IN THE FILE OR ELSE EVERYTHING AFTER IT WON'T RUN
}
