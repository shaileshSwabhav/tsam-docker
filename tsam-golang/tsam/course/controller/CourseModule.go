package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// CourseModuleController provides methods to do Update, Delete, Add, Get operations on course_modules.
type CourseModuleController struct {
	log     log.Logger
	service *service.CourseModuleService
	auth    *security.Authentication
}

// NewCourseModuleController creates new instance of CourseModuleController.
func NewCourseModuleController(service *service.CourseModuleService,
	log log.Logger, auth *security.Authentication) *CourseModuleController {
	return &CourseModuleController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *CourseModuleController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}/course-module",
		controller.AddCourseModule).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}/course-module/{courseModuleID}",
		controller.UpdateCourseModule).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}/course-module/{courseModuleID}",
		controller.DeleteCourseModule).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}/course-module",
		controller.GetCourseModule).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)

	controller.log.Info("Course Module Route Registered")

}

// AddCourseModule will add new module to course.
func (controller *CourseModuleController) AddCourseModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Course Module Called==============================")
	module := course.CourseModule{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	module.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	module.CourseID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	module.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddCourseModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully added for course")
}

// UpdateCourseModule will update module to course.
func (controller *CourseModuleController) UpdateCourseModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Course Module Called==============================")
	module := course.CourseModule{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	module.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	module.CourseID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	// Parse and set module ID.
	module.ID, err = parser.GetUUID("courseModuleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in UpdatedBy field.
	module.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = module.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdateCourseModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully updated for course")
}

// DeleteCourseModule will delete module to course.
func (controller *CourseModuleController) DeleteCourseModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Delete Course Module Called==============================")
	module := course.CourseModule{}

	parser := web.NewParser(r)

	var err error

	// Parse and set tenant ID.
	module.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	module.CourseID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	// Parse and set module ID.
	module.ID, err = parser.GetUUID("courseModuleID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse module id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in DeletedBy field.
	module.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteCourseModule(&module)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Module successfully deleted")
}

// GetCourseModule will fetch the modules for the given course.
func (controller *CourseModuleController) GetCourseModule(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Course Module Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	courseID, err := parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	modules := &[]course.CourseModuleDTO{}
	var totalCount int

	err = controller.service.GetCourseModule(modules, tenantID, courseID, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, modules)
}
