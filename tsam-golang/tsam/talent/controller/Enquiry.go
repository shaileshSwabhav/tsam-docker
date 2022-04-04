package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	talenq "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/talent/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentEnquiryController provides methods to  update, delete, add, get all, get all by salesperson,
// get one by id and get all by company requiremtn for enquiry.
type TalentEnquiryController struct {
	EnquiryService *service.EnquiryService
}

// NewEnquiryController creates and returns an instance of TalentEnquiryController.
func NewEnquiryController(enquiryservice *service.EnquiryService) *TalentEnquiryController {
	return &TalentEnquiryController{
		EnquiryService: enquiryservice,
	}
}

// RegisterRoutes registers all endpoints of this controller.
func (controller *TalentEnquiryController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Get all enquiries with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/limit/{limit}/offset/{offset}",
		controller.GetEnquiries).Methods(http.MethodGet)

	// Search all enquiries with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/search/limit/{limit}/offset/{offset}",
		controller.GetAllSearchEnquiries).Methods(http.MethodPost)

	// Get one enquiry.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/{enquiryID}",
		controller.GetEnquiry).Methods(http.MethodGet)

	// Add one enquiry.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/credential/{credentialID}",
		controller.AddEnquiry).Methods(http.MethodPost)

	// Add enquiries.
	router.HandleFunc("/tenant/{tenantID}/enquiries/credential/{credentialID}",
		controller.AddEnquiries).Methods(http.MethodPost)

	// Add enquiry from excel .
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry-excel/credential/{credentialID}",
		controller.AddEnquiryFromExcel).Methods(http.MethodPost)

	// Add one enquiry from enquiry form.
	addEnquiryForm := router.HandleFunc("/tenant/{tenantID}/talent-enquiry-form", controller.AddEnquiryForm).Methods(http.MethodPost)

	// Add one enquiry from application form.
	addApplicationForm := router.HandleFunc("/tenant/{tenantID}/talent-enquiry-application-form", controller.AddApplicationForm).Methods(http.MethodPost)

	// Update one enquiry.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/{enquiryID}/credential/{credentialID}",
		controller.UpdateEnquiry).Methods(http.MethodPut)

	// Convert enquiry to talent.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/{enquiryID}/convert-to-talent/credential/{credentialID}",
		controller.ConvertToTalent).Methods(http.MethodPut)

	// Delete one enqyiry.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/{enquiryID}/credential/{credentialID}",
		controller.DeleteEnquiry).Methods(http.MethodDelete)

	// Update all enquiries' salesperson.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/saleperson/{salesPersonID}/credential/{credentialID}",
		controller.UpdateEnquiriesSalesperson).Methods(http.MethodPut)

	// Get enquiries by waiting list.
	router.HandleFunc("/tenant/{tenantID}/talent-enquiry/waiting-list/limit/{limit}/offset/{offset}",
		controller.GetEnquiriesByWaitingList).Methods(http.MethodGet)

	// Exculde routes.
	*exclude = append(*exclude, addEnquiryForm, addApplicationForm)

	log.NewLogger().Info("Enquiry Route Registered")
}

// GetEnquiries godoc
// GetEnquiries Return All Enquiries
// @Summary Get All Enquiries
// @Description GetEnquiries Return all Enquiry By Limit and Offset
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200 {array} []talenq.TalentEnquiry
// @Router /enquiry/{limit}/{offset} [get]
func (controller *TalentEnquiryController) GetEnquiries(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetEnquiries called=======================================")

	// Create bucket.
	enquiries := []talenq.DTO{}

	// Create bucket for total enquiries count.
	var totalCount int

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

	err = controller.EnquiryService.GetEnquiries(&enquiries, tenantID, limit, offset, &totalCount, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, enquiries)
}

// GetEnquiry Return Specific Enquiry
// GetEnquiry godoc
// GetEnquiry Return Enquiry
// @Summary Get Enquiry
// @Description GetEnquiry Return Enquiry By Enquiry ID
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200 {array} talenq.TalentEnquiry
// @Router /enquiry/{enquiryID} [get]
func (controller *TalentEnquiryController) GetEnquiry(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetEnquiry called=======================================")

	// Create bucket.
	enquiry := talenq.DTO{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set enquiry ID.
	enquiry.ID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.EnquiryService.GetEnquiry(&enquiry, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, enquiry)
}

// GetAllSearchEnquiries godoc
// GetAllSearchEnquiries Return Enquiries Based On Search Paramater
// @Summary Get Enquiry
// @Description GetAllSearchEnquiries Return Enquiries Based On Search Paramater
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200 {array} []talenq.TalentEnquiry
// @Router /enquiry/search/{limit}/{offset} [post]
func (controller *TalentEnquiryController) GetAllSearchEnquiries(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetAllSearchEnquiries called=======================================")

	// Create bucket.
	enquiries := []talenq.DTO{}

	// Create bucket for search criteria.
	enquirySearch := talenq.Search{}

	// Create bucket for total enquiries count.
	var totalCount int

	// Fill the r.Form.
	r.ParseForm()

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the enquirySearch variable with given data.
	err = web.UnmarshalJSON(r, &enquirySearch)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewValidationError(errors.ErrorCodeInvalidJSON))
		return
	}

	// Get limit and offset from param and convert it to int.
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all search enquiries service method.
	if err := controller.EnquiryService.GetAllSearchEnquiries(&enquiries, &enquirySearch, tenantID, limit, offset, &totalCount, r.Form); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, enquiries)
}

