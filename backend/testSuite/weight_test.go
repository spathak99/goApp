package testSuite

import (
	"goApp/backend/db"
	"testing"
	weightHandlers "goApp/backend/handlers/weightHandlers"
	"github.com/stretchr/testify/assert"
)

// WeightTestHelper is the helper function for the weight update test
func WeightTestHelper(query1 string, query2 string) (int, int) {

	var weight int
	row := db.DB.QueryRow(query1, "testingaccount")
	err := row.Scan(&weight)

	var goalWeight int
	row = db.DB.QueryRow(query2, "testingaccount")
	err = row.Scan(&goalWeight)
	
	if(err != nil){
		panic(err)
	}

	return weight, goalWeight
}

// TestWeightUpdate tests if the users weights are updated as intended
func TestWeightUpdate(t *testing.T) {
	//Test 1
	mockData1 := []byte(`{
        "username":"testingaccount",
        "goalweight": 245,
        "bodyweight": 190
    }`)

	//Test 2
	mockData2 := []byte(`{
        "username":"testingaccount",
        "goalweight": 220,
        "bodyweight": 330
    }`)

	Query1 := "select bodyweight from users where username=$1"
	Query2 := "select goalweight from users where username=$1"
	//Test 1
	resp1,_ := Test_Helper(mockData1,weightHandlers.UpdateWeights,"/update_weight")
	weight1, goalWeight1 := WeightTestHelper(Query1, Query2)
	assert.Equal(t, 200, resp1)
	assert.Equal(t, 190, weight1)
	assert.Equal(t, 245, goalWeight1)

	//Test 2
	resp2,_ := Test_Helper(mockData2,weightHandlers.UpdateWeights,"/update_weight")
	weight2, goalWeight2 := WeightTestHelper(Query1, Query2)
	assert.Equal(t, 200, resp2)
	assert.Equal(t, 330, weight2)
	assert.Equal(t, 220, goalWeight2)
}



