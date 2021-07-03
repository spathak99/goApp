package handlers

import (
	"encoding/json"
	"net/http"	
	"fmt"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// UpdateCalories updates the calorie goals and calorie counts for the user
func UpdateCalories(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := &types.Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	//Updates Calories left
	var currCals int
	row := db.DB.QueryRow("select caloriesleft from users where username=$1", creds.Username)
	err = row.Scan(&currCals)
	calsLeft := creds.CaloriesLeft
	if (int(creds.CaloriesLeft)) < 0 {
		calsLeft = 0
	}

	//updates calories left
	query := fmt.Sprintf("UPDATE users SET caloriesleft = '%f' WHERE username = '%s';",calsLeft, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//updates calorie goal
	query = fmt.Sprintf("UPDATE users SET caloriegoal = '%f' WHERE username = '%s';", creds.CalorieGoal, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Response
	w.Write([]byte(`{"response":"Succesfully updated calories"}`))
	return
}