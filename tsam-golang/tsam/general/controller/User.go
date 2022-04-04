package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// UserController gives method for listing for all users that is salesperson and admin.
type UserController struct {
	UserService *service.UserService
}

// NewUserController creates new instance  UserController.
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (controller *UserController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Add one user.
	router.HandleFunc("/tenant/{tenantID}/user/credential/{credentialID}", controller.AddUser).
		Methods(http.MethodPost)

	// Update user.
	router.HandleFunc("/tenant/{tenantID}/user/{userID}/credential/{credentialID}", controller.UpdateUser).
		Methods(http.MethodPut)

	// Delete user.
	router.HandleFunc("/tenant/{tenantID}/user/{userID}/credential/{credentialID}", controller.DeleteUser).
		Methods(http.MethodDelete)

	// Get specific type of user.
	router.HandleFunc("/tenant/{tenantID}/user/limit/{limit}/offset/{offset}", controller.GetSpecificUsers).
		Methods(http.MethodGet)

	// Get list of all users.
	router.HandleFunc("/tenant/{tenantID}/user-list", controller.GetUserList).
		Methods(http.MethodGet)

	// Get list of salespeople.
	router.HandleFunc("/tenant/{tenantID}/salesperson-list", controller.GetSalesPeopleList).
		Methods(http.MethodGet)

	// Get list of user credentials.
	router.HandleFunc("/tenant/{tenantID}/user/credential-list", controller.GetUserCredentialList).
		Methods(http.MethodGet)

	log.NewLogger().Info("==========User routes registered=============")
}

// GetSpecificUsers gets all users by search criteria from database.
func (controller *UserController) GetSpecificUsers(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetSpecificUsers called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form.
	r.ParseForm()

	// Limit, offset & totalCount for pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	users := &[]general.User{}
	err = controller.UserService.GetSpecificUsers(tenantID, users, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, users)
}

// GetUserList returns list of all salesperson and admin.
func (controller *UserController) GetUserList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetListOfAllUsers called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	users := []list.User{}

	err = controller.UserService.GetUserList(tenantID, &users)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, users)
}

// GetSalesPeopleList returns list of all salespeople.
func (controller *UserController) GetSalesPeopleList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetListOfAllSalesPeople called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	salesPeople := []list.User{}

	err = controller.UserService.GetSalesPeopleList(tenantID, &salesPeople)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, salesPeople)
}

// GetUserCredentialList returns list of all users' credentials.
func (controller *UserController) GetUserCredentialList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetListOfAllUserCredentials called==============================")

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	credentials := []list.Credential{}

	err = controller.UserService.GetUserCredentialList(tenantID, &credentials)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing response with OK status.
	web.RespondJSON(w, http.StatusOK, credentials)
}

// AddUser adds new User to the table.
func (controller *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddUser called==============================")
	param := mux.Vars(r)
	user := general.User{}

	err := web.UnmarshalJSON(r, &user)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.CreatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = user.Validate()
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.UserService.AddUser(&user)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "User added successfully")
}

// UpdateUser updates the specific user by id.
func (controller *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateUser called==============================")
	param := mux.Vars(r)
	user := general.User{}

	err := web.UnmarshalJSON(r, &user)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.UpdatedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.ID, err = util.ParseUUID(param["userID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = user.Validate()
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.UserService.UpdateUser(&user)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "User updated successfully")
}

// DeleteUser deletes specific user by id.
func (controller *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteUser called==============================")
	param := mux.Vars(r)
	user := general.User{}
	var err error

	user.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	user.ID, err = util.ParseUUID(param["userID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = controller.UserService.DeleteUser(&user)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "User deleted successfully")
}
