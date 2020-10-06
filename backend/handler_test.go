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


func TestSignin(t *testing.T){
    data := []byte(`{
        "username":"fake_account"
        "password":"password"
    }`)
    signin_url := "http://localhost:8000/signin"

    //Start Server
    go startServer()


    //Test1
    req, err := http.NewRequest("POST", signin_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    //defer resp.Body.Close()
    assert.Equal(t, "400 Bad Request", resp.Status) 

    
    //Test 2
    data = []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

    req, err = http.NewRequest("POST", signin_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client = &http.Client{}
    resp, err = client.Do(req)
    if err != nil {
        panic(err)
    }
    //defer resp.Body.Close()
    assert.Equal(t, "200 OK", resp.Status)
}


/*
Testing Calories Update
*/
func TestCalorieUpdate(t *testing.T){
   
    //Calories URL
    calories_url := "http://localhost:8000/update_calories"

    //Test1
    data := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 5000
    }`)
    req, err := http.NewRequest("POST", calories_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    assert.Equal(t, "200 OK", resp.Status)
    
}





