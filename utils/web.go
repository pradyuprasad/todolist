package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var store *sessions.CookieStore

func Init() {

	// load the env file

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// start
	store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

}

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
	fmt.Println("REACHED NEWTODO")
	http.ServeFile(w, r, "static/protected/newtodo.html")
}

func AuthRequired(next http.Handler) http.Handler {
	fmt.Println("running the auth function")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("unable to authenticate") // debug statement
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		fmt.Println("authentication done") // debug statement
		next.ServeHTTP(w, r)
	})
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
		//fmt.Println("this is being run") // debug statement
		var storedPassword string
		if err := results.Scan(&storedPassword); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Stored Password:", storedPassword)

		if storedPassword == password {
			session, err := store.Get(r, "user-session")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			session.Values["authenticated"] = true
			session.Values["username"] = username

			if err := session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("sending user to todo")
			http.Redirect(w, r, "/newtodo", http.StatusSeeOther)

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
