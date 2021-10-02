package models

import (
	"html"
	"strings"
	"time"
)

type Tutor struct {
	ID       		 uint32        `gorm:"primary_key;auto_increment" json:"id"`
	FirstName  	string        `gorm:"size 255; not null" json:"first_name"`
	MiddleName string `gorm:"size 255" json:"middle_name"`
	LastName  string        `gorm:"size 255; not null" json:"last_name"`
	Email    		string        `gorm:"size 100;not null;unique" json:"email"`
	CreatedAt 	time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (tutor *Tutor) Prepare() {
	tutor.ID = 0
	tutor.FirstName = html.EscapeString(strings.TrimSpace(tutor.FirstName))
	tutor.MiddleName= html.EscapeString(strings.TrimSpace(tutor.MiddleName))
	tutor.LastName = html.EscapeString(strings.TrimSpace(tutor.LastName))
	tutor.Email = html.EscapeString(strings.TrimSpace(tutor.LastName))
	tutor.CreatedAt = time.Now()
	tutor.UpdatedAt	= time.Now()
}