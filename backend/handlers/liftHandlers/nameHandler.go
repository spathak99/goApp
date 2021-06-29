package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"	
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
	"github.com/lib/pq"
)

//GetLiftNames gets all the types of lifts that have been logged
func GetLiftNames(w http.ResponseWriter, r *http.Request){

	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//DB Query
	var username = creds["username"].(string)
	var liftlog []string
	row := db.DB.QueryRow("select lifts from exerciselog where username=$1", username)
	err = row.Scan(pq.Array(&liftlog))

	//Get Names
	var temp []string
	for _, lift := range liftlog {
		srcjson := []byte(lift)
		var helper types.Lift
		err := json.Unmarshal(srcjson, &helper)
		if err != nil {
			panic(err)
		}
		temp = append(temp,helper.Name)
	}

	//Reponse
	var liftNames = helpers.Unique(temp)
	response, err := json.Marshal(liftNames)
	w.Write([]byte(fmt.Sprintf(`{"lifts": %s }`, response)))
}
