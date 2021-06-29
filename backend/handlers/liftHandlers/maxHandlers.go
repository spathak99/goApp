package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"	
)

//GetUserMax gets the max of the user
func GetUserMax(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Gets lifts of user
	var lifts string
	username := creds["username"].(string)
	row := db.DB.QueryRow("select lifts from userlifts where username=$1", username)
	err = row.Scan(&lifts)

	//Write
	w.Write([]byte(lifts))
}

//EstimateMax calculates the estimated one rep max
func EstimateMax(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Calculation and response
	reps := float64(creds.Reps)
	weight := float64(creds.Weight)
	repMax := int((weight) / (1.0278 - 0.0278*((10.0-creds.RPE)+reps)))
	ret := strconv.Itoa(repMax)
	w.Write([]byte(ret))
}