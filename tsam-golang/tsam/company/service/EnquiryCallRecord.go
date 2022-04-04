package services

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	company "github.com/techlabs/swabhav/tsam/models/company"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// EnquiryCallRecordService provides method to update, delete, add, get all, get one for enquiry call records.
type EnquiryCallRecordService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// enquiryCallRecordAssociations provides preload associations array for enquiry call record.
var enquiryCallRecordAssociations []string = []string{"Purpose", "Outcome"}

// NewEnquiryCallRecordService returns new instance of EnquiryCallRecordService.
func NewEnquiryCallRecordService(db *gorm.DB, repository repository.Repository) *EnquiryCallRecordService {
	return &EnquiryCallRecordService{
		DB:         db,
		Repository: repository,
	}
}

// AddEnquiryCallRecord adds one enquiry call record to database.
func (service *EnquiryCallRecordService) AddEnquiryCallRecord(callRecord *company.CallRecord) error {
	// Get credential id from CreatedBy field of callRecord(set in controller).
	credentialID := callRecord.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(callRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, callRecord.TenantID); err != nil {
		return err
	}

	// Validate enquiry id.
	if err := service.doesCompanyEnquiryExist(callRecord.TenantID, callRecord.EnquiryID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(callRecord); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add enquiry call record to database.
	if err := service.Repository.Add(uow, callRecord); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry call record could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetEnquiryCallRecords gets all enquiry call records from database.
func (service *EnquiryCallRecordService) GetEnquiryCallRecords(enquiryCallRecords *[]company.CallRecordDTO,
	tenantID uuid.UUID, enquiryID uuid.UUID, uows ...*repository.UnitOfWork) error {
	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate enquiry id.
	if err := service.doesCompanyEnquiryExist(tenantID, enquiryID); err != nil {
		return err
	}

	// Get enquiry call records from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, enquiryCallRecords, "date_time",
		repository.Filter("enquiry_id=?", enquiryID),
		repository.PreloadAssociations(enquiryCallRecordAssociations)); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// GetEnquiryCallRecord gets one enquiry call record form database.
func (service *EnquiryCallRecordService) GetEnquiryCallRecord(callRecord *company.CallRecord) error {
	// Validate tenant id.
	if err := service.doesTenantExist(callRecord.TenantID); err != nil {
		return err
	}

	// Validate enquiry id.
	if err := service.doesCompanyEnquiryExist(callRecord.TenantID, callRecord.EnquiryID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get enquiry call record.
	if err := service.Repository.GetForTenant(uow, callRecord.TenantID, callRecord.ID, callRecord,
		repository.PreloadAssociations(enquiryCallRecordAssociations)); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateEnquiryCallRecord updates enquiry call record in Database.
func (service *EnquiryCallRecordService) UpdateEnquiryCallRecord(callRecord *company.CallRecord) error {
	// Get credential id from UpdatedBy field of callRecord(set in controller).
	credentialID := callRecord.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(callRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, callRecord.TenantID); err != nil {
		return err
	}

	// Validate enquiry id.
	if err := service.doesCompanyEnquiryExist(callRecord.TenantID, callRecord.EnquiryID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(callRecord); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting enquiry call record already present in database.
	tempEnquiryCallRecord := company.CallRecord{}

	// Get enquiry call record for getting created_by field of enquiry call record from database.
	if err := service.Repository.GetForTenant(uow, callRecord.TenantID, callRecord.ID, &tempEnquiryCallRecord); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp enquiry call record to enquiry call record to be updated.
	callRecord.CreatedBy = tempEnquiryCallRecord.CreatedBy

	// Update enquiry call record.
	if err := service.Repository.Save(uow, callRecord); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry call reecord could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteEnquiryCallRecord deletes one enquiry call record form database.
func (service *EnquiryCallRecordService) DeleteEnquiryCallRecord(callRecord *company.CallRecord) error {
	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of callRecord(set in controller).
	credentialID := callRecord.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(callRecord.TenantID); err != nil {
		return err
	}

	// Validate enquiry call record id.
	if err := service.doesEnquiryCallRecordExist(callRecord.ID, callRecord.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, callRecord.TenantID); err != nil {
		return err
	}

	// Validate enquiry id.
	if err := service.doesCompanyEnquiryExist(callRecord.TenantID, callRecord.EnquiryID); err != nil {
		return err
	}

	// Update enquiry call record for updating deleted_by and deleted_at field of enquiry call record.
	if err := service.Repository.UpdateWithMap(uow, &company.CallRecord{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("tenant_id=? AND id=?", callRecord.TenantID, callRecord.ID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Enquiry call record could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *EnquiryCallRecordService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *EnquiryCallRecordService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if credential(parent credential) exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyEnquiryExist validates if company exists or not in database.
func (service *EnquiryCallRecordService) doesCompanyEnquiryExist(tenantID uuid.UUID, enquiryID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent company qnuiry exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, company.Enquiry{},
		repository.Filter("`id` = ?", enquiryID))
	if err := util.HandleError("Invalid company enquiry ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesEnquiryCallRecordExist validates if enquiry call record exists or not in database.
func (service *EnquiryCallRecordService) doesEnquiryCallRecordExist(enquiryCallRecordID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check enquiry call record exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, company.CallRecord{},
		repository.Filter("`id` = ?", enquiryCallRecordID))
	if err := util.HandleError("Invalid enquiry call record ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doForeignKeysExist validates if purpose ad outcome exists or not in database.
func (service *EnquiryCallRecordService) doForeignKeysExist(callRecord *company.CallRecord) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if purpose exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Purpose{}, repository.Filter("`id` = ?", callRecord.PurposeID))
	if err := util.HandleError("Invalid purpose ID", exists, err); err != nil {
		return err
	}

	// Check if outcome exists or not.
	exists, err = repository.DoesRecordExist(uow.DB, general.Outcome{}, repository.Filter("`id` = ?", callRecord.OutcomeID))
	if err := util.HandleError("Invalid outcome ID", exists, err); err != nil {
		return err
	}
	return nil
}
