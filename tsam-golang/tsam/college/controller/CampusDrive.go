package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/college/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"

	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CampusDriveController provides methods to do update, delete, add, get all and get one operations on campus drive.
type CampusDriveController struct {
	CampusDriveService *service.CampusDriveService
}

// NewCampusDriveController creates new instance of CampusDriveController.
func NewCampusDriveController(campusDriveService *service.CampusDriveService) *CampusDriveController {
	return &CampusDriveController{
		CampusDriveService: campusDriveService,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *CampusDriveController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// Add one campus drive.
	router.HandleFunc("/tenant/{tenantID}/campus-drive/credential/{credentialID}",
		controller.AddCampusDrive).Methods(http.MethodPost)

	// Get all campus drive with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/campus-drive/limit/{limit}/offset/{offset}",
		controller.GetAllCampusDrives).Methods(http.MethodGet)

	// Get one campus drive by id.
	router.HandleFunc("/tenant/{tenantID}/campus-drive/{campusDriveID}",
		controller.GetCampusDrive).Methods(http.MethodGet)

	// Update one campus drive.
	router.HandleFunc("/tenant/{tenantID}/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.UpdateCampusDrive).Methods(http.MethodPut)

	// Delete one campus drive.
	router.HandleFunc("/tenant/{tenantID}/campus-drive/{campusDriveID}/credential/{credentialID}",
		controller.DeleteCampusDrive).Methods(http.MethodDelete)

	// Get one campus drive by code.
	campusDriveByCode := router.HandleFunc("/tenant/{tenantID}/campus-drive/code/{campusDriveCode}",
		controller.GetCampusDriveByCode).Methods(http.MethodGet)

	// Exculde routes.
	*exclude = append(*exclude, campusDriveByCode)

	log.NewLogger().Info("CampusDrive Routes Registered")
}

// GetAllCampusDrives gets all campus drives.
func (controller *CampusDriveController) GetAllCampusDrives(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetAllCampusDrives call**************************************")

	// Create bucket.
	campusDrives := []college.CampusDriveDTO{}

	// Create bucket for total campus drive count.
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

	// Call get campus drives method.
	err = controller.CampusDriveService.GetAllCampusDrives(&campusDrives, tenantID, limit, offset, &totalCount, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status and total count in header to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, campusDrives)
}

// GetCampusDrive gets campus drive by calling the get service.
func (controller *CampusDriveController) GetCampusDrive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCampusDrive called=======================================")

	// Create bucket.
	campusDrive := college.CampusDrive{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set campus drive ID.
	campusDrive.ID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	campusDrive.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.CampusDriveService.GetCampusDrive(&campusDrive); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, campusDrive)
}

// UpdateCampusDrive updates the campus drive by calling the update service.
func (controller *CampusDriveController) UpdateCampusDrive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================UpdateCampusDrive called=======================================")

	// Create bucket.
	campusDrive := college.CampusDrive{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the campus drive variable with given data.
	err := web.UnmarshalJSON(r, &campusDrive)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = campusDrive.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set campus drive ID to campus drive.
	campusDrive.ID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID to campus drive.
	campusDrive.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field of campus drive.
	campusDrive.UpdatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update service method.
	err = controller.CampusDriveService.UpdateCampusDrive(&campusDrive)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Campus drive updated successfully")
}

// AddCampusDrive adds new campus drive by calling the add service.
func (controller *CampusDriveController) AddCampusDrive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================AddCampusDrive called=======================================")

	// Create bucket.
	campusDrive := college.CampusDrive{}

	// Get params from api.
	params := mux.Vars(r)

	// Fill the campus dtive variable with given data.
	err := web.UnmarshalJSON(r, &campusDrive)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = campusDrive.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	campusDrive.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field of campus drive.
	campusDrive.CreatedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	if err = controller.CampusDriveService.AddCampusDrive(&campusDrive); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Campus drive added successfully")
}

// DeleteCampusDrive deletes campus drive by calling the delete service.
func (controller *CampusDriveController) DeleteCampusDrive(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================DeleteCampusDrive called=======================================")

	// Create bucket.
	campusDrive := college.CampusDrive{}

	// Get params from api.
	params := mux.Vars(r)

	// Declare err.
	var err error

	// Parse and set campus drive ID.
	campusDrive.ID, err = util.ParseUUID(params["campusDriveID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive id", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	campusDrive.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse credentialID and assign to campus drive's DeletedBy field.
	campusDrive.DeletedBy, err = util.ParseUUID(params["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.CampusDriveService.DeleteCampusDrive(&campusDrive)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Campus drive deleted successfully")

}

// GetCampusDriveByCode gets campus drive by its code.
func (controller *CampusDriveController) GetCampusDriveByCode(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===============================GetCampusDriveByCode called=======================================")

	// Create bucket.
	campusDrive := college.CampusDrive{}

	// Declare err.
	var err error

	// Get params from api.
	params := mux.Vars(r)

	// Parse and set campus drive ID.
	campusDrive.Code = (params["campusDriveCode"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse campus drive code", http.StatusBadRequest))
		return
	}

	// Parse and set tenant ID.
	campusDrive.TenantID, err = util.ParseUUID(params["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get service method.
	if err := controller.CampusDriveService.GetCampusDriveByCode(&campusDrive); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, campusDrive)
}
