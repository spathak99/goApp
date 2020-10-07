package main

import (
    "net/http"
    "testing"
    //"fmt"
    //"encoding/json"
    "io/ioutil"
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


/*
Testing Calories Update
*/
func TestCalorieUpdate(t *testing.T){

    //Signin
    signin_url := "http://localhost:8000/signin"
    data := []byte(`{
        "username":"testingaccount",
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


    calorie_url := "http://localhost:8000/update_calories"

    //Test 1
    data = []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    req, err = http.NewRequest("POST",calorie_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    handler2 := http.HandlerFunc(UpdateCalories)
    handler2.ServeHTTP(w, req)
    resp = w.Result()

    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
   
    assert.Equal(t, 200, resp.StatusCode)

    //Test 2
    data2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 4000
    }`)

    req, err = http.NewRequest("POST",calorie_url, bytes.NewBuffer(data2))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    handler2.ServeHTTP(w, req)
    resp = w.Result()

    _, err = ioutil.ReadAll(resp.Body)

    if err != nil {
        panic(err)
    }

    assert.Equal(t, 200, resp.StatusCode)


    //Test 2
    data3 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": -29
    }`)

    req, err = http.NewRequest("POST",calorie_url, bytes.NewBuffer(data3))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    handler2.ServeHTTP(w, req)
    resp = w.Result()

    _, err = ioutil.ReadAll(resp.Body)

    if err != nil {
        panic(err)
    }
    assert.Equal(t, 200, resp.StatusCode)
}



/*
Test Bio Update
*/

func TestDescriptionUpdate(t *testing.T){
        //Signin
        signin_url := "http://localhost:8000/signin"
        data := []byte(`{
            "username":"testingaccount",
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
    
    
        bio_url := "http://localhost:8000/update_bio"
    
        //Test 1
        data = []byte(`{
            "username":"testingaccount",
            "password":"password",
            "description":"Test Bio 1",
            "goalweight": 200,
            "bodyweight": 188,
            "caloriegoal": 4000,
            "caloriesleft": 10
        }`)
    
        req, err = http.NewRequest("POST",bio_url, bytes.NewBuffer(data))
        req.Header.Set("X-Custom-Header", "myvalue")
        req.Header.Set("Content-Type", "application/json")
    
        if err != nil {
            panic(err)
        }
    
        handler2 := http.HandlerFunc(UpdateDescription)
        handler2.ServeHTTP(w, req)
        resp = w.Result()
    
        _, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }
       
        assert.Equal(t, 200, resp.StatusCode)
    
        //Test 2
        data2 := []byte(`{
            "username":"testingaccount",
            "password":"password",
            "description":"Test Bio 2",
            "goalweight": 200,
            "bodyweight": 188,
            "caloriegoal": 4000,
            "caloriesleft": 10
        }`)
    
        req, err = http.NewRequest("POST",bio_url, bytes.NewBuffer(data2))
        req.Header.Set("X-Custom-Header", "myvalue")
        req.Header.Set("Content-Type", "application/json")
    
        if err != nil {
            panic(err)
        }
    
        handler2.ServeHTTP(w, req)
        resp = w.Result()
    
        _, err = ioutil.ReadAll(resp.Body)
    
        if err != nil {
            panic(err)
        }
    
        assert.Equal(t, 200, resp.StatusCode)
}