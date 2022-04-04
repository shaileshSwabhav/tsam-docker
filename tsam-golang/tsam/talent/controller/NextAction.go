package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentNextActionController provides method to update, delete, add, get all, get one for talent next actions.
type TalentNextActionController struct {
	TalentNextActionService *service.TalentNextActionService
}

// NewTalentNextActionController creates new instance of TalentNextActionController.
func NewTalentNextActionController(talentNextActionService *service.TalentNextActionService) *TalentNextActionController {
	return &TalentNextActionController{
		TalentNextActionService: talentNextActionService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *TalentNextActionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/talent-next-action",
		controller.GetAllTalentNextActions).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/talent-next-action/credential/{credentialID}",
		controller.AddTalentNextAction).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/talent-next-action/{talentNextActionID}",
		controller.GetTalentNextAction).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/talent-next-action/{talentNextActionID}/credential/{credentialID}",
		controller.UpdateTalentNextAction).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/talent-next-action/{talentNextActionID}/credential/{credentialID}",
		controller.DeleteTalentNextAction).Methods(http.MethodDelete)

	log.NewLogger().Info("TalentNextAction Routes Registered")
}

// AddTalentNextAction adds one talent next action.
func (controller *TalentNextActionController) AddTalentNextAction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalentNextAction called=======================================")

	talentNextAction := talent.NextAction{}
	var err error

	params := mux.Vars(r)

	// Fill the talentNextAction variable with given data.
	err = web.UnmarshalJSON(r, &talentNextAction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talentNextAction.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of talentNextAction.
	talentNextAction.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentNextAction.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Validate compulsory fields.
	err = talentNextAction.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.TalentNextActionService.AddTalentNextAction(&talentNextAction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Next action added successfully")
}

// GetAllTalentNextActions gets all talent next actions.
func (controller *TalentNextActionController) GetAllTalentNextActions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentNextActions called=======================================")

	talentNextActions := []*talent.NextActionDTO{}

	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Call get talent next actions method.
	err = controller.TalentNextActionService.GetAllTalentNextActions(&talentNextActions, tenantID, talentID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentNextActions)
}

// GetTalentNextAction gets one talent next action.
func (controller *TalentNextActionController) GetTalentNextAction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentNextAction called=======================================")

	talentNextAction := talent.NextAction{}
	var err error

	params := mux.Vars(r)

	// Parse and set tenant ID.
	talentNextAction.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set talent ID.
	talentNextAction.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set talentNextAction ID.
	talentNextAction.ID, err = util.ParseUUID(params["talentNextActionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent next action id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TalentNextActionService.GetTalentNextAction(&talentNextAction); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talentNextAction)
}

// UpdateTalentNextAction updates talent next action.
func (controller *TalentNextActionController) UpdateTalentNextAction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateTalentNextAction called=======================================")

	talentNextAction := talent.NextAction{}

	params := mux.Vars(r)

	// Fill the talentNextAction variable with given data.
	err := web.UnmarshalJSON(r, &talentNextAction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set talentNextAction ID to talentNextAction.
	talentNextAction.ID, err = util.ParseUUID(params["talentNextActionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent next action id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to talentNextAction.
	talentNextAction.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of talentNextAction.
	talentNextAction.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talentNextAction ID to talentNextAction.
	talentNextAction.TalentID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = talentNextAction.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.TalentNextActionService.UpdateTalentNextAction(&talentNextAction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Next action updated successfully")
}

// DeleteTalentNextAction deletes one talent next action.
func (controller *TalentNextActionController) DeleteTalentNextAction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteTalentNextAction called=======================================")

	talentNextAction := talent.NextAction{}

	params := mux.Vars(r)

	var err error

	// Parse and set tenant ID.
	talentNextAction.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to talentNextAction's DeletedBy field.
	talentNextAction.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set talentNextAction ID to talentNextAction.
	talentNextAction.TalentID, err = util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set talentNextAction ID.
	talentNextAction.ID, err = util.ParseUUID(params["talentNextActionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent next action id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.TalentNextActionService.DeleteTalentNextAction(&talentNextAction); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Next action deleted successfully")
}
