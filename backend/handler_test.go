package main

import (
    "net/http"
    "testing"
    "io/ioutil"
    "net/url"
    "strings"
    "github.com/stretchr/testify/assert"
)


func TestSignIn(t *testing.T) {
    data := url.Values{
		"username":{"Shardool"},
		"password":{"Pathak"},
		"description":{"My Bio"},
		"goalweight": {"165"},
        "bodyweight": {"145"},
		"caloriegoal":{"3500"},
		"caloriesleft":{"1600"},
    }

    url := "http://localhost:8000/signin"
    client := &http.Client{}

    //Start Server
    go startServer()
    req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)

    //Check Response OK
    if err != nil {
            panic(err)
    }
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    //Check Expected Response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            panic(err)
    }
}

