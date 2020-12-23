package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lib/pq"
)

// UpdateCustomProgram initializes the first custom program
func UpdateCustomProgram(w http.ResponseWriter, r *http.Request) {
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

// GetCustomProgram grabs the users custom program
func GetCustomProgram(w http.ResponseWriter, r *http.Request) {
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

	//Grab program from database
	var TempProgram CustomProgram
	row := db.QueryRow(`select * from customprograms where username=$1`, creds.Username)
	err = row.Scan(&TempProgram.Username, &TempProgram.ProgramDict, pq.Array(&TempProgram.WorkoutDays))

	//Switch to helper struct
	var program CustomProgramHelper
	program.Username = TempProgram.Username
	program.ProgramDict = []byte(TempProgram.ProgramDict)
	program.WorkoutDays = TempProgram.WorkoutDays

	//Write Response
	ret, err := json.Marshal(program)
	w.Write(ret)
}
