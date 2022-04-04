package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	fct "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentFeedbackService provides methods to do different CRUD operations on talent_batch_feedback table.
type TalentFeedbackService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewTalentFeedbackService returns a new instance Of FeedbackService.
func NewTalentFeedbackService(db *gorm.DB, repository repository.Repository) *TalentFeedbackService {
	return &TalentFeedbackService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Question",
			"Option",
		},
	}
}

// AddBatchFeedback will add talent's feedback to table
func (ser *TalentFeedbackService) AddBatchFeedback(feedback *bat.TalentFeedback, uows ...*repository.UnitOfWork) error {

	// check if all foreign keys exist
	err := ser.doesForeignKeyExist(feedback, feedback.CreatedBy)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(ser.DB, false)
	} else {
		uow = uows[0]
	}

	// 	// assign answer to feedback
	// if feedback.OptionID != nil {
	// 	tempFeedbackOption := general.FeedbackOption{}
	// 	err = ser.Repository.GetRecordForTenant(uow, feedback.TenantID, &tempFeedbackOption,
	// 		repository.Filter("`id` = ?", *feedback.OptionID), repository.Select("`value`"))
	// 	if err != nil {
	// 		if length == 0 {
	// 			uow.RollBack()
	// 		}
	// 		return err
	// 	}
	// 	feedback.Answer = tempFeedbackOption.Value
	// }

	err = ser.Repository.Add(uow, feedback)
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

