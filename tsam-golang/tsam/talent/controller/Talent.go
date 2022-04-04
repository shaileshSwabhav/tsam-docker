package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/techlabs/swabhav/tsam/repository"
	service "github.com/techlabs/swabhav/tsam/talent/service"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentController provides methods to update, delete, add, get all, get all by salesperson,
// get one by id and get all by company requirement for talent.
type TalentController struct {
	TalentService *service.TalentService
}

// NewTalentController creates new instance of TalentController.
func NewTalentController(talentService *service.TalentService) *TalentController {
	return &TalentController{
		TalentService: talentService,
	}
}

// RegisterRoutes registers all endpoint to router.
func (controller *TalentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	//***************************************TALENT CRUD API*******************************************
	// Get all talents with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent/limit/{limit}/offset/{offset}",
		controller.GetTalents).Methods(http.MethodGet)

	// Add one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/credential/{credentialID}",
		controller.AddTalent).Methods(http.MethodPost)

	// Add talents.
	router.HandleFunc("/tenant/{tenantID}/talents/credential/{credentialID}",
		controller.AddTalents).Methods(http.MethodPost)

	// Add talent from excel .
	router.HandleFunc("/tenant/{tenantID}/talent-excel/credential/{credentialID}",
		controller.AddTalentFromExcel).Methods(http.MethodPost)

	// Get one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}", controller.GetTalent).Methods(http.MethodGet)

	// Get one talent for checking eligibility.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/eligible", controller.GetEligibleTalent).Methods(http.MethodGet)

	// Update one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/credential/{credentialID}",
		controller.UpdateTalent).Methods(http.MethodPut)

	// Delete one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/credential/{credentialID}",
		controller.DeleteTalent).Methods(http.MethodDelete)

	//***************************************TALENT OTHER API*******************************************
	// Search all talents with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent/search/limit/{limit}/offset/{offset}",
		controller.GetAllSearchTalents).Methods(http.MethodPost)

	// Get all talents by company requirement with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent/requirement/{requirementID}/limit/{limit}/offset/{offset}",
		controller.GetTalentsByCompanyRequirement).Methods(http.MethodGet)

	// Get all talents by batch id with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/talent/limit/{limit}/offset/{offset}",
		controller.GetTalentsForBatch).Methods(http.MethodGet)

	// Update all talents' salesperson.
	router.HandleFunc("/tenant/{tenantID}/talent/saleperson/{salesPersonID}/credential/{credentialID}",
		controller.UpdateTalentsSalesperson).Methods(http.MethodPut)

	// Get all batches of one talent.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/batch",
		controller.GetTalentBatches).Methods(http.MethodGet)

	// Get talents by waiting list's ids.
	router.HandleFunc("/tenant/{tenantID}/talent/waiting-list/limit/{limit}/offset/{offset}",
		controller.GetTalentsByWaitingList).Methods(http.MethodGet)

	// Get talents by professional summary report.
	router.HandleFunc("/tenant/{tenantID}/talent/pro-summary-report/limit/{limit}/offset/{offset}",
		controller.GetTalentsForProSummaryReport).Methods(http.MethodGet)

	// Get talents by professional summary report's technology talent count.
	router.HandleFunc("/tenant/{tenantID}/talent/pro-summary-report-tech-count/limit/{limit}/offset/{offset}",
		controller.GetTalentsForProSummaryReportTechnologyCount).Methods(http.MethodGet)

	// Get talents by fresher summary report.
	router.HandleFunc("/tenant/{tenantID}/talent/fresher-summary-report/limit/{limit}/offset/{offset}",
		controller.GetTalentsForFresherSummaryReport).Methods(http.MethodGet)

	// Get talents by fresher summary report.
	router.HandleFunc("/tenant/{tenantID}/talent/package-summary-report/limit/{limit}/offset/{offset}",
		controller.GetTalentsForPackageSummaryReport).Methods(http.MethodGet)

	// Get talents by campus drive id.
	router.HandleFunc("/tenant/{tenantID}/talent/campus-drive/{campusDriveID}/limit/{limit}/offset/{offset}",
		controller.GetTalentsByCampusDrive).Methods(http.MethodGet)

	// Get talents by seminar id.
	router.HandleFunc("/tenant/{tenantID}/talent/seminar/{seminarID}/limit/{limit}/offset/{offset}",
		controller.GetTalentsBySeminar).Methods(http.MethodGet)

	// Get talents for excel download with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent/excel-download/limit/{limit}/offset/{offset}",
		controller.GetTalentsForExcelDownload).Methods(http.MethodGet)

	// Get searched talents for excel download with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent-search/excel-download/limit/{limit}/offset/{offset}",
		controller.GetSearchedTalentsForExcelDownload).Methods(http.MethodPost)

	// Add credential for updating password.
	router.HandleFunc("/add-credential-for-talent/credential/{credentialID}",
		controller.AddCredentialForTalent).Methods(http.MethodPost)

	// Update the experience in months field for all talents.
	router.HandleFunc("/update-experience-in-months-of-talents",
		controller.UpdateExperienceInMonthsOfTalent).Methods(http.MethodPut)

	log.NewLogger().Info("Talent Routes Registered")
}

