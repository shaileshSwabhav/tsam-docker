package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/web"
)

// ModuleProgrammingConceptController provides methods to do CRUD operations.
type ModuleProgrammingConceptController struct {
	ModuleProgrammingConceptService *service.ModuleProgrammingConceptService
	log                             log.Logger
	auth                            *security.Authentication
}

// NewModuleProgrammingConceptController creates new instance of concept module controller.
func NewModuleProgrammingConceptController(moduleProgrammingConceptService *service.ModuleProgrammingConceptService,
	log log.Logger, auth *security.Authentication) *ModuleProgrammingConceptController {
	return &ModuleProgrammingConceptController{
		ModuleProgrammingConceptService: moduleProgrammingConceptService,
		log:                             log,
		auth:                            auth,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ModuleProgrammingConceptController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add programming concept module.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module",
		controller.AddModuleProgrammingConcepts).Methods(http.MethodPost)

	// Update programming concept module.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module/module/{moduleID}",
		controller.UpdateModuleProgrammingConcepts).Methods(http.MethodPut)

	// Delete programming concept module.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module/{moduleProgrammingConceptID}",
		controller.DeleteModuleProgrammingConcept).Methods(http.MethodDelete)

	// Get all programming concepts modules for concept tree.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module-tree",
		controller.GetAllModuleProgrammingConceptsForConceptTree).Methods(http.MethodGet)

	// Get all programming concepts modules.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module",
		controller.GetAllModuleProgrammingConcepts).Methods(http.MethodGet)

	// Get all programming concepts modules for talent score.
	router.HandleFunc("/tenant/{tenantID}/talent/{talentID}/programming-concept-module",
		controller.GetAllModuleProgrammingConceptsForTalentScore).Methods(http.MethodGet)

	// Get all programming concepts modules for an assignment.
	router.HandleFunc("/tenant/{tenantID}/batch-topic-assignment/{assignmentID}/programming-concepts-modules",
		controller.GetModuleConceptsForAssignment).Methods(http.MethodGet)

	// Get one programming concept module.
	router.HandleFunc("/tenant/{tenantID}/programming-concept-module/{moduleProgrammingConceptID}",
		controller.GetModuleProgrammingConcept).Methods(http.MethodGet)

	log.NewLogger().Info("Programming Concept Route Registered")
}

// AddModuleProgrammingConcepts will add new programming concept modules to the table.
func (controller *ModuleProgrammingConceptController) AddModuleProgrammingConcepts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddModuleProgrammingConcepts called==============================")

	parser := web.NewParser(r)

	// Create bucket for pragramming conceptModules.
	conceptModules := []programming.ModuleProgrammingConcepts{}

	// Unmarshal JSON.
	err := web.UnmarshalJSON(r, &conceptModules)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate compulsary concept modules fields.
	for _, conceptModule := range conceptModules {
		err := conceptModule.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set credential id.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call add service method.
	err = controller.ModuleProgrammingConceptService.AddModuleProgrammingConcepts(tenantID, credentialID, &conceptModules)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept modules added successfully")
}

// UpdateModuleProgrammingConcepts will update the specified programming concept module record in the table.
func (controller *ModuleProgrammingConceptController) UpdateModuleProgrammingConcepts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateModuleProgrammingConcepts called==============================")

	parser := web.NewParser(r)

	// Create bucket for programming concept module.
	conceptModules := []programming.ModuleProgrammingConcepts{}

	// Unmarshal JSON.
	err := web.UnmarshalJSON(r, &conceptModules)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Validate compulsary concept modules fields.
	for _, conceptModule := range conceptModules {
		err := conceptModule.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set credential id.
	credentialID, err := controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set module id.
	moduleID, err := parser.GetUUID("moduleID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.ModuleProgrammingConceptService.UpdateModuleProgrammingConcepts(moduleID, tenantID, credentialID,
		&conceptModules)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept modules updated successfully")
}

