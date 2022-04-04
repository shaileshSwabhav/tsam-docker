package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchTopicAssignmentController provides methods to do Update, Delete, Add, Get operations on batch_topic_assignments.
type BatchTopicAssignmentController struct {
	log     log.Logger
	service *service.BatchTopicAssignmentService
	auth    *security.Authentication
}

// NewBatchTopicAssignmentController creates new instance of BatchTopicAssignmentController.
func NewBatchTopicAssignmentController(service *service.BatchTopicAssignmentService,
	log log.Logger, auth *security.Authentication) *BatchTopicAssignmentController {
	return &BatchTopicAssignmentController{
		service: service,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes registers all endpoints To router.
func (controller *BatchTopicAssignmentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add.
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{topicID}/assignment",
		controller.AddBatchTopicAssignment).Methods(http.MethodPost)

	// update
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{topicID}/assignment/{topicAssignmentID}",
		controller.UpdateTopicAssignment).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/{topicID}/assignment/{topicAssignmentID}",
		controller.DeleteTopicAssignment).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic-assignment-list",
		controller.GetSessionAssignmentList).Methods(http.MethodGet)

	//get assignments by subTopic ids
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/assignments",
		controller.GetAllTopicAssignments).Methods(http.MethodGet)

	// get with score can be optional using params if needed.
	// future endpoint : ("/tenant/{tenantID}/batch/{batchID}/faculty/{facultyID}/topic-assignments"
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic-assignments",
		controller.GetAllAssignmentsWithSubmissions).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/batch-topic/{batchTopicID}/programming-assignment",
	// 	controller.GetSessionProgrammingAssignment).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/programming-assignment",
	// 	controller.GetBatchProgrammingAssignment).Methods(http.MethodGet)

	// router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/topic/programming-assignment",
	// 	controller.GetLatestBatchProgrammingAssignment).Methods(http.MethodGet)

	controller.log.Info("Batch Topic Assignment Routes Registered")
}

// AddBatchTopicAssignment will add new assignment to topic assigned in batch.
func (controller *BatchTopicAssignmentController) AddBatchTopicAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add Topic Programming Assignmnet Called==============================")

	topicAssignment := batch.TopicAssignment{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &topicAssignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	topicAssignment.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topicAssignment.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topicAssignment.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	topicAssignment.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = topicAssignment.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.AddBatchTopicAssignment(&topicAssignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Assignment successfully added for batch topic")
}

// UpdateTopicAssignment will add update assignment assigned to topic.
func (controller *BatchTopicAssignmentController) UpdateTopicAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Update Topic Programming Assignmnet Called==============================")

	topicAssignment := batch.TopicAssignment{}
	parser := web.NewParser(r)

	err := web.UnmarshalJSON(r, &topicAssignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	topicAssignment.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topicAssignment.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topicAssignment.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	topicAssignment.ID, err = parser.GetUUID("topicAssignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to topic assignment id", http.StatusBadRequest))
		return
	}

	topicAssignment.UpdatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = topicAssignment.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.UpdatedTopicAssignment(&topicAssignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Assignment successfully updated for batch topic")
}

// DeleteTopicAssignment will add delete assignment assigned to batch topic.
func (controller *BatchTopicAssignmentController) DeleteTopicAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Delete Topic Programming Assignmnet Called==============================")

	topicAssignment := batch.TopicAssignment{}
	parser := web.NewParser(r)
	var err error

	topicAssignment.TenantID, err = parser.GetUUID("tenantID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	topicAssignment.TopicID, err = parser.GetUUID("topicID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse topic id", http.StatusBadRequest))
		return
	}

	topicAssignment.BatchID, err = parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to session batch id", http.StatusBadRequest))
		return
	}

	topicAssignment.ID, err = parser.GetUUID("topicAssignmentID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to session assignment id", http.StatusBadRequest))
		return
	}

	topicAssignment.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.service.DeleteTopicAssignment(&topicAssignment)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Assignment successfully deleted for batch topic")
}

// GetSessionProgrammingAssignment will return all the batch session programming assignments.
func (controller *BatchTopicAssignmentController) GetAllTopicAssignments(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get All Topic Assignments Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batchID id", http.StatusBadRequest))
		return
	}
	var assignments []batch.TopicAssignmentDTO

	err = controller.service.GetAllTopicAssignments(tenantID, batchID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)

}

// GetSessionProgrammingAssignment will return all the batch session programming assignments.
func (controller *BatchTopicAssignmentController) GetSessionProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Programming Assignment Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch session ID.
	batchSessionID, err := parser.GetUUID("batchSessionID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch session id", http.StatusBadRequest))
		return
	}

	var assignments []batch.TopicAssignmentDTO
	err = controller.service.GetSessionProgrammingAssignment(tenantID, batchSessionID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}

// GetLatestBatchProgrammingAssignment will fetch all the programming assignments for the given batch session.
func (controller *BatchTopicAssignmentController) GetLatestBatchProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Latest Batch Programming Assignment Called==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set batch ID.
	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	var assignments []batch.TopicAssignmentDTO
	err = controller.service.GetLatestBatchProgrammingAssignment(tenantID, batchID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}

// GetBatchProgrammingAssignment will return all the batch session programming assignments.
func (controller *BatchTopicAssignmentController) GetBatchProgrammingAssignment(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Get Batch Programming Assignment Called==============================")
	// parser := web.NewParser(r)

	// tenantID, err := parser.GetTenantID()
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
	// 	return
	// }

	// // Parse and set batch session ID.
	// batchID, err := parser.GetUUID("batchID")
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
	// 	return
	// }

	// var assignments []batch.BatchTopicAssignmentDTO
	// err = controller.service.GetBatchProgrammingAssignment(tenantID, batchID, &assignments, parser)
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }
	// // Writing Response with OK Status to ResponseWriter.
	// web.RespondJSON(w, http.StatusOK, assignments)
}

// GetSessionAssignmentList will return session_assignment list.
func (controller *BatchTopicAssignmentController) GetSessionAssignmentList(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("============================== GetSessionAssignmentList Call ==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	var assignments []batch.TopicAssignmentDTO
	err = controller.service.GetSessionAssignmentList(tenantID, batchID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}

// GetAllAssignmentsWithSubmissions will return session_assignment list.
func (controller *BatchTopicAssignmentController) GetAllAssignmentsWithSubmissions(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("============================== GetAllAssignmentsWithSubmissions Call ==============================")
	parser := web.NewParser(r)

	tenantID, err := parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	batchID, err := parser.GetUUID("batchID")
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse batch id", http.StatusBadRequest))
		return
	}

	var assignments []batch.TopicAssignmentDTO
	err = controller.service.GetAllAssignmentsWithSubmissions(tenantID, batchID, &assignments, parser)
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, assignments)
}
