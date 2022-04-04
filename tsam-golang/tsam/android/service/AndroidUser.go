package services

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

type AndroidUserService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

func NewAndroidUserService(db *gorm.DB, repo repository.Repository) *AndroidUserService {
	return &AndroidUserService{
		DB:         db,
		Repository: repo,
	}
}

func (service *AndroidUserService) AddAndroidUser(user *general.AndroidUser) error {
	uow := repository.NewUnitOfWork(service.DB, false)

	var err error
	user.Password, err = util.EncryptPassword(user.Password)

	if err != nil {
		return err
	}
	err = service.Repository.Add(uow, &user)
	if err != nil {
		uow.DB.Rollback()
		return err

	}
	uow.Commit()
	return nil
}

func (service *AndroidUserService) Login(androidUser *general.AndroidUser, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	tempPassword := androidUser.Password

	tempCredentials := []general.AndroidUser{}
	isCredentialFound := false

	err = service.Repository.GetAllForTenant(uow, tenantID, &tempCredentials,
		repository.Filter("`email` = ? ", androidUser.Email))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, cred := range tempCredentials {

		isCredentialFound = util.DoPasswordsMatch(cred.Password, tempPassword)

		if isCredentialFound {
			*androidUser = cred
			goto passwordMatched
		}
	}
	// If loop ends and there is no password match, return error.
	return errors.NewValidationError("Login Failed! Email and Password did not match")

passwordMatched:

	// // Check if deptartment exist for role. If exist get the department ID
	// tempDepartment := general.Department{}
	// exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
	// 	repository.Filter("`role_id` = ?", androidUser.RoleID))
	// if err != nil {
	// 	return err
	// }
	// if exist {
	// 	err = service.Repository.GetRecordForTenant(uow, tenantID, &tempDepartment,
	// 		repository.Select("`id`"), repository.Filter("`role_id` = ?", androidUser.RoleID))
	// 	if err != nil {
	// 		uow.RollBack()
	// 		return err
	// 	}
	// 	if util.IsUUIDValid(tempDepartment.ID) {
	// 		tokenDTO.DepartmentID = &tempDepartment.ID
	// 	}
	// 	// androidUser.DepartmentID = &tempDepartment.ID
	// }

	// loginSession := general.LoginSession{
	// 	CredentialID: androidUser.ID,
	// 	// RoleID:       androidUser.RoleID,
	// 	TenantBase: general.TenantBase{
	// 		TenantID:  tenantID,
	// 		CreatedBy: androidUser.ID,
	// 	},
	// }

	// loginSessionService := NewLoginSessionService(service.DB, service.Repository)

	// err = loginSessionService.AddLoginSession(&loginSession, uow)
	// if err != nil {
	// 	uow.RollBack()
	// 	return errors.NewValidationError(err.Error())
	// }

	// tokenDTO.CredentialID = androidUser.ID
	// tokenDTO.FirstName = androidUser.FirstName
	// tokenDTO.Email = androidUser.Email
	// tokenDTO.TenantID = androidUser.TenantID
	// tokenDTO.RoleID = androidUser.RoleID
	// tokenDTO.LoginSessionID = &loginSession.ID

	// if androidUser.LastName != nil {
	// 	tokenDTO.LastName = *androidUser.LastName
	// }

	uow.Commit()
	return nil
}
func (service *AndroidUserService) doesTenantExist(tenantID uuid.UUID) error {
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
