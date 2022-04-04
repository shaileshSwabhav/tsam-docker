package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/programming/service"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProgrammingQuestionTypeController provides method to update, delete, add, get all for ProgrammingQuestionType.
type ProgrammingQuestionTypeController struct {
	ProgrammingQuestionTypeService *service.ProgrammingQuestionTypeService
}

// NewProgrammingQuestionTypeController creates new instance of ProgrammingQuestionTypeController.
func NewProgrammingQuestionTypeController(programmingQuestionTypeService *service.ProgrammingQuestionTypeService) *ProgrammingQuestionTypeController {
	return &ProgrammingQuestionTypeController{
		ProgrammingQuestionTypeService: programmingQuestionTypeService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *ProgrammingQuestionTypeController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one question type.
	router.HandleFunc("/tenant/{tenantID}/programming-question-type/credential/{credentialID}",
		controller.AddProgrammingQuestionType).Methods(http.MethodPost)

	// Add multiple question types.
	router.HandleFunc("/tenant/{tenantID}/programming-question-types/credential/{credentialID}",
		controller.AddProgrammingQuestionTypes).Methods(http.MethodPost)

	// Update one question type.
	router.HandleFunc("/tenant/{tenantID}/programming-question-type/{programmingTypeID}/credential/{credentialID}",
		controller.UpdateProgrammingQuestionType).Methods(http.MethodPut)

	// Delete one question type.
	router.HandleFunc("/tenant/{tenantID}/programming-question-type/{programmingTypeID}/credential/{credentialID}",
		controller.DeleteProgrammingQuestionType).Methods(http.MethodDelete)

	// Get question types with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/programming-question-type/limit/{limit}/offset/{offset}",
		controller.GetProgrammingQuestionTypes).Methods(http.MethodGet)

	// Get question type list.
	router.HandleFunc("/tenant/{tenantID}/programming-question-type",
		controller.GeProgrammingQuestionTypeList).Methods(http.MethodGet)

	log.NewLogger().Info("Programming Question Type Routes Registered")
}

// AddProgrammingQuestionType will add programming question type to the table.
func (controller *ProgrammingQuestionTypeController) AddProgrammingQuestionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingQuestionType called==============================")
	param := mux.Vars(r)
	programmingType := new(programming.ProgrammingQuestionType)

	err := web.UnmarshalJSON(r, programmingType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	programmingType.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	programmingType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = programmingType.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionTypeService.AddProgrammingQuestionType(programmingType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question type added successfully")
}

// AddProgrammingQuestionTypes will add multiple programming question types to the table.
func (controller *ProgrammingQuestionTypeController) AddProgrammingQuestionTypes(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddProgrammingQuestionTypes called==============================")
	param := mux.Vars(r)
	programmingTypes := new([]programming.ProgrammingQuestionType)

	err := web.UnmarshalJSON(r, programmingTypes)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	for _, question := range *programmingTypes {
		err = question.Validate()
		if err != nil {
			log.NewLogger().Error(err.Error())
			web.RespondError(w, err)
			return
		}
	}

	err = controller.ProgrammingQuestionTypeService.AddProgrammingQuestionTypes(programmingTypes, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question types added successfully")
}

// UpdateProgrammingQuestionType will update specified programming question type.
func (controller *ProgrammingQuestionTypeController) UpdateProgrammingQuestionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateProgrammingQuestionType called==============================")
	param := mux.Vars(r)
	programmingType := new(programming.ProgrammingQuestionType)

	err := web.UnmarshalJSON(r, programmingType)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	programmingType.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	programmingType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	programmingType.ID, err = util.ParseUUID(param["programmingTypeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = programmingType.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionTypeService.UpdateProgrammingQuestionType(programmingType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question updated successfully")
}

// DeleteProgrammingQuestionType will update specified programming question type.
func (controller *ProgrammingQuestionTypeController) DeleteProgrammingQuestionType(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteProgrammingQuestionType called==============================")
	param := mux.Vars(r)
	var err error
	programmingType := new(programming.ProgrammingQuestionType)

	programmingType.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	programmingType.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	programmingType.ID, err = util.ParseUUID(param["programmingTypeID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionTypeService.DeleteProgrammingQuestionType(programmingType)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Programming question deleted successfully")
}

// GetProgrammingQuestionTypes will return all the feedback question type with limit and offset.
func (controller *ProgrammingQuestionTypeController) GetProgrammingQuestionTypes(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetProgrammingQuestionTypes called==============================")
	param := mux.Vars(r)
	programmingTypes := new([]programming.ProgrammingQuestionTypeDTO)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Parse form
	r.ParseForm()

	err = controller.ProgrammingQuestionTypeService.GetProgrammingQuestionTypes(programmingTypes, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, programmingTypes)
}

// GeProgrammingQuestionTypeList will return all the feedback question types.
func (controller *ProgrammingQuestionTypeController) GeProgrammingQuestionTypeList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllProgrammingQGeProgrammingQuestionTypeListuestionTypes called==============================")
	param := mux.Vars(r)
	programmingTypes := new([]programming.ProgrammingQuestionTypeDTO)

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.ProgrammingQuestionTypeService.GeProgrammingQuestionTypeList(programmingTypes, r.Form, tenantID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, programmingTypes)
}
