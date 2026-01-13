package models

import "time"

// Client represents a shopper account.
type Client struct {
	ID            int64
	Email         *string
	Username      *string
	Phone         *string
	GoogleID      *string
	PasswordHash  *string
	DateOfBirth   *time.Time
	IsActive      bool
	TokenHash     *string
	TokenExpiry   *time.Time
	OTPCode       *string
	OTPCodeExpiry *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
