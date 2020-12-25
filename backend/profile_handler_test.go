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

// TestSignin Test if signin works
func TestSignin(t *testing.T) {
	//Start Server
	go startServer()

	//Test 1
	badSigninData := []byte(`{
        "username":"fake_account"
        "password":"password"
    }`)

	//HTTP Request
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
	OKSigninData := []byte(`{
        "username":"testingaccount",
        "password":"password"
	}`)

	//HTTP Request
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

// CalTestHelper is the helper for the calorie test
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

// TestCalUpdate tests if the calorie values are updated correctly
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

// DescTestHelper is a helper for the description update test
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

// TestDescUpdate tests if the user bio update works as intended
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
        "description":"Test Bio 3",
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
	assert.Equal(t, "Test Bio 3", desc2)
}

// WeightTestHelper is the helper function for the weight update test
func WeightTestHelper(data []byte, query1 string, query2 string) (int, int, int) {
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
	row := db.QueryRow(query1, "testingaccount")
	err = row.Scan(&weight)

	var goalWeight int
	row = db.QueryRow(query2, "testingaccount")
	err = row.Scan(&goalWeight)

	return weight, goalWeight, resp.StatusCode
}

// TestWeightsUpdate tests if the users weights are updated as intended
func TestWeightsUpdate(t *testing.T) {
	//Test 1
	mockData1 := []byte(`{
        "username":"testingaccount",
        "password":"password",
        "description":"Test Bio 1",
        "goalweight": 245,
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
        "bodyweight": 330,
        "caloriegoal": 4000,
        "caloriesleft": 10
    }`)

	Query1 := "select bodyweight from users where username=$1"
	Query2 := "select goalweight from users where username=$1"
	//Test 1
	weight1, goalWeight1, resp1 := WeightTestHelper(mockData1, Query1, Query2)
	assert.Equal(t, 200, resp1)
	assert.Equal(t, 190, weight1)
	assert.Equal(t, 245, goalWeight1)

	//Test 2
	weight2, goalWeight2, resp2 := WeightTestHelper(mockData2, Query1, Query2)
	assert.Equal(t, 200, resp2)
	assert.Equal(t, 330, weight2)
	assert.Equal(t, 220, goalWeight2)
}

// FollowTestHelper is the helper function for the follow test
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

// TestFollower tests if the follow and unfollow handlers work as intended
func TestFollower(t *testing.T) {
	//Test 1
	mockData1 := []byte(`{
        "follower":"testingaccount",
        "following":"Shardool"
    }`)

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

	//Test 2
	mockData2 := []byte(`{
        "follower":"testingaccount",
        "following":"Bijon"
    }`)

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

// LikesTestHelper is a helper function for the like/unlike post tests
func LikesTestHelper(data []byte, f http.HandlerFunc, route string, query string) ([]string, int) {
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

	//DB query
	var likes []string
	row := db.QueryRow(query)
	err = row.Scan(pq.Array(&likes))

	return likes, resp.StatusCode
}

// TestLikes tests liking and unliking posts
func TestLikes(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "id":"5492C1CA32B7"
    }`)

	query := fmt.Sprintf("select likes from posts where id='%s'", "5492C1CA32B7")

	//Test 1
	likes, resp := LikesTestHelper(mockData1, LikePost, "/like_post", query)
	assert.Equal(t, 200, resp)
	assert.Contains(t, likes, "testingaccount")

	likes2, resp2 := LikesTestHelper(mockData1, Unlike, "/unlike_post", query)
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, likes2, "testingaccount")
}

// CustomProgramTestHelper helps with the program test
func CustomProgramTestHelper(data []byte, f http.HandlerFunc, route string, query string) (CustomProgram, int) {
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

	//DB query
	var temp CustomProgram
	row := db.QueryRow(query, "testingaccount")
	err = row.Scan(&temp.Username, &temp.ProgramDict, pq.Array(&temp.WorkoutDays))
	return temp, resp.StatusCode
}

// TestCustomProgram tests if the program is updated correctly
func TestCustomProgram(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount",
		"programdict": {"Test_Key": "Test_Value"},
		"workoutdays":["monday","wednesday","friday"]
	}`)

	query := "select * from customprograms where username='%s'"

	program, resp := CustomProgramTestHelper(mockData1, UpdateCustomProgram, "/update_custom_program", query)
	assert.Equal(t, 200, resp)
	assert.Equal(t, "testingaccount", program.Username)
	assert.Equal(t, `{"Test_Key": "Test_Value"}`, program.ProgramDict)
	assert.Contains(t, program.WorkoutDays, "monday")
	assert.Contains(t, program.WorkoutDays, "wednesday")
	assert.Contains(t, program.WorkoutDays, "friday")
}
