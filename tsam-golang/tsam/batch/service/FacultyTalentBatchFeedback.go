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

// FacultyFeedbackService provides methods to do different CRUD operations on faculty_talent_batch_feedback table.
type FacultyFeedbackService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewFacultyFeedbackService returns a new instance Of FeedbackService.
func NewFacultyFeedbackService(db *gorm.DB, repository repository.Repository) *FacultyFeedbackService {
	return &FacultyFeedbackService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Question",
			"Option",
			// "Talent",
		},
	}
}

// AddBatchFeedback will add faculties feedback to table
func (ser *FacultyFeedbackService) AddBatchFeedback(feedback *bat.FacultyTalentFeedback, uows ...*repository.UnitOfWork) error {

	// check if all foreign keys exist
	err := ser.doesForeignKeyExist(feedback, feedback.CreatedBy)
	if err != nil {
		return err
	}

	// check if talent exist
	err = ser.doesTalentExist(feedback.TenantID, feedback.TalentID)
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

// AddBatchFeedbacks will add faculties feedback for multiple questions
func (ser *FacultyFeedbackService) AddBatchFeedbacks(feedbacks *[]bat.FacultyTalentFeedback, tenantID,
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

// DeleteFacultyTalentBatchFeedback will delete specified feedback of a talent in table
func (ser *FacultyFeedbackService) DeleteFacultyTalentBatchFeedback(feedback *bat.FacultyTalentFeedback) error {

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

	// check if talent_feedback exist
	err = ser.doesTalentExist(feedback.TenantID, feedback.TalentID)
	if err != nil {
		return err
	}

	// check if talent_feedback exist
	err = ser.doesBatchExist(feedback.TenantID, feedback.BatchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.UpdateWithMap(uow, bat.FacultyTalentFeedback{}, map[interface{}]interface{}{
		"DeletedBy": feedback.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=? AND `talent_id`=?", feedback.BatchID, feedback.TalentID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllFacultyBatchFeedback returns all the feedback for specified batch // -> FacultyTalentFeedback (admin login)
func (ser *FacultyFeedbackService) GetAllFacultyBatchFeedback(feedbacks *[]bat.FacultyTalentBatchFeedbackDTO,
	tenantID, batchID uuid.UUID, form url.Values) error {
	// get for admin login

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
	tempTalents := []bat.TalentDTO{}
	err = ser.Repository.GetAll(uow, &tempTalents,
		repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id`=talents.`id`"),
		ser.addSearchQueries(form), repository.Filter("batch_talents.`batch_id`=?", batchID),
		repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("talents.`first_name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	tempFaculty := fct.Faculty{}
	err = ser.Repository.GetRecord(uow, &tempFaculty,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempTalent := range tempTalents {
		tempDTO := bat.FacultyTalentBatchFeedbackDTO{}
		tempDTO.Talent = tempTalent
		tempDTO.Faculty = tempFaculty

		// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.BatchFeedbacks, "`created_by`",
		// 	repository.Filter("`batch_id`=? AND talent_id=?", batchID, tempTalent.ID),
		// 	repository.PreloadAssociations(ser.associations))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = ser.Repository.GetAll(uow, &tempDTO.BatchFeedbacks,
			repository.Filter("`batch_id`=? AND `talent_id`=?", batchID, tempTalent.ID),
			repository.PreloadAssociations(ser.associations),
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= faculty_talent_batch_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculty_talent_batch_feedback.`tenant_id`=?", tenantID),
			repository.OrderBy("feedback_questions.`order`"))
		// , repository.Filter("feedback_questions.`is_active` = true")
		// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = faculty_talent_batch_feedback.`option_id`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		(*feedbacks) = append((*feedbacks), tempDTO)
	}

	uow.Commit()
	return nil
}

// GetFacultyBatchFeedback returns all the feedback for specified batch // -> FacultyTalentFeedback (faculty login)
func (ser *FacultyFeedbackService) GetFacultyBatchFeedback(feedbacks *[]bat.FacultyTalentBatchFeedbackDTO,
	tenantID, batchID, facultyID uuid.UUID, requestForm url.Values) error {

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

	// check if faculty exist
	err = ser.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, true)
	tempTalents := []bat.TalentDTO{}
	err = ser.Repository.GetAll(uow, &tempTalents,
		repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id`=talents.`id`"),
		ser.addSearchQueries(requestForm), repository.Filter("batch_talents.`batch_id`=?", batchID),
		repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("talents.`first_name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	tempFaculty := fct.Faculty{}
	err = ser.Repository.GetRecord(uow, &tempFaculty,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempTalent := range tempTalents {
		tempDTO := bat.FacultyTalentBatchFeedbackDTO{}
		tempDTO.Talent = tempTalent
		tempDTO.Faculty = tempFaculty

		// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.BatchFeedbacks, "`created_by`",
		// 	repository.Filter("`batch_id`=? AND talent_id=? AND faculty_id=?", batchID, tempTalent.ID, facultyID),
		// 	repository.PreloadAssociations(ser.associations))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = ser.Repository.GetAll(uow, &tempDTO.BatchFeedbacks,
			repository.Filter("`batch_id`=? AND talent_id=? AND faculty_id=?", batchID, tempTalent.ID, facultyID),
			repository.PreloadAssociations(ser.associations),
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= faculty_talent_batch_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculty_talent_batch_feedback.`tenant_id`=?", tenantID),
			repository.OrderBy("feedback_questions.`order`"))
		// , repository.Filter("feedback_questions.`is_active` = true")
		// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = faculty_talent_batch_feedback.`option_id`"))
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
func (ser *FacultyFeedbackService) doesForeignKeyExist(feedback *bat.FacultyTalentFeedback, credentialID uuid.UUID) error {

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

	// check if question exist
	err = ser.doesFeedbackQuestionExist(feedback.TenantID, feedback.QuestionID)
	if err != nil {
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (ser *FacultyFeedbackService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (ser *FacultyFeedbackService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (ser *FacultyFeedbackService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ? AND `is_active` = true", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (ser *FacultyFeedbackService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.TalentDTO{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (ser *FacultyFeedbackService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record for the given tenant.
func (ser *FacultyFeedbackService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// addSearchQueries will add any search queries specified in query params
func (ser *FacultyFeedbackService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talents.`id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no talent_feedback record for the given tenant.
// func (ser *FacultyFeedbackService) doesTalentFeedbackExist(tenantID, feedbackID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.FacultyTalentFeedback{},
// 		repository.Filter("`id` = ?", feedbackID))
// 	if err := util.HandleError("Invalid feedback ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// // UpdateBatchFeedback will update specified feedback in table
// func (ser *FacultyFeedbackService) UpdateBatchFeedback(feedback *bat.FacultyTalentFeedback) error {

// 	// check if all foreign keys exist
// 	err := ser.doesForeignKeyExist(feedback, feedback.UpdatedBy)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent_feedback exist
// 	err = ser.doesTalentFeedbackExist(feedback.TenantID, feedback.ID)
// 	if err != nil {
// 		return err
// 	}

// 	// extract talentID and make talent nil
// 	// feedback.TalentID = feedback.Talent.ID
// 	// feedback.Talent = nil

// 	// check if talent exist
// 	err = ser.doesTalentExist(feedback.TenantID, feedback.TalentID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, false)

// 	tempFeedback := bat.FacultyTalentFeedback{}
// 	err = ser.Repository.GetRecordForTenant(uow, feedback.TenantID, &tempFeedback,
// 		repository.Filter("`id` = ?", feedback.ID), repository.Select("`created_by`"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	feedback.CreatedBy = tempFeedback.CreatedBy

// 	err = ser.Repository.Save(uow, feedback)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	uow.Commit()
// 	return nil
// }

// // DeleteBatchFeedback will delete specified feedback in table
// func (ser *FacultyFeedbackService) DeleteBatchFeedback(feedback *bat.FacultyTalentFeedback) error {

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
// 	err = ser.doesTalentFeedbackExist(feedback.TenantID, feedback.ID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, false)

// 	err = ser.Repository.UpdateWithMap(uow, bat.FacultyTalentFeedback{}, map[interface{}]interface{}{
// 		"DeletedBy": feedback.DeletedBy,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`id` = ?", feedback.ID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	uow.Commit()
// 	return nil
// }

// // GetAllFeedbackForTalent returns all the feedback for specified talent
// func (ser *FacultyFeedbackService) GetAllFeedbackForTalent(feedback *[]bat.FacultyTalentFeedback, tenantID,
// 	batchID, talentID uuid.UUID) error {

// 	// check if tenant exist
// 	err := ser.doesTenantExist(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist
// 	err = ser.doesBatchExist(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent exist
// 	err = ser.doesTalentExist(tenantID, talentID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, true)
// 	err = ser.Repository.GetAllForTenant(uow, tenantID, feedback,
// 		repository.Filter("batch_id=? AND talent_id=?", batchID, talentID),
// 		repository.PreloadAssociations(ser.associations))
// 	if err != nil {
// 		uow.RollBack()
// 		return err

// 	}
// 	uow.Commit()
// 	return nil
// }

// // GetFacultyTalentBatchFeedback will return faculty feedback for specified talent
// func (ser *FacultyFeedbackService) GetFacultyTalentBatchFeedback(feedbacks *[]bat.FacultyTalentBatchFeedbackDTO,
// 	tenantID, batchID, facultyID, talentID uuid.UUID) error {

// 	// check if tenant exist
// 	err := ser.doesTenantExist(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist
// 	err = ser.doesBatchExist(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if faculty exist
// 	err = ser.doesFacultyExist(tenantID, facultyID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent exist
// 	err = ser.doesTalentExist(tenantID, talentID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, true)
// 	tempTalent := bat.TalentDTO{}
// 	err = ser.Repository.GetRecord(uow, &tempTalent,
// 		repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id`=talents.`id`"),
// 		repository.Filter("batch_talents.`batch_id`=?", batchID), repository.OrderBy("talents.`first_name`"),
// 		repository.Filter("talents.`id`=?", talentID),
// 		repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
// 		repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	tempDTO := bat.FacultyTalentBatchFeedbackDTO{}
// 	tempDTO.Talent = tempTalent

// 	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.BatchFeedbacks, "`created_by`",
// 		repository.Filter("`batch_id`=? AND talent_id=? AND faculty_id=?", batchID, talentID, facultyID),
// 		repository.PreloadAssociations(ser.associations))
// 	if err != nil {
// 		uow.RollBack()
// 		return err

// 	}
// 	(*feedbacks) = append((*feedbacks), tempDTO)

// 	uow.Commit()
// 	return nil
// }
