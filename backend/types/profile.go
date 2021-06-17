package types

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
	Name         string   `json:"name",db"name"`
}