package models

import (
	"html"
	"strings"
	"time"
)

type Course struct {
	ID            		 uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name           	string `gorm:"size 255;not null;" json:"name"`
	CourseTutor          Tutor  `json:"tutor"`
	CourseTutorID 	uint32 		`gorm:"not null" json:"tutor_id"`
	CreatedAt      	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (course *Course) Prepare()  {
	course.ID = 0
	course.Name = html.EscapeString(strings.TrimSpace(course.Name))
	course.CourseTutor = Tutor{}
	course.CreatedAt= time.Now()
	course.UpdatedAt = time.Now()
}