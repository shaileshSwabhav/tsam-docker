package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DashboardService provides all details to be shown on faculty dashboard.
type DashboardService struct {
	db          *gorm.DB
	repo        repository.Repository
	batchStatus []string
}

// NewFacultyDashboardService returns new instance of FacultyDashboardService.
func NewFacultyDashboardService(db *gorm.DB, repository repository.Repository) *DashboardService {
	return &DashboardService{
		db:          db,
		repo:        repository,
		batchStatus: []string{"Ongoing", "Upcoming", "Finished"},
	}
}

// GetFacultyDashboardDetails gets all details required for FacultyDashboard.
func (service *DashboardService) GetFacultyDashboardDetails(facultyDashboard *dashboard.FacultyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.db, true)

	// Get count of total faculty.
	if err := service.repo.GetCount(uow, faculty.Faculty{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.TotalFaculty = uint(totalCount)

	if err := service.getAllBatchesCompleted(facultyDashboard); err != nil {
		return err
	}

	if err := service.getAllLiveBatches(facultyDashboard); err != nil {
		return err
	}

	if err := service.getAllUpcomingBatches(facultyDashboard); err != nil {
		return err
	}

	return nil

}

// GetFacultyBatchDetails will return total count of all batches for specified faculty.
func (service *DashboardService) GetFacultyBatchDetails(tenantID, facultyID uuid.UUID,
	facultyBatch *faculty.FacultyBatch, requestForm url.Values) error {

	err := service.getFacultyBatchDetails(tenantID, facultyID, facultyBatch, requestForm)
	if err != nil {
		return err
	}
	return nil
}

// GetOngoingBatchDetails will return ongoing batch details for specified faculty.
func (service *DashboardService) GetOngoingBatchDetails(tenantID, facultyID uuid.UUID,
	batchDetails *[]faculty.OngoingBatchDetails, requestForm url.Values) error {

	err := service.getOngoingBatchDetails(tenantID, facultyID, batchDetails, requestForm)
	if err != nil {
		return err
	}
	return nil
}

// GetBarchartData will return data for barchart.
func (service *DashboardService) GetBarchartData(tenantID, facultyID uuid.UUID,
	barchartData *faculty.BarchartData, requestForm url.Values) error {

	err := service.getBarchartData(tenantID, facultyID, barchartData, requestForm)
	if err != nil {
		return err
	}
	return nil
}

// GetPiechartData will return project and its count in a week.
func (service *DashboardService) GetPiechartData(tenantID, facultyID uuid.UUID,
	piechartData *[]faculty.PiechartData, requestForm url.Values) error {

	err := service.getPiechartData(tenantID, facultyID, piechartData, requestForm)
	if err != nil {
		return nil
	}
	return nil
}

// GetTaskList will get activites from timesheet for specified credentialID.
func (service *DashboardService) GetTaskList(tenantID, credentialID uuid.UUID,
	activities *[]admin.TimesheetActivity, requestForm url.Values) error {

	err := service.getTaskList(tenantID, credentialID, activities, requestForm)
	if err != nil {
		return err
	}

	return nil
}

func (service *DashboardService) GetTalentFeedbackScore(talentFeedback *faculty.Feedback,
	tenantID, batchID uuid.UUID, requestForm url.Values) error {

	err := service.getTalentFeedbackScore(talentFeedback, tenantID, batchID, requestForm)
	if err != nil {
		return err
	}

	return nil
}

// GetWeeklyAvgRating will return average rating of faculty for specified batch.
func (service *DashboardService) GetWeeklyAvgRating(rating *faculty.WeeklyAvgRating,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.getWeeklyAvgRating(rating, tenantID, batchID, parser)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionFeedbackRating will return session wise feedback rating.
func (service *DashboardService) GetSessionFeedbackRating(rating *[]faculty.WeeklyAvgRating,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.getSessionFeedbackRating(rating, tenantID, batchID, parser)
	if err != nil {
		return err
	}

	return nil
}

// Gets all fresher requirements.
func (service *DashboardService) getAllLiveBatches(facultyDashboard *dashboard.FacultyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	facultyDashboard.LiveBatches.TotalBatches = 786
	facultyDashboard.LiveBatches.TotalTalentsJoined = 786
	return nil
}

// Gets all live requirements.
func (service *DashboardService) getAllUpcomingBatches(facultyDashboard *dashboard.FacultyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	facultyDashboard.UpcomingBatches.TotalBatches = 786
	facultyDashboard.UpcomingBatches.TotalTalentsJoined = 786
	return nil
}

// Gets all live requirements.
func (service *DashboardService) getAllBatchesCompleted(facultyDashboard *dashboard.FacultyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	facultyDashboard.BatchesCompleted.TotalBatches = 786
	facultyDashboard.BatchesCompleted.TotalTalentsJoined = 786
	return nil
}

// getFacultyBatchDetails will return total count of all batches for specified faculty.
func (service *DashboardService) getFacultyBatchDetails(tenantID, facultyID uuid.UUID,
	facultyBatch *faculty.FacultyBatch, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	// Start Transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	err = service.getOngoingBatches(uow, tenantID, facultyID, &facultyBatch.OngoingBatches)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getUpcomingBatches(uow, tenantID, facultyID, &facultyBatch.UpcomingBatches)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getFinishedBatches(uow, tenantID, facultyID, &facultyBatch.FinishedBatches)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getTotalStudents(uow, tenantID, facultyID, &facultyBatch.TotalStudents)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getCompletedTrainingHours(uow, tenantID, facultyID, facultyBatch)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// getOngoingBatches will return count of all the ongoing batches for specified faculty.
func (service *DashboardService) getOngoingBatches(uow *repository.UnitOfWork,
	tenantID, facultyID uuid.UUID, totalCount *int) error {

	err := service.repo.GetCount(uow, new(batch.Batch), totalCount,
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = ? AND batches.`tenant_id` = ?"+
			" AND batch_modules.`deleted_at` IS NULL", facultyID, service.batchStatus[0], tenantID),
		repository.GroupBy("batches.`id`"))
	if err != nil {
		return err
	}
	return nil
}

// getUpcomingBatches will return count of all the upcoming batches for specified faculty.
func (service *DashboardService) getUpcomingBatches(uow *repository.UnitOfWork,
	tenantID, facultyID uuid.UUID, totalCount *int) error {

	err := service.repo.GetCount(uow, new(batch.Batch), totalCount,
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = ? AND batches.`tenant_id` = ?"+
			" AND batch_modules.`deleted_at` IS NULL", facultyID, service.batchStatus[1], tenantID),
		repository.GroupBy("batches.`id`"))
	if err != nil {
		return err
	}
	return nil
}

// getFinishedBatches will return count of all the finished batches for specified faculty.
func (service *DashboardService) getFinishedBatches(uow *repository.UnitOfWork,
	tenantID, facultyID uuid.UUID, totalCount *int) error {

	err := service.repo.GetCount(uow, new(batch.Batch), totalCount,
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = ? AND batches.`tenant_id` = ?"+
			" AND batch_modules.`deleted_at` IS NULL", facultyID, service.batchStatus[2], tenantID),
		repository.GroupBy("batches.`id`"))
	if err != nil {
		return err
	}
	return nil
}

