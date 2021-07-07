package testSuite

import (
	"goApp/backend/db"
	"testing"
	descriptionHandlers "goApp/backend/handlers/descriptionHandlers"	
	"github.com/stretchr/testify/assert"
)

// DescTestHelper is a helper for the description update test
func DescTestHelper() string {
	var desc string
	row := db.DB.QueryRow("select description from users where username=$1", "testingaccount")
	err := row.Scan(&desc)
	if(err != nil){
		panic(err)
	}

	return desc
}

//TestDescriptionUpdate tests if the user bio update works as intended
func TestDescriptionUpdate(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "description":"Test Bio 1"
    }`)

	mockData2 := []byte(`{
        "username":"testingaccount",
        "description":"Test Bio 4"
    }`)

	//Test 1
	resp1,_ := Test_Helper(mockData1,descriptionHandlers.UpdateDescription,"/update_bio")
	desc1 := DescTestHelper()
	assert.Equal(t, 200, resp1)
	assert.Equal(t, "Test Bio 1", desc1)

	//Test 2
	resp2,_ := Test_Helper(mockData2,descriptionHandlers.UpdateDescription,"/update_bio")
	assert.Equal(t, 200, resp2)
	desc2 := DescTestHelper()
	assert.Equal(t, "Test Bio 4", desc2)
}