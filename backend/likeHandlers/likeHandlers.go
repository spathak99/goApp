package likeHandlers


import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/lib/pq"
	_ "github.com/lib/pq"	
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// LikePost adds a username to the list of likes on a given post
func LikePost(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Creds
	creds := &types.Like{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check if user already liked post
	var likes []string
	row := db.DB.QueryRow("select likes from posts where id=$1", creds.ID)
	err = row.Scan(pq.Array(&likes))
	for _, username := range likes {
		if username == creds.Username {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Add username to list of user who like post
	query := fmt.Sprintf("UPDATE posts SET likes = likes || '%s'::text WHERE ID = '%s';", creds.Username, creds.ID)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Unlike removes a user from the list of likes on a post
func Unlike(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Creds
	creds := &types.Like{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check if user likes post
	var isLiked = false
	var likes []string
	row := db.DB.QueryRow("select likes from posts where id=$1", creds.ID)
	err = row.Scan(pq.Array(&likes))
	for _, username := range likes {
		if username == creds.Username {
			isLiked = true
			break
		}
	}

	//Return
	if !isLiked {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Remove username from list of users who like post
	query := fmt.Sprintf("UPDATE posts SET likes = ARRAY_REMOVE(likes,'%s'::text) WHERE ID = '%s';", creds.Username, creds.ID)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}