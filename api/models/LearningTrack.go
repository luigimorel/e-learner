package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type LearningTrack struct {
	ID        				uint32   			`gorm:"primary_key;auto_increment;" json:"id"`
	Title      				string   				`gorm:"size 255;not null;unique" json:"title"`
	Course 	  			Course   				`json:"course"`
	CourseID 			uint32			`gorm:"not null" json:"course_id"`
	Creator       		Tutor 			  `json:"creator"`
	CreatorID 			uint32			`gorm:"not null" json:"creator_id"`
	Learner				Student	`json:"student"` 
	LearnerID 			uint32 `gorm:"not null" json:"learner_id"`
	CreatedAt 		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 		time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (learningTrack *LearningTrack) Prepare()  {
	learningTrack.ID = 0
	learningTrack.Title = html.EscapeString(strings.TrimSpace(learningTrack.Title))
	learningTrack.Course = Course{}
	learningTrack.Creator = Tutor{}
	learningTrack.CreatedAt = time.Now()
	learningTrack.UpdatedAt = time.Now()
}

//A function to validate the input using the new errors func that needs to be created

func (lt *LearningTrack) SaveLearningTrack(db *gorm.DB) (*LearningTrack, error) {
	var err error 
	err = db.Debug().Model(&Course{}).Create(&lt).Error
	if err != nil {
		return &LearningTrack{}, nil 
	}

	if lt.ID  != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ? ", lt.ID).Take(&lt.Creator).Error
		if err != nil {
			return &LearningTrack{}, err
		}
	} 

	return lt, nil 
}


func (lt *LearningTrack) FindAllLearningTracks(db *gorm.DB) (*[]LearningTrack, error )  {
	var err error 
	learningTracks := []LearningTrack{}
	err = db.Debug().Model(&Course{}).Limit(100).Find(&learningTracks).Error
	if err != nil {
		return &[]LearningTrack{}, nil 
	}

	if len(learningTracks) > 0 {
		for i := range learningTracks {
			err := db.Debug().Model(&Course{}).Where("id = ?").Take(&learningTracks[i].Creator).Error
			if err != nil {
				return &[]LearningTrack{}, err 
			}
		}
	}
	
	return &learningTracks, nil 
}


func (lt *LearningTrack) FindLearningTrackById(db *gorm.DB, lt_id uint64) (*LearningTrack, error) {
	var err error

	err = db.Debug().Model(&Course{}).Where("id = ?", lt_id).Take(&lt).Error
	if err != nil {
		return &LearningTrack{}, nil 
	}

	if lt.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", lt.CreatorID).Take(&lt.Creator).Error
		if err != nil {
			return &LearningTrack{}, err 
		}
	}
	return lt, nil 
}

func (lt *LearningTrack) UpdateALearningTrack(db *gorm.DB)(*LearningTrack, error) {
	var err error 

	err = db.Debug().Model(&LearningTrack{}).Where("id = ?").Updates(LearningTrack{Title: lt.Title,Course: lt.Course, Learner: lt.Learner,  UpdatedAt: time.Now()}).Error
	if err != nil {
		return &LearningTrack{}, nil 
	}

	if lt.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", lt.CreatorID).Take(&lt.Creator).Error
		if err != nil {
			return &LearningTrack{}, err 
		}
	}

	return lt, nil 
}


func (lt *LearningTrack) DeleteALearningTrack(db *gorm.DB, cid uint64, uid uint32) (uint64, error) {
	db = db.Debug().Model(&LearningTrack{}).Where("id = ? and creator_id = ?", cid, uid).Take(&LearningTrack{}).Delete(&LearningTrack{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Course not found")
		}
		return 0, db.Error
	}
	return uint64(db.RowsAffected), nil 
}