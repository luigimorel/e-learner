package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/morelmiles/school-mgt-system/api/models"
)

var students = []models.Student{
	{
		FirstName: "Luigi ",
		LastName:  "Morel",
		Email:     "morel@gmail.com",
		Password:  "password",
	},
	{
		FirstName: "Kasita",
		LastName:  "John",
		Email:     "kasita@gmail.com",
		Password:  "password123",
	},
}

var modules = []models.Module{
	{
		Content: "Lorem ipsum",
		Title:   "Lorem ipsum title",
	},
}

var creators = []models.Tutor{
	{
		FirstName: "Luigi",
		LastName:  "Morlel",
		Email:     "kl@test.com",
		Password:  "password",
	},
	{
		FirstName: "Kand",
		LastName:  "iwe",
		Email:     "kl@test.com",
		Password:  "password",
	},
}

var learningTrack = []models.LearningTrack{
	{
		Title: "Learning track one ",
	},
}

var courses = []models.Course{
	{
		Name:  "Machines learning",
		Title: "AI",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Student{}, &models.Course{}, &models.LearningTrack{}, &models.Module{}, &models.Tutor{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Student{}, &models.Course{}, &models.LearningTrack{}, &models.Module{}, &models.Tutor{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range students {
		err = db.Debug().Model(&models.Student{}).Create(&students[i]).Error
		if err != nil {
			log.Fatalf("cannot seed students table: %v", err)
		}
		modules[i].CreatorID = creators[i].ID

		err = db.Debug().Model(&models.Course{}).Create(&creators[i]).Error
		if err != nil {
			log.Fatalf("cannot seed creaors table: %v", err)
		}

		learningTrack[i].CourseID = courses[i].ID
		err = db.Debug().Model(&models.LearningTrack{}).Create(&learningTrack[i]).Error
		if err != nil {
			log.Fatalf("cannot seed learning track table: %v", err)
		}
	}
}
