
package testSuite

import (
	"testing"
	"strconv"
	"goApp/backend/liftHandlers"
	"github.com/stretchr/testify/assert"
)



//CalculateMaxTest tests if a one rep max estimate is valid
func CalculateMaxTest(t *testing.T) {
	mockData1 := []byte(`{
		"weight": 405,
		"reps": 3,
		"rpe":7.5
	}`)

	//Test 1
	resp, ret := TstHelper(mockData1, liftHandlers.EstimateMax, "/estimate_max")
	ERM, _ := strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 462)

	mockData2 := []byte(`{
		"weight": 225,
		"reps": 3,
		"rpe":9.5
	}`)

	//Test 2
	resp, ret = TstHelper(mockData2, liftHandlers.EstimateMax, "/estimate_max")
	ERM, _ = strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 241)

	mockData3 := []byte(`{
		"weight": 365,
		"reps": 1,
		"rpe":10
	}`)

	//Test 3
	resp, ret = TstHelper(mockData3, liftHandlers.EstimateMax, "/estimate_max")
	ERM, _ = strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 365)
}