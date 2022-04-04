package service

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchTalentService provides method to add, update and get method for batch talent.
type BatchTalentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewBatchTalentService returns a new instance of BatchTalentService.
func NewBatchTalentService(db *gorm.DB, repository repository.Repository) *BatchTalentService {
	return &BatchTalentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Talent",
		},
	}
}

// AddTalentsToBatch adds student to a particular batch
func (service *BatchTalentService) AddTalentsToBatch(talents *[]bat.MappedTalent, tenantID, batchID, credentialID uuid.UUID) error {

	// Validates all the foreign keys for batch_talents table which are talentID, batchID, tenantID and credentialID
	err := service.validateMappedForeignKeys(talents, batchID, tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	batch := bat.Batch{}
	batch.ID = batchID

	// Start Transcation.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Gets the record of the batch
	err = service.Repository.GetRecordForTenant(uow, tenantID, &batch)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Checks if batch can be allocated to more students
	if batch.TotalStudents != nil && batch.TotalIntake != nil {
		if *batch.TotalIntake == *batch.TotalStudents {
			log.NewLogger().Error("Current batch is full")
			return errors.NewValidationError("Current batch is full")
		}
		if *batch.TotalIntake < *batch.TotalStudents+uint8(len(*talents)) {
			return errors.NewValidationError("Only " + strconv.Itoa(int(*batch.TotalIntake-*batch.TotalStudents)) + " slot left")
		}
	}

	duplicateTalent := &talent.Talent{}

	// Check if JSON has duplicate entries
	talentID, found := service.checkFieldUniquenessInJSON(talents)
	if found {
		err := service.Repository.GetRecordForTenant(uow, tenantID, duplicateTalent, repository.Filter("`id` = ?", talentID),
			repository.Select([]string{"`first_name`", "`last_name`"}))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		log.NewLogger().Error(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " selected multiple times")
		return errors.NewValidationError(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " selected multiple times")
	}

	for _, talent := range *talents {

		// Check if talent is already assigned with same batch
		exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &bat.MappedTalent{},
			repository.Filter("batch_id=? AND talent_id=?", batchID, talent.TalentID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		if exists {
			err := service.Repository.GetRecordForTenant(uow, tenantID, duplicateTalent, repository.Filter("`id` = ?", talent.TalentID),
				repository.Select([]string{"`first_name`", "`last_name`"}))
			if err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
			log.NewLogger().Error(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " is already added to this batch")
			return errors.NewValidationError(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " is already added to this batch")
		}

		// Assign talent to the specified batch
		talent.BatchID = batchID
		talent.TenantID = tenantID
		talent.CreatedBy = credentialID

		err = service.Repository.Add(uow, &talent)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}

	}

	// Update totalStudents in batch
	err = service.Repository.UpdateWithMap(uow, batch, map[interface{}]interface{}{
		"TotalStudents": *batch.TotalStudents + uint8(len(*talents)),
	})
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// UpdateBatchTalent updates batch talent.
func (service BatchTalentService) UpdateBatchTalent(batchTalent *bat.MappedTalent) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(batchTalent, batchTalent.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if batch talent exist.
	err = service.doesBatchTalentExist(batchTalent.TenantID, batchTalent.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting batch talent already present in database.
	tempBatchTalent := bat.MappedTalent{}

	// Get batch talent for getting created_by field of batch talent from database.
	if err := service.Repository.GetForTenant(uow, batchTalent.TenantID, batchTalent.ID, &tempBatchTalent); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from batch talent blog to batch talent to be updated.
	batchTalent.CreatedBy = tempBatchTalent.CreatedBy

	err = service.Repository.Save(uow, batchTalent)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateSuspensionDateBatchTalent updates suspension date of batch talent.
func (service BatchTalentService) UpdateSuspensionDateBatchTalent(batchTalent *bat.UpdateBatchTalentSuspension,
	tenantID, credentialID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if batch talent exists.
	err = service.doesBatchTalentExist(tenantID, batchTalent.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Get today's date.
	currentDate := time.Now().Format("2006-01-02")

	// If is suspend is true then set suspension date.
	if batchTalent.IsSuspend {
		err = service.Repository.UpdateWithMap(uow, &bat.MappedTalent{}, map[interface{}]interface{}{
			"SuspensionDate": currentDate,
			"UpdatedBy":      credentialID,
		}, repository.Filter("`id`=?", batchTalent.ID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	// If is suspend is false then make suspension date as null.
	if !batchTalent.IsSuspend {
		err = service.Repository.UpdateWithMap(uow, &bat.MappedTalent{}, map[interface{}]interface{}{
			"SuspensionDate": nil,
			"UpdatedBy":      credentialID,
		}, repository.Filter("`id`=?", batchTalent.ID))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateIsActiveBatchTalent updates is active of batch talent.
func (service BatchTalentService) UpdateIsActiveBatchTalent(batchTalent *bat.UpdateBatchTalentIsActive,
	tenantID, credentialID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if batch talent exists.
	err = service.doesBatchTalentExist(tenantID, batchTalent.ID)
	if err != nil {
		return err
	}

	// Validate batch ID.
	err = service.doesBatchExist(tenantID, batchTalent.BatchID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for batch.
	batch := bat.Batch{}
	batch.ID = batchTalent.BatchID

	// Gets the record of the batch.
	err = service.Repository.GetRecordForTenant(uow, tenantID, &batch)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// If is active is true then increment total talent count of batch by one.
	if batchTalent.IsActive {

		// Checks if batch can be allocated to more students.
		if batch.TotalStudents != nil && *batch.TotalIntake != 0 {
			if *batch.TotalIntake == *batch.TotalStudents {
				log.NewLogger().Error("Current batch is full")
				return errors.NewValidationError("Current batch is full")
			}
		}

		// Increment totalStudents in batch.
		err = service.Repository.UpdateWithMap(uow, batch, map[interface{}]interface{}{
			"TotalStudents": *batch.TotalStudents + uint8(1),
		})
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
	}

	// If is active is false then decrement total talent count of batch by one.
	if !batchTalent.IsActive {

		// Decrement totalStudents in batch.
		err = service.Repository.UpdateWithMap(uow, batch, map[interface{}]interface{}{
			"TotalStudents": *batch.TotalStudents - uint8(1),
		})
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
	}

	// Update batch talent.
	err = service.Repository.UpdateWithMap(uow, &bat.MappedTalent{}, map[interface{}]interface{}{
		"IsActive":  batchTalent.IsActive,
		"UpdatedBy": credentialID,
	}, repository.Filter("`id`=?", batchTalent.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchMultipleTalentDetails will return details of all talents for specified batch.
func (service *BatchTalentService) GetBatchMultipleTalentDetails(tenantID, batchID uuid.UUID,
	sessionTalents *[]bat.BatchTalentDTO, parser *web.Parser) error {

	now := time.Now()

	defer func() {
		fmt.Println("=================== duration ->", time.Since(now))
	}()

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

	err = service.Repository.GetAll(uow, sessionTalents,
		repository.Join("INNER JOIN talents ON batch_talents.`talent_id` = talents.`id` AND "+
			"batch_talents.`tenant_id` = talents.`tenant_id`"), repository.Filter("batch_talents.`tenant_id` = ?", tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL"),
		repository.Filter("`batch_id` = ?", batchID), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Batch", "Talent"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	var totalSessionsCount uint
	var sessionTotalHours bat.SessionTotalHours

	err = service.getBatchTotalSessions(uow, &totalSessionsCount, tenantID, batchID, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getBatchSessionTotalHours(uow, tenantID, batchID, &sessionTotalHours, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	talentAttendanceFieldChannel := make(chan error, 1)
	talentFeedbackFieldChannel := make(chan error, 1)

	for index := range *sessionTalents {
		(*sessionTalents)[index].TotalSessionsCount = totalSessionsCount
		(*sessionTalents)[index].TotalHours = sessionTotalHours.TotalHours

		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.BatchSessionTalent{},
			repository.Filter("`batch_id` = ? AND `talent_id` = ?", batchID, (*sessionTalents)[index].TalentID))
		if err != nil {
			return err
		}
		if exist {
			// get all talent fields.
			// can be done using channels.
			go service.getTalentAttendanceFields(uow, tenantID, batchID, &(*sessionTalents)[index], talentAttendanceFieldChannel, parser)
			go service.getTalentFeedbackFields(uow, tenantID, batchID, &(*sessionTalents)[index], talentFeedbackFieldChannel)
			// if err != nil {
			// 	uow.RollBack()
			// 	return err
			// }

			// can be done using channels.
			err = <-talentAttendanceFieldChannel
			if err != nil {
				uow.RollBack()
				return err
			}

			err = <-talentFeedbackFieldChannel
			if err != nil {
				uow.RollBack()
				return err
			}
		}
	}

	uow.Commit()
	return nil
}

// GetBatchTalentDetails will return details of one talent for specified batch.
func (service *BatchTalentService) GetBatchTalentDetails(tenantID uuid.UUID, batchTalent *bat.BatchTalentDTO,
	parser *web.Parser) error {

	now := time.Now()

	defer func() {
		fmt.Println("=================== duration ->", time.Since(now))
	}()

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExist(tenantID, batchTalent.BatchID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(tenantID, batchTalent.TalentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get batch talent.
	err = service.Repository.GetRecord(uow, batchTalent,
		repository.Select("batches.`meet_link` AS batch_meet_link, batches.`telegram_link` AS batch_telegram_link,"+
			"batches.`start_date` AS start_date, batches.`estimated_end_date` AS estimated_end_date, `batch_talents`.*"),
		repository.Join("INNER JOIN talents ON batch_talents.`talent_id` = talents.`id` AND "+
			"batch_talents.`tenant_id` = talents.`tenant_id`"),
		// repository.Join("JOIN batch_timing on batch_timing.`batch_id` = batch_talents.`batch_id`"),
		repository.Join("JOIN batches on batches.`id` = batch_talents.`batch_id`"),
		// repository.Join("JOIN faculties on batches.`faculty_id` = faculties.`id`"),
		repository.Filter("batch_talents.`tenant_id` = ?", tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL AND "+
			"batches.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`batch_id` = ?", batchTalent.BatchID),
		repository.Filter("batch_talents.`talent_id` = ?", batchTalent.TalentID),
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Batch", "Talent"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	var totalSessionsCount uint
	var totalSessionsCompletedCount uint
	var sessionTotalHours bat.SessionTotalHours
	var sessionTotalCompletedHours bat.SessionTotalHours

	// Get total sessions of batch.
	err = service.getBatchTotalSessions(uow, &totalSessionsCount, tenantID, batchTalent.BatchID, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get total completed sessions of batch.
	err = service.getBatchTotalCompletedSessions(uow, &totalSessionsCompletedCount, tenantID, batchTalent.BatchID, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get total hours of batch.
	err = service.getBatchSessionTotalHours(uow, tenantID, batchTalent.BatchID, &sessionTotalHours, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get total completed hours of batch.
	err = service.getBatchSessionTotalCompletedHours(uow, tenantID, batchTalent.BatchID, &sessionTotalCompletedHours, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get batch timings.
	tempBatchTimings := []bat.Timing{}
	err = service.Repository.GetAll(uow, &tempBatchTimings,
		repository.Join("JOIN days on batch_timing.`day_id` = days.`id`"),
		repository.Filter("batch_timing.`batch_id`=?", batchTalent.BatchID),
		repository.Filter("batch_timing.`tenant_id`=? AND days.`tenant_id`=?", tenantID, tenantID),
		repository.Filter("batch_timing.`deleted_at` IS NULL AND days.`deleted_at` IS NULL"),
		repository.OrderBy("days.`order`"),
		repository.PreloadAssociations([]string{"Day"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Give batch timings to batch talent.
	batchTalent.BatchTimings = tempBatchTimings

	talentAttendanceFieldChannel := make(chan error, 1)
	talentFeedbackFieldChannel := make(chan error, 1)

	batchTalent.TotalSessionsCount = totalSessionsCount
	batchTalent.TotalSessionsCompleted = totalSessionsCompletedCount
	batchTalent.TotalHours = sessionTotalHours.TotalHours
	batchTalent.TotalCompletedHours = sessionTotalCompletedHours.TotalHours

	// Check if batch session talent exists or not.
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.BatchSessionTalent{},
		repository.Filter("`batch_id` = ? AND `talent_id` = ?", batchTalent.BatchID, batchTalent.TalentID))
	if err != nil {
		return err
	}

	if exist {
		// get all talent fields.
		// can be done using channels.
		go service.getTalentAttendanceFields(uow, tenantID, batchTalent.BatchID, batchTalent, talentAttendanceFieldChannel, parser)
		go service.getTalentFeedbackFields(uow, tenantID, batchTalent.BatchID, batchTalent, talentFeedbackFieldChannel)
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		// can be done using channels.
		err = <-talentAttendanceFieldChannel
		if err != nil {
			uow.RollBack()
			return err
		}

		err = <-talentFeedbackFieldChannel
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetBatchesOfTalent will return batches of one talent.
func (service *BatchTalentService) GetBatchesOfTalent(tenantID, talentID uuid.UUID,
	batchTalent *[]bat.MinimumBatchTalentForTalent, parser *web.Parser) error {

	now := time.Now()

	defer func() {
		fmt.Println("=================== duration ->", time.Since(now))
	}()

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(tenantID, talentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get batch talent.
	err = service.Repository.GetAll(uow, batchTalent,
		repository.Select("courses.`name` AS course_name, batches.`status` AS batch_status, batch_talents.*"),
		repository.Join("JOIN batches ON batch_talents.`batch_id` = batches.`id`"),
		repository.Join("JOIN courses ON batches.`course_id` = courses.`id`"),
		repository.Filter("batch_talents.`tenant_id` = ? AND batches.`tenant_id`=?", tenantID, tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id` = ?", talentID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTalentListOfBatch will return batches of one talent.
func (service *BatchTalentService) GetTalentListOfBatch(tenantID, batchID uuid.UUID,
	batchTalents *[]list.Talent, parser *web.Parser) error {

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
	defer uow.Commit()

	// Get batch talent.
	err = service.Repository.GetAll(uow, batchTalents,
		repository.Select("`talents`.`id`,`talents`.`email`,`talents`.`first_name`,`talents`.`last_name`"),
		repository.Join("INNER JOIN `batch_talents` ON `talents`.`id` = `batch_talents`.`talent_id`"),
		repository.Filter("`batch_talents`.`tenant_id` = ? AND `talents`.`tenant_id`=?", tenantID, tenantID),
		repository.Filter("`batch_talents`.`deleted_at` IS NULL AND `talents`.`deleted_at` IS NULL"),
		repository.Filter("`batch_talents`.`batch_id` = ? AND `batch_talents`.`is_active` = ?", batchID, true),
		repository.OrderBy("`talents`.`first_name`"))
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

// doesForeignKeyExist will check if all foreign keys are valid.
func (service BatchTalentService) doesForeignKeyExist(batchTalent *bat.MappedTalent, credentialID uuid.UUID) error {

	// Check if tenant exists.
	err := service.doesTenantExist(batchTalent.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exists.
	err = service.doesCredentialExist(batchTalent.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if batch exists.
	err = service.doesBatchExist(batchTalent.TenantID, batchTalent.BatchID)
	if err != nil {
		return err
	}

	// Check if talent exists.
	err = service.doesTalentExist(batchTalent.TenantID, batchTalent.TalentID)
	if err != nil {
		return err
	}

	return nil
}

func (service *BatchTalentService) getBatchTotalSessions(uow *repository.UnitOfWork, totalCount *uint,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.Repository.GetCount(uow, bat.Session{}, totalCount,
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.batch_session_id = batch_sessions.id"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_sessions.`batch_id` = ?", batchID),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.GroupBy("batch_sessions.`id`"))
	if err != nil {
		return err
	}

	return nil
}

func (service *BatchTalentService) getBatchTotalCompletedSessions(uow *repository.UnitOfWork, totalCount *uint,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	// err := service.Repository.GetCount(uow, bat.MappedSession{}, totalCount,
	// 	repository.Join("INNER JOIN course_sessions ON course_sessions.`id` = batch_sessions.`course_session_id` AND "+
	// 		"course_sessions.`tenant_id` = batch_sessions.`tenant_id`"),
	// 	repository.Filter("batch_sessions.`tenant_id` = ?", tenantID),
	// 	repository.Filter("batch_sessions.`deleted_at` IS NULL AND course_sessions.`deleted_at` IS NULL"),
	// 	repository.Filter("course_sessions.`session_id` IS NULL AND batch_sessions.`batch_id` = ?", batchID),
	// 	repository.Filter("batch_sessions.`is_completed` = 1"))
	err := service.Repository.GetCountForTenant(uow, tenantID, bat.Session{}, totalCount,
		repository.Filter("batch_sessions.`batch_id`=?", batchID),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_sessions.`tenant_id`=?", tenantID),
		repository.Filter("batch_sessions.`is_session_taken`=?", true))
	if err != nil {
		return err
	}

	return nil
}

func (service *BatchTalentService) getBatchSessionTotalHours(uow *repository.UnitOfWork,
	tenantID, batchID uuid.UUID, sessionTotalHours *bat.SessionTotalHours, parser *web.Parser) error {

	err := service.Repository.Scan(uow, sessionTotalHours, repository.Table("batch_session_topics"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_hours`"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ?", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`deleted_at` IS NULL", batchID))
	if err != nil {
		return err
	}
	return nil
}

func (service *BatchTalentService) getBatchSessionTotalCompletedHours(uow *repository.UnitOfWork,
	tenantID, batchID uuid.UUID, sessionTotalHours *bat.SessionTotalHours, parser *web.Parser) error {

	err := service.Repository.Scan(uow, sessionTotalHours, repository.Table("batch_session_topics"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_hours`"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ?", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_topics.`is_completed` = ?", 1),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`deleted_at` IS NULL", batchID))
	if err != nil {
		return err
	}
	return nil
}

func (service *BatchTalentService) getTalentAttendanceFields(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	sessionTalent *bat.BatchTalentDTO, channel chan error, parser *web.Parser) {
	//, channel chan error

	// get count of sessions attended by talent.
	err := service.getTotalSessionsAttendedByTalent(uow, &sessionTalent.SessionsAttended, tenantID, batchID, sessionTalent.TalentID, parser)
	if err != nil {
		channel <- err
		return
		// return err
	}

	// get total session hours attended by talent.
	err = service.getTalentTotalAttendedHours(uow, tenantID, batchID, sessionTalent, parser)
	if err != nil {
		channel <- err
		return
		// return err
	}

	channel <- nil
	// return nil
}

func (service *BatchTalentService) getTalentFeedbackFields(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	sessionTalent *bat.BatchTalentDTO, channel chan error) {
	//, channel chan error

	// get total session feedbacks givent to talent.
	err := service.getFeedbacksGivenToTalent(uow, tenantID, batchID, sessionTalent.TalentID, &sessionTalent.TotalFeedbacksGiven)
	if err != nil {
		channel <- err
		return
		// return err
	}

	// get talent average rating.
	err = service.getTalentAverageRating(uow, tenantID, batchID, sessionTalent)
	if err != nil {
		channel <- err
		return
		// return err
	}

	channel <- nil
	// return nil
}

func (service *BatchTalentService) getTotalSessionsAttendedByTalent(uow *repository.UnitOfWork, totalCount *uint,
	tenantID, batchID, talentID uuid.UUID, parser *web.Parser) error {

	err := service.Repository.GetCount(uow, bat.BatchSessionTalent{}, totalCount,
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id`=batch_session_talents.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_session_talents.`batch_id` = ? AND `talent_id` = ? AND `is_present` = ?", batchID, talentID, true),
		repository.Filter("batch_session_talents.`tenant_id` = ? AND batch_session_talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID))
	if err != nil {
		return err
	}

	return nil
}

func (service *BatchTalentService) getTalentTotalAttendedHours(uow *repository.UnitOfWork,
	tenantID, batchID uuid.UUID, sessionTalent *bat.BatchTalentDTO, parser *web.Parser) error {

	// err := service.Repository.Scan(uow, sessionTotalHours, repository.Table("batch_sessions"),
	// repository.Select("SUM(course_sessions.`hours`) AS `total_hours`"),
	// repository.Join("INNER JOIN course_sessions ON course_sessions.`id` = batch_sessions.`course_session_id` AND "+
	// 	"course_sessions.`tenant_id` = batch_sessions.`tenant_id`"), repository.Filter("batch_sessions.`tenant_id` = ?", tenantID),
	// repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`deleted_at` IS NULL", batchID),
	// repository.Filter("course_sessions.`deleted_at` IS NULL AND course_sessions.`session_id` IS NULL"))

	err := service.Repository.Scan(uow, sessionTalent,
		repository.Table("batch_session_talents"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `attended_hours`"),
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_talents.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`batch_session_id` = batch_sessions.`id`"),
		repository.Filter("batch_session_talents.`batch_id` = ? AND batch_session_talents.`talent_id` = ?",
			batchID, sessionTalent.TalentID),
		repository.Filter("batch_session_topics.`is_completed`=?", true),
		repository.Filter("batch_session_talents.`is_present` = ?", 1),
		repository.Filter("batch_session_talents.`tenant_id` = ? AND batch_session_talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_sessions.`deleted_at` IS NULL"))
	if err != nil {
		return err
	}
	return nil
}

func (service *BatchTalentService) getTalentAverageRating(uow *repository.UnitOfWork, tenantID,
	batchID uuid.UUID, sessionTalent *bat.BatchTalentDTO) error {
	err := service.Repository.Scan(uow, sessionTalent, repository.Table("batch_session_talents"),
		repository.Select("AVG(batch_session_talents.`average_rating`) AS `average_rating`"),
		repository.Filter("batch_session_talents.`tenant_id` = ? AND batch_session_talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_talents.`batch_id` = ? AND batch_session_talents.`talent_id` = ?",
			batchID, sessionTalent.TalentID))
	if err != nil {
		return err
	}
	return nil
}

func (service *BatchTalentService) getFeedbacksGivenToTalent(uow *repository.UnitOfWork, tenantID,
	batchID, talentID uuid.UUID, totalCount *uint) error {
	err := service.Repository.GetCountForTenant(uow, tenantID, bat.FacultyTalentBatchSessionFeedback{}, totalCount,
		repository.Filter("`batch_id` = ? AND `talent_id` = ?", batchID, talentID),
		repository.GroupBy("`batch_session_id`"))
	if err != nil {
		return err
	}
	return nil
}

// addSearchQueriesForBatchDetails adds all search queries if any when getAll is called
func (service *BatchTalentService) addSearchQueriesForBatchDetails(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("batch_sessions.`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addSearchQueries will add query processors based on queryparams.
func (service *BatchTalentService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	fmt.Println(" =================================== requestForm ->", requestForm)

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("batch_session_talents.`talent_id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	if batchSessionID, ok := requestForm["batchSessionID"]; ok {
		util.AddToSlice("batch_session_talents.`batch_session_id`", "= ?", "AND", batchSessionID, &columnNames, &conditions, &operators, &values)
	}

	if isPresent, ok := requestForm["isPresent"]; ok {
		util.AddToSlice("batch_session_talents.`is_present`", "= ?", "AND", isPresent, &columnNames, &conditions, &operators, &values)
	}

	if averageRating, ok := requestForm["averageRating"]; ok {
		util.AddToSlice("batch_session_talents.`average_rating`", "= ?", "AND", averageRating, &columnNames, &conditions, &operators, &values)
	}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("batch_talents.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	if suspensionDate, ok := requestForm["suspensionDate"]; ok {
		util.AddToSlice("batch_talents.`suspension_date`", "= ?", "AND", suspensionDate, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service BatchTalentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service BatchTalentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBatchExist returns error if there is no batch record in table for the given tenant.
func (service BatchTalentService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist returns error if there is no talent record in table for the given tenant.
func (service BatchTalentService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBatchTalentExist returns error if there is no batch talent record in table for the given tenant.
func (service BatchTalentService) doesBatchTalentExist(tenantID, batchTalentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.MappedTalent{},
		repository.Filter("`id` = ?", batchTalentID))
	if err := util.HandleError("Invalid batch talent ID", exists, err); err != nil {
		return err
	}
	return nil
}

// validateMappedForeignKeys validates all the foreign keys for the batch
func (service *BatchTalentService) validateMappedForeignKeys(talents *[]bat.MappedTalent, batchID, tenantID, credentialID uuid.UUID) error {
	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// Validate batchID
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	// // validate credentialID
	// err = service.doesCredentialExist(tenantID, credentialID)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }
	// validate talentID
	for _, talent := range *talents {
		exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.TalentDTO{}, repository.Filter("`id` = ?", talent.TalentID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		if !exists {
			return errors.NewValidationError("Talent not found")
		}
	}
	return nil
}

// checkFieldUniquenessInJSON checks if the slice already checkFieldUniquenessInJSON specified talent
func (service *BatchTalentService) checkFieldUniquenessInJSON(talents *[]bat.MappedTalent) (uuid.UUID, bool) {

	totalTalents := len(*talents)
	talentMap := make(map[uuid.UUID]uint, totalTalents)

	for _, talent := range *talents {

		talentMap[talent.TalentID] = talentMap[talent.TalentID] + 1
		if talentMap[talent.TalentID] > 1 {

			return talent.TalentID, true
		}
	}

	return uuid.Nil, false
}
