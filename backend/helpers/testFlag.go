package helpers

import (
	"net/http"
	"strconv"
)

//Gets testing flag
func GetFlag(r http.Request) bool{
	//Params
	keys, _ := r.URL.Query()["test_flag"]

	//Get flag
	flag,err := strconv.ParseBool(keys[0])
	if(err != nil){
		panic(err)
	}

	return flag
}
