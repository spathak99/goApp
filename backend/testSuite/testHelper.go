
package testSuite

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"goApp/backend/signinHandlers"
)


//Helper for testing
func TstHelper(data []byte, f http.HandlerFunc, route string) (int,string) {
	var baseURL = "http://localhost:8000"

	//Signin
	signinData := []byte(`{
		"username":"testingaccount",
		"password":"password"
	}`)

	//Request
	req, err := http.NewRequest("POST", baseURL+"/signin", bytes.NewBuffer(signinData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(signinHandlers.Signin)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	print(resp.StatusCode)

	//TEST
	req, err = http.NewRequest("POST", baseURL+route, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	handler = http.HandlerFunc(f)
	handler.ServeHTTP(w, req)
	resp = w.Result()
	res := w.Body.String()

	return resp.StatusCode,res
}