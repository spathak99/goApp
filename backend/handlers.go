package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"net/http"
	 "fmt" 
	_ "github.com/lib/pq"    
	"github.com/gorilla/sessions"
)


/*
	Session keys
*/
var (
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)



/*
	Update profile description
*/
func UpdateDescription(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	
	query := fmt.Sprintf("UPDATE users SET description = '%s' WHERE username = '%s';",creds.Description,creds.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


/*
Calculate Calories left for the day and update as needed
*/

func Calories(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	//TODO: Implement Queries and calorie calculation, along with meals stuff
}



/*
	Update to your current weight
*/
func UpdateWeights(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	
	query := fmt.Sprintf("UPDATE users SET bodyweight = '%f' WHERE username = '%s';",creds.Bodyweight,creds.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query2 := fmt.Sprintf("UPDATE users SET goalweight = '%f' WHERE username = '%s';",creds.GoalWeight,creds.Username)
	if _, err = db.Query(query2); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


/*
	End session for suer
*/
func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)
}


/*
	Signup a new user
*/
func Signup(w http.ResponseWriter, r *http.Request){	
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	query := "insert into users values ($1, $2,$3,$4,$5,$6)"
	if _, err = db.Query(query, 
		creds.Username, 
		string(hashedPassword),
		string(creds.Description),
		creds.GoalWeight,
		creds.Bodyweight,
		creds.CalorieGoal); 
		err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
}

/*
	Signin existing users
*/

func Signin(w http.ResponseWriter, r *http.Request){
	//Start Session
	session, _ := store.Get(r, "cookie-name")

	//User authentication below
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	result := db.QueryRow("select password from users where username=$1", creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	storedCreds := &Profile{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	//Verify Session
	session.Values["authenticated"] = true
    session.Save(r, w)
}





