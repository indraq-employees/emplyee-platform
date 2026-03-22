package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/models"
	"employee-platform/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService struct {
	users     *mongo.Collection
	profiles  *mongo.Collection
	emailOTPs *mongo.Collection
	jwtSecret string
	mailer    *MailerService
}

func NewAuthService(db *mongo.Database, jwtSecret string, mailer *MailerService) *AuthService {
	return &AuthService{
		users:     db.Collection("users"),
		profiles:  db.Collection("employee_profiles"),
		emailOTPs: db.Collection("email_otps"),
		jwtSecret: jwtSecret,
		mailer:    mailer,
	}
}

func (s *AuthService) SendAdminSignupOTP(ctx context.Context, req dto.AdminSignupRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	count, err := s.users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}
	if !s.mailer.IsConfigured() {
		return errors.New("smtp is not configured")
	}

	otp, err := generateOTP()
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.emailOTPs.UpdateOne(
		ctx,
		bson.M{"email": email, "purpose": models.OTPPurposeAdminSignup},
		bson.M{
			"$set": bson.M{
				"email":         email,
				"purpose":       models.OTPPurposeAdminSignup,
				"otp":           otp,
				"expiresAt":     now.Add(10 * time.Minute),
				"consumedAt":    nil,
				"signupPayload": req,
				"updatedAt":     now,
			},
			"$setOnInsert": bson.M{
				"createdAt": now,
			},
		},
		optionsUpsert(),
	)
	if err != nil {
		return err
	}

	return s.mailer.SendHTML(email, "Admin signup verification OTP", fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;line-height:1.6">
			<h2>Verify your admin signup</h2>
			<p>Your OTP is <strong style="font-size:22px;letter-spacing:4px;">%s</strong></p>
			<p>This OTP expires in 10 minutes.</p>
		</div>`, otp))
}

func (s *AuthService) VerifyAdminSignupOTP(ctx context.Context, req dto.VerifyAdminSignupOTPRequest) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var otpDoc models.EmailOTP
	if err := s.emailOTPs.FindOne(ctx, bson.M{
		"email":   email,
		"purpose": models.OTPPurposeAdminSignup,
	}).Decode(&otpDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", errors.New("otp request not found")
		}
		return nil, "", err
	}

	if otpDoc.ConsumedAt != nil {
		return nil, "", errors.New("otp already used")
	}
	if time.Now().After(otpDoc.ExpiresAt) {
		return nil, "", errors.New("otp expired")
	}
	if strings.TrimSpace(req.OTP) != otpDoc.OTP {
		return nil, "", errors.New("invalid otp")
	}
	if otpDoc.SignupPayload == nil {
		return nil, "", errors.New("signup payload not found")
	}

	count, err := s.users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return nil, "", err
	}
	if count > 0 {
		return nil, "", errors.New("email already exists")
	}

	hash, err := utils.HashPassword(otpDoc.SignupPayload.Password)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	user := &models.User{
		ID:               primitive.NewObjectID(),
		FirstName:        otpDoc.SignupPayload.FirstName,
		LastName:         otpDoc.SignupPayload.LastName,
		Email:            email,
		PasswordHash:     hash,
		Role:             models.RoleAdmin,
		IsEmailVerified:  true,
		IsActive:         true,
		ProfileCompleted: true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if _, err = s.users.InsertOne(ctx, user); err != nil {
		return nil, "", err
	}

	_, _ = s.emailOTPs.UpdateByID(ctx, otpDoc.ID, bson.M{
		"$set": bson.M{
			"consumedAt": now,
			"updatedAt":  now,
		},
	})

	token, err := utils.GenerateJWT(user.ID.Hex(), string(user.Role), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) SendLoginOTP(ctx context.Context, req dto.LoginRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("invalid credentials")
		}
		return err
	}

	if !user.IsActive {
		return errors.New("account is inactive")
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return errors.New("invalid credentials")
	}

	if !s.mailer.IsConfigured() {
		return errors.New("smtp is not configured")
	}

	otp, err := generateOTP()
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = s.emailOTPs.UpdateOne(
		ctx,
		bson.M{"email": email, "purpose": models.OTPPurposeAdminLogin},
		bson.M{
			"$set": bson.M{
				"email":       email,
				"purpose":     models.OTPPurposeAdminLogin,
				"otp":         otp,
				"expiresAt":   now.Add(10 * time.Minute),
				"consumedAt":  nil,
				"loginUserId": user.ID,
				"updatedAt":   now,
			},
			"$setOnInsert": bson.M{
				"createdAt": now,
			},
		},
		optionsUpsert(),
	)
	if err != nil {
		return err
	}

	return s.mailer.SendHTML(email, "Login verification OTP", fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;line-height:1.6">
			<h2>Verify your login</h2>
			<p>Your OTP is <strong style="font-size:22px;letter-spacing:4px;">%s</strong></p>
			<p>This OTP expires in 10 minutes.</p>
		</div>`, otp))
}

func (s *AuthService) VerifyLoginOTP(ctx context.Context, req dto.VerifyLoginOTPRequest) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var otpDoc models.EmailOTP
	if err := s.emailOTPs.FindOne(ctx, bson.M{
		"email":   email,
		"purpose": models.OTPPurposeAdminLogin,
	}).Decode(&otpDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", errors.New("otp request not found")
		}
		return nil, "", err
	}

	if otpDoc.ConsumedAt != nil {
		return nil, "", errors.New("otp already used")
	}
	if time.Now().After(otpDoc.ExpiresAt) {
		return nil, "", errors.New("otp expired")
	}
	if strings.TrimSpace(req.OTP) != otpDoc.OTP {
		return nil, "", errors.New("invalid otp")
	}
	if otpDoc.LoginUserID == nil {
		return nil, "", errors.New("login request is invalid")
	}

	var user models.User
	if err := s.users.FindOne(ctx, bson.M{"_id": otpDoc.LoginUserID}).Decode(&user); err != nil {
		return nil, "", err
	}

	now := time.Now()
	_, _ = s.users.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
		},
	})
	_, _ = s.emailOTPs.UpdateByID(ctx, otpDoc.ID, bson.M{
		"$set": bson.M{
			"consumedAt": now,
			"updatedAt":  now,
		},
	})

	user.LastLoginAt = &now

	token, err := utils.GenerateJWT(user.ID.Hex(), string(user.Role), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*models.User, *models.EmployeeProfile, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, nil, err
	}

	var user models.User
	if err := s.users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		return nil, nil, err
	}

	var profile models.EmployeeProfile
	err = s.profiles.FindOne(ctx, bson.M{"userId": objID}).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &user, nil, nil
		}
		return nil, nil, err
	}

	return &user, &profile, nil
}

func generateOTP() (string, error) {
	var builder strings.Builder
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		builder.WriteString(n.String())
	}
	return builder.String(), nil
}

func optionsUpsert() *options.UpdateOptions {
	upsert := true
	return &options.UpdateOptions{Upsert: &upsert}
}