package server

import (
	"net/http"
	"github.com/rs/cors"
	"goApp/backend/db"
	signinHandlers "goApp/backend/handlers/signinHandlers"
	weightHandlers "goApp/backend/handlers/weightHandlers"
	descriptionHandlers "goApp/backend/handlers/descriptionHandlers"
	calorieHandlers "goApp/backend/handlers/calorieHandlers"
	followHandlers "goApp/backend/handlers/followHandlers"
	postHandlers "goApp/backend/handlers/postHandlers"
	fuzzySearch "goApp/backend/handlers/fuzzySearch"
	customProgramHandlers "goApp/backend/handlers/customProgramHandlers"
	likeHandlers "goApp/backend/handlers/likeHandlers"
	userHandlers "goApp/backend/handlers/userHandlers"
	feedHandlers "goApp/backend/handlers/feedHandlers"
	liftHandlers "goApp/backend/handlers/liftHandlers"
	_ "github.com/lib/pq"
)


//StartServer begins the server
func StartServer() {
	print("Starting Server")
	mux:= http.NewServeMux()

	//Routes
	mux.HandleFunc("/signin", signinHandlers.Signin)
	mux.HandleFunc("/signup", signinHandlers.Signup)
	mux.HandleFunc("/logout", signinHandlers.Logout)
	mux.HandleFunc("/get_all_users", userHandlers.GetUsers)
	mux.HandleFunc("/update_bio", descriptionHandlers.UpdateDescription)
	mux.HandleFunc("/update_weight", weightHandlers.UpdateWeights)
	mux.HandleFunc("/update_calories", calorieHandlers.UpdateCalories)
	mux.HandleFunc("/get_user_data", userHandlers.GetUserData)
	mux.HandleFunc("/get_followers", followHandlers.GetFollowers)
	mux.HandleFunc("/get_following", followHandlers.GetFollowing)
	mux.HandleFunc("/follow", followHandlers.Follow)
	mux.HandleFunc("/unfollow", followHandlers.Unfollow)
	mux.HandleFunc("/make_post",postHandlers.MakePost)
	mux.HandleFunc("/get_feed", feedHandlers.GetFeed)
	mux.HandleFunc("/like_post", likeHandlers.LikePost)
	mux.HandleFunc("/unlike_post", likeHandlers.Unlike)
	mux.HandleFunc("/initial_custom_program", customProgramHandlers.InitializeProgram)
	mux.HandleFunc("/update_custom_program", customProgramHandlers.UpdateCustomProgram)
	mux.HandleFunc("/get_custom_program", customProgramHandlers.GetCustomProgram)
	mux.HandleFunc("/get_post", postHandlers.GetPost)
	mux.HandleFunc("/get_personal_feed", feedHandlers.GetPersonalFeed)
	mux.HandleFunc("/search", fuzzySearch.FuzzySearch)
	mux.HandleFunc("/update_name", descriptionHandlers.UpdateName)
	mux.HandleFunc("/initialize_lifts", liftHandlers.InitializeLifts)
	mux.HandleFunc("/get_user_max", liftHandlers.GetUserMax)
	mux.HandleFunc("/update_lifts", liftHandlers.UpdateLifts)
	mux.HandleFunc("/estimate_max", liftHandlers.EstimateMax)
	mux.HandleFunc("/logexercise", liftHandlers.LogExercise)
	mux.HandleFunc("/get_lifts",liftHandlers.GetLiftNames)
	mux.HandleFunc("/grablog", liftHandlers.GrabLog)
	
	//Launch
	db.InitDB()
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8000", handler)
}
