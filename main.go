package main

import (
	"database/sql"
	"net/http"
	"log"
	_ "github.com/lib/pq"
)

const hashCost = 8
var db *sql.DB

func main() {
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/update_bio",UpdateDescription)
	http.HandleFunc("/update_weight",UpdateWeight)
	initDB()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDB(){
	var err error
	db, err = sql.Open("postgres", "dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
}