package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	general "github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	talenq "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// WaitingListService provides method to update, delete, add, get for waiting list.
type WaitingListService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// WaitingListAssociationNames provides preload associations array for waiting list.
var WaitingListAssociationNames []string = []string{
	"CompanyBranch", "CompanyRequirement", "Course", "Batch", "CompanyRequirement.Branch",
}

// NewWaitingListService returns new instance of WaitingListService.
func NewWaitingListService(db *gorm.DB, repository repository.Repository) *WaitingListService {
	return &WaitingListService{
		DB:         db,
		Repository: repository,
	}
}

// AddWaitingList adds one waiting list to database for talent.
func (service *WaitingListService) AddWaitingList(waitingList *tal.WaitingList) error {
	// Get credential id from CreatedBy field of waitingList(set in controller).
	credentialID := waitingList.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(waitingList.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, waitingList.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(waitingList, false); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// If waiting list is being added for enquiry.
	if waitingList.TalentID == nil {
		// Check id talent with same email exists or not.
		tempTalent := tal.Talent{}
		if err := service.Repository.GetRecordForTenant(uow, waitingList.TenantID, &tempTalent,
			repository.Filter("`email`=?", waitingList.Email),
			repository.Select("id, is_active")); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}

		}

		if tempTalent.ID != uuid.Nil {
			waitingList.TalentID = &tempTalent.ID
		}
	}

	// If company requirement id exists then check for same email and requirement id confict.
	if waitingList.CompanyRequirementID != nil {
		exists, err := service.doesEmailAndRequirementExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and company requirement exists")
		}
	}

	// If batch id exists then check for same email and batch id confict.
	if waitingList.BatchID != nil {
		exists, err := service.doesEmailAndBatchExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and batch exists")
		}
	}

	// Add waiting list to database.
	if err := service.Repository.Add(uow, waitingList); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Waiting List could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// AddWaitingListForApplicationForm adds one waiting list to database for enquiry.
func (service *WaitingListService) AddWaitingListForApplicationForm(waitingList *tal.WaitingList, enquiryID,
	tenantID uuid.UUID, email string, uows ...*repository.UnitOfWork) error {

	// // Get credential id from CreatedBy field of waitingList(set in controller).
	// credentialID := waitingList.CreatedBy

	// Set tenant id for waiting list.
	waitingList.TenantID = tenantID

	// Set enquiry id for waiting list.
	waitingList.EnquiryID = &enquiryID

	// Set enquiry email for waiting list.
	waitingList.Email = &email

	// // Validate tenant id.
	// if err := service.doesTenantExist(waitingList.TenantID); err != nil {
	// 	return err
	// }

	// // Validate credential id.
	// if err := service.doesCredentialExist(credentialID, waitingList.TenantID); err != nil {
	// 	return err
	// }

	// Validate foreign keys.
	if err := service.doForeignKeysExist(waitingList, true); err != nil {
		return err
	}

	// Starting new transaction.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Check id talent with same email exists or not.
	tempTalent := tal.Talent{}
	if err := service.Repository.GetRecordForTenant(uow, tenantID, &tempTalent,
		repository.Filter("`email`=?", email),
		repository.Select("id, is_active")); err != nil {
		if err != gorm.ErrRecordNotFound {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}
	}

	if tempTalent.ID != uuid.Nil {
		waitingList.TalentID = &tempTalent.ID
	}

	// If company requirement id exists then check for same email and requirement id confict.
	if waitingList.CompanyRequirementID != nil {
		exists, err := service.doesEmailAndRequirementExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and company requirement exists")
		}
	}

	// If batch id exists then check for same email and batch id confict.
	if waitingList.BatchID != nil {
		exists, err := service.doesEmailAndBatchExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and batch exists")
		}
	}

	// Add waiting list to database.
	if err := service.Repository.Add(uow, waitingList); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Waiting List could not be added", http.StatusInternalServerError)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// GetWaitingLists gets all waiting lists by talent id form database.
