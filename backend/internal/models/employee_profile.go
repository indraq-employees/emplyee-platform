package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmployeeProfile struct {
	ID                  primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID              primitive.ObjectID  `bson:"userId" json:"userId"`
	EmployeeID          string              `bson:"employeeId" json:"employeeId"`
	NickName            string              `bson:"nickName,omitempty" json:"nickName,omitempty"`
	Department          string              `bson:"department,omitempty" json:"department,omitempty"`
	Location            string              `bson:"location,omitempty" json:"location,omitempty"`
	Designation         string              `bson:"designation,omitempty" json:"designation,omitempty"`
	JobRole             string              `bson:"jobRole,omitempty" json:"jobRole,omitempty"`
	EmploymentType      string              `bson:"employmentType,omitempty" json:"employmentType,omitempty"`
	EmployeeStatus      string              `bson:"employeeStatus,omitempty" json:"employeeStatus,omitempty"`
	SourceOfHire        string              `bson:"sourceOfHire,omitempty" json:"sourceOfHire,omitempty"`
	DateOfJoining       *time.Time          `bson:"dateOfJoining,omitempty" json:"dateOfJoining,omitempty"`
	ReportingManagerID  *primitive.ObjectID `bson:"reportingManagerId,omitempty" json:"reportingManagerId,omitempty"`
	DateOfBirth         *time.Time          `bson:"dateOfBirth,omitempty" json:"dateOfBirth,omitempty"`
	MaritalStatus       string              `bson:"maritalStatus,omitempty" json:"maritalStatus,omitempty"`
	AboutMe             string              `bson:"aboutMe,omitempty" json:"aboutMe,omitempty"`
	Expertise           string              `bson:"expertise,omitempty" json:"expertise,omitempty"`
	UAN                 string              `bson:"uan,omitempty" json:"uan,omitempty"`
	PAN                 string              `bson:"pan,omitempty" json:"pan,omitempty"`
	WorkPhoneNumber     string              `bson:"workPhoneNumber,omitempty" json:"workPhoneNumber,omitempty"`
	PersonalMobileNumber string             `bson:"personalMobileNumber,omitempty" json:"personalMobileNumber,omitempty"`
	Extension           string              `bson:"extension,omitempty" json:"extension,omitempty"`
	PersonalEmailAddress string             `bson:"personalEmailAddress,omitempty" json:"personalEmailAddress,omitempty"`
	SeatingLocation     string              `bson:"seatingLocation,omitempty" json:"seatingLocation,omitempty"`
	Tags                []string            `bson:"tags,omitempty" json:"tags,omitempty"`
	PresentAddress      string              `bson:"presentAddress,omitempty" json:"presentAddress,omitempty"`
	PermanentAddress    string              `bson:"permanentAddress,omitempty" json:"permanentAddress,omitempty"`
	AddedBy             *primitive.ObjectID `bson:"addedBy,omitempty" json:"addedBy,omitempty"`
	ModifiedBy          *primitive.ObjectID `bson:"modifiedBy,omitempty" json:"modifiedBy,omitempty"`
	AddedAt             time.Time           `bson:"addedAt" json:"addedAt"`
	ModifiedAt          time.Time           `bson:"modifiedAt" json:"modifiedAt"`
}