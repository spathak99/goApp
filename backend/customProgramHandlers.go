package main

import (
	"encoding/json"
	"io/ioutil"
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
	creds := &CustomProgramHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	program := CustomProgram{}
	program.Username = creds.Username
	program.ProgramDict = string(creds.ProgramDict)
	program.WorkoutDays = creds.WorkoutDays

	//DB Query
	query := "insert into customprograms values ($1,$2,$3)"
	if _, err = db.Query(query,
		program.Username,
		program.ProgramDict,
		pq.Array(program.WorkoutDays)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateCustomProgram updates the db with a new jsonified program
func UpdateCustomProgram() {
	//Replace program in db
}
