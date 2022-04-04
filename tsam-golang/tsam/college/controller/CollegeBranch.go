package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	colg "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CollegeBranchController provides methods to do Update, Delete, Add, Get operations on CollegeBranch.
type CollegeBranchController struct {
	CollegeBranchService *service.CollegeBranchService
}

// NewCollegeBranchController creates new instance of CollegeBranchController.
func NewCollegeBranchController(collegeBranchService *service.CollegeBranchService) *CollegeBranchController {
	return &CollegeBranchController{
		CollegeBranchService: collegeBranchService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *CollegeBranchController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	collegeBranchList := router.HandleFunc("/tenant/{tenantID}/college/branch/list",
		controller.GetCollegeBranchList).Methods(http.MethodGet)

	collegeBranchWithLimitList := router.HandleFunc("/tenant/{tenantID}/college/branch/list/limit/{limit}/offset/{offset}",
		controller.GetCollegeBranchListWithLimit).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/college/branch/limit/{limit}/offset/{offset}",
		controller.GetAllBranches).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/college/branch/sales-person/{salesPersonID}/limit/{limit}/offset/{offset}",
		controller.GetAllBranchesForSalesPerson).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/branch",
		controller.GetAllBranchesOfCollege).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/branch/credential/{credentialID}",
		controller.AddCollegeBranch).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/branch/{branchID}",
		controller.GetCollegeBranch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/branch/{branchID}/credential/{credentialID}",
		controller.UpdateCollegeBranch).Methods(http.MethodPut)

	router.HandleFunc("/tenant/{tenantID}/college/{collegeID}/branch/{branchID}/credential/{credentialID}",
		controller.DeleteCollegeBranch).Methods(http.MethodDelete)

	//exculde routes
	*exclude = append(*exclude, collegeBranchList, collegeBranchWithLimitList)

	log.NewLogger().Info("College Branch Routes Registered")
}

// GetAllBranches returns all branches in database as response
func (controller *CollegeBranchController) GetAllBranches(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllBranches Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	fmt.Println("========Form==============================", r.Form)
	allBranches := &[]*colg.Branch{}
	err = controller.CollegeBranchService.GetAllBranches(tenantID, allBranches, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allBranches)
}

// GetCollegeBranchList returns all branches with name and ID from database.
func (controller *CollegeBranchController) GetCollegeBranchList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCollegeBranchList Called==============================")

	params := mux.Vars(r)

	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	allBranches := &[]*list.Branch{}
	err = controller.CollegeBranchService.GetCollegeBranchList(tenantID, allBranches)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, allBranches)
}

// GetCollegeBranchList returns all branches with name and ID from database.
func (controller *CollegeBranchController) GetCollegeBranchListWithLimit(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCollegeBranchListWithLimit Called==============================")

	params := mux.Vars(r)

	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	allBranches := &[]*list.Branch{}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	err = controller.CollegeBranchService.GetCollegeBranchListWithLimit(tenantID, allBranches, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allBranches)
}

// GetAllBranchesForSalesPerson returns all branches in where the given sales person is assigned.
func (controller *CollegeBranchController) GetAllBranchesForSalesPerson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllBranchesForSalesPerson Called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	salesPersonID, err := util.ParseUUID(mux.Vars(r)["salesPersonID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse sales person id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	allBranches := &[]*colg.Branch{}
	err = controller.CollegeBranchService.GetAllBranchesForSalesPerson(tenantID, salesPersonID, allBranches, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allBranches)
}

// GetAllBranchesOfCollege returns all branches of a speicifc college as response
func (controller *CollegeBranchController) GetAllBranchesOfCollege(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllBranchesOfCollege Called==============================")

	params := mux.Vars(r)
	// Parse tenant ID from url.
	tenantID, err := util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse college ID from url.
	collegeID, err := util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college id", http.StatusBadRequest))
		return
	}

	collegeBranches := &[]*colg.Branch{}
	err = controller.CollegeBranchService.GetAllBranchesOfCollege(tenantID, collegeID, collegeBranches)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, collegeBranches)
}

// GetCollegeBranch returns a specific college branch as response
func (controller *CollegeBranchController) GetCollegeBranch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetCollegeBranch Called==============================")

	var err error
	params := mux.Vars(r)
	branch := &colg.Branch{}

	// Assign Tenant ID
	branch.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Assign college ID
	branch.CollegeID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Assign college branch ID
	branch.ID, err = util.ParseUUID(params[paramCollegeBranchID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call get service
	err = controller.CollegeBranchService.GetCollegeBranch(branch)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, branch)
}

// UpdateCollegeBranch updates the specific college branch
func (controller *CollegeBranchController) UpdateCollegeBranch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateCollegeBranch Called==============================")

	branch := &colg.Branch{}
	params := mux.Vars(r)

	// Parse college branch from request.
	err := web.UnmarshalJSON(r, branch)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate College
	err = branch.ValidateCollegeBranch()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID to college.
	branch.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// UpdatedBy will be set with the given credentialID
	branch.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set college ID to branch.
	branch.CollegeID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college id", http.StatusBadRequest))
		return
	}

	// Parse and set branch ID to branch.
	branch.ID, err = util.ParseUUID(params[paramCollegeBranchID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college branch id", http.StatusBadRequest))
		return
	}

	err = controller.CollegeBranchService.UpdateCollegeBranch(branch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "College Branch Updated")
}

// AddCollegeBranch adds a new college branch
func (controller *CollegeBranchController) AddCollegeBranch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddCollegeBranch Called==============================")

	branch := &colg.Branch{}
	params := mux.Vars(r)

	// Fill the college branch variable with given data.
	err := web.UnmarshalJSON(r, branch)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	branch.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	branch.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set CollegeID to branch.
	branch.CollegeID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college id", http.StatusBadRequest))
		return
	}

	// Validate college branch.
	err = branch.ValidateCollegeBranch()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service.
	err = controller.CollegeBranchService.AddCollegeBranch(branch)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, branch.ID)
}

// DeleteCollegeBranch deletes the specific College Branch
func (controller *CollegeBranchController) DeleteCollegeBranch(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteCollege Called==============================")

	branch := &colg.Branch{}
	params := mux.Vars(r)
	var err error

	// Parse and set tenant ID.
	branch.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to branch's DeletedBy field.
	branch.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse college ID and assign to branch
	branch.CollegeID, err = util.ParseUUID(params[paramCollegeID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college id", http.StatusBadRequest))
		return
	}

	// Parse branch ID and assign it to branch
	branch.ID, err = util.ParseUUID(params[paramCollegeBranchID])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse college's branch id", http.StatusBadRequest))
		return
	}

	err = controller.CollegeBranchService.DeleteCollegeBranch(branch)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "College Branch Deleted")
}
