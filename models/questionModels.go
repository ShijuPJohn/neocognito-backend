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
	ID                string    `json:"id" bson:"_id,omitempty"`
	Name              string    `json:"name" bson:"name"`
	QIDList           []string  `json:"qid_list,omitempty" bson:"qid_list" validate:"required"`
	CorrectAnswerList [][]int   `json:"correct_answer_list,omitempty" bson:"correct_answer_list" validate:"required"`
	MarkList          []float64 `json:"mark_list,omitempty" bson:"mark_list" validate:"required"`
	Subject           string    `json:"subject,omitempty" bson:"subject" validate:"required"`
	Tags              []string  `json:"tags,omitempty" bson:"tags" validate:"required"`
	Exam              string    `json:"exam,omitempty" bson:"exam"`
	TimeDuration      string    `json:"time,omitempty" bson:"time" validate:""`
	Description       string    `json:"explanation,omitempty" bson:"explanation"`
	CreatedAt         time.Time `json:"created_at,omitempty" bson:"created_at"`
	EditedAt          time.Time `json:"edited_at,omitempty" bson:"edited_at"`
	CreatedById       string    `json:"created_by,omitempty" bson:"created_by,omitempty"`
	EditedByIds       []string  `json:"edited_by,omitempty" bson:"edited_by,omitempty"`
}
type QuestionAnswerData struct {
	Correct             []int   `json:"correct_answer_list" bson:"correct_answer_list"`
	Selected            []int   `json:"selected_answer_list" bson:"selected_answer_list"`
	QuestionsTotalMark  float64 `json:"questions_total_mark" bson:"questions_total_mark"`
	QuestionsScoredMark float64 `json:"questions_scored_mark" bson:"questions_scored_mark"`
}
type TestSession struct {
	ID                 string                         `json:"id" bson:"_id,omitempty"`
	Finished           bool                           `json:"finished" bson:"finished"`
	TakenByID          string                         `json:"taken_by_id,omitempty" bson:"taken_by_id" validate:"required"`
	QuestionSetID      string                         `json:"question_set_id" bson:"question_set_id"`
	QuestionAnswerData map[string]*QuestionAnswerData `json:"question_answer_data" bson:"question_answer_data"`
	CurrentQuestionNum int                            `json:"current_questionNum,omitempty" bson:"current_questionNum" validate:"required"`
	QuestionIDsOrdered []string                       `json:"question_ids_ordered,omitempty" bson:"question_ids_ordered" validate:"required"`
	StartedTime        *time.Time                     `json:"started_time,omitempty" bson:"started_time"`
	FinishedTime       *time.Time                     `json:"finished_time,omitempty" bson:"finished_time"`
	Mode               string                         `json:"mode,omitempty" bson:"mode" validate:"oneof=practice exam timed-practice"`
	TotalMarks         float64                        `json:"total_marks" bson:"total_marks"`
	ScoredMarks        float64                        `json:"scored_marks" bson:"scored_marks"`
}

//UserActivity
//Error Reports
//Feedback
