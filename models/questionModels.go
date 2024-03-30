package models

import "time"

type Question struct {
	ID             string    `json:"id" bson:"_id,omitempty"`
	Question       string    `json:"question" bson:"question" validate:"required"`
	MediaPresent   string    `json:"mediaPresent" bson:"mediaPresent"` //Type of Media Present if any
	MediaURL       string    `json:"mediaURL" bson:"mediaURL"`
	Category       string    `json:"category" bson:"category" validate:"required"`
	Subject        string    `json:"subject" bson:"subject" validate:"required"`
	Exam           string    `json:"exam" bson:"exam"`
	Difficulty     int       `json:"difficulty" bson:"difficulty"`
	QuestionType   string    `json:"questionType" bson:"questionType" validate:"oneof=m-choice m-select numeric"`
	Options        []string  `json:"options" bson:"options" validate:"required"`
	CorrectOptions []int     `json:"correctOptions" bson:"correctOptions" validate:"required"`
	Explanation    string    `json:"explanation" bson:"explanation"`
	CreatedAt      time.Time `json:"createdAt"`
	EditedAt       time.Time `json:"editedAt"`
	CreatedById    string    `json:"createdBy" bson:"createdBy"`
	EditedByIds    []string  `json:"editedBy" bson:"editedBy"`
}
