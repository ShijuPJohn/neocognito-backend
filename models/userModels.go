package models

import "time"

type User struct {
	ID                string     `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string     `json:"name,omitempty" validate:"required,min=2,max=128" bson:"name"`
	Email             string     `json:"email,omitempty" validate:"required,email" bson:"email"`
	Password          string     `json:"password,omitempty" validate:"required,min=8,max=64" bson:"password"`
	Role              string     `json:"role" bson:"role"`
	CreatedAt         *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	PasswordChangedAt *time.Time `json:"password_changed_at,omitempty" bson:"password_changed_at"`
	Verified          bool       `json:"verified,omitempty" bson:"verified"`
	LinkedIn          string     `json:"linkedin,omitempty" bson:"linkedIn"`
	Facebook          string     `json:"facebook,omitempty" bson:"facebook"`
	Instagram         string     `json:"instagram,omitempty" bson:"instagram"`
	ProfilePic        string     `json:"profile_pic,omitempty" bson:"profile_pic"`
	About             string     `json:"about,omitempty" bson:"about"`
}
