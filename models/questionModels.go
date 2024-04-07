package models

import (
	"time"
)

type Question struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	Question       string    `json:"question,omitempty" bson:"question" validate:"required"`
	Category       string    `json:"category,omitempty" bson:"category" validate:"required"`
	Subject        string    `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags           []string  `json:"tags,omitempty" bson:"tags"`
	Exam           string    `json:"exam,omitempty" bson:"exam"`
	Language       string    `json:"language,omitempty" bson:"language" validate:"required"`
	Difficulty     int       `json:"difficulty,omitempty" bson:"difficulty"`
	QuestionType   string    `json:"questionType,omitempty" bson:"questionType" validate:"oneof=m-choice m-select numeric"`
	Options        []string  `json:"options,omitempty" bson:"options" validate:"required"`
	CorrectOptions []int     `json:"correctOptions,omitempty" bson:"correctOptions" validate:"required"`
	Explanation    string    `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	EditedAt       time.Time `json:"editedAt,omitempty"`
	CreatedById    string    `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	EditedByIds    []string  `json:"editedBy,omitempty" bson:"editedBy,omitempty"`
}

type QuestionSet struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Name          string    `json:"name" bson:"name"`
	Questions     []string  `json:"questions" bson:"questions" validate:""`
	Category      string    `json:"category,omitempty" bson:"category" validate:"required"`
	Subject       string    `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags          []string  `json:"tags,omitempty" bson:"tags" validate:"required"`
	Exam          string    `json:"exam,omitempty" bson:"exam"`
	Language      string    `json:"language,omitempty" bson:"language" validate:"required"`
	Difficulty    int       `json:"difficulty,omitempty" bson:"difficulty"`
	Explanation   string    `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	EditedAt      time.Time `json:"editedAt,omitempty"`
	CreatedById   string    `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	EditedByIds   []string  `json:"editedBy,omitempty" bson:"editedBy,omitempty"`
	TotalAttempts int       `json:"totalAttempts,omitempty" bson:"totalAttempts"`
	MarksObtained []int     `json:"marksObtained,omitempty" bson:"marksObtained"`
}
