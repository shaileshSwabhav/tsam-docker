package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// RoleService provides methods to update, delete, add, get, get all and get all by role id for role
type RoleService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewRoleService returns new instance of RoleService.
func NewRoleService(db *gorm.DB, repo repository.Repository) *RoleService {
	return &RoleService{
		DB:         db,
		Repository: repo,
	}
}

// GetRole returns one role by specific role id form database
func (service *RoleService) GetRole(role *general.Role, roleID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	//giving tenant id to role
	role.TenantID = tenantID

	//validate tenant id
	if err := service.ValidateTenant(uow, tenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get role from database
	err := service.Repository.GetForTenant(uow, tenantID, roleID, role)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetRoles returns all roles from database
func (service *RoleService) GetRoles(roles *[]general.Role, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	//validate tenant id
	if err := service.ValidateTenant(uow, tenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get all roles form database
	err := service.Repository.GetAllForTenant(uow, tenantID, roles)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// AddRole adds one role in database
func (service *RoleService) AddRole(role *general.Role, tenantID uuid.UUID, credentialID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, false)

	//giving tenant id to role
	role.TenantID = tenantID

	//validate tenant id
	if err := service.ValidateTenant(uow, role.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//Check if role name exists or not
	err := service.DoesRoleNameExist(uow, role)
	if err != nil {
		uow.RollBack()
		return err
	}

	//give credential id to created_by of role
	role.CreatedBy = credentialID

	//add role in database
	err = service.Repository.Add(uow, role)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Role could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateRole updates one role by specific role id in database
func (service *RoleService) UpdateRole(role *general.Role, tenantID uuid.UUID, roleID uuid.UUID, credentialID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, false)

	//give role id to role
	role.ID = roleID

	//give tenant id to role
	role.TenantID = tenantID

	//validate tenant id
	if err := service.ValidateTenant(uow, role.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate role id
	if err := service.ValidateRoleID(uow, role.ID, role.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//Check if role name exists or not
	err := service.DoesRoleNameExist(uow, role)
	if err != nil {
		uow.RollBack()
		return err
	}

	//give credential id to updated_by of role
	role.UpdatedBy = credentialID

	//update role in database
	err = service.Repository.Update(uow, role)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Role could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteRole deletes one role by specific role id in database
func (service *RoleService) DeleteRole(role *general.Role, tenantID uuid.UUID, roleID uuid.UUID, credentialID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, false)

	//give role id to role
	role.ID = roleID

	//give tenant id to role
	role.TenantID = tenantID

	//validate tenant id
	if err := service.ValidateTenant(uow, role.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//validate role id
	if err := service.ValidateRoleID(uow, role.ID, role.TenantID); err != nil {
		uow.RollBack()
		return err
	}

	//get role for updating deleted_by field of role
	if err := service.Repository.GetForTenant(uow, tenantID, roleID, role); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//give credential id to deleted_by of role
	role.DeletedBy = credentialID

	//update role for updating deleted_by field of role
	if err := service.Repository.UpdateWithMap(uow, role, map[string]interface{}{"DeletedBy": credentialID}); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Role could not be delted", http.StatusInternalServerError)
	}

	//delete role from database
	err := service.Repository.Delete(uow, &role)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete role", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DoesRoleNameExist returns true if the role name already exists in database
func (service *RoleService) DoesRoleNameExist(uow *repository.UnitOfWork, role *general.Role) error {
	//check for same role name conflict
	exists, err := repository.DoesRecordExistForTenant(uow.DB, role.TenantID, &general.Role{},
		repository.Filter("role_name=? AND id!=?", role.RoleName, role.ID))
	if err := util.HandleIfExistsError("Role name exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// ValidateTenant validates if tenant exists or not in database
func (service *RoleService) ValidateTenant(uow *repository.UnitOfWork, tenantID uuid.UUID) error {
	//check if tenant(parent tenant) exists or not
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// ValidateRoleID validates if role exists or not in database
func (service *RoleService) ValidateRoleID(uow *repository.UnitOfWork, roleID uuid.UUID, tenantID uuid.UUID) error {
	//check role exists or not
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Role{},
		repository.Filter("`id` = ?", roleID))
	if err := util.HandleError("Invalid role ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
