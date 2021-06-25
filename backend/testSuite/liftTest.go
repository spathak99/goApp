package testSuite

import (
	"goApp/backend/db"
	"goApp/backend/types"
	"testing"
	"fmt"
	"strconv"
	"encoding/json"
	"goApp/backend/liftHandlers"
	"github.com/stretchr/testify/assert"
)

// LiftTestHelper helps with the lift test
func LiftTestHelper(query string) (types.UserLifts) {
	//DB query
	var temp types.UserLifts
	row := db.DB.QueryRow(query)
	err := row.Scan(&temp.Username, &temp.Lifts)
	if(err != nil){
		panic(err)
	}
	return temp
}

//UpdateLiftsTest tests if the lift dictionary was properly updated
func UpdateLiftsTest(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount",
		"lifts": {
			"Deadlift": {
				"Current Max": 450,
				"Estimated Max": 475
			},
			"Squat": {
				"Current Max": 350,
				"Estimated Max": 365
			},
			"Bench": {
				"Current Max": 250,
				"Estimated Max": 260
			}
		}
	}`)
	query := "select * from userlifts where username='testingaccount'"
	resp,_ := TstHelper(mockData1, liftHandlers.UpdateLifts, "/update_lifts")
	lifts := LiftTestHelper(query)

	//Test 1
	assert.Equal(t, 200, resp)
	assert.Equal(t, "testingaccount", lifts.Username)

	liftmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(lifts.Lifts), &liftmap); err != nil {
		panic(err)
	}
	dlMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Deadlift"].(map[string]interface{})["Current Max"]))
	dlERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Deadlift"].(map[string]interface{})["Estimated Max"]))
	sqMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Squat"].(map[string]interface{})["Current Max"]))
	sqERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Squat"].(map[string]interface{})["Estimated Max"]))
	bMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Bench"].(map[string]interface{})["Current Max"]))
	bERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Bench"].(map[string]interface{})["Estimated Max"]))

	//Test 2
	assert.Equal(t, dlMax, 450)
	assert.Equal(t, dlERM, 475)
	assert.Equal(t, sqMax, 350)
	assert.Equal(t, sqERM, 365)
	assert.Equal(t, bMax, 250)
	assert.Equal(t, bERM, 260)
}
