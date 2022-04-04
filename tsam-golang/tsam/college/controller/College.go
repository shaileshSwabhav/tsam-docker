package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	colg "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CollegeController provides methods to do Update, Delete, Add, Get operations on college.
type CollegeController struct {
	CollegeService *service.CollegeService
}

// NewCollegeController creates new instance of CollegeController.
func NewCollegeController(collegeService *service.CollegeService) *CollegeController {
	return &CollegeController{
		CollegeService: collegeService,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *CollegeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/college/list", controller.GetCollegeList).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/college/credential/{credentialID}", controller.AddCollege).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/colleges/credential/{credentialID}", controller.AddColleges).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/college/limit/{limit}/offset/{offset}", controller.GetAllColleges).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}", controller.GetCollege).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/credential/{credentialID}", controller.UpdateCollege).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/credential/{credentialID}", controller.DeleteCollege).Methods(http.MethodDelete)
	log.NewLogger().Info("College Routes Registered")
}

// GetAllColleges gets all colleges by calling the get all service
func (controller *CollegeController) GetAllColleges(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllColleges Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)
	allColleges := &[]colg.College{}

	// Fills the form.
	r.ParseForm()

	err = controller.CollegeService.GetAllColleges(tenantID, allColleges, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allColleges)
}

// GetCollege gets college by calling the get service
func (controller *CollegeController) GetCollege(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCollege Called==============================")
	var err error
	param := mux.Vars(r)
	college := &colg.College{}

	// Assign Tenant ID
	college.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Assign college ID
	college.ID, err = util.ParseUUID(param[paramCollegeID])
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Call get service
	err = controller.CollegeService.GetCollege(college)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, college)
}

// UpdateCollege updates the college by calling the update service
func (controller *CollegeController) UpdateCollege(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateCollege Called==============================")
	college := &colg.College{}
	params := mux.Vars(r)

	// Parse college from request.
	err := web.UnmarshalJSON(r, college)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Validate College
	err = college.ValidateCollege()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID to college.
	college.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID
	college.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	college.ID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.CollegeService.UpdateCollege(college)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "College Updated")
}

// AddCollege adds new college by calling the add service
func (controller *CollegeController) AddCollege(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddCollege Called==============================")
	college := &colg.College{}
	params := mux.Vars(r)

	// Fill the college variable with given data.
	err := web.UnmarshalJSON(r, college)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	college.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	college.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate College
	err = college.ValidateCollege()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service to add college to DB.
	err = controller.CollegeService.AddCollege(college)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, college.ID)
}

// DeleteCollege deletes college by calling the delete service
func (controller *CollegeController) DeleteCollege(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteCollege Called==============================")

	college := &colg.College{}
	params := mux.Vars(r)
	var err error

	// Parse and set tenant ID.
	college.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to college's DeletedBy field.
	college.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse college ID and assign to college
	college.ID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college id", http.StatusBadRequest))
		return
	}

	// Calling service.
	err = controller.CollegeService.DeleteCollege(college)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "College deleted")
}

// GetCollegeList returns listing of all the colleges
func (controller *CollegeController) GetCollegeList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCollegeList call==============================")
	colleges := []list.College{}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.CollegeService.GetCollegeList(&colleges, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, colleges)
}

// AddColleges if used to add multiple colleges.
func (controller *CollegeController) AddColleges(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddColleges called==============================")

	params := mux.Vars(r)
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	credentialID, err := util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	colleges := []*colg.College{}
	// Parse College from request.
	err = web.UnmarshalJSON(r, &colleges)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewValidationError(err.Error()))
		return
	}

	// Validate every College entry.
	for _, college := range colleges {
		if err := college.ValidateCollege(); err != nil {
			web.RespondError(w, err)
			return
		}
		college.CreatedBy = credentialID
		college.TenantID = tenantID
	}
	fmt.Println("================service about to be called")
	err = controller.CollegeService.AddColleges(colleges)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// IDCollection will have the list of the UUIDs of the newly added Colleges.
	IDCollection := []uuid.UUID{}
	for _, college := range colleges {
		if college != nil {
			IDCollection = append(IDCollection, college.ID)
		}
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, IDCollection)
}
