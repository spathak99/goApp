package main

import (
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/sahilm/fuzzy"
)

// SearchUsers uses a fuzzy search to search for users to follow
func SearchUsers(w http.ResponseWriter, r *http.Request) {
	//Authentication
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

	const bold = "\033[1m%s\033[0m"
	pattern := "mnr"
	data := GrabData()
	matches := fuzzy.Find(pattern, data)

	for _, match := range matches {
		for i := 0; i < len(match.Str); i++ {
			if contains(i, match.MatchedIndexes) {
				fmt.Print(fmt.Sprintf(bold, string(match.Str[i])))
			} else {
				fmt.Print(string(match.Str[i]))
			}
		}
		fmt.Println()
	}
}

// Contains checks if the element is contained
func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

// GrabData grabs the data required for the fuzzy search
func GrabData() []string {
	var users []string
	row := db.QueryRow("select username from users;")
	err := row.Scan(pq.Array(&users))

	if err != nil {
		print(err)
	}
	return users
}
