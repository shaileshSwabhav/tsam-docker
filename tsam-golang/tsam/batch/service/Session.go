package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/web"
)

// SessionService provide method like Add, Update, Delete, Get for batch_sessions.
type SessionService struct {
	db          *gorm.DB
	repo        repository.Repository
	association []string
}

// NewSessionService creates a new instance of BatchSessionService.
func NewSessionService(db *gorm.DB, repo repository.Repository) *SessionService {
	return &SessionService{
		db:          db,
		repo:        repo,
		association: []string{},
	}
}

func (service *SessionService) GetModuleTopics(tenantID, batchID, moduleID uuid.UUID, moduleTopics *[]batch.ModuleTopics, parser *web.Parser) error {
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

	// SELECT m.id, m.topic_name, b.id FROM module_topics AS m
	// LEFT JOIN batch_session_topics AS b ON m.id = b.topic_id
	// WHERE m.module_id = "fcf4a3e1-eff4-4fc0-a3b7-6013f3d5b552"
	// GROUP BY m.id;

	batchSessionTopic := "module_topics.*, b.`id` AS `batch_session_topic_id`, b.`batch_session_id` AS batch_session_id," +
		"b.`order` AS batch_session_topic_order,b.`is_completed` AS batch_session_topic_is_completed,b.`total_time` AS batch_session_topic_total_time," +
		"b.`initial_date` AS batch_session_initial_date,b.`completed_date` AS batch_session_completion_date"

	err = service.repo.GetAll(uow, moduleTopics, repository.Table("module_topics"),
		repository.Join("LEFT JOIN batch_session_topics AS b ON module_topics.`id` = b.topic_id"+
			" AND `b`.`tenant_id` = `module_topics`.`tenant_id` AND b.`deleted_at` IS NULL"),
		repository.Filter("module_topics.`topic_id` IS NULL AND module_topics.`module_id`= ?", moduleID),
		repository.Filter("module_topics.`tenant_id`= ? AND b.`deleted_at` IS NULL AND module_topics.`deleted_at` IS NULL", tenantID),
		repository.GroupBy("module_topics.`id`"),
		repository.OrderBy("module_topics.`order`"),
		repository.Select(batchSessionTopic))
	if err != nil {
		return err
	}
	for i := 0; i < len(*moduleTopics); i++ {
		// moduleSubTopics := []batch.ModuleTopics{}

		batchSessionTopic := "module_topics.*, b.`id` AS `batch_session_topic_id`, b.`batch_session_id` AS batch_session_id," +
			"b.`order` AS batch_session_topic_order,b.`is_completed` AS batch_session_topic_is_completed,b.`total_time` AS batch_session_topic_total_time," +
			"b.`initial_date` AS batch_session_initial_date,b.`completed_date` AS batch_session_completion_date"

		// subquery, err := service.repo.SubQuery(uow, &batch.SessionTopicDTO{},
		// 	repository.Model(&batch.SessionTopicDTO{}),
		// 	repository.Filter("batch_id=?", batchID))
		// if err != nil {
		// 	return err
		// }

		err = service.repo.GetAll(uow, &(*moduleTopics)[i].SubTopics, repository.Table("module_topics"),
			repository.Join("LEFT JOIN batch_session_topics AS b ON module_topics.`id` = b.`sub_topic_id` AND "+
				"`b`.`tenant_id` = `module_topics`.`tenant_id` AND b.`deleted_at` IS NULL"),
			repository.Filter("module_topics.`module_id`=? AND module_topics.`topic_id`=?", moduleID, (*moduleTopics)[i].ID), // AND b.`batch_id` = ?, batchID
			repository.GroupBy("module_topics.`id`"),
			repository.OrderBy("module_topics.`order`"),
			repository.Select(batchSessionTopic))
		// repository.Filter("`b`.batch_id=?", batchID),
		if err != nil {
			return err
		}

		// (*moduleTopics)[i].SubTopics = moduleSubTopics
	}
	uow.Commit()
	return nil
}

