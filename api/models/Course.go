package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Course struct {
	ID            			 uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name           		string `gorm:"size 255;not null;" json:"name"`
	Title 					string	`gorm:"size 255;not null;" json:"title"`
	Creator          Tutor  `json:"tutor"`
	CreatorID 			uint32 		`gorm:"not null" json:"tutor_id"`
	CreatedAt      		time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 			time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (course *Course) Prepare()  {
	course.ID = 0
	course.Name = html.EscapeString(strings.TrimSpace(course.Name))
	course.Title = html.EscapeString(strings.TrimSpace(course.Title))
	course.Creator = Tutor{}
	course.CreatedAt= time.Now()
	course.UpdatedAt = time.Now()
}

//A function to validate the input using the new errors func that needs to be created

func (course *Course) SaveCourse(db *gorm.DB) (*Course, error) {
	var err error 
	err = db.Debug().Model(&Course{}).Create(&course).Error
	if err != nil {
		return &Course{}, nil 
	}

	if course.ID  != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ? ", course.ID).Take(&course.Creator).Error
		if err != nil {
			return &Course{}, err
		}
	} 

	return course, nil 
}

func (course *Course) FindAllCourses(db *gorm.DB) (*[]Course, error )  {
	var err error 
	courses := []Course{}
	err = db.Debug().Model(&Course{}).Limit(100).Find(&courses).Error
	if err != nil {
		return &[]Course{}, nil 
	}

	if len(courses) > 0 {
		for i := range courses {
			err := db.Debug().Model(&Course{}).Where("id = ?").Take(&courses[i].Creator).Error
			if err != nil {
				return &[]Course{}, err 
			}
		}
	}
	
	return &courses, nil 
}


func (course *Course) FindCourseById(db *gorm.DB, cid uint64) (*Course, error) {
	var err error

	err = db.Debug().Model(&Course{}).Where("id = ?", cid).Take(&course).Error
	if err != nil {
		return &Course{}, nil 
	}

	if course.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", course.CreatorID).Take(&course.Creator).Error
		if err != nil {
			return &Course{}, err 
		}
	}
	return course, nil 
}

func (course *Course) UpdateACourse(db *gorm.DB)(*Course, error) {
	var err error 

	err = db.Debug().Model(&Course{}).Where("id = ?").Updates(Course{Title: course.Title, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Course{}, nil 
	}

	if course.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", course.CreatorID).Take(&course.Creator).Error
		if err != nil {
			return &Course{}, err 
		}
	}
	return course, nil 
}

func (course *Course) DeleteACourse(db *gorm.DB, cid uint64, uid uint32) (uint64, error) {
	db = db.Debug().Model(&Course{}).Where("id = ? and tutor_id = ?", cid, uid).Take(&Course{}).Delete(&Course{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Course not found")
		}
		return 0, db.Error
	}
	return uint64(db.RowsAffected), nil 
}