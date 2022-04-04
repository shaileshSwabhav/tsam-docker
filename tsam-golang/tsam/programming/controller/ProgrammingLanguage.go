package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingLanguageController provide method to update, delete, add, get for programming language.
type ProgrammingLanguageController struct {
	ProgrammingLanguageService *service.ProgrammingLanguageService
}

// NewProgrammingLanguageController creates new instance of ProgrammingLanguageController.
func NewProgrammingLanguageController(programmingLanguageService *service.ProgrammingLanguageService) *ProgrammingLanguageController {
	return &ProgrammingLanguageController{
		ProgrammingLanguageService: programmingLanguageService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *ProgrammingLanguageController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all programming languages by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/programming-language/limit/{limit}/offset/{offset}",
		controller.GetProgrammingLanguages).Methods(http.MethodGet)

	// Get programming language list.
	programmingLanguageList := router.HandleFunc("/tenant/{tenantID}/programming-language",
		controller.GetProgrammingLanguageList).Methods(http.MethodGet)

	// Get one programming language.
	router.HandleFunc("/tenant/{tenantID}/programming-language/{languageID}",
		controller.GetProgrammingLanguage).Methods(http.MethodGet)

	// Add one programming language.
	router.HandleFunc("/tenant/{tenantID}/programming-language/credential/{credentialID}",
		controller.AddProgrammingLanguage).Methods(http.MethodPost)

	// Update programming language.
	router.HandleFunc("/tenant/{tenantID}/programming-language/{languageID}/credential/{credentialID}",
		controller.UpdateProgrammingLanguage).Methods(http.MethodPut)

	// Delete programming language.
	router.HandleFunc("/tenant/{tenantID}/programming-language/{languageID}/credential/{credentialID}",
		controller.DeleteProgrammingLanguage).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, programmingLanguageList)

	log.NewLogger().Info("Programming language Route Registered")
}

// GetProgrammingLanguages returns all the programming languages.
func (controller *ProgrammingLanguageController) GetProgrammingLanguages(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingLanguages called==============================")

	// Create bucket.
	languages := []general.ProgrammingLanguage{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all service method.
	err = controller.ProgrammingLanguageService.GetProgrammingLanguages(&languages, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, languages)
}

// GetProgrammingLanguageList returns all the programming languages (without pagination).
func (controller *ProgrammingLanguageController) GetProgrammingLanguageList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingLanguageList call**************************************")

	// Create bucket.
	languages := []general.ProgrammingLanguage{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.ProgrammingLanguageService.GetProgrammingLanguageList(&languages, tenantID)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, languages)
}

// GetProgrammingLanguage returns the specified programming language by id.
func (controller *ProgrammingLanguageController) GetProgrammingLanguage(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetProgrammingLanguage call**************************************")

	// Create bucket.
	language := general.ProgrammingLanguage{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	language.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting language id from param and parsing it to uuid.
	language.ID, err = util.ParseUUID(mux.Vars(r)["languageID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse programming language id", http.StatusBadRequest))
		return
	}

	// Call get language by id service method.
	err = controller.ProgrammingLanguageService.GetProgrammingLanguage(&language)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, language)
}

// AddProgrammingLanguage adds new programming language.
func (controller *ProgrammingLanguageController) AddProgrammingLanguage(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddProgrammingLanguage call**************************************")

	// Create bucket.
	language := general.ProgrammingLanguage{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = language.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	language.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	language.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.ProgrammingLanguageService.AddProgrammingLanguage(&language)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming language added successfully")
}

// UpdateProgrammingLanguage updates the specified programming language by id.
func (controller *ProgrammingLanguageController) UpdateProgrammingLanguage(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateProgrammingLanguage call**************************************")

	// Create bucket.
	language := general.ProgrammingLanguage{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	language.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting langauge id from param and parsing it to uuid.
	language.ID, err = util.ParseUUID(mux.Vars(r)["languageID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse programming language id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	language.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate programming language fields.
	err = language.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.ProgrammingLanguageService.UpdateProgrammingLanguage(&language)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming language updated successfully")
}

// DeleteProgrammingLanguage deletes specific programming language by id.
func (controller *ProgrammingLanguageController) DeleteProgrammingLanguage(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteProgrammingLanguage call**************************************")

	// Create bucket.
	language := general.ProgrammingLanguage{}

	var err error

	param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	language.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	language.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting programming language id from param and parsing it to uuid.
	language.ID, err = util.ParseUUID(param["languageID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse programming language id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.ProgrammingLanguageService.DeleteProgrammingLanguage(&language)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Programming language deleted successfully")
}