// SkipPendingSession will skip the pending sessions to next date.
// func (service *SessionService) SkipPendingSession(tenantID, credentialID, batchID uuid.UUID,
// 	parser *web.Parser) error {

// 	// check if foreign key exist
// 	err := service.doesForeignKeyExist(tenantID, credentialID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.db, false)

// 	params := batch.SessionParameter{
// 		BatchSessions:       &[]batch.Session{},
// 		Modules:             &[]batch.ModuleDTO{},
// 		ModuleBatchSessions: &[]batch.Session{},
// 		ModuleSessionTopics: &[]batch.SessionTopic{},
// 		BatchID:             batchID,
// 		SessionTopics:       &[]batch.SessionTopic{},
// 		Parser:              parser,
// 		FacultyID:           uuid.Nil,
// 	}
// 	// err = service.getBatchAndModules(uow, tenantID, batchID, params)
// 	err = service.getBatchModules(uow, params, tenantID)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	err = service.repo.GetAllInOrderForTenant(uow, tenantID, params.SessionTopics, "`batch_session_topics`.`order`",
// 		repository.Filter("`batch_session_topics`.`batch_id` = ? AND `batch_session_topics`.`is_completed` = ?", batchID, false))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	err = service.repo.GetAllInOrderForTenant(uow, tenantID, params.BatchSessions, "`batch_sessions`.`date`",
// 		repository.Filter("`batch_sessions`.`batch_id` = ? AND `batch_sessions`.`is_session_taken` = ?", batchID, false))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// err = service.skipSessionPlan(uow, params)
// 	// if err != nil {
// 	// 	uow.RollBack()
// 	// 	return err
// 	// }

