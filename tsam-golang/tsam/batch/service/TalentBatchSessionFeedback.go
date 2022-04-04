package service

import (
	"fmt"
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

// TalentSessionFeedbackService provides methods to do different CRUD operations on talent_batch_session_feedback table.
type TalentSessionFeedbackService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewTalentSessionFeedbackService returns a new instance Of SessionFeedbackService.
func NewTalentSessionFeedbackService(db *gorm.DB, repository repository.Repository) *TalentSessionFeedbackService {
	return &TalentSessionFeedbackService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Question",
			"Option",
		},
	}
}

// GetTalentSessionFeedback will get session feedback from talent to faculty for all batch-session
func (service *TalentSessionFeedbackService) GetTalentSessionFeedback(sessionFeedback *[]bat.TalentSessionFeedbackDTO,
	tenantID, batchID uuid.UUID, requestForm url.Values) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
	// repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	uow := repository.NewUnitOfWork(service.DB, true)
	tempSessions := []bat.Session{}
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &tempSessions, "date",
		repository.Filter("batch_id=?", batchID))
	if err != nil {
		uow.RollBack()
		return err
	}
	for _, tempSessions := range tempSessions {
		tempDTO := bat.TalentSessionFeedbackDTO{}
		tempDTO.BatchSession = tempSessions

		tempFeedBack := []bat.TalentBatchSessionFeedback{}
		err = service.Repository.GetAll(uow, &tempFeedBack,
			repository.Join("LEFT JOIN batch_sessions ON batch_sessions.`batch_id` = talent_batch_session_feedback.`batch_id`"),
			repository.PreloadAssociations(service.associations),
			repository.PreloadAssociations([]string{"Faculty"}),
			repository.Filter("batch_sessions.`batch_id`=?", batchID),
			repository.Filter("talent_batch_session_feedback.`batch_session_id`=?", tempSessions.ID),
			repository.Filter("batch_sessions.`tenant_id`=? AND batch_sessions.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("talent_batch_session_feedback.`batch_session_id`"))
		if err != nil {
			fmt.Println(err)
			uow.RollBack()
			return err
		}

		tempDTO.BatchSessionFeedback = &tempFeedBack
		// (*tempDTO.BatchSessionFeedback) = append((*tempDTO.BatchSessionFeedback), tempFeedBack...)

		(*sessionFeedback) = append((*sessionFeedback), tempDTO)

	}
	uow.Commit()
	return nil
}

// AddBatchSessionFeedback will add talent's feedback to specified session
func (service *TalentSessionFeedbackService) AddBatchSessionFeedback(sessionFeedback *bat.TalentBatchSessionFeedback,
	uows ...*repository.UnitOfWork) error {

	// check if foreign exist
	err := service.doesForeignKeyExist(sessionFeedback, sessionFeedback.CreatedBy)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// 	// assign answer to sessionFeedback
	// if sessionFeedback.OptionID != nil {
	// 	tempFeedbackOption := general.FeedbackOption{}
	// 	err = ser.Repository.GetRecordForTenant(uow, sessionFeedback.TenantID, &tempFeedbackOption,
	// 		repository.Filter("`id` = ?", *sessionFeedback.OptionID), repository.Select("`value`"))
	// 	if err != nil {
	// 		if length == 0 {
	// 			uow.RollBack()
	// 		}
	// 		return err
	// 	}
	// 	sessionFeedback.Answer = tempFeedbackOption.Value
	// }

	err = service.Repository.Add(uow, sessionFeedback)
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

// AddBatchSessionFeedbacks will add multiple feedbacks of talent for specified batch-session
func (service *TalentSessionFeedbackService) AddBatchSessionFeedbacks(sessionFeedbacks *[]bat.TalentBatchSessionFeedback,
	tenantID, credentialID, batchID, batchSessionID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, feedback := range *sessionFeedbacks {
		feedback.TenantID = tenantID
		feedback.CreatedBy = credentialID
		feedback.BatchID = batchID
		feedback.BatchSessionID = batchSessionID
		// feedback.BatchTopicID = batchTopicID

		err := service.AddBatchSessionFeedback(&feedback, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// DeleteTalentBatchSessionFeedback will delete all the feedback of talent of specified session
func (service *TalentSessionFeedbackService) DeleteTalentBatchSessionFeedback(sessionFeedback *bat.TalentBatchSessionFeedback) error {

	// check if tenant exist
	err := service.doesTenantExist(sessionFeedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(sessionFeedback.TenantID, sessionFeedback.DeletedBy)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(sessionFeedback.TenantID, sessionFeedback.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesBatchSessionExist(sessionFeedback.TenantID, sessionFeedback.BatchID, sessionFeedback.BatchSessionID)
	if err != nil {
		return err
	}

	// check if faculty exist
	// err = ser.doesFacultyExist(sessionFeedback.TenantID, sessionFeedback.FacultyID)
	// if err != nil {
	// 	return err
	// }

	// check if talent exist
	err = service.doesTalentExist(sessionFeedback.TenantID, sessionFeedback.TalentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, bat.TalentBatchSessionFeedback{}, map[interface{}]interface{}{
		"DeletedBy": sessionFeedback.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=? AND `batch_session_id`=? AND `talent_id`=?",
		sessionFeedback.BatchID, sessionFeedback.BatchSessionID, sessionFeedback.TalentID))
	// , sessionFeedback.FacultyID
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetTalentBatchSessionFeedback will return all the feedback for specified batch
func (service *TalentSessionFeedbackService) GetTalentBatchSessionFeedback(sessionFeedbacks *[]bat.TalentBatchSessionFeedbackDTO,
	tenantID, batchID, batchSessionID uuid.UUID, requestForm url.Values) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesBatchSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	tempFaculties := []fct.Faculty{}
	err = service.Repository.GetAll(uow, &tempFaculties,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Join("LEFT JOIN batch_sessions ON batches.`id` = batch_sessions.`batch_id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batch_sessions.`id`=?", batchSessionID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempFaculty := range tempFaculties {
		tempDTO := bat.TalentBatchSessionFeedbackDTO{}
		tempDTO.Faculty = tempFaculty

		tempTalents := []bat.TalentDTO{}
		err = service.Repository.GetAll(uow, &tempTalents,
			repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id`=talents.`id`"),
			service.addSearchQueries(requestForm), repository.Filter("batch_talents.`batch_id`=?", batchID),
			repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.OrderBy("talents.`first_name`"))
		if err != nil {
			uow.RollBack()
			return err
		}
		for _, tempTalent := range tempTalents {
			tempTalentFeedback := bat.BatchSessionFeedbackDTO{}
			tempTalentFeedback.Talent = tempTalent

			err = service.Repository.GetAll(uow, &tempTalentFeedback.SessionFeedbacks,
				repository.Filter("`batch_id`=? AND `talent_id`=? AND `batch_session_id`=?",
					batchID, tempTalent.ID, batchSessionID),
				repository.PreloadAssociations(service.associations),
				repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= talent_batch_session_feedback.`question_id`"),
				repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
				repository.Filter("talent_batch_session_feedback.`tenant_id`=?", tenantID),
				// repository.Filter("feedback_questions.`is_active` = true"),
				repository.OrderBy("feedback_questions.`order`"))
			// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = talent_batch_session_feedback.`option_id`"))
			if err != nil {
				uow.RollBack()
				return err
			}

			tempDTO.Feedbacks = append(tempDTO.Feedbacks, tempTalentFeedback)
		}
		(*sessionFeedbacks) = append((*sessionFeedbacks), tempDTO)
	}

	uow.Commit()
	return nil
}

// GetSpecifiedTalentBatchFeedback will return batch feedback for specified talent.
func (service *TalentSessionFeedbackService) GetSpecifiedTalentBatchFeedback(sessionFeedbacks *[]bat.SingleTalentBatchFeedbackDTO,
	tenantID, batchID, talentID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Create bucket for batch sessions.
	tempBatchSessions := []bat.Session{}

	// Get all batch sessions.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &tempBatchSessions, "`date`",
		repository.Filter("batch_sessions.`batch_id`=?", batchID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
		repository.Filter("batch_sessions.`is_session_taken`=?", true))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get all talent feedbacks for faculty for each batch session.
	for _, tempBatchSession := range tempBatchSessions {

		// Create bucket for SingleTalentBatchFeedbackDTO.
		tempSingleTalentBatchFeedbackDTO := bat.SingleTalentBatchFeedbackDTO{}

		// Give batch session id to tempSingleTalentBatchFeedbackDTO.
		tempSingleTalentBatchFeedbackDTO.BatchSessionID = tempBatchSession.ID
		tempSingleTalentBatchFeedbackDTO.Date = tempBatchSession.Date

		err = service.Repository.GetAll(uow, &tempSingleTalentBatchFeedbackDTO.SessionFeedbacks,
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= talent_batch_session_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("`batch_session_id`=?", tempBatchSession.ID),
			repository.Filter("talent_batch_session_feedback.`tenant_id`=? AND talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
			repository.PreloadAssociations(service.associations),
			repository.OrderBy("feedback_questions.`order`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		// Push into sessionFeedbacks.
		(*sessionFeedbacks) = append((*sessionFeedbacks), tempSingleTalentBatchFeedbackDTO)
	}

	uow.Commit()
	return nil
}

// GetSpecifiedTalentBatchSessionFeedback will return batch feedback for specified talent got one
// batch session.
func (service *TalentSessionFeedbackService) GetSpecifiedTalentBatchSessionFeedback(sessionFeedbacks *[]bat.TalentBatchSessionFeedback,
	tenantID, batchID, talentID, batchSessionID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesBatchSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all talent batch session feedbacks.
	err = service.Repository.GetAll(uow, &sessionFeedbacks,
		repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= talent_batch_session_feedback.`question_id`"),
		repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
		repository.Filter("talent_batch_session_feedback.`batch_id`=?", batchID),
		repository.Filter("talent_batch_session_feedback.`batch_session_id`=?", batchSessionID),
		repository.Filter("talent_batch_session_feedback.`talent_id`=?", talentID),
		repository.Filter("talent_batch_session_feedback.`tenant_id`=?", tenantID),
		repository.Filter("talent_batch_session_feedback.`deleted_at` IS NULL"),
		repository.OrderBy("feedback_questions.`order`"))

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

// doesForeignKeyExist checks if all foreign keys exist
func (service *TalentSessionFeedbackService) doesForeignKeyExist(sessionFeedback *bat.TalentBatchSessionFeedback, credentialID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(sessionFeedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(sessionFeedback.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(sessionFeedback.TenantID, sessionFeedback.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesBatchSessionExist(sessionFeedback.TenantID, sessionFeedback.BatchID, sessionFeedback.BatchSessionID)
	if err != nil {
		return err
	}

	// check if talent exist
	err = service.doesTalentExist(sessionFeedback.TenantID, sessionFeedback.TalentID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = service.doesFacultyExist(sessionFeedback.TenantID, sessionFeedback.FacultyID)
	if err != nil {
		return err
	}

	// check if question exist
	err = service.doesFeedbackQuestionExist(sessionFeedback.TenantID, sessionFeedback.QuestionID)
	if err != nil {
		return err
	}

	// check if option exist
	if sessionFeedback.OptionID != nil {
		err = service.doesFeedbackOptionExist(sessionFeedback.TenantID,
			sessionFeedback.QuestionID, *sessionFeedback.OptionID)
		if err != nil {
			return err
		}
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *TalentSessionFeedbackService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *TalentSessionFeedbackService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record for the given tenant.
func (service *TalentSessionFeedbackService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no session record for the given tenant.
func (service *TalentSessionFeedbackService) doesBatchSessionExist(tenantID, batchID, batchSessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Session{},
		repository.Filter("`id` = ? AND `batch_id` = ?", batchSessionID, batchID))
	if err := util.HandleError("Invalid batch session ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no faculty record for the given tenant.
func (service *TalentSessionFeedbackService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (service *TalentSessionFeedbackService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ? AND `is_active` = true", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback option record for the given tenant.
func (service *TalentSessionFeedbackService) doesFeedbackOptionExist(tenantID, feedbackQuestionID, feedbackOptionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackOption{},
		repository.Filter("id=? AND question_id=?", feedbackOptionID, feedbackQuestionID))
	if err := util.HandleError("Invalid feedback option ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (service *TalentSessionFeedbackService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.TalentDTO{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// addSearchQueries will add any search queries specified in query params
func (service *TalentSessionFeedbackService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talents.`id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}
	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("faculty_id", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// // DeleteFacultyBatchSessionFeedback will delete all the feedback of faculty of specified session
// func (ser *TalentSessionFeedbackService) DeleteFacultyBatchSessionFeedback(sessionFeedback *bat.TalentBatchSessionFeedback) error {

// 	// check if tenant exist
// 	err := ser.doesTenantExist(sessionFeedback.TenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if credential exist
// 	err = ser.doesCredentialExist(sessionFeedback.TenantID, sessionFeedback.DeletedBy)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist
// 	err = ser.doesBatchExist(sessionFeedback.TenantID, sessionFeedback.BatchID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if session exist
// 	err = ser.doesSessionExist(sessionFeedback.TenantID, sessionFeedback.SessionID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent exist
// 	err = ser.doesFacultyExist(sessionFeedback.TenantID, sessionFeedback.FacultyID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, false)

// 	err = ser.Repository.UpdateWithMap(uow, bat.TalentBatchSessionFeedback{}, map[interface{}]interface{}{
// 		"DeletedBy": sessionFeedback.DeletedBy,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`batch_id`=? AND `session_id`=? AND `faculty_id`=?",
// 		sessionFeedback.BatchID, sessionFeedback.SessionID, sessionFeedback.FacultyID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	uow.Commit()
// 	return nil
// }
