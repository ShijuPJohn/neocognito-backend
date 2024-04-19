package models

import "time"

type User struct {
	ID                      string              `json:"id,omitempty" bson:"_id,omitempty"`
	Name                    string              `json:"name,omitempty" validate:"required,min=2,max=128" bson:"name"`
	Email                   string              `json:"email,omitempty" validate:"required,email" bson:"email"`
	Password                string              `json:"password,omitempty" validate:"required,min=8,max=64" bson:"password"`
	Role                    string              `json:"role" bson:"role"`
	CreatedAt               time.Time           `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt               time.Time           `json:"updatedAt,omitempty" bson:"updatedAt"`
	PasswordChangedAt       time.Time           `json:"passwordChangedAt,omitempty" bson:"passwordChangedAt"`
	Verified                bool                `json:"verified,omitempty" bson:"verified"`
	LinkedIn                string              `json:"linkedIn,omitempty" bson:"linkedIn"`
	Facebook                string              `json:"facebook,omitempty" bson:"facebook"`
	Instagram               string              `json:"instagram,omitempty" bson:"instagram"`
	ProfilePic              string              `json:"profilePic,omitempty" bson:"profilePic"`
	About                   string              `json:"about,omitempty" bson:"about"`
	SectionScore            []map[string]string `json:"sectionScore,omitempty" bson:"sectionScore"`
	DailyAttemptedQuestions map[string][]string `json:"dailyAttemptedQuestions,omitempty" bson:"dailyAttemptedQuestions"`
}

type UserSignup struct {
	ID              string `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string `json:"name,omitempty" validate:"required" bson:"name"`
	Email           string `json:"email,omitempty" validate:"required,email" bson:"email"`
	Password        string `json:"password,omitempty" validate:"required" bson:"password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required" bson:"password"`
	LinkedIn        string `json:"linkedIn,omitempty" bson:"linkedIn"`
	Facebook        string `json:"facebook,omitempty" bson:"facebook"`
	Instagram       string `json:"instagram,omitempty" bson:"instagram"`
	ProfilePic      string `json:"profilePic,omitempty" bson:"profilePic"`
	About           string `json:"about,omitempty" bson:"about"`
}
