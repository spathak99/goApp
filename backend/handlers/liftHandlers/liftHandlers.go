package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)


// InitializeLifts initializes the first set of lifts entered by the user
func InitializeLifts(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := types.UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := "insert into userlifts values ($1,$2)"

	if _, err = db.DB.Query(query,
		lifts.Username,
		string(lifts.Lifts)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateLifts updates the max lifts of the user
func UpdateLifts(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := types.UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := fmt.Sprintf("UPDATE userlifts SET lifts= '%s' WHERE username = '%s';", lifts.Lifts, lifts.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Write Response
	w.Write([]byte(`{"response":"Succesfully updated program"}`))
}





