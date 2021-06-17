package liftHandlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"	
	"goApp/backend/db"
	"goApp/backend/types"
	"goApp/backend/helpers"
	"github.com/lib/pq"
)


//LogExercise logs a lift for the user
func LogExercise(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input, _ := json.Marshal(creds)
	var currLifts []string
	row := db.DB.QueryRow("select lifts from exerciselog where username=$1", creds.Username)
	err = row.Scan(pq.Array(&currLifts))

	//Query
	if len(currLifts) == 0 {
		var inputarr [1]string
		inputarr[0] = string(input)
		query := "insert into exerciselog values ($1,$2)"
		if _, err = db.DB.Query(query,
			creds.Username,
			pq.Array(inputarr)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		query := fmt.Sprintf("UPDATE exerciselog SET lifts = lifts || '%s'::text WHERE username = '%s'", string(input), creds.Username)
		if _, err = db.DB.Query(query); err != nil {
			print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}


//GrabLog gets the logged exercises by lift
func GrabLog(w http.ResponseWriter, r *http.Request) {
	//Authentication
	helpers.Authenticate(w,r)

	//Credentials
	creds := &types.Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query
	var liftlog []string
	row := db.DB.QueryRow("select lifts from exerciselog where username=$1", creds.Username)
	err = row.Scan(pq.Array(&liftlog))


	//Loop
	var trendLifts []string
	for _, lift := range liftlog {
		srcjson := []byte(lift)
		var helper types.Lift
		err := json.Unmarshal(srcjson, &helper)
		if err != nil {
			panic(err)
		}

		if helper.Name == creds.Name {
			out, err := json.Marshal(helper)
			if err != nil {
				panic(err)
			}
			trendLifts = append(trendLifts, string(out))
		}
	}

	//Write Response
	w.Write([]byte(fmt.Sprint(trendLifts)))
}