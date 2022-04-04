package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BatchService provide method like Add, Update, Delete, GetByID, GetAll for batch
// batchAssociationNames field will contain details about the sub-structs in batch for preload and other operations.
type BatchService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewBatchService creates a new instance of BatchService
func NewBatchService(db *gorm.DB, repo repository.Repository) *BatchService {

	return &BatchService{
		DB:         db,
		Repository: repo,
		associations: []string{
			"Course", "Course.Eligibility",
			"Course.Eligibility.Technologies",
			"Eligibility", "Eligibility.Technologies",
			// "Faculty",
			"SalesPerson",
			"Requirement", "Requirement.Branch", "Timing.Day",
			// "Timing",
		},
	}
}

// AddBatch adds new batch branch to database.
func (service *BatchService) AddBatch(batch *bat.Batch) error {

	// Extract IDs
	service.extractID(batch)

	// Checks all the foreign keys of the batch
	err := service.doForeignKeysExist(batch, batch.CreatedBy)
	if err != nil {
		return err
	}

	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Validates credential ID and checks if credentialID has permission to add batch
	// err = credentialService.ValidatePermission(batch.TenantID, batch.CreatedBy, "/batch/master", "add")
	// if err != nil {
	// 	return err
	// }

	// initializes totalStudents and totalIntake to zero if empty
	var totalStudents uint8 = 0
	if batch.TotalStudents == nil {
		batch.TotalStudents = &totalStudents
	}

	if batch.TotalIntake == nil {
		*batch.TotalStudents = 0
	}

	if batch.Eligibility != nil {
		batch.Eligibility.TenantID = batch.TenantID
		batch.Eligibility.CreatedBy = batch.CreatedBy
	}

	service.initializeBatchTiming(batch)

	// Starting transaction.
	uow := repository.NewUnitOfWork(service.DB, false)
	batch.Code, err = util.GenerateUniqueCode(uow.DB, batch.BatchName, "`code` = ?", &bat.Batch{})
	if err != nil {
		return err
	}

	// Add batch to the DB.
	err = service.Repository.Add(uow, batch)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBatch update the data of batch
func (service *BatchService) UpdateBatch(batch *bat.Batch) error {

	// Extract IDs
	service.extractID(batch)

	// Check if batch exists
	err := service.doesBatchExists(batch.TenantID, batch.ID)
	if err != nil {
		return err
	}

	// Validate batch code
	exists, err := repository.DoesRecordExistForTenant(service.DB, batch.TenantID, &bat.Batch{},
		repository.Filter("code=? AND `id` NOT IN (?) ", batch.Code, batch.ID))

	if err := util.HandleIfExistsError("Batch with similar code already exists", exists, err); err != nil {
		return err
	}
	// checks all the foreign keys of the batch
	err = service.doForeignKeysExist(batch, batch.UpdatedBy)
	if err != nil {
		return err
	}

	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Validates credential ID and checks if credentialID has permission to update batch
	// err = credentialService.ValidatePermission(batch.TenantID, batch.UpdatedBy, "/batch/master", "update")
	// if err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// }

	// check totalIntake
	if batch.TotalIntake != nil && batch.TotalStudents != nil {
		if *batch.TotalStudents > *batch.TotalIntake {
			return errors.NewValidationError("Total intake cannot be less than total students")
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// get batch details to set created_by field
	tempBatch := bat.Batch{}
	err = service.Repository.GetRecordForTenant(uow, batch.TenantID, &tempBatch, repository.Select("created_by"),
		repository.Filter("`id` = ?", batch.ID))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	batch.CreatedBy = tempBatch.CreatedBy

	// update eligibility for the batch
	err = service.updateBatchEligibility(uow, batch)
	if err != nil {
		uow.RollBack()
		return err
	}

	// update batch-timing
	err = service.updateBatchTiming(uow, batch, batch.UpdatedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// update batch details
	err = service.Repository.Save(uow, batch)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Update all the waiting list entries' isActive field to false if batch is inactive or batch status is finished.
	if !(*batch.IsActive) || (batch.Status == "Finished") {
		err = service.updateWaitingList(uow, batch)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// DeleteBatch delete the data of batch
func (service *BatchService) DeleteBatch(batch *bat.Batch) error {

	// Check if tenant exists
	err := service.doesTenantExists(batch.TenantID)
	if err != nil {
		return err
	}

	// Check if batch exists
	err = service.doesBatchExists(batch.TenantID, batch.ID)
	if err != nil {
		return err
	}

	// credentialService := genService.NewCredentialService(service.DB, service.Repository)

	// // Validates credential ID and checks if credentialID has permission to delete batch
	// err = credentialService.ValidatePermission(batch.TenantID, batch.DeletedBy, "/batch/master", "delete")
	// if err != nil {
	// 	return err
	// }

	// Start Transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// deletes all the talents from current batch
	err = service.deleteBatchTalent(uow, batch.ID, batch.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// deletes all the sessions for the current batch
	err = service.deleteBatchSession(uow, batch.ID, batch.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete all timings for specified batch
	err = service.deleteBatchTiming(uow, batch.ID, batch.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete all session-feedback for specified batch
	err = service.deleteFeedback(uow, batch.ID, batch.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// get batch record based on batchID
	tempBatch := bat.Batch{}
	err = service.Repository.GetRecordForTenant(uow, batch.TenantID, &tempBatch,
		repository.Filter("`id` = ?", batch.ID), repository.Select([]string{"`eligibility_id`"}))
	// repository.PreloadAssociations([]string{"Eligibility", "Eligibility.Technologies"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.deleteBatchAssociation(uow, &tempBatch, batch.DeletedBy)
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete batch
	err = service.Repository.UpdateWithMap(uow, bat.Batch{}, map[string]interface{}{
		"DeletedBy": batch.DeletedBy,
		"DeletedAt": time.Now(),
		"IsActive":  false,
	}, repository.Filter("`id` = ?", batch.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Update all the waiting list entries' isActive field to false.
	err = service.Repository.UpdateWithMap(uow, tal.WaitingList{}, map[string]interface{}{
		"UpdatedBy": batch.DeletedBy,
		"IsActive":  false,
	},
		repository.Filter("`batch_id`=?", batch.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllBatches returns all batches from database
func (service *BatchService) GetAllBatches(batches *[]bat.BatchDTO, tenantID uuid.UUID,
	parser *web.Parser, totalCount, totalTalents *int) error {

	// Check if tenant exists
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	limit, offset := parser.ParseLimitAndOffset()

	// get all batches ordered by `created_by DESC` order
	// "batches.`status` DESC", `status` and

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.PreloadWithCustomCondition(repository.Preload{
		Schema: "Timing", Queryprocessors: []repository.QueryProcessor{
			repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
			repository.Filter("days.`deleted_at` IS NULL AND batch_timing.`tenant_id` = ?", tenantID),
			repository.OrderBy("days.`order`")}}), repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.PreloadAssociations(service.associations),
		repository.OrderBy("batches.`created_at` DESC, batches.`is_active` DESC"),
		repository.Paginate(limit, offset, totalCount))

	err = service.Repository.GetAll(uow, batches, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	batchSessionsChannel := make(chan error, 1)
	batchValuesChannel := make(chan error, 1)

	for index := range *batches {
		// go service.getBatchSessionsCount(uow, &(*batches)[index].TotalSessionCount,
		// 	&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID, batchSessionsChannel)

		go service.getSessionCount(uow, &(*batches)[index].TotalSessionCount,
			&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID, batchSessionsChannel)

		// Added by sejal****************************************************************************
		// Range batches for getting extra fields.
		go service.getValuesForBatches(uow, &(*batches)[index], tenantID, batchValuesChannel)

		err = <-batchSessionsChannel
		if err != nil {
			uow.RollBack()
			return err
		}

		err = <-batchValuesChannel
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.Repository.GetAll(uow, &(*batches)[index].Faculty,
			repository.Join("INNER JOIN batch_modules ON faculties.id = batch_modules.faculty_id"+
				" AND faculties.tenant_id = batch_modules.tenant_id"),
			repository.Join("INNER JOIN batches ON batches.id = batch_modules.batch_id"),
			repository.Filter("batch_modules.batch_id=?", &(*batches)[index].ID),
			repository.GroupBy("faculties.`id`"))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	queryProcessors = []repository.QueryProcessor{}

	queryProcessors = append(queryProcessors, repository.Join("INNER JOIN batches ON batches.id = batch_talents.batch_id AND "+
		" batches.`tenant_id` = batch_talents.`tenant_id`"), repository.Filter("batch_talents.`tenant_id` = ?", tenantID))
	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.GroupBy("batch_talents.`talent_id`"))

	err = service.Repository.GetCount(uow, bat.MappedTalent{}, totalTalents, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetActiveBatchList returns listing of all the active batches
func (service *BatchService) GetActiveBatchList(batches *[]list.Batch, parser *web.Parser, tenantID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, repository.Filter("batches.`deleted_at` IS NULL AND batches.`is_active` = ? AND "+
		" batches.`tenant_id` = ?", true, tenantID))
	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)

	// get all batches ordered by `status` and `created_by DESC`
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batches, "batches.`created_at` DESC", queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// GetUpcomingBatches returns upcoming batches.
func (service *BatchService) GetUpcomingBatches(batches *[]bat.UpcomingBatch,
	tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}
	limit, offset := parser.ParseLimitAndOffset()
	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query processors for filter, pagination and preloads.
	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, service.addSearchQueries(parser.Form)...)

	queryProcessors = append(queryProcessors,
		repository.Filter("batches.`tenant_id`=? AND batches.`status`=?", tenantID, "upcoming"),
		repository.OrderBy("batches.`start_date` DESC"), repository.PreloadAssociations([]string{"Course"}),
		repository.Paginate(limit, offset, totalCount))

	// Get all requirements form database.
	err = service.Repository.GetAll(uow, batches, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *batches {

		// Count total talents who have been enrolled for this batch.
		err = service.Repository.Scan(uow, &(*batches)[index],
			repository.Table("batches"),
			repository.Select("COUNT(DISTINCT batch_talents.`talent_id`) AS total_enrolled"),
			repository.Join("INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
			repository.Filter("batches.`id` = ? AND batches.`tenant_id` = ?", (*batches)[index].ID, tenantID),
			repository.Filter("batch_talents.`tenant_id` = ?", tenantID),
			repository.Filter("batches.`deleted_at` IS NULL AND batch_talents.`deleted_at` IS NULL"))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetBatchDetails returns specific deatils of a single batch.
func (service *BatchService) GetBatchDetails(batch *bat.BatchDetails, tenantID uuid.UUID, parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exists.
	err = service.doesBatchExists(tenantID, batch.ID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get batch.
	err = service.Repository.GetRecord(uow, batch,
		repository.Filter("batches.`id`=?", batch.ID),
		repository.Filter("batches.`tenant_id`=?", tenantID),
		repository.Filter("batches.`deleted_at` IS NULL"),
		repository.PreloadAssociations([]string{"Course", "BatchTimings", "BatchTimings.Day"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get total talents of batch.
	var totalTalentsCount uint

	err = service.Repository.GetCount(uow, bat.MappedTalent{}, &totalTalentsCount,
		repository.Filter("batch_talents.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_talents.`batch_id`=?", batch.ID),
		repository.Filter("batch_talents.`is_active` = ?", true))
	if err != nil {
		uow.RollBack()
		return err
	}

	batch.TotalIntake = uint8(totalTalentsCount)

	// Get batch sessions.
	batchSessions := []bat.Session{}

	err = service.Repository.GetAll(uow, &batchSessions,
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.batch_session_id = batch_sessions.id"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_sessions.`batch_id` = ?", batch.ID),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.GroupBy("batch_sessions.`id`"),
		repository.OrderBy("batch_sessions.date"))
	if err != nil {
		return err
	}

	if len(batchSessions) > 0 {
		batch.StartDate = batchSessions[0].Date
		batch.EndDate = &batchSessions[len(batchSessions)-1].Date
	}

	batch.TotalSessionsCount = uint(len(batchSessions))

	// Get total hours of batch session.
	var sessionTotalHours bat.SessionTotalHours

	err = service.Repository.Scan(uow, &sessionTotalHours,
		repository.Table("batch_session_topics"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_hours`"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ?", tenantID),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`deleted_at` IS NULL", batch.ID))
	if err != nil {
		return err
	}

	batch.TotalHours = sessionTotalHours.TotalHours

	// Get faculty.
	var tempBatchSessions []bat.BatchDetailsSessionDTO

	err = service.Repository.GetAll(uow, &tempBatchSessions,
		repository.Join("INNER JOIN batches ON batch_sessions.batch_id = batches.id"),
		repository.Filter("batch_sessions.`batch_id` = ?", batch.ID),
		repository.Filter("batches.`tenant_id` = ? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.GroupBy("batch_sessions.`faculty_id`"),
		repository.PreloadAssociations([]string{"Faculty"}))
	if err != nil {
		return err
	}

	for i := range tempBatchSessions {
		if tempBatchSessions[i].FacultyID != nil {
			batch.Faculty = append(batch.Faculty, *tempBatchSessions[i].Faculty)
		}
	}

	// Get total sessions completed.
	var totalSessionsCompletedCount uint

	err = service.getBatchTotalCompletedSessions(uow, &totalSessionsCompletedCount, tenantID, batch.ID, parser)
	if err != nil {
		uow.RollBack()
		return err
	}
	batch.TotalSessionsCompleted = totalSessionsCompletedCount

	// Get total hours completed.
	var sessionTotalCompletedHours bat.SessionTotalHours

	err = service.getBatchSessionTotalCompletedHours(uow, tenantID, batch.ID, &sessionTotalCompletedHours, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	batch.TotalCompletedHours = sessionTotalCompletedHours.TotalHours

	uow.Commit()
	return nil
}

// GetBatch returns particular batch by ID
func (service *BatchService) GetBatch(batch *bat.Batch) error {
	// Check if tenant exists
	err := service.doesTenantExists(batch.TenantID)
	if err != nil {
		return err
	}
	// Check if batch exists
	err = service.doesBatchExists(batch.TenantID, batch.ID)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, batch.TenantID, batch.ID, batch,
		repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return err
	}
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================
// func (service *BatchService) getBatchSessionTotalHours(uow *repository.UnitOfWork,
// 	tenantID, batchID uuid.UUID, sessionTotalHours *bat.SessionTotalHours) error {

// 	err := service.Repository.Scan(uow, sessionTotalHours, repository.Table("batch_session_topics"),
// 		repository.Select("SUM(batch_session_topics.total_time) AS total_hours"),
// 		repository.Filter("batch_session_topics.tenant_id = ?", tenantID),
// 		repository.Filter("batch_session_topics.batch_id = ? AND batch_session_topics.deleted_at IS NULL", batchID))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (service *BatchService) updateWaitingList(uow *repository.UnitOfWork, batch *bat.Batch) error {
	err := service.Repository.UpdateWithMap(uow, tal.WaitingList{}, map[string]interface{}{
		"UpdatedBy": batch.UpdatedBy,
		"IsActive":  false,
	},
		repository.Filter("`batch_id`=?", batch.ID))
	if err != nil {
		return err
	}
	return nil
}

// updateBatchAssociation Update batch's Dependencies
func (service *BatchService) updateBatchAssociation(uow *repository.UnitOfWork, batch *bat.Batch, credentialID uuid.UUID) error {

	if batch.Eligibility != nil {
		if err := service.Repository.ReplaceAssociations(uow, batch.Eligibility, "Technologies", batch.Eligibility.Technologies); err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
	}

	return nil
}

// deleteBatchAssociation soft deletes all associations of batch
func (service *BatchService) deleteBatchAssociation(uow *repository.UnitOfWork, batch *bat.Batch, credentialID uuid.UUID) error {

	if batch.EligibilityID != nil {
		err := service.Repository.UpdateWithMap(uow, bat.Eligibility{}, map[interface{}]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		}, repository.Filter("`id`=?", batch.EligibilityID))
		if err != nil {
			return err
		}
		// if len((*batch.Eligibility).Technologies) > 0 {
		// 	if err := service.Repository.RemoveAssociations(uow, batch.Eligibility, "Technologies", batch.Eligibility.Technologies); err != nil {
		// 		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		// 	}
		// }
	}

	return nil
}

// updateBatchEligibility
func (service *BatchService) updateBatchEligibility(uow *repository.UnitOfWork, batch *bat.Batch) error {

	// get batch record
	tempBatch := &bat.Batch{}
	err := service.Repository.GetAllForTenant(uow, batch.TenantID, tempBatch,
		repository.Filter("`id`=?", batch.ID),
		repository.PreloadAssociations([]string{"Eligibility", "Eligibility.Technologies"}))
	if err != nil {
		return err
	}
	// previously no eligibility exists
	if tempBatch.EligibilityID == nil && batch.Eligibility == nil {
		return nil
	}

	// previously eligibility exists but now eligibility is removed
	if batch.Eligibility == nil {

		if len((*tempBatch.Eligibility).Technologies) > 0 {
			if err := service.Repository.RemoveAssociations(uow, tempBatch.Eligibility, "Technologies",
				tempBatch.Eligibility.Technologies); err != nil {
				return err
			}
		}

		err = service.deleteBatchAssociation(uow, tempBatch, batch.UpdatedBy)
		if err != nil {
			return err
		}

		// set batch.eligibilityID to null
		err = service.Repository.UpdateWithMap(uow, &bat.Batch{}, map[interface{}]interface{}{
			"EligibilityID": nil,
		}, repository.Filter("`id` = ?", batch.ID))
		if err != nil {
			return err
		}
		return nil
	}

	if tempBatch.EligibilityID != nil {
		batch.Eligibility.TenantID = batch.TenantID
		batch.Eligibility.UpdatedBy = batch.UpdatedBy

		// update batch associations
		err = service.updateBatchAssociation(uow, batch, batch.UpdatedBy)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
		}
		batch.Eligibility.Technologies = nil

		err = service.Repository.Update(uow, batch.Eligibility)
		if err != nil {
			return err
		}

		batch.EligibilityID = &batch.Eligibility.ID
		batch.Eligibility = nil
		return nil
	}

	batch.Eligibility.TenantID = batch.TenantID
	batch.Eligibility.CreatedBy = batch.UpdatedBy
	err = service.Repository.Add(uow, batch.Eligibility)
	if err != nil {
		return err
	}
	batch.EligibilityID = &batch.Eligibility.ID
	batch.Eligibility = nil

	return nil
}

func (service *BatchService) updateBatchTiming(uow *repository.UnitOfWork,
	batch *bat.Batch, credentialID uuid.UUID) error {

	batchTimingMap := make(map[uuid.UUID]uint)
	batchTimings := &batch.Timing
	tempBatchTimings := &[]bat.Timing{}

	err := service.Repository.GetAllForTenant(uow, batch.TenantID, tempBatchTimings,
		repository.Filter("batch_id=?", batch.ID))
	if err != nil {
		return err
	}

	// populate all entries for existing batch timing (existing)
	for _, tempBatchTime := range *tempBatchTimings {
		batchTimingMap[tempBatchTime.ID] = batchTimingMap[tempBatchTime.ID] + 1
	}

	for _, batchTime := range *batchTimings {

		if util.IsUUIDValid(batchTime.ID) {
			batchTimingMap[batchTime.ID] = batchTimingMap[batchTime.ID] + 1
		} else { // add new batch-timing
			batchTime.CreatedBy = credentialID
			batchTime.BatchID = batch.ID
			batchTime.TenantID = batch.TenantID
			batchTime.DayID = batchTime.Day.ID
			batchTime.Day = nil
			err = service.Repository.Add(uow, &batchTime)
			if err != nil {
				return err
			}
		}

		// update existing records
		if batchTimingMap[batchTime.ID] > 1 {
			batchTime.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &batchTime)
			if err != nil {
				return err
			}
			batchTimingMap[batchTime.ID] = 0
		}
	}

	for _, tempBatchTime := range *tempBatchTimings {
		if batchTimingMap[tempBatchTime.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, bat.Timing{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempBatchTime.ID))
			if err != nil {
				return err
			}
		}
	}

	batch.Timing = nil
	return nil
}

// deleteBatchSession deletes all the sessions for the current batch
func (service *BatchService) deleteBatchSession(uow *repository.UnitOfWork,
	batchID, credentialID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, bat.Session{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.UpdateWithMap(uow, bat.SessionTopic{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.UpdateWithMap(uow, bat.TopicAssignment{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ?", batchID))
	if err != nil {
		return err
	}

	return nil
}

// deleteBatchTalent deletes all the talents from the current batch
func (service *BatchService) deleteBatchTalent(uow *repository.UnitOfWork,
	batchID, credentialID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, bat.MappedTalent{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}
	return nil
}

// deleteBatchTiming deletes all the timings for specified batch
func (service *BatchService) deleteBatchTiming(uow *repository.UnitOfWork,
	batchID, credentialID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, bat.Timing{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	return nil
}

// deleteFeedback deletes all the session-feedbacks for specified batch
func (service *BatchService) deleteFeedback(uow *repository.UnitOfWork,
	batchID, credentialID uuid.UUID) error {

	err := service.Repository.UpdateWithMap(uow, bat.FacultyTalentBatchSessionFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.UpdateWithMap(uow, bat.FacultyTalentFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.UpdateWithMap(uow, bat.TalentFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.UpdateWithMap(uow, bat.TalentBatchSessionFeedback{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=?", batchID))
	if err != nil {
		return err
	}
	return nil
}

// getSessionCount will return count of total topics and completed topics for the specified batch
func (service *BatchService) getSessionCount(uow *repository.UnitOfWork,
	totalSession, completedSession *uint, tenantID, batchID uuid.UUID, batchSessionsChannel chan error) {

	err := service.Repository.GetCount(uow, bat.Session{}, totalSession,
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`batch_session_id` = batch_sessions.`id` AND "+
			"batch_session_topics.`tenant_id` = batch_sessions.`tenant_id`"),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`batch_id` = ?", batchID), repository.GroupBy("batch_session_topics.`batch_session_id`"))
	if err != nil {
		batchSessionsChannel <- err
	}

	err = service.Repository.GetCountForTenant(uow, tenantID, bat.Session{}, completedSession,
		repository.Filter("`batch_sessions`.`is_session_taken` = ? AND `batch_sessions`.`batch_id` = ?", 1, batchID))
	if err != nil {
		batchSessionsChannel <- err
	}

	batchSessionsChannel <- nil
}

// getBatchSessionsCount will return count of total sessions and completed sessions for the specified batch
// func (service *BatchService) getBatchSessionsCount(uow *repository.UnitOfWork,
// 	totalSession, completedSession *uint, tenantID, batchID uuid.UUID, batchSessionsChannel chan error) {

// 	err := service.Repository.GetCount(uow, bat.MappedSession{},
// 		totalSession, repository.Join("INNER JOIN course_sessions ON course_sessions.`id`=batch_sessions.`course_session_id`"+
// 			" AND course_sessions.`tenant_id`=batch_sessions.`tenant_id`"),
// 		repository.Filter("course_sessions.`session_id` IS NULL AND course_sessions.`deleted_at` IS NULL"),
// 		repository.Filter("batch_sessions.`batch_id`=? AND batch_sessions.`tenant_id`=?", batchID, tenantID))
// 	if err != nil {
// 		// return err
// 		batchSessionsChannel <- err
// 	}

// 	// go service.getTotalSessionCount(uow, totalSession, tenantID, batchID, batchSessionsChannel)

// 	err = service.Repository.GetCount(uow, bat.MappedSession{},
// 		completedSession, repository.Join("INNER JOIN course_sessions ON course_sessions.`id`=batch_sessions.`course_session_id`"+
// 			" AND course_sessions.`tenant_id`=batch_sessions.`tenant_id`"),
// 		repository.Filter("course_sessions.`session_id` IS NULL AND course_sessions.`deleted_at` IS NULL"),
// 		repository.Filter("batch_sessions.`is_completed` = true"),
// 		repository.Filter("batch_sessions.`batch_id`=? AND batch_sessions.`tenant_id`=?", batchID, tenantID))
// 	if err != nil {
// 		// return err
// 		batchSessionsChannel <- err
// 	}

// 	// return nil
// 	batchSessionsChannel <- nil
// }

// func (service *BatchService) getTotalSessionCount(uow *repository.UnitOfWork,
// 	totalSession *uint, tenantID, batchID uuid.UUID, batchSessionsChannel chan error) {

// 	err := service.Repository.GetCount(uow, bat.MappedSession{},
// 		totalSession, repository.Join("INNER JOIN course_sessions ON course_sessions.`id`=batch_sessions.`course_session_id`"+
// 			" AND course_sessions.`tenant_id`=batch_sessions.`tenant_id`"),
// 		repository.Filter("course_sessions.`session_id` IS NULL AND course_sessions.`deleted_at` IS NULL"),
// 		repository.Filter("batch_sessions.`batch_id`=? AND batch_sessions.`tenant_id`=?", batchID, tenantID))
// 	if err != nil {
// 		// return err
// 		batchSessionsChannel <- err
// 	}
// 	batchSessionsChannel <- nil
// }

// doForeignKeysExist validates all the foreign keys for the batch
func (service *BatchService) doForeignKeysExist(batch *bat.Batch, credentialID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExists(batch.TenantID)
	if err != nil {
		return err
	}

	// check if batch name exist.
	err = service.doesBatchNameExists(batch.TenantID, batch.ID, batch.BatchName)
	if err != nil {
		return err
	}

	// Validate salesPersonID
	err = service.doesSalesPersonExist(batch.TenantID, batch.SalesPersonID)
	if err != nil {
		return err
	}

	// Validate facultyID
	// err = service.doesFacultyExist(batch.TenantID, batch.FacultyID)
	// if err != nil {
	// 	return err
	// }

	// Validate courseID
	err = service.doesCourseExist(batch.TenantID, batch.CourseID)
	if err != nil {
		return err
	}

	// check if requirementID is valid
	if batch.RequirementID != nil {
		err = service.doesRequirementExists(batch.TenantID, *batch.RequirementID)
		if err != nil {
			return err
		}
	}

	return nil
}

// doesTenantExists validates tenantID
func (service *BatchService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBatchExists validates batchID
func (service *BatchService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{}, repository.Filter("`id`=?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// Validate duplicate batch name
func (service *BatchService) doesBatchNameExists(tenantID, batchID uuid.UUID, batchName string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` != ? AND `batch_name`=?", batchID, batchName))
	if err := util.HandleIfExistsError("Batch name already exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSalesPersonExist validates salesPersonID
func (service *BatchService) doesSalesPersonExist(tenantID, salesPersonID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{}, repository.Filter("`id`=?", salesPersonID))
	if err := util.HandleError("SalesPerson not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFacultyExist validates facultyID
// func (service *BatchService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fclt.Faculty{}, repository.Filter("`id`=?", facultyID))
// 	if err := util.HandleError("Faculty not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// doesCourseExist validates courseID
func (service *BatchService) doesCourseExist(tenantID, courseID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Course{}, repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Course not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSessionExists checks if session id's are valid
func (service *BatchService) doesRequirementExists(tenantID, requirementID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Requirement{},
		repository.Filter("`id`=?", requirementID))
	if err := util.HandleError("Invalid requirement ID", exists, err); err != nil {
		return err
	}
	return nil
}

// Extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the college branch entity is being added or updated.
func (service *BatchService) extractID(batch *bat.Batch) {

	batch.SalesPersonID = batch.SalesPerson.ID

	// batch.FacultyID = batch.Faculty.ID

	batch.CourseID = batch.Course.ID

	if batch.Requirement != nil {
		batch.RequirementID = &batch.Requirement.ID
	}

}

func (service *BatchService) initializeBatchTiming(batch *bat.Batch) {
	for index := range batch.Timing {
		batch.Timing[index].CreatedBy = batch.CreatedBy
		batch.Timing[index].TenantID = batch.TenantID
		batch.Timing[index].DayID = batch.Timing[index].Day.ID
		batch.Timing[index].Day = nil
	}
}

func (service *BatchService) getBatchTotalCompletedSessions(uow *repository.UnitOfWork, totalCount *uint,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	// err := service.Repository.GetCount(uow, bat.MappedSession{}, totalCount,
	// 	repository.Join("INNER JOIN course_sessions ON course_sessions.`id` = batch_sessions.`course_session_id` AND "+
	// 		"course_sessions.`tenant_id` = batch_sessions.`tenant_id`"),
	// 	repository.Filter("batch_sessions.`tenant_id` = ?", tenantID),
	// 	repository.Filter("batch_sessions.`deleted_at` IS NULL AND course_sessions.`deleted_at` IS NULL"),
	// 	repository.Filter("course_sessions.`session_id` IS NULL AND batch_sessions.`batch_id` = ?", batchID),
	// 	repository.Filter("batch_sessions.`is_completed` = 1"))
	err := service.Repository.GetCountForTenant(uow, tenantID, bat.Session{}, totalCount,
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_sessions.`batch_id`=?", batchID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_sessions.`tenant_id`=?", tenantID),
		repository.Filter("batch_sessions.`is_session_taken`=?", true))
	if err != nil {
		return err
	}

	return nil
}

func (service *BatchService) getBatchSessionTotalCompletedHours(uow *repository.UnitOfWork,
	tenantID, batchID uuid.UUID, sessionTotalHours *bat.SessionTotalHours, parser *web.Parser) error {

	err := service.Repository.Scan(uow, sessionTotalHours, repository.Table("batch_session_topics"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_hours`"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		service.addSearchQueriesForBatchDetails(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ?", tenantID),
		repository.Filter("batch_session_topics.`is_completed` = ?", 1),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`deleted_at` IS NULL", batchID))
	if err != nil {
		return err
	}
	return nil
}

// addBatchSearchQueries adds all search queries if any when getAll is called
func (service *BatchService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("batches.`id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}

	if !util.IsEmpty(requestForm.Get("batchName")) {
		util.AddToSlice("`batches`.`batch_name`", "LIKE ?", "AND", "%"+requestForm.Get("batchName")+"%", &columnNames, &conditions, &operators, &values)
	}

	if batchStatus, ok := requestForm["batchStatus"]; ok {
		util.AddToSlice("`status`", "IN(?)", "AND", batchStatus, &columnNames, &conditions, &operators, &values)
	}

	if batchObjective, ok := requestForm["batchObjective"]; ok {
		util.AddToSlice("`batch_objective`", "= ?", "AND", batchObjective, &columnNames, &conditions, &operators, &values)
	}

	if startDate, ok := requestForm["startDate"]; ok {
		util.AddToSlice("`start_date`", ">= ?", "AND", startDate, &columnNames, &conditions, &operators, &values)
	}

	if estimatedEndDate, ok := requestForm["estimatedEndDate"]; ok {
		util.AddToSlice("`estimated_end_date`", "<= ?", "AND", estimatedEndDate, &columnNames, &conditions, &operators, &values)
	}

	if courseID, ok := requestForm["courseID"]; ok {
		util.AddToSlice("`course_id`", "= ?", "AND", courseID, &columnNames, &conditions, &operators, &values)
	}

	if salesPersonID, ok := requestForm["salesPersonID"]; ok {
		util.AddToSlice("`sales_person_id`", "= ?", "AND", salesPersonID, &columnNames, &conditions, &operators, &values)
	}

	if facultyID, ok := requestForm["facultyID"]; ok {
		if util.IsEmpty(requestForm.Get("isViewAllBatches")) || requestForm.Get("isViewAllBatches") == "0" {
			// util.AddToSlice("`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
			queryProcessors = append(queryProcessors, repository.Join("INNER JOIN `batch_modules` ON batch_modules.`batch_id` = batches.`id` "+
				"AND batch_modules.`tenant_id` = batches.`tenant_id`"),
				repository.Filter("batch_modules.`faculty_id` = ? AND batch_modules.`deleted_at` IS NULL", facultyID))
		}
	}

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("batches.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("batches.`id`"))

	return queryProcessors
	// return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addSearchQueriesForBatchDetails adds all search queries if any when getAll is called
func (service *BatchService) addSearchQueriesForBatchDetails(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("batch_sessions.`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// getValuesForBatches gets values for batches by firing individual query for each batch.
func (service *BatchService) getValuesForBatches(uow *repository.UnitOfWork, batch *bat.BatchDTO, tenantID uuid.UUID,
	batchValuesChannel chan error) {
	//********************************************APPLICANTS*********************************************************************

	// Create bucket for applicants.
	var waitingListCount uint16 = 0

	// Get all count of applicants form database.
	// If batch is active and not finished then get only active waiting list entries.
	if *batch.IsActive && *batch.Status != "Finished" {
		err := service.Repository.GetCountForTenant(uow, tenantID, tal.WaitingList{}, &waitingListCount,
			repository.Filter("`batch_id`=? AND is_active=?", batch.ID, true))
		if err != nil {
			// return err
			batchValuesChannel <- err
		}
	} else {
		err := service.Repository.GetCountForTenant(uow, tenantID, tal.WaitingList{}, &waitingListCount,
			repository.Filter("`batch_id`=?", batch.ID))
		if err != nil {
			// return err
			batchValuesChannel <- err
		}
	}

	// If no applicants then dont assign any value to applicants fields.
	if waitingListCount == 0 {
		// return nil
		batchValuesChannel <- nil
	}

	// Give waiting list to requirement.
	batch.TotalApplicants = &waitingListCount
	// return nil
	batchValuesChannel <- nil
}

// // GetAllBatchesForFaculty returns batches for particular faculty
// func (service *BatchService) GetAllBatchesForFaculty(batches *[]bat.BatchDTO, form url.Values, tenantID, facultyID uuid.UUID,
// 	limit, offset int, totalCount, totalTalents *int) error {

// 	// check if tenant exists
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batches, "batches.`status` DESC",
// 		repository.Filter("faculty_id=?", facultyID), service.addBatchSearchQueryParams(form),
// 		repository.OrderBy("batches.`created_by` DESC, batches.`is_active` DESC"), repository.PreloadAssociations(service.associations),
// 		repository.Paginate(limit, offset, totalCount))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	for index := range *batches {
// 		err = service.Repository.GetAllInOrder(uow, &(*batches)[index].Timing, "days.`order`",
// 			repository.PreloadAssociations([]string{"Day"}), repository.Join("INNER JOIN days ON days.`id`=batch_timing.`day_id`"),
// 			repository.Filter("batch_timing.`batch_id`=?", (*batches)[index].ID),
// 			repository.Filter("batch_timing.`tenant_id`=? AND batch_timing.`deleted_at` IS NULL", tenantID),
// 			repository.Filter("days.`tenant_id`=? AND days.`deleted_at` IS NULL", tenantID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}

// 		err = service.getBatchSessionsCount(uow, &(*batches)[index].TotalSessionCount,
// 			&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	err = service.getTotalTalents(uow, totalTalents, tenantID, form)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// // GetAllBatchesForSalesPerson returns batches for particular salesperson
// func (service *BatchService) GetAllBatchesForSalesPerson(batches *[]bat.BatchDTO, form url.Values, tenantID, salesPersonID uuid.UUID,
// 	limit, offset int, totalCount, totalTalents *int) error {
// 	// Check if tenant exists
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}
// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batches, "batches.`status` DESC",
// 		repository.Filter("sales_person_id=?", salesPersonID), service.addBatchSearchQueryParams(form),
// 		repository.OrderBy("batches.`created_by` DESC, batches.`is_active` DESC"), repository.PreloadAssociations(service.associations),
// 		repository.Paginate(limit, offset, totalCount))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	for index := range *batches {
// 		err = service.Repository.GetAllInOrder(uow, &(*batches)[index].Timing, "days.`order`",
// 			repository.PreloadAssociations([]string{"Day"}), repository.Join("INNER JOIN days ON days.`id`=batch_timing.`day_id`"),
// 			repository.Filter("batch_timing.`batch_id`=?", (*batches)[index].ID),
// 			repository.Filter("batch_timing.`tenant_id`=? AND batch_timing.`deleted_at` IS NULL", tenantID),
// 			repository.Filter("days.`tenant_id`=? AND days.`deleted_at` IS NULL", tenantID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}

// 		err = service.getBatchSessionsCount(uow, &(*batches)[index].TotalSessionCount,
// 			&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	err = service.getTotalTalents(uow, totalTalents, tenantID, form)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// // SearchBatches returns all batches from database
// func (service *BatchService) SearchBatches(batches *[]*bat.Batch, batchSearch *bat.Search, tenantID uuid.UUID,
// 	limit, offset int, totalCount *int) error {
// 	// Check if tenant exists
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Start Transaction.
// 	uow := repository.NewUnitOfWork(service.DB, true)
// 	err = service.Repository.GetAllForTenant(uow, tenantID, batches, service.addBatchSearchQueries(batchSearch),
// 		repository.PreloadAssociations(service.associations), repository.Paginate(limit, offset, totalCount))
// 	if err != nil {
// 		uow.RollBack()
// 		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
// 	}

// 	for index := range *batches {
// 		err = service.Repository.GetAllInOrder(uow, &(*batches)[index].Timing, "days.`order`",
// 			repository.PreloadAssociations([]string{"Day"}), repository.Join("INNER JOIN days ON days.`id`=batch_timing.`day_id`"),
// 			repository.Filter("batch_timing.`batch_id`=?", (*batches)[index].ID),
// 			repository.Filter("batch_timing.`tenant_id`=? AND batch_timing.`deleted_at` IS NULL", tenantID),
// 			repository.Filter("days.`tenant_id`=? AND days.`deleted_at` IS NULL", tenantID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}

// 		err = service.getBatchSessionsCount(uow, &(*batches)[index].TotalSessionCount,
// 			&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	return nil
// }

// returns error if there is no credential record in table for the given tenant.
// func (service *BatchService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
// 		repository.Filter("`id`=?", credentialID))
// 	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // getTotalTalents will return total number of talents assigned to batches
// func (service *BatchService) getTotalTalents(uow *repository.UnitOfWork,
// 	totalTalents *int, tenantID uuid.UUID, requestForm url.Values) error {

// 	var totalCount int

// 	tempBatches := []bat.Batch{}
// 	err := service.Repository.GetAllForTenant(uow, tenantID, &tempBatches,
// 		service.addBatchSearchQueryParams(requestForm))
// 	if err != nil {
// 		return err
// 	}

// 	for _, batch := range tempBatches {
// 		err := service.Repository.GetCountForTenant(uow, tenantID, bat.MappedTalent{},
// 			&totalCount, repository.Filter("`batch_id` = ?", batch.ID), repository.GroupBy("`talent_id`"))
// 		if err != nil {
// 			return err
// 		}
// 		*totalTalents += totalCount
// 	}

// 	return nil
// }

// AddTalentsToBatch adds student to a particular batch
// func (service *BatchService) AddTalentsToBatch(talents *[]bat.MappedTalent, tenantID, batchID, credentialID uuid.UUID) error {

// 	// Validates all the foreign keys for batch_talents table which are talentID, batchID, tenantID and credentialID
// 	err := service.validateMappedForeignKeys(talents, batchID, tenantID, credentialID)
// 	if err != nil {
// 		return err
// 	}

// 	credentialService := genService.NewCredentialService(service.DB, service.Repository)

// 	// Validates credential ID and checks if credentialID has permission to add talent to batch
// 	// Firstly validates the credential ID and then gets the permissions for the role of that credential
// 	// If credentialID does not have permission to add then an error will be returned
// 	err = credentialService.ValidatePermission(tenantID, credentialID, "/batch/master", "add")
// 	if err != nil {
// 		return err
// 	}

// 	batch := bat.Batch{}
// 	batch.ID = batchID

// 	// Start Transcation.
// 	uow := repository.NewUnitOfWork(service.DB, false)

// 	// Gets the record of the batch
// 	err = service.Repository.GetRecordForTenant(uow, tenantID, &batch)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Checks if batch can be allocated to more students
// 	if batch.TotalStudents != nil && batch.TotalIntake != 0 {
// 		if batch.TotalIntake == *batch.TotalStudents {
// 			return errors.NewValidationError("Current batch is full")
// 		}
// 		if batch.TotalIntake < *batch.TotalStudents+uint8(len(*talents)) {
// 			return errors.NewValidationError("Only " + strconv.Itoa(int(batch.TotalIntake-*batch.TotalStudents)) + " slot left")
// 		}
// 	}

// 	duplicateTalent := &tal.Talent{}

// 	// Check if JSON has duplicate entries
// 	talentID, found := service.checkFieldUniquenessInJSON(talents)
// 	if found {
// 		err := service.Repository.GetRecordForTenant(uow, tenantID, duplicateTalent, repository.Filter("`id` = ?", talentID),
// 			repository.Select([]string{"`first_name`", "`last_name`"}))
// 		if err != nil {
// 			return err
// 		}
// 		return errors.NewValidationError(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " selected multiple times")
// 	}

// 	for _, talent := range *talents {

// 		// Check if talent is already assigned with same batch
// 		exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &bat.MappedTalent{},
// 			repository.Filter("batch_id=? AND talent_id=?", batchID, talent.TalentID))
// 		if err != nil {
// 			return err
// 		}
// 		if exists {
// 			err := service.Repository.GetRecordForTenant(uow, tenantID, duplicateTalent, repository.Filter("`id` = ?", talent.TalentID),
// 				repository.Select([]string{"`first_name`", "`last_name`"}))
// 			if err != nil {
// 				return err
// 			}
// 			return errors.NewValidationError(duplicateTalent.FirstName + " " + duplicateTalent.LastName + " is already added to this batch")
// 		}

// 		// Assign talent to the specified batch
// 		talent.BatchID = batchID
// 		talent.TenantID = tenantID
// 		talent.CreatedBy = credentialID
// 		talent.IsActive = true

// 		err = service.Repository.Add(uow, &talent)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}

// 	}

// 	// Update totalStudents in batch
// 	err = service.Repository.UpdateWithMap(uow, batch, map[interface{}]interface{}{
// 		"TotalStudents": *batch.TotalStudents + uint8(len(*talents)),
// 	})
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	uow.Commit()
// 	return nil
// }

// // DeleteTalentFromBatch soft deletes talent from the batch (i.e. soft deletes record from batch_talents table)
// func (service *BatchService) DeleteTalentFromBatch(tenantID, credentialID, talentID, batchID uuid.UUID) error {

// 	// Check if tenant exists
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	credentialService := genService.NewCredentialService(service.DB, service.Repository)

// 	// Validates credential ID and checks if credentialID has permission to delete talent from batch
// 	// Firstly validates the credential ID and then gets the permissions for the role of that credential
// 	// If credentialID does not have permission to delete then an error will be returned
// 	err = credentialService.ValidatePermission(tenantID, credentialID, "/batch/master", "delete")
// 	if err != nil {
// 		return err
// 	}

// 	// Check if batch exists
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent exist
// 	err = service.doesTalentExist(tenantID, talentID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if talent exist for specified batch
// 	err = service.doesTalentExistForBatch(tenantID, talentID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	// Start Transaction.
// 	uow := repository.NewUnitOfWork(service.DB, false)

// 	err = service.deleteTalentFromFeedback(uow, tenantID, credentialID, talentID, batchID)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Update DeletedBy in batch_talents
// 	err = service.Repository.UpdateWithMap(uow, &bat.MappedTalent{}, map[interface{}]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Update totalStudents in batch
// 	batch := bat.Batch{}

// 	// GetBatch for current batchID
// 	err = service.Repository.GetRecordForTenant(uow, tenantID, &batch,
// 		repository.Filter("`id` = ?", batchID), repository.Select([]string{"`total_students`"}))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Update totalStudents for current batchID
// 	if batch.TotalStudents != nil {
// 		// fmt.Println("****** stud ->", *batch.TotalStudents)
// 		err = service.Repository.UpdateWithMap(uow, &bat.Batch{}, map[interface{}]interface{}{
// 			"TotalStudents": *batch.TotalStudents - 1,
// 		}, repository.Filter("`id` = ?", batchID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	uow.Commit()
// 	return nil
// }

// deleteTalentFromFeedback will delete all the records were talent is added for batch.
// func (service *BatchService) deleteTalentFromFeedback(uow *repository.UnitOfWork, tenantID, credentialID,
// 	talentID, batchID uuid.UUID) error {

// 	// delete talent from faculty-batch-session-feedback.
// 	err := service.Repository.UpdateWithMap(uow, new(bat.FacultyTalentBatchSessionFeedback), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	// delete talent from faculty-batch-feedback.
// 	err = service.Repository.UpdateWithMap(uow, new(bat.FacultyTalentFeedback), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	// delete talent from talent-batch-session-feedback.
// 	err = service.Repository.UpdateWithMap(uow, new(bat.TalentBatchSessionFeedback), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	// delete talent from talent-batch-feedback.
// 	err = service.Repository.UpdateWithMap(uow, new(bat.TalentFeedback), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	// delete talent from aha-moment.
// 	err = service.Repository.UpdateWithMap(uow, new(bat.AhaMoment), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	// delete talent from aha-moment-response.
// 	err = service.Repository.UpdateWithMap(uow, new(bat.AhaMomentResponse), map[string]interface{}{
// 		"DeletedBy": credentialID,
// 		"DeletedAt": time.Now(),
// 	}, repository.Filter("`talent_id` = ? AND `batch_id` = ?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// validateMappedForeignKeys validates all the foreign keys for the batch
// func (service *BatchService) validateMappedForeignKeys(talents *[]bat.MappedTalent, batchID, tenantID, credentialID uuid.UUID) error {

// 	// Check if tenant exists
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Validate batchID
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	// validate talentID
// 	for _, talent := range *talents {
// 		exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.TalentDTO{}, repository.Filter("`id` = ?", talent.TalentID))
// 		if err != nil {
// 			return err
// 		}
// 		if !exists {
// 			return errors.NewValidationError("Talent not found")
// 		}
// 	}

// 	return nil
// }

// // doesTalentExist validates TalentID
// func (service *BatchService) doesTalentExist(tenantID, talentID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{}, repository.Filter("`id`=?", talentID))
// 	if err := util.HandleError("Talent not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // doesTalentExistForBatch checks if talent is assigned to specified batch
// func (service *BatchService) doesTalentExistForBatch(tenantID, talentID, batchID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.MappedTalent{},
// 		repository.Filter("talent_id=? AND batch_id=?", talentID, batchID))
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return errors.NewValidationError("Talent not found in this batch")
// 	}
// 	return nil
// }

// // checkFieldUniquenessInJSON checks if the slice already checkFieldUniquenessInJSON specified talent
// func (service *BatchService) checkFieldUniquenessInJSON(talents *[]bat.MappedTalent) (uuid.UUID, bool) {

// 	totalTalents := len(*talents)
// 	talentMap := make(map[uuid.UUID]uint, totalTalents)

// 	for _, talent := range *talents {

// 		talentMap[talent.TalentID] = talentMap[talent.TalentID] + 1
// 		if talentMap[talent.TalentID] > 1 {

// 			return talent.TalentID, true
// 		}
// 	}

// 	return uuid.Nil, false
// }
