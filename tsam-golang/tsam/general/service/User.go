package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// UserService gives method for listing for all users that is salesperson and admin.
type UserService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewUserService returns a new instance Of UserService.
func NewUserService(db *gorm.DB, repository repository.Repository) *UserService {
	return &UserService{
		DB:           db,
		Repository:   repository,
		associations: []string{"Country", "State", "Role"},
	}
}

// GetSpecificUsers returns all users by search criteria in database.
func (service *UserService) GetSpecificUsers(tenantID uuid.UUID, users *[]general.User, form url.Values,
	limit, offset int, totalCount *int) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		uow.RollBack()
		return err
	}

	// Get all sales people from database.
	err := service.Repository.GetAllInOrder(uow, users, "`first_name`",
		service.addSearchQueries(form),
		repository.PreloadAssociations(service.associations),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetUserList gets list of all salesperson and admin
func (service *UserService) GetUserList(tenantID uuid.UUID, users *[]list.User) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Repo get all call.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, users, "users.`first_name`",
		repository.Filter("`is_active` = ?", true))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// GetUserCredentialList gets list of all users from credential table.
func (service *UserService) GetUserCredentialList(tenantID uuid.UUID, credentials *[]list.Credential) error {
	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Repo get all call.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAll(uow, &credentials,
		repository.Join("INNER JOIN roles ON credentials.`role_id` = roles.`id` AND "+
			"credentials.`tenant_id` = roles.`tenant_id` AND roles.`deleted_at` IS NULL"),
		// Remove select if full credential is needed.
		repository.Select([]string{"credentials.`id`", "credentials.`first_name`", "credentials.`last_name`"}),
		// Remove admin if not needed.
		repository.Filter("roles.`deleted_at` IS NULL AND (roles.`role_name`=? OR roles.`role_name`=?) AND "+
			"credentials.`tenant_id` = ? AND credentials.`deleted_at` IS NULL", "salesperson", "admin", tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// Niranjan

// AddUser adds new sales person to database.
func (service *UserService) AddUser(user *general.User) error {
	var err error
	tenantID := user.TenantID
	credentialID := user.CreatedBy
	uow := repository.NewUnitOfWork(service.DB, false)

	// Assign User Code.
	user.Code, err = util.GenerateUniqueCode(uow.DB, user.FirstName, "`code` = ?", &general.User{})
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	roleName := user.Role.RoleName

	// Extract id from objects.
	if err = service.extractID(user); err != nil {
		uow.RollBack()
		return err
	}

	// Validate tenant id.
	if err = service.doForeignKeysExist(credentialID, user); err != nil {
		uow.RollBack()
		return err
	}

	// Check if email already exists.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.User{},
		repository.Filter("`email` = ? AND `role_id` = ? AND `id` NOT IN (?)",
			user.Email, user.RoleID, user.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	if exists {
		// Get users that exist.
		if err := service.Repository.GetRecordForTenant(uow, tenantID, user,
			repository.Filter("`email` = ?", user.Email)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		// Create credential for user.
		if err := service.AddCredential(uow, user, credentialID, user.RoleID, roleName); err != nil {
			uow.RollBack()
			return err
		}

		uow.Commit()

		// Return error to user for same email.
		return errors.NewValidationError("Email already exists. Login added!")

	}
	// Add user to database.
	if err := service.Repository.Add(uow, user); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("User could not be added", http.StatusInternalServerError)
	}

	// Create credential for user.
	if err := service.AddCredential(uow, user, credentialID, user.RoleID, roleName); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateUser updates sales person in database.
func (service *UserService) UpdateUser(user *general.User) error {
	uow := repository.NewUnitOfWork(service.DB, false)
	tenantID := user.TenantID

	err := service.validateFieldUniqueness(user)
	if err != nil {
		uow.RollBack()
		return err
	}

	roleName := user.Role.RoleName

	// Extract id from objects.
	if err = service.extractID(user); err != nil {
		uow.RollBack()
		return err
	}

	// Check if tenant exists.
	err = service.doForeignKeysExist(user.UpdatedBy, user)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Create bucket for getting user already present in database.
	tempUser := general.User{}

	// Get user for getting created_by field of user from database.
	if err := service.Repository.GetForTenant(uow, tenantID, user.ID, &tempUser,
		repository.Select("`created_by`")); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("User not found")
	}

	// Assign initial created_by id as we are using "save".
	user.CreatedBy = tempUser.CreatedBy

	// Update user.
	if err := service.Repository.Save(uow, user); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("User could not be updated")
	}

	// updates is_active field of user credential.
	err = service.updateUserCredential(uow, user, roleName)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Update credentials #Niranjan.
	uow.Commit()
	return nil
}

// DeleteUser deletes sales person from database.
func (service *UserService) DeleteUser(user *general.User) error {
	uow := repository.NewUnitOfWork(service.DB, false)
	tenantID := user.TenantID
	credentialID := user.DeletedBy

	// Validate credential id.
	err := service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Validate user id.
	err = service.doesUserExist(tenantID, user.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Update user for updating deleted_by field of user.
	err = service.Repository.UpdateWithMap(uow, user, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	})
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("User could not be deleted", http.StatusInternalServerError)
	}

	// Delete credentials.
	exist, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, new(general.Credential),
		repository.Filter("`user_id` = ?", user.ID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	if exist {
		err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		}, repository.Filter("`credentials`.`user_id` = ?", user.ID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("User credential could not be deleted", http.StatusInternalServerError)
		}
	} else {
		err = service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		}, repository.Filter("`credentials`.`sales_person_id` = ?", user.ID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("User credential could not be deleted", http.StatusInternalServerError)
		}
	}

	uow.Commit()
	return nil
}

// GetSalesPeopleList returns list of all salespeople.
func (service *UserService) GetSalesPeopleList(tenantID uuid.UUID, salesPeople *[]list.User) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		uow.RollBack()
		return err
	}

	// Get all sales people from database.
	// err := service.Repository.GetAllInOrderForTenant(uow, tenantID, salesPeople, "first_name")
	err := service.Repository.GetAllInOrder(uow, salesPeople, "`first_name`",
		repository.Join("INNER JOIN roles ON users.`role_id` = roles.`id` AND "+
			"users.`tenant_id` = roles.`tenant_id` AND roles.`deleted_at` IS NULL"),
		repository.Filter("users.`tenant_id` = ? AND roles.`role_name` = ? AND users.`is_active` = ?",
			tenantID, "salesperson", true))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// AddCredential adds credential to credential table according to salesperson details and generates password.
