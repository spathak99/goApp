package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"net/http"
	_ "github.com/lib/pq"    
	"github.com/gorilla/sessions"
)


/*
	Session keys
*/
var (
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)


/*
	User Field
*/
type Credentials struct {
	Password string `json:"password", db:"password"`
	Description string `json:"description",db:"description"`
	Username string `json:"username", db:"username"`
}


/*
	Test Session Method
*/
func Secret(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Print secret message
    print(w, "The cake is a lie!")
}


/*
	End session for suer
*/

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)
}


/*
	Signup a new user
*/
func Signup(w http.ResponseWriter, r *http.Request){	
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if _, err = db.Query("insert into users values ($1, $2,$3)", creds.Username, string(hashedPassword),string(creds.Description)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


/*
	Signin existing users
*/

func Signin(w http.ResponseWriter, r *http.Request){
	//Start Session
	session, _ := store.Get(r, "cookie-name")

	//User authentication below
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	result := db.QueryRow("select password from users where username=$1", creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	//Verify Session
	session.Values["authenticated"] = true
    session.Save(r, w)
}





