package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleAdmin      UserRole = "admin"
	RoleManager    UserRole = "manager"
	RoleEmployee   UserRole = "employee"
)

type User struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	FirstName        string              `bson:"firstName" json:"firstName"`
	LastName         string              `bson:"lastName" json:"lastName"`
	Email            string              `bson:"email" json:"email"`
	PasswordHash     string              `bson:"passwordHash" json:"-"`
	Role             UserRole            `bson:"role" json:"role"`
	IsEmailVerified  bool                `bson:"isEmailVerified" json:"isEmailVerified"`
	IsActive         bool                `bson:"isActive" json:"isActive"`
	ProfileCompleted bool                `bson:"profileCompleted" json:"profileCompleted"`
	LastLoginAt      *time.Time          `bson:"lastLoginAt,omitempty" json:"lastLoginAt,omitempty"`
	CreatedBy        *primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedBy        *primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy,omitempty"`
	CreatedAt        time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time           `bson:"updatedAt" json:"updatedAt"`
}