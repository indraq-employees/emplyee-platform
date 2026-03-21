package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/models"
	"employee-platform/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	users     *mongo.Collection
	profiles  *mongo.Collection
	jwtSecret string
}

func NewAuthService(db *mongo.Database, jwtSecret string) *AuthService {
	return &AuthService{
		users:     db.Collection("users"),
		profiles:  db.Collection("employee_profiles"),
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) AdminSignup(ctx context.Context, req dto.AdminSignupRequest) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	count, err := s.users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return nil, "", err
	}
	if count > 0 {
		return nil, "", errors.New("email already exists")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	user := &models.User{
		ID:               primitive.NewObjectID(),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            email,
		PasswordHash:     hash,
		Role:             models.RoleAdmin,
		IsEmailVerified:  true,
		IsActive:         true,
		ProfileCompleted: true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	_, err = s.users.InsertOne(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), string(user.Role), s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*models.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if !user.IsActive {
		return nil, "", errors.New("account is inactive")
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	now := time.Now()
	_, _ = s.users.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"lastLoginAt": now,
			"updatedAt":   now,
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