package types


// Lift is struct used for logging trends and calculating maxes
type Lift struct {
	Username string
	Name     string
	Weight   int
	Reps     int
	Sets     int
	RPE      float64
	Date     string
	PR       bool
}


