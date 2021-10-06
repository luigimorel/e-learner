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

func (server *Server) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if course id is valid 
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} 

	// Check if the course exists
	course := models.Course{}
	err = server.DB.Debug().Model(models.Course{}).Where("id = ?", cid).Take(&course).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("course not found"))
		return
	}	

	//Check if the token is valid
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("unauthorized"))
	}

	// If user tries to update a course not belonging to them 
	if uid != course.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	//Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return 
	}
	
	//Start processing the requested data 
	courseUpdate := models.Course{}
	err = json.Unmarshal(body, &courseUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	//Check if the req user id is equal to one gotten from the token 
	if uid != courseUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return 
	}

	courseUpdate.Prepare()
	err = courseUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return 
	}

	courseUpdate.ID = course.ID //tells the model the id to update 

	courseUpdated, err := courseUpdate.UpdateACourse(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
	}
	responses.JSON(w, http.StatusOK, courseUpdated)

}

func (server *Server) DeleteCourse(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)	

	//Check if the id is valid 
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return 
	}
	//Check if the user is authenticated 
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return 
	}
	//Check if the course exists 
	course := models.Course{}
	err = server.DB.Debug().Model(models.Course{}).Where("id = ?", cid).Take(&course).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("course not found"))
		return
	}
	// Check if the user is auth-e and owner of the course 
	if uid != course.CreatorID {
		responses.ERROR(w, http.StatusNotFound, errors.New("unauthorized"))
		return
	}
	_, err = course.DeleteACourse(server.DB, cid, uid)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Entity", fmt.Sprintf("%d", cid))
		responses.JSON(w, http.StatusNoContent, "")
} 