// UpdateEnquiry Update The Enquiry
// @Summary Update Enquiry
// @Description Update Enquiry
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200
// @Router /enquiry/{enquiryID} [put]
func (controller *TalentEnquiryController) UpdateEnquiry(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateEnquiry called=======================================")

	// Create bucket.
	enquiry := talenq.Enquiry{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = enquiry.Validate(false)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set enquiry ID to enquiry.
	enquiry.ID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to enquiry.
	enquiry.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of enquiry.
	enquiry.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.EnquiryService.UpdateEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry updated successfully")
}

// ConvertToTalent
// @Summary ConvertToTalent
// @Description converts enquiry to talent
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200
// @Router /tenant/{tenantID}/talent-enquiry/{enquiryID}/convert-to-talent/credential/{credentialID} [put]
func (controller *TalentEnquiryController) ConvertToTalent(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================ConvertToTalent called=======================================")

	// Create bucket.
	enquiry := talenq.Enquiry{}

	// Get params from api.
	params := mux.Vars(r)
	var err error

	// Parse and set enquiry ID to enquiry.
	enquiry.ID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to enquiry.
	enquiry.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of enquiry.
	enquiry.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = controller.EnquiryService.ConvertToTalent(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry converted to talent successfully.")
}

// AddEnquiry Add New Enquiry
// @Summary Add New Enquiry
// @Description Add New Enquiry
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200
// @Router /enquiry [post]
func (controller *TalentEnquiryController) AddEnquiry(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddEnquiry called=======================================")

	// Create bucket.
	enquiry := talenq.Enquiry{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiry.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of enquiry.
	enquiry.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.EnquiryService.AddEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry added successfully")
}

// AddEnquiryForm add one enquiry, enquiry coming from enquiry form.
func (controller *TalentEnquiryController) AddEnquiryForm(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddEnquiryForm called=======================================")

	// Create bucket.
	enquiry := talenq.Enquiry{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiry)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiry.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.EnquiryService.AddEnquiryForm(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry added successfully")
}

// AddApplicationForm add one enquiry and waiting list, enquiry coming from enquiry form.
func (controller *TalentEnquiryController) AddApplicationForm(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddApplicationForm called=======================================")

	// Create bucket.
	applicationForm := talenq.ApplicationForm{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &applicationForm)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.EnquiryService.AddApplicationForm(&applicationForm, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry added successfully")
}

// DeleteEnquiry Delete Enquiry
// @Summary DeleteEnquiry Enquiry
// @Description Delete Enquiry By ID
// @Tags Enquiry-Master
// @Accept  json
// @Produce  json
// @Success 200
// @Router /enquiry/{enquiryID} [delete]
func (controller *TalentEnquiryController) DeleteEnquiry(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteEnquiry called=======================================")

	// Create bucket.
	enquiry := talenq.Enquiry{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set enquiry ID.
	enquiry.ID, err = util.ParseUUID(params["enquiryID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse enquiry id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	enquiry.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to enquiry's DeletedBy field.
	enquiry.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.EnquiryService.DeleteEnquiry(&enquiry)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry deleted successfully")
}

// UpdateEnquiriesSalesperson update one or more enquiries' salesperson
func (controller *TalentEnquiryController) UpdateEnquiriesSalesperson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateEnquiriesSalesperson called=======================================")

	// Create bucket.
	enquiries := []talenq.EnquiryUpdate{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiries variable with given data.
	err := web.UnmarshalJSON(r, &enquiries)
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

	// getting credential id from param and parsing it to uuid.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.EnquiryService.UpdateEnquiriesSalesperson(&enquiries, salesPersonID, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiries' salesperson updated successfully")
}

// GetEnquiriesByWaitingList returns all enquiries by waiting lists.
func (controller *TalentEnquiryController) GetEnquiriesByWaitingList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetEnquiriesByWaitingList call==============================")

	// Create bucket.
	enquiries := []talenq.DTO{}

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
	err = controller.EnquiryService.GetEnquiriesByWaitingList(&enquiries, tenantID, limit, offset, &totalCount, &totalLifetimeValue, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, enquiries)
}

// AddEnquiries adds multiple enquiries.
func (controller *TalentEnquiryController) AddEnquiries(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddEnquiries called=======================================")

	// Create bucket for enquiry successfully added count.
	enquiryAddedCount := 0

	// Create bucket for enquiries to be added.
	enquiries := []talenq.EnquiryExcel{}

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

	// Fill the enquiries variable with given data.
	err = web.UnmarshalJSON(r, &enquiries)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	for _, enquiry := range enquiries {
		err := enquiry.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
			return
		}
	}

	// Call add enquiries service method.
	err = controller.EnquiryService.AddEnquiries(&enquiries, &enquiryAddedCount, tenantID, &errorList, credentialID)
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
		Message:   fmt.Sprintf("Enquiries added: %v", enquiryAddedCount),
	}

	// If error while adding individual enquiry send error list.
	if len(errorList) > 0 {
		web.RespondJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiries added successfully")
}

// AddEnquiryFromExcel adds one enquiry from excel.
func (controller *TalentEnquiryController) AddEnquiryFromExcel(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddEnquiryFromExcel called=======================================")

	// Create bucket.
	enquiryExcel := talenq.EnquiryExcel{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the enquiry variable with given data.
	err := web.UnmarshalJSON(r, &enquiryExcel)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := enquiryExcel.Validate(); err != nil {
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

	// Parse and set credentialID in CreatedBy field of enquiry.
	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add enquiry from excel service method.
	if err = controller.EnquiryService.AddEnquiryFromExcel(&enquiryExcel, tenantID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Enquiry from excel added successfully")
}
