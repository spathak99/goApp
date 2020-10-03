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
        "username":"testingaccount",
        "password":"password"
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

    //Test2
    data2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 15
    }`)
    req2, err := http.NewRequest("POST", calories_url, bytes.NewBuffer(data2))
    req2.Header.Set("X-Custom-Header", "myvalue")
    req2.Header.Set("Content-Type", "application/json")

    client2 := &http.Client{}
    resp2, err := client2.Do(req2)
    if err != nil {
        panic(err)
    }
    assert.Equal(t, "200 OK", resp2.Status)


    //Test3
    data3 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 300
    }`)
    req3, err := http.NewRequest("POST", calories_url, bytes.NewBuffer(data3))
    req3.Header.Set("X-Custom-Header", "myvalue")
    req3.Header.Set("Content-Type", "application/json")

    client3 := &http.Client{}
    resp3, err := client3.Do(req3)
    if err != nil {
        panic(err)
    }
    assert.Equal(t, "200 OK", resp3.Status)
}





