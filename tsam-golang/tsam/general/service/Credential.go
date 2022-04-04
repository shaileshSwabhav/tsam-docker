package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CredentialService provides methods to add, get operation login.
type CredentialService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCredentialService returns new instance of LoginService
func NewCredentialService(db *gorm.DB, repository repository.Repository) *CredentialService {
	return &CredentialService{
		DB:         db,
		Repository: repository,
	}
}

// Login will verify user before login.
func (service *CredentialService) Login(tokenDTO *general.TokenDTO, credential *general.Credential, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	tempPassword := credential.Password

	tempCredentials := []general.Credential{}
	isCredentialFound := false

	err = service.Repository.GetAllForTenant(uow, tenantID, &tempCredentials,
		repository.Filter("`email` = ? AND `is_active` = '1'", credential.Email))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, cred := range tempCredentials {

		isCredentialFound = util.DoPasswordsMatch(cred.Password, tempPassword)

		if isCredentialFound {
			*credential = cred
			goto passwordMatched
		}
	}
	// If loop ends and there is no password match, return error.
	return errors.NewValidationError("Login Failed! Email and Password did not match")

passwordMatched:

	// Check if deptartment exist for role. If exist get the department ID
	tempDepartment := general.Department{}
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
		repository.Filter("`role_id` = ?", credential.RoleID))
	if err != nil {
		return err
	}
	if exist {
		err = service.Repository.GetRecordForTenant(uow, tenantID, &tempDepartment,
			repository.Select("`id`"), repository.Filter("`role_id` = ?", credential.RoleID))
		if err != nil {
			uow.RollBack()
			return err
		}
		if util.IsUUIDValid(tempDepartment.ID) {
			tokenDTO.DepartmentID = &tempDepartment.ID
		}
		// credential.DepartmentID = &tempDepartment.ID
	}

	loginSession := general.LoginSession{
		CredentialID: credential.ID,
		// RoleID:       credential.RoleID,
		TenantBase: general.TenantBase{
			TenantID:  tenantID,
			CreatedBy: credential.ID,
		},
	}

	loginSessionService := NewLoginSessionService(service.DB, service.Repository)

	err = loginSessionService.AddLoginSession(&loginSession, uow)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError(err.Error())
	}

	tokenDTO.CredentialID = credential.ID
	tokenDTO.FirstName = credential.FirstName
	tokenDTO.Email = credential.Email
	tokenDTO.TenantID = credential.TenantID
	tokenDTO.RoleID = credential.RoleID
	tokenDTO.LoginSessionID = &loginSession.ID

	if credential.LastName != nil {
		tokenDTO.LastName = *credential.LastName
	}

	uow.Commit()
	return nil
}

// Logout will set session end time and logout.
func (service *CredentialService) Logout(credential *general.Credential, tenantID, credentialID uuid.UUID,
	requestForm url.Values) error {

	var loginSessionID uuid.UUID

	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// get loginsessionID
	if _, ok := requestForm["loginSessionID"]; ok {
		loginSessionID, err = util.ParseUUID(requestForm.Get("loginSessionID"))
		if err != nil {
			return err
		}
	}

	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, credential,
		repository.Filter("`id` = ?", credentialID))
	if err != nil {
		return err
	}
	if !exists {
		log.NewLogger().Error("CredentialID not found")
		return errors.NewValidationError("CredentialID not found")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	loginSessionService := NewLoginSessionService(service.DB, service.Repository)

	err = loginSessionService.EndLoginSession(uow, credentialID, tenantID, loginSessionID)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError(err.Error())
	}

	uow.Commit()
	return nil
}

// AddCredential will add the new user to login table.
// when adding a credential all its fields must be filled in the calling function itself
// eg : 	password := "faculty" // change to util.GeneratePassword()
// credentials := general.Credential{
// 	FirstName: faculty.FirstName,
// 	LastName:  faculty.LastName,
// 	Email:     faculty.Email,
// 	Password:  password,
// 	FacultyID: &faculty.ID,
// 	Contact:   faculty.Contact,
// 	RoleID:    role.ID,
// }
// credentials.TenantID = tenantID
// credentials.CreatedBy = credentialID
func (service *CredentialService) AddCredential(credential *general.Credential, uow *repository.UnitOfWork) error {

	err := credential.ValidateCredentials()
	if err != nil {
		return err
	}

	// check if tenant exist
	// err = service.doesTenantExist(credential.TenantID)
	// if err != nil {
	// 	return err
	// }

	// // check if credential exist
	// err = service.doesCredentialExist(credential.TenantID, credential.CreatedBy)
	// if err != nil {
	// 	return err
	// }

	// checks if roleID exists in roles table
	role := general.Role{}
	exists, err := repository.DoesRecordExist(service.DB, &role, repository.Filter("`id`=?", credential.RoleID))
	if err != nil {
		return errors.NewValidationError(err.Error())
	}
	if !exists {
		return err
	}
	// credential.ID = util.GenerateUUID()
	// login.Password = tempPassword
	err = service.doesEmailExistForRole(credential.TenantID, credential.RoleID, credential.Email)
	if err != nil {
		return err
	}
	err = service.doesLoginIDExist(credential)
	if err != nil {
		return err
	}
	credential.Password, err = util.EncryptPassword(credential.Password)
	if err != nil {
		return err
	}
	err = service.Repository.Add(uow, credential)
	if err != nil {
		return err
	}
	return nil
}

