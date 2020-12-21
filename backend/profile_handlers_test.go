package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var baseURL = "http://localhost:8000"

/*
Testing Signin
*/
func TestSignin(t *testing.T) {
	//Start Server
	go startServer()

	//Test Data
	badSigninData := []byte(`{
        "username":"fake_account"
        "password":"password"
    }`)

	OKSigninData := []byte(`{
        "username":"testingaccount",
        "password":"password"
    }`)

	//Test 1
	req, err := http.NewRequest("POST", baseURL+"/signin", bytes.NewBuffer(badSigninData))
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
	req, err = http.NewRequest("POST", baseURL+"/signin", bytes.NewBuffer(OKSigninData))
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
func CalTestHelper(data []byte) int {
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
	handler := http.HandlerFunc(Signin)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	print(resp.StatusCode)

	//TEST
	req, err = http.NewRequest("POST", baseURL+"/update_calories", bytes.NewBuffer(data))
	if err != nil {
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
func TestCalUpdate(t *testing.T) {
	//Testing Data
	calorieData1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	calorieData2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 4000
    }`)

	calorieData3 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"enjoy workoiut",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": -29
    }`)

	//Test 1
	resp1 := CalTestHelper(calorieData1)
	assert.Equal(t, 200, resp1)

	//Test 2
	resp2 := CalTestHelper(calorieData2)
	assert.Equal(t, 200, resp2)

	//Test 3
	resp3 := CalTestHelper(calorieData3)
	assert.Equal(t, 200, resp3)
}

/*
Description Test Helper
*/
func DescTestHelper(data []byte) (string, int) {
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
	handler := http.HandlerFunc(Signin)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	print(resp.StatusCode)

	//TEST
	req, err = http.NewRequest("POST", baseURL+"/update_bio", bytes.NewBuffer(data))
	if err != nil {
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
	row := db.QueryRow("select description from users where username=$1", "testingaccount")
	err = row.Scan(&desc)
	return desc, resp.StatusCode
}

/*
Test Bio Update
*/
func TestDescUpdate(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	mockData2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 2",
        "goalweight": 200,
        "bodyweight": 188,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	//Test 1
	desc1, resp1 := DescTestHelper(mockData1)
	assert.Equal(t, 200, resp1)
	assert.Equal(t, "Test Bio 1", desc1)

	//Test 2
	desc2, resp2 := DescTestHelper(mockData2)
	assert.Equal(t, 200, resp2)
	assert.Equal(t, "Test Bio 2", desc2)
}

/*
Weight Test Helper
*/
func WeightTestHelper(data []byte, query string) (int, int) {
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
	handler := http.HandlerFunc(Signin)
	handler.ServeHTTP(w, req)
	resp := w.Result()
	print(resp.StatusCode)

	//TEST
	req, err = http.NewRequest("POST", baseURL+"/update_weight", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	//Serve HTTP
	handler = http.HandlerFunc(UpdateWeights)
	handler.ServeHTTP(w, req)
	resp = w.Result()

	//Resp Body + DB Query
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var weight int
	row := db.QueryRow(query, "testingaccount")
	err = row.Scan(&weight)
	return weight, resp.StatusCode
}

/*
Test Weight Update
*/

func TestWeightsUpdate(t *testing.T) {
	//Test 1
	mockData1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 200,
        "bodyweight": 190,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	//Test 2
	mockData2 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 2",
        "goalweight": 220,
        "bodyweight": 190,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	//Test 1
	weight1, resp1 := WeightTestHelper(mockData1, "select bodyweight from users where username=$1")
	assert.Equal(t, 200, resp1)
	assert.Equal(t, 190, weight1)

	//Test 2
	weight2, resp2 := WeightTestHelper(mockData2, "select goalweight from users where username=$1")
	assert.Equal(t, 200, resp2)
	assert.Equal(t, 220, weight2)
}

/*
   Follow/Unfollow Test Helper
*/
func FollowTestHelper(data []byte, f http.HandlerFunc, route string, query1 string, query2 string) ([]string, []string, int) {
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
	handler := http.HandlerFunc(Signin)
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

	//Resp Body
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//DB Queries
	var following []string
	row := db.QueryRow(query1)
	err = row.Scan(pq.Array(&following))

	var followers []string
	row = db.QueryRow(query2)
	err = row.Scan(pq.Array(&followers))

	return following, followers, resp.StatusCode
}

/*
Test Follow/Unfollow
*/
func TestFollower(t *testing.T) {
	mockData1 := []byte(`{
        "follower":"testingaccount",
        "following":"Shardool"
    }`)

	//Test 1
	followerQuery := fmt.Sprintf("select following from users where username='%s'", "testingaccount")
	followedQuery := fmt.Sprintf("select followers from users where username='%s'", "Shardool")

	following, followers, resp1 := FollowTestHelper(mockData1, Follow, "/follow", followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Shardool")
	assert.Contains(t, followers, "testingaccount")

	following2, followers2, resp2 := FollowTestHelper(mockData1, Unfollow, "/unfollow", followerQuery, followedQuery)
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, following2, "Shardool")
	assert.NotContains(t, followers2, "testingaccount")

	mockData2 := []byte(`{
        "follower":"testingaccount",
        "following":"Bijon"
    }`)

	//Test 2
	followerQuery = fmt.Sprintf("select following from users where username='%s'", "testingaccount")
	followedQuery = fmt.Sprintf("select followers from users where username='%s'", "Bijon")

	following, followers, resp1 = FollowTestHelper(mockData2, Follow, "/follow", followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Bijon")
	assert.Contains(t, followers, "testingaccount")

	following2, followers2, resp2 = FollowTestHelper(mockData2, Unfollow, "/unfollow", followerQuery, followedQuery)
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, following2, "Bijon")
	assert.NotContains(t, followers2, "testingaccount")
}
