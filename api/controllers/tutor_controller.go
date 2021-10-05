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

func (server *Server) CreateTutor(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tutor := models.Tutor{}
	err = json.Unmarshal(body, &tutor)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tutor.Prepare()
	err = tutor.Validate("")
	if err != nil {
	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	return
	}

	tutorCreated, err := tutor.SaveTutor(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return 
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, tutorCreated.ID))
	responses.JSON(w, http.StatusCreated, err)
}

func (server *Server) GetTutors(w http.ResponseWriter, r *http.Request) {

	tutor :=models.Tutor{}
	
	tutors, err := tutor.FindAllTutors(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err )
		return 
	}
	responses.JSON(w, http.StatusOK, tutors)
} 

func (server *Server) GetTutor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"],10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} 	
	tutor := models.Tutor{}

	studentResponse, err := tutor.FindTutorById(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err )
		return
	}
	responses.JSON(w, http.StatusOK, studentResponse)
} 

func (server *Server) UpdateTutor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if student id is valid 
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	} 

		//Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return 
	}

		//Start processing the requested data 
	tutor := models.Tutor{}
	err = json.Unmarshal(body, &tutor)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	
	//Check if the token is valid
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("unauthorized"))
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return 
	}

	tutor.Prepare()
	err = tutor.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return 
	}

	tutorUpdated, err := tutor.UpdateTutor(server.DB, uint32(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
	}
	responses.JSON(w, http.StatusOK, tutorUpdated)

}

func (server *Server) DeleteTutor(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)	
	
	tutor := models.Tutor{}

	//Check if the id is valid 
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return 
	}

	//Check if the student is authenticated 
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return 
	}
	if tokenID != 0 && tokenID != uint32(cid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = tutor.DeleteTutor(server.DB, uint32(cid))
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
		w.Header().Set("Entity", fmt.Sprintf("%d", cid))
		responses.JSON(w, http.StatusNoContent, "")
} 