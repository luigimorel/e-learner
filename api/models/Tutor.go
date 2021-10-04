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

type Tutor struct {
	ID       		 uint32        `gorm:"primary_key;auto_increment" json:"id"`
	FirstName  	string        `gorm:"size 255; not null" json:"first_name"`
	MiddleName string `gorm:"size 255" json:"middle_name"`
	LastName  string        `gorm:"size 255; not null" json:"last_name"`
	Email    		string        `gorm:"size 100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
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



func HashTutorPassword(password string) ([]byte, error)  {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyTutorPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (t *Tutor) BeforeSave() error  {
	hashedPassword, err := HashTutorPassword(t.Password)
	if err != nil {
		return err
	}

	t.Password = string(hashedPassword)
	return nil 
}

func (t *Tutor) Validate(action string) error  {
	switch strings.ToLower(action) {
	case "update":
		if t.Email == "" {
			return errors.New("email is required")
		}
		if t.Password == " " {
			return errors.New("password is required")
		}

		return nil
	case "login" : 
		if t.Password == " " {
			return errors.New("password is required")
		}
		if t.Email == "" {
			return errors.New("email is required")
		}
		if err := checkmail.ValidateFormat(t.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil 
	default: 
		if t.Password == " " {
			return errors.New("password is required")
		}
		if t.Email == "" {
			return errors.New("email is required")
		}
		if err := checkmail.ValidateFormat(t.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil 
	}
}

func (t *Tutor) SaveUser(db *gorm.DB) (*Tutor, error) {
	var err error = db.Debug().Create(&t).Error
	if err != nil {
		return &Tutor{}, err
	}
	return t, nil
}

func (t *Tutor) FindAllStudents(db *gorm.DB) (*[]Tutor, error){
	var err error 
	tutors := []Tutor{}
	err =  db.Debug().Model(&Tutor{}).Limit(100).Find(&tutors).Error
	if err != nil {
		return &[]Tutor{}, err
	}

	return &tutors, err
}

func (t*Tutor) FindStudentById(db *gorm.DB, tid uint32) (*Tutor, error) {
	
	var err error = db.Debug().Model(&Tutor{}).Where("id = ?").Take(&t).Error
	if err != nil {
		return &Tutor{}, err 
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Tutor{}, errors.New("tutor not found")
	}
	return t, err
}

func (t *Tutor) UpdateStudent(db *gorm.DB, tid uint32) (*Tutor, error)  {
	
	// Hash the password
	err := t.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&Tutor{}).Where("id = ?", tid).Take(&Tutor{}).UpdateColumns(
		map[string]interface{}{
			"password" : t.Password, 
			"email" : t.Email, 
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Tutor{}, err
	}

	//Display the update user 
	err = db.Debug().Model(&Tutor{}).Where("id = ?", tid).Take(&t).Error
	if err != nil {
		return &Tutor{}, err 
	}

	return t, nil 
}

func (t *Tutor) DeleteUser(db *gorm.DB, tid uint32) (int64, error) {
	db = db.Debug().Model(&Tutor{}).Where("id = ?", tid).Take(&Tutor{}).Delete(&Tutor{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil 
}