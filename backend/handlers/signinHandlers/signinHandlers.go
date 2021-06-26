package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
	"golang.org/x/crypto/bcrypt"
)


// Signin signs in the user and authenticates them
func Signin(w http.ResponseWriter, r *http.Request) {
	//Start Session
	session, _ := helpers.Store.Get(r, helpers.CookieName)



	//User authentication below
	creds := &types.Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		session.Values["authenticated"] = false
		session.Save(r, w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Grabs password
	result := db.DB.QueryRow("select password from users where username=$1", creds.Username)
	if err != nil {
		session.Values["authenticated"] = false
		session.Save(r, w)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Decode Creds
	storedCreds := &types.Profile{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		session.Values["authenticated"] = false
		session.Save(r, w)
		if err == sql.ErrNoRows {

			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Hash Password
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		session.Values["authenticated"] = false
		session.Save(r, w)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Verify Session
	session.Values["authenticated"] = true
	session.Save(r, w)
}

// Logout logs the user out and ends the session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := helpers.Store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

// Signup creates a new entry in the users table in the db
func  Signup(w http.ResponseWriter, r *http.Request) {

	//Decode Creds
	creds := &types.Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Logs new user
	query := "insert into users values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if _, err = db.DB.Query(query,
		creds.Username,
		string(hashedPassword),
		string(creds.Description),
		creds.GoalWeight,
		creds.Bodyweight,
		creds.CalorieGoal,
		creds.CaloriesLeft,
		pq.Array(creds.Followers),
		pq.Array(creds.Following),
		string(creds.Program),
		string(creds.Name)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
