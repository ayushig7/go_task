package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
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

	// Get the user ID from the URL parameter
	userID := mux.Vars(r)["id"]

	// Fetch the user from the database based on the provided userID
	var user User
	Database.Preload("Tweets", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC") // Sort tweets by created_at in descending order
	}).First(&user, userID)

	if user.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user.Tweets)
}
