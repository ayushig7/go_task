package main

import (
	"encoding/json"
	"net/http"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var emp Employee
	json.NewDecoder(r.Body).Decode(&emp)
	Database.Create(&emp)
	json.NewEncoder(w).Encode(emp)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {

}

func GetEmployeeById(w http.ResponseWriter, r *http.Request) {

}
