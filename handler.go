package main

import (
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"github.com/gorilla/mux"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var emp Employee
	json.NewDecoder(r.Body).Decode(&emp)
	Database.Create(&emp)
	json.NewEncoder(w).Encode(emp)
}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var employees []Employee
	Database.Find(&employees)
	json.NewEncoder(w).Encode(employees)
}

func GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var employee Employee
	Database.First(&employee, mux.Vars(r)["eid"])
	json.NewEncoder(w).Encode(employee)
}

func AddTweet(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get the user ID from the query parameter
    queryParams := r.URL.Query()
    userID := queryParams.Get("id")

    // Fetch the user from the database based on the provided userID
    var user User
    Database.Preload("Tweets").First(&user, userID)

    if user.ID == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    var tweet Tweet
    err := json.NewDecoder(r.Body).Decode(&tweet)
    if err != nil {
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // Set the UserID for the tweet
    tweet.UserID = user.ID

    // Set the CreatedAt time for the tweet
    tweet.CreatedAt = time.Now()

    // Create a new tweet and append it to the user's tweets
    user.Tweets = append(user.Tweets, tweet)

    Database.Save(&user)

    json.NewEncoder(w).Encode(user)
}

func GetUserTweets(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Get the user ID from the query parameter
    queryParams := r.URL.Query()
    userID := queryParams.Get("id")
    fmt.Println("User ID from query parameter:", userID)


    // Query the database for tweets with a matching user_id
    var tweets []Tweet
    Database.Where("user_id = ?", userID).Order("created_at DESC").Find(&tweets)

    if len(tweets) == 0 {
        http.Error(w, "No tweets found for the user", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(tweets)
}
func GetAllTweets(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Query the database to fetch all tweets along with user information
   type CustomTweet struct {
       ID         uint      `json:"id"`
       Content    string    `json:"content"`
       CreatedAt  time.Time `json:"created_at"`
       UserName   string    `json:"user_name"`
       UserEmail  string    `json:"user_email"`
       UserID  string    `json:"user_id"`
   }

   var tweetsWithUsers []CustomTweet
   Database.Table("tweets").Select("tweets.id, tweets.content, tweets.created_at, users.name as user_name,users.id as user_id, users.email as user_email").
       Joins("INNER JOIN users ON tweets.user_id = users.id").
       Order("tweets.created_at DESC").Find(&tweetsWithUsers)


    if len(tweetsWithUsers) == 0 {
        http.Error(w, "No tweets found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(tweetsWithUsers)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Query the database to fetch all tweets
    var users []User
    Database.Find(&users)

    if len(users) == 0 {
        http.Error(w, "No tweets found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(users)
}
func GetUserByID(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Get the user ID from the query parameter
    queryParams := r.URL.Query()
    userID := queryParams.Get("id")
    fmt.Println("User ID from query parameter:", userID)


    // Query the database for tweets with a matching user_id
    var users []User
    Database.Where("id = ?", userID).Find(&users)

    if len(users) == 0 {
        http.Error(w, "No tweets found for the user", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(users)
}


