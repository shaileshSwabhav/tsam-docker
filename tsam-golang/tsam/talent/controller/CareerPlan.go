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

// CareerPlanController provides methods to update, delete, add, get methods for career plan.
type CareerPlanController struct {
	CareerPlanService *service.CareerPlanService
}

// NewCareerPlanController Create New Instance Of CareerPlanController.
func NewCareerPlanController(careerPlanServic *service.CareerPlanService) *CareerPlanController {
	return &CareerPlanController{
		CareerPlanService: careerPlanServic,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *CareerPlanController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all career plans.
	router.HandleFunc("/tenant/{tenantID}/career-plan/talent/{talentID}", controller.GetCareerPlans).Methods(http.MethodGet)

	// Add career plans.
	router.HandleFunc("/tenant/{tenantID}/career-plan/talent/{talentID}/credential/{credentialID}", controller.AddCareerPlans).Methods(http.MethodPost)

	// Update one career plan.
	router.HandleFunc("/tenant/{tenantID}/career-plan/{careerPlanID}/talent/{talentID}/credential/{credentialID}", controller.UpdateCareerPlan).Methods(http.MethodPut)

	// Delete one career plan.
	router.HandleFunc("/tenant/{tenantID}/career-plan/{careerObjectiveID}/talent/{talentID}/credential/{credentialID}", controller.DeleteCareerPlan).Methods(http.MethodDelete)

	log.NewLogger().Info("Career Objective Route Registered")
}

// AddCareerPlans adds one career objective.
func (controller *CareerPlanController) AddCareerPlans(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddCareerPlans called=======================================")

	// Create bucket.
	careerPlans := []talent.CareerPlan{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &careerPlans)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	for _, careerPlan := range careerPlans {
		if err := careerPlan.Validate(); err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	// Parse tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse tenant ID.
	talentID, err := util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse credential id.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.CareerPlanService.AddCareerPlans(&careerPlans, tenantID, talentID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career plan added successfully")
}

// GetCareerPlans returns all career plans.
func (controller *CareerPlanController) GetCareerPlans(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCareerPlans called=======================================")

	// Create bucket
	careerPlans := []talent.CareerPlan{}

	// Getting talent id from param and parsing it to uuid.
	talentID, err := util.ParseUUID(mux.Vars(r)["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Getting talent id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get career objectives method.
	if err := controller.CareerPlanService.GetCareerPlans(&careerPlans, tenantID, talentID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, careerPlans)
}

// UpdateCareerPlan updates one career plan by specific career plan id.
func (controller *CareerPlanController) UpdateCareerPlan(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateCareerPlan called=======================================")

	// Create bucket.
	careerPlan := talent.CareerPlan{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &careerPlan)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = careerPlan.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set talent ID to careerPlan.
	careerPlan.ID, err = util.ParseUUID(params["careerPlanID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Career Objective id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to careerPlan.
	careerPlan.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in DeletedBy field of careerPlan.
	careerPlan.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	if err = controller.CareerPlanService.UpdateCareerPlan(&careerPlan); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career plan updated successfully")
}

// DeleteCareerPlan deletes one career plan by specific career plan id.
func (controller *CareerPlanController) DeleteCareerPlan(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteCareerPlan called=======================================")

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set talent ID to careerPlan.
	careerObjectiveID, err := util.ParseUUID(params["careerObjectiveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Career Objective id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to careerPlan.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of careerPlan.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	if err := controller.CareerPlanService.DeleteCareerPlan(careerObjectiveID, tenantID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Career plan deleted successfully")
}
