package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"` // Never output in JSON responses
	Role      string    `json:"role" validate:"required,oneof=starthub investor donator collaborator"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterUserRequest represents the request body for user registration
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=starthub investor donator collaborator"`
}

// UserResponse represents the safe user data returned in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
