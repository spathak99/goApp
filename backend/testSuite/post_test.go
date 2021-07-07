

package testSuite

import (
	postHandlers "goApp/backend/handlers/postHandlers"
	"testing"
	"github.com/stretchr/testify/assert"
)

//TestPosts tests is a post can be made
func TestPosts(t *testing.T) {
	mockData1 := []byte(`{
		"id":"none",
		"username":"testingaccount",
		"contents": "Mr Test Account PR'd on squat for 3 sets of 3 at 450lbs",
		"media":"www.linktomedia.xyz",
		"date":"10/02/2021",
		"likes":[]
	}`)

	//Test 1
	resp,_ := Test_Helper(mockData1, postHandlers.MakePost, "/make_post")
	assert.Equal(t, resp, 200)

	mockData2 := []byte(`{
		"id":"none",
		"username":"testingaccount",
		"contents": "Mr Test Account PR'd on deadlift for 1 set of 6  at 405lbs",
		"media":"www.linktomedia242.xyz",
		"date":"10/02/2021",
		"likes":[]
	}`)

	//Test 2
	resp,_ = Test_Helper(mockData2, postHandlers.MakePost, "/make_post")
	assert.Equal(t, resp, 200)
}