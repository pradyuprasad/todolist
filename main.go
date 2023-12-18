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
	router.HandleFunc("/loggedin", utils.LoginPOST).Methods("GET")
	router.HandleFunc("/logout", utils.LogoutHandle).Methods("GET")

	SubRouter := router.PathPrefix("/").Subrouter()
	SubRouter.Use(utils.AuthRequired)
	SubRouter.HandleFunc("/newtodo", utils.NewTodoGET).Methods("GET")
	SubRouter.HandleFunc("/newtodo", utils.NewTodoPOST).Methods("POST")
	SubRouter.HandleFunc("/mytodos", utils.MyTodosGET).Methods("GET")
	SubRouter.HandleFunc("/test", utils.TestGET).Methods("GET")
	SubRouter.HandleFunc("/test-api", utils.TestAPI).Methods("GET")
	SubRouter.HandleFunc("/get_todos", utils.GetTodosAPI).Methods("GET")
	SubRouter.HandleFunc("/delete_todo/{todoID}", utils.DeleteTODOS).Methods("DELETE")

	LoginRouter := router.PathPrefix("/").Subrouter()
	LoginRouter.Use(utils.NotLoggedin)
	LoginRouter.HandleFunc("/login", utils.LoginGET).Methods("GET")
	LoginRouter.HandleFunc("/login", utils.LoginPOST).Methods("POST")
	utils.Serverrun(router) // ALWAYS RUN THIS AS THE LAST THING IN THE FILE OR ELSE EVERYTHING AFTER IT WON'T RUN
}
