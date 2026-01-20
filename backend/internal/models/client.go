package models

import "time"

// Client represents a shopper account.
type Client struct {
	ID            int64      `json:"client_id"`
	Email         *string    `json:"email"`
	Username      *string    `json:"username"`
	Phone         *string    `json:"phone"`
	GoogleID      *string    `json:"google_id"`
	PasswordHash  *string    `json:"-"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	IsActive      bool       `json:"is_active"`
	Activated     bool       `json:"activated"`
	TokenHash     *string    `json:"-"`
	TokenExpiry   *time.Time `json:"-"`
	OTPCode       *string    `json:"-"`
	OTPCodeExpiry *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
