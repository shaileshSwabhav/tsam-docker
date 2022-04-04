package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchSessionTalentService provides method to Update, Delete, Add, Get Method For batch_session_talents.
type BatchSessionTalentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewSessionBatchTalentService returns a new instance of BatchSessionTalent.
func NewSessionBatchTalentService(db *gorm.DB, repository repository.Repository) *BatchSessionTalentService {
	return &BatchSessionTalentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Batch", "BatchSession", "Talent",
		},
	}
}

func (service *BatchSessionTalentService) GetTalentTopicDetails(tenantID, batchID, talentID uuid.UUID, talentSessionDetails *[]batch.TalentDetailsDTO, parser *web.Parser) error {
	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, talentSessionDetails,
		repository.Filter("`batch_id` = ?", batchID),
		repository.OrderBy("`date`"),
		repository.GroupBy("`id`"))
	// service.addSearchQueries(parser.Form),
	// repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *talentSessionDetails {
		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.FacultyTalentBatchSessionFeedback{},
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID,
				(*talentSessionDetails)[index].ID))
		if err != nil {
			return err
		}

		(*talentSessionDetails)[index].IsFeedbackGiven = exist

		exist, err = repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchSessionTalent{},
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID,
				(*talentSessionDetails)[index].ID))
		if err != nil {
			return err
		}
		if !exist {
			(*talentSessionDetails)[index].IsPresent = false
			continue
		}

		err = service.Repository.Scan(uow, &(*talentSessionDetails)[index], repository.Table("batch_session_talents"),
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID,
				(*talentSessionDetails)[index].ID), repository.Select("batch_session_talents.`is_present`"))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// AddTalentAttendance will add talent's attendance details to the table.
