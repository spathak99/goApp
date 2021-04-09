package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var baseURL = "http://localhost:8000"

//Helper for testing
func TstHelper(data []byte, f http.HandlerFunc, route string) (int,string) {
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
	res := w.Body.String()

	return resp.StatusCode,res
}


// TestSignin Test if signin works
func TestSignin(t *testing.T) {
	//Start Server
	go startServer()

	//Test 1
	badSigninData := []byte(`{
        "username":"fake_account",
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
	assert.Equal(t, 401, resp.StatusCode)

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
	resp1,_ := TstHelper(calorieData1,UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp1)

	//Test 2
	resp2,_ := TstHelper(calorieData2,UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp2)

	//Test 3
	resp3,_ := TstHelper(calorieData3,UpdateCalories,"/update_calories")
	assert.Equal(t, 200, resp3)
}

// DescTestHelper is a helper for the description update test
func DescTestHelper() string {
	var desc string
	row := db.QueryRow("select description from users where username=$1", "testingaccount")
	err := row.Scan(&desc)
	if(err != nil){
		panic(err)
	}

	return desc
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
	resp1,_ := TstHelper(mockData1,UpdateDescription,"/update_bio")
	desc1 := DescTestHelper()
	assert.Equal(t, 200, resp1)
	assert.Equal(t, "Test Bio 1", desc1)

	//Test 2
	resp2,_ := TstHelper(mockData2,UpdateDescription,"/update_bio")
	assert.Equal(t, 200, resp2)
	desc2 := DescTestHelper()
	assert.Equal(t, "Test Bio 3", desc2)
}

// WeightTestHelper is the helper function for the weight update test
func WeightTestHelper(query1 string, query2 string) (int, int) {

	var weight int
	row := db.QueryRow(query1, "testingaccount")
	err := row.Scan(&weight)

	var goalWeight int
	row = db.QueryRow(query2, "testingaccount")
	err = row.Scan(&goalWeight)
	
	if(err != nil){
		panic(err)
	}

	return weight, goalWeight
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
	resp1,_ := TstHelper(mockData1,UpdateWeights,"/update_weight")
	weight1, goalWeight1 := WeightTestHelper(Query1, Query2)
	assert.Equal(t, 200, resp1)
	assert.Equal(t, 190, weight1)
	assert.Equal(t, 245, goalWeight1)

	//Test 2
	resp2,_ := TstHelper(mockData2,UpdateWeights,"/update_weight")
	weight2, goalWeight2 := WeightTestHelper(Query1, Query2)
	assert.Equal(t, 200, resp2)
	assert.Equal(t, 330, weight2)
	assert.Equal(t, 220, goalWeight2)
}

// FollowTestHelper is the helper function for the follow test
func FollowTestHelper(query1 string, query2 string) ([]string, []string) {
	//DB Queries
	var following []string
	row := db.QueryRow(query1)
	err := row.Scan(pq.Array(&following))

	var followers []string
	row = db.QueryRow(query2)
	err = row.Scan(pq.Array(&followers))

	if(err != nil){
		panic(err)
	}

	return following, followers
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

	resp1,_ := TstHelper(mockData1,Follow,"/follow")
	following, followers := FollowTestHelper(followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Shardool")
	assert.Contains(t, followers, "testingaccount")

	resp2,_ := TstHelper(mockData1,Unfollow,"/unfollow")
	following2, followers2 := FollowTestHelper(followerQuery, followedQuery)	
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

	resp1,_ = TstHelper(mockData2,Follow,"/follow")
	following, followers = FollowTestHelper(followerQuery, followedQuery)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, following, "Bijon")
	assert.Contains(t, followers, "testingaccount")

	resp2,_ = TstHelper(mockData2,Unfollow,"/unfollow")
	following2, followers2 = FollowTestHelper(followerQuery, followedQuery)		
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, following2, "Bijon")
	assert.NotContains(t, followers2, "testingaccount")
}

// LikesTestHelper is a helper function for the like/unlike post tests
func LikesTestHelper(query string) []string {
	//DB query
	var likes []string
	row := db.QueryRow(query)
	err := row.Scan(pq.Array(&likes))
	if(err != nil){
		panic(err)
	}

	return likes
}

// TestLikes tests liking and unliking posts
func TestLikes(t *testing.T) {
	mockData1 := []byte(`{
        "username":"testingaccount",
        "id":"5492C1CA32B7"
    }`)

	query := fmt.Sprintf("select likes from posts where id='%s'", "5492C1CA32B7")

	//Test 1
	resp1,_ := TstHelper(mockData1, LikePost, "/like_post")
	likes := LikesTestHelper(query)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, likes, "testingaccount")

	resp2,_ := TstHelper(mockData1, Unlike, "/unlike_post")
	likes2 := LikesTestHelper(query)
	assert.Equal(t, 200, resp2)
	assert.NotContains(t, likes2, "testingaccount")
}


// FuzzyTestHelper returns the query for the test
func FuzzyTestHelper(res string) []string {
	var usernames []string
	var users []Profile

	
	err := json.Unmarshal([]byte(res), &users)
	if(err != nil){
		panic(err)
	}

	for _, entry := range users {
		usernames = append(usernames, entry.Username)
	}
	return usernames
}

// TestFuzzySearch tests if the program can search for users
func TestFuzzySearch(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount",
		"query":"Shard"
	}`)

	resp1,res1 := TstHelper(mockData1, FuzzySearch, "/search")
	usernames := FuzzyTestHelper(res1)
	assert.Equal(t, 200, resp1)
	assert.Contains(t, usernames, "Shardool")
	assert.Contains(t, usernames, "Shardel")

	mockData2 := []byte(`{
		"username":"testingaccount",
		"query":"Shardool Pa"
	}`)

	resp2,res2 := TstHelper(mockData2, FuzzySearch, "/search")
	usernames2 := FuzzyTestHelper(res2)
	assert.Equal(t, 200, resp2)
	assert.Contains(t, usernames2, "Shardool")
	assert.NotContains(t, usernames2, "Shardel")
}

