package testSuite

import (
	"testing"
	"github.com/stretchr/testify/assert"
	calorieHandlers "goApp/backend/handlers/calorieHandlers"
)


//TestCalorieUpdate tests if the calorie values are updated correctly
func TestCalorieUpdate(t *testing.T) {
	//Testing Data
	calorieData1 := []byte(`{
        "username":"testingaccount",
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	calorieData2 := []byte(`{
        "username":"testingaccount",
        "caloriegoal": 4000,
        "caloriesleft": 4000
    }`)

	calorieData3 := []byte(`{
        "username":"testingaccount",
        "caloriegoal": 4000,
        "caloriesleft": -29
    }`)

	//Test 1
	resp1,_ := Test_Helper(calorieData1,calorieHandlers.UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp1)

	//Test 2
	resp2,_ := Test_Helper(calorieData2,calorieHandlers.UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp2)

	//Test 3
	resp3,_ := Test_Helper(calorieData3,calorieHandlers.UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp3)
}