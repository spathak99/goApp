package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Session Key and Store
var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//Follow adds the users from the following/follower list in the respective db entries
func Follow(w http.ResponseWriter, r *http.Request) {

	//Authentication
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	creds := &FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check is user is already following other user
	var following []string
	row := db.QueryRow("select following from users where username=$1", creds.Follower)
	err = row.Scan(pq.Array(&following))
	for _, username := range following {
		if username == creds.Following {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Update DB on following end
	query := fmt.Sprintf("UPDATE users SET following = following || '%s'::text WHERE username = '%s';", creds.Following, creds.Follower)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update DB on followers end
	query = fmt.Sprintf("UPDATE users SET followers = followers || '%s'::text WHERE username = '%s';", creds.Follower, creds.Following)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Unfollow removes the users from the following/follower list in the respective db entries
func Unfollow(w http.ResponseWriter, r *http.Request) {

	//Authentication
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	creds := &FollowRelation{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Update DB on following end
	query := fmt.Sprintf("UPDATE users SET following = ARRAY_REMOVE(following,'%s'::text) WHERE username = '%s';", creds.Following, creds.Follower)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update DB on followers end
	query = fmt.Sprintf("UPDATE users SET followers = ARRAY_REMOVE(followers,'%s'::text)WHERE username = '%s';", creds.Follower, creds.Following)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateDescription updates the bio of the given user in their db entry
func UpdateDescription(w http.ResponseWriter, r *http.Request) {

	//Authentication
	session, _ := store.Get(r, "cookie-name")
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
	query := fmt.Sprintf("UPDATE users SET description = '%s' WHERE username = '%s';", creds.Description, creds.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateCalories updates the calorie goals and calorie counts for the user
func UpdateCalories(w http.ResponseWriter, r *http.Request) {

	//Authentication
	session, _ := store.Get(r, "cookie-name")
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
	if (int(creds.CaloriesLeft) > currCals) || (int(creds.CaloriesLeft) < 0) {
		rp := []byte(`{
			"response":"Cannot increase calories left or make them negative"
		}`)
		w.Write(rp)
		return
	}

	query2 := fmt.Sprintf("UPDATE users SET caloriesleft = '%f' WHERE username = '%s';", creds.CaloriesLeft, creds.Username)
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

// UpdateWeights updates the goal weight and current body weight of the user
func UpdateWeights(w http.ResponseWriter, r *http.Request) {

	//Authentication
	session, _ := store.Get(r, "cookie-name")
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
	query := fmt.Sprintf("UPDATE users SET bodyweight = '%f' WHERE username = '%s';", creds.Bodyweight, creds.Username)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Update Goalweight
	query2 := fmt.Sprintf("UPDATE users SET goalweight = '%f' WHERE username = '%s';", creds.GoalWeight, creds.Username)
	if _, err = db.Query(query2); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Logout logs the user out and ends the session
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

// Signup creates a new entry in the users table in the db
func Signup(w http.ResponseWriter, r *http.Request) {
	creds := &Profile{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	//Logs new user
	query := "insert into users values ($1, $2,$3,$4,$5,$6,$7,$8,$9,$10)"

	if _, err = db.Query(query,
		creds.Username,
		string(hashedPassword),
		string(creds.Description),
		creds.GoalWeight,
		creds.Bodyweight,
		creds.CalorieGoal,
		creds.CaloriesLeft,
		pq.Array(creds.Followers),
		pq.Array(creds.Following),
		string(creds.Program)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Signin signs in the user and authenticates them
func Signin(w http.ResponseWriter, r *http.Request) {
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
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	//Verify Session
	session.Values["authenticated"] = true
	session.Save(r, w)
}

// GetUserData grabs profile struct data for the given user
func GetUserData(w http.ResponseWriter, r *http.Request) {

	//Authentication
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
		&user.CalorieGoal, &user.CaloriesLeft,
		pq.Array(&user.Followers), pq.Array(&user.Following), &user.Program)

	w.Header().Set("Content-Type", "application/json")

	ret, err := json.Marshal(user)
	w.Write(ret)
}

// MakePost creates a post and adds that post to the posts table
func MakePost(w http.ResponseWriter, r *http.Request) {

	//Authentication
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

	//Generate Random ID
	n := 6
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	uniqueID := fmt.Sprintf("%X", b)

	//Query
	query := "insert into posts values ($1,$2,$3,$4,$5,$6)"
	if _, err = db.Query(query,
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

// Reverse
func reverse(posts []Post) []Post {
	newList := make([]Post, len(posts))
	for i, j := 0, len(posts)-1; i <= j; i, j = i+1, j-1 {
		newList[i], newList[j] = posts[j], posts[i]
	}
	return newList
}

// LikePost adds a username to the list of likes on a given post
func LikePost(w http.ResponseWriter, r *http.Request) {

	//Authentication
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

	//Check if user already liked post
	var likes []string
	row := db.QueryRow("select likes from posts where id=$1", creds.ID)
	err = row.Scan(pq.Array(&likes))
	for _, username := range likes {
		if username == creds.Username {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//Add username to list of user who like post
	query := fmt.Sprintf("UPDATE posts SET likes = likes || '%s'::text WHERE ID = '%s';", creds.Username, creds.ID)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

// Unlike removes a user from the list of likes on a post
func Unlike(w http.ResponseWriter, r *http.Request) {

	//Authentication
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

	//Check if user likes post
	var isLiked = false
	var likes []string
	row := db.QueryRow("select likes from posts where id=$1", creds.ID)
	err = row.Scan(pq.Array(&likes))
	for _, username := range likes {
		if username == creds.Username {
			isLiked = true
			break
		}
	}
	if !isLiked {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Remove username from list of users who like post
	query := fmt.Sprintf("UPDATE posts SET likes = ARRAY_REMOVE(likes,'%s'::text) WHERE ID = '%s';", creds.Username, creds.ID)
	if _, err = db.Query(query); err != nil {
		print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetFeed grabs the news feed for a given user
func GetFeed(w http.ResponseWriter, r *http.Request) {

	//Authentication
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

	//Gets all users that current user is following
	var following []string
	row := db.QueryRow("select following from users where username=$1", creds.Username)
	err = row.Scan(pq.Array(&following))

	//Gets all posts for news feed
	var postList []Post
	fllwng := strings.Join(following, "','")
	sqlRaw := fmt.Sprintf(`select * from posts where username in ('%s')`, fllwng)
	rows, err := db.Query(sqlRaw)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currPost Post
		err = rows.Scan(&currPost.ID, &currPost.Username,
			&currPost.Contents, &currPost.Media,
			&currPost.Date)
		if err != nil {
			panic(err)
		}
		postList = append(postList, currPost)
	}
	postList = reverse(postList)
	ret, err := json.Marshal(postList)
	w.Write(ret)
}
