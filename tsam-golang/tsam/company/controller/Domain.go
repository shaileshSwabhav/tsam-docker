package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/company/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	model "github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DomainController Provide method to Update, Delete, Add, Get Method For Domain.
type DomainController struct {
	DomainService *services.DomainService
}

// NewDomainController Create New Instance Of DomainController.
func NewDomainController(ser *services.DomainService) *DomainController {
	return &DomainController{
		DomainService: ser,
	}
}

// RegisterRoutes Register All Endpoint To Router.
func (controller *DomainController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	router.HandleFunc("/tenant/{tenantID}/domain/credential/{credentialID}", controller.AddDomain).Methods(http.MethodPost)
	router.HandleFunc("/tenant/{tenantID}/domain/{domainID}/credential/{credentialID}", controller.UpdateDomain).Methods(http.MethodPut)
	router.HandleFunc("/tenant/{tenantID}/domain/{domainID}/credential/{credentialID}", controller.DeleteDomain).Methods(http.MethodDelete)
	router.HandleFunc("/tenant/{tenantID}/domain/limit/{limit}/offset/{offset}", controller.GetAllDomains).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/domain/list", controller.GetDomainList).Methods(http.MethodGet)
	router.HandleFunc("/tenant/{tenantID}/domain/{domainID}", controller.GetDomain).Methods(http.MethodGet)
	log.NewLogger().Info("All routes registered..")
}

// AddDomain Add New Domain
func (controller *DomainController) AddDomain(w http.ResponseWriter, r *http.Request) {

	log.NewLogger().Info("Add Domain API Call")
	domain := &model.Domain{}

	// Fill the domain variable with given data.
	err := web.UnmarshalJSON(r, domain)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	domain.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credential ID
	domain.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Validate domain
	err = domain.ValidateDomain()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// call add service
	err = controller.DomainService.AddDomain(domain)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	//Writing response with ok status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, domain.ID)
}

//UpdateDomain Update the domain
func (controller *DomainController) UpdateDomain(w http.ResponseWriter, r *http.Request) {

	log.NewLogger().Info("Update Domain API Call")
	domain := &model.Domain{}

	// Parse domain from request.
	err := web.UnmarshalJSON(r, domain)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Parse and set tenant ID.
	domain.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credential ID
	domain.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set domain ID
	domain.ID, err = util.ParseUUID(mux.Vars(r)["domainID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse domain id", http.StatusBadRequest))
		return
	}

	// Validate domain
	err = domain.ValidateDomain()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// call update service
	err = controller.DomainService.UpdateDomain(domain)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	//Writing response with ok status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Domain updated")
}

//DeleteDomain delete Domain
func (controller *DomainController) DeleteDomain(w http.ResponseWriter, r *http.Request) {

	log.NewLogger().Info("Delete Domain API Call")
	domain := &model.Domain{}
	var err error

	// Parse and set tenant ID.
	domain.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credential ID
	domain.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Parse and set domain ID
	domain.ID, err = util.ParseUUID(mux.Vars(r)["domainID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse domain id", http.StatusBadRequest))
		return
	}

	// call delete service
	err = controller.DomainService.DeleteDomain(domain)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	//Writing response with ok status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Domain deleted")
}

// GetAllDomains returns all Domain
func (controller *DomainController) GetAllDomains(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get All Domain API Call")

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Fill the r.Form
	r.ParseForm()

	// limit,offset & totalCount for pagination
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	domains := &[]*model.Domain{}
	err = controller.DomainService.GetAllDomains(tenantID, domains, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, domains)
}

// GetAllDomains returns all Domain
func (controller *DomainController) GetDomainList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Get All Domain API Call")

	// Parse and set tenant ID.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	domains := &[]*model.Domain{}
	err = controller.DomainService.GetDomainList(tenantID, domains)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, domains)
}

// GetDomain returns specific Domain
func (controller *DomainController) GetDomain(w http.ResponseWriter, r *http.Request) {

	log.NewLogger().Info("Get Domain API Call")

	var err error
	domain := &model.Domain{}

	// Parse and set tenant ID.
	domain.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set domain ID
	domain.ID, err = util.ParseUUID(mux.Vars(r)["domainID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse domain id", http.StatusBadRequest))
		return
	}

	err = controller.DomainService.GetDomain(domain)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status
	web.RespondJSON(w, http.StatusOK, domain)
}
