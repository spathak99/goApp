package testSuite

import (
	"goApp/backend/db"
	"testing"
	"goApp/backend/descriptionHandlers"
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

// TestDescUpdate tests if the user bio update works as intended
func TestDescUpdate(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "description":"Test Bio 1"
    }`)

	mockData2 := []byte(`{
        "username":"testingaccount",
        "description":"Test Bio 3"
    }`)

	//Test 1
	resp1,_ := TstHelper(mockData1,descriptionHandlers.UpdateDescription,"/update_bio")
	desc1 := DescTestHelper()
	assert.Equal(t, 200, resp1)
	assert.Equal(t, "Test Bio 1", desc1)

	//Test 2
	resp2,_ := TstHelper(mockData2,descriptionHandlers.UpdateDescription,"/update_bio")
	assert.Equal(t, 200, resp2)
	desc2 := DescTestHelper()
	assert.Equal(t, "Test Bio 3", desc2)
}
