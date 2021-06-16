
package main

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)


//Cookie for authentication
var (
	cookie = securecookie.GenerateRandomKey(32)
	store  = sessions.NewCookieStore(cookie)
	cookieName   = "cookie-name"
)


// Reverse a list
func reverse(source []Post) []Post {
	destination := make([]Post, len(source))
	for i, j := 0, len(source)-1; i <= j; i, j = i+1, j-1 {
		destination[i], destination[j] = source[j], source[i]
	}
	return destination
}


//Unique returns unique array
func Unique(source []string) []string {
	keys := make(map[string]bool)
	destination := []string{}
	for _, entry := range source {
	    if _, value := keys[entry]; !value {
		keys[entry] = true
		destination = append(destination, entry)
	    }
	}    
	return destination
    }
    
   
//Check if user is authenticated
func authenticate(w http.ResponseWriter,r *http.Request){
	session, _ := store.Get(r, cookieName)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}