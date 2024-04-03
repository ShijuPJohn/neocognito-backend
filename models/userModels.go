package models

import "time"

type User struct {
	ID                string              `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string              `json:"name" validate:"required,min=2,max=128" bson:"name"`
	Email             string              `json:"email" validate:"required,email" bson:"email"`
	Password          string              `json:"password" validate:"required,min=8,max=64" bson:"password"`
	Role              string              `json:"role" bson:"role"`
	CreatedAt         time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt" bson:"updatedAt"`
	PasswordChangedAt time.Time           `json:"passwordChangedAt" bson:"passwordChangedAt"`
	Verified          bool                `json:"verified" bson:"verified"`
	LinkedIn          string              `json:"linkedIn,omitempty" bson:"linkedIn"`
	Facebook          string              `json:"facebook,omitempty" bson:"facebook"`
	Instagram         string              `json:"instagram,omitempty" bson:"instagram"`
	ProfilePic        string              `json:"profilePic,omitempty" bson:"profilePic"`
	About             string              `json:"about,omitempty" bson:"about"`
	SectionScore      []map[string]string `json:"sectionScore,omitempty" bson:"sectionScore"`
}

type UserSignup struct {
	ID              string `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string `json:"name" validate:"required" bson:"name"`
	Email           string `json:"email" validate:"required,email" bson:"email"`
	Password        string `json:"password" validate:"required" bson:"password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required" bson:"password"`
	LinkedIn        string `json:"linkedIn,omitempty" bson:"linkedIn"`
	Facebook        string `json:"facebook,omitempty" bson:"facebook"`
	Instagram       string `json:"instagram,omitempty" bson:"instagram"`
	ProfilePic      string `json:"profilePic,omitempty" bson:"profilePic"`
	About           string `json:"about,omitempty" bson:"about"`
}
