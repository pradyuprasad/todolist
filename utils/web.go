package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	db, err := DBopen()
	if err != nil {
		http.Error(w, "Error opening Database", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	var username = r.FormValue("username")

	if !ValidateLoginUsername(username) {
		// Send a response with JavaScript to show the popup
		invalidUsernameJS := `
            <script>
                alert('Invalid username');
            </script>
        `
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, invalidUsernameJS)
		http.ServeFile(w, r, "static/login.html")
		return

	}

	results, err := db.Query("select password from users where username = ?", username)

	if err != nil {
		http.Error(w, "Error with Database", http.StatusInternalServerError)
		return
	}

	defer results.Close()

	var password = r.FormValue("password")
	fmt.Println(username, password)

	if results.Next() {
		fmt.Println("this is being run")
		var storedPassword string
		if err := results.Scan(&storedPassword); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		fmt.Println("Stored Password:", storedPassword)

		if storedPassword == password {
			http.Redirect(w, r, "/static/loggedin.html", http.StatusFound)

		} else {

			invalidPasswordJS := `
            <script>
                alert('Invalid Password');
            </script>
        `
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, invalidPasswordJS)
			http.ServeFile(w, r, "static/login.html")

		}
	} else {

		invalidUsernameJS := `
            <script>
                alert('Invalid username');
            </script>
        `
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, invalidUsernameJS)
		http.ServeFile(w, r, "static/login.html")

	}

}

func LoginGET(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "static/login.html")

}

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

func CountRows(rows *sql.Rows) (int, error) {
	var count = 0
	for rows.Next() {

		count++

	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func ValidateLoginUsername(username string) bool {
	trimmedUsername := strings.TrimSpace(username)

	if trimmedUsername == "" {
		return false
	}

	return true

}