// 	err = service.updateSessionTopicPlan(uow, params, tenantID, credentialID)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// UpdateBatchSession will add/update topics to session.
func (service *SessionService) UpdateBatchSession(tenantID, credentialID, batchID uuid.UUID,
	sessionTopics *[]batch.SessionTopic, parser *web.Parser) error {

	removeCompletedSessionTopics(sessionTopics)

	params := batch.SessionParameter{
		SessionTopics: sessionTopics,
		InitialDate:   (*sessionTopics)[0].CompletedDate,
		Modules:       &[]batch.ModuleDTO{},
		BatchSessions: &[]batch.Session{},
		Parser:        parser,
		BatchID:       batchID,
		TenantID:      tenantID,
		CredentialID:  credentialID,
	}

	// check foreign keys
	err := service.checkForeignkey(params)
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

	// err = service.createSessionPlan(params)
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

	if len(*params.BatchSessions) == 0 {
		return errors.NewValidationError("batch session could not be created.")
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// createSessionPlan will create a session plan for specified batch.
// func (service *SessionService) createSessionPlan(params batch.SessionParameter) error {

// 	dayTimingMap := make(map[uint]time.Duration)

// 	err := getBatchTimeDifference(params.BatchID, dayTimingMap)
// 	if err != nil {
// 		return err
// 	}

// 	batchTotalHours := getTotalHours(params.SessionTopics)
// 	datesArray, err := getBatchDatesArray(params.BatchID, *params.BatchID.StartDate,
// 		batchTotalHours, dayTimingMap)
// 	if err != nil {
// 		return err
// 	}

// 	err = generateSessionPlan(datesArray, dayTimingMap, params, true)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func generateSessionPlan(datesArray []string, dayTimingMap map[uint]time.Duration, params batch.SessionParameter,
// 	isGenerate bool) error {

// 	fmt.Println(" ============================= dayTimingMap ->", dayTimingMap)

// 	var tempSessionTopics []batch.SessionTopic
// 	var tempBatchSession batch.Session
// 	// var currentFacultyID uuid.UUID

// 	// moduleFacultyMap := getModuleFacultyMap(params.Modules)
// 	// fmt.Println(" ========================= moduleFacultyMap ->", moduleFacultyMap)

// 	for _, date := range datesArray {

// 		if isGenerate {
// 			tempBatchSession = batch.Session{
// 				BatchID: params.BatchID,
// 				Date:    date,
// 			}
// 		}

// 		currentDate, err := time.Parse("2006-01-02", date)
// 		if err != nil {
// 			return err
// 		}

// 		currentDay := int(currentDate.Weekday())
// 		allocationDurationLeft := dayTimingMap[uint(currentDay)].Minutes()

// 		if len(*params.SessionTopics) == 0 {
// 			break
// 		}

// 		// currentFacultyID = moduleFacultyMap[(*params.SessionTopics)[0].ModuleID]

// 		allocateParams := batch.AllocateParameter{
// 			SessionTopics:          &tempSessionTopics,
// 			BatchSession:           &tempBatchSession,
// 			AllocationDurationLeft: allocationDurationLeft,
// 			CurrentDate:            currentDate,
// 			// ModuleFacultyMap:       moduleFacultyMap,
// 			// CurrentFacultyID:       currentFacultyID,
// 		}

// 		allocateSessionTopic(params, allocateParams, isGenerate)

// 		if isGenerate {
// 			if len(tempBatchSession.BatchSessionTopic) > 0 {
// 				*params.BatchSessions = append(*params.ModuleBatchSessions, tempBatchSession)
// 			}
// 		}
// 	}

// 	if !isGenerate {
// 		*params.SessionTopics = tempSessionTopics
// 	}

// 	return nil
// }

// func allocateSessionTopic(params batch.SessionParameter, allocateParams batch.AllocateParameter, isGenerate bool) {
// 	for i := 0; i < len(*params.SessionTopics); i++ {

// 		// if allocateParams.ModuleFacultyMap[(*params.SessionTopics)[i].ModuleID] != allocateParams.CurrentFacultyID {
// 		// 	allocateParams.CurrentFacultyID = allocateParams.ModuleFacultyMap[(*params.SessionTopics)[i].ModuleID]
// 		// 	return
// 		// }

// 		allocateParams.AllocationDurationLeft -= float64((*params.SessionTopics)[i].TotalTime)

// 		if allocateParams.AllocationDurationLeft >= 0.0 {
// 			completedDate := allocateParams.CurrentDate.Format("2006-01-02")
// 			(*params.SessionTopics)[i].CompletedDate = &completedDate

// 			if !isGenerate {
// 				(*allocateParams.SessionTopics) = append((*allocateParams.SessionTopics), (*params.SessionTopics)[i])
// 			}

// 			if isGenerate {
// 				allocateParams.BatchSession.BatchSessionTopic = append(allocateParams.BatchSession.BatchSessionTopic, (*params.SessionTopics)[i])
// 			}

// 			(*params.SessionTopics) = append((*params.SessionTopics)[:i], (*params.SessionTopics)[i+1:]...)
// 			i = i - 1
// 			continue
// 		}
// 		return
// 	}
// }

// skipSessionPlan will update a session plan for specified batch.
// func (service *SessionService) skipSessionPlan(uow *repository.UnitOfWork,
// 	params batch.SessionParameter) error {

// 	dayTimingMap := make(map[uint]time.Duration)

// 	err := getBatchTimeDifference(params.BatchID, dayTimingMap)
// 	if err != nil {
// 		return err
// 	}

// 	batchTotalHours := getTotalHours(params.SessionTopics)

// 	datesArray, err := getBatchDatesArray(params.BatchID, *params.BatchID.StartDate,
// 		batchTotalHours, dayTimingMap)
// 	if err != nil {
// 		return err
// 	}

// 	today, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
// 	if err != nil {
// 		return err
// 	}

// 	for i := 0; i < len(datesArray); i++ {
// 		date, err := time.Parse("2006-01-02", datesArray[i])
// 		if err != nil {
// 			return err
// 		}

// 		// splice date if it is before today's date.
// 		if date.Before(today) || today.Equal(date) {
// 			datesArray = append(datesArray[:i], datesArray[i+1:]...)
// 			i = i - 1
// 		}
// 	}

// 	err = generateSessionPlan(datesArray, dayTimingMap, params, false)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// getBatchTimeDifference will calculate time difference for specified days.
// func getBatchTimeDifference(tempBatch *batch.Batch, dayTimingMap map[uint]time.Duration) error {

// 	var fromTime, toTime time.Time
// 	var err error

// 	for _, timing := range tempBatch.Timing {

// 		if timing.FromTime != nil {
// 			fromTime, err = time.Parse("15:04:05", *timing.FromTime)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if timing.ToTime != nil {
// 			toTime, err = time.Parse("15:04:05", *timing.ToTime)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		dayTimingMap[timing.Day.Order] = toTime.Sub(fromTime)
// 	}

// 	return nil
// }

// // getBatchDatesArray will generate dates array to allocate sub-topics.
// func getBatchDatesArray(tempBatch *batch.Batch, initialDate string, totalHours uint,
// 	dayTimingMap map[uint]time.Duration) ([]string, error) {

// 	var totalWeeklyHours uint
// 	var datesArray []string

// 	for _, value := range dayTimingMap {
// 		totalWeeklyHours += uint(value.Hours())
// 	}

// 	// fmt.Println("================ totalHours ->", totalHours, " ================== totalWeeklyHours ->", totalWeeklyHours)

// 	totalWeeks := math.Ceil(float64(totalHours)/float64(totalWeeklyHours)) + 1 // 1 week added has a buffer week

// 	startDate, err := time.Parse("2006-01-02", initialDate[:10])
// 	if err != nil {
// 		return nil, err
// 	}

// 	datesArray = calculateBatchDates(tempBatch, startDate, totalWeeks)

// 	return datesArray, nil
// }

// // calculateBatchDates will calculate dates for session plan.
// func calculateBatchDates(tempBatch *batch.Batch, startDate time.Time, totalWeeks float64) (datesArray []string) {

// 	// fmt.Println(" ============================== startDate ->", startDate)

// 	weekStartDate := getFirstDateOfWeek(startDate)
// 	// fmt.Println(" ============================== weekStartDate ->", weekStartDate)

// 	// fmt.Println(" ============================== totalWeeks ->", totalWeeks)

// 	for i := 0; i < int(totalWeeks); i++ {
// 		if i > 0 {
// 			weekStartDate = weekStartDate.AddDate(0, 0, 7)
// 		}
// 		for j := 0; j < len(tempBatch.Timing); j++ {
// 			if i == 0 && startDate.After(weekStartDate.AddDate(0, 0, int(tempBatch.Timing[j].Day.Order))) {
// 				continue
// 			}
// 			date := weekStartDate.AddDate(0, 0, int(tempBatch.Timing[j].Day.Order)-1).Format("2006-01-02")
// 			datesArray = append(datesArray, date)
// 		}
// 	}

// 	return
// }

// doesBatchSessionExists validates batchID
// func (service *BatchSessionService) doesBatchSessionExists(tenantID, batchID, batchSessionID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchSession{},
// 		repository.Filter("`id` = ? AND `batch_id` = ?", batchSessionID, batchID))
// 	if err := util.HandleError("Batch session not found", exists, err); err != nil {
// 		return err
// 	}
// 	return nil
// }

// ================================ USE FOR REFERENCE ONLY ================================

// GetBatchSessionPlan will get batch sessions for the specified batch.
// func (service *SessionService) GetBatchSessionPlan(batchSessions *batch.SessionDTO,
// 	tenantID, batchID uuid.UUID, parser *web.Parser) error {

// 	date := parser.Form.Get("date")
// 	// now := time.Now()

// 	// fmt.Println(" ================= date ->", date)

// 	if util.IsEmpty(date) {
// 		return errors.NewValidationError("session date must be specified")
// 	}

// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.db, true)

// 	err = service.repo.GetRecordForTenant(uow, tenantID, batchSessions,
// 		repository.Filter("batch_sessions.`batch_id` = ?", batchID),
// 		// repository.Filter("batch_sessions.`date` = ?", date),
// 		service.addSearchQueries(parser.Form),
// 		repository.PreloadAssociations([]string{"BatchSessionPrerequisiteDTO"}))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	channel := make(chan error, 2)
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go service.getCourseModules(uow, &batchSessions.Module, tenantID, batchID, channel, func() { wg.Done() },
// 		repository.Filter("batch_session_topics.`completed_date` = ?", date))

// 	wg.Add(1)
// 	go service.getCourseModules(uow, &batchSessions.PendingModule, tenantID, batchID, channel, func() { wg.Done() },
// 		repository.Filter("batch_session_topics.`completed_date` < ?", date))

// 	go func() {
// 		defer close(channel)
// 		wg.Wait()
// 	}()

// 	for err := range channel {
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	// defer func() {
// 	// 	fmt.Println(" ===================== duration ->", time.Since(now))
// 	// 	fmt.Println("Go routine Number is --------------------------", runtime.NumGoroutine())
// 	// }()

// 	err = service.getAttendanceCompleted(uow, tenantID, batchID, batchSessions.ID, &batchSessions.IsAttendanceMarked)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// // GetAllBatchSessions will return all batch session and its related topics and subtopics.
// func (service *SessionService) GetAllBatchSessions(batchSessions *[]batch.SessionDTO,
// 	tenantID, batchID uuid.UUID, parser *web.Parser) error {

// 	// now := time.Now()

// 	// check if tenant exist.
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist.
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.db, true)

// 	err = service.repo.GetAllInOrder(uow, batchSessions, "batch_sessions.`date`",
// 		repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id` "+
// 			"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
// 		repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
// 		repository.Filter("batch_sessions.`batch_id` = ?", batchID),
// 		repository.GroupBy("batch_session_topics.`batch_session_id`"),
// 		repository.PreloadAssociations([]string{"BatchSessionPrerequisiteDTO"}))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	channel := make(chan error, len(*batchSessions))
// 	var wg sync.WaitGroup

// 	for index := range *batchSessions {
// 		wg.Add(1)
// 		go service.getCourseModules(uow, &(*batchSessions)[index].Module, tenantID, batchID, channel,
// 			func() { wg.Done() }, repository.Filter("batch_session_topics.`batch_session_id` = ?", (*batchSessions)[index].ID))
// 	}

// 	go func() {
// 		defer close(channel)
// 		wg.Wait()
// 	}()

// 	for err := range channel {
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	// defer func() {
// 	// 	fmt.Println(" ===================== duration ->", time.Since(now))
// 	// 	fmt.Println("Go routine Number is --------------------------", runtime.NumGoroutine())
// 	// }()

// 	uow.Commit()
// 	return nil
// }

// GetBatchSessionsCounts will return all counts for batch session plan for specified batch.
// func (service *SessionService) GetBatchSessionsCounts(batchSessions *batch.SessionCounts,
// 	tenantID, batchID uuid.UUID, parser *web.Parser) error {

// 	// Check if tenant exist.
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if batch exist.
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.db, true)

// 	// Get count of modules.
// 	err = service.repo.GetCountForTenant(uow, tenantID, &batch.SessionTopic{}, &batchSessions.ModuleCount,
// 		repository.Filter("batch_session_topics.`deleted_at` IS NULL"),
// 		repository.Filter("batch_session_topics.`batch_id`=?", batchID), service.addSearchQueries(parser.Form),
// 		repository.GroupBy("`module_id`"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Get count of topics.
// 	err = service.repo.GetCountForTenant(uow, tenantID, &batch.SessionTopic{}, &batchSessions.TopicCount,
// 		repository.Filter("batch_session_topics.`deleted_at` IS NULL"),
// 		repository.Filter("batch_session_topics.`batch_id`=?", batchID),
// 		repository.GroupBy("`topic_id`"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Get count of assignments.
// 	err = service.repo.GetCountForTenant(uow, tenantID, &batch.TopicAssignment{}, &batchSessions.AssignmentCount,
// 		repository.Filter("batch_topic_assignments.`deleted_at` IS NULL"),
// 		repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
// 	)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Get count of batch sessions.
// 	err = service.repo.GetCountForTenant(uow, tenantID, &batch.Session{}, &batchSessions.SessionCount,
// 		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
// 		repository.Filter("batch_sessions.`batch_id`=?", batchID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	// Get total hours of batch.
// 	var sessionTotalHours batch.SessionTotalHours
// 	err = service.repo.Scan(uow, &sessionTotalHours,
// 		repository.Table("batch_session_topics"),
// 		repository.Select("SUM(batch_session_topics.`total_time`) AS `total_hours`"),
// 		repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_session_topics.`tenant_id`=?", tenantID),
// 		repository.Filter("batch_session_topics.`batch_id`=?", batchID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	batchSessions.TotalBatchHours = uint(sessionTotalHours.TotalHours)

// 	err = service.getBatchProjectCount(uow, tenantID, batchID, &batchSessions.ProjectCount)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	uow.Commit()
// 	return nil
// }

// getCourseModules will get all modules from course_modules for specified batch.
// func (service *SessionService) getCourseModules(uow *repository.UnitOfWork, modules *[]course.ModuleDTO,
// 	tenantID, batchID uuid.UUID, channel chan error, onExit func(), queryProcessor ...repository.QueryProcessor) {

// 	defer onExit()

// 	// parser *web.Parser,
// 	// service.addSearchQueries(parser.Form),
// 	tempQueryProcessor := append(queryProcessor, repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`module_id` = modules.`id` AND "+
// 		" batch_session_topics.`tenant_id` = modules.`tenant_id`"),
// 		repository.Filter("modules.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
// 		repository.Filter("batch_session_topics.`batch_id` = ?", batchID), repository.PreloadAssociations([]string{"Resources"}),
// 		repository.OrderBy("batch_session_topics.`order`"), repository.GroupBy("batch_session_topics.`module_id`"))

// 	err := service.repo.GetAll(uow, modules, tempQueryProcessor...)
// 	if err != nil {
// 		channel <- err
// 	}

// 	for index := range *modules {

// 		err = service.getModuleTopics(uow, &(*modules)[index].ModuleTopics, tenantID, batchID,
// 			(*modules)[index].ID, queryProcessor...)
// 		if err != nil {
// 			channel <- err
// 		}

// 		for j := range (*modules)[index].ModuleTopics {

// 			err = service.getModuleSubTopics(uow, &(*modules)[index].ModuleTopics[j].SubTopics,
// 				tenantID, batchID, (*modules)[index].ID, (*modules)[index].ModuleTopics[j].ID, queryProcessor...)
// 			if err != nil {
// 				channel <- err
// 			}
// 		}
// 	}

// 	channel <- nil
// }

// getModuleTopics will get all module topics for specified batch and module.
// func (service *SessionService) getModuleTopics(uow *repository.UnitOfWork, moduleTopics *[]*course.ModuleTopicDTO,
// 	tenantID, batchID, moduleID uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

// 	tempQueryProcessor := append(queryProcessor, repository.Join("INNER JOIN batch_session_topics ON module_topics.`id` = batch_session_topics.`topic_id` "+
// 		"AND batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
// 		repository.Filter("batch_session_topics.`deleted_at` IS NULL"), repository.Filter("batch_session_topics.`batch_id` = ?", batchID),
// 		repository.Filter("batch_session_topics.`module_id` = ? AND module_topics.`tenant_id` = ?", moduleID, tenantID),
// 		repository.PreloadAssociations([]string{
// 			"TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
// 		}), repository.PreloadWithCustomCondition(repository.Preload{
// 			Schema: "BatchTopicAssignment",
// 			Queryprocessors: []repository.QueryProcessor{
// 				repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
// 				repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
// 			}}), repository.PreloadWithCustomCondition(repository.Preload{
// 			Schema: "BatchSessionTopic",
// 			Queryprocessors: []repository.QueryProcessor{
// 				repository.Filter("batch_session_topics.`batch_id`=?", batchID),
// 			}}), repository.OrderBy("batch_session_topics.`order`"), repository.GroupBy("batch_session_topics.`topic_id`"))

// 	err := service.repo.GetAll(uow, moduleTopics, tempQueryProcessor...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// getModuleTopics will get all sub topics for specified batch and module.
// func (service *SessionService) getModuleSubTopics(uow *repository.UnitOfWork, moduleTopics *[]*course.ModuleTopicDTO,
// 	tenantID, batchID, moduleID, topicID uuid.UUID, queryProcessor ...repository.QueryProcessor) error {

// 	tempQueryProcessor := append(queryProcessor, repository.Join("INNER JOIN batch_session_topics ON module_topics.`id` = batch_session_topics.`sub_topic_id` "+
// 		"AND batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
// 		repository.Filter("batch_session_topics.`module_id` = ? AND module_topics.`tenant_id` = ?",
// 			moduleID, tenantID), repository.Filter("batch_session_topics.`deleted_at` IS NULL"),
// 		repository.Filter("batch_session_topics.`topic_id` = ?", topicID), repository.Filter("batch_session_topics.`batch_id` = ?", batchID),
// 		repository.PreloadAssociations([]string{
// 			"BatchSessionTopic", "TopicProgrammingConcept", "TopicProgrammingConcept.ProgrammingConcept",
// 		}), repository.PreloadWithCustomCondition(repository.Preload{
// 			Schema: "BatchTopicAssignment",
// 			Queryprocessors: []repository.QueryProcessor{
// 				repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
// 				repository.PreloadAssociations([]string{"ProgrammingQuestion"}),
// 			},
// 		}), repository.PreloadWithCustomCondition(repository.Preload{
// 			Schema: "BatchSessionTopic",
// 			Queryprocessors: []repository.QueryProcessor{
// 				repository.Filter("batch_session_topics.`batch_id`=?", batchID),
// 			},
// 		}), repository.OrderBy("batch_session_topics.`order`"))

// 	err := service.repo.GetAll(uow, moduleTopics, tempQueryProcessor...)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// GetBatchSessionTopicList will return list of topics and sub-topics for specified batch.
// func (service *SessionService) GetBatchSessionTopicList(tenantID, batchID uuid.UUID,
// 	moduleTopics *[]course.ModuleTopicDTO) error {

// 	// check if tenant exist.
// 	err := service.doesTenantExists(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// check if batch exist.
// 	err = service.doesBatchExists(tenantID, batchID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.db, true)

// 	err = service.repo.GetAllInOrder(uow, moduleTopics, "batch_session_topics.`order`",
// 		repository.Join("INNER JOIN batch_session_topics ON module_topics.`id` = batch_session_topics.`topic_id` "+
// 			"AND module_topics.`tenant_id` = batch_session_topics.`tenant_id`"),
// 		// repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = batch_session_topics.`batch_session_id` "+
// 		// 	"AND batch_sessions.`tenant_id` = batch_session_topics.`tenant_id`"),
// 		// repository.Filter("batch_sessions.`tenant_id` = ? AND batch_sessions.`deleted_at` IS NULL", tenantID),
// 		repository.Filter("batch_session_topics.`tenant_id` = ? AND batch_session_topics.`deleted_at` IS NULL", tenantID),
// 		repository.GroupBy("batch_session_topics.`topic_id`"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	for index := range *moduleTopics {
// 		err = service.repo.GetAllInOrder(uow, &(*moduleTopics)[index].SubTopics, "batch_session_topics.`order`",
// 			repository.Join("INNER JOIN batch_session_topics ON batch_session_topics.`sub_topic_id` = module_topics.`id` "+
// 				"AND batch_session_topics.`tenant_id` = module_topics.`tenant_id`"),
// 			repository.Filter("batch_session_topics.`deleted_at` IS NULL AND batch_session_topics.`tenant_id` = ?", tenantID),
// 			repository.OrderBy("batch_session_topics.`order`"))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}

// 	return nil
// }