// DeleteModuleProgrammingConcept will delete the specified programming concept module record from the table.
func (controller *ModuleProgrammingConceptController) DeleteModuleProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteModuleProgrammingConcept called==============================")

	parser := web.NewParser(r)

	// Create bucket for pragramming concept.
	conceptModule := programming.ModuleProgrammingConcepts{}

	// Create error varoable.
	var err error

	// Parse and set tenant id.
	conceptModule.TenantID, err = parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set credential id.
	conceptModule.DeletedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Prade and set concept module id.
	conceptModule.ID, err = parser.GetUUID("moduleProgrammingConceptID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call delete service method.
	err = controller.ModuleProgrammingConceptService.DeleteModuleProgrammingConcept(&conceptModule)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Programming concept module deleted successfully")
}

// GetAllModuleProgrammingConceptsForConceptTree will return all the records from programming concept modules table for concept tree.
func (controller *ModuleProgrammingConceptController) GetAllModuleProgrammingConceptsForConceptTree(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllModuleProgrammingConceptsForConceptTree called==============================")

	parser := web.NewParser(r)

	// Create bucket for pragramming concept.
	conceptModules := []programming.ModuleProgrammingConcepts{}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Create variables for total.
	var totalCount int

	// Call get all service method.
	err = controller.ModuleProgrammingConceptService.GetAllModuleProgrammingConceptsForConceptTree(&conceptModules, tenantID, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, conceptModules)
}

// GetAllModuleProgrammingConcepts will return all the records from programming concept modules table.
func (controller *ModuleProgrammingConceptController) GetAllModuleProgrammingConcepts(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllModuleProgrammingConcepts called==============================")

	parser := web.NewParser(r)

	// Create bucket for pragramming concept.
	conceptModules := []programming.ModuleProgrammingConceptsDTO{}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Create variables for total.
	var totalCount int

	// Call get all service method.
	err = controller.ModuleProgrammingConceptService.GetAllModuleProgrammingConcepts(&conceptModules, tenantID, parser, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, conceptModules)
}

// GetModuleProgrammingConcept will return specified record from programming concept module table.
func (controller *ModuleProgrammingConceptController) GetModuleProgrammingConcept(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetModuleProgrammingConcept called==============================")

	parser := web.NewParser(r)

	// Create error varoable.
	var err error

	// Create bucket for pragramming concept module.
	conceptModule := programming.ModuleProgrammingConcepts{}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set concept module id.
	conceptModule.ID, err = parser.GetUUID("moduleProgrammingConceptID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call get one service method.
	err = controller.ModuleProgrammingConceptService.GetModuleProgrammingConcept(&conceptModule, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, conceptModule)
}

// GetModuleConceptsForAssignment will concept_modules for an assignment.
func (controller *ModuleProgrammingConceptController) GetModuleConceptsForAssignment(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetModuleConceptsForAssignment called==============================")

	parser := web.NewParser(r)

	var err error

	// Create bucket for pragramming concept module.
	conceptModule := []programming.ModuleProgrammingConceptsDTO{}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set concept module id.
	assignmentID, err := parser.GetUUID("assignmentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call get one service method.
	err = controller.ModuleProgrammingConceptService.GetModuleConceptsForAssignment(&conceptModule,
		tenantID, assignmentID, parser)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, conceptModule)
}

// GetAllModuleProgrammingConceptsForTalentScore will return all the records from programming concept modules table for talent score.
func (controller *ModuleProgrammingConceptController) GetAllModuleProgrammingConceptsForTalentScore(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllModuleProgrammingConceptsForTalentScore called==============================")

	parser := web.NewParser(r)

	// Create bucket for pragramming concept.
	conceptModules := []programming.ModuleProgrammingConceptsForTalentScore{}

	// Parse and set tenant id.
	tenantID, err := parser.GetTenantID()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set talent id.
	talentID, err := parser.GetUUID("talentID")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call get all service method.
	err = controller.ModuleProgrammingConceptService.GetAllModuleProgrammingConceptsForTalentScore(&conceptModules, tenantID, talentID, parser)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, conceptModules)
}
