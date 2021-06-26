
package testSuite

import (
	"goApp/backend/types"
	"testing"
	"encoding/json"
	fuzzySearch"goApp/backend/handlers/fuzzySearch"
	"github.com/stretchr/testify/assert"
)


// FuzzyTestHelper returns the query for the test
func FuzzyTestHelper(res string) []string {
	var usernames []string
	var users []types.Profile

	
	err := json.Unmarshal([]byte(res), &users)
	if(err != nil){
		panic(err)
	}

	for _, entry := range users {
		usernames = append(usernames, entry.Username)
	}
	return usernames
}

// FuzzySearchTest tests if the program can search for users
func FuzzySearchTest(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount",
		"query":"Shard"
	}`)

	//Test 1
	resp1,res1 := TstHelper(mockData1, fuzzySearch.FuzzySearch, "/search")
	usernames := FuzzyTestHelper(res1)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, usernames, "Shardool")
	assert.Contains(t, usernames, "Shardel")

	mockData2 := []byte(`{
		"username":"testingaccount",
		"query":"Shardool Pa"
	}`)

	//Test 2
	resp2,res2 := TstHelper(mockData2, fuzzySearch.FuzzySearch, "/search")
	usernames2 := FuzzyTestHelper(res2)
	assert.Equal(t, 200, resp2)
	assert.Contains(t, usernames2, "Shardool")
	assert.NotContains(t, usernames2, "Shardel")
}
