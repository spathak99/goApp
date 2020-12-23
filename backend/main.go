package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const hashCost = 8

var db *sql.DB
var server *http.Server

func main() {
	startServer()
}

func startServer() {
	print("Starting Server")
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/update_bio", UpdateDescription)
	http.HandleFunc("/update_weight", UpdateWeights)
	http.HandleFunc("/update_calories", UpdateCalories)
	http.HandleFunc("/get_user_data", GetUserData)
	http.HandleFunc("/follow", Follow)
	http.HandleFunc("/unfollow", Unfollow)
	http.HandleFunc("/make_post", MakePost)
	http.HandleFunc("/get_feed", GetFeed)
	http.HandleFunc("/like_post", LikePost)
	http.HandleFunc("/unlike_post", Unlike)
	http.HandleFunc("/update_custom_program", UpdateCustomProgram)
	http.HandleFunc("/get_custom_program", GetCustomProgram)
	initDB()
	server = &http.Server{
		Addr:    ":8000",
		Handler: http.DefaultServeMux,
	}
	log.Fatal(server.ListenAndServe())
}

func stopServer() {
	server.Close()
}

func initDB() {
	var err error
	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}
