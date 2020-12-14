package main

import (
    "net/http"
    "testing"
    "io/ioutil"
    "bytes"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
)

var base_url = "http://localhost:8000"

/*
Testing Signin
*/
func TestSignin(t *testing.T){
    //Start Server
    go startServer()
    
    //Test Data
    Bad_Signin_Data := []byte(`{
        "username":"fake_account"
        "password":"password"
    }`)

    OK_Signin_Data := []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)


    //Test 1
    req, err := http.NewRequest("POST",base_url + "/signin", bytes.NewBuffer(Bad_Signin_Data))
    if err != nil {
        t.Error(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp := w.Result()

    //Assert
    assert.Equal(t, 400, resp.StatusCode)
   
    

    //Test 2
    req, err = http.NewRequest("POST",base_url + "/signin", bytes.NewBuffer(OK_Signin_Data))
    if err != nil {
        t.Error(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    w = httptest.NewRecorder()
    handler = http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp = w.Result()

    //Assert
    assert.Equal(t, 200, resp.StatusCode)
}


/*
Calorie Test Helper
*/
func CalTestHelper(data []byte) int{
    //Signin         
    signin_data := []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

    //Request
    req, err := http.NewRequest("POST",base_url + "/signin", bytes.NewBuffer(signin_data))
    if err != nil {
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp := w.Result()
    print(resp.StatusCode)


    //TEST 
    req, err = http.NewRequest("POST",base_url + "/update_calories", bytes.NewBuffer(data))
    if err != nil{
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    handler = http.HandlerFunc(UpdateCalories)
    handler.ServeHTTP(w, req)
    resp = w.Result()

    //Resp Body
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    //Assert
    return resp.StatusCode
}



/*
Testing Calories Update
*/
func TestCalUpdate(t *testing.T){
    //Testing Data
    Calorie_Data_1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    Calorie_Data_2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 4000
    }`)

    Calorie_Data_3 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": -29
    }`)


    //Test 1
    resp1 := CalTestHelper(Calorie_Data_1)
    assert.Equal(t, 200, resp1)

    //Test 2
    resp2 := CalTestHelper(Calorie_Data_2)
    assert.Equal(t, 200, resp2)

    //Test 3
    resp3 := CalTestHelper(Calorie_Data_3)
    assert.Equal(t, 200, resp3)
}


/*
Description Test Helper
*/
func DescTestHelper(data []byte) (string,int){
    //Signin         
    signin_data := []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

    //Request
    req, err := http.NewRequest("POST",base_url + "/signin", bytes.NewBuffer(signin_data))
    if err != nil {
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp := w.Result()
    print(resp.StatusCode)


    //TEST 
    req, err = http.NewRequest("POST",base_url + "/update_bio", bytes.NewBuffer(data))
    if err != nil{
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    handler = http.HandlerFunc(UpdateDescription)
    handler.ServeHTTP(w, req)
    resp = w.Result()

    //Resp Body
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    var desc string
    row := db.QueryRow("select description from users where username=$1","testingaccount")
    err = row.Scan(&desc)
    return desc,resp.StatusCode
}


/*
Test Bio Update
*/
func TestDescUpdate(t *testing.T){
    Mock_Data_1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    Mock_Data_2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 2",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    //Test 1
    desc1,resp1 := DescTestHelper(Mock_Data_1)
    assert.Equal(t, 200, resp1)
    assert.Equal(t,"Test Bio 1",desc1)

    //Test 2
    desc2,resp2 := DescTestHelper(Mock_Data_2)
    assert.Equal(t, 200, resp2)
    assert.Equal(t,"Test Bio 2",desc2)
}


/*
Weight Test Helper
*/
func WeightTestHelper(data []byte) (int,int){
    //Signin         
    signin_data := []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

    //Request
    req, err := http.NewRequest("POST",base_url + "/signin", bytes.NewBuffer(signin_data))
    if err != nil {
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    w := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(w, req)
    resp := w.Result()
    print(resp.StatusCode)


    //TEST 
    req, err = http.NewRequest("POST",base_url + "/update_weight", bytes.NewBuffer(data))
    if err != nil{
        panic(err)
    }
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    //Serve HTTP
    handler = http.HandlerFunc(UpdateWeights)
    handler.ServeHTTP(w, req)
    resp = w.Result()

    //Resp Body
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    var weight int
    row := db.QueryRow("select bodyweight from users where username=$1","testingaccount")
    err = row.Scan(&weight)
    return weight,resp.StatusCode
}



/*
Test Weight Update
*/

func TestWeightsUpdate(t *testing.T){
    //Test 1
    Mock_Data_1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 200,
        "bodyweight": 190,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    //Test 2
    Mock_Data_2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 2",
        "goalweight": 220,
        "bodyweight": 190,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    //Test 1
    weight1,resp1 := WeightTestHelper(Mock_Data_1)
    assert.Equal(t, 200, resp1)
    assert.Equal(t,190,weight1)

    //Test 2
    weight2,resp2 := WeightTestHelper(Mock_Data_2)
    assert.Equal(t, 200, resp2)
    assert.Equal(t,220,weight2)
}