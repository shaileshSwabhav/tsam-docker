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

// CourseTopicConceptController provides methods to do Update, Delete, Add, Get operations on course_topic_programming_concepts.
type CourseTopicConceptController struct {
	log     log.Logger
	service *service.CourseTopicConceptService
	auth    *security.Authentication
}

// NewCourseTopicConceptController creates new instance of CourseTopicConceptController.
func NewCourseTopicConceptController(service *service.CourseTopicConceptService,
	log log.Logger, auth *security.Authentication) *CourseTopicConceptController {
	return &CourseTopicConceptController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *CourseTopicConceptController) RegisterRoutes(router *mux.Router) {

	// add
	router.HandleFunc("/tenant/{tenantID}/module-topic/{moduleTopicID}/programming-concept",
		controller.AddTopicProgrammingConcept).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/module-topic/{moduleTopicID}/programming-concept/{courseConceptID}",
		controller.UpdateCourseProgrammingConcept).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/module-topic/{moduleTopicID}/programming-assignment/{courseconceptID}",
		controller.DeleteCourseProgrammingConcept).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/module-topic/{moduleTopicID}/programming-concepts",
		controller.GetTopicProgrammingConcepts).Methods(http.MethodGet)

	router.Use(controller.auth.Middleware)
	controller.log.Info("Course Programming Concept Routes Registered")
}

// AddTopicProgrammingConcept will add new course topic programming concept to the table.
func (controller *CourseTopicConceptController) AddTopicProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================AddTopicProgrammingConcept call==============================")

	courseConcept := course.TopicProgrammingConcept{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &courseConcept)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	courseConcept.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	courseConcept.TopicID, err = parser.GetUUID("moduleTopicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	courseConcept.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddTopicConcept(&courseConcept)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "successfully added  course topic programming concept")
}

// UpdateCourseProgrammingConcept will update programming_Concept for specified course.
func (controller *CourseTopicConceptController) UpdateCourseProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================UpdateCourseProgrammingConcept call==============================")

	topicProgrammingConcept := course.TopicProgrammingConcept{}

	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &topicProgrammingConcept)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	topicProgrammingConcept.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	topicProgrammingConcept.TopicID, err = parser.GetUUID("moduleTopicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	topicProgrammingConcept.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	controller.service.UpdateCourseProgrammingAssignment(&topicProgrammingConcept)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "successfully updated  course topic programming concept")
}

// DeleteCourseProgrammingConcept will delete programming_concept for specified course.
func (controller *CourseTopicConceptController) DeleteCourseProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================DeleteCourseProgrammingConcept call==============================")

	courseConcept := course.TopicProgrammingConcept{}

	parser := web.NewParser(r)

	var err error
	courseConcept.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	courseConcept.TopicID, err = parser.GetUUID("moduleTopicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	courseConcept.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteCourseProgrammingConcept(&courseConcept)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "successfully deleted course topic programming concept")
}

// GetTopicProgrammingConcepts will fetch the programming concepts with limit and offset.
func (controller *CourseTopicConceptController) GetTopicProgrammingConcepts(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================GetTopicProgrammingConcepts call==============================")

	courseConcept := []course.TopicProgrammingConceptDTO{}

	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set course ID.
	moduleTopicID, err := parser.GetUUID("moduleTopicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse course id", http.StatusBadRequest))
		return
	}

	var totalCount int
	err = controller.service.GetTopicProgrammingConcepts(tenantID, moduleTopicID, &courseConcept, parser, &totalCount)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, courseConcept)
}
