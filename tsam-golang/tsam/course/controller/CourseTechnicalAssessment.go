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

// AssessmentController provides method to update, delete, add, get method for aha_moment.
type AssessmentController struct {
	AssessmentService *service.AssessmentService
	log               log.Logger
	auth              *security.Authentication
}

// NewAssessmentController creates new instance of AssessmentController.
func NewAssessmentController(assessmentService *service.AssessmentService, log log.Logger, auth *security.Authentication) *AssessmentController {
	return &AssessmentController{
		AssessmentService: assessmentService,
		log:               log,
		auth:              auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *AssessmentController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/course-technical-assessment/faculty/{facultyID}",
		controller.AddAssessments).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/course-technical-assessment/{assessmentID}/faculty/{facultyID}",
		controller.UpdateAssessment).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/course-technical-assessment/{assessmentID}",
		controller.DeleteAssessment).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/course-technical-assessment", controller.GetAllAssessments).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/course-technical-assessment/faculty/{facultyID}",
		controller.GetAssessmentsForFaculty).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)
	controller.log.Info("Course Technical Assessment Routes Registered")
}

// AddAssessments will add multiple course technical assessments to the table.
func (controller *AssessmentController) AddAssessments(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================AddAssessments called===========================")
	// param := mux.Vars(r)
	assessments := []course.CourseTechnicalAssessment{}

	err := web.UnmarshalJSON(r, &assessments)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}
	parser := web.NewParser(r)

	// util.ParseUUID(param["credentialID"])
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(param["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(param["facultyID"])
	facultyID, err := parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
		return
	}

	// courseID, err := util.ParseUUID(param["courseID"])
	// if err != nil {
	// 	controller.log.Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	for _, assessment := range assessments {
		err = assessment.Validate()
		if err != nil {
			controller.log.Error(err)
			web.RespondError(w, err)
			return
		}
	}

	err = controller.AssessmentService.AddAssessments(&assessments, tenantID, credentialID, facultyID)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Course Technical Assessment successfully added")
}

// UpdateAssessment will update specified assessment in the table.
func (controller *AssessmentController) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================UpdateAssessment called===========================")
	// param := mux.Vars(r)
	assessment := course.CourseTechnicalAssessment{}

	err := web.UnmarshalJSON(r, &assessment)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}
	parser := web.NewParser(r)
	// util.ParseUUID(param["credentialID"])
	assessment.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	assessment.TenantID, err = parser.GetTenantID()
	// util.ParseUUID(param["tenantID"])
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	assessment.ID, err = parser.GetUUID("assessmentID")
	// util.ParseUUID(param["assessmentID"])
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse assessment id", http.StatusBadRequest))
		return
	}

	assessment.FacultyID, err = parser.GetUUID("facultyID")
	// util.ParseUUID(param["facultyID"])
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
		return
	}

	// courseID, err := util.ParseUUID(param["courseID"])
	// if err != nil {
	// 	controller.log.Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	err = assessment.Validate()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.AssessmentService.UpdateAssessment(&assessment)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Course Technical Assessment successfully updated")
}

// DeleteAssessment will delete specified assessment from the table.
func (controller *AssessmentController) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================DeleteAssessment called===========================")
	// param := mux.Vars(r)
	var err error
	assessment := course.CourseTechnicalAssessment{}

	parser := web.NewParser(r)

	// util.ParseUUID(param["credentialID"])
	assessment.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	// util.ParseUUID(param["tenantID"])
	assessment.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(param["assessmentID"])
	assessment.ID, err = parser.GetUUID("assessmentID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse assesment id", http.StatusBadRequest))
		return
	}

	err = controller.AssessmentService.DeleteAssessment(&assessment)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Course Technical Assessment successfully deleted")
}

// GetAllAssessments will return all the technical assessments.
func (controller *AssessmentController) GetAllAssessments(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================GetAllAssessments called===========================")
	// param := mux.Vars(r)
	// var err error
	assessments := []course.CourseTechnicalAssessmentDTO{}

	parser := web.NewParser(r)
	// util.ParseUUID(param["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	err = controller.AssessmentService.GetAllAssessments(&assessments, tenantID)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, assessments)
}

// GetAssessmentsForFaculty will return all the technical assessments for the specified faculty.
func (controller *AssessmentController) GetAssessmentsForFaculty(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("===========================GetAssessmentsForFaculty called===========================")
	// param := mux.Vars(r)
	// var err error
	assessments := []course.CourseTechnicalAssessmentDTO{}

	parser := web.NewParser(r)

	// util.ParseUUID(param["tenantID"])
	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// util.ParseUUID(param["facultyID"])
	facultyID, err := parser.GetUUID("facultyID")
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, errors.NewHTTPError("unable to parse faculty id", http.StatusBadRequest))
		return
	}

	err = controller.AssessmentService.GetAssessmentsForFaculty(&assessments, tenantID, facultyID)
	if err != nil {
		controller.log.Error(err)
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, assessments)
}
