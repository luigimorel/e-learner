package models

import (
	"html"
	"strings"
	"time"
)

type Module struct {
	ID        			uint32 				`gorm:"primary_key;auto_increment" json:"id"`
	Content   		string 			`gorm:"text;not null" json:"module_content"`
	Author    		Tutor				  `json:"author"`
	AuthorID  		uint32 			`gorm:"not null;" json:"author_id"`
	MainCourse 	  Course 		`json:"main_course"`
	MainCourseID  uint32 	`gorm:"not null;" json:"main_course_id"`
	CreatedAt 		time.Time		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (module *Module) Prepare()  {
	module.ID = 0
	module.Content = html.EscapeString(strings.TrimSpace(module.Content))
	module.Author = Tutor{}
	module.MainCourse = Course{}
	module.CreatedAt = time.Now()
	module.UpdatedAt = time.Now()
}