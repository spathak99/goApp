package postHandlers

import (
	"encoding/json"
	"net/http"	
	"fmt"
	"crypto/rand"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// MakePost creates a post and adds that post to the posts table
func MakePost(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.Post{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Generate Random ID
	len := 6
	id := make([]byte, len)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	uniqueID := fmt.Sprintf("%X", id)

	//Query
	query := "insert into posts values ($1,$2,$3,$4,$5,$6)"
	if _, err = db.DB.Query(query,
		uniqueID,
		creds.Username,
		string(creds.Contents),
		string(creds.Media),
		string(creds.Date),
		pq.Array(creds.Likes)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetPost grabs an individual post by ID
func GetPost(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query post
	id := creds["id"].(string)
	row := db.DB.QueryRow(`select * from posts where id=$1`, id)
	var currPost types.Post
	err = row.Scan(&currPost.ID, &currPost.Username,
		&currPost.Contents, &currPost.Media,
		&currPost.Date, pq.Array(&currPost.Likes))

	//Write Response
	w.Header().Set("Content-Type", "application/json")
	ret, err := json.Marshal(currPost)
	w.Write(ret)
}