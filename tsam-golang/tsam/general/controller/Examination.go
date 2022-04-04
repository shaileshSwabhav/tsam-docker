package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ExaminationController provides method to update, delete, add, get method for examination.
type ExaminationController struct {
	log                log.Logger
	auth               *security.Authentication
	ExaminationService *service.ExaminationService
}

// NewExaminationController creates new instance of ExaminationController.
func NewExaminationController(examinationService *service.ExaminationService, log log.Logger, auth *security.Authentication) *ExaminationController {
	return &ExaminationController{
		ExaminationService: examinationService,
		log:                log,
		auth:               auth,
	}
}

// RegisterRoutes excludes routes from middleware.
func (controller *ExaminationController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all examinations.
	examinationList := router.HandleFunc("/tenant/{tenantID}/examination",
		controller.GetExaminations).Methods(http.MethodGet)

	// Get one examination by examination id.
	router.HandleFunc("/tenant/{tenantID}/examination/{examinationID}",
		controller.GetExamination).Methods(http.MethodGet)

	// Add one examination.
	router.HandleFunc("/tenant/{tenantID}/examination",
		controller.AddExamination).Methods(http.MethodPost)

	// Add multiple examinations.
	router.HandleFunc("/tenant/{tenantID}/examinations",
		controller.AddExaminations).Methods(http.MethodPost)

	// Update one examination.
	router.HandleFunc("/tenant/{tenantID}/examination/{examinationID}",
		controller.UpdateExamination).Methods(http.MethodPut)

	// Delete one examination.
	router.HandleFunc("/tenant/{tenantID}/examination/{examinationID}",
		controller.DeleteExamination).Methods(http.MethodDelete)

	// Exculde routes.
	*exclude = append(*exclude, examinationList)

	controller.log.Info("Examination Routes Registered")
}

// GetExaminations returns all examinations.
func (controller *ExaminationController) GetExaminations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetExaminations call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examinations := []general.Examination{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get examinations method.
	err = controller.ExaminationService.GetExaminations(&examinations, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResonseWriter.
	web.RespondJSON(w, http.StatusOK, examinations)
}

// GetExamination returns examination by id.
func (controller *ExaminationController) GetExamination(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************GetExamination call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examination := general.Examination{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	examination.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting examination id from param and parsing it to uuid.
	examination.ID, err = parser.GetUUID("examinationID")
	if err != nil {
		controller.log.Error("unable to parse examination id")
		web.RespondError(w, errors.NewHTTPError("unable to parse examination id", http.StatusBadRequest))
		return
	}

	// Call get examination method.
	err = controller.ExaminationService.GetExamination(&examination)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ReponseWriter.
	web.RespondJSON(w, http.StatusOK, examination)
}

// AddExamination adds new examination.
func (controller *ExaminationController) AddExamination(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddExamination call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examination := general.Examination{}

	//param := mux.Vars(r)

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &examination)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := examination.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	examination.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	examination.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add examination method.
	err = controller.ExaminationService.AddExamination(&examination)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResposeWriter.
	web.RespondJSON(w, http.StatusOK, "Examination added successfully")
}

// AddExaminations adds multiple examinations.
func (controller *ExaminationController) AddExaminations(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************AddExaminations call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examinationIDs := []uuid.UUID{}
	examinations := []general.Examination{}

	//param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &examinations)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate all compulsary fields of examination.
	for _, examination := range examinations {
		err = examination.Validate()
		if err != nil {
			controller.log.Error(err.Error())
			web.RespondError(w, err)
			return
		}

	}

	// Call add examinations method.
	err = controller.ExaminationService.AddExaminations(&examinations, &examinationIDs, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponeWriter.
	web.RespondJSON(w, http.StatusOK, "Examinations added successfully")
}

// UpdateExamination updates examination.
func (controller *ExaminationController) UpdateExamination(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************UpdateExamination call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examination := general.Examination{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	examination.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	examination.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting examination id from param and parsing it to uuid.
	examination.ID, err = parser.GetUUID("examinationID")
	if err != nil {
		controller.log.Error("unable to parse examination id")
		web.RespondError(w, errors.NewHTTPError("unable to parse examination id", http.StatusBadRequest))
		return
	}

	// Unmarshal JSON.
	err = web.UnmarshalJSON(r, &examination)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate examination.
	err = examination.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update examination method.
	err = controller.ExaminationService.UpdateExamination(&examination)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Examination updated successfully")
}

// DeleteExamination deletes examination.
func (controller *ExaminationController) DeleteExamination(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("********************************DeleteExamination call**************************************")
	parser := web.NewParser(r)
	// Create bucket.
	examination := general.Examination{}

	//param := mux.Vars(r)

	var err error

	// Getting tenant id from param and parsing it to uuid.
	examination.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	examination.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting examination id from param and parsing it to uuid.
	examination.ID, err = parser.GetUUID("examinationID")
	if err != nil {
		controller.log.Error("unable to parse examination id")
		web.RespondError(w, errors.NewHTTPError("unable to parse examination id", http.StatusBadRequest))
		return
	}

	// Call delete examination method.
	err = controller.ExaminationService.DeleteExamination(&examination)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Examination deleted successfully")
}
