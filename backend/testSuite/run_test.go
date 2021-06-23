package testSuite

import(
	testing "testing"
	"github.com/joho/godotenv"
	"os"
	"log"
)

var baseURL string 

func TestHandlers(t *testing.T){
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
      
	baseURL = os.Getenv("BaseURL")
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