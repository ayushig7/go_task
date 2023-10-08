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
    gorm.Model
    Name   string  `json:"name" gorm:"column:name"`
    Email  string  `json:"email" gorm:"column:email"`
    Tweets []Tweet `gorm:"foreignkey:UserID"`
}

type Tweet struct {
    gorm.Model
    Content   string    `json:"content" gorm:"column:content"`
    UserID    uint      `json:"-"`
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

