package testSuite

import(testing "testing")

func TestHandlers(t *testing.T){
	SigninTest(t)
	UpdateCalorieTest(t)
	DescTest(t)
	NewsFeedTest(t)
	PersonalFeedTest(t)
	FollowTest(t)
	FuzzySearchTest(t)
	UpdateLiftsTest(t)
	LikesTest(t)
	LogTest(t)
	PostTest(t)
	CalculateMaxTest(t)
	WeightUpdateTest(t)
}