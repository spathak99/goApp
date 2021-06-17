package types

import (
	"encoding/json"
)

// UserLifts for the user
type UserLifts struct {
	Username string `json:"id", db:"id"`
	Lifts    string `json:"lifts",db"lifts"`
}

// UserLiftsHelper helps with parsing the json and does not communicate w/ DB
type UserLiftsHelper struct {
	Username string
	Lifts    json.RawMessage
}