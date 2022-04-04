package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// LoginController provides methods to do Get(Login) operations on Login.
type LoginController struct {
	LoginService *service.CredentialService
	Auth         *security.Authentication
}

// NewCredentialController Create New Instance Of LoginController.
func NewCredentialController(service *service.CredentialService, auth *security.Authentication) *LoginController {
	return &LoginController{
		LoginService: service,
		Auth:         auth,
	}
}

// RegisterRoutes will register the login routes.
func (controller *LoginController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Login.
	login := router.HandleFunc("/tenant/{tenantID}/login",
		controller.Login).Methods(http.MethodPost)

	// Update password.
	updatePassword := router.HandleFunc("/tenant/{tenantID}/login/password",
		controller.UpdateCredentialPassword).Methods(http.MethodPut)

	// Verify old password.
	router.HandleFunc("/tenant/{tenantID}/login/verify-password",
		controller.PasswordVerification).Methods(http.MethodPost)

	// Change password after login.
	router.HandleFunc("/tenant/{tenantID}/login/change-password/credential/{credentialID}",
		controller.ChangePassword).Methods(http.MethodPut)

	// Logout.
	router.HandleFunc("/tenant/{tenantID}/logout/credential/{credentialID}",
		controller.Logout).Methods(http.MethodPost)

	// Validates session when the page is being redirected from login page.
	router.HandleFunc("/tenant/{tenantID}/session",
		controller.CheckSessionValidity).Methods(http.MethodGet)

	// Use for testing purpose only.
	addLogin := router.HandleFunc("/tenant/{tenantID}/addLogin/credential/{credentialID}",
		controller.AddLogin).Methods(http.MethodPost)

	// Delete credential/
	router.HandleFunc("/tenant/{tenantID}/delete/credential/{credentialID}",
		controller.DeleteCredential).Methods(http.MethodDelete)

	// Get credentials by role name.
	router.HandleFunc("/tenant/{tenantID}/credential-by-role",
		controller.GetAllCredentialsByRole).Methods(http.MethodGet)

	*exclude = append(*exclude, login, addLogin, updatePassword)

	log.NewLogger().Info("Login Routes Registered")
}

// Login will login the user.
func (controller *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================Login called==============================")
	credential := general.Credential{}
	tokenDTO := general.TokenDTO{}

	err := web.UnmarshalJSON(r, &credential)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse requested data", http.StatusBadRequest))
		return
	}

	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	err = controller.LoginService.Login(&tokenDTO, &credential, tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	loginID := controller.getLoginID(&credential)
	// Creating login token.
	token, err := controller.Auth.GenerateToken(credential.ID.String(), loginID.String(), credential.Email, credential.RoleID.String())
	if err != nil {
		log.NewLogger().Error(err.Error())
		errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		return
	}

	tokenDTO.LoginID = loginID
	tokenDTO.Token = token

	// tokenDTO := general.TokenDTO{
	// 	CredentialID: credential.ID,
	// 	FirstName:    credential.FirstName,
	// 	Email:        credential.Email,
	// 	TenantID:     credential.TenantID,
	// 	LoginID:      controller.getLoginID(&credential),
	// 	RoleID:       credential.RoleID,
	// 	DepartmentID: credential.DepartmentID,
	// 	Token:        token,
	// }
	// if credential.LastName != nil {
	// 	tokenDTO.LastName = *credential.LastName
	// }
	// if credential.DepartmentID != nil && util.IsUUIDValid(*credential.DepartmentID) {
	// 	tokenDTO.DepartmentID = credential.DepartmentID
	// }

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, tokenDTO)
}

// Logout will end user's session & logout.
func (controller *LoginController) Logout(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================Logout called==============================")
	var err error
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant ID", http.StatusBadRequest))
		return
	}
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse login ID", http.StatusBadRequest))
		return
	}
	credential := general.Credential{}
	// credential.ID = credentialID
	// credential.TenantID = tenantID

	r.ParseForm()

	err = controller.LoginService.Logout(&credential, tenantID, credentialID, r.Form)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Log out successful")
}

