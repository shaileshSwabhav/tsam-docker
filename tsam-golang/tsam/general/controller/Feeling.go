package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// FeelingController provides methods to do CRUD operations.
type FeelingController struct {
	log            log.Logger
	auth           *security.Authentication
	FeelingService *service.FeelingService
}

// NewFeelingController creates new instance of feeling type controller.
func NewFeelingController(feelingService *service.FeelingService, log log.Logger, auth *security.Authentication) *FeelingController {
	return &FeelingController{
		FeelingService: feelingService,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *FeelingController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all feelings by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/feeling",
		controller.GetAllFeelings).Methods(http.MethodGet)

	// Add feeling.
	router.HandleFunc("/tenant/{tenantID}/feeling",
		controller.AddFeeling).Methods(http.MethodPost)

	// Update feeling.
	router.HandleFunc("/tenant/{tenantID}/feeling/{feelingID}",
		controller.UpdateFeeling).Methods(http.MethodPut)

	// Delete feeling.
	router.HandleFunc("/tenant/{tenantID}/feeling/{feelingID}",
		controller.DeleteFeeling).Methods(http.MethodDelete)

	// Get feeling list.
	router.HandleFunc("/tenant/{tenantID}/feeling-list",
		controller.GetFeelingsList).Methods(http.MethodGet)

	// Get feeling levels by feeling.
	router.HandleFunc("/tenant/{tenantID}/feeling/{feelingID}/level",
		controller.GetFeelingLevels).Methods(http.MethodGet)

	controller.log.Info("Feeling Routes Registered")
}

// AddFeeling will add the feeling and its feeling levels.
func (controller *FeelingController) AddFeeling(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================AddFeeling called===========================")
	parser := web.NewParser(r)
	feeling := general.Feeling{}

	err := web.UnmarshalJSON(r, &feeling)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	feeling.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credentail id", http.StatusBadRequest))
		return
	}

	feeling.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = feeling.ValidateFeeling()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeelingService.AddFeeling(&feeling)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feeling added successfully")
}

// UpdateFeeling will update the existing feeling.
func (controller *FeelingController) UpdateFeeling(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================UpdateFeeling called===========================")
	parser := web.NewParser(r)
	feeling := general.Feeling{}

	err := web.UnmarshalJSON(r, &feeling)
	if err != nil {
		controller.log.Error("unable to parse requested data")
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	feeling.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credentail id", http.StatusBadRequest))
		return
	}

	feeling.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	feeling.ID, err = parser.GetUUID("feelingID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse feeling id", http.StatusBadRequest))
		return
	}

	err = feeling.ValidateFeeling()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FeelingService.UpdateFeeling(&feeling)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feeling updated successfully")
}

// DeleteFeeling will delete feeling and its feeling levels.
func (controller *FeelingController) DeleteFeeling(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================DeleteFeeling called===========================")
	parser := web.NewParser(r)
	var err error
	feeling := general.Feeling{}

	feeling.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credentail id", http.StatusBadRequest))
		return
	}

	feeling.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	feeling.ID, err = parser.GetUUID("feelingID")
	if err != nil {
		controller.log.Error("unable to parse feeling id")
		web.RespondError(w, errors.NewHTTPError("unable to parse feeling id", http.StatusBadRequest))
		return
	}

	err = controller.FeelingService.DeleteFeeling(&feeling)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Feeling deleted successfully")
}

// GetFeelingsList will return list of feelings from feelings table.
func (controller *FeelingController) GetFeelingsList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================GetFeelingsList called===========================")
	feelings := []general.Feeling{}
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.FeelingService.GetFeelingsList(&feelings, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, feelings)
}

// GetFeelingLevels will get the levels of feelings.
func (controller *FeelingController) GetFeelingLevels(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================GetFeelingLevels called===========================")
	feelingLevels := []general.FeelingLevel{}
	parser := web.NewParser(r)
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	feelingID, err := parser.GetUUID("feelingID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse feeling id", http.StatusBadRequest))
		return
	}

	err = controller.FeelingService.GetFeelingLevels(&feelingLevels, tenantID, feelingID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, feelingLevels)
}

// GetAllFeelings returns all feelings.
func (controller *FeelingController) GetAllFeelings(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===============================GetAllFeelings called=======================================")
	parser := web.NewParser(r)
	// Create bucket.
	feelings := []general.Feeling{}

	// Create bucket for total talents count.
	var totalCount int

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error("unable to parse tenant id")
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get career objectives method.
	if err := controller.FeelingService.GetAllFeelings(&feelings, tenantID, parser, &totalCount); err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, feelings)
}
