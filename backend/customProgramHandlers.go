package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

// InitializeProgram initializes the first custom program
func InitializeProgram(w http.ResponseWriter, r *http.Request) {
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
	program.StartDate = creds.StartDate

	//DB Query
	query := "insert into customprograms values ($1,$2,$3,$4)"

	if _, err = db.Query(query,
		program.Username,
		string(program.ProgramDict),
		pq.Array(program.WorkoutDays),
		string(program.StartDate)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateCustomProgram updates the program of the user
func UpdateCustomProgram(w http.ResponseWriter, r *http.Request) {
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
	program.StartDate = creds.StartDate

	//DB Query 1
	query := fmt.Sprintf("UPDATE customprograms SET programdict= '%s' WHERE username = '%s';", program.ProgramDict, program.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//DB Query 2
	justString := "{" + strings.Join(program.WorkoutDays, ",") + "}"
	query2 := fmt.Sprintf("UPDATE customprograms SET workoutdays= '%s' WHERE username = '%s';", justString, program.Username)
	if _, err = db.Query(query2); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//DB Query 3
	query3 := fmt.Sprintf("UPDATE customprograms SET startdate = '%s' WHERE username = '%s';", program.StartDate, program.Username)
	if _, err = db.Query(query3); err != nil {
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

// GetCustomProgram grabs the users custom program
func GetCustomProgram(w http.ResponseWriter, r *http.Request) {
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
	var program CustomProgramHelper
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
	err = row.Scan(&TempProgram.Username, &TempProgram.ProgramDict, pq.Array(&TempProgram.WorkoutDays), &TempProgram.StartDate)

	//Switch to helper struct
	program.Username = TempProgram.Username
	program.ProgramDict = []byte(TempProgram.ProgramDict)
	program.WorkoutDays = TempProgram.WorkoutDays
	program.StartDate = TempProgram.StartDate

	//Write Response
	ret, err := json.Marshal(program)
	w.Write(ret)
}
