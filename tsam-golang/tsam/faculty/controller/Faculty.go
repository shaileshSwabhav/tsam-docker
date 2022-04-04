package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/faculty/service"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"

	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/web"
)

// FacultyController Provide method to Update, Delete, Add, Get Method For faculty.
type FacultyController struct {
	FacultyService *service.FacultyService
	log            log.Logger
	auth           *security.Authentication
}

// NewFacultyController Create New Instance Of FacultyController.
func NewFacultyController(ser *service.FacultyService, log log.Logger, auth *security.Authentication) *FacultyController {
	return &FacultyController{
		FacultyService: ser,
		log:            log,
		auth:           auth,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *FacultyController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/faculty",
		controller.AddFaculty).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/faculties",
		controller.AddFaculties).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/faculty/{facultyID}",
		controller.UpdateFaculty).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/faculty/{facultyID}",
		controller.DeleteFaculty).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/faculty-list", controller.GetFacultyList).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/faculty", controller.GetFaculties).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/faculty/search/limit/{limit}/offset/{offset}",
	// 	controller.GetSearchedFaculties).Methods(http.MethodPost)

	router.HandleFunc("/tenant/{tenantID}/faculty/{facultyID}",
		controller.GetFaculty).Methods(http.MethodGet)

	// Get faculty for batch.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/faculty",
		controller.GetFacultyForBatch).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/faculty/credential/list",
		controller.GetFacultyListFromCredential).Methods(http.MethodGet)

	controller.log.Info("Faculty Route Registered")
}

//AddFaculty Add New faculty
func (controller *FacultyController) AddFaculty(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddFaculty call==============================")

	faculty := fct.Faculty{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &faculty)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested Data", http.StatusBadRequest))
		return
	}

	faculty.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	faculty.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = faculty.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	err = controller.FacultyService.AddFaculty(&faculty)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Faculty successfully added")
}

//AddFaculties adds new faculty by calling add service
func (controller *FacultyController) AddFaculties(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddFaculties call==============================")

	facultiesIDs := []uuid.UUID{}
	faculties := []fct.Faculty{}

	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// util.ParseUUID(params["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	// util.ParseUUID(params["credentialID"])
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse faculties from request
	err = web.UnmarshalJSON(r, faculties)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = controller.FacultyService.AddFaculties(&faculties, &facultiesIDs, tenantID, credentialID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tempIDs := []uuid.UUID{}
	for _, faculty := range faculties {
		tempIDs = append(tempIDs, faculty.ID)
	}
	fmt.Println(tempIDs)

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Faculties successfully added")
}

//UpdateFaculty Update The faculty
func (controller *FacultyController) UpdateFaculty(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateFaculty Call==============================")

	faculty := fct.Faculty{}

	var err error

	// params := mux.Vars(r)
	parser := web.NewParser(r)

	// util.ParseUUID(params["tenantID"])
	faculty.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["facultyID"])
	faculty.ID, err = parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse faculty id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["credentialID"])
	faculty.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = web.UnmarshalJSON(r, &faculty)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	err = faculty.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// faculty.ID = facultyID
	// faculty.TenantID = tenantID
	err = controller.FacultyService.UpdateFaculty(&faculty)
	if err != nil {
		fmt.Println(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Faculty Updated")

}

//DeleteFaculty will delete the specified faculty.
func (controller *FacultyController) DeleteFaculty(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteFaculty Call==============================")

	var err error

	// params := mux.Vars(r)
	faculty := fct.Faculty{}

	parser := web.NewParser(r)

	// util.ParseUUID(params["tenantID"])
	faculty.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["facultyID"])
	faculty.ID, err = parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["credentialID"])
	faculty.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.FacultyService.DeleteFaculty(&faculty)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Faculty Deleted")

}

// GetFacultyList returns faculty list.
func (controller *FacultyController) GetFacultyList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetFacultyList call==============================")

	faculties := &[]*list.Faculty{}

	parser := web.NewParser(r)
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.FacultyService.GetFacultyList(faculties, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, faculties)
}

// GetFaculties returns all faculties from DB as response.
func (controller *FacultyController) GetFaculties(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetAllFaculty call==============================")

	// faculty := &[]fct.Faculty{}
	faculty := &[]fct.FacultyDTO{}

	parser := web.NewParser(r)
	// util.ParseUUID(mux.Vars(r)[tenantID])

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse Tenant ID", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	// r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	// limit, offset := web.GetLimitAndOffset(r)

	err = controller.FacultyService.GetFaculties(faculty, tenantID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, faculty)
}

// //GetSearchedFaculties Return All Searched faculties
// func (controller *FacultyController) GetSearchedFaculties(w http.ResponseWriter, r *http.Request) {
// controller.log.Info("==============================GetSearchFaculty Call==============================")
// 	searchFaculty := &fct.Search{}
// 	faculty := &[]fct.Faculty{}
// 	tenantID, err := util.ParseUUID(mux.Vars(r)[tenantID])
// 	if err != nil {
// 	controller.log.Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}
// 	err = web.UnmarshalJSON(r, searchFaculty)
// 	if err != nil {
// 	controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// limit, offset and totalCount are used for pagination
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)
// 	err = controller.FacultyService.GetSearchedFaculties(faculty, searchFaculty, tenantID, limit, offset, &totalCount)
// 	if err != nil {
// 	controller.log.Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// Writing Response with OK Status to ResponseWriter
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, faculty)
// }

//GetFaculty Return Specific faculty
func (controller *FacultyController) GetFaculty(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetFaculty Call==============================")

	faculty := &fct.Faculty{}

	var err error
	// params := mux.Vars(r)
	// util.ParseUUID(params["tenantID"])
	parser := web.NewParser(r)

	faculty.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(params["facultyID"])
	faculty.ID, err = parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse faculty id", http.StatusBadRequest))
		return
	}

	// faculty.ID = facultyID
	err = controller.FacultyService.GetFaculty(faculty)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, faculty)
}

// GetFacultyForBatch Return Specific faculty for batch.
func (controller *FacultyController) GetFacultyForBatch(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetFacultyForBatch Call==============================")

	faculty := &list.Faculty{}

	var err error
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse batch id", http.StatusBadRequest))
		return
	}

	err = controller.FacultyService.GetFacultyForBatch(faculty, tenantID, batchID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, faculty)
}

//GetFacultyListFromCredential returns faculty list from credential table
func (controller *FacultyController) GetFacultyListFromCredential(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetFacultyListFromCredential Call==============================")

	faculty := &[]list.FacultyCredentialDTO{}
	// params := mux.Vars(r)
	parser := web.NewParser(r)
	// util.ParseUUID(params["tenantID"])

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.FacultyService.GetFacultyListFromCredential(faculty, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, faculty)
}
