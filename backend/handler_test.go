package main

import (
    "net/http"
    "bytes"
    "net/http/httptest"
    "testing"
)


func TestSignIn(t *testing.T) {
    jsonData := []byte(`
		"username":"Shardool",
		"password":"Pathak",
		"description":"My Bio",
		"goalweight": 165,
		"bodyweight": 188,
		"caloriegoal":3500,
		"caloriesleft":1600
    }`)

    req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonData))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	print(rr.Body.String())
}

func TestSignIn(t *testing.T) {
    jsonData := []byte(`
		"username":"Shardool",
		"password":"Pathak",
		"description":"My Bio",
		"goalweight": 165,
		"bodyweight": 188,
		"caloriegoal":3500,
		"caloriesleft":1600
    }`)

    req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonData))

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Signin)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

	print(rr.Body.String())
}