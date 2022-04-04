package general

import (
	"net/http"

	"github.com/gorilla/mux"
	services "github.com/techlabs/swabhav/tsam/android/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/security"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

type AndroidController struct {
	AndroidUserService *services.AndroidUserService
	log                log.Logger
	auth               *security.Authentication
}

func NewAndroidUserController(AndroidUserService *services.AndroidUserService, log log.Logger, auth *security.Authentication) *AndroidController {
	return &AndroidController{
		AndroidUserService: AndroidUserService,
		log:                log,
		auth:               auth,
	}
}

func (controller *AndroidController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	addAndroidUser := router.HandleFunc("/tenant/{tenantID}/androidUser", controller.AddAndroidUser).Methods(http.MethodPost)

	// login
	login := router.HandleFunc("/tenant/{tenantID}/android-login",
		controller.Login).Methods(http.MethodPost)
	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.UpdateAndroidUser).Methods(http.MethodPut)

	// delete
	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.DeleteAndroidUser).Methods(http.MethodDelete)

	// get
	// router.HandleFunc("/tenant/{tenantID}/company/{companyID}", controller.GetCompanyByI).Methods(http.MethodGet)
	// router.HandleFunc("/tenant/{tenantID}/company", controller.GetAllCompanies).Methods(http.MethodGet)

	*exclude = append(*exclude, addAndroidUser, login)
	controller.log.Info("Android User Routes Registered")
}
func (controller *AndroidController) AddAndroidUser(w http.ResponseWriter, r *http.Request) {
	controller.log.Info("==============================Add AndroidUser API Called==============================")

	androidUser := &general.AndroidUser{}

	// Fill the androidUser variable with given data.
	err := web.UnmarshalJSON(r, androidUser)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	parser := web.NewParser(r)
	// Parse and set tenant ID.
	//  util.ParseUUID(mux.Vars(r)["tenantID"])
	androidUser.TenantID, err = parser.GetTenantID()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse and set credentialID in CreatedBy field.
	// util.ParseUUID(mux.Vars(r)["credentialID"])
	// androidUser.CreatedBy, err = controller.auth.ExtractCredentialIDFromToken(r)
	// if err != nil {
	// 	controller.log.Error(err.Error())
	// 	web.RespondError(w, err)
	// 	return
	// }

	// Validate androidUser
	err = androidUser.Validate()
	if err != nil {
		controller.log.Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// call Service
	err = controller.AndroidUserService.AddAndroidUser(androidUser)
	if err != nil {
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "User added")
}
func (controller *AndroidController) Login(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================Login called==============================")
	credential := general.AndroidUser{}
	// tokenDTO := general.TokenDTO{}

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

	err = controller.AndroidUserService.Login(&credential, tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Creating login token.
	// token, err := controller.auth.GenerateToken(credential.ID.String(), credential.Email, credential.RoleID.String())
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	// 	return
	// }

	// tokenDTO.LoginID = controller.getLoginID(&credential)
	// tokenDTO.Token = token

	web.RespondJSON(w, http.StatusOK, "Login Success")
}

// func (controller *AndroidController) getLoginID(login *general.Credential) uuid.UUID {

// 	switch {
// 	case login.CollegeID != nil:
// 		return *login.CollegeID
// 	case login.CompanyID != nil:
// 		return *login.CompanyID
// 	case login.TalentID != nil:
// 		return *login.TalentID
// 	case login.FacultyID != nil:
// 		return *login.FacultyID
// 	case login.UserID != nil:
// 		return *login.UserID
// 	case login.SalesPersonID != nil:
// 		return *login.SalesPersonID
// 	case login.EmployeeID != nil:
// 		return *login.EmployeeID
// 	}
// 	return uuid.Nil
// }
