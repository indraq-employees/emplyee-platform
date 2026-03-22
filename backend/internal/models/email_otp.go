package models

import (
	"time"

	"employee-platform/backend/internal/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTPPurpose string

const (
	OTPPurposeAdminSignup OTPPurpose = "admin_signup"
	OTPPurposeAdminLogin  OTPPurpose = "admin_login"
)

type EmailOTP struct {
	ID            primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	Email         string                  `bson:"email" json:"email"`
	Purpose       OTPPurpose              `bson:"purpose" json:"purpose"`
	OTP           string                  `bson:"otp" json:"-"`
	ExpiresAt     time.Time               `bson:"expiresAt" json:"expiresAt"`
	ConsumedAt    *time.Time              `bson:"consumedAt,omitempty" json:"consumedAt,omitempty"`
	SignupPayload *dto.AdminSignupRequest `bson:"signupPayload,omitempty" json:"-"`
	LoginUserID   *primitive.ObjectID     `bson:"loginUserId,omitempty" json:"-"`
	CreatedAt     time.Time               `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time               `bson:"updatedAt" json:"updatedAt"`
}