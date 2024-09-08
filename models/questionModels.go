package models

import (
	"time"
)

type Question struct {
	ID             string     `json:"id,omitempty" bson:"_id,omitempty"`
	Question       string     `json:"question,omitempty" bson:"question" validate:"required"`
	Subject        string     `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags           []string   `json:"tags,omitempty" bson:"tags"`
	Exam           string     `json:"exam,omitempty" bson:"exam"`
	Language       string     `json:"language,omitempty" bson:"language" validate:"required"`
	Difficulty     int        `json:"difficulty,omitempty" bson:"difficulty" validate:"oneof=1 2 3 4 5"`
	QuestionType   string     `json:"question_type,omitempty" bson:"question_type" validate:"oneof=m-choice m-select numeric"`
	Options        []string   `json:"options,omitempty" bson:"options" validate:"required"`
	CorrectOptions []int      `json:"correct_options,omitempty" bson:"correct_options" validate:"required"`
	Explanation    string     `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt      *time.Time `json:"created_at,omitempty" bson:"created_at"`
	EditedAt       *time.Time `json:"edited_at,omitempty" bson:"edited_at"`
	CreatedById    string     `json:"created_by,omitempty" bson:"created_by,omitempty"`
	EditedByIds    []string   `json:"edited_by,omitempty" bson:"edited_by,omitempty"`
}

type QuestionSet struct {
	ID                     string             `json:"id" bson:"_id,omitempty"`
	Name                   string             `json:"name" bson:"name"`
	QuestionsIDsAndAnswers []map[string][]int `json:"all_question_ids,omitempty" bson:"all_question_ids" validate:"required"`
	Subject                string             `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags                   []string           `json:"tags,omitempty" bson:"tags" validate:"required"`
	Exam                   string             `json:"exam,omitempty" bson:"exam"`
	TimeDuration           string             `json:"time,omitempty" bson:"time" validate:""`
	Description            string             `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt              time.Time          `json:"created_at,omitempty" bson:"created_at"`
	EditedAt               time.Time          `json:"edited_at,omitempty" bson:"edited_at"`
	CreatedById            string             `json:"created_by,omitempty" bson:"created_by,omitempty"`
	EditedByIds            []string           `json:"edited_by,omitempty" bson:"edited_by,omitempty"`
}

type TestSession struct {
	ID                     string             `json:"id" bson:"_id,omitempty"`
	Finished               bool               `json:"finished" bson:"finished"`
	Started                bool               `json:"started" bson:"started"`
	Name                   string             `json:"name" bson:"name"`
	Tags                   []string           `json:"tags" bson:"tags"`
	QuestionSetId          string             `json:"question_set_id,omitempty" bson:"question_set_id" validate:"required"`
	TakenById              string             `json:"taken_by_id,omitempty" bson:"taken_by_id" validate:"required"`
	QuestionsIDsAndAnswers []map[string][]int `json:"allQuestionsId,omitempty" bson:"allQuestionsId" validate:"required"`
	SelectedAnswers        [][]int            `json:"answered,omitempty" bson:"answered"`
	CurrentQuestionNum     int                `json:"currentQuestionNum,omitempty" bson:"currentQuestionNum" validate:"required"`
	QuestionIDsOrdered     []string           `json:"questionIDsOrdered,omitempty" bson:"questionIDsOrdered" validate:"required"`
	NCorrectlyAnswered     int                `json:"nCorrectlyAnswered,omitempty" bson:"nCorrectlyAnswered" validate:""`
	StartedTime            time.Time          `json:"started_time,omitempty" bson:"started_timet"`
	FinishedTime           time.Time          `json:"finished_time,omitempty" bson:"finished_time"`
	Mode                   string             `json:"mode,omitempty" bson:"mode" validate:"oneof=practice exam timed-practice"`
}
