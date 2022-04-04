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
	"github.com/techlabs/swabhav/tsam/web"
)

// FacultySessionFeedbackService provides methods to do different CRUD operations on faculty_talent_session_feedback table.
type FacultySessionFeedbackService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewFacultySessionFeedbackService returns a new instance Of SessionFeedbackService.
func NewFacultySessionFeedbackService(db *gorm.DB, repository repository.Repository) *FacultySessionFeedbackService {
	return &FacultySessionFeedbackService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Question", "Option",
			// "Talent",
		},
	}
}

func (service *FacultySessionFeedbackService) GetTalentFeedbackDetails(tenantID, batchID, sessionID uuid.UUID, talentSessionFeedback *[]bat.FacultyTalentBatchSessionFeedback, parser *web.Parser) error {
	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesSessionExist(tenantID, batchID, sessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAll(uow, talentSessionFeedback,
		service.addSearchQueriesForFacultyFeedbacks(parser.Form),
		repository.PreloadAssociations(service.associations),
		repository.Join("INNER Join feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" And faculty_talent_batch_session_feedback.`tenant_id` = feedback_questions.`tenant_id`"),
		repository.OrderBy("`order`"),
		repository.Filter("batch_id=?", batchID),
		repository.Filter("batch_session_id=?", sessionID),
		repository.Filter("feedback_questions.`tenant_id`=?", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// AddBatchSessionFeedback adds session feedback of talent in the table
func (service *FacultySessionFeedbackService) AddBatchSessionFeedback(sessionFeedback *bat.FacultyTalentBatchSessionFeedback, uows ...*repository.UnitOfWork) error {

	// check if foreign exist
	err := service.doesForeignKeyExist(sessionFeedback, sessionFeedback.CreatedBy)
	if err != nil {
		return err
	}

	// check if talent exist
	err = service.doesTalentExist(sessionFeedback.TenantID, sessionFeedback.TalentID)
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

// AddBatchSessionFeedbacks will add session feedback for multiple questions
func (service *FacultySessionFeedbackService) AddBatchSessionFeedbacks(feedbacks *[]bat.FacultyTalentBatchSessionFeedback, tenantID,
	batchID, batchSessionID, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)

	var talentID uuid.UUID
	var averageRating float64
	var totalRating int

	err := service.isBatchSessionCompleted(batchSessionID, tenantID)
	if err != nil {
		uow.RollBack()
		return err
	}
	for _, feedback := range *feedbacks {
		talentID = feedback.TalentID
		feedback.TenantID = tenantID
		feedback.CreatedBy = credentialID
		feedback.BatchID = batchID
		feedback.BatchSessionID = batchSessionID
		// feedback.BatchTopicID = batchTopicID

		if feedback.Option != nil {
			totalRating += feedback.Option.Key
		}

		feedback.Option = nil

		err := service.AddBatchSessionFeedback(&feedback, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	// get max score from feedback_questions table and calculate avgerage. #shailesh

	// check if talent record exist in batch_sessions_talents table.
	exist, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, bat.BatchSessionTalent{},
		repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID,
			talentID, batchSessionID))
	if err != nil {
		return err
	}
	if exist {
		averageRating = float64(totalRating / len(*feedbacks))

		err = service.Repository.UpdateWithMap(uow, bat.BatchSessionTalent{}, map[string]interface{}{
			"AverageRating": averageRating,
			"UpdatedBy":     credentialID,
		}, repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID,
			talentID, batchSessionID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// DeleteFacultyTalentBatchSessionFeedback will delete all the feedback for talent of specified session
func (service *FacultySessionFeedbackService) DeleteFacultyTalentBatchSessionFeedback(feedback *bat.FacultyTalentBatchSessionFeedback) error {

	// check if tenant exist
	err := service.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(feedback.TenantID, feedback.DeletedBy)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(feedback.TenantID, feedback.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesSessionExist(feedback.TenantID, feedback.BatchID, feedback.BatchSessionID)
	if err != nil {
		return err
	}

	// check if talent exist
	err = service.doesTalentExist(feedback.TenantID, feedback.TalentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, bat.FacultyTalentBatchSessionFeedback{}, map[string]interface{}{
		"DeletedBy": feedback.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=? AND `batch_session_id`=? AND `talent_id`=?",
		feedback.BatchID, feedback.BatchSessionID, feedback.TalentID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllFacultyBatchSessionFeedback will return all the feedback for specified batch-sesion
func (service *FacultySessionFeedbackService) GetAllFacultyBatchSessionFeedback(feedbacks *[]bat.FacultyTalentBatchSessionFeedbackDTO,
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
	err = service.doesSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
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

	tempFaculty := fct.Faculty{}
	err = service.Repository.GetRecord(uow, &tempFaculty,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempTalent := range tempTalents {
		tempDTO := bat.FacultyTalentBatchSessionFeedbackDTO{}
		tempDTO.Talent = tempTalent
		tempDTO.Faculty = tempFaculty

		// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.SessionFeedbacks, "`created_by`",
		// 	repository.Filter("`batch_id`=? AND `session_id`=? AND talent_id=?", batchID, sessionID, tempTalent.ID),
		// 	repository.PreloadAssociations(ser.associations))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = service.Repository.GetAll(uow, &tempDTO.SessionFeedbacks,
			repository.Filter("`batch_id`=? AND `batch_session_id`=? AND `talent_id`=?", batchID, batchSessionID, tempTalent.ID),
			repository.PreloadAssociations(service.associations),
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= faculty_talent_batch_session_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?", tenantID),
			// repository.Filter("feedback_questions.`is_active` = true"),
			repository.OrderBy("feedback_questions.`order`"))
		// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = faculty_talent_batch_session_feedback.`option_id`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		(*feedbacks) = append((*feedbacks), tempDTO)
	}

	uow.Commit()
	return nil
}

// GetFacultyBatchSessionFeeback will return batch-session feedback for specified faculty
func (service *FacultySessionFeedbackService) GetFacultyBatchSessionFeeback(feedbacks *[]bat.FacultyTalentBatchSessionFeedbackDTO,
	tenantID, batchID, batchSessionID, facultyID uuid.UUID, requestForm url.Values) error {

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
	err = service.doesSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
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

	tempFaculty := fct.Faculty{}
	err = service.Repository.GetRecord(uow, &tempFaculty,
		repository.Join("LEFT JOIN batches ON batches.`faculty_id`=faculties.`id`"),
		repository.Filter("batches.`id`=?", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempTalent := range tempTalents {
		tempDTO := bat.FacultyTalentBatchSessionFeedbackDTO{}
		tempDTO.Talent = tempTalent
		tempDTO.Faculty = tempFaculty

		// err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempDTO.SessionFeedbacks, "`created_by`",
		// 	repository.Filter("`batch_id`=? AND `session_id`=? AND `talent_id`=? AND `faculty_id`=?",
		// 		batchID, sessionID, tempTalent.ID, facultyID),
		// 	repository.PreloadAssociations(ser.associations))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = service.Repository.GetAll(uow, &tempDTO.SessionFeedbacks,
			repository.Filter("`batch_id`=? AND `batch_session_id`=? AND `talent_id`=? AND `faculty_id`=?",
				batchID, batchSessionID, tempTalent.ID, facultyID),
			repository.PreloadAssociations(service.associations),
			repository.Join("LEFT JOIN feedback_questions ON feedback_questions.`id`= faculty_talent_batch_session_feedback.`question_id`"),
			repository.Filter("feedback_questions.`tenant_id`=? AND feedback_questions.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?", tenantID),
			// repository.Filter("feedback_questions.`is_active` = true"),
			repository.OrderBy("feedback_questions.`order`"))
		// repository.Join("LEFT JOIN feedback_options ON feedback_options.`id` = faculty_talent_batch_session_feedback.`option_id`"))
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

func (service *FacultySessionFeedbackService) isBatchSessionCompleted(batchSessionID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Session{},
		repository.Filter("`id`=?", batchSessionID),
		// repository.Filter("`is_completed` = ?", true))
		repository.Filter("`is_session_taken` = ?", true))
	if err := util.HandleError("Feedback can only after session is completed", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesForeignKeyExist checks if all foreign keys exist
func (service *FacultySessionFeedbackService) doesForeignKeyExist(feedback *bat.FacultyTalentBatchSessionFeedback, credentialID uuid.UUID) error {

	// check if tenant exist
	err := service.doesTenantExist(feedback.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(feedback.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if feedback exit
	err = service.doesFeebackExist(feedback)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExist(feedback.TenantID, feedback.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesSessionExist(feedback.TenantID, feedback.BatchID, feedback.BatchSessionID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = service.doesFacultyExist(feedback.TenantID, feedback.FacultyID)
	if err != nil {
		return err
	}

	// check if question exist
	err = service.doesFeedbackQuestionExist(feedback.TenantID, feedback.QuestionID)
	if err != nil {
		return err
	}

	// check if option exist
	if feedback.OptionID != nil {
		err = service.doesFeedbackOptionExist(feedback.TenantID,
			feedback.QuestionID, *feedback.OptionID)
		if err != nil {
			return err
		}
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *FacultySessionFeedbackService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *FacultySessionFeedbackService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is feedback record in table for the given tenant.
func (service *FacultySessionFeedbackService) doesFeebackExist(sessionFeedback *bat.FacultyTalentBatchSessionFeedback) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, sessionFeedback.TenantID, bat.FacultyTalentBatchSessionFeedback{},
		repository.Filter("`batch_id` = ? AND `batch_session_id` = ? AND `talent_id` = ? AND `faculty_id` = ?",
			sessionFeedback.BatchID, sessionFeedback.BatchSessionID, sessionFeedback.TalentID, sessionFeedback.FacultyID))
	if err := util.HandleIfExistsError("Feedback already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no feedback question record for the given tenant.
func (service *FacultySessionFeedbackService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ? AND `is_active` = true", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no feedback option record for the given tenant.
func (service *FacultySessionFeedbackService) doesFeedbackOptionExist(tenantID, feedbackQuestionID, feedbackOptionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.FeedbackOption{},
		repository.Filter("id=? AND question_id=?", feedbackOptionID, feedbackQuestionID))
	if err := util.HandleError("Invalid feedback option ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record for the given tenant.
func (service *FacultySessionFeedbackService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.TalentDTO{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no faculty record for the given tenant.
func (service *FacultySessionFeedbackService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fct.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record for the given tenant.
func (service *FacultySessionFeedbackService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no session record for the given tenant.
func (service *FacultySessionFeedbackService) doesSessionExist(tenantID, batchID, batchSessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Session{},
		repository.Filter("`id` = ? AND `batch_id` = ?", batchSessionID, batchID))
	if err := util.HandleError("Invalid session ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent_feedback record for the given tenant.
// func (ser *FacultySessionFeedbackService) doesSessionFeedbackExist(tenantID, sessionFeedbackID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, bat.FacultyTalentBatchSessionFeedback{},
// 		repository.Filter("`id` = ?", sessionFeedbackID))
// 	if err := util.HandleError("Invalid session feedback ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// addSearchQueries will add any search queries specified in query params
func (service *FacultySessionFeedbackService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talents.`id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addSearchQueries will add any search queries specified in query params
func (service *FacultySessionFeedbackService) addSearchQueriesForFacultyFeedbacks(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("`talent_id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
