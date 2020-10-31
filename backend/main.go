package main

import (
	"database/sql"
	"net/http"
	"log"
	_ "github.com/lib/pq"
)

const hashCost = 8
var db *sql.DB
var server *http.Server

func main() {
	startServer();
}

func startServer(){
	print("Starting Server")
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/update_bio",UpdateDescription)
	http.HandleFunc("/update_weight",UpdateWeights)
	http.HandleFunc("/update_calories",UpdateCalories)
	http.HandleFunc("/get_user_data",GetUserData)
	http.HandleFunc("/follow",Follow)
	http.HandleFunc("/unfollow",Unfollow)
	http.HandleFunc("/make_post",MakePost)
	http.HandleFunc("/get_feed",GetFeed)
	initDB()
	server = &http.Server{
        Addr:    ":8000",
        Handler: http.DefaultServeMux,
    }
	log.Fatal(server.ListenAndServe())
}

func stopServer(){
	server.Close()
}

func initDB(){
	var err error
	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}