func (service *BatchSessionTalentService) AddTalentAttendance(talentAttendance *batch.BatchSessionTalent,
	uows ...*repository.UnitOfWork) error {

	// check if foreign keys exist.
	err := service.doesForeignKeyExist(talentAttendance, talentAttendance.CreatedBy)
	if err != nil {
		return err
	}

	// check if talent attendance exist.
	err = service.doesTalentAttendanceExist(*talentAttendance)
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

	err = service.Repository.Add(uow, talentAttendance)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		// uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddTalentAttendance will add talent's attendance details to the table.
func (service *BatchSessionTalentService) AddTalentsAttendance(tenantID, credentialID, batchSessionID, batchID uuid.UUID,
	talentAttendance *[]batch.BatchSessionTalent) error {

	uow := repository.NewUnitOfWork(service.DB, false)
	for _, talent := range *talentAttendance {
		talent.TenantID = tenantID
		talent.BatchID = batchID
		talent.BatchSessionID = batchSessionID
		talent.CreatedBy = credentialID
		err := service.AddTalentAttendance(&talent, uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// UpdateTalentAttendance will update attendance for specified talent.
func (service *BatchSessionTalentService) UpdateTalentAttendance(talentAttendance *batch.BatchSessionTalent) error {

	// check if foreign keys exist.
	err := service.doesForeignKeyExist(talentAttendance, talentAttendance.CreatedBy)
	if err != nil {
		return err
	}

	// check if talent attendance exist.
	err = service.doesBatchTopicTalentExist(talentAttendance.TenantID, talentAttendance.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, talentAttendance)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchSessionTalents will get all the batch session talents.
func (service *BatchSessionTalentService) GetBatchSessionTalents(tenantID, batchID, batchSessionID uuid.UUID,
	sessionTalents *[]batch.BatchSessionTalentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// check if batch session exist.
	err = service.doesBatchSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, sessionTalents,
		repository.Filter("`batch_id` = ? AND "+"`batch_session_id` = ?", batchID, batchSessionID),
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}
	for index := range *sessionTalents {
		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.FacultyTalentBatchSessionFeedback{},
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, ((*sessionTalents)[index].TalentID),
				batchSessionID))
		if err != nil {
			return err
		}

		(*sessionTalents)[index].IsFeedbackGiven = exist

		exist, err = repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchSessionTalent{},
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, ((*sessionTalents)[index].TalentID),
				batchSessionID))
		if err != nil {
			return err
		}
		if !exist {
			*(*sessionTalents)[index].IsPresent = false
			continue
		}

		err = service.Repository.Scan(uow, &(*sessionTalents)[index], repository.Table("batch_session_talents"),
			repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, ((*sessionTalents)[index].TalentID),
				batchSessionID), repository.Select("batch_session_talents.`is_present`"))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetAllBatchSessionForTalent will get all the batch sessions talents for specific talent.
func (service *BatchSessionTalentService) GetAllBatchSessionForTalent(tenantID, batchID, talentID uuid.UUID,
	sessionTalents *[]batch.BatchSessionTalent, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, sessionTalents,
		repository.Filter("`batch_id` = ? AND `talent_id` = ?", batchID, talentID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAverageRatingForTalentByWeek will get average rating for batch for current week.
func (service *BatchSessionTalentService) GetAverageRatingForTalentByWeek(tenantID, batchID, talentID uuid.UUID,
	averageRatingBatch *batch.AverageRatingBatch, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecord(uow, averageRatingBatch,
		repository.Select("SUM(answer) as average_rating, SUM(max_score) as max_score"),
		repository.Join("INNER JOIN feedback_questions on feedback_questions.`id` = faculty_talent_batch_session_feedback.`question_id`"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_questions.`tenant_id` = ?", tenantID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_sessions.`tenant_id` = ?", tenantID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL AND faculty_talent_batch_session_feedback.`tenant_id` = ?", tenantID),
		repository.Filter("batch_sessions.`batch_id` = ? AND faculty_talent_batch_session_feedback.`talent_id` = ?", batchID, talentID),
		repository.Filter("YEARWEEK(batch_sessions.`date`, 1) = YEARWEEK(CURDATE(), 1)"),
		repository.GroupBy("faculty_talent_batch_session_feedback.`talent_id`"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uow.Commit()
			return nil
		} 
		return errors.NewValidationError("Internal server error")
	}

	uow.Commit()
	return nil
}

// GetOneBatchSessionTalentsForTalent will get one the batch sessions talents for specific talent.
func (service *BatchSessionTalentService) GetOneBatchSessionTalentsForTalent(tenantID, batchID, talentID,
	batchSessionID uuid.UUID, sessionTalent *batch.BatchSessionTalent, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// check if batch session exist.
	err = service.doesBatchSessionExist(tenantID, batchID, batchSessionID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, sessionTalent,
		repository.Filter("`batch_id` = ? AND `talent_id` = ? AND batch_session_id=?", batchID, talentID, batchSessionID))
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

// func (service *BatchSessionTalentService) IsTalentandFeedbackPresent(tenantID, batchID, talentID, sessionID uuid.UUID, isFeedbackGiven, isPresent *bool,
// 	out interface{}, uow *repository.UnitOfWork) error {
// 	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.FacultyTalentBatchSessionFeedback{},
// 		repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID, sessionID))
// 	if err != nil {
// 		return err
// 	}

// 	// (*sessionTalents)[index].IsFeedbackGiven = exist
// 	isFeedbackGiven = &exist

// 	exist, err = repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchSessionTalent{},
// 		repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID,
// 			sessionID))
// 	if err != nil {
// 		return err
// 	}
// 	if !exist {
// 		// *(*sessionTalents)[index].IsPresent = false
// 		*isPresent = false
// 		// continue
// 	}

// 	err = service.Repository.Scan(uow, &out, repository.Table("batch_session_talents"),
// 		repository.Filter("`batch_id` = ? AND `talent_id` = ? AND `batch_session_id` = ?", batchID, talentID,
// 			sessionID), repository.Select("batch_session_talents.`is_present`"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	return nil
// }

// doesForeignKeyExist will check if all foreign keys are valid.
func (service *BatchSessionTalentService) doesForeignKeyExist(talentAttendance *batch.BatchSessionTalent, credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(talentAttendance.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(talentAttendance.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExist(talentAttendance.TenantID, talentAttendance.BatchID)
	if err != nil {
		return err
	}

	// check if batch session exist.
	err = service.doesBatchSessionExist(talentAttendance.TenantID, talentAttendance.BatchID, talentAttendance.BatchSessionID)
	if err != nil {
		return err
	}

	// check if talent exist.
	err = service.doesTalentExist(talentAttendance.TenantID, talentAttendance.TalentID)
	if err != nil {
		return err
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *BatchSessionTalentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *BatchSessionTalentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *BatchSessionTalentService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no batch sessison record in table for the given tenant.
func (service *BatchSessionTalentService) doesBatchSessionExist(tenantID, batchID, batchSessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Session{},
		repository.Filter("`id` = ? AND `batch_id` = ?", batchSessionID, batchID))
	if err := util.HandleError("Invalid batch session ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no talent record in table for the given tenant.
func (service *BatchSessionTalentService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no batch topic talent record in table for the given tenant.
func (service *BatchSessionTalentService) doesBatchTopicTalentExist(tenantID, sesionTalentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", sesionTalentID))
	if err := util.HandleError("Invalid talent attendance ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no batch session talent record in table for the given tenant.
func (service *BatchSessionTalentService) doesTalentAttendanceExist(talentAttendance batch.BatchSessionTalent) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, talentAttendance.TenantID,
		batch.BatchSessionTalent{}, repository.Filter("`batch_id` = ? AND `batch_session_id` = ? AND `talent_id` = ?",
			talentAttendance.BatchID, talentAttendance.BatchSessionID, talentAttendance.TalentID))
	if err := util.HandleIfExistsError("Talent attendance already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueries will add query processors based on queryparams.
func (service *BatchSessionTalentService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if isPresent, ok := requestForm["isPresent"]; ok {
		util.AddToSlice("batch_session_talents.`is_present`", "= ?", "AND", isPresent, &columnNames, &conditions, &operators, &values)
	}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("batch_session_talents.`talent_id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
