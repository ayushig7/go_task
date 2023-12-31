package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandlerRouting() {
	r := mux.NewRouter()
	r.HandleFunc("/employees", GetEmployees).Methods("GET")
	r.HandleFunc("/employee", CreateEmployee).Methods("POST")
	r.HandleFunc("/employee/{eid}", GetEmployeeById).Methods("GET")
	r.HandleFunc("/addTweet", AddTweet).Methods("POST")
	r.HandleFunc("/getUserTweets", GetUserTweets).Methods("GET")
	r.HandleFunc("/getAllTweets", GetAllTweets).Methods("GET")
	r.HandleFunc("/getAllUsers", GetAllUsers).Methods("GET")
	r.HandleFunc("/getUserByID", GetUserByID).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
