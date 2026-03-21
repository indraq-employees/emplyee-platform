package services

import (
	"context"
	"time"

	"employee-platform/backend/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProfileService struct {
	users    *mongo.Collection
	profiles *mongo.Collection
}

func NewProfileService(db *mongo.Database) *ProfileService {
	return &ProfileService{
		users:    db.Collection("users"),
		profiles: db.Collection("employee_profiles"),
	}
}

func (s *ProfileService) CompleteProfile(ctx context.Context, userID string, req dto.CompleteProfileRequest) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"nickName":             req.NickName,
		"department":           req.Department,
		"location":             req.Location,
		"designation":          req.Designation,
		"jobRole":              req.JobRole,
		"employmentType":       req.EmploymentType,
		"employeeStatus":       req.EmployeeStatus,
		"sourceOfHire":         req.SourceOfHire,
		"maritalStatus":        req.MaritalStatus,
		"aboutMe":              req.AboutMe,
		"expertise":            req.Expertise,
		"uan":                  req.UAN,
		"pan":                  req.PAN,
		"workPhoneNumber":      req.WorkPhoneNumber,
		"personalMobileNumber": req.PersonalMobileNumber,
		"extension":            req.Extension,
		"personalEmailAddress": req.PersonalEmailAddress,
		"seatingLocation":      req.SeatingLocation,
		"tags":                 req.Tags,
		"presentAddress":       req.PresentAddress,
		"permanentAddress":     req.PermanentAddress,
		"modifiedAt":           now,
	}

	if req.DateOfJoining != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateOfJoining); err == nil {
			update["dateOfJoining"] = parsed
		}
	}
	if req.DateOfBirth != "" {
		if parsed, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			update["dateOfBirth"] = parsed
		}
	}

	_, err = s.profiles.UpdateOne(ctx, bson.M{"userId": objID}, bson.M{"$set": update})
	if err != nil {
		return err
	}

	_, err = s.users.UpdateByID(ctx, objID, bson.M{"$set": bson.M{"profileCompleted": true, "updatedAt": now}})
	return err
}