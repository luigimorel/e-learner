package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/morelmiles/school-mgt-system/api/auth"
	"github.com/morelmiles/school-mgt-system/api/models"
	"github.com/morelmiles/school-mgt-system/api/responses"
	"github.com/morelmiles/school-mgt-system/api/utils"
)

func (server *Server) CreateCourse(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	course := models.Course{}
	err = json.Unmarshal(body, &course)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	course.Prepare()
	err = course.Validate()
	if err != nil {

	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	return
	}

	uid,err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return 
	}
	if uid != course.CreatorID  {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return 
	}

	courseCreated, err := course.SaveCourse(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return 
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, courseCreated.ID))
	responses.JSON(w, http.StatusCreated, err)
}

func (server *Server) GetCourses(w http.ResponseWriter, r *http.Request) {

	course :=models.Course{}
	
	courses, err := course.FindAllCourses(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err )
		return 
	}
	responses.JSON(w, http.StatusOK, courses)
} 

func (server *Server) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"],10,64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} 	
	course := models.Course{}

	courseResponse, err := course.FindCourseById(server.DB, cid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err )
		return
	}
	responses.JSON(w, http.StatusOK, courseResponse)
} 