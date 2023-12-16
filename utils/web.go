package utils

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func LoginPOST(w http.ResponseWriter, r *http.Request) error {
	db, err := DBopen()

	if err != nil {
		http.Error(w, "Error opening Database", http.StatusInternalServerError)
		return err
	}

	r.ParseForm()
	var username = r.FormValue("username")
	var password = r.FormValue("password")

	_, err = db.Exec("INSERT into users (username, password) VALUES (?, ?)", username, password)

	if err != nil {
		http.Error(w, "Error creating User", http.StatusInternalServerError)
		return err
	}
	http.Redirect(w, r, "./static/createduser.html", http.StatusFound)
	return nil

}

func LoginGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login.html")
}

func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello, Go!")

}

func Serverrun(router *mux.Router) {
	// always keep these two together v
	/*fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)*/
	// always keep these 2 together ^ idk why time to find out

	err := http.ListenAndServe(":8000", router)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}