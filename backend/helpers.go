
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
	name   = "cookie-name"
)


// Reverse a list
func reverse(posts []Post) []Post {
	newList := make([]Post, len(posts))
	for i, j := 0, len(posts)-1; i <= j; i, j = i+1, j-1 {
		newList[i], newList[j] = posts[j], posts[i]
	}
	return newList
}


//Unique returns unique array
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	newList := []string{}
	for _, entry := range slice {
	    if _, value := keys[entry]; !value {
		keys[entry] = true
		newList = append(newList, entry)
	    }
	}    
	return newList
    }
    
   
//Authenticate 
func authenticate(w http.ResponseWriter,r *http.Request){
	session, _ := store.Get(r, name)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}