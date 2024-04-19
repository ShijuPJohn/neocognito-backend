package models

import (
	"time"
)

type Question struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Question string `json:"question,omitempty" bson:"question" validate:"required"`
	//Category       string    `json:"category,omitempty" bson:"category" validate:"required"`
	Subject        string    `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags           []string  `json:"tags,omitempty" bson:"tags"`
	Exam           string    `json:"exam,omitempty" bson:"exam"`
	Language       string    `json:"language,omitempty" bson:"language" validate:"required"`
	Difficulty     int       `json:"difficulty,omitempty" bson:"difficulty"`
	QuestionType   string    `json:"questionType,omitempty" bson:"questionType" validate:"oneof=m-choice m-select numeric"`
	Options        []string  `json:"options,omitempty" bson:"options" validate:"required"`
	CorrectOptions int       `json:"correctOptions,omitempty" bson:"correctOptions" validate:"required"`
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
	Mode          string    `json:"mode,omitempty" bson:"mode" validate:"required,oneof=practice exam timed"`
	Subject       string    `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags          []string  `json:"tags,omitempty" bson:"tags" validate:"required"`
	Exam          string    `json:"exam,omitempty" bson:"exam"`
	Language      string    `json:"language,omitempty" bson:"language" validate:"required"`
	TimeDuration  string    `json:"time,omitempty" bson:"time" validate:""`
	Difficulty    int       `json:"difficulty,omitempty" bson:"difficulty"`
	Description   string    `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	EditedAt      time.Time `json:"editedAt,omitempty"`
	CreatedById   string    `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	EditedByIds   []string  `json:"editedBy,omitempty" bson:"editedBy,omitempty"`
	TotalAttempts int       `json:"totalAttempts,omitempty" bson:"totalAttempts"`
	MarksObtained []int     `json:"marksObtained,omitempty" bson:"marksObtained"`
}

type QTest struct {
	ID                 string           `json:"id" bson:"_id,omitempty"`
	Finished           bool             `json:"finished" bson:"finished"`
	Started            bool             `json:"started" bson:"started"`
	Name               string           `json:"name" bson:"name"`
	Tags               []string         `json:"tags" bson:"tags"`
	QuestionSetId      string           `json:"questionSetId,omitempty" bson:"questionSetId" validate:"required"`
	TakenById          string           `json:"takenById,omitempty" bson:"takenById" validate:"required"`
	NTotalQuestions    int              `json:"nTotalQuestions,omitempty" bson:"nTotalQuestions" validate:"required"`
	AllQuestionsIDs    map[string][]int `json:"allQuestionsId,omitempty" bson:"allQuestionsId" validate:"required"`
	CurrentQuestionNum int              `json:"currentQuestionNum,omitempty" bson:"currentQuestionNum" validate:"required"`
	QuestionIDsOrdered []string         `json:"questionIDsOrdered,omitempty" bson:"questionIDsOrdered" validate:"required"`
	//AnsweredQuestionsIDs []string       `json:"answeredQuestionsId,omitempty" bson:"answeredQuestionsId" validate:"required"`
	//NTotalAnswered       int            `json:"nTotalAnswered,omitempty" bson:"nTotalAnswered" validate:"required"`
	NCorrectlyAnswered int       `json:"nCorrectlyAnswered,omitempty" bson:"nCorrectlyAnswered" validate:""`
	Rank               int       `json:"rank,omitempty" bson:"rank" validate:""`
	TakenAtTime        time.Time `json:"takenAt,omitempty" bson:"takenAt"`
	Mode               string    `json:"mode,omitempty" bson:"mode" validate:"oneof=practice exam timed-practice"`
}
