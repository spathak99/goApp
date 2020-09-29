package main

import (
    "net/http"
    "testing"
    "strings"
    "io/ioutil"
    "net/url"
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

    signin_url := "http://localhost:8000/signin"


    //Start Server
    go startServer()
    client := &http.Client{}
    req, _ := http.NewRequest("POST",signin_url, strings.NewReader(data.Encode()))
    resp, err := client.Do(req)
    
    //Check Response OK
    if err != nil {
            panic(err)
    }
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    //Check Expected Response
    _, err := ioutil.ReadAll(resp.Body)
    if err != nil {
            panic(err)
    }
}


