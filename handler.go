package main

import (
	"encoding/json"
	"net/http"

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