func (service *UserService) AddCredential(uow *repository.UnitOfWork, user *general.User,
	credenetialID, roleID uuid.UUID, roleName string) error {
	// Create bucket for credential.
	credential := general.Credential{}

	// Give salesperson details to credential.
	credential.CreatedBy = credenetialID
	credential.FirstName = user.FirstName
	credential.LastName = &user.LastName
	credential.Email = user.Email
	credential.Contact = user.Contact
	credential.RoleID = roleID
	credential.TenantID = user.TenantID

	if roleName == "Admin" {
		credential.UserID = &user.ID
	}
	if roleName == "SalesPerson" {
		credential.SalesPersonID = &user.ID
	}

	// Change in future #niranjan
	//credential.Password = util.GeneratePassword()
	credential.Password = strings.ToLower(user.FirstName)

	// Add credential to database.
	loginService := NewCredentialService(uow.DB, service.Repository)
	if err := loginService.AddCredential(&credential, uow); err != nil {
		return err
	}
	return nil
}

// updateUserCredential updates specified user record in credentials table.
func (service *UserService) updateUserCredential(uow *repository.UnitOfWork, user *general.User, roleName string) error {

	if roleName == "Admin" {
		err := service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
			"UpdatedBy": user.UpdatedBy,
			"IsActive":  user.IsActive,
		}, repository.Filter("`user_id` = ?", user.ID))
		if err != nil {
			return err
		}
	}
	if roleName == "SalesPerson" {
		err := service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
			"UpdatedBy": user.UpdatedBy,
			"IsActive":  user.IsActive,
		}, repository.Filter("`sales_person_id` = ?", user.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doForeignKeysExist checks if all foregin key exist and if not returns error.
func (service *UserService) doForeignKeysExist(credentialID uuid.UUID, user *general.User) error {

	tenantID := user.TenantID
	//  Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	//  Check if credential exists.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	err = service.doesRoleExist(tenantID, user.RoleID)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *UserService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *UserService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesUserExist validates if user exists or not in database.
func (service *UserService) doesUserExist(tenantID, userID uuid.UUID) error {
	//check user exists or not
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id` = ?", userID))
	if err := util.HandleError("Invalid user ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesRoleExist validates if role exists or not in database.
func (service *UserService) doesRoleExist(tenantID, roleID uuid.UUID) error {
	//check parent role exists or not
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Role{},
		repository.Filter("`id` = ?", roleID))
	if err := util.HandleError("Invalid role ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// validateFieldUniqueness will check the table for any repitition for unique fields.
func (service *UserService) validateFieldUniqueness(user *general.User) error {
	// Return error if any record has the same email in DB.
	exists, err := repository.DoesRecordExistForTenant(service.DB, user.TenantID, &general.User{},
		repository.Filter("`email` = ? AND `role_id` = ? AND `id` NOT IN (?)", user.Email, user.RoleID, user.ID))
	if err := util.HandleIfExistsError("Record already exists with the same email.", exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// extractID extracts ID from object.
func (service *UserService) extractID(user *general.User) error {
	if user.Country != nil {
		user.CountryID = &user.Country.ID
	}
	if user.State != nil {
		user.StateID = &user.State.ID
	}
	user.RoleID = user.Role.ID
	return nil
}

// adds search queries accordingly.
func (service *UserService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In user search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if firstName, ok := requestForm["firstName"]; ok {
		util.AddToSlice("`users`.`first_name`", "LIKE ?", "AND", "%"+firstName[0]+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if lastName, ok := requestForm["lastName"]; ok {
		util.AddToSlice("`users`.`last_name`", "LIKE ?", "AND", "%"+lastName[0]+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if email, ok := requestForm["email"]; ok {
		util.AddToSlice("`users`.`email`", "LIKE ?", "AND", "%"+email[0]+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if roleID, ok := requestForm["roleID"]; ok {
		util.AddToSlice("`users`.`role_id`", "= ?", "AND", roleID,
			&columnNames, &conditions, &operators, &values)
	}
	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`users`.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
