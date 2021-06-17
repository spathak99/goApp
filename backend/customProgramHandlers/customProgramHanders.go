package customProgramHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// InitializeProgram initializes the first custom program
func InitializeProgram(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.CustomProgram{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		fmt.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	query := "insert into customprograms values ($1,$2)"

	if _, err = db.DB.Query(query,
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
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.CustomProgram{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		fmt.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query 1
	input := string(creds.ProgramList)
	query := fmt.Sprintf("UPDATE customprograms SET programlist= '%s' WHERE username = '%s';", input, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
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
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Grab program from database
	var program types.CustomProgramHelper
	username := creds["username"].(string)
	row := db.DB.QueryRow(`select * from customprograms where username=$1`, username)
	err = row.Scan(&program.Username, &program.ProgramList)

	//Write Responso
	ret, err := json.Marshal(program)
	w.Write(ret)
}