package main

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
)

// InitializeProgram initializes the first custom program
func InitializeProgram(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	//Credentials
	creds := &CustomProgram{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	query := "insert into customprograms values ($1,$2,$3)"
	if _, err = db.Query(query,
		creds.Username,
		creds.ProgramDict,
		pq.Array(creds.WorkoutDays)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateCustomProgram updates the db with a new jsonified program
func UpdateCustomProgram() {
	//Replace program in db
}
