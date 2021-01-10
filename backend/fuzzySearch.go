package main

import (
	"encoding/json"
	"net/http"
)

// FuzzySearch does a fuzzy search of the name of the user
func FuzzySearch(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	//Creds
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//DB Query
	var users []SearchInfo
	rows, err := db.Query("select username, name from users where name like $1", creds.Name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user SearchInfo
		err = rows.Scan(&user.Username, &user.Name)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	ret, err := json.Marshal(users)
	w.Write(ret)
}
