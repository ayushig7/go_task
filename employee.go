package main

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	EmpName   string
	EmpSalary float64
	Email     string
}