// FeedTestHelper helps with grabbing the news feed
func FeedTestHelper(res string) ([]string) {
	var usernames []string
	var posts []Post

	err := json.Unmarshal([]byte(res), &posts)
	if(err != nil){
		panic(err)
	}

	for _, entry := range posts {
		usernames = append(usernames, entry.Username)
	}

	return usernames
}

//TestNewsFeed checks if a feed can be grabbed for a user
func TestNewsFeed(t *testing.T) {
	mockData := []byte(`{
		"username":"testingaccount"
	}`)

	query1 := "select following from users where username='testingaccount'"
	var following []string
	row := db.QueryRow(query1)
	err := row.Scan(pq.Array(&following))
	if err != nil {
		panic(err)
	}

	resp,res := TstHelper(mockData, GetFeed, "/get_feed")
	usernames := FeedTestHelper(res)
	assert.Equal(t, resp, 200)
	for _, username := range usernames {
		assert.Contains(t, following, username)
	}
}


// LiftTestHelper helps with the lift test
func LiftTestHelper(query string) (UserLifts) {
	//DB query
	var temp UserLifts
	row := db.QueryRow(query)
	err := row.Scan(&temp.Username, &temp.Lifts)
	if(err != nil){
		panic(err)
	}
	return temp
}

//TestLiftUpdate tests if the lift dictionary was properly updated
func TestLiftUpdate(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount",
		"lifts": {
			"Deadlift": {
				"Current Max": 450,
				"Estimated Max": 475
			},
			"Squat": {
				"Current Max": 350,
				"Estimated Max": 365
			},
			"Bench": {
				"Current Max": 250,
				"Estimated Max": 260
			}
		}
	}`)
	query := "select * from userlifts where username='testingaccount'"
	resp,_ := TstHelper(mockData1, UpdateLifts, "/update_lifts")
	lifts := LiftTestHelper(query)
	assert.Equal(t, 200, resp)
	assert.Equal(t, "testingaccount", lifts.Username)

	liftmap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(lifts.Lifts), &liftmap); err != nil {
		panic(err)
	}
	dlMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Deadlift"].(map[string]interface{})["Current Max"]))
	dlERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Deadlift"].(map[string]interface{})["Estimated Max"]))
	sqMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Squat"].(map[string]interface{})["Current Max"]))
	sqERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Squat"].(map[string]interface{})["Estimated Max"]))
	bMax, _ := strconv.Atoi(fmt.Sprint(liftmap["Bench"].(map[string]interface{})["Current Max"]))
	bERM, _ := strconv.Atoi(fmt.Sprint(liftmap["Bench"].(map[string]interface{})["Estimated Max"]))

	assert.Equal(t, dlMax, 450)
	assert.Equal(t, dlERM, 475)
	assert.Equal(t, sqMax, 350)
	assert.Equal(t, sqERM, 365)
	assert.Equal(t, bMax, 250)
	assert.Equal(t, bERM, 260)
}


//TestMaxCalculator tests if a one rep max estimate is valid
func TestMaxCalculator(t *testing.T) {
	mockData1 := []byte(`{
		"weight": 405,
		"reps": 3,
		"rpe":7.5
	}`)

	resp, ret := TstHelper(mockData1, EstimateMax, "/estimate_max")
	ERM, _ := strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 462)

	mockData2 := []byte(`{
		"weight": 225,
		"reps": 3,
		"rpe":9.5
	}`)

	resp, ret = TstHelper(mockData2, EstimateMax, "/estimate_max")
	ERM, _ = strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 241)

	mockData3 := []byte(`{
		"weight": 365,
		"reps": 1,
		"rpe":10
	}`)

	resp, ret = TstHelper(mockData3, EstimateMax, "/estimate_max")
	ERM, _ = strconv.Atoi(ret)
	assert.Equal(t, resp, 200)
	assert.Equal(t, ERM, 365)
}

//TestPost tests is a post can be made
func TestPost(t *testing.T) {
	mockData1 := []byte(`{
		"id":"none",
		"username":"testingaccount",
		"contents": "Mr Test Account PR'd on squat for 3 sets of 3 at 450lbs",
		"media":"www.linktomedia.xyz",
		"date":"10/02/2021",
		"likes":[]
	}`)

	resp,_ := TstHelper(mockData1, MakePost, "/make_post")
	assert.Equal(t, resp, 200)

	mockData2 := []byte(`{
		"id":"none",
		"username":"testingaccount",
		"contents": "Mr Test Account PR'd on deadlift for 1 set of 6  at 405lbs",
		"media":"www.linktomedia242.xyz",
		"date":"10/02/2021",
		"likes":[]
	}`)

	resp,_ = TstHelper(mockData2, MakePost, "/make_post")
	assert.Equal(t, resp, 200)
}


//TestLog tests is a post can be made
func TestLog(t *testing.T) {
	mockData1 := []byte(`{
		"username":"testingaccount"
	}`)

	statusCode,_ := TstHelper(mockData1, GetLiftNames, "/get_lifts")
	assert.Equal(t, statusCode, 200)

	mockData2 := []byte(`{
		"username":"testingaccount",
		"name":"Fake Lift",
		"weight":135,
		"reps": 1,
		"sets": 3,
		"rpe": 9.5,
		"date": "09/05/2020",
		"pr":true
	}`)

	statusCode,_ = TstHelper(mockData2, LogExercise, "/logexercise")
	assert.Equal(t, statusCode, 200)
}


//TestPersonalFeed checks if the feed for the user is retrieved
func TestPersonalFeed(t *testing.T) {
	mockData := []byte(`{
		"username":"Shardool"
	}`)
	resp,_ := TstHelper(mockData, GetPersonalFeed, "/get_personal_feed")
	assert.Equal(t, resp, 200)
}