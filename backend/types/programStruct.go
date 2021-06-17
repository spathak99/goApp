package types

import (
	"encoding/json"	
)

// Program is the struct for the program that the user can choose to upload
type Program struct {
	Username    string `json:"username",db"username"`
	ProgramFile string `json:"programfile",db"programfile"`
	StartDate   string `json:"startdate",db"startdate"`
}

// CustomProgram is the struct for a program that the user can chose to create
type CustomProgram struct {
	Username    string          `json:"username",db"username"`
	ProgramList json.RawMessage `json:"programlist",db"programlist"`
}

// CustomProgram helper writes the program as a string
type CustomProgramHelper struct {
	Username    string `json:"username",db"username"`
	ProgramList string `json:"programlist",db"programlist"`
}
