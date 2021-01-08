package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// FuzzySearch does a fuzzy search of the name of the user
func FuzzySearch(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	//Creds
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Get Seperate Names
	var FullName []string
	var FirstName string
	var LastName string

	FullName = strings.Split(creds.Name, " ")
	FirstName = FullName[0]
	LastName = FullName[1]

	//DB Query
	var usernames []string
	rows, err := db.Query("select username from users where name like '%%s' and name like '%s%'", FirstName, LastName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		usernames = append(usernames, username)
	}

	var names []string
	rows, err = db.Query("select name from users where name like '%%s' and name like '%s%'", FirstName, LastName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		names = append(names, name)
	}
}
