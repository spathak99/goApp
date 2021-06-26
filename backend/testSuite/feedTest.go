package testSuite


import (
	"encoding/json"
	"testing"
	"goApp/backend/types"
	"goApp/backend/db"
	feedHandlers "goApp/backend/handlers/feedHandlers"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// FeedTestHelper helps with grabbing the news feed
func FeedTestHelper(res string) ([]string) {
	var usernames []string
	var posts []types.Post

	err := json.Unmarshal([]byte(res), &posts)
	if(err != nil){
		panic(err)
	}

	for _, entry := range posts {
		usernames = append(usernames, entry.Username)
	}

	return usernames
}

//NewsFeedTest checks if a feed can be grabbed for a user
func NewsFeedTest(t *testing.T) {
	mockData := []byte(`{
		"username":"testingaccount"
	}`)

	query1 := "select following from users where username='testingaccount'"
	var following []string
	row := db.DB.QueryRow(query1)
	err := row.Scan(pq.Array(&following))
	if err != nil {
		panic(err)
	}

	resp,res := TstHelper(mockData,feedHandlers.GetFeed, "/get_feed")
	usernames := FeedTestHelper(res)
	assert.Equal(t, resp, 200)
	for _, username := range usernames {
		assert.Contains(t, following, username)
	}
}


//PersonalFeedTest checks if the feed for the user is retrieved
func PersonalFeedTest(t *testing.T) {
	mockData := []byte(`{
		"username":"Shardool"
	}`)
	resp,_ := TstHelper(mockData, feedHandlers.GetPersonalFeed, "/get_personal_feed")
	assert.Equal(t, resp, 200)
}

