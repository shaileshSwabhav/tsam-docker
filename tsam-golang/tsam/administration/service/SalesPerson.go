package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	genService "github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SalesPersonService provides methods to do Update, Delete, Add, Get operations on SalesPerson.
type SalesPersonService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewSalesPersonService returns new instance of SalesPersonService.
func NewSalesPersonService(db *gorm.DB, repository repository.Repository) *SalesPersonService {
	return &SalesPersonService{
		DB:         db,
		Repository: repository,
	}
}

// AddSalesPerson adds new sales person to database.
func (ser *SalesPersonService) AddSalesPerson(salesPerson *general.User, credentialID uuid.UUID,
	tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(ser.DB, false)

	//giving tenant id to salesperson
	salesPerson.TenantID = tenantID

	// Assign Salesperson Code
	var codeError error
	salesPerson.Code, codeError = util.GenerateUniqueCode(uow.DB, salesPerson.FirstName, "`code` = ?", &general.User{})
	if codeError != nil {
		log.NewLogger().Error(codeError.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	//validate compulsary fields
	if err := salesPerson.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	//validate tenant id
	if err := ser.ValidateTenant(uow, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate credential id
	if err := ser.ValidateCredential(uow, credentialID, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate foreign key ids
	if err := ser.ValidateRole(uow, salesPerson.TenantID, salesPerson.RoleID); err != nil {
		uow.RollBack()
		return err
	}

	//give credential id to created_by of salesperson
	salesPerson.CreatedBy = credentialID

	//check if email already exists
	exists, err := ser.DoesEmailExist(uow, salesPerson)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if exists {
		//get salesPerson that exists
		if err := ser.Repository.GetRecordForTenant(uow, salesPerson.TenantID, salesPerson,
			repository.Filter("email=?", salesPerson.Email)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		//create credential for salesPerson
		if err := ser.AddCredential(uow, salesPerson, credentialID, salesPerson.RoleID); err != nil {
			uow.RollBack()
			return err
		}

		uow.Commit()

		//return error to user for same email
		return errors.NewValidationError("Email already exists")

	}
	// Add salesperson to database
	if err := ser.Repository.Add(uow, salesPerson); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("SalesPerson could not be added", http.StatusInternalServerError)
	}

	//create credential for salesPerson
	if err := ser.AddCredential(uow, salesPerson, credentialID, salesPerson.RoleID); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateSalesPerson updates sales person in database.
func (ser *SalesPersonService) UpdateSalesPerson(salesPerson *general.User, salesPersonID uuid.UUID,
	tenantID uuid.UUID, credentialID uuid.UUID) error {
	uow := repository.NewUnitOfWork(ser.DB, false)

	//give salesPerson id to salesPerson
	salesPerson.ID = salesPersonID

	//give tenant id to salesPerson
	salesPerson.TenantID = tenantID

	//validate compulsary fields
	if err := salesPerson.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	//validate tenant id
	if err := ser.ValidateTenant(uow, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate credential id
	if err := ser.ValidateCredential(uow, credentialID, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate salesPerson id
	if err := ser.ValidateSalespersonID(uow, salesPerson.ID, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate foreign key ids
	if err := ser.ValidateRole(uow, salesPerson.TenantID, salesPerson.RoleID); err != nil {
		uow.RollBack()
		return err
	}

	//Check if email exists or not
	exists, err := ser.DoesEmailExist(uow, salesPerson)
	if err != nil || exists {
		uow.RollBack()
		return err
	}

	//create bucket for getting salesPerson already present in database
	tempSalesPerson := general.User{}

	//get salesPerson for getting created_by field of salesPerson from database
	if err := ser.Repository.GetForTenant(uow, tenantID, salesPersonID, &tempSalesPerson); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//give created_by id from temp salesPerson to salesPerson to be updated
	salesPerson.CreatedBy = tempSalesPerson.CreatedBy

	//give credential id to updated_by of salesPerson
	salesPerson.UpdatedBy = credentialID

	//update sales person
	if err := ser.Repository.Save(uow, salesPerson); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Sales person could not be updated")
	}

	uow.Commit()
	return nil
}

// DeleteSalesPerson deletes sales person from database.
func (ser *SalesPersonService) DeleteSalesPerson(salesPerson *general.User, tenantID uuid.UUID,
	salesPersonID uuid.UUID, credentialID uuid.UUID) error {
	uow := repository.NewUnitOfWork(ser.DB, false)

	//give salesPerson id to salesPerson
	salesPerson.ID = salesPersonID

	//give tenant id to salesPerson
	salesPerson.TenantID = tenantID

	//validate tenant id
	if err := ser.ValidateTenant(uow, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate credential id
	if err := ser.ValidateCredential(uow, credentialID, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate salesPerson id
	if err := ser.ValidateSalespersonID(uow, salesPerson.ID, salesPerson.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get salesPerson for updating deleted_by field of salesPerson
	if err := ser.Repository.GetForTenant(uow, tenantID, salesPersonID, salesPerson); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//give salesPerson id to deleted_by of salesPerson
	salesPerson.DeletedBy = credentialID

	//update salesPerson for updating deleted_by field of salesPerson
	if err := ser.Repository.UpdateWithMap(uow, salesPerson, map[string]interface{}{"DeletedBy": credentialID}); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("SalesPerson could not be delted", http.StatusInternalServerError)
	}

	// delete salesPerson from database
	if err := ser.Repository.Delete(uow, salesPerson); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete salesPerson", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetSalesPerson returns all SalesPersones.
func (ser *SalesPersonService) GetSalesPerson(salesPersonID uuid.UUID, salesPerson *general.User, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(ser.DB, true)

	//giving tenant id to salesPerson
	salesPerson.TenantID = tenantID

	//validate tenant id
	if err := ser.ValidateTenant(uow, tenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get sales person
	err := ser.Repository.Get(uow, salesPersonID, salesPerson, repository.PreloadAssociations([]string{"Country", "State"}))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//to get only date for date of birth
	if salesPerson.DateOfBirth != nil {
		tempDateOfBirth := *salesPerson.DateOfBirth
		tempDateOfBirth = tempDateOfBirth[:10]
		salesPerson.DateOfBirth = &tempDateOfBirth
	}

	//to get only date for date of joining
	if salesPerson.DateOfJoining != nil {
		tempDateOfJoining := *salesPerson.DateOfJoining
		tempDateOfJoining = tempDateOfJoining[:10]
		salesPerson.DateOfJoining = &tempDateOfJoining
	}

	uow.Commit()
	return nil
}

// GetAllSalesPeople returns all sales people in database.
func (ser *SalesPersonService) GetAllSalesPeople(tenantID uuid.UUID, salesPeople *[]general.User) error {
	uow := repository.NewUnitOfWork(ser.DB, true)

	//validate tenant id
	if err := ser.ValidateTenant(uow, tenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get all sales people from database
	// queryProcessors := []repository.QueryProcessor{}
	// queryProcessors = append(queryProcessors, repository.Join("left join roles on users.role_id = roles.id"))
	// queryProcessors = append(queryProcessors, repository.Filter("roles.role_name=? AND users.tenant_id=? AND roles.tenant_id=?",
	// 	"salesperson", tenantID, tenantID))
	// queryProcessors = append(queryProcessors, repository.PreloadAssociations([]string{"Country", "State"}))
	err := ser.Repository.GetAllInOrder(uow, salesPeople, "first_name",
		repository.Join("left join roles on users.role_id = roles.id"),
		repository.Filter("roles.role_name=? AND users.tenant_id=? AND roles.tenant_id=?", "salesperson", tenantID, tenantID),
		repository.PreloadAssociations([]string{"Country", "State"}))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// //get all sales people from database
	// err := ser.Repository.GetAllInOrderForTenant(uow, tenantID, salesPeople, "first_name",
	// 	repository.PreloadAssociations([]string{"Country", "State"}))
	// if err != nil {
	// 	uow.RollBack()
	// 	log.NewLogger().Error(err.Error())
	// 	return errors.NewValidationError("Record not found")
	// }

	if len(*salesPeople) != 0 {
		for i := 0; i < len(*salesPeople); i++ {
			//to get only date for date of birth
			if (*salesPeople)[i].DateOfBirth != nil {
				tempDateOfBirth := *((*salesPeople)[i].DateOfBirth)
				tempDateOfBirth = tempDateOfBirth[:10]
				(*salesPeople)[i].DateOfBirth = &tempDateOfBirth
			}

			//to get only date for date of joining
			if (*salesPeople)[i].DateOfJoining != nil {
				tempDateOfJoining := *((*salesPeople)[i].DateOfJoining)
				tempDateOfJoining = tempDateOfJoining[:10]
				(*salesPeople)[i].DateOfJoining = &tempDateOfJoining
			}
		}
	}

	uow.Commit()
	return nil
}

// GetSalesPeopleList returns all sales people in database.
func (ser *SalesPersonService) GetSalesPeopleList(tenantID uuid.UUID, salesPeople *[]list.User) error {
	uow := repository.NewUnitOfWork(ser.DB, true)

	//validate tenant id
	if err := ser.ValidateTenant(uow, tenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get all sales people from database
	// err := ser.Repository.GetAllInOrderForTenant(uow, tenantID, salesPeople, "first_name")
	err := ser.Repository.GetAllInOrder(uow, salesPeople, "first_name",
		repository.Join("left join roles on users.role_id = roles.id"),
		repository.Filter("roles.role_name=? AND users.tenant_id=? AND roles.tenant_id=? AND roles.`deleted_at` IS NULL",
			"salesperson", tenantID, tenantID),
		repository.Filter("users.`deleted_at` IS NULL AND users.`tenant_id`=?", tenantID))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// ValidateTenant validates if tenant exists or not in database
func (ser *SalesPersonService) ValidateTenant(uow *repository.UnitOfWork, tenantID uuid.UUID) error {
	//check if tenant(parent tenant) exists or not
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// ValidateCredential validates if credental exists or not in database
func (ser *SalesPersonService) ValidateCredential(uow *repository.UnitOfWork, credentialID uuid.UUID, tenantID uuid.UUID) error {
	//check if credential exists or not
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Credential{}, repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// ValidateSalespersonID if salesperson exists or not in database
func (ser *SalesPersonService) ValidateSalespersonID(uow *repository.UnitOfWork,
	salespersonID uuid.UUID, tenantID uuid.UUID) error {
	//check salesperson exists or not
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.User{},
		repository.Filter("`id` = ?", salespersonID))
	if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// ValidateRole validates if role exists or not in database
func (ser *SalesPersonService) ValidateRole(uow *repository.UnitOfWork, tenantID uuid.UUID, roleID uuid.UUID) error {
	//check parent role exists or not
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Role{},
		repository.Filter("`id` = ?", roleID))
	if err := util.HandleError("Invalid role ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

//AddCredential adds credential to credential table according to salesperson details and generates password
func (ser *SalesPersonService) AddCredential(uow *repository.UnitOfWork, salesperson *general.User,
	credenetialID uuid.UUID, roleID uuid.UUID) error {
	//create bucket for credential
	credential := general.Credential{}

	//give salesperson details to credential
	credential.FirstName = salesperson.FirstName
	credential.LastName = &salesperson.LastName
	credential.Email = salesperson.Email
	//credential.Password = util.GeneratePassword()
	credential.Contact = salesperson.Contact
	credential.RoleID = roleID
	credential.TenantID = salesperson.TenantID
	credential.SalesPersonID = &salesperson.ID
	credential.Password = "salesperson"
	credential.CreatedBy = credenetialID

	//add credential to database
	loginService := genService.NewCredentialService(ser.DB, ser.Repository)
	if err := loginService.AddCredential(&credential, uow); err != nil {
		return err
	}
	return nil
}

//DoesEmailExist check for same email conflict, if email exists return true
func (ser *SalesPersonService) DoesEmailExist(uow *repository.UnitOfWork, salesperson *general.User) (bool, error) {
	var count int
	if err := ser.Repository.GetCountForTenant(uow, salesperson.TenantID, &general.User{}, &count,
		repository.Filter("email=? AND id NOT IN (?)", salesperson.Email, salesperson.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if count != 0 { //email already present
		log.NewLogger().Error("Validate email:Email already exists")
		return true, nil
	}

	return false, nil
}
