
package testSuite

import (	
	"bytes"
	"net/http"
	"net/http/httptest"
	"gopkg.in/h2non/gock.v1"
	"goApp/backend/server"
)


//Helper for testing
func TstHelper(data []byte, f http.HandlerFunc, route string) (int,string) {
	go server.StartServer()

	OKSigninData := []byte(`{
		"username":"testingaccount",
		"password":"password"	
	}`)

	defer gock.Off()
	
	gock.New(baseURL).
		Post("/signin").
		Reply(200).
		JSON(OKSigninData)
    
	req, _ := http.Post(baseURL+"/signin", "application/json",bytes.NewBuffer(OKSigninData))
	print(req.StatusCode)

	
	//Request
	res,err := http.NewRequest("POST", baseURL+route, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	res.Header.Set("X-Custom-Header", "myvalue")
	res.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	handler.ServeHTTP(w, res)
	resp := w.Result()
	resBody := w.Body.String()

	return resp.StatusCode,resBody
}