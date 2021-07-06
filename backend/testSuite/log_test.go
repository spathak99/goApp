package testSuite

import (
	"testing"
	liftHandlers "goApp/backend/handlers/liftHandlers"
	"github.com/stretchr/testify/assert"
)

//TestLog tests is a post can be made
func TestLog(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount"
	}`)

	//Test1
	statusCode,_ := Test_Helper(mockData1, liftHandlers.GetLiftNames, "/get_lifts")
	assert.Equal(t, statusCode, 200)

	mockData2 := []byte(`{
		"username":"testingaccount",
		"name":"Fake Lift",
		"weight":135,
		"reps": 1,
		"sets": 3,
		"rpe": 9.5,
		"date": "09/05/2020",
		"pr":true
	}`)

	//Test 2
	statusCode,_ = Test_Helper(mockData2,liftHandlers.LogExercise, "/logexercise")
	assert.Equal(t, statusCode, 200)
}