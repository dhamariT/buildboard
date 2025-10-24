package models

import (
	"time"

	"gorm.io/gorm"
)

// EarlyStartUser represents a user who has signed up for early access.
type EarlyStartUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Email     string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	FirstName string `gorm:"size:100" json:"firstName,omitempty"`
	LastName  string `gorm:"size:100" json:"lastName,omitempty"`
	IPAddress string `gorm:"size:45" json:"ipAddress"` // IPv6 max length is 45 chars

	// OTP fields for email verification
	OTPHash        string     `gorm:"size:6" json:"-"`                 // 6-character OTP (not exposed in JSON)
	OTPExpiresAt   *time.Time `json:"-"`                               // OTP expiration time
	OTPVerifiedAt  *time.Time `json:"otpVerifiedAt,omitempty"`         // When OTP was verified
	OTPAttempts    int        `gorm:"default:0" json:"-"`              // Failed OTP verification attempts
	OTPLastAttempt *time.Time `json:"-"`                               // Last OTP verification attempt
	IsVerified     bool       `gorm:"default:false" json:"isVerified"` // Whether email is verified

	// Email engagement tracking
	EmailSent       bool       `gorm:"default:false" json:"emailSent"`
	EngagementToken string     `gorm:"uniqueIndex;size:64" json:"-"` // Unique token for tracking
	ReadAt          *time.Time `json:"readAt,omitempty"`             // First time email was read
	ReadCount       int        `gorm:"default:0" json:"readCount"`   // Number of times email was viewed
	LastReadAt      *time.Time `json:"lastReadAt,omitempty"`         // Most recent view timestamp
	ReaderIP        string     `gorm:"size:45" json:"-"`             // IP address from first read
	ReaderClient    string     `gorm:"size:500" json:"-"`            // Email client from first read
}

// BeforeCreate hook is called before creating a new record.
// It can be used to set default values or perform validation.
func (e *EarlyStartUser) BeforeCreate(tx *gorm.DB) error {
	// Add any pre-creation logic here if needed
	return nil
}

// AfterCreate hook is called after creating a new record.
// It can be used to trigger background tasks or notifications.
func (e *EarlyStartUser) AfterCreate(tx *gorm.DB) error {
	// Add any post-creation logic here if needed
	return nil
}
