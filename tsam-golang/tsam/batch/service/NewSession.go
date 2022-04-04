package service

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// NewCreateBatchSessionPlan will create a plan for batch session.
func (service *SessionService) NewCreateBatchSessionPlan(tenantID, credentialID, batchID uuid.UUID,
	sessionTopics *[]batch.SessionTopic, parser *web.Parser) error {

	if util.IsEmpty(parser.Form.Get("facultyID")) {
		return errors.NewValidationError("Session plan can only be created by faculty.")
	}

	facultyID, err := util.ParseUUID(parser.Form.Get("facultyID"))
	if err != nil {
		return err
	}

	params := batch.SessionParameter{
		BatchSessions:       &[]batch.Session{},
		Modules:             &[]batch.ModuleDTO{},
		ModuleBatchSessions: &[]batch.Session{},
		ModuleSessionTopics: &[]batch.SessionTopic{},
		BatchID:             batchID,
		SessionTopics:       sessionTopics,
		Parser:              parser,
		FacultyID:           facultyID,
		TenantID:            tenantID,
		CredentialID:        credentialID,
		IsCreate:            true,
		IsSkip:              false,
	}

	// check foreign keys
	err = service.checkForeignkey(params)
	if err != nil {
		return err
	}

	// check if session plan exist for batch
	err = service.doesSessionPlanForFacultyExist(params.TenantID, params.BatchID, params.FacultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	// err = service.getBatchAndModules(uow, tenantID, batchID, params)
	err = service.getBatchModules(uow, params)
	if err != nil {
		uow.RollBack()
		return err
	}

	// check if unassigned module exist.
	// err = service.doesUnassignedModuleExist(params)
	// if err != nil {
	// 	return err
	// }

	err = service.createOrUpdatePlan(uow, params, true)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.updateBatch(uow, tenantID, params.BatchID, credentialID)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	fmt.Println(" ================= Session plan successfully created ================= ")
	return nil
}

type BatchDate struct {
	StartDate string
	EndDate   string
}

// updateBatch will update estimated_date for current batch.
func (service *SessionService) updateBatch(uow *repository.UnitOfWork, tenantID, batchID, credentialID uuid.UUID) error {

	// var tempBatchSession batch.Session
	var dates BatchDate

	err := service.repo.Scan(uow, &dates, repository.Table("batch_sessions"),
		repository.Select("min(`date`) AS start_date, max(`date`) AS end_date"),
		repository.Filter("`batch_id` = ? AND `tenant_id` = ?", batchID, tenantID))
	if err != nil {
		return err
	}

	err = service.repo.UpdateWithMap(uow, &batch.Batch{}, map[string]interface{}{
		// "StartDate":        dates.StartDate,
		"EstimatedEndDate": dates.EndDate,
		"UpdatedBy":        credentialID,
	}, repository.Filter("`id` = ?", batchID))
	if err != nil {
		return err
	}

	return nil
}

// SkipPendingSession will skip the pending sessions to next date.
func (service *SessionService) NewSkipPendingSession(tenantID, credentialID, batchID uuid.UUID,
	parser *web.Parser) error {

	// check if foreign key exist
	err := service.doesForeignKeyExist(tenantID, credentialID, batchID)
	if err != nil {
		return err
	}

	if util.IsEmpty(parser.Form.Get("facultyID")) {
		return errors.NewValidationError("Session plan can only be created/updated by faculty.")
	}

	facultyID, err := util.ParseUUID(parser.Form.Get("facultyID"))
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	params := batch.SessionParameter{
		BatchSessions:       &[]batch.Session{},
		Modules:             &[]batch.ModuleDTO{},
		ModuleBatchSessions: &[]batch.Session{},
		ModuleSessionTopics: &[]batch.SessionTopic{},
		BatchID:             batchID,
		SessionTopics:       &[]batch.SessionTopic{},
		Parser:              parser,
		FacultyID:           facultyID,
		TenantID:            tenantID,
		CredentialID:        credentialID,
		IsCreate:            false,
		IsSkip:              true,
	}

	// err = service.getBatchAndModules(uow, tenantID, batchID, params)
	err = service.getBatchModules(uow, params)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getBatchSessionTopics(uow, params)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.repo.GetAllInOrderForTenant(uow, tenantID, params.BatchSessions, "`batch_sessions`.`date`",
		repository.Filter("`batch_sessions`.`batch_id` = ? AND `batch_sessions`.`faculty_id` = ? AND `batch_sessions`.`is_session_taken` = ?",
			batchID, facultyID, false))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.createOrUpdatePlan(uow, params, false)
	if err != nil {
		uow.RollBack()
		return err
	}

	// uow.Commit()
	fmt.Println(" ================= Session plan successfully skipped ================= ")
	// return errors.NewValidationError(" === some error === ")
	return nil
}

// NewUpdateBatchSession will add/update topics to session.
func (service *SessionService) NewUpdateBatchSession(tenantID, credentialID, batchID uuid.UUID,
	sessionTopics *[]batch.SessionTopic, parser *web.Parser) error {

	if util.IsEmpty(parser.Form.Get("facultyID")) {
		return errors.NewValidationError("Session plan can only be created by faculty.")
	}

	facultyID, err := util.ParseUUID(parser.Form.Get("facultyID"))
	if err != nil {
		return err
	}

	removeCompletedSessionTopics(sessionTopics)

	params := batch.SessionParameter{
		BatchSessions:       &[]batch.Session{},
		Modules:             &[]batch.ModuleDTO{},
		ModuleBatchSessions: &[]batch.Session{},
		ModuleSessionTopics: &[]batch.SessionTopic{},
		BatchID:             batchID,
		SessionTopics:       sessionTopics,
		Parser:              parser,
		FacultyID:           facultyID,
		TenantID:            tenantID,
		CredentialID:        credentialID,
		IsCreate:            false,
		IsSkip:              false,
	}

	// check foreign keys
	err = service.checkForeignkey(params)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	err = service.deleteSessionTopic(uow, params)
	if err != nil {
		uow.RollBack()
		return err
	}

	// err = service.getBatchAndModules(uow, tenantID, batchID, params)
	err = service.getBatchModules(uow, params)
	if err != nil {
		uow.RollBack()
		return err
	}

	fmt.Println(" ==== len Modules ->", len(*params.Modules))

	err = service.repo.GetAllInOrderForTenant(uow, tenantID, params.BatchSessions, "`batch_sessions`.`date`",
		repository.Filter("`batch_sessions`.`batch_id` = ? AND `batch_sessions`.`is_session_taken` = ?", batchID, false))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.createOrUpdatePlan(uow, params, false)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	fmt.Println(" ================= Session plan successfully updated ================= ")
	// return errors.NewValidationError(" === some error === ")
	return nil
}

func (service *SessionService) deleteSessionTopic(uow *repository.UnitOfWork, params batch.SessionParameter) error {
	// fmt.Println(" ================= deleteSessionTopic ================= ")

	var tempSessionTopics []batch.SessionTopic

	err := service.repo.GetAll(uow, &tempSessionTopics,
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND "+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`is_completed` = ?"+
			" AND batch_sessions.`faculty_id` = ? AND batch_session_topics.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL",
			params.BatchID, 0, params.FacultyID, params.TenantID), repository.OrderBy("`completed_date`, `order`"),
		repository.GroupBy("batch_session_topics.`id`"))
	if err != nil {
		return err
	}

	if len(tempSessionTopics) == 0 {
		return nil
	}

	var topicMap = make(map[uuid.UUID]uint)

	for _, tempSessionTopic := range tempSessionTopics {
		if util.IsUUIDValid(tempSessionTopic.ID) {
			topicMap[tempSessionTopic.ID]++
		}
	}

	for _, sessionTopic := range *params.SessionTopics {
		if !util.IsUUIDValid(sessionTopic.ID) {
			continue
		}
		topicMap[sessionTopic.ID]++
	}

	// fmt.Println(" ================== topic map -> ", topicMap)

	// delete those entries where count od session topic ID is 1.
	for key, value := range topicMap {
		if value == 1 {
			err = service.repo.UpdateWithMap(uow, batch.SessionTopic{}, map[string]interface{}{
				"DeletedAt": time.Now(),
				"DeletedBy": params.CredentialID,
			}, repository.Filter("`id` = ?", key))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// MarkSubTopicAsComplete will mark specified session as completed.
func (service *SessionService) MarkSubTopicAsComplete(batchSessionTopic *batch.SessionTopic) error {

	// check if foreign keys exist.
	err := service.doesForeignKeyExist(batchSessionTopic.TenantID, batchSessionTopic.UpdatedBy, batchSessionTopic.BatchID)
	if err != nil {
		return err
	}

	// check if sub-topic exist.
	err = service.doesBatchSessionTopicExists(batchSessionTopic.TenantID,
		batchSessionTopic.BatchID, batchSessionTopic.TopicID, batchSessionTopic.SubTopicID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)

	// tempBatchSessionTopic := batch.BatchSessionTopic{}
	credentialID := batchSessionTopic.UpdatedBy
	date := batchSessionTopic.CompletedDate

	err = service.repo.GetRecordForTenant(uow, batchSessionTopic.TenantID, batchSessionTopic,
		repository.Filter("`batch_id` = ? AND `topic_id` = ? AND `sub_topic_id` = ?",
			batchSessionTopic.BatchID, batchSessionTopic.TopicID, batchSessionTopic.SubTopicID))
	if err != nil {
		uow.RollBack()
		return err
	}

	*batchSessionTopic.IsCompleted = true
	batchSessionTopic.UpdatedBy = credentialID
	batchSessionTopic.CompletedDate = date

	err = service.repo.UpdateWithMap(uow, batch.SessionTopic{}, map[string]interface{}{
		"IsCompleted":   true,
		"CompletedDate": batchSessionTopic.CompletedDate,
		"UpdatedBy":     batchSessionTopic.UpdatedBy,
	}, repository.Filter("`id` = ?", batchSessionTopic.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.markBatchSessionCompletion(uow, batchSessionTopic)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.markBatchStatus(uow, batchSessionTopic)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// markBatchSessionCompletion will mark batch_session's is_completed flag.
func (service *SessionService) markBatchSessionCompletion(uow *repository.UnitOfWork,
	batchSessionTopic *batch.SessionTopic) error {

	exist, err := repository.DoesRecordExistForTenant(uow.DB, batchSessionTopic.TenantID, batch.SessionTopic{},
		repository.Filter("`batch_session_id` = ? AND `batch_id` = ? AND `completed_date` = ? AND `is_completed` = ?",
			batchSessionTopic.SessionID, batchSessionTopic.BatchID, batchSessionTopic.CompletedDate, false))
	if err != nil {
		return err
	}

	if !exist {
		err = service.repo.UpdateWithMap(uow, batch.Session{}, map[string]interface{}{
			"IsCompleted":    true,
			"IsSessionTaken": true,
			"UpdatedBy":      batchSessionTopic.UpdatedBy,
		}, repository.Filter("`id` = ?", batchSessionTopic.SessionID))
		if err != nil {
			return err
		}
		return nil
	}

	err = service.repo.UpdateWithMap(uow, batch.Session{}, map[string]interface{}{
		"IsSessionTaken": true,
		"UpdatedBy":      batchSessionTopic.UpdatedBy,
	}, repository.Filter("`date` = ?", batchSessionTopic.CompletedDate))
	if err != nil {
		return err
	}

	return nil
}

func (service *SessionService) markBatchStatus(uow *repository.UnitOfWork, batchSessionTopic *batch.SessionTopic) error {

	exist, err := repository.DoesRecordExistForTenant(service.db, batchSessionTopic.TenantID, batch.SessionTopic{},
		repository.Filter("`batch_id` = ? AND `is_completed` = ?", batchSessionTopic.BatchID, 0))
	if err != nil {
		return err
	}
	if !exist {
		err = service.repo.UpdateWithMap(uow, batch.Batch{}, map[string]interface{}{
			"Status":    "Finished",
			"UpdatedBy": batchSessionTopic.UpdatedBy,
		}, repository.Filter("`id` = ?", batchSessionTopic.BatchID))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteSessionPlan will delete specified session and its topics.
// code should be changed/improved #shailesh
func (service *SessionService) DeleteSessionPlan(session *batch.Session) error {

	// check if tenant exist
	err := service.doesTenantExists(session.TenantID)
	if err != nil {
		return err
	}

	// check if batch exist
	err = service.doesBatchExists(session.TenantID, session.BatchID)
	if err != nil {
		return err
	}

	// check if faculty exist
	err = service.doesFacultyExists(session.TenantID, session.FacultyID)
	if err != nil {
		return err
	}

	exist, err := repository.DoesRecordExistForTenant(service.db, session.TenantID, batch.Session{},
		repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`faculty_id` = ?"+
			" AND batch_sessions.`is_session_taken` = ?", session.BatchID, session.FacultyID, 1))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("session plan cannot be deleted as its sessions are already completed")
	}

	uow := repository.NewUnitOfWork(service.db, false)

	var tempSessionTopics []batch.SessionTopic

	err = service.repo.GetAll(uow, &tempSessionTopics, repository.Select("batch_session_topics.`id`"),
		repository.Join("INNER JOIN `batch_sessions` ON batch_sessions.`id` = batch_session_topics.`batch_session_id`"+
			" AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`faculty_id` = ?", session.BatchID, session.FacultyID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_session_topics.`tenant_id` = ?", session.TenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	var tempTopicQuestions []batch.TopicAssignment

	err = service.repo.GetAll(uow, &tempTopicQuestions, repository.Select("batch_topic_assignments.`id`"),
		repository.Join("INNER JOIN `batch_session_topics` ON batch_session_topics.`topic_id` = batch_topic_assignments.`topic_id`"+
			" AND batch_session_topics.`tenant_id` = batch_topic_assignments.`tenant_id`"),
		repository.Join("INNER JOIN `batch_sessions` ON batch_sessions.`id` = batch_session_topics.`batch_session_id`"+
			" AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_topic_assignments.`batch_id` = ? AND batch_sessions.`faculty_id` = ?", session.BatchID, session.FacultyID),
		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_topic_assignments.`tenant_id` = ?"+
			" AND batch_sessions.`deleted_at` IS NULL", session.TenantID), repository.GroupBy("batch_topic_assignments.`id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	var sessionTopicID []uuid.UUID

	for _, sessionTopic := range tempSessionTopics {
		sessionTopicID = append(sessionTopicID, sessionTopic.ID)
	}

	// delete session topic
	err = service.repo.UpdateWithMap(uow, batch.SessionTopic{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": session.DeletedBy,
	}, repository.Filter("`id` IN (?)", sessionTopicID))
	if err != nil {
		uow.RollBack()
		return err
	}

	var topicQuestionID []uuid.UUID

	for _, topicQuestion := range tempTopicQuestions {
		topicQuestionID = append(topicQuestionID, topicQuestion.ID)
	}

	// delete topic assignments
	err = service.repo.UpdateWithMap(uow, batch.TopicAssignment{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": session.DeletedBy,
	}, repository.Filter("`id` IN (?)", topicQuestionID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// delete session
	err = service.repo.UpdateWithMap(uow, batch.Session{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": session.DeletedBy,
	}, repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`faculty_id` = ?", session.BatchID, session.FacultyID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// NewGetBatchSessionPlan will get batch sessions for the specified batch.
func (service *SessionService) NewGetBatchSessionPlan(batchSession *batch.SessionDTO,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	if util.IsEmpty(parser.Form.Get("sessionDate")) {
		return errors.NewValidationError("session date must be specified")
	}

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	// service.addSessionSearchQueries(parser.Form),
	err = service.repo.GetRecord(uow, batchSession,
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`batch_session_id` = batch_sessions.`id` AND "+
			" batch_session_topics.`tenant_id` = batch_sessions.`tenant_id`"), repository.Filter("batch_sessions.`batch_id` = ?", batchID),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		service.addSearchQueries(parser.Form), service.addSessionTopicSearchQueries(parser.Form, false),
		repository.PreloadAssociations([]string{"BatchSessionPrerequisiteDTO"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	channel := make(chan error, 2)
	var wg sync.WaitGroup

	currentParam := batch.GetModuleParams{
		TenantID:  tenantID,
		BatchID:   batchID,
		Parser:    parser,
		IsPending: false,
		Channel:   channel,
		Modules:   &batchSession.Module,
	}

	wg.Add(1)
	go service.newGetCourseModules(uow, currentParam, func() { wg.Done() })

	if !util.IsEmpty(parser.Form.Get("pendingCompletedDate")) {

		pendingParam := batch.GetModuleParams{
			TenantID:  tenantID,
			BatchID:   batchID,
			Parser:    parser,
			IsPending: true,
			Channel:   channel,
			Modules:   &batchSession.PendingModule,
		}

		wg.Add(1)
		go service.newGetCourseModules(uow, pendingParam, func() { wg.Done() })
	}

	go func() {
		defer close(channel)
		wg.Wait()
	}()

	for err := range channel {
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	err = service.getAttendanceCompleted(uow, tenantID, batchID, batchSession)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getFeedbackGivenFlag(uow, tenantID, batchID, batchSession)
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// NewGetAllBatchSessions will return all batch session and its related topics and subtopics.
func (service *SessionService) NewGetAllBatchSessions(batchSessions *[]batch.SessionDTO,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	// now := time.Now()

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	// service.addSessionSearchQueries(parser.Form),
	err = service.repo.GetAllInOrder(uow, batchSessions, "batch_sessions.`date`",
		repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id` "+
			"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		service.addSearchQueries(parser.Form), repository.Filter("batch_sessions.`batch_id` = ?", batchID),
		repository.GroupBy("batch_sessions.`date`"),
		repository.PreloadAssociations([]string{"BatchSessionPrerequisiteDTO"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	channel := make(chan error, len(*batchSessions))
	var wg sync.WaitGroup

	for index := range *batchSessions {
		wg.Add(1)

		param := batch.GetModuleParams{
			TenantID:    tenantID,
			BatchID:     batchID,
			Parser:      parser,
			IsPending:   false,
			Channel:     channel,
			Modules:     &(*batchSessions)[index].Module,
			SessionDate: (*batchSessions)[index].Date,
		}

		go service.newGetCourseModules(uow, param, func() { wg.Done() })
	}

	go func() {
		defer close(channel)
		wg.Wait()
	}()

	for err := range channel {
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	// defer func() {
	// 	fmt.Println(" ===================== duration ->", time.Since(now))
	// 	fmt.Println("Go routine Number is --------------------------", runtime.NumGoroutine())
	// }()

	uow.Commit()
	return nil
}

// GetBatchSessions will get sessions for specified batch.
func (service *SessionService) GetBatchSessions(batchSessions *[]batch.Session, tenantID,
	batchID uuid.UUID, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	// repository.Join("INNER JOIN batch_modules ON batch_sessions.`batch_id` = batch_modules.`batch_id` "+
	// "AND batch_sessions.`tenant_id` = batch_modules.`tenant_id`"), AND"+
	// " batch_modules.`deleted_at` IS NULL

	err = service.repo.GetAllInOrder(uow, batchSessions, "batch_sessions.`date`",
		repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id` "+
			"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"), service.addSearchQueries(parser.Form),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`batch_id` = ?", batchID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
		repository.GroupBy("batch_sessions.`date`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// NewGetBatchSessionsCounts will return all counts for batch session plan for specified batch.
func (service *SessionService) NewGetBatchSessionsCounts(batchSessions *batch.SessionCounts,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	// Get count of modules.
	err = service.getSessionModuleCount(uow, tenantID, batchID, &batchSessions.ModuleCount, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get count of topics.
	err = service.getSessionTopicCount(uow, tenantID, batchID, &batchSessions.TopicCount, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get count of assignments.
	err = service.getSessionAssignmentCount(uow, tenantID, batchID, &batchSessions.AssignmentCount, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get count of batch sessions.
	err = service.getSessionCount(uow, tenantID, batchID, &batchSessions.SessionCount, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get total hours of batch.
	err = service.getTotalBatchHours(uow, batchSessions, tenantID, batchID, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getBatchProjectCount(uow, tenantID, batchID, &batchSessions.ProjectCount, parser)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetBatchSessionWithTopicNameList will return list of topic names and sub-topic names for specified batch.
func (service *SessionService) GetBatchSessionWithTopicNameList(tenantID, batchID uuid.UUID,
	batchSessions *[]batch.SessionWithTopicNameDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAllInOrder(uow, batchSessions, "batch_sessions.`date`",
		repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_sessions.`batch_id`=?", batchID),
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"Faculty"}),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchSessionTopics",
			Queryprocessors: []repository.QueryProcessor{
				repository.OrderBy("batch_session_topics.`order`"),
				repository.PreloadAssociations([]string{"Topic", "SubTopic"}),
			}}),
		repository.GroupBy("batch_sessions.`date`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Give batch module timings to batch session.
	for index := range *batchSessions {
		batchSessionDate, err := time.Parse("2006-01-02", (*batchSessions)[index].Date)
		if err != nil {
			uow.RollBack()
			return err
		}

		batchSessionWeekDay := batchSessionDate.Weekday()
		tempBatchModuleTimings := batch.ModuleTiming{}
		// err = service.repo.GetAll(uow, &tempBatchModuleTimings,
		// 	repository.Table("batch_sessions"),
		// 	repository.Select("batch_module_timings.*"),
		// 	repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id`"),
		// 	repository.Join("JOIN batch_module_timings on batch_module_timings.`module_id` = batch_session_topics.`module_id`"),
		// 	repository.Join("JOIN days on days.`id` = batch_module_timings.`day_id`"),
		// 	repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
		// 	repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		// 	repository.Filter("batch_module_timings.`tenant_id` = ? AND batch_module_timings.`deleted_at` IS NULL", tenantID),
		// 	repository.Filter("days.`tenant_id` = ? AND days.`deleted_at` IS NULL", tenantID),
		// 	repository.Filter("batch_sessions.`date`=?", (*batchSessions)[i].Date),
		// 	repository.Filter("days.`day`=?", batchSessionWeekDay),
		// 	repository.Filter("batch_sessions.`batch_id`=?", batchID),
		// 	repository.GroupBy("batch_module_timings.`id`"))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }

		err = service.repo.GetRecord(uow, &tempBatchModuleTimings,
			repository.Join("INNER JOIN batch_session_topics on batch_module_timings.module_id = batch_session_topics.module_id"),
			repository.Join("INNER JOIN batch_sessions ON batch_sessions.id = batch_session_topics.batch_session_id"),
			repository.Join("INNER JOIN days on days.id = batch_module_timings.day_id"),
			repository.Filter("batch_session_topics.tenant_id = ? AND batch_session_topics.deleted_at IS NULL", tenantID),
			repository.Filter("batch_sessions.tenant_id = ? AND batch_sessions.deleted_at IS NULL", tenantID),
			repository.Filter("batch_module_timings.tenant_id = ? AND batch_module_timings.deleted_at IS NULL", tenantID),
			repository.Filter("days.tenant_id = ? AND days.deleted_at IS NULL", tenantID),
			repository.Filter("batch_sessions.date=?", (*batchSessions)[index].Date),
			repository.Filter("days.day=?", batchSessionWeekDay.String()),
			repository.Filter("batch_sessions.batch_id=?", batchID),
			repository.GroupBy("batch_module_timings.id"))
		if err != nil {
			uow.RollBack()
			return err
		}

		(*batchSessions)[index].ModuleTiming = tempBatchModuleTimings
	}

	uow.Commit()
	return nil
}

// GetSessionTopics will return topics, subtopics for specified module.
func (service *SessionService) GetSessionModuleTopics(tenantID, batchID uuid.UUID,
	moduleTopics *[]course.ModuleTopicDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// check if batch exist.
	err = service.doesBatchExists(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAllInOrder(uow, moduleTopics, "`batch_session_topics`.`order`",
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`topic_id` = module_topics.`id` AND "+
			"batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND "+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND module_topics.`tenant_id` = ?", tenantID),
		repository.Filter("batch_session_topics.`batch_id` = ?", batchID),
		repository.GroupBy("batch_session_topics.`topic_id`"), service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{
			"Module", "TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
		}), repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchTopicAssignment",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
				repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
			},
		}))
	if err != nil {
		uow.RollBack()
		return err
	}

	for j := range *moduleTopics {
		param := batch.GetModuleParams{
			TenantID:  tenantID,
			BatchID:   batchID,
			Parser:    parser,
			IsPending: false,
			Channel:   nil,
			Modules:   nil,
		}
		err = service.newGetModuleSubTopics(uow, &(*moduleTopics)[j].SubTopics, (*moduleTopics)[j].Module.ID, (*moduleTopics)[j].ID, param)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (service *SessionService) createOrUpdatePlan(uow *repository.UnitOfWork,
	params batch.SessionParameter, isCreate bool) error {

	var builder strings.Builder
	var strLen int
	var err error

	for index := range *params.Modules {
		fmt.Println(" ========================= before if ========================= ", index, !params.IsCreate)
		if index > 0 { // !params.IsCreate
			fmt.Println(" ========================= inside if ========================= ")
			if (*params.Modules)[index].StartDate == nil {
				fmt.Println(" ========================= start date not found ========================= ")
				(*params.Modules)[index].StartDate, err = getNextStartDate((*(*params.Modules)[index-1].EstimatedEndDate)[:10])
				if err != nil {
					return err
				}
			} else {
				fmt.Println(" ========================= start date found ========================= ")
				isParallel, err := isSessionParallel(index, params)
				if err != nil {
					return err
				}

				if !isParallel {
					(*params.Modules)[index].StartDate, err = getNextStartDate((*(*params.Modules)[index-1].EstimatedEndDate)[:10])
					if err != nil {
						return err
					}
				}
			}
		}
		fmt.Println(" ========================= end if ========================= ")

		params.FacultyID = (*params.Modules)[index].FacultyID
		params.ModuleBatchSessions = &[]batch.Session{}
		params.ModuleSessionTopics = &[]batch.SessionTopic{}

		err := service.getModuleSessionTopic(uow, params, (*params.Modules)[index].ModuleID, isCreate)
		if err != nil {
			return err
		}
		fmt.Println(" ================== module start date ->", (*params.Modules)[index].StartDate)

		err = service.createBatchModulePlan(params, &(*params.Modules)[index])
		if err != nil {
			return err
		}

		if isCreate {

			if len(*params.ModuleBatchSessions) == 0 {
				strLen, err = builder.WriteString("Session plan for " + (*params.Modules)[index].Module.ModuleName + " module could not be created \n")
				if err != nil {
					return err
				}
				continue
			}

			err = service.addSessionTopicPlan(uow, params, &(*params.Modules)[index], params.TenantID, params.CredentialID)
			if err != nil {
				return err
			}

			err = service.updateBatchModule(uow, &(*params.Modules)[index], params.TenantID, params.CredentialID)
			if err != nil {
				return err
			}

			continue
		}

		if len(*params.ModuleSessionTopics) == 0 {
			strLen, err = builder.WriteString("Session plan for " + (*params.Modules)[index].Module.ModuleName + " module could not be created \n")
			if err != nil {
				return err
			}
			continue
		}

		err = service.updateSessionTopicPlan(uow, params, &(*params.Modules)[index], params.TenantID, params.CredentialID)
		if err != nil {
			return err
		}

		err = service.updateBatchModule(uow, &(*params.Modules)[index], params.TenantID, params.CredentialID)
		if err != nil {
			return err
		}

	}

	if strLen > 0 {
		return errors.NewValidationError(builder.String())
	}

	// return errors.NewValidationError(" ==== some error === ")
	return nil
}

func (service *SessionService) updateBatchModule(uow *repository.UnitOfWork,
	module *batch.ModuleDTO, tenantID, credentialID uuid.UUID) error {

	err := service.repo.UpdateWithMap(uow, batch.Module{}, map[string]interface{}{
		"StartDate":        module.StartDate,
		"EstimatedEndDate": module.EstimatedEndDate,
		"UpdatedBy":        credentialID,
	}, repository.Filter("`id` = ?", module.ID))
	if err != nil {
		return err
	}
	return nil
}

// addSessionTopicPlan will add all the session plans to the table.
func (service *SessionService) addSessionTopicPlan(uow *repository.UnitOfWork,
	params batch.SessionParameter, module *batch.ModuleDTO, tenantID, credentialID uuid.UUID) error {

	isCompleted := false
	// fmt.Println(" ===================================== len(*batchSessions) ->", len(*batchSessions))

	for _, batchSession := range *params.ModuleBatchSessions {

		batchSession.CreatedBy = credentialID
		batchSession.TenantID = tenantID
		batchSession.IsSessionTaken, batchSession.IsCompleted = &isCompleted, &isCompleted

		for index := range batchSession.BatchSessionTopic {
			batchSession.BatchSessionTopic[index].InitialDate = batchSession.BatchSessionTopic[index].CompletedDate
			batchSession.BatchSessionTopic[index].IsCompleted = &isCompleted
			batchSession.BatchSessionTopic[index].CreatedBy = credentialID
			batchSession.BatchSessionTopic[index].TenantID = tenantID
		}

		module.EstimatedEndDate = &batchSession.Date

		err := service.repo.Add(uow, &batchSession)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateSessionTopicPlan will update incomplete sessions date.
func (service *SessionService) updateSessionTopicPlan(uow *repository.UnitOfWork,
	params batch.SessionParameter, module *batch.ModuleDTO, tenantID, credentialID uuid.UUID) error {

	isBatchSessionFound := false
	var previousDate string

	for _, sessionTopic := range *params.ModuleSessionTopics {
		// fmt.Println(" =================== sessionTopic date ->", *sessionTopic.InitialDate)
		sessionTopic.UpdatedBy = credentialID
		sessionTopic.TenantID = tenantID
		module.EstimatedEndDate = sessionTopic.CompletedDate
		*sessionTopic.InitialDate = (*sessionTopic.InitialDate)[:10]

		// fmt.Println(" ============ (*sessionTopic.InitialDate) -> ", *sessionTopic.InitialDate)
		// fmt.Println(" ============ (*sessionTopic.CompletedDate) -> ", *sessionTopic.CompletedDate)

		isBatchSessionFound = false

		for _, batchSession := range *params.BatchSessions {
			// fmt.Println(" =================== batchSession.Date ->", batchSession.Date)

			if batchSession.Date == *sessionTopic.CompletedDate {
				sessionTopic.SessionID = batchSession.ID
				isBatchSessionFound = true

				if previousDate != (*sessionTopic.CompletedDate) {
					previousDate = *sessionTopic.CompletedDate

					err := service.repo.UpdateWithMap(uow, &batch.Session{}, map[string]interface{}{
						"UpdatedBy": credentialID,
						"UpdatedAt": time.Now(),
					}, repository.Filter("`id` = ?", batchSession.ID))
					if err != nil {
						return err
					}
				}
				break
			}
		}

		// fmt.Println(" ============ isBatchSessionFound -> ", isBatchSessionFound)

		if !isBatchSessionFound {
			err := service.addBatchSession(uow, params, &sessionTopic, tenantID, credentialID)
			if err != nil {
				return err
			}

			// taking last batch session as it is appended at end in the function
			lastIndex := len(*params.BatchSessions) - 1
			sessionTopic.SessionID = (*params.BatchSessions)[lastIndex].ID

			// fmt.Println(" ============== session ID ->", sessionTopic.SessionID)
		}

		// fmt.Println(" ============== util.IsUUIDValid(sessionTopic.ID) -> ", util.IsUUIDValid(sessionTopic.ID))

		if util.IsUUIDValid(sessionTopic.ID) {
			// update session topics
			err := service.repo.Update(uow, &sessionTopic)
			if err != nil {
				return err
			}
			continue
		}

		isCompleted := false
		sessionTopic.IsCompleted = &isCompleted

		err := service.repo.Add(uow, &sessionTopic)
		if err != nil {
			return err
		}
	}

	return nil
}

// addBatchSession will add new batch-session.
func (service *SessionService) addBatchSession(uow *repository.UnitOfWork, params batch.SessionParameter,
	sessionTopic *batch.SessionTopic, tenantID, credentialID uuid.UUID) error {

	tempBatchSession := batch.Session{
		BatchID:   params.BatchID,
		Date:      *sessionTopic.CompletedDate,
		FacultyID: params.FacultyID,
	}

	booleanValue := false
	tempBatchSession.IsCompleted, tempBatchSession.IsSessionTaken = &booleanValue, &booleanValue
	tempBatchSession.CreatedBy = credentialID
	tempBatchSession.TenantID = tenantID

	err := service.repo.Add(uow, &tempBatchSession)
	if err != nil {
		return err
	}

	(*params.BatchSessions) = append((*params.BatchSessions), tempBatchSession)
	return nil
}

// createBatchModulePlan will create session plan for single module
func (service *SessionService) createBatchModulePlan(params batch.SessionParameter, module *batch.ModuleDTO) error {

	dayTimingMap := make(map[uint]time.Duration)

	err := getModuleTimeDifference(module, dayTimingMap)
	if err != nil {
		return err
	}

	// fmt.Println(" ====================== *module.StartDate ->", *module.StartDate)

	moduleTotalHours := getTotalHours(params.ModuleSessionTopics)
	datesArray, err := getModuleDatesArray(module, *module.StartDate, moduleTotalHours, dayTimingMap)
	if err != nil {
		return err
	}

	var today time.Time

	if params.IsSkip {
		today, err = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		if err != nil {
			return err
		}

		for i := 0; i < len(datesArray); i++ {
			date, err := time.Parse("2006-01-02", datesArray[i])
			if err != nil {
				return err
			}

			// splice date if it is before today's date.
			if date.Before(today) || today.Equal(date) {
				datesArray = append(datesArray[:i], datesArray[i+1:]...)
				i = i - 1
			}
		}
	}

	// fmt.Println(" ========================= datesArray ->", datesArray)

	module.StartDate = &datesArray[0]
	// fmt.Println(" ========================= module.StartDate ->", *module.StartDate)
	// fmt.Println(" ========================= module.Module.ID ->", module.Module.ID)

	params.CurrentModuleID = module.Module.ID
	err = generateModuleSessionPlan(datesArray, dayTimingMap, params)
	if err != nil {
		return err
	}

	return nil
}

func generateModuleSessionPlan(datesArray []string, dayTimingMap map[uint]time.Duration,
	params batch.SessionParameter) error {

	// fmt.Println(" ============================= dayTimingMap ->", dayTimingMap)

	var tempSessionTopics []batch.SessionTopic
	var tempBatchSession batch.Session

	for _, date := range datesArray {

		if params.IsCreate {
			tempBatchSession = batch.Session{
				BatchID: params.BatchID,
				Date:    date,
			}
		}

		currentDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return err
		}

		currentDay := int(currentDate.Weekday())
		allocationDurationLeft := dayTimingMap[uint(currentDay)].Minutes()

		if len(*params.SessionTopics) == 0 && params.IsCreate {
			break
		}

		if len(*params.ModuleSessionTopics) == 0 && !params.IsCreate {
			break
		}

		allocateParams := batch.AllocateParameter{
			ModuleSessionTopics:    &tempSessionTopics,
			BatchSession:           &tempBatchSession,
			AllocationDurationLeft: allocationDurationLeft,
			CurrentDate:            currentDate,
		}

		allocateModuleSessionTopic(params, allocateParams)

		if params.IsCreate {
			tempBatchSession.FacultyID = params.FacultyID
			if len(tempBatchSession.BatchSessionTopic) > 0 {
				*params.ModuleBatchSessions = append(*params.ModuleBatchSessions, tempBatchSession)
			}
		}
	}

	if !params.IsCreate {
		*params.ModuleSessionTopics = append(*params.ModuleSessionTopics, tempSessionTopics...)
	}

	// fmt.Println(" ============= len(*params.ModuleSessionTopics)", len(*params.ModuleSessionTopics))

	return nil
}

func allocateModuleSessionTopic(params batch.SessionParameter, allocateParams batch.AllocateParameter) {
	for i := 0; i < len(*params.ModuleSessionTopics); i++ {

		if params.CurrentModuleID == (*params.ModuleSessionTopics)[i].ModuleID {
			allocateParams.AllocationDurationLeft -= float64((*params.ModuleSessionTopics)[i].TotalTime)

			if allocateParams.AllocationDurationLeft >= 0.0 {
				completedDate := allocateParams.CurrentDate.Format("2006-01-02")
				(*params.ModuleSessionTopics)[i].CompletedDate = &completedDate

				if !params.IsCreate {
					// fmt.Println(" =========== id ->", (*params.ModuleSessionTopics)[i].ID)
					if !util.IsUUIDValid((*params.ModuleSessionTopics)[i].ID) {
						(*params.ModuleSessionTopics)[i].InitialDate = &completedDate
					}
					// fmt.Println(" =========== initial date ->", *(*params.ModuleSessionTopics)[i].InitialDate)
					(*allocateParams.ModuleSessionTopics) = append((*allocateParams.ModuleSessionTopics),
						(*params.ModuleSessionTopics)[i])
				}

				if params.IsCreate {
					allocateParams.BatchSession.BatchSessionTopic = append(allocateParams.BatchSession.BatchSessionTopic,
						(*params.ModuleSessionTopics)[i])
				}

				(*params.ModuleSessionTopics) = append((*params.ModuleSessionTopics)[:i], (*params.ModuleSessionTopics)[i+1:]...)
				i = i - 1
				continue
			}
			return
		}
	}
}

// getModuleSessionTopic will return all topics for specified module
func (service *SessionService) getModuleSessionTopic(uow *repository.UnitOfWork, params batch.SessionParameter,
	moduleID uuid.UUID, isCreate bool) error {
	for index := range *params.SessionTopics {
		if (*params.SessionTopics)[index].ModuleID == moduleID {
			*params.ModuleSessionTopics = append(*params.ModuleSessionTopics, (*params.SessionTopics)[index])
		}
	}

	// if topics are not found for specified module then fetch those topics
	if !isCreate && len(*params.ModuleSessionTopics) == 0 {
		fmt.Println(" ================ topics not found ================ ")
		err := service.getSessionTopics(uow, params, params.ModuleSessionTopics, moduleID)
		if err != nil {
			return err
		}
	}

	// return errors.NewValidationError(" === error!!!!!!! === ")
	return nil
}

// getSessionTopics will get all session topics for specified module.
func (service *SessionService) getSessionTopics(uow *repository.UnitOfWork, params batch.SessionParameter,
	moduleSessionTopics *[]batch.SessionTopic, moduleID uuid.UUID) error {

	//  AND `is_completed` = ?
	err := service.repo.GetAllInOrderForTenant(uow, params.TenantID, moduleSessionTopics, `order`,
		repository.Filter("`batch_id` = ? AND `module_id` = ?", params.BatchID, moduleID))
	if err != nil {
		return err
	}

	return nil
}

// getTotalHours calculates total hours for given topics.
func getTotalHours(topics *[]batch.SessionTopic) (totalHours uint) {
	for _, topic := range *topics {
		totalHours += topic.TotalTime
	}
	return
}

// getModuleDatesArray will generate dates array to allocate sub-topics.
func getModuleDatesArray(module *batch.ModuleDTO, initialDate string, totalHours uint,
	dayTimingMap map[uint]time.Duration) ([]string, error) {

	var totalWeeklyHours uint
	var datesArray []string

	for _, value := range dayTimingMap {
		totalWeeklyHours += uint(value.Hours())
	}

	// fmt.Println("================ totalHours ->", totalHours, " ================== totalWeeklyHours ->", totalWeeklyHours)

	totalWeeks := math.Ceil(float64(totalHours)/float64(totalWeeklyHours)) + 1 // 1 week added has a buffer week

	// #shailesh
	startDate, err := time.Parse("2006-01-02", (*module.StartDate)[:10])
	if err != nil {
		return nil, err
	}

	datesArray = calculateModuleDates(module, startDate, totalWeeks)

	return datesArray, nil
}

// calculateModuleDates will calculate dates for session plan.
func calculateModuleDates(module *batch.ModuleDTO, startDate time.Time, totalWeeks float64) (datesArray []string) {

	fmt.Println(" ============================== startDate ->", startDate)

	weekStartDate := getFirstDateOfWeek(startDate)
	// fmt.Println(" ============================== weekStartDate ->", weekStartDate)

	// fmt.Println(" ============================== totalWeeks ->", totalWeeks)

	for i := 0; i < int(totalWeeks); i++ {
		if i > 0 {
			weekStartDate = weekStartDate.AddDate(0, 0, 7)
		}
		for j := 0; j < len(module.ModuleTiming); j++ {
			if i == 0 && startDate.After(weekStartDate.AddDate(0, 0, int(module.ModuleTiming[j].Day.Order))) {
				// fmt.Println(" === found === ")
				continue
			}
			date := weekStartDate.AddDate(0, 0, int(module.ModuleTiming[j].Day.Order)-1).Format("2006-01-02")
			datesArray = append(datesArray, date)
		}
	}

	return
}

// getModuleTimeDifference will calculate time difference for specified days.
func getModuleTimeDifference(module *batch.ModuleDTO, dayTimingMap map[uint]time.Duration) error {

	var fromTime, toTime time.Time
	var err error

	for _, timing := range module.ModuleTiming {

		// fmt.Println(" ====== fromTime -> ", timing.FromTime)
		// fmt.Println(" ====== toTime -> ", timing.ToTime)

		fromTime, err = time.Parse("15:04:05", timing.FromTime)
		if err != nil {
			return err
		}

		toTime, err = time.Parse("15:04:05", timing.ToTime)
		if err != nil {
			return err
		}

		if toTime.Sub(fromTime) <= 0 {
			return errors.NewValidationError("to time should be greater than from time")
		}

		dayTimingMap[timing.Day.Order] = toTime.Sub(fromTime)
	}

	return nil
}

// getFirstDateOfWeek will return first day of the week for current date.
func getFirstDateOfWeek(startDate time.Time) (weekStartDate time.Time) {

	offset := int(time.Monday - startDate.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return
}

func removeCompletedSessionTopics(sessionTopics *[]batch.SessionTopic) {
	for i := 0; i < len(*sessionTopics); i++ {
		if (*sessionTopics)[i].IsCompleted != nil && *(*sessionTopics)[i].IsCompleted {
			*sessionTopics = append((*sessionTopics)[:i], (*sessionTopics)[i+1:]...)
			i = i - 1
		}
	}
}

// getNextStartDate will add 1 day to current date(estimated end date of previous module).
func getNextStartDate(date string) (*string, error) {
	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	date = startDate.AddDate(0, 0, 1).Format("2006-01-02")
	return &date, nil
}

func isSessionParallel(index int, params batch.SessionParameter) (bool, error) {

	if index > 0 {
		previousEndDate, err := time.Parse("2006-01-02", (*(*params.Modules)[index-1].EstimatedEndDate)[:10])
		if err != nil {
			return false, err
		}
		currentStartDate, err := time.Parse("2006-01-02", (*(*params.Modules)[index-1].StartDate)[:10])
		if err != nil {
			return false, err
		}

		return previousEndDate.After(currentStartDate), nil
	}

	return false, nil
}

// ==================================== create session plan private methods ====================================

// newGetCourseModules will get all modules from course_modules for specified batch.
func (service *SessionService) newGetCourseModules(uow *repository.UnitOfWork,
	param batch.GetModuleParams, onExit func()) {

	defer onExit()

	var queryProcessor = []repository.QueryProcessor{
		repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`module_id` = modules.`id` AND " +
			" batch_session_topics.`tenant_id` = modules.`tenant_id`"),
		repository.Join("INNER JOIN batch_modules ON batch_modules.`module_id` = modules.`id` AND" +
			" batch_modules.`tenant_id` = modules.`tenant_id`"), service.addSearchQueries(param.Parser.Form),
		service.addSessionTopicSearchQueries(param.Parser.Form, param.IsPending),
		repository.Filter("modules.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL AND"+
			" batch_modules.`deleted_at` IS NULL", param.TenantID),
		repository.Filter("batch_session_topics.`batch_id` = ?", param.BatchID),
		repository.PreloadAssociations([]string{"Resources"}), repository.OrderBy("batch_session_topics.`order`"),
		repository.GroupBy("batch_session_topics.`module_id`")}

	if param.IsPending {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`is_completed` = ?", false))
	}

	if !util.IsEmpty(param.SessionDate) {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`completed_date` = ?", param.SessionDate))
	}

	err := service.repo.GetAll(uow, param.Modules, queryProcessor...)
	if err != nil {
		param.Channel <- err
	}

	for index := range *param.Modules {

		err = service.newGetModuleTopics(uow, &(*param.Modules)[index].ModuleTopics, (*param.Modules)[index].ID, param)
		if err != nil {
			param.Channel <- err
		}

		for j := range (*param.Modules)[index].ModuleTopics {

			err = service.newGetModuleSubTopics(uow, &(*param.Modules)[index].ModuleTopics[j].SubTopics,
				(*param.Modules)[index].ID, (*param.Modules)[index].ModuleTopics[j].ID, param)
			if err != nil {
				param.Channel <- err
			}
		}
	}

	param.Channel <- nil
}

// newGetModuleTopics will get all module topics for specified batch and module.
func (service *SessionService) newGetModuleTopics(uow *repository.UnitOfWork, moduleTopics *[]*course.ModuleTopicDTO,
	moduleID uuid.UUID, param batch.GetModuleParams) error {

	var queryProcessor = []repository.QueryProcessor{
		repository.Join("INNER JOIN batch_session_topics ON module_topics.`id` = batch_session_topics.`topic_id` " +
			"AND batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
		repository.Join("INNER JOIN batch_sessions ON batch_session_topics.`batch_session_id` = batch_sessions.`id` " +
			"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_session_topics.`batch_id` = ? AND batch_session_topics.`deleted_at` IS NULL", param.BatchID),
		repository.Filter("batch_session_topics.`module_id` = ? AND module_topics.`tenant_id` = ?", moduleID, param.TenantID),
		service.addSessionTopicSearchQueries(param.Parser.Form, param.IsPending), service.addSearchQueries(param.Parser.Form),
		repository.OrderBy("batch_session_topics.`order`"),
		repository.PreloadAssociations([]string{
			"TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
		}), repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchTopicAssignment",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("batch_topic_assignments.`batch_id`= ? AND batch_topic_assignments.`assigned_date` IS NULL",
					param.BatchID), repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
			}}), repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchSessionTopic",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("batch_session_topics.`batch_id`=?", param.BatchID),
				repository.OrderBy("batch_session_topics.`order`"),
			}}), repository.GroupBy("batch_session_topics.`topic_id`")}

	if param.IsPending {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`is_completed` = ?", false))
	}

	if !util.IsEmpty(param.SessionDate) {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`completed_date` = ?", param.SessionDate))
	}

	err := service.repo.GetAll(uow, moduleTopics, queryProcessor...)
	if err != nil {
		return err
	}

	return nil
}

// newGetModuleSubTopics will get all sub topics for specified batch and module.
func (service *SessionService) newGetModuleSubTopics(uow *repository.UnitOfWork, moduleTopics *[]*course.ModuleTopicDTO,
	moduleID, topicID uuid.UUID, param batch.GetModuleParams) error {

	// tenantID, batchID, isPending bool, parser *web.Parser

	var queryProcessor = []repository.QueryProcessor{repository.Join("INNER JOIN batch_session_topics ON module_topics.`id` = batch_session_topics.`sub_topic_id` " +
		"AND batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
		repository.Join("INNER JOIN batch_sessions ON batch_session_topics.`batch_session_id` = batch_sessions.`id` " +
			"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Filter("batch_session_topics.`module_id` = ? AND module_topics.`tenant_id` = ?",
			moduleID, param.TenantID), repository.Filter("batch_session_topics.`deleted_at` IS NULL"),
		repository.Filter("batch_session_topics.`topic_id` = ?", topicID), service.addSearchQueries(param.Parser.Form),
		repository.Filter("batch_session_topics.`batch_id` = ?", param.BatchID),
		service.addSessionTopicSearchQueries(param.Parser.Form, param.IsPending),
		repository.OrderBy("batch_session_topics.`order`"),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchSessionTopic",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("batch_session_topics.`batch_id`= ?",
					param.BatchID), repository.OrderBy("batch_session_topics.`order`"),
			}}),
		repository.PreloadAssociations([]string{
			"TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
		}), repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchTopicAssignment",
			Queryprocessors: []repository.QueryProcessor{
				repository.Filter("batch_topic_assignments.`batch_id`=?", param.BatchID),
				repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
			}})}

	if param.IsPending {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`is_completed` = ?", false))
	}

	if !util.IsEmpty(param.SessionDate) {
		queryProcessor = append(queryProcessor, repository.Filter("batch_session_topics.`completed_date` = ?", param.SessionDate))
	}

	err := service.repo.GetAll(uow, moduleTopics, queryProcessor...)
	if err != nil {
		return err
	}

	return nil
}

// getBatchAndModules will get batch and its related modules
// func (service *SessionService) getBatchAndModules(uow *repository.UnitOfWork,
// 	tenantID, batchID uuid.UUID, params batch.SessionParameter) error {
// 	// , modules *[]batch.ModuleDTO

// 	// err := service.getBatch(uow, params.BatchID, tenantID, batchID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// err := service.getBatchModules(uow, params.Modules, params, tenantID, batchID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// getBatch will get specified batch record.
// func (service *SessionService) getBatch(uow *repository.UnitOfWork, tempBatch *batch.Batch,
// 	tenantID, batchID uuid.UUID) error {

// 	err := service.repo.GetRecordForTenant(uow, tenantID, tempBatch,
// 		repository.Filter("`id` = ?", batchID))
// 	// repository.PreloadWithCustomCondition(repository.Preload{
// 	// 	Schema: "Timing",
// 	// 	Queryprocessors: []repository.QueryProcessor{
// 	// 		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
// 	// 		repository.Filter("batch_timing.`tenant_id` = ? AND days.`deleted_at` IS NULL", tenantID),
// 	// 		repository.OrderBy("days.`order`"),
// 	// 	},
// 	// }), repository.PreloadAssociations([]string{"Timing.Day"})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// getBatchModules will get all batch_modules for specified batch.
func (service *SessionService) getBatchModules(uow *repository.UnitOfWork, params batch.SessionParameter) error {

	var queryProcessor []repository.QueryProcessor

	queryProcessor = append(queryProcessor, repository.Filter("batch_modules.`batch_id` = ?",
		params.BatchID), repository.PreloadAssociations([]string{"Module"}),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ModuleTiming",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND" +
					" days.`tenant_id` = batch_module_timings.`tenant_id`"),
				repository.Filter("batch_module_timings.`tenant_id` = ? AND batch_module_timings.`deleted_at` IS NULL", params.TenantID),
				repository.OrderBy("days.`order`"), repository.PreloadAssociations([]string{"Day"}),
			}}))

	if params.IsCreate {
		queryProcessor = append(queryProcessor, service.addSearchQueries(params.Parser.Form))
	}

	if !params.IsCreate {
		queryProcessor = append(queryProcessor, repository.Filter("batch_modules.`estimated_end_date` IS NOT NULL"))
	}

	err := service.repo.GetAllInOrderForTenant(uow, params.TenantID, params.Modules, "`order`", queryProcessor...)
	if err != nil {
		return err
	}

	return nil
}

func (service *SessionService) getBatchSessionTopics(uow *repository.UnitOfWork, params batch.SessionParameter) error {

	err := service.repo.GetAllInOrder(uow, params.SessionTopics, "`batch_session_topics`.`order`",
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND "+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"), repository.Filter("batch_sessions.`faculty_id` = ? AND "+
			"batch_sessions.`deleted_at` IS NULL AND batch_session_topics.`tenant_id` = ?", params.FacultyID, params.TenantID),
		repository.Filter("`batch_session_topics`.`batch_id` = ? AND `batch_session_topics`.`is_completed` = ?",
			params.BatchID, false))
	if err != nil {
		return err
	}
	return nil
}

func (service *SessionService) getAttendanceCompleted(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	batchSession *batch.SessionDTO) error {

	var err error

	batchSession.IsAttendanceMarked, err = repository.DoesRecordExistForTenant(service.db, tenantID, batch.BatchSessionTalent{},
		repository.Filter("`batch_id` = ? AND  `batch_session_id`=?", batchID, batchSession.ID))
	if err != nil {
		return err
	}

	return nil
}

func (service *SessionService) getFeedbackGivenFlag(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	batchSession *batch.SessionDTO) error {

	var totalCount int

	err := service.getBatchTalentCount(uow, tenantID, batchID, batchSession.ID, &totalCount)
	if err != nil {
		return err
	}

	err = service.repo.Scan(uow, batchSession, repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select("( CASE WHEN COUNT(DISTINCT `talent_id`) = ? THEN true ELSE false END ) AS is_feedback_given",
			totalCount), repository.Filter("`batch_id` = ? AND `batch_session_id` = ?", batchID, batchSession.ID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL AND "+
			"faculty_talent_batch_session_feedback.`tenant_id` = ?", tenantID))
	if err != nil {
		return err
	}

	return nil
}

func (service *SessionService) getBatchTalentCount(uow *repository.UnitOfWork, tenantID,
	batchID, sessionID uuid.UUID, totalCount *int) error {

	// join required as when talent is deleted no action is taken in batch_talents table.
	err := service.repo.GetCount(uow, batch.MappedTalent{}, totalCount,
		repository.Join("INNER JOIN talents ON talents.`id` = batch_talents.`talent_id` AND"+
			" talents.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN batch_session_talents ON talents.`id` = batch_session_talents.`talent_id` AND"+
			" batch_session_talents.`tenant_id` = talents.`tenant_id`"), repository.Filter("batch_talents.`tenant_id` = ?", tenantID),
		repository.Filter("batch_talents.`batch_id` = ? AND batch_talents.`deleted_at` IS NULL", batchID),
		repository.Filter("talents.`deleted_at` IS NULL AND batch_talents.`is_active` = ?", true),
		repository.Filter("batch_talents.`suspension_date` IS NULL AND batch_session_talents.`batch_session_id`=?", sessionID),
		repository.Filter("batch_session_talents.`deleted_at` IS NULL AND batch_session_talents.`is_present` = ?", true))
	if err != nil {
		return err
	}

	return nil
}

func (service *SessionService) getSessionCount(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	totalCount *uint, parser *web.Parser) error {

	err := service.repo.GetCountForTenant(uow, tenantID, &batch.Session{}, totalCount,
		repository.Filter("batch_sessions.`deleted_at` IS NULL"), service.addSearchQueries(parser.Form),
		repository.Filter("batch_sessions.`batch_id`=?", batchID))
	if err != nil {
		return err
	}
	return nil
}

func (service *SessionService) getSessionModuleCount(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	totalCount *uint, parser *web.Parser) error {

	err := service.repo.GetCount(uow, &batch.SessionTopic{}, totalCount,
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND"+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"), service.addSearchQueries(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_topics.`batch_id`=?", batchID), repository.GroupBy("batch_session_topics.`module_id`"))
	if err != nil {
		return err
	}
	return nil
}

func (service *SessionService) getSessionTopicCount(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	totalCount *uint, parser *web.Parser) error {

	err := service.repo.GetCount(uow, &batch.SessionTopic{}, totalCount,
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND "+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"), service.addSearchQueries(parser.Form),
		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_session_topics.`batch_id`=?", batchID), repository.GroupBy("`topic_id`"))
	if err != nil {
		return err
	}
	return nil
}

func (service *SessionService) getSessionAssignmentCount(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	totalCount *uint, parser *web.Parser) error {

	err := service.repo.GetCount(uow, &batch.TopicAssignment{}, totalCount,
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`batch_id` = batch_topic_assignments.`batch_id` AND "+
			" batch_sessions.`tenant_id` = batch_topic_assignments.`tenant_id`"),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_topic_assignments.`tenant_id` = ?", tenantID),
		service.addSearchQueries(parser.Form), repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
		repository.GroupBy("batch_topic_assignments.`id`"))
	if err != nil {
		return err
	}
	return nil
}

func (service *SessionService) getTotalBatchHours(uow *repository.UnitOfWork, batchSessions *batch.SessionCounts,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.repo.Scan(uow, &batchSessions, repository.Table("batch_session_topics"),
		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_batch_hours`"),
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` AND "+
			" batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"), service.addSearchQueries(parser.Form),
		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_session_topics.`tenant_id`=?", tenantID),
		repository.Filter("batch_sessions.`deleted_at` IS NULL AND batch_session_topics.`batch_id`=?", batchID))
	if err != nil {
		return err
	}
	return nil
}

// getBatchProjectCount will count total projects assigned for specified batch.
func (service *SessionService) getBatchProjectCount(uow *repository.UnitOfWork, tenantID, batchID uuid.UUID,
	totalCount *uint, parser *web.Parser) error {

	// , service.addSearchQueries(parser.Form)

	err := service.repo.GetCountForTenant(uow, tenantID, &batch.Project{}, totalCount,
		repository.Filter("batch_projects.`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	return nil
}

// addSearchQueries adds search criteria.
func (service *SessionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if isCompleted, ok := requestForm["isCompleted"]; ok {
		util.AddToSlice("`batch_session_topics`.`is_completed`", "= ?", "AND", isCompleted, &columnNames, &conditions, &operators, &values)
	}

	if moduleID, ok := requestForm["moduleID"]; ok {
		util.AddToSlice("`batch_session_topics`.`module_id`", "IN(?)", "AND", moduleID, &columnNames, &conditions, &operators, &values)
	}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addSessionTopicSearchQueries adds search criteria.
func (service *SessionService) addSessionTopicSearchQueries(requestForm url.Values, isPending bool) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if !isPending {
		if sessionTopicCompletedDate, ok := requestForm["sessionTopicCompletedDate"]; ok {
			util.AddToSlice("CAST(`batch_session_topics`.`completed_date` AS DATE)", "= ?", "AND", sessionTopicCompletedDate, &columnNames, &conditions, &operators, &values)
		}
	}

	if isPending {
		if pendingCompletedDate, ok := requestForm["pendingCompletedDate"]; ok {
			util.AddToSlice("CAST(`batch_session_topics`.`completed_date` AS DATE)", "< ?", "AND", pendingCompletedDate, &columnNames, &conditions, &operators, &values)
		}
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

func (service *SessionService) checkForeignkey(params batch.SessionParameter) error {

	err := service.doesForeignKeyExist(params.TenantID, params.CredentialID, params.BatchID)
	if err != nil {
		return err
	}

	// check if session topic foreign keys exist.
	err = service.doesSessionForeignKeyExist(params.TenantID, params.BatchID, params.SessionTopics)
	if err != nil {
		return err
	}

	// if isAdd {

	// check if order is valid.
	// err = checkBatchSessionTopicOrder(params.SessionTopics)
	// if err != nil {
	// 	return err
	// }

	// }

	return nil
}

// doesForeignKeyExist checks if all foreign keys are valid.
func (service *SessionService) doesForeignKeyExist(tenantID, credentialID, batchID uuid.UUID) error {

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

	return nil
}

func (service *SessionService) doesSessionForeignKeyExist(tenantID, batchID uuid.UUID,
	sessionTopics *[]batch.SessionTopic) error {

	isModuleChanged := false
	isTopicChanged := false

	for index, sessionTopic := range *sessionTopics {

		isModuleChanged = false
		isTopicChanged = false

		if index != 0 && sessionTopic.ModuleID == (*sessionTopics)[index-1].ModuleID {
			isModuleChanged = true
		}

		if index != 0 && sessionTopic.TopicID == (*sessionTopics)[index-1].TopicID {
			isModuleChanged = true
		}

		fmt.Println(" ============================== isModuleChanged ->", isModuleChanged)
		fmt.Println(" ============================== isTopicChanged ->", isTopicChanged)

		if !isModuleChanged {
			// check if module exist
			err := service.doesModuleExist(tenantID, sessionTopic.ModuleID)
			if err != nil {
				return err
			}

			// check if batch module exist
			err = service.doesBatchModuleExists(tenantID, batchID, sessionTopic.ModuleID)
			if err != nil {
				return err
			}
		}

		if !isTopicChanged {
			// check if module topic exist
			err := service.doesModuleTopicExist(tenantID, sessionTopic.ModuleID, sessionTopic.TopicID)
			if err != nil {
				return err
			}

		}

		// check if sub topic exist
		err := service.doesModuleTopicExist(tenantID, sessionTopic.ModuleID, sessionTopic.SubTopicID)
		if err != nil {
			return err
		}

	}
	return nil
}

// check if unassigned module exist.
// func (service *SessionService) doesUnassignedModuleExist(params batch.SessionParameter) error {

// 	for _, sessionTopic := range *params.SessionTopics {
// 		isModuleFound := false
// 		for _, module := range *params.Modules {
// 			if module.ModuleID == sessionTopic.ModuleID {
// 				isModuleFound = true
// 				break
// 			}
// 		}
// 		if !isModuleFound {
// 			return errors.NewValidationError("session plan cannot be created for modules that are not assigned to faculty")
// 		}
// 	}
// 	return nil
// }

func (service *SessionService) doesModuleExist(tenantID, module uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, course.Module{},
		repository.Filter("`id` = ?", module))
	if err := util.HandleError("Module not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *SessionService) doesModuleTopicExist(tenantID, moduleID, topicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, course.ModuleTopic{},
		repository.Filter("`id` = ? AND `module_id` = ?", topicID, moduleID))
	if err := util.HandleError("Topic not found", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *SessionService) doesBatchModuleExists(tenantID, batchID, moduleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Module{},
		repository.Filter("`batch_id` = ? AND `module_id` = ?", batchID, moduleID))
	if err := util.HandleError("Module not found in batch", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExists validates tenantID
func (service *SessionService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.db, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *SessionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchExists validates batchID
func (service *SessionService) doesBatchExists(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFacultyExists validates facultyID
func (service *SessionService) doesFacultyExists(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, faculty.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Faculty not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSessionPlanForBatchExist will check if session plan exist for specified batch.
// func (service *SessionService) doesSessionPlanForBatchExist(tenantID, batchID, sessionID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Session{},
// 		repository.Filter("`batch_id` = ?", batchID))
// 	if err := util.HandleIfExistsError("Session plan for batch already exist", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (service *SessionService) doesBatchSessionTopicExists(tenantID, batchID, topicID, subTopicID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.SessionTopic{},
		repository.Filter("`batch_id` = ? AND `topic_id` = ? AND`sub_topic_id` = ?", batchID, topicID, subTopicID))
	if err := util.HandleError("Topic not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSessionPlanForFacultyExist will check if session plan exist for specified batch.
func (service *SessionService) doesSessionPlanForFacultyExist(tenantID, batchID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Session{},
		repository.Filter("`batch_id` = ? AND `faculty_id` = ?", batchID, facultyID))
	if err := util.HandleIfExistsError("Session plan for batch already exist", exists, err); err != nil {
		return err
	}
	return nil
}

// checkBatchSessionTopicOrder will check topic order.
// func checkBatchSessionTopicOrder(sessionTopcis *[]batch.SessionTopic) error {

// 	// orderMap := make(map[uint]uint)
// 	var moduleStartIndex, moduleEndIndex int = 0, 0
// 	var topicStartIndex, topicEndIndex int = 0, 0
// 	var moduleID, topicID uuid.UUID

// 	for index, sessionTopic := range *sessionTopcis {
// 		// fmt.Printf("Order for %d sessionTopic is -> %d\n", (index + 1), sessionTopic.Order)

// 		if (index + 1) != int(sessionTopic.Order) {
// 			return errors.NewValidationError("invalid order given.")
// 		}
// 	}

// 	moduleID = (*sessionTopcis)[0].ModuleID
// 	topicID = (*sessionTopcis)[0].TopicID

// 	// validating module order.
// 	for i := 1; i < len(*sessionTopcis); i++ {
// 		if (*sessionTopcis)[i].TopicID == topicID {
// 			topicEndIndex = i
// 			continue
// 		}
// 		checkModuleOrder(moduleStartIndex, moduleEndIndex, moduleID, sessionTopcis)
// 		moduleStartIndex, moduleEndIndex = i, i
// 		topicID = (*sessionTopcis)[i].TopicID
// 	}

// 	// validating topic order.
// 	for i := 1; i < len(*sessionTopcis); i++ {
// 		if (*sessionTopcis)[i].TopicID == topicID {
// 			topicEndIndex = i
// 			continue
// 		}
// 		checkTopicOrder(topicStartIndex, topicEndIndex, topicID, sessionTopcis)
// 		// fmt.Println("topicstartindex ->", topicStartIndex, " topicendindex ->", topicEndIndex, " topicID ->", topicID)
// 		topicStartIndex, topicEndIndex = i, i
// 		topicID = (*sessionTopcis)[i].TopicID
// 	}

// 	return nil
// }

// func checkModuleOrder(startIndex, endIndex int,
// 	moduleID uuid.UUID, sessionTopics *[]batch.SessionTopic) error {

// 	for i := 0; i < len(*sessionTopics); i++ {
// 		if !(i <= startIndex && i <= endIndex) {
// 			if (*sessionTopics)[i].ModuleID == moduleID {
// 				return errors.NewValidationError("invalid module order")
// 			}
// 		}
// 	}

// 	return nil
// }

// func checkTopicOrder(startIndex, endIndex int,
// 	topicID uuid.UUID, sessionTopics *[]batch.SessionTopic) error {

// 	for i := 0; i < len(*sessionTopics); i++ {
// 		if !(i <= startIndex && i <= endIndex) {
// 			if (*sessionTopics)[i].TopicID == topicID {
// 				return errors.NewValidationError("invalid topic order")
// 			}
// 		}
// 	}

// 	return nil
// }
