package main

import "encoding/json"

// Profile is the user struct
type Profile struct {
	Username     string   `json:"username", db:"username"`
	Password     string   `json:"password", db:"password"`
	Description  string   `json:"description",db:"description"`
	GoalWeight   float32  `json:"goalweight", db:"goalweight"`
	Bodyweight   float32  `json:"bodyweight", db:"bodyweight"`
	CalorieGoal  float32  `json:"caloriegoal",db"caloriegoal`
	CaloriesLeft float32  `json:"caloriesleft",db"caloriesleft`
	Followers    []string `json:"followers",db"followers"`
	Following    []string `json:"following",db"following"`
	Program      string   `json:"program",db"program"`
}

// Post struct
type Post struct {
	ID       string   `json:"id", db:"id"`
	Username string   `json:"username",db"username"`
	Contents string   `json:"contents",db:"contents"`
	Media    string   `json:"media",db:"media`
	Date     string   `json:"date",db:"date"`
	Likes    []string `json:"likes",db:"likes"`
}

// Program is the struct for the program that the user can choose to upload
type Program struct {
	Username    string `json:"username",db"username"`
	ProgramFile string `json:"programfile",db"programfile"`
	StartDate   string `json:"startdate",db"startdate"`
}

// CustomProgramHelper is the struct that does not communicate with the DB
type CustomProgramHelper struct {
	Username    string
	ProgramDict json.RawMessage
	WorkoutDays []string
}

// CustomProgram is the struct for a program that the user can chose to create
type CustomProgram struct {
	Username    string   `json:"username",db"username"`
	ProgramDict string   `json:"programdict",db"programdict"`
	WorkoutDays []string `json:"workoutdays",db"workoutdays"`
}

// FollowRelation is the struct for follower/following relationship
type FollowRelation struct {
	Follower  string
	Following string
}

// Like Struct
type Like struct {
	Username string
	ID       string
}
