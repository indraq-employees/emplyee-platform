package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/models"
	"employee-platform/backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeService struct {
	users    *mongo.Collection
	profiles *mongo.Collection
}

func NewEmployeeService(db *mongo.Database) *EmployeeService {
	return &EmployeeService{
		users:    db.Collection("users"),
		profiles: db.Collection("employee_profiles"),
	}
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, actorUserID string, req dto.CreateEmployeeRequest) (*models.User, *models.EmployeeProfile, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	count, err := s.users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return nil, nil, "", err
	}
	if count > 0 {
		return nil, nil, "", errors.New("email already exists")
	}

	role := models.RoleEmployee
	if req.Role == string(models.RoleManager) {
		role = models.RoleManager
	}

	actorID, err := primitive.ObjectIDFromHex(actorUserID)
	if err != nil {
		return nil, nil, "", err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, nil, "", err
	}

	now := time.Now()
	userID := primitive.NewObjectID()
	user := &models.User{
		ID:               userID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            email,
		PasswordHash:     hash,
		Role:             role,
		IsEmailVerified:  true,
		IsActive:         true,
		ProfileCompleted: false,
		CreatedBy:        &actorID,
		UpdatedBy:        &actorID,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	profile := &models.EmployeeProfile{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		EmployeeID: s.generateEmployeeID(),
		AddedBy:    &actorID,
		ModifiedBy: &actorID,
		AddedAt:    now,
		ModifiedAt: now,
	}

	if _, err := s.users.InsertOne(ctx, user); err != nil {
		return nil, nil, "", err
	}
	if _, err := s.profiles.InsertOne(ctx, profile); err != nil {
		_, _ = s.users.DeleteOne(ctx, bson.M{"_id": userID})
		return nil, nil, "", err
	}

	loginLink := fmt.Sprintf("http://localhost:3000/login?email=%s", email)
	return user, profile, loginLink, nil
}

func (s *EmployeeService) ListEmployees(ctx context.Context) ([]models.User, error) {
	cursor, err := s.users.Find(ctx, bson.M{"role": bson.M{"$in": []string{string(models.RoleEmployee), string(models.RoleManager)}}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *EmployeeService) generateEmployeeID() string {
	return fmt.Sprintf("EMP%s", time.Now().Format("20060102150405"))
}