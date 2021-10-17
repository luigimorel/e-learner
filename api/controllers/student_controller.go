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

func (server *Server) CreateStudent(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	student := models.Student{}
	err = json.Unmarshal(body, &student)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	student.Prepare()
	err = student.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	studentCreated, err := student.SaveStudent(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, studentCreated.ID))
	responses.JSON(w, http.StatusCreated, err)
}

func (server *Server) GetStudents(w http.ResponseWriter, r *http.Request) {

	student := models.Student{}

	students, err := student.FindAllStudents(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, students)
}

func (server *Server) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	student := models.Student{}

	studentResponse, err := student.FindStudentById(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, studentResponse)
}

func (server *Server) UpdateStudent(w http.ResponseWriter, r *http.Request) {
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
	student := models.Student{}
	err = json.Unmarshal(body, &student)
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

	student.Prepare()
	err = student.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	studentUpdated, err := student.UpdateStudent(server.DB, uint32(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
	}
	responses.JSON(w, http.StatusOK, studentUpdated)

}

func (server *Server) DeleteStudent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	student := models.Student{}

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

	_, err = student.DeleteStudent(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", cid))
	responses.JSON(w, http.StatusNoContent, "")
}
