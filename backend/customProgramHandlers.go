package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// InitializeProgram initializes the first custom program
func InitializeProgram(w http.ResponseWriter, r *http.Request) {
	//Authentication
	authenticate(w,r)

	//Credentials
	creds := &CustomProgram{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		fmt.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	query := "insert into customprograms values ($1,$2)"

	if _, err = db.Query(query,
		creds.Username,
		string(creds.ProgramList),
	); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateCustomProgram updates the program of the user
func UpdateCustomProgram(w http.ResponseWriter, r *http.Request) {
	//Authentication
	authenticate(w,r)

	//Credentials
	creds := &CustomProgram{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		fmt.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query 1
	input := string(creds.ProgramList)
	print(input)
	query := fmt.Sprintf("UPDATE customprograms SET programlist= '%s' WHERE username = '%s';", input, creds.Username)
	if _, err = db.Query(query); err != nil {
		fmt.Printf("%s", err)
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
	authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Grab program from database
	var program CustomProgramHelper
	username := creds["username"].(string)
	row := db.QueryRow(`select * from customprograms where username=$1`, username)
	err = row.Scan(&program.Username, &program.ProgramList)

	//Write Responso
	ret, err := json.Marshal(program)
	w.Write(ret)
}
