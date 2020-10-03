package main

import (
    "net/http"
    "testing"
    "bytes"
    "github.com/stretchr/testify/assert"
)

/*
Testing Signin
*/
func TestSignIn(t *testing.T) {
    data := []byte(`{
        "username":"sunofthemoon",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 200
    }`)
    signin_url := "http://localhost:8000/signin"

    //Start Server
    go startServer()
    req, err := http.NewRequest("POST", signin_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    assert.Equal(t, "200 OK", resp.Status)
}



