package testSuite

import (
	"goApp/backend/db"
	"testing"
	"github.com/lib/pq"
	"fmt"
	followHandlers "goApp/backend/handlers/followHandlers"
	"github.com/stretchr/testify/assert"
)


// FollowTestHelper is the helper function for the follow test
func FollowTestHelper(query1 string, query2 string) ([]string, []string) {
	//DB Queries
	var following []string
	row := db.DB.QueryRow(query1)
	err := row.Scan(pq.Array(&following))

	var followers []string
	row = db.DB.QueryRow(query2)
	err = row.Scan(pq.Array(&followers))

	if(err != nil){
		panic(err)
	}

	return following, followers
}

// TestFollow tests if the follow and unfollow handlers work as intended
func TestFollow(t *testing.T) {
	//Test 1
	mockData1 := []byte(`{
		"follower":"testingaccount",
		"following":"Shardool"
	}`)

	followerQuery := fmt.Sprintf("select following from users where username='%s'", "testingaccount")
	followedQuery := fmt.Sprintf("select followers from users where username='%s'", "Shardool")

	resp1,_ := Test_Helper(mockData1,followHandlers.Follow,"/follow")
	following, followers := FollowTestHelper(followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Shardool")
	assert.Contains(t, followers, "testingaccount")

	resp2,_ := Test_Helper(mockData1,followHandlers.Unfollow,"/unfollow")
	following2, followers2 := FollowTestHelper(followerQuery, followedQuery)	
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, following2, "Shardool")
	assert.NotContains(t, followers2, "testingaccount")

	//Test 2
	mockData2 := []byte(`{
		"follower":"testingaccount",
		"following":"Bijon"
	}`)

	followerQuery = fmt.Sprintf("select following from users where username='%s'", "testingaccount")
	followedQuery = fmt.Sprintf("select followers from users where username='%s'", "Bijon")

	resp1,_ = Test_Helper(mockData2,followHandlers.Follow,"/follow")
	following, followers = FollowTestHelper(followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Bijon")
	assert.Contains(t, followers, "testingaccount")

	resp2,_ = Test_Helper(mockData2,followHandlers.Unfollow,"/unfollow")
	following2, followers2 = FollowTestHelper(followerQuery, followedQuery)		
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, following2, "Bijon")
	assert.NotContains(t, followers2, "testingaccount")
}

