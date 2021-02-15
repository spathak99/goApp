package main

import (
	"database/sql"
	"net/http"

	"github.com/rs/cors"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"
)

const hashCost = 8

var db *sql.DB
var server *http.Server

func main() {
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
		Secure:   false,
	}
	startServer()
}

func startServer() {
	print("Starting Server")

	mux := http.NewServeMux()
	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/signup", Signup)
	mux.HandleFunc("/logout", Logout)
	mux.HandleFunc("/get_all_users", getUsers)
	mux.HandleFunc("/update_bio", UpdateDescription)
	mux.HandleFunc("/update_weight", UpdateWeights)
	mux.HandleFunc("/update_calories", UpdateCalories)
	mux.HandleFunc("/get_user_data", GetUserData)
	mux.HandleFunc("/get_followers", GetFollowers)
	mux.HandleFunc("/get_following", GetFollowing)
	mux.HandleFunc("/follow", Follow)
	mux.HandleFunc("/unfollow", Unfollow)
	mux.HandleFunc("/make_post", MakePost)
	mux.HandleFunc("/get_feed", GetFeed)
	mux.HandleFunc("/like_post", LikePost)
	mux.HandleFunc("/unlike_post", Unlike)
	mux.HandleFunc("/initial_custom_program", InitializeProgram)
	mux.HandleFunc("/update_custom_program", UpdateCustomProgram)
	mux.HandleFunc("/get_custom_program", GetCustomProgram)
	mux.HandleFunc("/get_personal_feed", GetPersonalFeed)
	mux.HandleFunc("/search", FuzzySearch)
	mux.HandleFunc("/update_name", UpdateName)
	mux.HandleFunc("/initialize_lifts", InitializeLifts)
	mux.HandleFunc("/update_lifts", UpdateLifts)
	mux.HandleFunc("/estimate_max", EstimateMax)
	mux.HandleFunc("/logexercise", LogExercise)
	mux.HandleFunc("/grablog", GrabLog)
	initDB()
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8000", handler)

}

func stopServer() {
	server.Close()
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
