package main

import (
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
)

// FuzzySearch does a fuzzy search of the name of the user
func FuzzySearch(w http.ResponseWriter, r *http.Request) {
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

	//Creds
	creds := &SearchInfo{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	var users []Profile
	rows, err := db.Query("select * from users where name like '" + creds.Query + "%'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user Profile
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
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(users)
	w.Write(ret)
}