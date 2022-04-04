package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FacultyAssessmentService Provide method to Update, Delete, Add, Get Method For Faculty.
type FacultyAssessmentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewFacultyAssessmentService creates a new instance of FacultyAssessmentService
func NewFacultyAssessmentService(db *gorm.DB, repository repository.Repository) *FacultyAssessmentService {
	return &FacultyAssessmentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Question", "FeedbackQuestionGroup",
			"Faculty", "Option",
			"Credential", "Credential.Role",
		},
	}
}

// AddFacultyAssessment will add assessment of faculty for specified group.
func (service *FacultyAssessmentService) AddFacultyAssessment(feedback *fct.FacultyAssessment, uows ...*repository.UnitOfWork) error {

	// Check if all foreign keys exist.
	err := service.doesForeignKeyExist(feedback, feedback.CreatedBy)
	if err != nil {
		return err
	}

	// Extract feedback ids.
	// feedback.GroupID = feedback.FeedbackQuestionGroup.ID
	feedback.CredentialID = feedback.CreatedBy

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	if feedback.OptionID != nil {
		// Assign answer to feedback.
		tempFeedbackOption := general.FeedbackOption{}
		err = service.Repository.GetRecordForTenant(uow, feedback.TenantID, &tempFeedbackOption,
			repository.Filter("`id` = ?", *feedback.OptionID), repository.Select("`value`"))
		if err != nil {
			if length == 0 {
				uow.RollBack()
			}
			return err
		}

		feedback.Answer = tempFeedbackOption.Value
	}

	// Add feedback.
	err = service.Repository.Add(uow, feedback)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}

	if length == 0 {
		uow.Commit()
	}

	return nil
}

// AddFacultyAssessments will add assessment of faculty for all the groups.
func (service *FacultyAssessmentService) AddFacultyAssessments(feedbacks *[]fct.FacultyAssessment, tenantID,
	credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)

	for index := range *feedbacks {
		(*feedbacks)[index].TenantID = tenantID
		(*feedbacks)[index].CreatedBy = credentialID

		err := service.AddFacultyAssessment(&(*feedbacks)[index], uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()

	return nil
}

// DeleteFacultyAssessment will delete specified assessment of a faculty in the table.
func (service *FacultyAssessmentService) DeleteFacultyAssessment(feedback *fct.FacultyAssessment) error {

	// Check if tenant exist.
	err := service.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(feedback.TenantID, feedback.DeletedBy)
	if err != nil {
		return err
	}

	// Check if faculty exist.
	err = service.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, fct.FacultyAssessment{}, map[string]interface{}{
		"DeletedBy": feedback.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`faculty_id`=?", feedback.FacultyID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()

	return nil
}

// GetFacultyAssessment will get all the assessment for specified faculty
func (service *FacultyAssessmentService) GetFacultyAssessment(feedbacks *[]fct.FacultyAssessmentDTO,
	tenantID, facultyID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if faculty exist.
	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// repository.OrderBy("feedback_question_group.`group_name`")
	// repository.Join("LEFT JOIN credentials ON credentials.`id`=faculty_assessments.`created_by`"+
	// 		" AND credentials.`tenant_id` = faculty_assessments.`tenant_id`"),
	err = service.Repository.GetAll(uow, feedbacks,
		repository.Filter("faculty_assessments.`faculty_id`=?", facultyID),
		repository.PreloadAssociations(service.association),
		repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`=faculty_assessments.`question_id`"+
			" AND feedback_questions.`tenant_id` = faculty_assessments.`tenant_id`"),
		repository.Join("LEFT JOIN feedback_question_groups ON feedback_question_groups.`id`=faculty_assessments.`group_id`"+
			" AND feedback_question_groups.`tenant_id` = faculty_assessments.`tenant_id`"),
		repository.Filter("faculty_assessments.`tenant_id`=?", tenantID), repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("feedback_question_groups.`deleted_at` IS NULL AND feedback_questions.`deleted_at` IS NULL"),
		repository.OrderBy("feedback_question_groups.`order`, feedback_questions.`order`"))
	if err != nil {
		uow.RollBack()

		return err
	}

	uow.Commit()

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist checks if all foreign keys exist.
func (service *FacultyAssessmentService) doesForeignKeyExist(feedback *fct.FacultyAssessment, credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(feedback.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if faculty exist.
	err = service.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
	if err != nil {
		return err
	}

	// check if credential exist.
	// err = service.doesCredentialExist(feedback.TenantID, feedback.CredentialID)
	// if err != nil {
	// 	return err
	// }

	// check if question exist.
	err = service.doesFeedbackQuestionExist(feedback.TenantID, feedback.QuestionID)
	if err != nil {
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (service *FacultyAssessmentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}

	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *FacultyAssessmentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}

	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (service *FacultyAssessmentService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ?", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		return err
	}

	return nil
}

// returns error if there is no talent record for the given tenant.
func (service *FacultyAssessmentService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		return err
	}

	return nil
}
