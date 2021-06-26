package handlers

import (
	"encoding/json"
	"net/http"	
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)



//GetUsers gets all users except the current user
func GetUsers(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)
	
	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query and response
	var users []string
	sqlRaw := fmt.Sprintf(`select username from users`)
	rows, err := db.DB.Query(sqlRaw)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	ret, err := json.Marshal(users)
	w.Write(ret)
}


// GetUserData grabs profile struct data for the given user
func GetUserData(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)
	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Grabs user data from db
	var user types.Profile
	username := creds["username"].(string)
	row := db.DB.QueryRow("select * from users where username=$1", username)
	err = row.Scan(&user.Username, &user.Password, &user.Description,
		&user.GoalWeight, &user.Bodyweight,
		&user.CalorieGoal, &user.CaloriesLeft,
		pq.Array(&user.Followers), pq.Array(&user.Following),
		&user.Program, &user.Name)

	//Write Response
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(user)
	w.Write(ret)
}
