package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	general "github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentNextActionService provides method to update, delete, add, get all, get one for talent next actions.
type TalentNextActionService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewTalentNextActionService returns new instance of TalentNextActionService.
func NewTalentNextActionService(db *gorm.DB, repository repository.Repository) *TalentNextActionService {
	return &TalentNextActionService{
		DB:           db,
		Repository:   repository,
		associations: []string{"ActionType", "Courses", "Companies", "Technologies"},
	}
}

// AddTalentNextAction adds one action on talent to database.
func (service *TalentNextActionService) AddTalentNextAction(talentNextAction *talent.NextAction) error {
	// Get credential id from CreatedBy field of talentNextAction(set in controller).
	credentialID := talentNextAction.CreatedBy

	// Checks all the foreign keys in next action.
	err := service.doForeignKeysExist(credentialID, talentNextAction)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add talent next action to database.
	if err := service.Repository.Add(uow, talentNextAction); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent next action could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetAllTalentNextActions gets all talent next actions from database.
func (service *TalentNextActionService) GetAllTalentNextActions(talentNextActions *[]*talent.NextActionDTO,
	tenantID uuid.UUID, talentID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if talent exists.
	if err := service.doesTalentExist(tenantID, talentID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Get talent next actions from database.
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, talentNextActions, "`target_date`",
		repository.Filter("`talent_id`=?", talentID),
		repository.PreloadAssociations(service.associations)); err != nil {
		// No rollback since it is read-only.
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// // temp change as date does not get append in form
	// for index := range *talentNextActions {
	// 	(*talentNextActions)[index].FromDate = util.RemoveTimeStampFromDate(*(*talentNextActions)[index].FromDate)
	// 	(*talentNextActions)[index].ToDate = util.RemoveTimeStampFromDate(*(*talentNextActions)[index].ToDate)
	// 	(*talentNextActions)[index].TargetDate = util.RemoveTimeStampFromDate(*(*talentNextActions)[index].TargetDate)
	// }
	// No commit since it is read-only
	return nil
}

// GetTalentNextAction gets one talent_next_action record form database.
func (service *TalentNextActionService) GetTalentNextAction(talentNextAction *talent.NextAction) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(talentNextAction.TenantID); err != nil {
		return err
	}

	// Check if talent exists.
	if err := service.doesTalentExist(talentNextAction.TenantID, talentNextAction.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent next action.
	if err := service.Repository.GetForTenant(uow, talentNextAction.TenantID, talentNextAction.ID, talentNextAction,
		repository.PreloadAssociations(service.associations)); err != nil {
		// No rollback since it is read-only.
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}
	// No commit since it is read-only.
	return nil
}

// UpdateTalentNextAction updates talent next action in Database.
func (service *TalentNextActionService) UpdateTalentNextAction(talentNextAction *talent.NextAction) error {

	// Checks all the foreign keys in next action.
	err := service.doForeignKeysExist(talentNextAction.UpdatedBy, talentNextAction)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Check if talent next action record exist.
	err = service.doesTalentNextActionExist(talentNextAction.TenantID, talentNextAction.ID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting talent_next_action record already present in database.
	talentNextActionInDB := &talent.NextAction{}

	// Get talent next action for getting created_by field of talent next action from database.
	err = service.Repository.GetForTenant(uow, talentNextAction.TenantID, talentNextAction.ID,
		talentNextActionInDB,
		repository.Select("`created_by`"))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Next Action record not found")
	}

	// Assign existing created by value else save will update it.
	talentNextAction.CreatedBy = talentNextActionInDB.CreatedBy

	// Update target community associations.
	if err := service.updateTalentNextActionAssociation(uow, talentNextAction); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Target Community could not be updated", http.StatusInternalServerError)
	}

	// Update talent next action.
	err = service.Repository.Save(uow, talentNextAction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent next action could not be updated", http.StatusInternalServerError)
	}
	uow.Commit()
	return nil
}

// DeleteTalentNextAction deletes one talent_next_action record form database.
func (service *TalentNextActionService) DeleteTalentNextAction(talentNextAction *talent.NextAction) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(talentNextAction.TenantID); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Checks if the credential record exists.
	if err := service.doesCredentialExist(talentNextAction.TenantID, talentNextAction.DeletedBy); err != nil {
		return err
	}

	// Checks if next_action_record exists.
	if err := service.doesTalentNextActionExist(talentNextAction.TenantID, talentNextAction.ID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err := service.Repository.UpdateWithMap(uow, &talent.NextAction{}, map[interface{}]interface{}{
		"DeletedBy": talentNextAction.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id`=?", talentNextAction.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TalentNextActionService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesNextActionTypeExist validates if next action type exists or not in database.
func (service *TalentNextActionService) doesNextActionTypeExist(tenantID, nextActionTypeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.NextActionType{},
		repository.Filter("`id`=?", nextActionTypeID))
	if err := util.HandleError("Invalid next action type ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *TalentNextActionService) doesCredentialExist(tenantID uuid.UUID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *TalentNextActionService) doesTalentExist(tenantID uuid.UUID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id`=?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentNextActionExist validates if talent next action exists or not in database.
func (service *TalentNextActionService) doesTalentNextActionExist(tenantID uuid.UUID, talentNextActionID uuid.UUID) error {
	//check talent next action exists or not
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.NextAction{},
		repository.Filter("`id`=?", talentNextActionID))
	if err := util.HandleError("Invalid talent next action ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if talent next action exists or not in database.
func (service *TalentNextActionService) doForeignKeysExist(credentialID uuid.UUID, nextAction *talent.NextAction) error {
	tenantID := nextAction.TenantID
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if next action type exists.
	if err := service.doesNextActionTypeExist(tenantID, nextAction.ActionTypeID); err != nil {
		return err
	}

	// Check if compulsary fields exist based on action type.
	// Get action type by action id from database.
	nextActionType := talent.NextActionType{}
	err := service.Repository.GetForTenant(uow, tenantID, nextAction.ActionTypeID, &nextActionType)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Invalid next action type ID")
	}

	// Action type is Blog.
	if nextActionType.Type == "Blog" {

		// Technologies.
		if nextAction.Technologies == nil || (nextAction.Technologies != nil && len(nextAction.Technologies) == 0) {
			return errors.NewValidationError("Technologies must be specified")
		}

		// Target Date.
		if nextAction.TargetDate == nil || (nextAction.TargetDate != nil && util.IsEmpty(*nextAction.TargetDate)) {
			return errors.NewValidationError("Target Date must be specified")
		}
	}

	// Action type is Referral.
	if nextActionType.Type == "Referral" {

		// Referral Count.
		if nextAction.ReferralCount == nil || (nextAction.ReferralCount != nil && *nextAction.ReferralCount == 0) {
			return errors.NewValidationError("Referral count must be specified")
		}
	}

	// Action type is Teaching Assistant.
	if nextActionType.Type == "Teaching Assistant" {

		// Technologies.
		if nextAction.Technologies == nil || (nextAction.Technologies != nil && len(nextAction.Technologies) == 0) {
			return errors.NewValidationError("Technologies must be specified")
		}

		// Stipend.
		if nextAction.Stipend == nil || (nextAction.Stipend != nil && *nextAction.Stipend == 0) {
			return errors.NewValidationError("Stipend must be specified")
		}

		// Target Date.
		if nextAction.TargetDate == nil || (nextAction.TargetDate != nil && util.IsEmpty(*nextAction.TargetDate)) {
			return errors.NewValidationError("Target Date must be specified")
		}
	}

	// Action type is Teaching Assistant.
	if nextActionType.Type == "Teaching Assistant" {

		// Courses.
		if nextAction.Courses == nil || (nextAction.Courses != nil && len(nextAction.Courses) == 0) {
			return errors.NewValidationError("Courses must be specified")
		}

		// Target Date.
		if nextAction.TargetDate == nil || (nextAction.TargetDate != nil && util.IsEmpty(*nextAction.TargetDate)) {
			return errors.NewValidationError("Target Date must be specified")
		}
	}

	// Action type is Internship.
	if nextActionType.Type == "Internship" {

		// Technologies.
		if nextAction.Technologies == nil || (nextAction.Technologies != nil && len(nextAction.Technologies) == 0) {
			return errors.NewValidationError("Technologies must be specified")
		}

		// Stipend.
		if nextAction.Stipend == nil || (nextAction.Stipend != nil && *nextAction.Stipend == 0) {
			return errors.NewValidationError("Stipend must be specified")
		}

		// From Date.
		if nextAction.FromDate == nil || (nextAction.FromDate != nil && util.IsEmpty(*nextAction.FromDate)) {
			return errors.NewValidationError("From Date must be specified")
		}

		// To Date.
		if nextAction.ToDate == nil || (nextAction.ToDate != nil && util.IsEmpty(*nextAction.ToDate)) {
			return errors.NewValidationError("To Date must be specified")
		}

		// Target Date.
		if nextAction.TargetDate == nil || (nextAction.TargetDate != nil && util.IsEmpty(*nextAction.TargetDate)) {
			return errors.NewValidationError("Target Date must be specified")
		}
	}

	// Action type is Placement.
	if nextActionType.Type == "Placement" {

		// Companies.
		if nextAction.Companies == nil || (nextAction.Companies != nil && len(nextAction.Companies) == 0) {
			return errors.NewValidationError("Companies must be specified")
		}

		// Technologies.
		if nextAction.Technologies == nil || (nextAction.Technologies != nil && len(nextAction.Technologies) == 0) {
			return errors.NewValidationError("Technologies must be specified")
		}

		// Stipend.
		if nextAction.Stipend == nil || (nextAction.Stipend != nil && *nextAction.Stipend == 0) {
			return errors.NewValidationError("Stipend must be specified")
		}

		// Target Date.
		if nextAction.TargetDate == nil || (nextAction.TargetDate != nil && util.IsEmpty(*nextAction.TargetDate)) {
			return errors.NewValidationError("Target Date must be specified")
		}
	}

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credentials exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if talent exists.
	if err := service.doesTalentExist(tenantID, nextAction.TalentID); err != nil {
		return err
	}

	// Check if courses exist or not.
	if nextAction.Courses != nil && len(nextAction.Courses) > 0 {
		var courseIDs []uuid.UUID
		for _, tempCourse := range nextAction.Courses {
			courseIDs = append(courseIDs, tempCourse.ID)
		}
		// Get count for courseIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, course.Course{}, &count,
			repository.Filter("`id` IN (?)", courseIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(nextAction.Courses) {
			log.NewLogger().Error("Course ID is invalid")
			return errors.NewValidationError("Course ID is invalid")
		}
	}

	// Check if companies exist or not.
	if nextAction.Companies != nil && len(nextAction.Companies) > 0 {
		var companyIDs []uuid.UUID
		for _, tempCompany := range nextAction.Companies {
			companyIDs = append(companyIDs, tempCompany.ID)
		}
		// Get count for companyIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, company.Company{}, &count,
			repository.Filter("`id` IN (?)", companyIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(nextAction.Companies) {
			log.NewLogger().Error("Company ID is invalid")
			return errors.NewValidationError("Company ID is invalid")
		}
	}

	// Check if tecchnologies exist or not.
	if nextAction.Technologies != nil && len(nextAction.Technologies) > 0 {
		var technologyIDs []uuid.UUID
		for _, tempTechnology := range nextAction.Technologies {
			technologyIDs = append(technologyIDs, tempTechnology.ID)
		}
		// Get count for technologyIDs.
		var count int = 0
		err := service.Repository.GetCountForTenant(uow, tenantID, general.Technology{}, &count,
			repository.Filter("`id` IN (?)", technologyIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(nextAction.Technologies) {
			log.NewLogger().Error("Technology ID is invalid")
			return errors.NewValidationError("Technology ID is invalid")
		}
	}

	return nil
}

// updateTalentNextActionAssociation updates talent next action's associations.
func (service *TalentNextActionService) updateTalentNextActionAssociation(uow *repository.UnitOfWork, talentNextAction *talent.NextAction) error {
	// Replace courses of talent next action.
	if err := service.Repository.ReplaceAssociations(uow, talentNextAction, "Courses",
		talentNextAction.Courses); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace company branches of talent next action.
	if err := service.Repository.ReplaceAssociations(uow, talentNextAction, "Companies",
		talentNextAction.Companies); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Replace technologies of talent next action.
	if err := service.Repository.ReplaceAssociations(uow, talentNextAction, "Technologies",
		talentNextAction.Technologies); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// Make all association map nil to avoid inserts or updates while updating talent next action.
	talentNextAction.Courses = nil
	talentNextAction.Companies = nil
	talentNextAction.Technologies = nil

	return nil
}
