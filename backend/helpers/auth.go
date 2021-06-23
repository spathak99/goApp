package helpers

import (
	"net/http"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)


//Cookie for authentication
var Cookie = securecookie.GenerateRandomKey(32)
var Store  = sessions.NewCookieStore(Cookie)
var CookieName   = "cookie-name"

//Authenticate checks if user is logged in before allowing inside functions to run
func Authenticate(w http.ResponseWriter,r *http.Request){
	session, _ := Store.Get(r, CookieName)
	auth, _ := session.Values["authenticated"].(bool)
	if !auth {
		if _, ok := session.Values["authenticated"]; ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
}