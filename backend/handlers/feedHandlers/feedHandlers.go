package handlers

import (
	"encoding/json"
	"net/http"	
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"strings"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// GetFeed grabs the news feed for a given user
func GetFeed(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Gets all users that current user is following
	var following []string
	username := creds["username"].(string)
	row := db.DB.QueryRow("select following from users where username=$1", username)
	err = row.Scan(pq.Array(&following))

	//Gets all posts for news feed
	var postList []types.Post
	sqlRaw := fmt.Sprintf(`select * from posts where username in ('%s')`, strings.Join(following, "','"))
	rows, err := db.DB.Query(sqlRaw)
	if err != nil {
		panic(err)
	}

	//Read in posts from DB
	defer rows.Close()
	for rows.Next() {
		var currPost types.Post
		err = rows.Scan(&currPost.ID, &currPost.Username,
			&currPost.Contents, &currPost.Media,
			&currPost.Date, pq.Array(&currPost.Likes))
		if err != nil {
			panic(err)
		}
		postList = append(postList, currPost)
	}

	//Write response
	w.Header().Set("Content-Type", "application/json")
	postList = helpers.Reverse(postList)
	ret, err := json.Marshal(postList)
	w.Write(ret)
}

// GetPersonalFeed grabs posts that the user made
func GetPersonalFeed(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Get all feed for user
	var postList []types.Post
	username := creds["username"].(string)
	rows, err := db.DB.Query(`select * from posts where username=$1`, username)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currPost types.Post
		err = rows.Scan(&currPost.ID, &currPost.Username,
			&currPost.Contents, &currPost.Media,
			&currPost.Date, pq.Array(&currPost.Likes))
		if err != nil {
			panic(err)
		}
		postList = append(postList, currPost)
	}

	//Write Response
	w.Header().Set("Content-Type", "application/json")
	postList = helpers.Reverse(postList)
	ret, err := json.Marshal(postList)
	w.Write(ret)
}