// AddLogin used for testing purpose.
func (controller *LoginController) AddLogin(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================AddLogin called==============================")
	var err error
	credential := general.Credential{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	err = web.UnmarshalJSON(r, &credential)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	credential.TenantID = tenantID
	credential.CreatedBy = credentialID
	uow := repository.NewUnitOfWork(controller.LoginService.DB, false)
	err = controller.LoginService.AddCredential(&credential, uow)
	if err != nil {
		uow.RollBack()
		web.RespondError(w, err)
		return
	}
	uow.Commit()

	web.RespondJSON(w, http.StatusOK, credential)
}

// UpdateCredentialPassword updates the password of the credential.
func (controller *LoginController) UpdateCredentialPassword(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateCredentialPassword called==============================")
	var err error
	credential := general.Credential{}
	err = web.UnmarshalJSON(r, &credential)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}
	credential.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	if util.IsEmpty(credential.Email) || !util.ValidateEmail(credential.Email) {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Email must be specified and should be of type example@domain.com", http.StatusBadRequest))
		return
	}
	if util.IsEmpty(credential.Password) {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Password must be specified", http.StatusBadRequest))
		return
	}
	err = controller.LoginService.UpdateCredentialPassword(&credential)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Password updated successfully")

}

// PasswordVerification verifies the password.
func (controller *LoginController) PasswordVerification(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================PasswordVerification called==============================")

	// Create error variable.
	var err error

	// Create bucket.
	passwordChange := general.PasswordChange{}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &passwordChange)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := passwordChange.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	passwordChange.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call password verification service method.
	err = controller.LoginService.PasswordVerification(&passwordChange)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Password is correct")
}

// ChangePassword updates the password of the credential after login.
func (controller *LoginController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================ChangePassword called==============================")

	// Create error variable.
	var err error

	// Create bucket.
	passwordChange := general.PasswordChange{}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &passwordChange)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := passwordChange.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	passwordChange.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting Credential id from param and parsing it to uuid.
	passwordChange.CredentialID, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call update password with verification service method.
	err = controller.LoginService.ChangePassword(&passwordChange)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Password updated successfully")
}

// DeleteCredential deletes credential record.
// For testing purpose only.
func (controller *LoginController) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================DeleteCredential called==============================")
	var err error
	credential := general.Credential{}
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}
	credentialID, err := util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}
	userID, err := uuid.FromString("0722cc41-2aa7-46b9-a12b-79b6e5592d2c")
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	credential.UserID = &userID
	uow := repository.NewUnitOfWork(controller.LoginService.DB, false)
	err = controller.LoginService.DeleteCredential(&credential, tenantID, credentialID, userID, "user_id=?", uow)
	if err != nil {
		uow.RollBack()
		web.RespondError(w, err)
		return
	}
	uow.Commit()
	web.RespondJSON(w, http.StatusOK, credential)
}

// CheckSessionValidity checks the currently logged in credentials session expiry.
func (controller *LoginController) CheckSessionValidity(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================CheckSessionValidity called==============================")

	err := controller.Auth.CheckSessionValidity(r)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, true)
}

func (controller *LoginController) getLoginID(login *general.Credential) uuid.UUID {

	switch {
	case login.CollegeID != nil:
		return *login.CollegeID
	case login.CompanyID != nil:
		return *login.CompanyID
	case login.TalentID != nil:
		return *login.TalentID
	case login.FacultyID != nil:
		return *login.FacultyID
	case login.UserID != nil:
		return *login.UserID
	case login.SalesPersonID != nil:
		return *login.SalesPersonID
	case login.EmployeeID != nil:
		return *login.EmployeeID
	}
	return uuid.Nil
}

// GetAllCredentialsByRole will return all the records from credential table by role names.
func (controller *LoginController) GetAllCredentialsByRole(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetAllCredentialsByRole called==============================")
	param := mux.Vars(r)
	credentials := []general.Credential{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Fill the r.Form
	r.ParseForm()

	err = controller.LoginService.GetTargetCommunityList(&credentials, tenantID, r.Form)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, credentials)
}
