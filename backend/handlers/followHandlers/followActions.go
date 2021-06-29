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
	"goApp/backend/handlers"
)

//Wrapper
type extended handlers.HandlerWrapper


//Follow adds the users from the following/follower list in the respective db entries
func Follow(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)


	//Decode Creds
	creds := &types.FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check is user is already following other user
	var following []string
	row := db.DB.QueryRow("select following from users where username=$1", creds.Follower)
	err = row.Scan(pq.Array(&following))
	for _, username := range following {
		if username == creds.Following {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Update DB on following end
	query := fmt.Sprintf("UPDATE users SET following = following || '%s'::text WHERE username = '%s';", creds.Following, creds.Follower)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update DB on followers end
	query = fmt.Sprintf("UPDATE users SET followers = followers || '%s'::text WHERE username = '%s';", creds.Follower, creds.Following)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Unfollow removes the users from the following/follower list in the respective db entries
func Unfollow(w http.ResponseWriter, r *http.Request) {

	//Authentication
	helpers.Authenticate(w,r)

	//Decode Creds
	creds := &types.FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Update DB on following end
	query := fmt.Sprintf("UPDATE users SET following = ARRAY_REMOVE(following,'%s'::text) WHERE username = '%s';", creds.Following, creds.Follower)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update DB on followers end
	query = fmt.Sprintf("UPDATE users SET followers = ARRAY_REMOVE(followers,'%s'::text)WHERE username = '%s';", creds.Follower, creds.Following)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}