// getTotalStudents will calculate totalStudents for specified faculty.
func (service *DashboardService) getTotalStudents(uow *repository.UnitOfWork,
	tenantID, facultyID uuid.UUID, totalCount *int, queryProcessors ...repository.QueryProcessor) error {

	queryProcessors = append(queryProcessors, repository.Join("INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"+
		" AND batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"), repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL AND batch_modules.`deleted_at` IS NULL"),
		repository.Filter("batch_modules.`faculty_id` = ?", facultyID), repository.GroupBy("batch_talents.`talent_id`"))

	err := service.repo.GetCount(uow, &batch.Batch{}, totalCount, queryProcessors...)
	if err != nil {
		return err
	}
	return nil
}

// getCompletedTrainingHours will calculate total training hours of faculty.
func (service *DashboardService) getCompletedTrainingHours(uow *repository.UnitOfWork,
	tenantID, facultyID uuid.UUID, facultyBatch *faculty.FacultyBatch) error {

	// today := time.Now().Format("02-01-2006")

	// fmt.Println(" ================== getCompletedTrainingHours ================== ")

	exist, err := repository.DoesRecordExist(uow.DB, facultyBatch, repository.Table("batches"),
		repository.Join("INNER JOIN batch_session_topics ON batches.`id` = batch_session_topics.`batch_id` "+
			"AND batches.`tenant_id` = batch_session_topics.`tenant_id`"),
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		repository.Filter("CAST(batch_session_topics.`completed_date` AS DATE) >= ? AND CAST(batch_session_topics.`completed_date` AS DATE) < ?",
			util.GetBeginningOfMonth(time.Now()).Format("2006-01-02"), util.GetEndOfMonth(time.Now()).Format("2006-01-02")),
		repository.Filter("batch_modules.`deleted_at` IS NULL"),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`tenant_id` = ? AND batch_session_topics.`is_completed` = ?", facultyID, tenantID, 1),
		repository.Filter("batches.`deleted_at` IS NULL AND batch_session_topics.`deleted_at` IS NULL"),
		repository.GroupBy("batch_modules.`batch_id`"))
	if err != nil {
		return err
	}

	fmt.Println(" ==================================== exist ->", exist)

	if exist {
		err = service.repo.Scan(uow, facultyBatch, repository.Table("batches"),
			repository.Select("SUM(batch_session_topics.`total_time`) AS `completed_training_hrs`"),
			repository.Join("INNER JOIN batch_session_topics ON batches.`id` = batch_session_topics.`batch_id` "+
				"AND batches.`tenant_id` = batch_session_topics.`tenant_id`"),
			repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
				" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
			repository.Filter("CAST(batch_session_topics.`completed_date` AS DATE) >= ? AND  CAST(batch_session_topics.`completed_date` AS DATE) < ?",
				util.GetBeginningOfMonth(time.Now()).Format("2006-01-02"), util.GetEndOfMonth(time.Now()).Format("2006-01-02")),
			repository.Filter("batch_modules.`faculty_id` = ? AND batches.`tenant_id` = ? AND batch_session_topics.`is_completed` = ?",
				facultyID, tenantID, 1),
			repository.Filter("batches.`deleted_at` IS NULL AND batch_session_topics.`deleted_at` IS NULL"),
			repository.Filter("batch_modules.`deleted_at` IS NULL"), repository.GroupBy("batch_modules.`batch_id`"))
		if err != nil {
			return err
		}
	}

	return nil
}

