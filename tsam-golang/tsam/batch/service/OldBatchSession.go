package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// OBatchSessionService provide method like Add, Update, Delete, GetByID, GetAll for batch
type OBatchSessionService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// OldBatchSessionService creates a new instance of BatchSessionService
func OldBatchSessionService(db *gorm.DB, repo repository.Repository) *OBatchSessionService {
	return &OBatchSessionService{
		DB:         db,
		Repository: repo,
		association: []string{
			// "Session","Session.SubSessions",
			"Session.Resources", "Session.SubSessions.Resources",
		},
	}
}

// AddSessionForBatch will add sessions for the specified batch
func (service *OBatchSessionService) AddSessionForBatch(batchSessions *[]bat.MappedSession, tenantID, batchID, credentialID uuid.UUID) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// check if session exist
	// calculate start and end date (automatic calculation) -> pending
	// set isCompleted field for new batch-session
	for index := range *batchSessions {
		// isCompleted := false
		(*batchSessions)[index].CreatedBy = credentialID
		err = service.doesSessionExists(tenantID, (*batchSessions)[index].CourseSessionID)
		if err != nil {
			return err
		}
		// (*batchSessions)[index].IsCompleted = &isCompleted
		(*batchSessions)[index].TenantID = tenantID
		(*batchSessions)[index].BatchID = batchID

		err = service.Repository.Add(uow, &(*batchSessions)[index])
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateSessionForBatch will update single session assigned to batch
func (service *OBatchSessionService) UpdateSessionForBatch(batchSession *bat.MappedSession) error {

	// check tenant exist
	err := service.doesTenantExists(batchSession.TenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(batchSession.TenantID, batchSession.UpdatedBy)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(batchSession.TenantID, batchSession.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesSessionExists(batchSession.TenantID, batchSession.CourseSessionID)
	if err != nil {
		return err
	}

	// check if batch-session exist
	err = service.doesBatchSessionExists(batchSession.TenantID, batchSession.ID)
	if err != nil {
		return err
	}

	// batchSession.Session = nil

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, &batchSession)
	if err != nil {
		uow.RollBack()
		return err
	}

	// update sub-sessions if main session is completed.
	if *batchSession.IsCompleted {
		tempCourseSessions := []crs.CourseSession{}
		err = service.Repository.GetAll(uow, &tempCourseSessions,
			repository.Join("INNER JOIN old_batch_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`session_id` AND "+
				"old_batch_sessions.`tenant_id` = course_sessions.`tenant_id`"), repository.Filter("course_sessions.`deleted_at` IS NULL"),
			repository.Filter("old_batch_sessions.`deleted_at` IS NULL AND old_batch_sessions.`tenant_id` = ?", batchSession.TenantID),
			repository.Filter("old_batch_sessions.`id` = ?", batchSession.ID))
		if err != nil {
			uow.RollBack()
			return err
		}

		for _, courseSession := range tempCourseSessions {
			err = service.Repository.UpdateWithMap(uow, bat.MappedSession{}, map[string]interface{}{
				"UpdatedBy":   batchSession.UpdatedBy,
				"IsCompleted": true,
			}, repository.Filter("`batch_id` = ? AND `course_session_id` = ?", batchSession.BatchID, courseSession.ID))
			if err != nil {
				uow.RollBack()
				return err
			}
		}
	}

	err = service.checkAllSessionsCompleted(batchSession.TenantID, batchSession.BatchID, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateSessionsForBatch will update the sessions for specified batch
func (service *OBatchSessionService) UpdateSessionsForBatch(batchSessions *[]bat.MappedSession, tenantID,
	batchID, credentialID uuid.UUID) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check credential exist
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	for _, batchSession := range *batchSessions {
		// check session exist
		err = service.doesSessionExists(tenantID, batchSession.CourseSessionID)
		if err != nil {
			return err
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempBatchSessions := []bat.MappedSession{}

	// get all sessions for specified batch
	err = service.Repository.GetAllForTenant(uow, tenantID, &tempBatchSessions,
		repository.Filter("`batch_id` = ?", batchID))
	if err != nil {
		uow.RollBack()
		return err
	}

	sessionMap := make(map[uuid.UUID]uint)

	// populate sessionMap with already existing entries
	for _, tempBatchSession := range tempBatchSessions {
		sessionMap[tempBatchSession.CourseSessionID] = sessionMap[tempBatchSession.CourseSessionID] + 1
	}

	for _, batchSession := range *batchSessions {
		sessionMap[batchSession.CourseSessionID] = sessionMap[batchSession.CourseSessionID] + 1
		// is value is greater than 1 then it indicates that it is an existing record
		if sessionMap[batchSession.CourseSessionID] > 1 {
			// if ID is nil then get the record
			if !util.IsUUIDValid(batchSession.ID) {
				temp := bat.MappedSession{}
				err = service.Repository.GetRecordForTenant(uow, tenantID, &temp,
					repository.Filter("`batch_id` = ? AND `start_date` IS NULL AND `course_session_id` = ?", batchID,
						batchSession.CourseSessionID), repository.Select("`id`"))
				if err != nil {
					return err
				}
				batchSession.ID = temp.ID
			}
			batchSession.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &batchSession)
			if err != nil {
				uow.RollBack()
				return err
			}
			sessionMap[batchSession.CourseSessionID] = 0
		}

		// adding new record to the batch
		if sessionMap[batchSession.CourseSessionID] == 1 {
			batchSession.TenantID = tenantID
			batchSession.CreatedBy = credentialID
			err = service.Repository.Add(uow, &batchSession)
			if err != nil {
				uow.RollBack()
				return err
			}
			sessionMap[batchSession.CourseSessionID] = 0
		}
	}

	for _, tempSession := range tempBatchSessions {
		if sessionMap[tempSession.CourseSessionID] == 1 {
			// delete session from corresponding tables.
			err = service.deleteSessionFromFeedback(uow, tenantID, credentialID, batchID, tempSession.ID)
			if err != nil {
				uow.RollBack()
				return err
			}

			err = service.Repository.UpdateWithMap(uow, bat.MappedSession{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`batch_id` = ? AND `course_session_id` = ?", batchID, tempSession.CourseSessionID))
			if err != nil {
				uow.RollBack()
				return err
			}
		}
		sessionMap[tempSession.CourseSessionID] = 0
	}

	// check if all sessions are completed
	err = service.checkAllSessionsCompleted(tenantID, batchID, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteSessionForBatch will delete the specified session for the batch
func (service *OBatchSessionService) DeleteSessionForBatch(batchSession *bat.MappedSession) error {

	// check if tenant exist
	err := service.doesTenantExists(batchSession.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(batchSession.TenantID, batchSession.DeletedBy)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExists(batchSession.TenantID, batchSession.BatchID)
	if err != nil {
		return err
	}

	// check if session exist
	err = service.doesSessionExists(batchSession.TenantID, batchSession.CourseSessionID)
	if err != nil {
		return err
	}

	// check if batch-session exist
	err = service.doesBatchSessionExists(batchSession.TenantID, batchSession.ID)
	if err != nil {
		return err
	}

	exist, err := repository.DoesRecordExistForTenant(service.DB, batchSession.TenantID, &bat.TopicAssignment{},
		repository.Filter("`batch_session_id` = ?", batchSession.ID))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Session cannot be deleted as assignments are assigned to this session")
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempBatchSession := bat.MappedSession{}

	err = service.Repository.GetRecordForTenant(uow, batchSession.TenantID, &tempBatchSession,
		repository.Filter("`id` = ?", batchSession.ID), repository.Select("`course_session_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	batchSession.CourseSessionID = tempBatchSession.CourseSessionID

	// repository.DoesRecordExist(service.DB, batch.MappedSession{},
	// 	repository.Join("INNER JOIN course_sessions ON batch_sessions.`course_session_id` = course_sessions.`session_id` AND"+
	// 		"batch_sessions.`tenant__id` = course_sessions.`tenant_id`"), repository.Filter("course_sessions.`deleted_at` IS NULL"),
	// 	repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_sessions.`tenant_id` = ?", batchSession.TenantID),
	// 	repository.Filter("batch_sessions.`id` = ?", batchSession.ID))

	// get all sub-sessions for session to be deleted
	exists, err := repository.DoesRecordExistForTenant(service.DB, batchSession.TenantID, crs.CourseSession{},
		repository.Filter("`session_id`=?", batchSession.CourseSessionID))
	if err != nil {
		return err
	}
	if exists {
		tempCourseSessions := &[]crs.CourseSession{}

		// err = service.Repository.GetAll(uow, batch.MappedSession{},
		// repository.Join("INNER JOIN course_sessions ON batch_sessions.`course_session_id` = course_sessions.`session_id` AND"+
		// 	"batch_sessions.`tenant_id` = course_sessions.`tenant_id`"), repository.Filter("course_sessions.`deleted_at` IS NULL"),
		// repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_sessions.`tenant_id` = ?", batchSession.TenantID),
		// repository.Filter("batch_sessions.`id` = ?", batchSession.ID))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }
		err = service.Repository.GetAllForTenant(uow, batchSession.TenantID, tempCourseSessions,
			repository.Filter("`session_id`=?", batchSession.CourseSessionID), repository.Select("id"))
		if err != nil {
			uow.RollBack()
			return err
		}

		// deleting sub-sessions
		for _, tempSession := range *tempCourseSessions {
			err = service.Repository.UpdateWithMap(uow, bat.MappedSession{}, map[string]interface{}{
				"DeletedBy": batchSession.DeletedBy,
				"DeletedAt": time.Now(),
			}, repository.Filter("`batch_id`=? AND `course_session_id`=?", batchSession.BatchID, tempSession.ID))
			if err != nil {
				uow.RollBack()
				return err
			}
		}
	}

	err = service.deleteSessionFromFeedback(uow, batchSession.TenantID, batchSession.DeletedBy,
		batchSession.BatchID, batchSession.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete specified session
	err = service.Repository.UpdateWithMap(uow, bat.MappedSession{}, map[string]interface{}{
		"DeletedBy": batchSession.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id`=? AND `course_session_id`=?", batchSession.BatchID, batchSession.CourseSessionID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete programming_assignments assigned to batch_session.
	err = service.Repository.UpdateWithMap(uow, bat.TopicAssignment{}, map[string]interface{}{
		"DeletedBy": batchSession.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_session_id` = ?", batchSession.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchModules will fetch module wise sessions and sub-sessions.
// func (service *OBatchSessionService) GetBatchModules(tenantID, batchID uuid.UUID,
// 	batchModules *[]bat.CourseModuleDTO, parser *web.Parser) error {

// 	now := time.Now()

// 	defer func() {
// 		fmt.Println("=================== duration ->", time.Since(now))
// 	}()

// 	// check if tenant exist
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, true)

// 	err = service.Repository.GetAllInOrder(uow, batchModules, "`order`",
// 		repository.Join("INNER JOIN `batches` ON batches.`course_id` = course_modules.`course_id` AND"+
// 			" batches.`tenant_id` = course_modules.`tenant_id`"),
// 		repository.Filter("batches.`id` = ? AND batches.`deleted_at` IS NULL AND course_modules.`tenant_id` = ?",
// 			batchID, tenantID)) // repository.PreloadAssociations([]string{"Course"}),
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	var totalCount int

// 	err = service.getBatchTalentCount(uow, tenantID, batchID, &totalCount)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	fmt.Println(" ======================================= getBatchSessions totalCount ->", totalCount)

// 	channel := make(chan error, 1)

// 	for index := range *batchModules {
// 		go service.getBatchSessions(uow, &(*batchModules)[index].BatchSessions, tenantID, batchID, (*batchModules)[index].ID,
// 			totalCount, parser, channel)
// 		err = <-channel
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	return nil
// }

// GetSessionsForBatch returns all the sessions assigned for the specified batch
func (service *OBatchSessionService) GetSessionsForBatch(batchSessions *[]bat.MappedSessionDTO, batchID,
	tenantID uuid.UUID, requestForm url.Values) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, batchSessions, "course_sessions.`order`",
		repository.Join("LEFT JOIN course_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`id` AND"+
			" course_sessions.`tenant_id` = old_batch_sessions.`tenant_id`"),
		repository.Filter("`batch_id`=?", batchID), repository.Filter("course_sessions.`session_id` IS NULL"),
		repository.Filter("old_batch_sessions.`tenant_id`=? AND course_sessions.`deleted_at` IS NULL", tenantID),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema:          "Session",
			Queryprocessors: []repository.QueryProcessor{repository.OrderBy("course_sessions.`order`")},
		}, repository.Preload{
			Schema: "Session.SubSessions",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN old_batch_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`id` AND" +
					" course_sessions.`tenant_id` = old_batch_sessions.`tenant_id`"),
				repository.Filter("course_sessions.`tenant_id`=? AND old_batch_sessions.`deleted_at` IS NULL", tenantID),
				repository.Filter("`batch_id`=?", batchID), repository.OrderBy("course_sessions.`order`"),
			},
		}), repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	var totalCount int

	err = service.getBatchTalentCount(uow, tenantID, batchID, &totalCount)
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *batchSessions {

		// checks if all feedback is given for all talents in a particular session.
		err = service.Repository.Scan(uow, &(*batchSessions)[index], repository.Table("faculty_talent_batch_session_feedback"),
			repository.Select("( CASE WHEN COUNT(DISTINCT `talent_id`) = ? THEN true ELSE false END ) AS is_feedback_given",
				totalCount), repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, (*batchSessions)[index].ID),
			repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL AND "+
				"faculty_talent_batch_session_feedback.`tenant_id` = ?", tenantID))
		if err != nil {
			uow.RollBack()
			return err
		}

		// checks if all feedback is given for all talents in a particular session.
		err = service.Repository.Scan(uow, &(*batchSessions)[index], repository.Table("batch_sessions_talents"),
			repository.Select("( CASE WHEN COUNT(DISTINCT `talent_id`) = ? THEN true ELSE false END ) AS is_attendance_given",
				totalCount), repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, (*batchSessions)[index].ID),
			repository.Filter("batch_sessions_talents.`deleted_at` IS NULL AND "+
				"batch_sessions_talents.`tenant_id` = ?", tenantID))
		if err != nil {
			return err
		}
	}
	uow.Commit()
	return nil
}

// GetSessionsAndAssignmentsForBatch returns all the sessions and assignments assigned for the specified batch
func (service *OBatchSessionService) GetSessionsAndAssignmentsForBatch(batchSessions *[]bat.MappedSessionDTO, batchID,
	tenantID uuid.UUID, requestForm url.Values) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, batchSessions, "course_sessions.`order`",
		repository.Join("LEFT JOIN course_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`id` AND"+
			" course_sessions.`tenant_id` = old_batch_sessions.`tenant_id`"),
		repository.Filter("`batch_id`=?", batchID), repository.Filter("course_sessions.`session_id` IS NULL"),
		repository.Filter("old_batch_sessions.`tenant_id`=? AND course_sessions.`deleted_at` IS NULL", tenantID),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema:          "Session",
			Queryprocessors: []repository.QueryProcessor{repository.OrderBy("course_sessions.`order`")},
		}, repository.Preload{
			Schema: "Session.SubSessions",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN old_batch_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`id` AND" +
					" course_sessions.`tenant_id` = old_batch_sessions.`tenant_id`"),
				repository.Filter("course_sessions.`tenant_id`=? AND old_batch_sessions.`deleted_at` IS NULL", tenantID),
				repository.Filter("`batch_id`=?", batchID), repository.OrderBy("course_sessions.`order`"),
			},
		}), repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	assignmentAssociation := []string{
		"ProgrammingAssignment", "ProgrammingAssignment.ProgrammingQuestion",
		"ProgrammingAssignment.ProgrammingAssignmentSubTask",
		"ProgrammingAssignment.ProgrammingQuestion.ProgrammingQuestionTypes",
	}

	for index := range *batchSessions {
		err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &(*batchSessions)[index].SessionAssignment, "`order`",
			repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, (*batchSessions)[index].ID),
			repository.PreloadAssociations(assignmentAssociation))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetBatchSessionList returns list of all the sessions assigned for the specified batch
func (service *OBatchSessionService) GetBatchSessionList(batchSessions *[]list.BatchSession, batchID, tenantID uuid.UUID,
	requestForm url.Values) error {

	// check tenant exist
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check batch exist
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, batchSessions, "course_sessions.`order`",
		repository.Filter("`batch_id` = ?", batchID),
		repository.Join("LEFT JOIN course_sessions ON old_batch_sessions.`course_session_id` = course_sessions.`id`"),
		repository.Filter("course_sessions.`session_id` IS NULL"),
		repository.Filter("old_batch_sessions.`tenant_id`=? AND old_batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("course_sessions.`tenant_id`=? AND course_sessions.`deleted_at` IS NULL", tenantID),
		repository.Select([]string{"`old_batch_sessions`.`id`", "`old_batch_sessions`.`course_session_id` AS session_id",
			"`course_sessions`.`name`, old_batch_sessions.`start_date`, old_batch_sessions.`is_completed`, old_batch_sessions.`order`"}),
		service.addSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// func (service *OBatchSessionService) getBatchSessions(uow *repository.UnitOfWork, batchSessions *[]*bat.BatchSessionDTO,
// 	tenantID, batchID, courseModuleID uuid.UUID, totalCount int, parser *web.Parser, channel chan error) {

// 	err := service.Repository.GetAll(uow, batchSessions,
// 		repository.Join("INNER JOIN `course_sessions` ON course_sessions.`id` = old_batch_sessions.`course_session_id` AND"+
// 			" old_batch_sessions.`tenant_id` = course_sessions.`tenant_id`"),
// 		repository.Join("INNER JOIN `course_modules` ON course_modules.`id` = course_sessions.`course_module_id` AND"+
// 			" course_sessions.`tenant_id` = course_modules.`tenant_id`"),
// 		repository.Filter("`course_sessions`.`course_module_id` = ?", courseModuleID),
// 		repository.Filter("`batch_id`=?", batchID), repository.Filter("course_sessions.`session_id` IS NULL"),
// 		repository.PreloadWithCustomCondition(repository.Preload{
// 			Schema:          "Session",
// 			Queryprocessors: []repository.QueryProcessor{repository.OrderBy("course_sessions.`order`")},
// 		}, repository.Preload{
// 			Schema: "Session.SubSessions",
// 			Queryprocessors: []repository.QueryProcessor{
// 				repository.Join("INNER JOIN `old_batch_sessions` ON old_batch_sessions.`course_session_id` = `course_sessions`.`id` AND" +
// 					" course_sessions.`tenant_id` = `old_batch_sessions`.`tenant_id`"),
// 				// repository.Filter("`course_sessions`.`course_module_id` = ?", (*batchModules)[index].ID),
// 				repository.Filter("course_sessions.`tenant_id`=? AND `old_batch_sessions`.`deleted_at` IS NULL", tenantID),
// 				repository.Filter("`batch_id`=?", batchID), repository.OrderBy("course_sessions.`order`"),
// 			},
// 		}), repository.PreloadAssociations(service.association))
// 	if err != nil {
// 		channel <- err
// 		return
// 	}

// 	for index := range *batchSessions {

// 		// checks if all feedback is given for all talents in a particular session.
// 		go service.getFeedbackGivenFlag(uow, tenantID, batchID, (*batchSessions)[index].ID,
// 			(*batchSessions)[index], totalCount, channel)

// 		// checks if all feedback is given for all talents in a particular session.
// 		go service.getAttendanceFlag(uow, tenantID, batchID, (*batchSessions)[index].ID,
// 			(*batchSessions)[index], totalCount, channel)

// 		err = <-channel
// 		if err != nil {
// 			channel <- err
// 			return
// 		}

// 		err = <-channel
// 		if err != nil {
// 			channel <- err
// 			return
// 		}
// 	}

// 	channel <- nil
// }

// func (service *OBatchSessionService) getFeedbackGivenFlag(uow *repository.UnitOfWork, tenantID, batchID,
// 	batchSessionID uuid.UUID, batchSession *bat.BatchSessionDTO, totalCount int, channel chan error) {

// 	err := service.Repository.Scan(uow, batchSession, repository.Table("faculty_talent_batch_session_feedback"),
// 		repository.Select("( CASE WHEN COUNT(DISTINCT `talent_id`) = ? THEN true ELSE false END ) AS is_feedback_given",
// 			totalCount), repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, batchSessionID),
// 		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL AND "+
// 			"faculty_talent_batch_session_feedback.`tenant_id` = ?", tenantID))
// 	if err != nil {
// 		channel <- err
// 		return
// 	}
// 	channel <- nil
// }

// func (service *OBatchSessionService) getAttendanceFlag(uow *repository.UnitOfWork, tenantID, batchID,
// 	batchSessionID uuid.UUID, batchSession *bat.BatchSessionDTO, totalCount int, channel chan error) {

// 	err := service.Repository.Scan(uow, batchSession, repository.Table("batch_sessions_talents"),
// 		repository.Select("( CASE WHEN COUNT(DISTINCT `talent_id`) = ? THEN true ELSE false END ) AS is_attendance_given",
// 			totalCount), repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, batchSessionID),
// 		repository.Filter("batch_sessions_talents.`deleted_at` IS NULL AND "+
// 			"batch_sessions_talents.`tenant_id` = ?", tenantID))
// 	if err != nil {
// 		channel <- err
// 		return
// 	}
// 	channel <- nil
// }

func (service *OBatchSessionService) getBatchTalentCount(uow *repository.UnitOfWork, tenantID,
	batchID uuid.UUID, totalCount *int) error {

	// join required as when talent is deleted no action is taken in batch_talents table.
	err := service.Repository.GetCount(uow, bat.MappedTalent{}, totalCount,
		repository.Join("INNER JOIN talents ON talents.`id` = batch_talents.`talent_id` AND"+
			" talents.`tenant_id` = batch_talents.`tenant_id`"), repository.Filter("batch_talents.`tenant_id` = ?", tenantID),
		repository.Filter("batch_talents.`batch_id` = ? AND batch_talents.`deleted_at` IS NULL", batchID),
		repository.Filter("talents.`deleted_at` IS NULL AND batch_talents.`is_active` = ?", true),
		repository.Filter("batch_talents.`suspension_date` IS NULL"))
	if err != nil {
		return err
	}

	return nil
}

// checkAllSessionsCompleted will check if all sessions are completed and if completed then sets the status as 'Finished' for the batch
func (service *OBatchSessionService) checkAllSessionsCompleted(tenantID, batchID uuid.UUID, uow *repository.UnitOfWork) error {

	// check if all sessions are completed
	exist, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, bat.MappedSession{},
		repository.Filter("`batch_id`=? AND `is_completed`=? AND `start_date` IS NOT NULL", batchID, 0))
	if err != nil {
		return err
	}
	// fmt.Println("****exist ->", exist)
	if !exist {
		// all batch-sessions are completed
		err = service.Repository.UpdateWithMap(uow, bat.Batch{}, map[interface{}]interface{}{
			// "IsCompleted": 1,
			"status": "Finished",
		}, repository.Filter("`id`=?", batchID))
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteSessionFromFeedback will delete session records from corresponding tables.
func (service *OBatchSessionService) deleteSessionFromFeedback(uow *repository.UnitOfWork,
	tenantID, credentialID, batchID, sessionID uuid.UUID) error {

	// delete faculty-talent-batch-session-feedback
	err := service.Repository.UpdateWithMap(uow, new(bat.FacultyTalentBatchSessionFeedback), map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, sessionID))
	if err != nil {
		return err
	}

	// delete talent-batch-session-feedback
	err = service.Repository.UpdateWithMap(uow, new(bat.TalentBatchSessionFeedback), map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, sessionID))
	if err != nil {
		return err
	}

	// delete aha-moment.
	err = service.Repository.UpdateWithMap(uow, new(bat.AhaMoment), map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, sessionID))
	if err != nil {
		return err
	}

	// delete aha-moment.
	err = service.Repository.UpdateWithMap(uow, new(bat.AhaMomentResponse), map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	}, repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, sessionID))
	if err != nil {
		return err
	}

	return nil
}

