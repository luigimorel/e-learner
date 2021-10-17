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

func (server *Server) CreateModule(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	module := models.Module{}
	err = json.Unmarshal(body, &module)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	module.Prepare()
	err = module.Validate()
	if err != nil {

		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != module.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	moduleCreated, err := module.SaveModule(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, moduleCreated.ID))
	responses.JSON(w, http.StatusCreated, err)
}

func (server *Server) GetModules(w http.ResponseWriter, r *http.Request) {

	module := models.Module{}

	modules, err := module.FindAllModules(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, modules)
}

func (server *Server) GetModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	module := models.Module{}

	courseResponse, err := module.FindModuleById(server.DB, mid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, courseResponse)
}

func (server *Server) UpdateModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Check if module id is valid
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Check if the module exists
	module := models.Course{}
	err = server.DB.Debug().Model(models.Module{}).Where("id = ?", mid).Take(&module).Error
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
	if uid != module.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	//Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Start processing the requested data
	moduleUpdate := models.Module{}
	err = json.Unmarshal(body, &moduleUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	//Check if the req user id is equal to one gotten from the token
	if uid != moduleUpdate.CreatorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	moduleUpdate.Prepare()
	err = moduleUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	moduleUpdate.ID = module.ID //tells the model the id to update

	courseUpdated, err := moduleUpdate.UpdateAModule(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
	}
	responses.JSON(w, http.StatusOK, courseUpdated)

}

func (server *Server) DeleteModule(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//Check if the id is valid
	mid, err := strconv.ParseUint(vars["id"], 10, 64)
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
	module := models.Module{}
	err = server.DB.Debug().Model(models.Module{}).Where("id = ?", mid).Take(&module).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("module not found"))
		return
	}
	// Check if the user is auth-e and owner of the module
	if uid != module.CreatorID {
		responses.ERROR(w, http.StatusNotFound, errors.New("unauthorized"))
		return
	}
	_, err = module.DeleteACourse(server.DB, mid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", mid))
	responses.JSON(w, http.StatusNoContent, "")
}
