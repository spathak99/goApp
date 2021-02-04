package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
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
	http.HandleFunc("/initial_custom_program", InitializeProgram)
	http.HandleFunc("/update_custom_program", UpdateCustomProgram)
	http.HandleFunc("/get_custom_program", GetCustomProgram)
	http.HandleFunc("/get_personal_feed", GetPersonalFeed)
	http.HandleFunc("/search", FuzzySearch)
	http.HandleFunc("/update_name", UpdateName)
	/*c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
	})*/

	initDB()
	router := http.DefaultServeMux

	log.Fatal(http.ListenAndServe(":8000",
		handlers.LoggingHandler(os.Stdout, handlers.CORS(
			handlers.AllowedMethods([]string{"POST"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"X-Requested-With"}))(router))))

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
