package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-password"))

func NewTodoPOST(w http.ResponseWriter, r *http.Request) {

	db, err := DBopen()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()

	fmt.Println("db is", db)

	fmt.Println(r)

}

func NewTodoGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("Reached the 'newtodo' route. Requested URL: %s\n", r.URL.Path)
	fmt.Println("REACHED NEWTODO HAHAHAHAHAHA")
	http.ServeFile(w, r, "static/protected/newtodo.html")
}

func LoggedInGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/protected/loggedin.html")
}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	db, err := DBopen()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		// reload the login page again
		http.ServeFile(w, r, "static/login.html")
		return

	}
	// search for password
	results, err := db.Query("select password from users where username = ?", username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer results.Close()

	var password = r.FormValue("password")
	fmt.Println(username, password)

	//logic is basically that if results.Next() exists then it means that results is not zero length. Otherwise it is zero length
	if results.Next() {
		fmt.Println("this is being run")
		var storedPassword string
		if err := results.Scan(&storedPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Stored Password:", storedPassword)

		if storedPassword == password {

			session, err := store.Get(r, "user-name")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// something
			session.Values["authenticated"] = true
			session.Values["username"] = username
			err = session.Save(r, w) // ALWAYS SAVE THE SESSION BEFORE DOING ANYTHING ELSE
			fmt.Println("the saved session was", session)
			if err != nil {

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}
			http.Redirect(w, r, "/newtodo", http.StatusSeeOther)

		}

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

	invalidUsernameJS := `
            <script>
                alert('Invalid username');
            </script>
        `
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, invalidUsernameJS)
	http.ServeFile(w, r, "static/login.html")

}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "user-name")

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1

	session.Save(r, w)
	fmt.Fprintln(w, "Logged out")
	fmt.Println("Logged out")

}

func LoginGET(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "static/public/login.html")

}

func CreateUserPOST(w http.ResponseWriter, r *http.Request) {
	db, err := DBopen()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	var username = r.FormValue("username")
	var password = r.FormValue("password")
	fmt.Println(username, password)

	_, err = db.Exec("INSERT into users (username, password) VALUES (?, ?)", username, password)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	CreateUserGET(w, r)

}

func CreatedUserGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/public/createduser.html")

}

func CreateUserGET(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/public/createuser.html")
}

func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Hello, Go!")

}

func Serverrun(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/public")))) // forgot how this works need to ask GPT

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func ValidateLoginUsername(username string) bool {
	// returns true for valid and false for invalid

	// remove spaces using strings library
	trimmedUsername := strings.TrimSpace(username)
	// if the cut down string is empty then it is invalid
	if trimmedUsername == "" {
		return false
	}

	return true

}
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-name")
		fmt.Println("the session variable is", session)

		// Check if user is authenticated

		auth, ok := session.Values["authenticated"].(bool)
		if !ok {

			fmt.Println("not OK!")

			http.Error(w, "not OK!", http.StatusForbidden)

		} else if !auth {

			fmt.Println("not auth!")

			http.Error(w, "not auth", http.StatusForbidden)

		} else {

			fmt.Println("authentication done") // debug statement
			next.ServeHTTP(w, r)

		}

		/*if auth, ok := session.Values["authenticated"].(bool); !ok || !auth == "" {
			fmt.Println("unable to authenticate") // debug statement
			http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}*/

	})
}
