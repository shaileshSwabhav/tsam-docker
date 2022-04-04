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
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// SpecializationService provides methods to update, delete, add, get, get all and get all by degree for specialization.
type SpecializationService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewSpecializationService returns new instance of SpecializationService.
func NewSpecializationService(db *gorm.DB, repository repository.Repository) *SpecializationService {
	return &SpecializationService{
		DB:         db,
		Repository: repository,
	}
}

// AddSpecialization adds new specialization in database.
func (service *SpecializationService) AddSpecialization(specialization *general.Specialization, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate degree id.
	err = service.doesDegreeexist(specialization.TenantID, specialization.DegreeID)
	if err != nil {
		return err
	}

	// Validate if same specialization branch name exists.
	err = service.doesBranchNameExist(specialization)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(specialization.CreatedBy, specialization.TenantID)
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

	// Add specialization to database
	if err := service.Repository.Add(uow, specialization); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Specialization could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetSpecializationList gets all specializations from database.
func (service *SpecializationService) GetSpecializationList(specializations *[]list.Specialization, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get specializations from database
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, specializations, "`branch_name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSpecializations gets all specializations from database.
func (service *SpecializationService) GetSpecializations(specializations *[]general.SpecializationDTO, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get specializations from database
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, specializations, "`branch_name`",
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Degree"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetSpecializationsByDegree gets all specializations from database by specific degree id.
func (service *SpecializationService) GetSpecializationsByDegree(specializations *[]general.Specialization,
	tenantID uuid.UUID, degreeID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Validate degree id.
	err = service.doesDegreeexist(tenantID, degreeID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get specializations from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, specializations, "`branch_name`",
		repository.Filter("`degree_id`=?", degreeID)); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// AddSpecializations adds multiple specializations in database.
func (service *SpecializationService) AddSpecializations(specializations *[]general.Specialization, specializationIDs *[]uuid.UUID,
	tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*specializations); i++ {
		for j := 0; j < len(*specializations); j++ {
			if i != j && (*specializations)[i].BranchName == (*specializations)[j].BranchName && (*specializations)[i].DegreeID == (*specializations)[j].DegreeID {
				log.NewLogger().Error("Name:" + (*specializations)[j].BranchName + " exists")
				return errors.NewValidationError("Name:" + (*specializations)[j].BranchName + " exists")
			}
		}
	}

	// Add individual specialization.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, specialization := range *specializations {
		specialization.TenantID = tenantID
		specialization.CreatedBy = credentialID
		err := service.AddSpecialization(&specialization, uow)
		if err != nil {
			return err
		}
		*specializationIDs = append(*specializationIDs, specialization.ID)
	}

	uow.Commit()
	return nil
}

// GetSpecialization gets one specialization by specific specialization id from database.
func (service *SpecializationService) GetSpecialization(specialization *general.Specialization) error {
	// Validate tenant id.
	err := service.doesTenantExist(specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate specialization id.
	err = service.doesSpecializationExist(specialization.ID, specialization.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get specialization,
	if err := service.Repository.GetForTenant(uow, specialization.TenantID, specialization.ID, specialization); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateSpecialization updates one specialization by specific specialization id in database.
func (service *SpecializationService) UpdateSpecialization(specialization *general.Specialization) error {
	// Validate tenant ID.
	err := service.doesTenantExist(specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate specialization ID.
	err = service.doesSpecializationExist(specialization.ID, specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(specialization.UpdatedBy, specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate if same branch name exists.
	err = service.doesBranchNameExist(specialization)
	if err != nil {
		return err
	}

	// Validate degree id.
	err = service.doesDegreeexist(specialization.TenantID, specialization.DegreeID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update specialization.
	if err := service.Repository.Update(uow, specialization); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Specialization could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteSpecialization deletes one specialization by specific specialization id from database.
func (service *SpecializationService) DeleteSpecialization(specialization *general.Specialization) error {
	credentialID := specialization.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate specialization ID.
	err = service.doesSpecializationExist(specialization.ID, specialization.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, specialization.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update specialization for updating deleted_by and deleted_at fields of specialization.
	if err := service.Repository.UpdateWithMap(uow, specialization, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", specialization.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Specialization could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *SpecializationService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDegreeexist validates if degree exists or not in database.
func (service *SpecializationService) doesDegreeexist(tenantID uuid.UUID, degreeID uuid.UUID) error {
	// Check parent degree exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Degree{},
		repository.Filter("`id` = ?", degreeID))
	if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSpecializationExist if specialization exists or not in database.
func (service *SpecializationService) doesSpecializationExist(specializationID uuid.UUID, tenantID uuid.UUID) error {
	// Check specialization exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Specialization{},
		repository.Filter("`id`=?", specializationID))
	if err := util.HandleError("Invalid specialization ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *SpecializationService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check if credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{}, repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBranchNameExist checks if specialization's branch name already exists or not for the speicifc degree in database.
func (service *SpecializationService) doesBranchNameExist(specialization *general.Specialization) error {
	// Check for same branch name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, specialization.TenantID, &general.Specialization{},
		repository.Filter("`branch_name`=? AND `degree_id` = ? AND `id`!=?",
			specialization.BranchName, specialization.DegreeID, specialization.ID))
	if err := util.HandleIfExistsError("Branch name exists", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries adds search critera.
func (service *SpecializationService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if _, ok := requestForm["branchName"]; ok {
		util.AddToSlice("`branch_name`", "LIKE ?", "AND", "%"+requestForm.Get("branchName")+"%", &columnNames, &conditions, &operators, &values)
	}
	if degreeID, ok := requestForm["degreeID"]; ok {
		util.AddToSlice("`degree_id`", "= ?", "AND", degreeID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
