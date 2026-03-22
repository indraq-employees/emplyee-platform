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
)

type EmployeeService struct {
	users       *mongo.Collection
	profiles    *mongo.Collection
	mailer      *MailerService
	frontendURL string
}

func NewEmployeeService(db *mongo.Database, frontendURL string, mailer *MailerService) *EmployeeService {
	return &EmployeeService{
		users:       db.Collection("users"),
		profiles:    db.Collection("employee_profiles"),
		mailer:      mailer,
		frontendURL: frontendURL,
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
	if !s.mailer.IsConfigured() {
		return nil, nil, "", errors.New("smtp is not configured")
	}

	role := models.RoleEmployee
	if req.Role == string(models.RoleManager) {
		role = models.RoleManager
	}

	actorID, err := primitive.ObjectIDFromHex(actorUserID)
	if err != nil {
		return nil, nil, "", err
	}

	plainPassword := strings.TrimSpace(req.Password)
	if plainPassword == "" {
		plainPassword, err = generateTemporaryPassword(10)
		if err != nil {
			return nil, nil, "", err
		}
	}

	hash, err := utils.HashPassword(plainPassword)
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

	loginLink := fmt.Sprintf("%s/login?email=%s&invited=1", strings.TrimRight(s.frontendURL, "/"), email)
	emailBody := fmt.Sprintf(`
		<div style="font-family:Arial,sans-serif;line-height:1.6">
			<h2>You have been added as %s</h2>
			<p>Hello %s %s,</p>
			<p>Your account has been created successfully.</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Temporary Password:</strong> %s</p>
			<p><a href="%s" style="display:inline-block;padding:12px 18px;background:#111;color:#fff;text-decoration:none;border-radius:8px;">Go to Employee Login</a></p>
			<p>Please login and complete your profile.</p>
		</div>`,
		role,
		req.FirstName,
		req.LastName,
		email,
		plainPassword,
		loginLink,
	)

	if err := s.mailer.SendHTML(email, "Your employee account has been created", emailBody); err != nil {
		_, _ = s.profiles.DeleteOne(ctx, bson.M{"_id": profile.ID})
		_, _ = s.users.DeleteOne(ctx, bson.M{"_id": userID})
		return nil, nil, "", err
	}

	return user, profile, loginLink, nil
}

func (s *EmployeeService) ListEmployees(ctx context.Context) ([]models.User, error) {
	cursor, err := s.users.Find(ctx, bson.M{
		"role": bson.M{
			"$in": []string{string(models.RoleEmployee), string(models.RoleManager)},
		},
	})
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

func generateTemporaryPassword(length int) (string, error) {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789@#"

	var builder strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		builder.WriteByte(chars[n.Int64()])
	}
	return builder.String(), nil
}