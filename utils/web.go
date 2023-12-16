package utils

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func CreateUserPOST(w http.ResponseWriter, r *http.Request) {
	db, err := DBopen()

	if err != nil {
		http.Error(w, "Error opening Database", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	var username = r.FormValue("username")
	var password = r.FormValue("password")
	fmt.Println(username, password)

	_, err = db.Exec("INSERT into users (username, password) VALUES (?, ?)", username, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error creating User", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/static/createduser.html", http.StatusFound)

}

func CreateUserGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/createuser.html")
}

func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello, Go!")

}

func Serverrun(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