// AddBatchFeedbacks will add talent's feedback for facutly for multiple questions
func (ser *TalentFeedbackService) AddBatchFeedbacks(feedbacks *[]bat.TalentFeedback, tenantID,
	batchID, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(ser.DB, false)
	for _, feedback := range *feedbacks {
		feedback.TenantID = tenantID
		feedback.CreatedBy = credentialID
		feedback.BatchID = batchID

		err := ser.AddBatchFeedback(&feedback, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// DeleteTalentBatchFeedback will delete specified feedback of a faculty in table
func (ser *TalentFeedbackService) DeleteTalentBatchFeedback(feedback *bat.TalentFeedback) error {

	// check if tenant exist
	err := ser.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = ser.doesCredentialExist(feedback.TenantID, feedback.DeletedBy)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = ser.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = ser.doesBatchExist(feedback.TenantID, feedback.BatchID)
	if err != nil {
		return err
	}

	// check if talent exist
	err = ser.doesTalentExist(feedback.TenantID, feedback.TalentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.UpdateWithMap(uow, bat.TalentFeedback{}, map[interface{}]interface{}{
		"DeletedBy": feedback.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=? AND `faculty_id`=? AND `talent_id`=?", feedback.BatchID,
		feedback.FacultyID, feedback.TalentID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllTalentBatchFeeback returns all the feedback for specified batch (admin login)
func (ser *TalentFeedbackService) GetAllTalentBatchFeeback(feedbacks *[]bat.TalentBatchFeedbackDTO,
	tenantID, batchID uuid.UUID, requestFrom url.Values) error {

	// check if tenant exist
	err := ser.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = ser.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, true)
	tempFaculties := []fct.Faculty{}
	err = ser.Repository.GetAll(uow, &tempFaculties,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempFaculty := range tempFaculties {
		tempDTO := bat.TalentBatchFeedbackDTO{}
		tempDTO.Faculty = tempFaculty

		tempTalents := []bat.TalentDTO{}
		err = ser.Repository.GetAll(uow, &tempTalents,
			repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id`=talents.`id`"),
			ser.addSearchQueries(requestFrom), repository.Filter("batch_talents.`batch_id`=?", batchID),
			repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.OrderBy("talents.`first_name`"))
		if err != nil {
			uow.RollBack()
			return err
		}
		for _, tempTalent := range tempTalents {
			tempTalentFeedback := bat.BatchFeedbackDTO{}
			tempTalentFeedback.Talent = tempTalent

			// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempTalentFeedback.BatchFeedbacks, "`created_by`",
			// 	repository.Filter("`batch_id`=? AND `talent_id`=?", batchID, tempTalent.ID),
			// 	repository.PreloadAssociations(ser.associations))
			// if err != nil {
			// 	uow.RollBack()
			// 	return err
			// }

			err = ser.Repository.GetAll(uow, &tempTalentFeedback.BatchFeedbacks,
				repository.Filter("`batch_id`=? AND `talent_id`=?", batchID, tempTalent.ID),
				repository.PreloadAssociations(ser.associations),
				repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= talent_batch_feedback.`question_id`"),
				repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
				repository.Filter("talent_batch_feedback.`tenant_id`=?", tenantID),
				// repository.Filter("feedback_questions.`is_active` = true"),
				repository.OrderBy("feedback_questions.`order`"))
			// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = talent_batch_feedback.`option_id`"))
			if err != nil {
				uow.RollBack()
				return err
			}

			tempDTO.Feedbacks = (append(tempDTO.Feedbacks, tempTalentFeedback))
		}

		(*feedbacks) = append((*feedbacks), tempDTO)
	}

	uow.Commit()
	return nil
}

// GetTalentBatchFeedback returns all the feedback for specified batch // -> TalentFeedback (talent login)
func (ser *TalentFeedbackService) GetTalentBatchFeedback(feedbacks *[]bat.TalentBatchFeedbackDTO,
	tenantID, batchID, talentID uuid.UUID) error {

	// check if tenant exist
	err := ser.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = ser.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// check if talent exist
	err = ser.doesTalentExist(tenantID, talentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, true)
	tempFaculties := []fct.Faculty{}
	err = ser.Repository.GetAll(uow, &tempFaculties,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err

	}

	for _, tempFaculty := range tempFaculties {
		tempDTO := bat.TalentBatchFeedbackDTO{}
		tempDTO.Faculty = tempFaculty

		// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.BatchFeedbacks, "`created_by`",
		// 	repository.Filter("`batch_id`=? AND faculty_id=? AND talent_id=?", batchID, tempFaculty.ID, talentID),
		// 	repository.PreloadAssociations(ser.associations))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = ser.Repository.GetAll(uow, &tempDTO.BatchFeedbacks,
			repository.Filter("`batch_id`=? AND faculty_id=? AND talent_id=?", batchID, tempFaculty.ID, talentID),
			repository.PreloadAssociations(ser.associations),
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= talent_batch_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("talent_batch_feedback.`tenant_id`=?", tenantID),
			// repository.Filter("feedback_questions.`is_active` = true"),
			repository.OrderBy("feedback_questions.`order`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		(*feedbacks) = append((*feedbacks), tempDTO)
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist checks if all foreign keys exist
func (ser *TalentFeedbackService) doesForeignKeyExist(feedback *bat.TalentFeedback, credentialID uuid.UUID) error {

	// check if tenant exist
	err := ser.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = ser.doesCredentialExist(feedback.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = ser.doesBatchExist(feedback.TenantID, feedback.BatchID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = ser.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
	if err != nil {
		return err
	}

	// check if talent exist
	err = ser.doesTalentExist(feedback.TenantID, feedback.TalentID)
	if err != nil {
		return err
	}

	// check if question exist
	err = ser.doesFeedbackQuestionExist(feedback.TenantID, feedback.QuestionID)
	if err != nil {
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (ser *TalentFeedbackService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (ser *TalentFeedbackService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (ser *TalentFeedbackService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ? AND `is_active` = true", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (ser *TalentFeedbackService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.TalentDTO{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (ser *TalentFeedbackService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record for the given tenant.
func (ser *TalentFeedbackService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// addSearchQueries will add any search queries specified in query params
func (ser *TalentFeedbackService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talents.`id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// // DeleteBatchFeedbackForFaculty will delete specified feedback of a faculty in table
// func (ser *TalentFeedbackService) DeleteBatchFeedbackForFaculty(feedback *bat.TalentFeedback) error {

// 	// check if tenant exist
// 	err := ser.doesTenantExist(feedback.TenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if credential exist
// 	err = ser.doesCredentialExist(feedback.TenantID, feedback.DeletedBy)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent_feedback exist
// 	err = ser.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent_feedback exist
// 	err = ser.doesBatchExist(feedback.TenantID, feedback.BatchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, false)

// 	err = ser.Repository.UpdateWithMap(uow, bat.TalentFeedback{}, map[interface{}]interface{}{
// 		"DeletedBy": feedback.DeletedBy,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("batch_id=? AND faculty_id=?", feedback.BatchID, feedback.FacultyID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	uow.Commit()
// 	return nil
// }