// addSearchQueries will add query processors based on queryparams.
func (service *OBatchSessionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if batchSessionID, ok := requestForm["batchSessionID"]; ok {
		util.AddToSlice("old_batch_sessions.`id`", "= ?", "AND", batchSessionID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExists validates tenantID
func (service *OBatchSessionService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesBatchExists validates batchID
func (service *OBatchSessionService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{}, repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *OBatchSessionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSessionExists checks if session id's are valid
func (service *OBatchSessionService) doesSessionExists(tenantID, sessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, crs.CourseSession{},
		repository.Filter("`id`=?", sessionID))
	if err := util.HandleError("Invalid session ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchSessionExists checks if batch's session id's are valid
func (service *OBatchSessionService) doesBatchSessionExists(tenantID, batchSessionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.MappedSession{},
		repository.Filter("`id`=?", batchSessionID))
	if err := util.HandleError("Invalid batch session ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// // returns error if there is no faculty/salesperson record in table for the given tenant.
// func (service *BatchSessionService) doesBatchExistForLogin(tenantID, batchID, loginID uuid.UUID, loginCheck string) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
// 		repository.Filter("`id`=?", batchID), repository.Filter(loginCheck, loginID))
// 	if err := util.HandleError("Batch not found", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// // returns error if there is no talent record in table for the given tenant.
// func (service *BatchSessionService) doesBatchExistForTalent(tenantID, batchID, talentID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.MappedTalent{},
// 		repository.Filter("`batch_id`=? AND `talent_id` = ?", batchID, talentID))
// 	if err := util.HandleError("Batch not found", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// doesFacultyCredentialExist checks if faculty credential id is valid
// func (service *BatchSessionService) doesFacultyCredentialExist(tenantID, facultyCredentialID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
// 		repository.Filter("`id`=?", facultyCredentialID))
// 	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }
