package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// InitializeLifts initializes the first set of lifts entered by the user
func InitializeLifts(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, name)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	//Credentials
	creds := &UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := "insert into userlifts values ($1,$2)"

	if _, err = db.Query(query,
		lifts.Username,
		string(lifts.Lifts)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateLifts updates the max lifts of the user
func UpdateLifts(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, name)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	//Credentials
	creds := &UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := fmt.Sprintf("UPDATE userlifts SET lifts= '%s' WHERE username = '%s';", lifts.Lifts, lifts.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Write Response
	ret := []byte(`{
		"response":"Succesfully updated program"
	}`)
	w.Write(ret)
}

//EstimateMax calculates the estimated one rep max
func EstimateMax(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, name)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	//Credentials
	creds := &Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO: Estimate one rep max
}
