package testSuite

import (
	"goApp/backend/db"
	"testing"
	"github.com/lib/pq"
	"fmt"
	"goApp/backend/likeHandlers"
	"github.com/stretchr/testify/assert"
)

// LikesTestHelper is a helper function for the like/unlike post tests
func LikesTestHelper(query string) []string {
	//DB query
	var likes []string
	row := db.DB.QueryRow(query)
	err := row.Scan(pq.Array(&likes))
	if(err != nil){
		panic(err)
	}

	return likes
}

// TestLikes tests liking and unliking posts
func TestLikes(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "id":"5492C1CA32B7"
    }`)

	query := fmt.Sprintf("select likes from posts where id='%s'", "5492C1CA32B7")

	//Test 1
	resp1,_ := TstHelper(mockData1, likeHandlers.LikePost, "/like_post")
	likes := LikesTestHelper(query)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, likes, "testingaccount")

	resp2,_ := TstHelper(mockData1, likeHandlers.Unlike, "/unlike_post")
	likes2 := LikesTestHelper(query)
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, likes2, "testingaccount")
}