// getOngoingBatchDetails will return ongoing batch details for specified faculty.
func (service *DashboardService) getOngoingBatchDetails(tenantID, facultyID uuid.UUID,
	batchDetails *[]faculty.OngoingBatchDetails, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	tempBatches := []batch.Batch{}

	uow := repository.NewUnitOfWork(service.db, true)

	// repository.Select("batches.`batch_name`, courses.`name` AS `course_name`"),
	// 	repository.Join("INNER JOIN courses ON courses.`id` = batches.`course_id` AND courses.`tenant_id` = batches.`tenant_id`"),

	err = service.repo.GetAll(uow, &tempBatches,
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"), repository.Filter("`batches`.`tenant_id` = ?", tenantID),
		repository.Filter("`batches`.`status` = ? AND `batch_modules`.`faculty_id` = ?", service.batchStatus[0], facultyID),
		repository.Filter("`batch_modules`.`deleted_at` IS NULL"), repository.PreloadAssociations([]string{"Course"}),
		repository.GroupBy("batch_modules.`batch_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, tempBatch := range tempBatches {
		batch := faculty.OngoingBatchDetails{}

		batch.BatchID = tempBatch.ID
		batch.BatchName = tempBatch.BatchName
		batch.CourseName = tempBatch.Course.Name

		var totalCount int

		// get total students
		err = service.getTotalStudents(uow, tenantID, facultyID, &totalCount,
			repository.Filter("batches.`id` = ?", tempBatch.ID))
		if err != nil {
			uow.RollBack()
			return err
		}
		batch.TotalStudents = uint(totalCount)

		// get total and pending sessions.
		err = service.getBatchSessionCount(uow, &batch.TotalSession, &batch.PendingSession,
			tenantID, tempBatch.ID)
		if err != nil {
			uow.RollBack()
			return err
		}

		*batchDetails = append(*batchDetails, batch)
	}

	return nil
}

// getBatchSessionCount will return count of total sessions and completed sessions for the specified batch
func (service *DashboardService) getBatchSessionCount(uow *repository.UnitOfWork,
	totalSession, completedSession *uint, tenantID, batchID uuid.UUID) error {

	err := service.repo.GetCountForTenant(uow, tenantID, batch.Session{}, totalSession,
		repository.Filter("batch_sessions.`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	err = service.repo.GetCountForTenant(uow, tenantID, batch.Session{}, completedSession,
		repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`is_session_taken` = ?", batchID, true))
	if err != nil {
		return err
	}

	return nil
}

// getPiechartData will return project and its count in a week.
func (service *DashboardService) getPiechartData(tenantID, facultyID uuid.UUID,
	piechartData *[]faculty.PiechartData, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAll(uow, piechartData, repository.Table("timesheets"),
		repository.Select("projects.`name` AS `project_name`, COUNT(*) AS `total_count`,"+
			" SUM(timesheet_activities.`hours_needed`) AS `hours`"),
		repository.Join("INNER JOIN timesheet_activities ON timesheet_activities.`timesheet_id` = timesheets.`id`"+
			" AND timesheet_activities.`tenant_id` = timesheets.`tenant_id`"),
		repository.Join("INNER JOIN projects ON timesheet_activities.`project_id` = projects.`id`"+
			" AND timesheet_activities.`tenant_id` = projects.`tenant_id`"),
		repository.Join("INNER JOIN credentials ON timesheets.`credential_id` = credentials.`id`"+
			" AND timesheets.`tenant_id` = credentials.`tenant_id`"), service.addSearchQueries(requestForm),
		repository.Filter("credentials.`faculty_id` = ? AND timesheets.`tenant_id` = ?", facultyID, tenantID),
		repository.Filter("timesheet_activities.`deleted_at` IS NULL AND projects.`deleted_at` IS NULL AND "+
			"credentials.`deleted_at` IS NULL"),
		repository.GroupBy("timesheet_activities.`project_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// getBarchartData will return data for barchart.
func (service *DashboardService) getBarchartData(tenantID, facultyID uuid.UUID,
	barchartData *faculty.BarchartData, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	// Start transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	// total studetns.
	err = service.repo.GetCount(uow, new(batch.MappedTalent), &barchartData.TotalStudents,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id` AND"+
			" batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		service.addSearchQueries(requestForm), repository.Filter("batch_talents.`tenant_id` = ? ", tenantID),
		repository.Filter("batch_modules.`deleted_at` IS NULL AND batch_modules.`faculty_id` = ?", facultyID),
		repository.Filter("batches.`deleted_at` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// get fresher count.
	err = service.repo.GetCount(uow, new(batch.MappedTalent), &barchartData.Fresher,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id` AND"+
			" batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN talents ON talents.`id` = batch_talents.`talent_id` AND"+
			" talents.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		service.addSearchQueries(requestForm), repository.Filter("talents.`is_experience` = ?", false),
		repository.Filter("batch_talents.`tenant_id` = ? ", tenantID), repository.Filter("batch_modules.`faculty_id` = ?", facultyID),
		repository.Filter("talents.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL AND batch_modules.`deleted_at` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// get professional count.
	err = service.repo.GetCount(uow, &batch.MappedTalent{}, &barchartData.Professional,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id` AND"+
			" batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN talents ON talents.`id` = batch_talents.`talent_id` AND"+
			" talents.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id` AND"+
			" `batch_modules`.`tenant_id` = `batches`.`tenant_id`"),
		repository.Filter("talents.`is_experience` = ?", true), repository.Filter("batch_modules.`faculty_id` = ?", facultyID),
		service.addSearchQueries(requestForm), repository.Filter("batch_talents.`tenant_id` = ? ", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL AND batch_modules.`deleted_at` IS NULL"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

func (service *DashboardService) getTaskList(tenantID, credentialID uuid.UUID,
	activities *[]admin.TimesheetActivity, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAll(uow, activities,
		repository.Join("INNER JOIN timesheets ON timesheets.`id` = timesheet_activities.`timesheet_id` AND"+
			" timesheets.`tenant_id` = timesheet_activities.`tenant_id`"), service.addSearchQueries(requestForm),
		repository.Filter("timesheet_activities.`tenant_id` = ? AND timesheets.`credential_id` = ?", tenantID, credentialID),
		repository.Filter("timesheet_activities.`deleted_at` IS NULL AND timesheets.`deleted_at` IS NULL"),
		repository.OrderBy("`timesheets`.`date`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

func (service *DashboardService) getTalentFeedbackScore(feedback *faculty.Feedback,
	tenantID, batchID uuid.UUID, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.GetAllInOrder(uow, &feedback.Keywords, "feedback_questions.`order`",
		repository.Table("feedback_questions"), repository.Select("feedback_questions.`id`, `keyword` AS `name`"),
		repository.Join("INNER JOIN faculty_talent_batch_session_feedback AS f ON f.`question_id` = feedback_questions.`id` AND "+
			" f.`tenant_id` = feedback_questions.`tenant_id`"),
		repository.Join("INNER JOIN batches ON batches.`id` = f.`batch_id` AND f.`tenant_id` = batches.`tenant_id`"),
		repository.Filter("f.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL"),
		repository.Filter("`type` = ? AND feedback_questions.`tenant_id` = ? AND batches.`id` = ?",
			"Faculty_Session_Feedback", tenantID, batchID), repository.GroupBy("feedback_questions.`id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	selectTalent := "talents.`first_name`, talents.`last_name`, talents.`personality_type`, talents.`talent_type`"
	selectBatch := "batches.`batch_name`, batches.`id` AS batch_id"
	selectScore := "talents.`id` AS talent_id, ((SUM(feedback_options.`key`) / SUM(q.`max_score`)) * 10) AS score"

	err = service.repo.GetAll(uow, &feedback.TalentFeedback, repository.Table("talents"),
		repository.Select(selectTalent+","+selectBatch+","+selectScore),
		repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = talents.`id` AND"+
			" batch_talents.`tenant_id` = talents.`tenant_id`"),
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id` AND"+
			" batch_talents.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN faculty_talent_batch_session_feedback AS f ON f.`batch_id` = batches.`id` AND "+
			"f.`talent_id` = talents.`id` AND f.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN feedback_questions q ON q.`id` = f.`question_id` AND"+
			" q.`tenant_id` = f.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options ON f.`option_id` = feedback_options.`id` AND"+
			" feedback_options.`tenant_id` = f.`tenant_id`"),
		repository.Filter("f.`deleted_at` IS NULL AND batch_talents.`deleted_at` IS NULL AND batches.`deleted_at` IS NULL"+
			" AND q.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id`= ? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("batches.`id`= ? AND talents.`deleted_at` IS NULL AND batch_talents.`deleted_at` IS NULL", batchID),
		repository.GroupBy("talents.`id`"), repository.OrderBy("talents.`first_name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range feedback.TalentFeedback {
		err = service.repo.GetAll(uow, &feedback.TalentFeedback[index].SessionFeedback,
			repository.Table("faculty_talent_batch_session_feedback"),
			repository.Select("q.`keyword`, feedback_options.`key` AS `keyword_score`"),
			repository.Join("INNER JOIN feedback_questions AS q ON q.`id` = faculty_talent_batch_session_feedback.`question_id` AND"+
				" q.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id` AND"+
				" feedback_options.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Filter("faculty_talent_batch_session_feedback.`talent_id` = ? AND"+
				" faculty_talent_batch_session_feedback.`batch_id` = ?", feedback.TalentFeedback[index].TalentID, batchID),
			repository.Filter("q.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"+
				" AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL"), repository.OrderBy("q.`order`"))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// getWeeklyAvgRating will return average rating of faculty for specified batch.
func (service *DashboardService) getWeeklyAvgRating(rating *faculty.WeeklyAvgRating,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.Scan(uow, rating, repository.Table("talent_batch_session_feedback"),
		repository.Select("((SUM(o.`key`) / SUM(q.`max_score`)) * 10) AS rating"),
		repository.Join("INNER JOIN batch_sessions AS b ON b.`id` = talent_batch_session_feedback.`batch_session_id`"+
			" AND b.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options AS o ON o.`id` = talent_batch_session_feedback.`option_id`"+
			" AND o.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_questions q ON q.`id` = talent_batch_session_feedback.`question_id` AND"+
			" q.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("talent_batch_session_feedback.`batch_id` = ? AND talent_batch_session_feedback.`tenant_id` = ?"+
			"  AND CAST(b.`date` AS DATE) >= ? AND CAST(b.`date` AS DATE) <= ? AND"+
			" talent_batch_session_feedback.`deleted_at` IS NULL AND b.`deleted_at` IS NULL AND o.`deleted_at` IS NULL"+
			" AND q.`deleted_at` IS NULL", batchID, tenantID, util.GetBeginningOfWeek(time.Now()).Format("2006-01-02"),
			util.GetEndOfWeek(time.Now()).Format("2006-01-02")), service.addSearchQueries(parser.Form),
		repository.GroupBy("talent_batch_session_feedback.`batch_id`"))
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil
		}
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

func (service *DashboardService) getSessionFeedbackRating(weeklyRating *[]faculty.WeeklyAvgRating,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.repo.Scan(uow, weeklyRating, repository.Table("talent_batch_session_feedback"),
		repository.Select("((SUM(o.`key`) / SUM(q.`max_score`)) * 10) AS rating"),
		repository.Join("INNER JOIN batch_sessions AS b ON b.`id` = talent_batch_session_feedback.`batch_session_id`"+
			" AND b.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options AS o ON o.`id` = talent_batch_session_feedback.`option_id`"+
			" AND o.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_questions q ON q.`id` = talent_batch_session_feedback.`question_id` AND"+
			" q.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("q.`deleted_at` IS NULL AND b.`deleted_at` IS NULL AND o.`deleted_at` IS NULL"),
		repository.Filter("talent_batch_session_feedback.`batch_id` = ? AND talent_batch_session_feedback.`deleted_at` IS NULL AND"+
			" talent_batch_session_feedback.`tenant_id` = ?", batchID, tenantID), service.addSearchQueries(parser.Form),
		repository.GroupBy("WEEK(b.`date`)"))
	if err != nil {
		// send http error with status code 204 - No Content #shailesh
		if gorm.ErrRecordNotFound == err {
			return nil
		}
		uow.RollBack()
		return err
	}

	if !util.IsEmpty(parser.Form.Get("isDetailedRating")) && parser.Form.Get("isDetailedRating") == "1" {
		err = service.getDetailedRating(uow, weeklyRating, tenantID, batchID, parser)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

func (service *DashboardService) getDetailedRating(uow *repository.UnitOfWork, weeklyRating *[]faculty.WeeklyAvgRating,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	var feedbackQuestions []general.FeedbackQuestion

	err := service.repo.GetAllInOrder(uow, &feedbackQuestions, "feedback_questions.`order`",
		repository.Select("feedback_questions.`id`, feedback_questions.`keyword`"),
		repository.Join("INNER JOIN talent_batch_session_feedback AS t ON feedback_questions.`id` = t.`question_id`"+
			" AND t.`tenant_id` = t.`tenant_id`"),
		repository.Filter("t.`batch_id` = ? AND t.`tenant_id` = ? AND t.`deleted_at` IS NULL", batchID, tenantID),
		repository.GroupBy("t.`question_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *weeklyRating {
		var feedback faculty.FeedbackScore

		for j := range feedbackQuestions {
			err = service.repo.Scan(uow, &feedback, repository.Table("talent_batch_session_feedback"),
				repository.Select("q.`keyword`, ((SUM(o.`key`) / SUM(q.`max_score`)) * 10) AS keyword_score"),
				repository.Join("INNER JOIN batch_sessions AS b ON b.`id` = talent_batch_session_feedback.`batch_session_id`"+
					" AND b.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
				repository.Join("INNER JOIN feedback_options AS o ON o.`id` = talent_batch_session_feedback.`option_id`"+
					" AND o.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
				repository.Join("INNER JOIN feedback_questions q ON q.`id` = talent_batch_session_feedback.`question_id` AND"+
					" q.`tenant_id` = talent_batch_session_feedback.`tenant_id`"),
				repository.Filter("q.`deleted_at` IS NULL AND b.`deleted_at` IS NULL AND o.`deleted_at` IS NULL"+
					" AND talent_batch_session_feedback.`question_id` = ?", feedbackQuestions[j].ID),
				repository.Filter("talent_batch_session_feedback.`batch_id` = ? AND talent_batch_session_feedback.`deleted_at` IS NULL"+
					" AND talent_batch_session_feedback.`tenant_id` = ?", batchID, tenantID), service.addSearchQueries(parser.Form),
				repository.GroupBy("WEEK(b.`date`)"))
			if err != nil {
				if gorm.ErrRecordNotFound == err {
					return nil
				}
				uow.RollBack()
				return err
			}

			(*weeklyRating)[index].Feedback = append((*weeklyRating)[index].Feedback, feedback)
		}
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *DashboardService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.db, new(general.Tenant),
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no credential record in table.
func (service *DashboardService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, new(general.Credential),
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no faculty record in table.
func (service *DashboardService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, new(faculty.Faculty),
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		return err
	}
	return nil
}

// returns error if there is no batch record in table.
func (service *DashboardService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		return err
	}
	return nil
}

func (service *DashboardService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	fmt.Println("================================================requestForm", requestForm)

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("batches.`id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("talent_batch_session_feedback.`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	if questionID, ok := requestForm["questionID"]; ok {
		util.AddToSlice("talent_batch_session_feedback.`question_id`", "= ?", "AND", questionID, &columnNames, &conditions, &operators, &values)
	}

	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(`timesheets`.`date` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(`timesheets`.`date` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	if nextEstimatedDate, ok := requestForm["nextEstimatedDate"]; ok {
		util.AddToSlice("(`timesheet_activities`.`is_completed`", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		util.AddToSlice("`timesheet_activities`.`next_estimated_date`", "> ?)", "OR", nextEstimatedDate, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// func (service *FacultyDashboardService) beginningOfMonth(date time.Time) time.Time {
// 	return date.AddDate(0, 0, -date.Day()+1)
// }

// func (service *FacultyDashboardService) endOfMonth(date time.Time) time.Time {
// 	return date.AddDate(0, 1, -date.Day())
// }
