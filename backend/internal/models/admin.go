package models

import "time"

// Admin represents an administrative user.
type Admin struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	IsActive     bool
	LastLogin    *time.Time
	TokenHash    *string
	TokenExpiry  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
