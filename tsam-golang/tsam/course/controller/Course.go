package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/course/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// CourseController Provide method to Update, Delete, Add, Get Method For course.
type CourseController struct {
	log           log.Logger
	CourseService *service.CourseService
	auth          *security.Authentication
}

// NewCourseController Create New Instance Of CourseController.
func NewCourseController(courseService *service.CourseService, log log.Logger, auth *security.Authentication) *CourseController {
	return &CourseController{
		CourseService: courseService,
		log:           log,
		auth:          auth,
	}
}

// RegisterRoutes Register All Endpoints To Router excluding a few endpoints from token check.
func (controller *CourseController) RegisterRoutes(router *mux.Router) {

	unguarded := router.PathPrefix("/tenant/{tenantID}").Subrouter()
	// add
	router.HandleFunc("/tenant/{tenantID}/course",
		controller.AddCourse).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}",
		controller.UpdateCourse).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/course/{courseID}",
		controller.DeleteCourse).Methods(http.MethodDelete)

	// get
	unguarded.HandleFunc("/course-list",
		controller.GetCoursesList).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/course",
		controller.GetCourses).Methods(http.MethodGet)

	router.HandleFunc("/tenant/{tenantID}/course/{courseID}",
		controller.GetCourse).Methods(http.MethodGet)

	// Get course details for student login.
	router.HandleFunc("/tenant/{tenantID}/course/details/{courseID}",
		controller.GetCourseDetails).Methods(http.MethodGet)

	// Get course minimum details for student login.
	router.HandleFunc("/tenant/{tenantID}/course/talent/{talentID}",
		controller.GetCourseMinimumDetails).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)

	// router.HandleFunc("/tenant/{tenantID}/course/search/limit/{limit}/offset/{offset}", controller.GetSearchedCourses).Methods(http.MethodPost)

	controller.log.Info("Course Route Registered")
	// log.NewLogger().Info("Course Route Registered")
}

// AddCourse adds new course to DB
func (controller *CourseController) AddCourse(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddCourse call==============================")
	// log.NewLogger().Info("==============================AddCourse call==============================")
	course := &crs.Course{}
	parser := web.NewParser(r)

	var err error

	err = web.UnmarshalJSON(r, course)
	if err != nil {
		controller.log.Error(err.Error())
		// log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	course.TenantID, err = parser.GetTenantID()
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["credentialID"])
	course.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = course.Validate()
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.CourseService.AddCourse(course)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Programming assignment successfully added")
}

// UpdateCourse updates the specified course
func (controller *CourseController) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateCourse call==============================")
	// log.NewLogger().Info("==============================UpdateCourse call==============================")
	course := crs.Course{}
	parser := web.NewParser(r)
	var err error

	err = web.UnmarshalJSON(r, &course)
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(mux.Vars(r)["tenantID"])
	course.TenantID, err = parser.GetTenantID()
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["courseID"])
	course.ID, err = parser.GetUUID("courseID")
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse course id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["credentialID"])
	course.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		// log.NewLogger().Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = course.Validate()
	if err != nil {
		// controller.log.Error(err.Error())
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	err = controller.CourseService.UpdateCourse(&course)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Course updated")
}

// DeleteCourse deletes the specified course
func (controller *CourseController) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteCourse call==============================")
	// controller.log.Info("==============================DeleteCourse call==============================")
	course := &crs.Course{}
	parser := web.NewParser(r)
	var err error

	// util.ParseUUID(mux.Vars(r)["tenantID"])
	course.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["courseID"])
	course.ID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse course id", http.StatusBadRequest))
		return
	}

	course.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	// util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse credential id", http.StatusBadRequest))
		return
	}

	err = controller.CourseService.DeleteCourse(course)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Course Deleted")
}

// GetCoursesList returns courses list
func (controller *CourseController) GetCoursesList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCoursesList call==============================")
	courses := &[]list.Course{}

	parser := web.NewParser(r)

	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	err = controller.CourseService.GetCoursesList(courses, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, courses)
}

// GetCourses returns all the courses
func (controller *CourseController) GetCourses(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCourses call==============================")
	courses := &[]crs.CourseDTO{}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	var totalCount int

	err = controller.CourseService.GetCourses(courses, tenantID, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, courses)
}

// GetCourse returns specified course
func (controller *CourseController) GetCourse(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCourse call==============================")
	course := &crs.CourseDTO{}

	parser := web.NewParser(r)

	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	// util.ParseUUID(mux.Vars(r)["courseID"])
	course.ID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}
	err = controller.CourseService.GetCourse(tenantID, course)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, course)
}

// GetCourseDetails returns specific details of single course by id.
func (controller *CourseController) GetCourseDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCourseDetails call==============================")
	course := crs.CourseDetails{}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	// util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["courseID"])
	course.ID, err = parser.GetUUID("courseID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	err = controller.CourseService.GetCourseDetails(&course, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, course)
}

// GetCourseMinimumDetails returns minimum specific details of single course by id.
func (controller *CourseController) GetCourseMinimumDetails(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetCourseMinimumDetails call==============================")
	courses := []crs.CourseMinimumDetails{}

	parser := web.NewParser(r)

	// util.ParseUUID(mux.Vars(r)["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(mux.Vars(r)["talentID"])
	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse talent id", http.StatusBadRequest))
		return
	}

	err = controller.CourseService.GetCourseMinimumDetails(&courses, talentID, tenantID)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, courses)
}

// 	// Writing Response with OK Status to ResponseWriter
// 	web.RespondJSON(w, http.StatusOK, sessions)
// }

// GetSearchedCourses returns searched courses
// func (controller *CourseController) GetSearchedCourses(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("==============================GetSearchedCourses call==============================")
// 	courses := &[]crs.Course{}
// 	courseSearch := &crs.Search{}

// 	// parse course to be searched
// 	err := web.UnmarshalJSON(r, courseSearch)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}

// 	// parse tenantID
// 	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
// 		return
// 	}

// 	// limit,offset & totalCount for pagination
// 	var totalCount int
// 	limit, offset := web.GetLimitAndOffset(r)

// 	err = controller.CourseService.GetSearchedCourses(courses, courseSearch, tenantID, limit, offset, &totalCount)
// 	if err != nil {
// 		log.NewLogger().Error(err.Error())
// 		web.RespondError(w, err)
// 		return
// 	}
// 	// Writing Response with OK Status to ResponseWriter
// 	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, courses)
// }