func (service *WaitingListService) GetWaitingLists(waitingList *[]tal.WaitingListDTO, tenantID uuid.UUID,
	talentID, enquiryID *uuid.UUID, form url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate talent id.
	if talentID != nil {
		if err := service.doesTalentExist(*talentID, tenantID); err != nil {
			return err
		}
	}

	// Validate enquiry id.
	if enquiryID != nil {
		if err := service.doesEnquiryExist(*enquiryID, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// If talent id is present get waiting lits for talent.
	if talentID != nil {
		// Get waiting lists.
		if err := service.Repository.GetAllForTenant(uow, tenantID, waitingList,
			repository.Filter("`talent_id`=?", talentID),
			service.addSearchQueries(form),
			repository.PreloadAssociations(WaitingListAssociationNames)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}
	}

	// If enquiry id is present get waiting lits for enquiry.
	if enquiryID != nil {
		// Get waiting lists.
		if err := service.Repository.GetAllForTenant(uow, tenantID, waitingList,
			repository.Filter("`enquiry_id`=?", enquiryID),
			repository.PreloadAssociations(WaitingListAssociationNames)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}
	}

	uow.Commit()
	return nil
}

// GetTwoWaitingLists gets two waiting lists form database.
func (service *WaitingListService) GetTwoWaitingLists(twoWaitingLists *tal.TwoWaitingLists, tenantID uuid.UUID,
	form url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create waiting list for talent.
	talentWaitingList := &[]tal.WaitingList{}

	// Get waiting list for talent.
	if err := service.Repository.GetAllForTenant(uow, tenantID, talentWaitingList,
		repository.Filter("`talent_id` IS NOT NULL"),
		service.addSearchQueries(form),
	); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create waiting list for enquiry.
	enauiryWaitingList := &[]tal.WaitingList{}

	// Get waiting list for enquiry.
	if err := service.Repository.GetAllForTenant(uow, tenantID, enauiryWaitingList,
		repository.Filter("`enquiry_id` IS NOT NULL AND `talent_id` IS NULL"),
		service.addSearchQueries(form),
	); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give both the list to two waiting lists.
	twoWaitingLists.TalentWaitingList = talentWaitingList
	twoWaitingLists.EnquiryWaitingList = enauiryWaitingList

	uow.Commit()
	return nil
}

// UpdateWaitingList updates waiting list in database.
func (service *WaitingListService) UpdateWaitingList(waitingList *tal.WaitingList) error {

	// Get credential id from UpdatedBy field of waitingList(set in controller).
	credentialID := waitingList.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(waitingList.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, waitingList.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(waitingList, false); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// If company requirement id exists then check for same email and requirement id confict.
	if waitingList.CompanyRequirementID != nil {
		exists, err := service.doesEmailAndRequirementExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and company requirement exists")
		}
	}

	// If batch id exists then check for same email and batch id confict.
	if waitingList.BatchID != nil {
		exists, err := service.doesEmailAndBatchExist(uow, waitingList)
		if err != nil {
			uow.RollBack()
			return err
		}
		if exists {
			uow.RollBack()
			return errors.NewValidationError("Waiting list entry with same email and batch exists")
		}
	}

	// Create bucket for getting waiting list already present in database.
	tempWaitingList := tal.WaitingList{}

	// Get waiting list for getting created_by field of waiting list from database.
	if err := service.Repository.GetForTenant(uow, waitingList.TenantID, waitingList.ID, &tempWaitingList); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp waiting list to waiting list to be updated.
	waitingList.CreatedBy = tempWaitingList.CreatedBy

	// Give talent id from temp waiting list to waiting list to be updated.
	waitingList.TalentID = tempWaitingList.TalentID

	// Give enquiry id from temp waiting list to waiting list to be updated.
	waitingList.EnquiryID = tempWaitingList.EnquiryID

	// Update waiting list.
	if err := service.Repository.Save(uow, waitingList); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Waiting list could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// UpdateWaitingList updates waiting list in database.
func (service *WaitingListService) TransferWaitingList(updateWaitingList *tal.UpdateWaitingList, tenantID,
	credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, tenantID); err != nil {
		return err
	}

	// Validate company branch id.
	if updateWaitingList.CompanyBranchID != nil {
		if err := service.doesCompanyBranchExist(*updateWaitingList.CompanyBranchID, tenantID); err != nil {
			return err
		}
	}

	// Validate company requirement id.
	if updateWaitingList.RequirementID != nil {
		if err := service.doesCompanyRequirementExist(*updateWaitingList.RequirementID, tenantID); err != nil {
			return err
		}
	}

	// Validate course id.
	if updateWaitingList.CourseID != nil {
		if err := service.doesCourseExist(*updateWaitingList.CourseID, tenantID); err != nil {
			return err
		}
	}

	// Validate batch id.
	if updateWaitingList.BatchID != nil {
		if err := service.doesBatchExist(*updateWaitingList.BatchID, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Collect emails of all waiting lists entries to be transfered.
	var waitingListEmailsToBeReansfered []string

	for _, waitingList := range updateWaitingList.WaitingLists {
		waitingListEmailsToBeReansfered = append(waitingListEmailsToBeReansfered, *waitingList.Email)
	}

	// Collect IDs of all waiting list entries to be transfered
	var waitingListIDsToBeTransfered []uuid.UUID

	for _, temp := range updateWaitingList.WaitingLists {
		waitingListIDsToBeTransfered = append(waitingListIDsToBeTransfered, temp.ID)
	}

	// Create bucket for getting waiting list entries that are already present in batch or requirement to be transfered to.
	waitingListsAssigned := []tal.WaitingList{}

	// If transferring to batch check batch id.
	if updateWaitingList.BatchID != nil {
		// Get all waiting list entries that already have the batch id assigned to them.
		if err := service.Repository.GetAllForTenant(uow, tenantID, &waitingListsAssigned,
			repository.Filter("`email` IN (?)", waitingListEmailsToBeReansfered),
			repository.Filter("`batch_id` =?", updateWaitingList.BatchID)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

	}

	// If transferring to batch check requirement id.
	if updateWaitingList.RequirementID != nil {
		// Get all waiting list entries that already have the requirement id assigned to them.
		if err := service.Repository.GetAllForTenant(uow, tenantID, &waitingListsAssigned,
			repository.Filter("`email` IN (?)", waitingListEmailsToBeReansfered),
			repository.Filter("`company_requirement_id` =?", updateWaitingList.RequirementID)); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

	}

	// Collect emails of all waiting lists entries that are already assigned to the batch or requirement.
	var waitingListEmailsAlreadyAssigned []string

	// If there are waiting list entries that are already assigned then collect their emails.
	if len(waitingListsAssigned) != 0 {
		for _, temp := range waitingListsAssigned {
			waitingListEmailsAlreadyAssigned = append(waitingListEmailsAlreadyAssigned, *temp.Email)
		}
	}

	// fmt.Println("emails already present!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// for _, temp := range waitingListEmailsAlreadyAssigned {
	// 	fmt.Println(temp)
	// }

	// Create query processors for filters.
	var queryProcessors []repository.QueryProcessor

	// If waiting list entries already present then do not give the email filter.
	if len(waitingListsAssigned) != 0 {
		queryProcessors = append(queryProcessors, repository.Filter("`email` NOT IN (?)", waitingListEmailsAlreadyAssigned))
	}
	queryProcessors = append(queryProcessors, repository.Filter("`id` IN (?)", waitingListIDsToBeTransfered))

	// Update the waiting list entries' company branch id, company requirement id , course id and batch id.
	if err := service.Repository.UpdateWithMap(uow, &tal.WaitingList{}, map[string]interface{}{
		"UpdatedBy":            credentialID,
		"CompanyBranchID":      updateWaitingList.CompanyBranchID,
		"CompanyRequirementID": updateWaitingList.RequirementID,
		"CourseID":             updateWaitingList.CourseID,
		"BatchID":              updateWaitingList.BatchID,
		"IsActive":             true,
	},
		queryProcessors...,
	); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Waiting list transfered", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteWaitingList deletes one waiting list form database.
func (service *WaitingListService) DeleteWaitingList(waitingList *tal.WaitingList) error {

	// Get credential id from DeletedBy field of waitingList(set in controller).
	credentialID := waitingList.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(waitingList.TenantID); err != nil {
		return err
	}

	// Validate waiting list id.
	if err := service.doesWaitingListExist(waitingList.ID, waitingList.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, waitingList.TenantID); err != nil {
		return err
	}

	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update waiting list for updating deleted_by and deleted_at field of waiting list.
	if err := service.Repository.UpdateWithMap(uow, &tal.WaitingList{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", waitingList.TenantID, waitingList.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Waiting list could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *WaitingListService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *WaitingListService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesWaitingListExist validates if waiting list exists or not in database.
func (service *WaitingListService) doesWaitingListExist(waitingListID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.WaitingList{},
		repository.Filter("`id` = ?", waitingListID))
	if err := util.HandleError("Invalid waiting list ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *WaitingListService) doesTalentExist(talentID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesEnquiryExist validates if enquiry exists or not in database.
func (service *WaitingListService) doesEnquiryExist(enquiryID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talenq.Enquiry{},
		repository.Filter("`id` = ?", enquiryID))
	if err := util.HandleError("Invalid enquiry ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCompanyBranchExist validates if company branch exists or not in database.
func (service *WaitingListService) doesCompanyBranchExist(companyBranchID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Branch{},
		repository.Filter("`id` = ?", companyBranchID))
	if err := util.HandleError("Invalid company branch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCompanyRequirementExist validates if company requirement exists or not in database.
func (service *WaitingListService) doesCompanyRequirementExist(companyRequirementID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Requirement{},
		repository.Filter("`id` = ?", companyRequirementID))
	if err := util.HandleError("Invalid company requirement ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchExist validates if batch exists or not in database.
func (service *WaitingListService) doesBatchExist(batchID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCourseExist validates if batch exists or not in database.
func (service *WaitingListService) doesCourseExist(courseID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
		repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesEmailAndRequirementExist check for same email and compnay requirement id conflict, if conflict occurs return true.
func (service *WaitingListService) doesEmailAndRequirementExist(uow *repository.UnitOfWork, waitingList *tal.WaitingList) (bool, error) {
	var count int
	if err := service.Repository.GetCountForTenant(uow, waitingList.TenantID, &tal.WaitingList{}, &count,
		repository.Filter("`email`=? AND `company_requirement_id`=? AND `id`!=? AND `is_active`=?", waitingList.Email, waitingList.CompanyRequirementID, waitingList.ID, true)); err != nil {
		log.NewLogger().Error(err.Error())
		return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if count != 0 { //email already present
		log.NewLogger().Error("Waiting list entry with same email and company requirement exists")
		return true, nil
	}

	return false, nil
}

// doesEmailAndBatchExist check for same email and batch id conflict, if conflict occurs return true.
func (service *WaitingListService) doesEmailAndBatchExist(uow *repository.UnitOfWork, waitingList *tal.WaitingList) (bool, error) {
	var count int
	if err := service.Repository.GetCountForTenant(uow, waitingList.TenantID, &tal.WaitingList{}, &count,
		repository.Filter("`email`=? AND `batch_id`=? AND `id`!=? AND `is_active`=?", waitingList.Email, waitingList.BatchID, waitingList.ID, true)); err != nil {
		log.NewLogger().Error(err.Error())
		return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if count != 0 { //email already present
		log.NewLogger().Error("Waiting list entry with same email and batch exists")
		return true, nil
	}

	return false, nil
}

// doForeignKeysExist validates if purpose ad outcome exists or not in database.
func (service *WaitingListService) doForeignKeysExist(waitingList *tal.WaitingList, isEnquiry bool) error {
	// Check if talent exists or not.
	if waitingList.TalentID != nil {
		if err := service.doesTalentExist(*waitingList.TalentID, waitingList.TenantID); err != nil {
			return err
		}
	}

	// Check if enquiry exists or not.
	if waitingList.EnquiryID != nil && !isEnquiry {
		if err := service.doesEnquiryExist(*waitingList.EnquiryID, waitingList.TenantID); err != nil {
			return err
		}
	}

	// Check if company branch exists or not.
	if waitingList.CompanyBranchID != nil {
		if err := service.doesCompanyBranchExist(*waitingList.CompanyBranchID, waitingList.TenantID); err != nil {
			return err
		}
	}

	// Check if company requirement exists or not.
	if waitingList.CompanyRequirementID != nil {
		if err := service.doesCompanyRequirementExist(*waitingList.CompanyRequirementID, waitingList.TenantID); err != nil {
			return err
		}
	}

	// Check if course exists or not.
	if waitingList.CourseID != nil {
		if err := service.doesCourseExist(*waitingList.CourseID, waitingList.TenantID); err != nil {
			return err
		}
	}

	// Check if batch exists or not.
	if waitingList.BatchID != nil {
		if err := service.doesBatchExist(*waitingList.BatchID, waitingList.TenantID); err != nil {
			return err
		}
	}

	return nil
}

// addSearchQueries adds search queries.
func (service *WaitingListService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// Batch ID.
	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("`batch_id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}

	// Course ID.
	if courseID, ok := requestForm["courseID"]; ok {
		util.AddToSlice("`course_id`", "= ?", "AND", courseID, &columnNames, &conditions, &operators, &values)
	}

	// Company requirement ID.
	if companyRequirementID, ok := requestForm["companyRequirementID"]; ok {
		util.AddToSlice("`company_requirement_id`", "= ?", "AND", companyRequirementID, &columnNames, &conditions, &operators, &values)
	}

	// Is active.
	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
