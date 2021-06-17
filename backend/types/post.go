package types

// Post struct
type Post struct {
	ID       string   `json:"id", db:"id"`
	Username string   `json:"username",db"username"`
	Contents string   `json:"contents",db:"contents"`
	Media    string   `json:"media",db:"media`
	Date     string   `json:"date",db:"date"`
	Likes    []string `json:"likes",db:"likes"`
}
