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

func (server *Server) CreateLearningTrack(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	learningTrack := models.LearningTrack{}
	err = json.Unmarshal(body, &learningTrack)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	learningTrack.Prepare()
	err = learningTrack.Validate()
	if err != nil {

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != learningTrack.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	learningTrackCreated, err := learningTrack.SaveLearningTrack(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, learningTrackCreated.ID))
	responses.JSON(w, http.StatusCreated, err)
}

func (server *Server) GetLearningTracks(w http.ResponseWriter, r *http.Request) {

	learningTrack := models.LearningTrack{}

	learningTracks, err := learningTrack.FindAllLearningTracks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, learningTracks)
}

func (server *Server) GetLearningTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lt_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	learningTrack := models.LearningTrack{}

	courseResponse, err := learningTrack.FindLearningTrackById(server.DB, lt_id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, courseResponse)
}

func (server *Server) UpdateLearningTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if learningTrack id is valid
	lt_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the learningTrack exists
	learningTrack := models.LearningTrack{}
	err = server.DB.Debug().Model(models.LearningTrack{}).Where("id = ?", lt_id).Take(&learningTrack).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("learning track not found"))
		return
	}

	//Check if the token is valid
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("unauthorized"))
	}

	// If user tries to update a course not belonging to them
	if uid != learningTrack.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	//Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Start processing the requested data
	learningTrackUpdate := models.LearningTrack{}
	err = json.Unmarshal(body, &learningTrackUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	//Check if the req user id is equal to one gotten from the token
	if uid != learningTrackUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	learningTrackUpdate.Prepare()
	err = learningTrackUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	learningTrackUpdate.ID = learningTrack.ID //tells the model the id to update

	learningTrackUpdated, err := learningTrackUpdate.UpdateALearningTrack(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
	}
	responses.JSON(w, http.StatusOK, learningTrackUpdated)
}

func (server *Server) DeleteLearningTrack(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//Check if the id is valid
	lt_id, err := strconv.ParseUint(vars["id"], 10, 64)
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

	//Check if the module exists
	learningTrack := models.LearningTrack{}
	err = server.DB.Debug().Model(models.LearningTrack{}).Where("id = ?", lt_id).Take(&learningTrack).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("learning track not found"))
		return
	}

	// Check if the user is auth-ed and owner of the module
	if uid != learningTrack.CreatorID {
		responses.ERROR(w, http.StatusNotFound, errors.New("unauthorized"))
		return
	}

	_, err = learningTrack.DeleteALearningTrack(server.DB, lt_id, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", lt_id))
	responses.JSON(w, http.StatusNoContent, "")
}
