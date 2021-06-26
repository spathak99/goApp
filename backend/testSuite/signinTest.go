package testSuite

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"goApp/backend/server"
	signinHandlers "goApp/backend/handlers/signinHandlers"
	"github.com/stretchr/testify/assert"
)



// TestSignin Test if signin works
func SigninTest(t *testing.T) {
	//Start Server
	go server.StartServer()

	//Test 1
	badSigninData := []byte(`{
        "username":"fake_account",
        "password":"password"
    }`)

	//HTTP Request
	req, err := http.NewRequest("POST", baseURL+"/signin", bytes.NewBuffer(badSigninData))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(signinHandlers.Signin)
	handler.ServeHTTP(w, req)
	resp := w.Result()

	//Assert
	assert.Equal(t, 401, resp.StatusCode)

	//Test 2
	OKSigninData := []byte(`{
        "username":"testingaccount",
        "password":"password"
	}`)

	//HTTP Request
	req, err = http.NewRequest("POST", baseURL+"/signin", bytes.NewBuffer(OKSigninData))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	w = httptest.NewRecorder()
	handler = http.HandlerFunc(signinHandlers.Signin)
	handler.ServeHTTP(w, req)
	resp = w.Result()

	//Assert
	assert.Equal(t, 200, resp.StatusCode)
}