//***************************************TALENT CRUD API*******************************************

// GetTalents godoc
// GetTalents returns all talents
// @Summary Get All Talents
// @Description GetTalents Return all Talent By Limit and Offset
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param limit path string true "total talent return in response"
// @Param offset path string true "page number"
// @Success 200 {array} []tal.Talent
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent/{limit}/{offset} [Get]
func (controller *TalentController) GetTalents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalents called=======================================")

	// Create bucket.
	talents := []tal.DTO{}

	// Create bucket for total talents count.
	var totalCount int

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Fill the r.Form.
	r.ParseForm()

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get talents method.
	err = controller.TalentService.GetTalents(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// AddTalent adds one talent
// @Summary Add New Talent
// @Description Add New Talent
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param talent body tal.Talent true "Talent Data"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent [post]
func (controller *TalentController) AddTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalent called=======================================")

	// Create bucket.
	talent := tal.Talent{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &talent)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talent.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of talent.
	talent.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if _, err = controller.TalentService.AddTalent(&talent, false); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talent added successfully")
}

// AddTalentFromExcel adds one talent from excel.
// @Summary Add New Talent
// @Description Add New Talent
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param talent body tal.Talent true "Talent Data"
// @Success 200 {plain} plain
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent [post]
func (controller *TalentController) AddTalentFromExcel(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalentFromExcel called=======================================")

	// Create bucket.
	talentExcel := tal.TalentExcel{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &talentExcel)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := talentExcel.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of talent.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add talent from excel service method.
	if err = controller.TalentService.AddTalentFromExcel(&talentExcel, tenantID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talent from excel added successfully")
}

// AddTalents Add New Talent
// @Summary Add New Talents
// @Description Add Multiple New Talent
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param talents body []tal.Talent true "Multiple Talent Data"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Success 200 {array} []string
// @Router /talent/multiple [Post]
func (controller *TalentController) AddTalents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddTalents called=======================================")

	// Create bucket for talent successfully added count.
	talentAddedCount := 0

	// Create bucket for talents to be added.
	talents := []tal.TalentExcel{}

	// Create bucket for errors.
	var errorList []error

	// Get params from api.
	params := mux.Vars(r)

	// Parse tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID.
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Fill the talents variable with given data.
	err = web.UnmarshalJSON(r, &talents)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	for _, talent := range talents {
		err := talent.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	// Call add talents service method.
	err = controller.TalentService.AddTalents(&talents, &talentAddedCount, tenantID, &errorList, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	for _, err := range errorList {
		fmt.Println("++----------------ERRORS ARE ------------------------", err)
	}

	// Will define this struct if it really useful in future #Niranjan.
	payload := struct {
		ErrorList []error `json:"errorList"`
		Message   string  `json:"message"`
	}{
		ErrorList: errorList,
		Message:   fmt.Sprintf("Talents added: %v", talentAddedCount),
	}

	// If error while adding individual talent send error list.
	if len(errorList) > 0 {
		web.RespondJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talents added successfully")
}

// GetTalent returns talent by specific talent id
// @Summary Return Talent By Talent ID
// @Description GetTalent Return Talent By Talent ID
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param talentID path string true "Talent ID" Format(uuid.UUID)
// @Success 200 {array} []tal.Talent
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent/{talentID} [Get]
func (controller *TalentController) GetTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalent called=======================================")

	// Create bucket.
	talent := tal.DTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set talent ID.
	talent.ID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TalentService.GetTalent(&talent, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talent)
}

// GetEligibleTalent returns talent by specific talent id for getting fields for eligibility check.
func (controller *TalentController) GetEligibleTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetEligibleTalent called=======================================")

	// Create bucket.
	talent := tal.EligibleTalentDTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set talent ID.
	talent.ID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.TalentService.GetEligibleTalent(&talent, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, talent)
}

// UpdateTalent updates the talent by specific id
// GetTalent Return Specific Talent
// @Summary Update Talent
// @Description Update Talent By Talent ID
// @Tags Talent-Master
// @Accept  json
// @Param talent body tal.Talent true "Talent Data"
// @Param talentID path string true "Talent ID" Format(uuid.UUID)
// @Success 200 {array} []tal.Talent
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent/{talentID} [Put]
func (controller *TalentController) UpdateTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateTalent called=======================================")

	// Create bucket.
	talent := tal.Talent{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &talent)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = talent.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set talent ID to talent.
	talent.ID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to talent.
	talent.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of talent.
	talent.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.TalentService.UpdateTalent(&talent)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talent updated successfully")
}

