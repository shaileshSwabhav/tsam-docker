package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/administration/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SalesPersonController provides methods to update, delete, add, get, get all salesperson on SalesPerson.
type SalesPersonController struct {
	SalesPersonService *service.SalesPersonService
}

// NewSalesPersonController creates new instance of SalesPersonController.
func NewSalesPersonController(salesPersonService *service.SalesPersonService) *SalesPersonController {
	return &SalesPersonController{
		SalesPersonService: salesPersonService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (con *SalesPersonController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {
	// salesperson listing
	router.HandleFunc("/tenant/{tenantID}/salesperson/list", con.GetSalesPeopleList).Methods(http.MethodGet)

	// add one salesperson
	router.HandleFunc("/tenant/{tenantID}/salesperson/credential/{credentialID}", con.AddSalesPerson).Methods(http.MethodPost)

	//get all salespeople
	router.HandleFunc("/tenant/{tenantID}/salesperson", con.GetAllSalesPeople).Methods(http.MethodGet)

	//get one sales person
	router.HandleFunc("/tenant/{tenantID}/salesperson/{salespersonID}", con.GetSalesPerson).Methods(http.MethodGet)

	//update one salesperson
	router.HandleFunc("/tenant/{tenantID}/salesperson/{salespersonID}/credential/{credentialID}", con.UpdateSalesPerson).Methods(http.MethodPut)

	// delete one salesperson
	router.HandleFunc("/tenant/{tenantID}/salesperson/{salespersonID}/credential/{credentialID}", con.DeleteSalesPerson).Methods(http.MethodDelete)

	log.NewLogger().Info("Sales Person Routes Registered")
}

// GetAllSalesPeople returns all sales people in database as response
func (con *SalesPersonController) GetAllSalesPeople(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetAllSalesPeople API Called")

	//create bucket
	salesPeople := []general.User{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//call get all service method
	if err := con.SalesPersonService.GetAllSalesPeople(tenantID, &salesPeople); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, salesPeople)
}

// GetSalesPeopleList returns all sales people in database as response
func (con *SalesPersonController) GetSalesPeopleList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetSalesPeopleList API Called")

	//create bucket
	salesPeople := []list.User{}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//call get all service method
	if err := con.SalesPersonService.GetSalesPeopleList(tenantID, &salesPeople); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, salesPeople)
}

// GetSalesPerson returns a specific sales person as response
func (con *SalesPersonController) GetSalesPerson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("GetSalesPerson API Called")

	//create bucket
	salesPerson := general.User{}

	//getting id from param and parsing it to uuid
	salesPersonID, err := util.ParseUUID(mux.Vars(r)["salespersonID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse sales person id", http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//call get service method
	if err := con.SalesPersonService.GetSalesPerson(salesPersonID, &salesPerson, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, salesPerson)
}

// UpdateSalesPerson updates the specific sales person
func (con *SalesPersonController) UpdateSalesPerson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Update SalesPerson API Called")

	//create bucket
	salesPerson := general.User{}

	//getting id from param and parsing it to uuid
	salesPersonID, err := util.ParseUUID(mux.Vars(r)["salespersonID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse sales person id", http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//unmarshal json
	err = web.UnmarshalJSON(r, &salesPerson)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// //to get only date for date of birth
	// if salesPerson.DateOfBirth != nil {
	// 	if len(*salesPerson.DateOfBirth) > 10 {
	// 		tempDateOfBirth := *salesPerson.DateOfBirth
	// 		tempDateOfBirth = tempDateOfBirth[:10]
	// 		salesPerson.DateOfBirth = &tempDateOfBirth
	// 	}
	// }

	// //to get only date for date of joining
	// if salesPerson.DateOfJoining != nil {
	// 	if len(*salesPerson.DateOfJoining) > 10 {
	// 		tempDateOfJoining := *salesPerson.DateOfJoining
	// 		tempDateOfJoining = tempDateOfJoining[:10]
	// 		salesPerson.DateOfJoining = &tempDateOfJoining
	// 	}
	// }

	//call update service method
	if err := con.SalesPersonService.UpdateSalesPerson(&salesPerson, salesPersonID, tenantID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Sales Person Updated")
}

// AddSalesPerson adds a new sales person
func (con *SalesPersonController) AddSalesPerson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("AddSalesPerson API Called")

	//create bucket
	salesPerson := general.User{}

	//unmarshal json
	if err := web.UnmarshalJSON(r, &salesPerson); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//call add service method
	if err := con.SalesPersonService.AddSalesPerson(&salesPerson, credentialID, tenantID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, salesPerson.ID)
}

// DeleteSalesPerson deletes the specific Sales Person
func (con *SalesPersonController) DeleteSalesPerson(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("Delete SalesPerson API Called")

	//create bucket
	salesPerson := general.User{}

	//getting id from param and parsing it to uuid
	salesPersonID, err := util.ParseUUID(mux.Vars(r)["salespersonID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse sales person id", http.StatusBadRequest))
		return
	}

	//getting tenant id from param and parsing it to uuid
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	//getting credential id from param and parsing it to uuid
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	//call delete service method
	if err := con.SalesPersonService.DeleteSalesPerson(&salesPerson, tenantID, salesPersonID, credentialID); err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Sales Person Deleted")
}
