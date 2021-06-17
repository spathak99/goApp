package fuzzySearch

import (
	"encoding/json"
	"net/http"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
	"github.com/lib/pq"
)

// FuzzySearch does a fuzzy search of the name of the user
func FuzzySearch(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Creds
	creds := &types.SearchInfo{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	var users []types.Profile
	rows, err := db.DB.Query("select * from users where name like '" + creds.Query + "%'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user types.Profile
		err = rows.Scan(&user.Username, &user.Password, &user.Description,
			&user.GoalWeight, &user.Bodyweight,
			&user.CalorieGoal, &user.CaloriesLeft,
			pq.Array(&user.Followers), pq.Array(&user.Following),
			&user.Program, &user.Name)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	//Response
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(users)
	w.Write(ret)
}