// // UpdateCredentialStatus updates credential record
// // tenantID and credentialID should be validated by calling method.
// func (service *CredentialService) UpdateCredentialStatus(uow *repository.UnitOfWork, tenantID, credentialID uuid.UUID,
// 	isActive bool) error {

// 	return nil
// }

// UpdateCredentialPassword updates the password for the credential.
func (service *CredentialService) UpdateCredentialPassword(credential *general.Credential) error {

	// Email should be validated and entire credential struct is not required.

	err := service.doesTenantExist(credential.TenantID)
	if err != nil {
		return err
	}

	// Find credential with similar email.
	exists, err := repository.DoesRecordExistForTenant(service.DB, credential.TenantID, &general.Credential{},
		repository.Filter("email=?", credential.Email))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if !exists {
		log.NewLogger().Error("Email not found")
		return errors.NewValidationError("Email not found")
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	credential.Password, err = util.EncryptPassword(credential.Password)
	if err != nil {
		return err
	}
	err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[interface{}]interface{}{
		"password": credential.Password,
	}, repository.Filter("email=?", credential.Email))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// PasswordVerification verifies the password.
func (service *CredentialService) PasswordVerification(passwordChange *general.PasswordChange) error {

	// Validate tenant id.
	err := service.doesTenantExist(passwordChange.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create bucket for getting password from database for the given email and role ID.
	tempCredential := general.Credential{}

	// Get credential from database.
	err = service.Repository.GetAllForTenant(uow, passwordChange.TenantID, &tempCredential,
		repository.Filter("`email`=? AND role_id=?", passwordChange.Email, passwordChange.RoleID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Verify the password.
	if !util.DoPasswordsMatch(tempCredential.Password, passwordChange.Password) {
		return errors.NewValidationError("Password entered is wrong")
	}

	return nil
}

// ChangePassword updates the password of the credential after login.
func (service *CredentialService) ChangePassword(passwordChange *general.PasswordChange) error {

	// Validate tenant id.
	err := service.doesTenantExist(passwordChange.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential id.
	err = service.doesCredentialExist(passwordChange.TenantID, passwordChange.CredentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for all credentials with same email id and different role ids.
	tempCredentials := []general.Credential{}

	// Get all credentials with the same email and different role ids in credential table.
	err = service.Repository.GetAllForTenant(uow, passwordChange.TenantID, &tempCredentials,
		repository.Filter("`email`=?", passwordChange.Email))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Create flag for knowing if any credential has the new passowrd already.
	credentialsWithSamePasswordCount := 0

	// Create flag for knowing if the same credential has the new passowrd already.
	doesSameCredentialHaveSamePassword := false

	// Range all the credentials with same email.
	for _, tempCredential := range tempCredentials {

		// Compare the new password with the already present password.
		if util.DoPasswordsMatch(tempCredential.Password, passwordChange.Password) {
			credentialsWithSamePasswordCount = credentialsWithSamePasswordCount + 1

			// Check if credential id is same as temp credential id.
			if passwordChange.CredentialID == tempCredential.ID {
				doesSameCredentialHaveSamePassword = true
			}
		}
	}

	// If password alreday exists then return error.
	if credentialsWithSamePasswordCount > 0 {
		// If password already exists for same credential id then return password must be different error.
		if doesSameCredentialHaveSamePassword {
			return errors.NewValidationError("Please enter a new password")
		}

		// If password already exists for same email id then return password with email exists error.
		return errors.NewValidationError("Password with same email already exists")
	}

	// Encrypt new password.
	passwordChange.Password, err = util.EncryptPassword(passwordChange.Password)
	if err != nil {
		return err
	}

	// Update the credential with new password.
	err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[interface{}]interface{}{
		"password":  passwordChange.Password,
		"UpdatedBy": passwordChange.CredentialID,
	}, repository.Filter("id=?", passwordChange.CredentialID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteCredential soft deletes the login record.
// credential -> should be filled foreign key
// tenantID -> ID of tenant for which credential is to be deleted
// deletedBy -> it will have the ID of the login which is deleting the record
// recordID -> will have the id of the faculty, talent, etc
// eg - if deleting faculty then value will be facultyID
// condition -> name of the column which is to be deleted
// eg - for faculty condition will be faculty_id
// credentialID -> ID which is currently logged in and deleting the record
// eg - for faculty -> credential.facultyID should be mentioned
// 	DeleteCredential(credential{faculty_id = "cfe25758-f5fe..."}, "cfe25758-f5fe...", "cfe25758-f5fe...", "cfe25758-f5fe...", "faculty_id=?" unitOfWork)
// requires optimization
func (service *CredentialService) DeleteCredential(credential *general.Credential, tenantID, deletedBy, recordID uuid.UUID,
	condition string, uow *repository.UnitOfWork) error {

	// check tenant ID
	// err := service.doesTenantExist(tenantID)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }
	// check credential ID
	// exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
	// 	repository.Filter("`id`=?", deletedBy))
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }
	// if !exists {
	// 	log.NewLogger().Error("Credential ID not found")
	// 	return errors.NewValidationError("Credential ID not found")
	// }

	// check if foreign key record exists (talentID, facultyID, userID, etc)
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Credential{},
		repository.Filter(condition, recordID))
	if err != nil {
		return err
	}
	if !exists {
		return errors.NewValidationError("Record does not exists")
	}

	// get credential ID of the record which is to deleted
	// err = service.Repository.GetRecordForTenant(uow, tenantID, credential,
	// 	repository.Filter(condition, recordID))
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }
	// credential.DeletedBy = deletedBy
	err = service.Repository.UpdateWithMap(uow, general.Credential{}, map[string]interface{}{
		"DeletedBy": deletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter(condition, recordID))
	if err != nil {
		return err
	}
	// err = service.Repository.DeleteForTenant(uow, tenantID, credential)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// ValidatePermission will check if the credential has permission to add, update or delete.
//  ValidatePermission("cfe25758-f5fe...", "cfe25758-f5fe...", "batch-master", "add")
func (service *CredentialService) ValidatePermission(tenantID, credentialID uuid.UUID, menuName, permissionName string) error {

	// fmt.Println("====================================Validate permission================================================")

	// Check if credential for exists.
	err := service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	cred := general.Credential{}
	menus := general.Menu{}
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, &cred, repository.Select("`role_id`"), repository.Filter("`id` = ?", credentialID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	// fmt.Println("*****roleID ->", cred)
	err = service.Repository.GetRecordForTenant(uow, tenantID, &menus,
		repository.Filter("`role_id`=? AND `url`=?", cred.RoleID, menuName))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()

	switch permissionName {
	case "add":
		if !menus.Permission.Add {
			return errors.NewValidationError("Not Authorized to perform add operation")
		}
	case "update":
		if !menus.Permission.Update {
			return errors.NewValidationError("Not Authorized to perform update operation")
		}
	case "delete":
		if !menus.Permission.Delete {
			return errors.NewValidationError("Not Authorized to perform delete operation")
		}
	}

	return nil
}

// GetTargetCommunityList will return all the records from credential table by role names.
func (service *CredentialService) GetTargetCommunityList(credentials *[]general.Credential, tenantID uuid.UUID, form url.Values) error {

	// Get query params for role name and login id.
	roleNames := form["roleNames"]

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, credentials, "credentials.`first_name`, credentials.`last_name`",
		repository.Join("INNER JOIN roles ON roles.id = credentials.role_id"),
		repository.Filter("roles.tenant_id=? AND roles.deleted_at IS NULL", tenantID),
		repository.Filter("credentials.tenant_id=?", tenantID),
		repository.Filter("role_name IN(?)", roleNames))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesTenantExist validates if tenant exists.
func (service *CredentialService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if !exists {
		log.NewLogger().Error("Tenant Not Found")
		return errors.NewValidationError("Tenant Not Found")
	}
	return nil
}

// doesCredentialExist validates if credential exists.
func (service *CredentialService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if !exists {
		log.NewLogger().Error("Credential Not Found")
		return errors.NewValidationError("Credential Not Found")
	}
	return nil
}

// doesLoginIDExist will check if talentID, salespersonID, facultyID, userID, collegeID, companyID exists in their table.
func (service *CredentialService) doesLoginIDExist(credential *general.Credential) error {

	var doesLoginIDExists bool = false

	var tempTables = map[*uuid.UUID]string{
		credential.SalesPersonID: "sales_person_id=?",
		credential.TalentID:      "talent_id=?",
		credential.FacultyID:     "faculty_id=?",
		credential.UserID:        "user_id=?",
		credential.CollegeID:     "college_id=?",
		credential.CompanyID:     "company_id=?",
		credential.EmployeeID:    "employee_id=?",
	}
	for id, condition := range tempTables {
		// !reflect.ValueOf(id).IsNil()
		if id != nil {
			doesLoginIDExists = true
			exist, err := repository.DoesRecordExistForTenant(service.DB, credential.TenantID, general.Credential{},
				repository.Filter(condition, id))
			if err != nil {
				return err
			}
			if exist {
				return errors.NewValidationError("Credential already exists")
			}
			// using this method as tenantID's are not specified in the tables
			// if repository.DoesRecordExist(service.DB, tableName, repository.Filter("`id` = ?", id)) {
			// 	uow.RollBack()
			// 	return errors.NewValidationError("Login Exists")
			// }
		}
	}
	if !doesLoginIDExists {
		return errors.NewValidationError("Login ID is required")
	}
	return nil
}

// doesEmailExistForRole validates if email is already present for the given role or not.
func (service *CredentialService) doesEmailExistForRole(tenantID, roleID uuid.UUID, email string) error {

	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Credential{},
		repository.Filter("`role_id`=? AND `email`=?", roleID, email))
	if err := util.HandleIfExistsError("Credential with same email and role exist", exist, err); err != nil {
		return err
	}

	return nil
}
