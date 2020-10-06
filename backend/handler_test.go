package main

import (
    "net/http"
    "testing"
    "bytes"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
)

/*
Testing Signin
*/


func TestSignin(t *testing.T){
    //Start Server
    signin_url := "http://localhost:8000/signin"
    go startServer()
    

    //Test 1
    data := []byte(`{
        "username":"fake_account"
        "password":"password"
    }`)

    req, err := http.NewRequest("POST",signin_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    w := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp := w.Result()

    assert.Equal(t, 400, resp.StatusCode)
   
    //Test 2
    data = []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

    req, err = http.NewRequest("POST",signin_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    w = httptest.NewRecorder()
    handler = http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp = w.Result()

    assert.Equal(t, 200, resp.StatusCode)
}






