package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql database driver
	"github.com/morelmiles/school-mgt-system/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Init(Dbdriver, Dbuser, Dbpassword, Dbport, Dbhost, Dbname string) {
	var err error

	DbURL := fmt.Sprintf("%s:%s%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", Dbdriver, Dbuser, Dbpassword, Dbhost, Dbport, Dbname)
	server.DB, err = gorm.Open(Dbdriver, DbURL)
	if err != nil {
		log.Fatal(err)
	}
	server.DB.Debug().AutoMigrate(&models.Student{}, &models.Course{}, &models.LearningTrack{}, &models.Module{}, &models.Tutor{}) //Migrate the database

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
