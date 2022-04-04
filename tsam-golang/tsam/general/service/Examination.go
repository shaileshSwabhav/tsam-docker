package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ExaminationService provides methods to update, delete, add, get method for examination
type ExaminationService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewExaminationService returns new instance of ExaminationService.
func NewExaminationService(db *gorm.DB, repository repository.Repository) *ExaminationService {
	return &ExaminationService{
		DB:         db,
		Repository: repository,
	}
}

// AddExamination adds new examination to database.
func (service *ExaminationService) AddExamination(examination *general.Examination, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(examination.TenantID)
	if err != nil {
		return err
	}

	// Validate if same examination name exists.
	err = service.doesExaminationNameExist(examination)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(examination.CreatedBy, examination.TenantID)
	if err != nil {
		return err
	}

	// Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add examination.
	err = service.Repository.Add(uow, examination)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddExaminations adds multiple examinations to database.
func (service *ExaminationService) AddExaminations(examinations *[]general.Examination, examinationIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	//check for same name conflict
	for i := 0; i < len(*examinations); i++ {
		for j := 0; j < len(*examinations); j++ {
			if i != j && (*examinations)[i].Name == (*examinations)[j].Name {
				log.NewLogger().Error("Name:" + (*examinations)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*examinations)[j].Name + " exists")
			}
		}
	}

	// Add individual examination.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, examination := range *examinations {
		examination.TenantID = tenantID
		examination.CreatedBy = credentialID
		err := service.AddExamination(&examination, uow)
		if err != nil {
			return err
		}
		*examinationIDs = append(*examinationIDs, examination.ID)
	}

	uow.Commit()
	return nil
}

// UpdateExamination updates examination to database.
func (service *ExaminationService) UpdateExamination(examination *general.Examination) error {
	// Validate tenant ID.
	err := service.doesTenantExist(examination.TenantID)
	if err != nil {
		return err
	}

	// Validate examination ID.
	err = service.doeExaminationExist(examination.ID, examination.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(examination.UpdatedBy, examination.TenantID)
	if err != nil {
		return err
	}

	// Validate if same examination name exists.
	err = service.doesExaminationNameExist(examination)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update examination.
	err = service.Repository.Update(uow, examination)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteExamination deletes examination from database.
func (service *ExaminationService) DeleteExamination(examination *general.Examination) error {
	credentialID := examination.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(examination.TenantID)
	if err != nil {
		return err
	}

	// Validate examination ID.
	err = service.doeExaminationExist(examination.ID, examination.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, examination.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update examination for updating deleted_by and deleted_at fields of examination
	if err := service.Repository.UpdateWithMap(uow, examination, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", examination.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Examination could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetExamination returns one examination by specific id.
func (service *ExaminationService) GetExamination(examination *general.Examination) error {
	// Validate tenant id.
	err := service.doesTenantExist(examination.TenantID)
	if err != nil {
		return err
	}

	// Validate examination id.
	err = service.doeExaminationExist(examination.ID, examination.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one examination by examination id from database.
	err = service.Repository.GetForTenant(uow, examination.TenantID, examination.ID, examination)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetExaminations returns all examinations.
func (service *ExaminationService) GetExaminations(examinations *[]general.Examination, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all examinations from database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, examinations, "`name`")
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// doesExaminationNameExist returns true if the examination name already exists for the state in database.
func (service *ExaminationService) doesExaminationNameExist(examination *general.Examination) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, examination.TenantID, &general.Examination{},
		repository.Filter("`name`=? AND `id`!=?", examination.Name, examination.ID))
	if err := util.HandleIfExistsError("Name:"+examination.Name+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ExaminationService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doeExaminationExist validates if examination exists or not in database.
func (service *ExaminationService) doeExaminationExist(examinationID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Examination{},
		repository.Filter("`id` = ?", examinationID))
	if err := util.HandleError("Invalid examination ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if state exists or not in database.
func (service *ExaminationService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
