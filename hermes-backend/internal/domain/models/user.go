package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username" gorm:"uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"not null"` // Never expose password hash
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Role         Role       `json:"role" gorm:"not null;default:'USER'"`
	Active       bool       `json:"active" gorm:"not null;default:true"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// Role represents user roles for RBAC
type Role string

// User role constants
const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleGuest Role = "GUEST"
)

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserLogin represents the data needed to login
type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserUpdateRequest represents the data that can be updated for a user
type UserUpdateRequest struct {
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *Role   `json:"role"`
	Active    *bool   `json:"active"`
}

// UserPasswordChange represents the data needed to change password
type UserPasswordChange struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// TokenResponse represents the auth token response
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        User      `json:"user"`
}
