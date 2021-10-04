package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Student struct {
		ID 				    uint32 `gorm:"primary_key;auto_increment" json:"id"`
		FirstName	  string `gorm:"size 255; not null" json:"first_name"`
		MiddleName string `gorm:"size 255" json:"middle_name"`
		LastName	 string `gorm:"size 255; not null" json:"last_name"`
		Email 			  string `gorm:"size 100;not null;unique" json:"email"`
		Password  string    `gorm:"size:100;not null;" json:"password"`
		Progress			uint32 			`gorm:"not null" json:"progress"` 
		CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
		UpdatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (student *Student) Prepare() {
		student.ID = 0
		student.FirstName = html.EscapeString(strings.TrimSpace(student.FirstName))
		student.MiddleName = html.EscapeString(strings.TrimSpace(student.MiddleName))
		student.LastName =  html.EscapeString(strings.TrimSpace(student.LastName))
		student.Email = html.EscapeString(strings.TrimSpace(student.Email))
		student.CreatedAt = time.Now()
		student.UpdatedAt = time.Now()
}

func Hash(password string) ([]byte, error)  {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Student) BeforeSave() error  {
	hashedPassword, err := Hash(s.Password)
	if err != nil {
		return err
	}

	s.Password = string(hashedPassword)
	return nil 
}

func (s *Student) Validate(action string) error  {
	switch strings.ToLower(action) {
	case "update":
		if s.Email == "" {
			return errors.New("email is required")
		}
		if s.Password == " " {
			return errors.New("password is required")
		}

		return nil
	case "login" : 
		if s.Password == " " {
			return errors.New("password is required")
		}
		if s.Email == "" {
			return errors.New("email is required")
		}
		if err := checkmail.ValidateFormat(s.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil 
	default: 
		if s.Password == " " {
			return errors.New("password is required")
		}
		if s.Email == "" {
			return errors.New("email is required")
		}
		if err := checkmail.ValidateFormat(s.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil 
	}
}

func (s *Student) SaveUser(db *gorm.DB) (*Student, error) {
	var err error = db.Debug().Create(&s).Error
	if err != nil {
		return &Student{}, err
	}
	return s, nil
}

func (s *Student) FindAllStudents(db *gorm.DB) (*[]Student, error){
	var err error 
	students := []Student{}
	err =  db.Debug().Model(&Student{}).Limit(100).Find(&students).Error
	if err != nil {
		return &[]Student{}, err
	}

	return &students, err
}

func (s *Student) FindStudentById(db *gorm.DB, sid uint32) (*Student, error) {
	
	var err error = db.Debug().Model(&Student{}).Where("id = ?").Take(&s).Error
	if err != nil {
		return &Student{}, err 
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Student{}, errors.New("Student not found")
	}
	return s, err
}

func (s *Student) UpdateStudent(db *gorm.DB, sid uint32) (*Student, error)  {
	
	// Hash the password
	err := s.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&Student{}).Where("id = ?", sid).Take(&Student{}).UpdateColumns(
		map[string]interface{}{
			"password" : s.Password, 
			"email" : s.Email, 
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Student{}, err
	}

	//Display the update user 
	err = db.Debug().Model(&Student{}).Where("id = ?", sid).Take(&s).Error
	if err != nil {
		return &Student{}, err 
	}

	return s, nil 
}

func (s *Student) DeleteUser(db *gorm.DB, sid uint32) (int64, error) {
	db = db.Debug().Model(&Student{}).Where("id = ?", sid).Take(&Student{}).Delete(&Student{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil 
}