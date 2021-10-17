package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Module struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Content      string    `gorm:"text;not null" json:"module_content"`
	Title        string    `gorm:"size 255;not null;" json:"title"`
	Creator      Tutor     `json:"author"`
	CreatorID    uint32    `gorm:"not null;" json:"author_id"`
	MainCourse   Course    `json:"main_course"`
	MainCourseID uint32    `gorm:"not null;" json:"main_course_id"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (module *Module) Prepare() {
	module.ID = 0
	module.Content = html.EscapeString(strings.TrimSpace(module.Content))
	module.Title = html.EscapeString(strings.TrimSpace(module.Title))
	module.Creator = Tutor{}
	module.MainCourse = Course{}
	module.CreatedAt = time.Now()
	module.UpdatedAt = time.Now()
}

func (m *Module) Validate() error {
	if m.Title == "" {
		return errors.New("required title")
	}
	if m.Content == "" {
		return errors.New("required content")
	}
	if m.CreatorID < 1 {
		return errors.New("required author")
	}
	return nil
}

func (m *Module) SaveModule(db *gorm.DB) (*Module, error) {
	var err error
	err = db.Debug().Model(&Module{}).Create(&m).Error

	if err != nil {
		return &Module{}, nil
	}

	if m.ID != 0 {
		err = db.Debug().Model(&Module{}).Where("id = ?", m.ID).Take(&m.Creator).Error
		if err != nil {
			return &Module{}, err
		}
	}
	return m, nil

}

func (m *Module) FindAllModules(db *gorm.DB) (*[]Module, error) {
	var err error
	modules := []Module{}
	err = db.Debug().Model(&Course{}).Limit(100).Find(&modules).Error
	if err != nil {
		return &[]Module{}, nil
	}

	if len(modules) > 0 {
		for i := range modules {
			err := db.Debug().Model(&Course{}).Where("id = ?").Take(&modules[i].Creator).Error
			if err != nil {
				return &[]Module{}, err
			}
		}
	}

	return &modules, nil
}

func (m *Module) FindModuleById(db *gorm.DB, cid uint64) (*Module, error) {
	var err error

	err = db.Debug().Model(&Course{}).Where("id = ?", cid).Take(&m).Error
	if err != nil {
		return &Module{}, nil
	}

	if m.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", m.CreatorID).Take(&m.Creator).Error
		if err != nil {
			return &Module{}, err
		}
	}
	return m, nil
}

func (m *Module) UpdateAModule(db *gorm.DB) (*Module, error) {
	var err error

	err = db.Debug().Model(&Module{}).Where("id = ?").Updates(Module{Title: m.Title, Content: m.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Module{}, nil
	}

	if m.ID != 0 {
		err = db.Debug().Model(&Course{}).Where("id = ?", m.CreatorID).Take(&m.Creator).Error
		if err != nil {
			return &Module{}, err
		}
	}
	return m, nil
}

func (m *Module) DeleteACourse(db *gorm.DB, cid uint64, uid uint32) (uint64, error) {
	db = db.Debug().Model(&Module{}).Where("id = ? and tutor_id = ?", cid, uid).Take(&Module{}).Delete(&Module{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Module not found")
		}
		return 0, db.Error
	}
	return uint64(db.RowsAffected), nil
}
