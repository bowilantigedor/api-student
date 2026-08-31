package main

import "time"

type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Major     string    `json:"major"`
	GPA       float64   `json:"gpa"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStudentRequest struct {
	NIM   string  `json:"nim"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Major string  `json:"major"`
	GPA   float64 `json:"gpa"`
}

type ReplaceStudentRequest struct {
	NIM   string  `json:"nim"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Major string  `json:"major"`
	GPA   float64 `json:"gpa"`
}

type PatchStudentRequest struct {
	NIM   *string  `json:"nim"`
	Name  *string  `json:"name"`
	Email *string  `json:"email"`
	Major *string  `json:"major"`
	GPA   *float64 `json:"gpa"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

type WebResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
