package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchController Provide method to Update, Delete, Add, Get Method For Batch.
type BatchController struct {
	BatchService *service.BatchService
	log          log.Logger
	auth         *security.Authentication
}

// NewBatchController Create New Instance Of BatchController.
func NewBatchController(bs *service.BatchService, log log.Logger, auth *security.Authentication) *BatchController {
	return &BatchController{
		BatchService: bs,
		log:          log,
		auth:         auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *BatchController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch",
		controller.AddBatch).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}",
		controller.UpdateBatch).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}",
		controller.DeleteBatch).Methods(http.MethodDelete)

	// get
	batchList := router.HandleFunc("/tenant/{tenantID}/batch-list",
		controller.GetActiveBatchList).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch",
		controller.GetAllBatches).Methods(http.MethodGet)

	// Get upcoming batches.
	router.HandleFunc("/tenant/{tenantID}/batch/upcoming",
		controller.GetUpcomingBatches).Methods(http.MethodGet)

	// Get batch details for student login.
	router.HandleFunc("/tenant/{tenantID}/batch/details/{batchID}",
		controller.GetBatchDetails).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}",
		controller.GetBatch).Methods(http.MethodGet)

	// Add, Delete talent from or to Batch.
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/credential/{credentialID}",
	// 	controller.AddTalentsToBatch).Methods(http.MethodPost)
	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/{talentID}/credential/{credentialID}",
	// 	controller.DeleteTalentFromBatch).Methods(http.MethodDelete)

	// router.HandleFunc("/tenant/{tenantID}/batch/faculty/{facultyID}/limit/{limit}/offset/{offset}",
	// 	batchCtrl.GetAllBatchesForFaculty).Methods(http.MethodGet)
	// router.HandleFunc("/tenant/{tenantID}/batch/salesperson/{salesPersonID}/limit/{limit}/offset/{offset}",
	// 	batchCtrl.GetAllBatchesForSalesPerson).Methods(http.MethodGet)

	// search
	// router.HandleFunc("/tenant/{tenantID}/batch/search/limit/{limit}/offset/{offset}",
	// 	batchCtrl.SearchBatch).Methods(http.MethodPost)

	// Exculde routes.
	*exclude = append(*exclude, batchList)

	controller.log.Info("Batch route registered")
}

// AddBatch adds new batch by calling add service
func (controller *BatchController) AddBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddBatch call==============================")
	var err error
	batch := bat.Batch{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the batch variable with given data.
	err = web.UnmarshalJSON(r, &batch)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	batch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	// util.ParseUUID(params["credentialID"])
	batch.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate Batch
	err = batch.ValidateBatch()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Call add service to add batch to DB.
	err = controller.BatchService.AddBatch(&batch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batch.ID)
}

//UpdateBatch updates the batch by calling update service
func (controller *BatchController) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateBatch call==============================")
	var err error
	batch := bat.Batch{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// Fill the batch variable with given data.
	err = web.UnmarshalJSON(r, &batch)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set batch ID.
	// util.ParseUUID(params[paramBatchID])
	batch.ID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}
	// Parse and set tenant ID.
	// util.ParseUUID(params["tenantID"])
	batch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field.
	// util.ParseUUID(params["credentialID"])
	batch.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate Batch
	err = batch.ValidateBatch()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service to update batch to DB.
	err = controller.BatchService.UpdateBatch(&batch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Batch updated")

}

//DeleteBatch deletes batch by calling delete service
func (controller *BatchController) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteBatch call==============================")
	batch := bat.Batch{}
	var err error

	parser := web.NewParser(r)
	// util.ParseUUID(mux.Vars(r)[paramBatchID])
	batch.ID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	batch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	// util.ParseUUID(mux.Vars(r)["credentialID"])
	batch.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	err = controller.BatchService.DeleteBatch(&batch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Batch deleted")
}

// GetAllBatches returns all Batch
func (controller *BatchController) GetAllBatches(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllBatches call==============================")
	batches := &[]bat.BatchDTO{}

	parser := web.NewParser(r)
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	// total talents count
	var totalTalents int

	// Fills the form.
	// r.ParseForm()

	err = controller.BatchService.GetAllBatches(batches, tenantID, parser, &totalCount, &totalTalents)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// assign total talents assigned to batch in header
	web.SetNewHeader(w, "Total-Batch-Talents", strconv.Itoa(int(totalTalents)))
	// w.Header().Add("Access-Control-Expose-Headers", "Total-Batch-Talents")
	// w.Header().Set("Total-Batch-Talents", strconv.Itoa(int(totalTalents)))

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batches)
}

// GetUpcomingBatches returns upcoming batches.
func (controller *BatchController) GetUpcomingBatches(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetUpcomingBatches call==============================")

	// Craete bucket.
	batches := &[]bat.UpcomingBatch{}
	parser := web.NewParser(r)

	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get total, limit and offset from param and convert it to int.
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	r.ParseForm()

	// Call get all service method.
	err = controller.BatchService.GetUpcomingBatches(batches, tenantID,
		parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// web.RespondJSON(w, http.StatusOK, companyRequirements)
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batches)
}

// GetBatchDetails returns specific deatils of a single batch.
func (controller *BatchController) GetBatchDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatchDetails call==============================")

	// Craete bucket.
	batch := bat.BatchDetails{}

	parser := web.NewParser(r)
	// Getting tenant id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramTenantID])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting batch id from param and parsing it to uuid.
	// util.ParseUUID(mux.Vars(r)[paramBatchID])
	batch.ID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.BatchService.GetBatchDetails(&batch, tenantID, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batch)
}

