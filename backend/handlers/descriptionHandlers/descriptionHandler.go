package handlers

import (
	"encoding/json"
	"net/http"	
	"fmt"
	_ "github.com/lib/pq"
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
)

// UpdateDescription updates the bio of the given user in their db entry
func UpdateDescription(w http.ResponseWriter, r *http.Request) {

	helpers.Authenticate(w,r)

	//Decode Creds
	creds := &types.Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Updates description
	query := fmt.Sprintf("UPDATE users SET description = '%s' WHERE username = '%s';", creds.Description, creds.Username)
	if _, err = db.DB.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}