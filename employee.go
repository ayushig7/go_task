package main

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	EmpName   string  `json:"empname"`
	EmpSalary float64 `json:"salary"`
	Email     string  `json:"email"`
}

type User struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Tweets []Tweet `json:"tweets"`
}

type Tweet struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	UserID    int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
