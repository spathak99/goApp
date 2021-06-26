
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

// UpdateWeights updates the goal weight and current body weight of the user
func UpdateWeights(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := &types.Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Update Bodyweight
	query := fmt.Sprintf("UPDATE users SET bodyweight = '%f' WHERE username = '%s';", creds.Bodyweight, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update Goalweight
	query = fmt.Sprintf("UPDATE users SET goalweight = '%f' WHERE username = '%s';", creds.GoalWeight, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}