// DeleteTalent delete talent by specific id
// @Summary Delete Talent
// @Description Delete Talent By Talent ID
// @Tags Talent-Master
// @Accept  json
// @Param talentID path string true "Talent ID"
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Success 200 {plain} plain
// @Router /talent/{talentID} [Delete]
func (controller *TalentController) DeleteTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteTalent called=======================================")

	// Create bucket.
	talent := tal.Talent{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set talent ID.
	talent.ID, err = util.ParseUUID(params["talentID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	talent.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to talent's DeletedBy field.
	talent.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.TalentService.DeleteTalent(&talent)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talent deleted successfully")
}

//*******************************************TALENT OTHER API*********************************************

// GetAllSearchTalents return all talents by search criteria
// @Summary Return Talent Based On Search Parameter
// @Description SearchTalents Return Talents Based On Search Paramater
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param limit path string true "total talent return in response"
// @Param offset path string true "Page number"
// @Param talentSearch body tal.TalentSearch true "Search Parameter For Talent Search"
// @Success 200 {array} []tal.Talent
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent/search/{limit}/{offset} [Get]
func (controller *TalentController) GetAllSearchTalents(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetAllSearchTalents called=======================================")

	// Create bucket.
	talents := []tal.DTO{}

	// Create bucket for search criteria.
	talentSearch := tal.Search{}

	// Create bucket for total talents count.
	var totalCount int

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Fill the r.Form.
	r.ParseForm()

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the talentSearch variable with given data.
	if err := web.UnmarshalJSON(r, &talentSearch); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all search talents service method.
	if err := controller.TalentService.GetAllSearchTalents(&talents, &talentSearch, tenantID, limit, offset, &totalCount,
		&totalLifetimeValue, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsByCompanyRequirement return all talents by Company Requirement ID
// @Summary Get Talent By Company Requirement ID
// @Description GetTalentsByCompanyRequirement Return All Talents By Company Requirement ID
// @Tags Talent-Master
// @Accept  json
// @Produce  json
// @Param requirementID path string true "Company Requirement ID" Format(uuid.UUID)
// @Param limit path string true "total talent return in response"
// @Param offset path string true "Page number"
// @Success 200 {array} []tal.Talent
// @Failure 400 {object} errors.ValidationError
// @Failure 500 {object} errors.HTTPError
// @Router /talent/requirement/{requirementID}/talents/{limit}/{offset} [Get]
func (controller *TalentController) GetTalentsByCompanyRequirement(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetTalentsByCompanyRequirement called=======================================")

	// Create bucket.
	talents := []tal.DTO{}

	// Create bucket for total talents count.
	var totalCount int

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Get params from api.
	params := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting requirement id from param and parsing it to uuid.
	requirementID, err := util.ParseUUID(params["requirementID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requirement id", http.StatusBadRequest))
		return
	}

	// Call get talents by requirement service method.
	if err := controller.TalentService.GetTalentsByRequirementID(&talents, requirementID, tenantID, limit, offset, &totalCount); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsForBatch returns all talents for one batch
func (controller *TalentController) GetTalentsForBatch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForBatch call==============================")
	// batch := bat.Batch{}
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting batch id from param and parsing it to uuid.
	batchID, err := util.ParseUUID(mux.Vars(r)["batchID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForBatch(&talents, tenantID, batchID, limit, offset, &totalCount, &totalLifetimeValue)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsByWaitingList returns all talents by waiting lists.
func (controller *TalentController) GetTalentsByWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsByWaitingList call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Fill the r.Form.
	r.ParseForm()

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsByWaitingList(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsByCampusDrive returns all talents by campus drive id.
func (controller *TalentController) GetTalentsByCampusDrive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsByCampusDrive call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting campus drive id from param and parsing it to uuid.
	campusDriveID, err := util.ParseUUID(mux.Vars(r)["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsByCampusDrive(&talents, tenantID, campusDriveID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsBySeminar returns all talents by seminar id.
func (controller *TalentController) GetTalentsBySeminar(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsBySeminar call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting seminar id from param and parsing it to uuid.
	seminarID, err := util.ParseUUID(mux.Vars(r)["seminarID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse seminar id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsBySeminar(&talents, tenantID, seminarID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsForProSummaryReport returns all talents for professional summary report.
func (controller *TalentController) GetTalentsForProSummaryReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForProSummaryReport call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForProSummaryReport(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsForProSummaryReportTechnologyCount returns all talents for professional summary report by technology talent count.
func (controller *TalentController) GetTalentsForProSummaryReportTechnologyCount(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForProSummaryReportTechnologyCount call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForProSummaryReportTechnologyCount(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsForFresherSummaryReport returns all talents for fresher summary report.
func (controller *TalentController) GetTalentsForFresherSummaryReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForFresherSummaryReport call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForFresherSummaryReport(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetTalentsForPackageSummaryReport returns all talents for package summary report.
func (controller *TalentController) GetTalentsForPackageSummaryReport(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForPackageSummaryReport call==============================")

	// Create bucket.
	talents := []tal.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Create bucket for total lifetime value.
	var totalLifetimeValue tal.TotalLifetimeValueResult

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForPackageSummaryReport(&talents, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Give total lifetime value in header.
	web.SetNewHeader(w, "totalLifetimeValue", strconv.Itoa(int(totalLifetimeValue.TotalLifetimeValue)))

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// UpdateTalentsSalesperson update one or more talents' salesperson.
func (controller *TalentController) UpdateTalentsSalesperson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateTalentSalesperson called=======================================")

	// Create bucket.
	talents := []tal.TalentUpdate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the talent variable with given data.
	err := web.UnmarshalJSON(r, &talents)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Getting salesperson id from param and parsing it to uuid.
	salesPersonID, err := util.ParseUUID(params["salesPersonID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse salesPerson id", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.TalentService.UpdateTalentsSalesperson(&talents, salesPersonID, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Talents' salesperson updated successfully")
}

// GetTalentBatches gets all batches for one talent
func (controller *TalentController) GetTalentBatches(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentBatches call==============================")

	// Create bucket.
	batches := []bat.BatchDTO{}

	// Get params from api.
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

	// Call get batches by on etakent method.
	err = controller.TalentService.GetTalentBatches(&batches, tenantID, talentID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, batches)
}

// GetTalentsForExcelDownload returns talents for excel download.
func (controller *TalentController) GetTalentsForExcelDownload(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetTalentsForExcelDownload call==============================")

	// Create bucket.
	talents := []tal.ExcelDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all batch talents service method.
	err = controller.TalentService.GetTalentsForExcelDownload(&talents, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

// GetSearchedTalentsForExcelDownload gets searched talents from database.
func (controller *TalentController) GetSearchedTalentsForExcelDownload(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSearchedTalentsForExcelDownload call==============================")

	// Create bucket.
	talents := []tal.ExcelDTO{}

	// Create bucket for search criteria.
	talentSearch := tal.Search{}

	// Fill the talentSearch variable with given data.
	if err := web.UnmarshalJSON(r, &talentSearch); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all batch talents service method.
	err = controller.TalentService.GetSearchedTalentsForExcelDownload(&talents, &talentSearch, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, talents)
}

//===================================================Add Credential For Talent===========================================================================

func (controller *TalentController) AddCredentialForTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddCredentialForTalent called=======================================")
	allTalents := []tal.Talent{}

	// Parse and set credentialID in CreatedBy field of talent
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	uow := repository.NewUnitOfWork(controller.TalentService.DB, false)
	err = controller.TalentService.DB.Find(&allTalents).Error
	if err != nil {
		log.NewLogger().Error(err)
		return
	}
	for _, talent := range allTalents {
		fmt.Println("TENANT ID ================>", talent.TenantID)
		// talent.Password = util.GeneratePassword()
		// talent.Password = strings.ToLower(talent.FirstName)

		talent.CreatedBy = credentialID
		err := controller.TalentService.AddTalentAndCredentialForUpdatingPassword(uow, &talent)
		if err != nil {
			log.NewLogger().Error(err)
		}
	}
	uow.Commit()
	web.RespondJSON(w, http.StatusOK, "YAYYYY")
}

//===================================================UPDATE EXPERIENCE IN MONTHS OF ALL TALENTS===========================================================================

func (controller *TalentController) UpdateExperienceInMonthsOfTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateExperienceInMonthsOfTalent called=======================================")

	// Call update experience in months of all talents service method.
	err := controller.TalentService.UpdateExperienceInMonthsOfTalent()
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Experience in months of all talents updated")
}
