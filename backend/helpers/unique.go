package helpers

//Unique returns unique array
func Unique(source []string) []string {
	keys := make(map[string]bool)
	destination := []string{}
	for _, entry := range source {
	    if _, value := keys[entry]; !value {
		keys[entry] = true
		destination = append(destination, entry)
	    }
	}    
	return destination
    }
 