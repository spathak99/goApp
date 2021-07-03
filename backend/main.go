package main

import (
	"net/http"
	"github.com/gorilla/sessions"
	"goApp/backend/server"
	"goApp/backend/helpers"
	_ "github.com/lib/pq"
)

//Server
const hashCost = 8

func main() {
	helpers.Store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: false,
		Secure:   false,
	}
	server.StartServer()
}

//Allow cross origin sharing
func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

