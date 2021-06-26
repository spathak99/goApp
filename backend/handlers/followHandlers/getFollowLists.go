package handlers

import (
	"encoding/json"
	"net/http"	
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/helpers"
)

//GetFollowing gets a list of the people the user follows
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query
	username := creds["username"].(string)
	var following []string
	row := db.DB.QueryRow("select following from users where username=$1", username)
	err = row.Scan(pq.Array(&following))

	ret, err := json.Marshal(following)
	w.Write(ret)
}

//GetFollowers gets a list of your followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query
	username := creds["username"].(string)
	var followers []string
	row := db.DB.QueryRow("select followers from users where username=$1", username)
	err = row.Scan(pq.Array(&followers))

	ret, err := json.Marshal(followers)
	w.Write(ret)
}
