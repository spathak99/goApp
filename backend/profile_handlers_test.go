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
Signs in for the handler functions that occur after login
*/
func signin_helper(){
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
}



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
Test Bio Update
*/

func TestDescUpdate(t *testing.T){
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

        var desc string
        row := db.QueryRow("select description from users where username=$1","testingaccount")
        err = row.Scan(&desc)
        assert.Equal(t,desc,"Test Bio 1")

    
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

        var desc2 string
        row = db.QueryRow("select description from users where username=$1","testingaccount")
        err = row.Scan(&desc2)
        assert.Equal(t,desc2,"Test Bio 2")
}

/*
Test Weight Update
*/

func TestWeightsUpdate(t *testing.T){
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


    bio_url := "http://localhost:8000/update_weight"

    //Test 1
    data = []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 200,
        "bodyweight": 190,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

    req, err = http.NewRequest("POST",bio_url, bytes.NewBuffer(data))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    if err != nil {
        panic(err)
    }

    handler2 := http.HandlerFunc(UpdateWeights)
    handler2.ServeHTTP(w, req)
    resp = w.Result()

    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
   
    assert.Equal(t, 200, resp.StatusCode)

    var weight int
    row := db.QueryRow("select bodyweight from users where username=$1","testingaccount")
    err = row.Scan(&weight)
    assert.Equal(t,weight,190)


    //Test 2
    data2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 2",
        "goalweight": 220,
        "bodyweight": 190,
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

    var weight2 int
    row = db.QueryRow("select goalweight from users where username=$1","testingaccount")
    err = row.Scan(&weight2)
    assert.Equal(t,weight2,220)
}