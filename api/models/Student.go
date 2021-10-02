package models

import (
	"html"
	"strings"
	"time"
)

type Student struct {
		ID 				    uint32 `gorm:"primary_key;auto_increment" json:"id"`
		FirstName	  string `gorm:"size 255; not null" json:"first_name"`
		MiddleName string `gorm:"size 255" json:"middle_name"`
		LastName	 string `gorm:"size 255; not null" json:"last_name"`
		Email 			  string `gorm:"size 100;not null;unique" json:"email"`
		EnrolledCourse 			Course `json:"enrolled_course"`
		EnrolledCourseID uint32 	`gorm:"not null" json:"enrolled_course_id"`
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
		UpdatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (student *Student) Prepare() {
		student.ID = 0
		student.FirstName = html.EscapeString(strings.TrimSpace(student.FirstName))
		student.MiddleName = html.EscapeString(strings.TrimSpace(student.MiddleName))
		student.LastName =  html.EscapeString(strings.TrimSpace(student.LastName))
		student.Email = html.EscapeString(strings.TrimSpace(student.Email))
		student.EnrolledCourse = Course{}
		student.CreatedAt = time.Now()
		student.UpdatedAt = time.Now()
}