// GetActiveBatchList returns listing of all the batches
func (controller *BatchController) GetActiveBatchList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetActiveBatchList call==============================")
	batches := &[]list.Batch{}

	parser := web.NewParser(r)
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// r.ParseForm()

	err = controller.BatchService.GetActiveBatchList(batches, parser, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, batches)
}

// GetBatch returns one batch
func (controller *BatchController) GetBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetBatch call==============================")
	var err error

	// param := mux.Vars(r)
	parser := web.NewParser(r)
	batch := bat.Batch{}

	// util.ParseUUID(param[paramBatchID])
	batch.ID, err = parser.GetUUID("batchID")
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	batch.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	err = controller.BatchService.GetBatch(&batch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// batches := []bat.Batch{
	// 	*batch,
	// }
	web.RespondJSON(w, http.StatusOK, batch)
}

// // GetAllBatchesForFaculty returns batches for specified faculty
// func (batchCtrl *BatchController) GetAllBatchesForFaculty(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetAllBatchesForFaculty call==============================")
// 	batches := &[]bat.BatchDTO{}

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	facultyID, err := util.ParseUUID(mux.Vars(r)["facultyID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
// 		return
// 	}

// 	// limit,offset & totalCount for pagination
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	// total talents count
// 	var totalTalents int

// 	// Fills the form.
// 	r.ParseForm()

// 	err = batchCtrl.BatchService.GetAllBatchesForFaculty(batches, r.Form, tenantID, facultyID,
// 		limit, offset, &totalCount, &totalTalents)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// assign total talents assigned to batch in header
// 	w.Header().Add("Access-Control-Expose-Headers", "Total-Batch-Talents")
// 	w.Header().Set("Total-Batch-Talents", strconv.Itoa(int(totalTalents)))

// 	// Writing Response with OK Status to ResponseWriter
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batches)

// }

// // GetAllBatchesForSalesPerson returns batches for specified salesperson
// func (batchCtrl *BatchController) GetAllBatchesForSalesPerson(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================GetAllBatchesForSalesPerson call==============================")
// 	batches := &[]bat.BatchDTO{}

// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	salesPersonID, err := util.ParseUUID(mux.Vars(r)["salesPersonID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
// 		return
// 	}

// 	// limit,offset & totalCount for pagination
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	// total talents count
// 	var totalTalents int

// 	// Fills the form.
// 	r.ParseForm()

// 	err = batchCtrl.BatchService.GetAllBatchesForSalesPerson(batches, r.Form, tenantID, salesPersonID,
// 		limit, offset, &totalCount, &totalTalents)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// assign total talents assigned to batch in header
// 	w.Header().Add("Access-Control-Expose-Headers", "Total-Batch-Talents")
// 	w.Header().Set("Total-Batch-Talents", strconv.Itoa(int(totalTalents)))

// 	// Writing Response with OK Status to ResponseWriter
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batches)

// }

// AddTalentsToBatch Add Student to Batch By Batch ID
// AddTalentsToBatch godoc
// AddTalentsToBatch adds talents to particular batch
// @Summary adds talent to batch
// @Description AddTalentsToBatch adds talent to specified batch if they are not already added
// @Tags Batch-Master
// @Accept  json
// @Produce  json
// @Param talent body bat.MappedTalent true "talentID, batchID and dateOfJoining of talent"
// @Param tenantID path string true "ID of the tenant where the batch and talent belong"
// @Param batchID path string true "ID of the batch where talent is to be added"
// @Param credentialID path string true "ID of the user who is adding the talent to a batch"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /batch/{batchID}/talent/credential/{credentialID} [post]
// func (controller *BatchController) AddTalentsToBatch(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================AddTalentsToBatch call==============================")
// 	var err error
// 	talents := &[]bat.MappedTalent{}
// 	params := mux.Vars(r)
// 	batchID, err := util.ParseUUID(params[paramBatchID])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
// 		return
// 	}
// 	tenantID, err := util.ParseUUID(params["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	credentialID, err := util.ParseUUID(params["credentialID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
// 		return
// 	}
// 	err = web.UnmarshalJSON(r, talents)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	for _, talent := range *talents {
// 		// Validate mapped talent fields
// 		err = talent.ValidateMappedTalent()
// 		if err != nil {
// 			controller.log.Error(err.Error())
// 			web.RespondError(w, err)
// 			return
// 		}
// 	}
// 	err = controller.BatchService.AddTalentsToBatch(talents, tenantID, batchID, credentialID)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Talents added")
// }

// DeleteTalentFromBatch Delete Student to Batch By Batch ID and Student ID
// func (controller *BatchController) DeleteTalentFromBatch(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================DeleteStudentFromBatch call==============================")
// 	params := mux.Vars(r)
// 	tenantID, err := util.ParseUUID(params["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	batchID, err := util.ParseUUID(params[paramBatchID])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
// 		return
// 	}
// 	talentID, err := util.ParseUUID(params["talentID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
// 		return
// 	}
// 	credentialID, err := util.ParseUUID(params["credentialID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
// 		return
// 	}
// 	err = controller.BatchService.DeleteTalentFromBatch(tenantID, credentialID, talentID, batchID)
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Talent deleted")
// }

// // SearchBatch Add New Batch
// func (batchCtrl *BatchController) SearchBatch(w http.ResponseWriter, r *http.Request) {
// 	controller.log.Info("==============================SearchBatch call==============================")
// 	batchSearch := &bat.Search{}
// 	batches := &[]*bat.Batch{}
// 	var totalCount int
// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	err = web.UnmarshalJSON(r, batchSearch)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	limit, offset := web.GetLimitAndOffset(r)
// 	err = batchCtrl.BatchService.SearchBatches(batches, batchSearch, tenantID, limit, offset, &totalCount)
// 	if err != nil {
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, batches)

// }
