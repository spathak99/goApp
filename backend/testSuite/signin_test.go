package testSuite

import (
	"testing"
	"github.com/stretchr/testify/assert"
	signinHandler "goApp/backend/handlers/signinHandlers"
)

//Test Signin tests if signin is done properly
func TestSignin(t *testing.T){
	//Test 1
	badSigninData := []byte(`{
		"username":"fake_account",
		"password":"password"
	}`)

	//Test 1
	resp,_ := Test_Helper(badSigninData, signinHandler.Signin, "/signin")
	assert.Equal(t,401,resp)

	OKSigninData := []byte(`{
		"username":"testingaccount",
		"password":"password"
	}`)
	//Test 1

	resp,_ = Test_Helper(OKSigninData, signinHandler.Signin, "/signin")
	assert.Equal(t,200,resp)
}