package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)


// InitializeLifts initializes the first set of lifts entered by the user
func InitializeLifts(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := &UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := "insert into userlifts values ($1,$2)"

	if _, err = db.Query(query,
		lifts.Username,
		string(lifts.Lifts)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateLifts updates the max lifts of the user
func UpdateLifts(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := &UserLiftsHelper{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Shift to the struct that is compatible with the db
	lifts := UserLifts{}
	lifts.Username = creds.Username
	lifts.Lifts = string(creds.Lifts)

	//DB Query
	query := fmt.Sprintf("UPDATE userlifts SET lifts= '%s' WHERE username = '%s';", lifts.Lifts, lifts.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Write Response
	ret := []byte(`{
		"response":"Succesfully updated program"
	}`)
	w.Write(ret)
}

//GetUserMax gets the max of the user
func GetUserMax(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Gets lifts of user
	var lifts string
	username := creds["username"].(string)
	row := db.QueryRow("select lifts from userlifts where username=$1", username)
	err = row.Scan(&lifts)

	ret := []byte(lifts)
	w.Write(ret)
}

//EstimateMax calculates the estimated one rep max
func EstimateMax(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := &Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Calculation and response
	reps := float64(creds.Reps)
	weight := float64(creds.Weight)
	ERM := int((weight) / (1.0278 - 0.0278*((10.0-creds.RPE)+reps)))
	ret := strconv.Itoa(ERM)
	w.Write([]byte(ret))
}

//LogExercise logs a lift for the user
func LogExercise(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := &Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	inp, _ := json.Marshal(creds)
	var currLifts []string
	row := db.QueryRow("select lifts from exerciselog where username=$1", creds.Username)
	err = row.Scan(pq.Array(&currLifts))

	//Query
	if len(currLifts) == 0 {
		var inputarr [1]string
		inputarr[0] = string(inp)
		query := "insert into exerciselog values ($1,$2)"
		if _, err = db.Query(query,
			creds.Username,
			pq.Array(inputarr)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		query := fmt.Sprintf("UPDATE exerciselog SET lifts = lifts || '%s'::text WHERE username = '%s'", string(inp), creds.Username)
		if _, err = db.Query(query); err != nil {
			print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

//Unique returns unique array
func Unique(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

//GetLiftNames gets all the types of lifts that have been logged
func GetLiftNames(w http.ResponseWriter, r *http.Request){
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

	//Credentials
	creds := map[string]interface{}{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//DB Query
	var username = creds["username"].(string)
	var liftlog []string
	row := db.QueryRow("select lifts from exerciselog where username=$1", username)
	err = row.Scan(pq.Array(&liftlog))

	//Get Names
	var temp []string
	for _, lift := range liftlog {
		srcjson := []byte(lift)
		var helper Lift
		err := json.Unmarshal(srcjson, &helper)
		if err != nil {
			panic(err)
		}
		temp = append(temp,helper.Name)
	}
	var ret = Unique(temp)
	response, err := json.Marshal(ret)
	w.Write(response)
}

//GrabLog gets the logged exercises by lift
func GrabLog(w http.ResponseWriter, r *http.Request) {
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

	//Credentials
	creds := &Lift{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Query
	var liftlog []string
	row := db.QueryRow("select lifts from exerciselog where username=$1", creds.Username)
	err = row.Scan(pq.Array(&liftlog))


	//Loop
	var trendLifts []string
	for _, lift := range liftlog {
		srcjson := []byte(lift)
		var helper Lift
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


