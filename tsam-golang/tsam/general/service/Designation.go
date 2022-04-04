package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DesignationService provide method to update, delete, add, get method for Designation.
type DesignationService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewDesignationService returns new instance of DesignationService.
func NewDesignationService(db *gorm.DB, repository repository.Repository) *DesignationService {
	return &DesignationService{
		DB:         db,
		Repository: repository,
	}
}

// AddDesignation adds new designation to database.
func (service *DesignationService) AddDesignation(designation *general.Designation, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(designation.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesPositionExist(designation)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(designation.CreatedBy, designation.TenantID)
	if err != nil {
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add desgination to database.
	err = service.Repository.Add(uow, designation)
	if err != nil {
		uow.RollBack()
		return err
	}

	// If calling function does not pass the uow only then make commit.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddDesignations adds multiple designations.
func (service *DesignationService) AddDesignations(designations *[]general.Designation, designationIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*designations); i++ {
		for j := 0; j < len(*designations); j++ {
			if i != j && (*designations)[i].Position == (*designations)[j].Position {
				log.NewLogger().Error("Position:" + (*designations)[j].Position + " exists")
				return errors.NewValidationError("Position:" + (*designations)[j].Position + " exists")
			}
		}
	}

	// Add individual designation.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, designation := range *designations {
		designation.TenantID = tenantID
		designation.CreatedBy = credentialID
		err := service.AddDesignation(&designation, uow)
		if err != nil {

			return err
		}
		*designationIDs = append(*designationIDs, designation.ID)
	}

	uow.Commit()
	return nil
}

// UpdateDesignation updates designation to database.
func (service *DesignationService) UpdateDesignation(designation *general.Designation) error {
	// Validate tenant ID.
	err := service.doesTenantExist(designation.TenantID)
	if err != nil {
		return err
	}

	// Validate designation ID.
	err = service.doesDesignationExist(designation.ID, designation.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(designation.UpdatedBy, designation.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesPositionExist(designation)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update designation.
	err = service.Repository.Update(uow, designation)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteDesignation deletes designation from database.
func (service *DesignationService) DeleteDesignation(designation *general.Designation) error {
	credentialID := designation.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(designation.TenantID)
	if err != nil {
		return err
	}

	// Validate designation ID.
	err = service.doesDesignationExist(designation.ID, designation.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, designation.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update designation for updating deleted_by and deleted_at fields of designation
	if err := service.Repository.UpdateWithMap(uow, designation, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", designation.TenantID)); err != nil {

		uow.RollBack()
		return errors.NewHTTPError("Designation could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetDesignation returns one designation.
func (service *DesignationService) GetDesignation(designation *general.Designation, tenantID, designationID uuid.UUID) error {
	// Validate tenant ID
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}
	// validate designation ID
	exists, err := repository.DoesRecordExist(service.DB, designation, repository.Filter("`id` = ?", designationID))
	if err != nil {
		return err
	}
	if !exists {
		return errors.NewValidationError("Designation not found")
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, tenantID, designationID, designation)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetDesignations returns all designations.
func (service *DesignationService) GetDesignations(designation *[]general.Designation, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get all designations from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, designation, "`position`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetDesignationList returns designation list.
func (service *DesignationService) GetDesignationList(designation *[]general.Designation, tenantID uuid.UUID) error {
	// Validate tenant ID
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get designation list.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, designation, "`position`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search queries.
func (service *DesignationService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	if _, ok := requestForm["position"]; ok {
		return repository.Filter("`position` LIKE ?", "%"+requestForm.Get("position")+"%")
	}
	return nil
}

// doesPositionExist returns true if the position already exists for the designation in database.
func (service *DesignationService) doesPositionExist(designation *general.Designation) error {
	// Check for same position conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, designation.TenantID, &general.Designation{},
		repository.Filter("`position`=? AND `id`!=?", designation.Position, designation.ID))
	if err := util.HandleIfExistsError("Position:"+designation.Position+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *DesignationService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDesignationExist validates if designation exists or not in database.
func (service *DesignationService) doesDesignationExist(designationID uuid.UUID, tenantID uuid.UUID) error {
	// Check designation exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Designation{},
		repository.Filter("`id` = ?", designationID))
	if err := util.HandleError("Invalid designation ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *DesignationService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
