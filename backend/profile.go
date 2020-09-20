package main

/*
	User Field
*/
type Profile struct {
	Password string `json:"password", db:"password"`
	Description string `json:"description",db:"description"`
	Username string `json:"username", db:"username"`
	GoalWeight float32 `json:"goalweight", db:"goalweight"` 
	Bodyweight float32 `json:"bodyweight", db:"bodyweight"` 
	CalorieGoal float32 `json:"caloriegoal",db"caloriegoal`
}