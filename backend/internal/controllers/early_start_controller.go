package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/buildboard/backend/internal/models"
	"github.com/buildboard/backend/internal/services"
)

// EarlyStartController handles early access signup and verification.
type EarlyStartController struct {
	DB           *gorm.DB
	EmailService *services.EmailService
}

// NewEarlyStartController creates a new early start controller instance.
func NewEarlyStartController(db *gorm.DB, emailService *services.EmailService) *EarlyStartController {
	return &EarlyStartController{
		DB:           db,
		EmailService: emailService,
	}
}

type signupRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type verifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

const (
	maxOTPAttempts    = 5
	otpValidityPeriod = 15 * time.Minute
)

// generateOTP creates a 6-character alphanumeric OTP.
// It excludes visually similar characters for better user experience.
func generateOTP() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Exclude similar looking characters
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	otp := make([]byte, 6)
	for i := 0; i < 6; i++ {
		otp[i] = charset[int(b[i])%len(charset)]
	}
	return string(otp), nil
}

// validateEmailFormat validates email address format.
func validateEmailFormat(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// Signup handles early start signup requests.
// It creates or updates a user record and sends an OTP email for verification.
func (e *EarlyStartController) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate email format
	if err := validateEmailFormat(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	// Check if email already exists
	var existing models.EarlyStartUser
	err := e.DB.Where("email = ?", req.Email).First(&existing).Error

	if err == nil {
		// Email exists
		if existing.IsVerified {
			// Use generic message for security - don't reveal if email is already verified
			c.JSON(http.StatusOK, gin.H{"message": "If this email is valid, a verification code has been sent"})
			return
		}
		// Email exists but not verified - allow resending OTP
		// Check rate limiting on OTP attempts
		if existing.OTPLastAttempt != nil && time.Since(*existing.OTPLastAttempt) < 1*time.Minute {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Please wait before requesting a new code"})
			return
		}
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	// Generate OTP
	otp, err := generateOTP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}
	otpHash := string(otp) // Store as-is (6 characters)

	// Get client IP address
	clientIP := c.ClientIP()

	// Generate engagement tracking token
	engagementToken, err := services.GenerateEngagementToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	// Send OTP email
	emailSent := false
	if e.EmailService != nil {
		err := e.EmailService.SendOTPEmail(req.Email, otp, engagementToken)
		if err == nil {
			emailSent = true
		}
	}

	// Calculate OTP expiration
	expiresAt := time.Now().Add(otpValidityPeriod)
	now := time.Now()

	if existing.ID != 0 {
		// Update existing record with new OTP
		updates := map[string]interface{}{
			"first_name":       req.FirstName,
			"last_name":        req.LastName,
			"otp_hash":         otpHash,
			"otp_expires_at":   expiresAt,
			"otp_attempts":     0,
			"otp_last_attempt": now,
			"email_sent":       emailSent,
			"engagement_token": engagementToken,
			"read_count":       0,
			"read_at":          nil,
			"last_read_at":     nil,
			"reader_ip":        nil,
			"reader_client":    nil,
		}
		e.DB.Model(&existing).Updates(updates)
	} else {
		// Create new early start user record
		user := models.EarlyStartUser{
			Email:           req.Email,
			FirstName:       req.FirstName,
			LastName:        req.LastName,
			EmailSent:       emailSent,
			IPAddress:       clientIP,
			OTPHash:         otpHash,
			OTPExpiresAt:    &expiresAt,
			OTPAttempts:     0,
			OTPLastAttempt:  &now,
			IsVerified:      false,
			EngagementToken: engagementToken,
		}

		if err := e.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
			return
		}
	}

	// Generic success message for security
	c.JSON(http.StatusOK, gin.H{
		"message": "If this email is valid, a verification code has been sent. Please check your email.",
	})
}

// VerifyOTP verifies the OTP and completes the early start signup.
func (e *EarlyStartController) VerifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Normalize OTP to uppercase
	otp := strings.ToUpper(strings.TrimSpace(req.OTP))
	if len(otp) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Find user by email
	var user models.EarlyStartUser
	if err := e.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Use generic error message for security
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Check if already verified
	if user.IsVerified {
		c.JSON(http.StatusOK, gin.H{
			"message":    "Email already verified",
			"isVerified": true,
		})
		return
	}

	// Check OTP expiration
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code has expired. Please request a new one."})
		return
	}

	// Check max attempts
	if user.OTPAttempts >= maxOTPAttempts {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many failed attempts. Please request a new code."})
		return
	}

	// Verify OTP
	if otp != user.OTPHash {
		// Increment failed attempts
		now := time.Now()
		e.DB.Model(&user).Updates(map[string]interface{}{
			"otp_attempts":     user.OTPAttempts + 1,
			"otp_last_attempt": now,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// OTP is valid - mark as verified
	verifiedAt := time.Now()
	updates := map[string]interface{}{
		"is_verified":      true,
		"otp_verified_at":  verifiedAt,
		"otp_hash":         "", // Clear OTP
		"otp_expires_at":   nil,
		"otp_attempts":     0,
		"otp_last_attempt": nil,
	}

	if err := e.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Service temporarily unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Email verified successfully!",
		"email":         user.Email,
		"isVerified":    true,
		"otpVerifiedAt": verifiedAt,
	})
}

// Count returns the total number of early start signups.
func (e *EarlyStartController) Count(c *gin.Context) {
	var totalCount int64
	var verifiedCount int64

	if err := e.DB.Model(&models.EarlyStartUser{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get signup count"})
		return
	}

	if err := e.DB.Model(&models.EarlyStartUser{}).Where("is_verified = ?", true).Count(&verifiedCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get verified count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    totalCount,
		"verified": verifiedCount,
		"message":  "early start signups",
	})
}

// List returns paginated list of early start signups (admin endpoint).
func (e *EarlyStartController) List(c *gin.Context) {
	var users []models.EarlyStartUser

	// Parse pagination parameters
	page := 1
	limit := 50

	// Query with pagination
	offset := (page - 1) * limit
	if err := e.DB.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve signups"})
		return
	}

	// Get total count
	var total int64
	e.DB.Model(&models.EarlyStartUser{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// TrackEngagement tracks email engagement (open tracking pixel).
func (e *EarlyStartController) TrackEngagement(c *gin.Context) {
	filename := c.Param("filename")

	// Extract token from filename (remove .png extension)
	token := strings.TrimSuffix(filename, ".png")

	// Find user by engagement token
	var user models.EarlyStartUser
	if err := e.DB.Where("engagement_token = ?", token).First(&user).Error; err != nil {
		// Return 1x1 transparent pixel even if token not found
		c.Data(http.StatusOK, "image/png", []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
			0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
			0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
			0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		})
		return
	}

	// Update engagement tracking
	now := time.Now()
	updates := map[string]interface{}{
		"read_count":   user.ReadCount + 1,
		"last_read_at": now,
	}

	// Set first read time if not already set
	if user.ReadAt == nil {
		updates["read_at"] = now
		updates["reader_ip"] = c.ClientIP()
		updates["reader_client"] = c.GetHeader("User-Agent")
	}

	e.DB.Model(&user).Updates(updates)

	// Return 1x1 transparent PNG pixel
	c.Data(http.StatusOK, "image/png", []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	})
}
