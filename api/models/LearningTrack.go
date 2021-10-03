package models

import (
	"html"
	"strings"
	"time"
)

type LearningTrack struct {
	ID        		uint32   `gorm:"primary_key;auto_increment;" json:"id"`
	Name      	string   `gorm:"size 255;not null;unique" json:"name"`
	Course 	  Course   `json:"course"`
	CourseID 		uint32	`gorm:"not null" json:"course_id"`
	Assessor       Tutor `json:"tutor"`
	Learner			Student	`json:"student"` 
	LearnerID 		uint32 `gorm:"not null" json:"learner_id"`
	CreatedAt 		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (learningTrack *LearningTrack) Prepare()  {
	learningTrack.ID = 0
	learningTrack.Name = html.EscapeString(strings.TrimSpace(learningTrack.Name))
	learningTrack.Course = Course{}
	learningTrack.Assessor = Tutor{}
	learningTrack.CreatedAt = time.Now()
	learningTrack.UpdatedAt = time.Now()
}