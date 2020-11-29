package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"net/http"
	"github.com/lib/pq"
	 "fmt" 
	 "strings"
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
	Add Followers
*/
func Follow(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}	

    creds := &FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	query := fmt.Sprintf("UPDATE users SET following = following || '%s'::text WHERE username = '%s';",creds.Following,creds.Follower)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	query = fmt.Sprintf("UPDATE users SET followers = followers || '%s'::text WHERE username = '%s';",creds.Follower,creds.Following)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
	Remove Followers
*/

func Unfollow(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}	

    creds := &FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	query := fmt.Sprintf("UPDATE users SET following = ARRAY_REMOVE(following,'%s'::text) WHERE username = '%s';",creds.Following,creds.Follower)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	query = fmt.Sprintf("UPDATE users SET followers = ARRAY_REMOVE(followers,'%s'::text)WHERE username = '%s';",creds.Follower,creds.Following)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

/*
	Update profile description
*/
func UpdateDescription(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}

	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	
	//Updates description
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

func UpdateCalories(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	
	var currCals int 
	row := db.QueryRow("select caloriesleft from users where username=$1", creds.Username)
	err = row.Scan(&currCals)
	if((int(creds.CaloriesLeft) > currCals) || (int(creds.CaloriesLeft) < 0)){
	rp := []byte(`{
			"response":"Cannot increase calories left or make them negative"
		}`)
		w.Write(rp)
		return
	}
	
	query2 := fmt.Sprintf("UPDATE users SET caloriesleft = '%f' WHERE username = '%s';",creds.CaloriesLeft,creds.Username)
	if _, err = db.Query(query2); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rp := []byte(`{
		"response":"Succesfully updated calories"
	}`)
	w.Write(rp)
	return
}



/*
	Update to your Goal weight
*/
func UpdateGoal(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}
	
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}
	//Update Bodyweight
	query := fmt.Sprintf("UPDATE users SET bodyweight = '%f' WHERE username = '%s';",creds.Bodyweight,creds.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update Goalweight
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

	//Logs new user
	query := "insert into users values ($1, $2,$3,$4,$5,$6,$7,$8,$9)"

	if _, err = db.Query(query, 
		creds.Username, 
		string(hashedPassword),
		string(creds.Description),
		creds.GoalWeight,
		creds.Bodyweight,
		creds.CalorieGoal,
		creds.CaloriesLeft,
		pq.Array(creds.Followers),
		pq.Array(creds.Following)); 
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

	//Hash Password
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password),[]byte(creds.Password)); err != nil{
		w.WriteHeader(http.StatusUnauthorized)
	}

	//Verify Session
	session.Values["authenticated"] = true
    session.Save(r, w)
}


/*
Grab all profile fields
*/
func GetUserData(w http.ResponseWriter,r *http.Request){
	//Start session
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}
	
	//Credentials
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}	

    //Grabs user data from db
	var user Profile
	row := db.QueryRow("select * from users where username=$1", creds.Username)
	err = row.Scan(&user.Username, &user.Password, &user.Description,
				   &user.GoalWeight, &user.Bodyweight, 
				   &user.CalorieGoal,&user.CaloriesLeft,
				   pq.Array(&user.Followers),pq.Array(&user.Following))

	w.Header().Set("Content-Type", "application/json")

	ret, err := json.Marshal(user)
	w.Write(ret)
}


/*
	Upload a post
*/
func MakePost(w http.ResponseWriter, r *http.Request){
	//Start session
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}

	//Credentials
	creds := &Post{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}	

	query := "insert into posts values ($1, $2,$3,$4,$5,$6)"
	if _, err = db.Query(query, 
		creds.ID,
		creds.Username, 
		string(creds.Contents),
		string(creds.Media),
		string(creds.Date),
		pq.Array(creds.Likes)); 
		err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
}

/*
	Reverses news feed so it is in order
*/
func reverse(posts []Post) []Post {
	newList := make([]Post, len(posts))
	for i, j := 0, len(posts)-1; i <= j; i, j = i+1, j-1 {
		newList[i], newList[j] = posts[j], posts[i]
	}
	return newList
}



/*
	Like a post
*/
func Like_Post(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}
	
	//Creds
	creds := &Like{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}	

	//Add username to list of user who like post
	query := fmt.Sprintf("UPDATE posts SET likes = likes || '%s'::text WHERE ID = '%s';",creds.Username,creds.ID)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}


/*
	Grabs news feed for a user
*/
func GetFeed(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
	}

	//Credentials
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if(err != nil){
		w.WriteHeader(http.StatusBadRequest)
		return 
	}	

	//Gets all users that current user is following
	var following []string
	row := db.QueryRow("select following from users where username=$1", creds.Username)
	err = row.Scan(pq.Array(&following))
	

	//Gets all posts for news feed
	var post_list []Post
	fllwng := strings.Join(following, "','")
	sqlRaw := fmt.Sprintf(`select * from posts where username in ('%s')`, fllwng)
	rows, err := db.Query(sqlRaw)
	if(err != nil){
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var curr_post Post
		err = rows.Scan(&curr_post.ID,&curr_post.Username,
						&curr_post.Contents,&curr_post.Media,
						&curr_post.Date)
		if err != nil {
			panic(err)
		}
		post_list = append(post_list,curr_post)
	}
	post_list = reverse(post_list)
	ret, err := json.Marshal(post_list)
	w.Write(ret)
}


