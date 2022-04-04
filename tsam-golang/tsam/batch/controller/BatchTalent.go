package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchTalentController provides methods to do Update, Add, Get operations on batch talents.
type BatchTalentController struct {
	log     log.Logger
	service *service.BatchTalentService
	auth    *security.Authentication
}

// NewBatchTalentController creates new instance of BatchTalentController.
func NewBatchTalentController(service *service.BatchTalentService,
	log log.Logger, auth *security.Authentication) *BatchTalentController {
	return &BatchTalentController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *BatchTalentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add batch talent.
	router.HandleFunc("/tenant/{tenantID}/batch-talent/batch/{batchID}",
		controller.AddTalentsToBatch).Methods(http.MethodPost)

	// Update batch talent.
	router.HandleFunc("/tenant/{tenantID}/batch-talent/{batchTalentID}",
		controller.UpdateBatchTalent).Methods(http.MethodPut)

	// Update suspension date of batch talent.
	router.HandleFunc("/tenant/{tenantID}/batch-talent-suspension-date/{batchTalentID}",
		controller.UpdateSuspensionDateBatchTalent).Methods(http.MethodPut)

	// Update is active of batch talent.
	router.HandleFunc("/tenant/{tenantID}/batch-talent-is-active/{batchTalentID}",
		controller.UpdateIsActiveBatchTalent).Methods(http.MethodPut)

	// Get details of all talents in a batch.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-details",
		controller.GetBatchMultipleTalentDetails).Methods(http.MethodGet)

	// Get details of one talent in a batch.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}",
		controller.GetBatchTalentDetails).Methods(http.MethodGet)

	// Get batches of one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/batches-minimum",
		controller.GetBatchesOfTalent).Methods(http.MethodGet)

	// Get talent list of a batch
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent-list",
		controller.GetTalentListOfBatch).Methods(http.MethodGet)

	controller.log.Info("Batch Talent Routes Registered")
}

// AddTalentsToBatch adds talent to batch.
func (controller *BatchTalentController) AddTalentsToBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddTalentsToBatch call==============================")

	talents := &[]batch.MappedTalent{}
	params := mux.Vars(r)

	batchID, err := util.ParseUUID(params[paramBatchID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	err = web.UnmarshalJSON(r, talents)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	for _, talent := range *talents {
		// Validate mapped talent fields
		err = talent.ValidateMappedTalent()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}
	err = controller.service.AddTalentsToBatch(talents, tenantID, batchID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Talents added")
}

// UpdateBatchTalent updates batch talent.
func (controller *BatchTalentController) UpdateBatchTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateBatchTalent Called==============================")
	batchTalent := batch.MappedTalent{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &batchTalent)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	batchTalent.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch talent ID.
	batchTalent.ID, err = parser.GetUUID("batchTalentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch talent id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	batchTalent.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = batchTalent.ValidateMappedTalent()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateBatchTalent(&batchTalent)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Batch talent successfully updated")
}

// UpdateSuspensionDateBatchTalent updates suspension date of batch talent.
func (controller *BatchTalentController) UpdateSuspensionDateBatchTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateSuspensionDateBatchTalent Called==============================")
	batchTalent := batch.UpdateBatchTalentSuspension{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &batchTalent)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch talent ID.
	batchTalent.ID, err = parser.GetUUID("batchTalentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch talent id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateSuspensionDateBatchTalent(&batchTalent, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch talent successfully updated")
}

// UpdateIsActiveBatchTalent updates is active of batch talent.
func (controller *BatchTalentController) UpdateIsActiveBatchTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateIsActiveBatchTalent Called==============================")
	batchTalent := batch.UpdateBatchTalentIsActive{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &batchTalent)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch talent ID.
	batchTalent.ID, err = parser.GetUUID("batchTalentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch talent id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateIsActiveBatchTalent(&batchTalent, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Batch talent successfully updated")
}

// GetBatchMultipleTalentDetails will return details of all talents for specified batch.
func (controller *BatchTalentController) GetBatchMultipleTalentDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchMultipleTalentDetails Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	sessionTalents := &[]batch.BatchTalentDTO{}
	err = controller.service.GetBatchMultipleTalentDetails(tenantID, batchID, sessionTalents, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, sessionTalents)
}

// GetBatchTalentDetails will return details of one talent for specified batch.
func (controller *BatchTalentController) GetBatchTalentDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchTalentDetails Called==============================")
	parser := web.NewParser(r)

	batchTalent := batch.BatchTalentDTO{}

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchTalent.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	batchTalent.TalentID, err = parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	err = controller.service.GetBatchTalentDetails(tenantID, &batchTalent, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, batchTalent)
}

// GetBatchesOfTalent will return batches of one talent.
func (controller *BatchTalentController) GetBatchesOfTalent(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchesOfTalent Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	batchTalent := []batch.MinimumBatchTalentForTalent{}
	err = controller.service.GetBatchesOfTalent(tenantID, talentID, &batchTalent, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, batchTalent)
}

// GetTalentListOfBatch will all talent names in the specific batch.
func (controller *BatchTalentController) GetTalentListOfBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetTalentListOfBatch Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	batchTalents := []list.Talent{}
	err = controller.service.GetTalentListOfBatch(tenantID, batchID, &batchTalents, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, batchTalents)
}
