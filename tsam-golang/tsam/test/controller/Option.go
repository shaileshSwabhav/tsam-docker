package controller

import (
	"net/http"

	"github.com/techlabs/swabhav/tsam/repository"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	tst "github.com/techlabs/swabhav/tsam/models/test"
	"github.com/techlabs/swabhav/tsam/test/service"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/web"
)

// OptionController Provide method to Update, Delete, Add, Get Method For Option.
type OptionController struct {
	OptionService *service.OptionService
}

// NewOptionController Create New Instance Of OptionController.
func NewOptionController(ser *service.OptionService) *OptionController {
	return &OptionController{
		OptionService: ser,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (optionCon *OptionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/option", optionCon.GetOptions).Methods(http.MethodGet)
	router.HandleFunc("/option", optionCon.AddOption).Methods(http.MethodPost)
	router.HandleFunc("/option/mutiple", optionCon.AddOptions).Methods(http.MethodPost)
	router.HandleFunc("/option/{optionid}", optionCon.GetOption).Methods(http.MethodGet)
	router.HandleFunc("/option/{optionid}", optionCon.UpdateOption).Methods(http.MethodPut)
	router.HandleFunc("/option/{optionid}", optionCon.DeleteOption).Methods(http.MethodDelete)
	router.HandleFunc("/option/question/{questionid}", optionCon.GetOptionsByQuestionID).
		Methods(http.MethodGet)
	log.NewLogger().Info("Option Route Registered")
}

// GetOptions Return All option
func (optionCon *OptionController) GetOptions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get All Option API Call")
	options := &[]tst.Option{}
	err := optionCon.OptionService.GetOptions(options)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, options)
}

// GetOptionsByQuestionID Return All option By Question ID
func (optionCon *OptionController) GetOptionsByQuestionID(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get All Option By QUestion ID API Call")
	questionID, err := util.ParseUUID(mux.Vars(r)["questionid"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse question id", http.StatusBadRequest))
		return
	}
	options := &[]tst.Option{}
	err = optionCon.OptionService.GetOptions(options, repository.Filter("question_id =?", questionID))
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, options)
}

// GetOption Return Specific Option
func (optionCon *OptionController) GetOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get Option API Call")

	// Parse Option ID
	param := mux.Vars(r)
	optionID, err := util.ParseUUID(param["optionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable To Parse Option ID", http.StatusBadRequest))
		return
	}

	option := &tst.Option{}
	err = optionCon.OptionService.GetOption(option, &optionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, option)
}

// UpdateOption Update The Option
func (optionCon *OptionController) UpdateOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Update Option API Call")
	option := &tst.Option{}

	// Parse Data From Request
	param := mux.Vars(r)
	optionID, err := util.ParseUUID(param["optionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Option ID", http.StatusBadRequest))
		return
	}

	err = web.UnmarshalJSON(r, &option)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	option.ID = optionID

	err = optionCon.OptionService.UpdateOption(option)
	if err != nil {
		web.RespondError(w, err)
		return
	}
}

// AddOption Add New Option
func (optionCon *OptionController) AddOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add Option API Call")

	// Parse Option From Request & Add New ID.
	option := tst.Option{}
	err := web.UnmarshalJSON(r, &option)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	option.ID = util.GenerateUUID()

	// Add Option To Database
	err = optionCon.OptionService.AddOption(&option)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, option.ID)
}

// AddOptions Add Multiple New Option
func (optionCon *OptionController) AddOptions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Add Option API Call")
	optionsIDs := []uuid.UUID{}
	options := []tst.Option{}

	// Parse Options From Request
	err := web.UnmarshalJSON(r, &options)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}

	err = optionCon.OptionService.AddOptions(&options, &optionsIDs)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, optionsIDs)
}

// DeleteOption Delete Option
func (optionCon *OptionController) DeleteOption(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Delete Option API Call")
	param := mux.Vars(r)
	optionID, err := util.ParseUUID(param["optionid"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse Option ID", http.StatusBadRequest))
		return
	}
	err = optionCon.OptionService.DeleteOption(&optionID)
	if err != nil {
		web.RespondError(w, err)
		return
	}
}
