package service

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DashboardService provides all details to be shown on talent dashboard.
type DashboardService struct {
	db         *gorm.DB
	Repository repository.Repository
}

// NewTalentDashboardService returns new instance of TalentDashboardService.
func NewTalentDashboardService(db *gorm.DB, repository repository.Repository) *DashboardService {
	return &DashboardService{
		db:         db,
		Repository: repository,
	}
}

// GetTalentDashboardDetails gets all details required for TalentDashboard.
func (service *DashboardService) GetTalentDashboardDetails(talentDashboard *dashboard.TalentDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.db, true)

	// Get count of total talents.
	if err := service.Repository.GetCount(uow, tal.Talent{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.TotalTalents = uint(totalCount)

	if err := service.getAllB2CTalents(talentDashboard); err != nil {
		return err
	}

	talentDashboard.TotalFreshers = 786
	talentDashboard.TotalExperienced = 786
	talentDashboard.SwabhavTalents = 786
	talentDashboard.TotalInterestedInForeign = 786

	return nil

}

// Gets all fresher requirements.
func (service *DashboardService) getAllB2CTalents(talentDashboard *dashboard.TalentDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)
	b2c := dashboard.TalentSegregation{
		FirstYear:  786,
		SecondYear: 786,
		ThirdYear:  786,
		FourthYear: 786,
	}
	talentDashboard.B2CTalents = b2c.GetSumOfAllTalents()
	return nil
}

// GetFacultyFeedbackToTalentWeekWiseDashboardDetails will return minimum details for faculty feedback to be
// displayed on talent dashboard for current week and previous week.
func (service *DashboardService) GetFacultyFeedbackToTalentWeekWiseDashboardDetails(twoFeedbacks *dashboard.ThisAndPreviousWeekFacultyToTalentFeedackDashboard,
	tenantID, batchID, talentID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(talentID, tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	// Get all faculty feedbacks for talent for current week.
	err = service.Repository.GetAll(uow, &twoFeedbacks.ThisWeekFeedbacks,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select("SUM(faculty_talent_batch_session_feedback.`answer`)/COUNT(DISTINCT(batch_sessions.`id`)) AS answer,"+
			" feedback_questions.`question` as question"),
		repository.Join("JOIN batch_sessions on batch_sessions.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"),
		repository.Join("JOIN feedback_questions on faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"),
		repository.Filter("faculty_talent_batch_session_feedback.`batch_id`=?", batchID),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=?", talentID),
		repository.Filter("faculty_talent_batch_session_feedback.`option_id` IS NOT NULL"),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?", tenantID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL"),
		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL"),
		repository.Filter("YEARWEEK(batch_sessions.`date`, 1) = YEARWEEK(CURDATE(), 1)"),
		repository.GroupBy("feedback_questions.`id`"),
		repository.OrderBy("feedback_questions.`question`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get all faculty feedbacks for talent for previous week.
	err = service.Repository.GetAll(uow, &twoFeedbacks.PreviousWeekFeedbacks,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select("SUM(faculty_talent_batch_session_feedback.`answer`)/COUNT(DISTINCT(batch_sessions.`id`)) AS answer,"+
			" feedback_questions.`question` as question"),
		repository.Join("JOIN batch_sessions on batch_sessions.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"),
		repository.Join("JOIN feedback_questions on faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"),
		repository.Filter("faculty_talent_batch_session_feedback.`batch_id`=?", batchID),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=?", talentID),
		repository.Filter("faculty_talent_batch_session_feedback.`option_id` IS NOT NULL"),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?", tenantID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL"),
		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL"),
		repository.Filter("batch_sessions.`date` BETWEEN CURDATE()-INTERVAL 1 WEEK AND CURDATE()"),
		repository.OrderBy("feedback_questions.`question`"),
		repository.GroupBy("feedback_questions.`id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFacultyFeedbackToTalentDashboardDetails will return minimum details for faculty feedback to be
// displayed on talent dashboard.
func (service *DashboardService) GetFacultyFeedbackToTalentDashboardDetails(feedbacks *[]dashboard.FacultyFeedbackToTalentDashboard,
	tenantID, batchID, talentID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(talentID, tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	// Get all faculty feedback for talent for all batch sessions.
	err = service.Repository.GetAll(uow, &feedbacks,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select("SUM(faculty_talent_batch_session_feedback.`answer`)/COUNT(DISTINCT(batch_sessions.`id`)) AS answer,"+
			" feedback_questions.`question` as question"),
		repository.Join("JOIN batch_sessions on batch_sessions.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"),
		repository.Join("JOIN feedback_questions on faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"),
		repository.Filter("faculty_talent_batch_session_feedback.`batch_id`=?", batchID),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=?", talentID),
		repository.Filter("faculty_talent_batch_session_feedback.`option_id` IS NOT NULL"),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?", tenantID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL"),
		repository.Filter("batch_sessions.`deleted_at` IS NULL"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL"),
		repository.GroupBy("feedback_questions.`id`"),
		repository.OrderBy("feedback_questions.`question`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFacultyFeedbackRatingLeaderBoard will return minimum details for faculty feedback to be
// displayed on talent dashboard leader board.
func (service *DashboardService) GetFacultyFeedbackRatingLeaderBoard(feedbacks *[]dashboard.FacultyFeedbackRatingLeaderBoard,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch exist.
	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.db, true)

	// Create query precessors for sub query one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Select("AVG(`average_rating`) as rating, talents.`first_name` as first_name, talents.`last_name`"+
			" as last_name, batch_session_talents.`talent_id` as talent_id"),
		repository.Table("batch_session_talents"),
		repository.Join("INNER JOIN batch_sessions on batch_sessions.`id` = batch_session_talents.`batch_session_id`"+
			"AND batch_sessions.`tenant_id` = batch_session_talents.`tenant_id`"),
		repository.Join("JOIN talents on batch_session_talents.`talent_id` = talents.`id`"),
		service.addSearchQueriesForBatchSessions(parser.Form),
		repository.Filter("batch_session_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_session_talents.`tenant_id`=?", tenantID),
		repository.Filter("batch_session_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_session_talents.`batch_id`=?", batchID),
		repository.Filter("batch_session_talents.`attended_date` IS NOT NULL"),
		repository.Filter("batch_session_talents.`is_present`=?", true),
		repository.Filter("YEARWEEK(batch_session_talents.`attended_date`, 1) = YEARWEEK(CURDATE(), 1)"),
		repository.GroupBy("batch_session_talents.`talent_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, dashboard.FacultyFeedbackRatingLeaderBoard{},
		queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	// rawQuery := "SELECT DENSE_RANK() OVER(ORDER BY  rating desc) as `rank`, sub.* FROM ? as sub"
	rawQuery := "SELECT sub.* FROM ? as sub"

	if _, ok := parser.Form["limit"]; ok {
		rawQuery = rawQuery + " LIMIT " + parser.Form.Get("limit")
	}

	// Get leader board.
	if err := service.Repository.GetAll(uow, &feedbacks,
		repository.RawQuery(rawQuery, subQueryOne)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	uow.Commit()
	return nil
}

// GetWeeklyRating will return weekly rating for talents
func (service *DashboardService) GetWeeklyRating(performanceDetails *[]tal.PerformanceDetails,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.Repository.Scan(uow, performanceDetails, repository.Table("talents"),
		repository.Select("talents.`first_name`, talents.`last_name`, talents.`id` AS talent_id"),
		repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = talents.`id` AND"+
			" batch_talents.`tenant_id` = talents.`tenant_id`"), repository.Filter("batch_talents.`batch_id` = ?", batchID),
		service.addSearchQueriesForBatchbatchTalents(parser.Form),
		repository.Filter("talents.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	var queryProcessors []repository.QueryProcessor

	for index := range *performanceDetails {
		queryProcessors = []repository.QueryProcessor{}
		queryProcessors = append(queryProcessors, repository.Table("faculty_talent_batch_session_feedback"),
			repository.Select("((SUM(fo.`key`) / SUM(fq.`max_score`)) * 10) AS rating"),
			repository.Join("INNER JOIN batch_sessions ON faculty_talent_batch_session_feedback.`batch_session_id` = batch_sessions.`id`"+
				" AND batch_sessions.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Join("INNER JOIN feedback_questions AS fq ON fq.`id` = faculty_talent_batch_session_feedback.`question_id`"+
				" AND fq.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Join("INNER JOIN feedback_options AS fo ON fo.`id` = faculty_talent_batch_session_feedback.`option_id`"+
				" AND fo.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Filter("faculty_talent_batch_session_feedback.`tenant_id` = ? AND"+
				" faculty_talent_batch_session_feedback.`batch_id` = ? AND fq.`deleted_at` IS NULL AND fo.`deleted_at` IS NULL AND"+
				" batch_sessions.`deleted_at` IS NULL AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID, batchID),
			repository.Filter("faculty_talent_batch_session_feedback.`talent_id` = ?", (*performanceDetails)[index].TalentID))
		queryProcessors = append(queryProcessors, service.addSessionTalentSearchQueries(parser.Form)...)
		queryProcessors = append(queryProcessors, repository.GroupBy("WEEK(batch_sessions.`date`), faculty_talent_batch_session_feedback.`talent_id`"))

		err = service.Repository.Scan(uow, &(*performanceDetails)[index].Score, queryProcessors...)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetFeedbackRating will fetch feedback rating of all talents
func (service *DashboardService) GetFeedbackRating(performanceDetails *[]tal.PerformanceDetails,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors, repository.Table("talents"),
		repository.Select("talents.`first_name`, talents.`last_name`, talents.`id` AS talent_id, "+
			"((SUM(fo.`key`) / SUM(fq.`max_score`)) * 10) AS average_rating"),
		repository.Join("INNER JOIN faculty_talent_batch_session_feedback ON"+
			" faculty_talent_batch_session_feedback.`talent_id` = talents.`id`"+
			" AND talents.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN batch_sessions AS bs ON bs.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"+
			" AND bs.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_questions AS fq ON fq.`id` = faculty_talent_batch_session_feedback.`question_id`"+
			" AND fq.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options AS fo ON fo.`id` = faculty_talent_batch_session_feedback.`option_id`"+
			" AND fo.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = talents.`id` AND"+
			" batch_talents.`tenant_id` = talents.`tenant_id`"), repository.Filter("batch_talents.`batch_id` = ?", batchID),
		repository.Filter("faculty_talent_batch_session_feedback.`deleted_at` IS NULL AND bs.`deleted_at` IS NULL AND"+
			" fq.`deleted_at` IS NULL AND fo.`deleted_at` IS NULL"),
		repository.Filter("talents.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL"+
			" AND talents.`deleted_at` IS NULL", tenantID),
			service.addSearchQueriesForBatchbatchTalents(parser.Form),
		)
	queryProcessors = append(queryProcessors, service.addSessionTalentSearchQueries(parser.Form)...)
	queryProcessors = append(queryProcessors, repository.GroupBy("faculty_talent_batch_session_feedback.`talent_id`"))

	err = service.Repository.Scan(uow, performanceDetails, queryProcessors...)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTalentConceptRatingWithBatchTopicAssignment will get talent concept ratings for each batch topic assignemnt for all talents.
func (service *DashboardService) GetTalentConceptRatingWithBatchTopicAssignment(batchTalents *[]tal.ConceptRatingWithAssignment,
	tenantID, batchID uuid.UUID, parser *web.Parser) error {

	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesBatchExist(batchID, tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, true)

	err = service.Repository.Scan(uow, batchTalents, 
		repository.Table("talents"),
		repository.Select("talents.`first_name`, talents.`last_name`, talents.`id` AS talent_id"),
		repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = talents.`id` AND"+
			" batch_talents.`tenant_id` = talents.`tenant_id`"), 
		repository.Filter("batch_talents.`batch_id` = ?", batchID),
		service.addSearchQueriesForBatchbatchTalents(parser.Form),
		repository.Filter("talents.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL", tenantID))
	if err != nil {
		uow.RollBack()
		return err
	}

	var queryProcessors []repository.QueryProcessor

	for i := range *batchTalents {
		queryProcessors = []repository.QueryProcessor{}
		queryProcessors = append(queryProcessors, repository.Table("talent_concept_ratings"),
		repository.Select("talent_concept_ratings.`score`, batch_topic_assignments.`id` AS assignment_id"),
		repository.Join("INNER JOIN talent_assignment_submissions on talent_assignment_submissions.`id` = "+
			"talent_concept_ratings.`talent_submission_id`"),
		repository.Join("INNER JOIN batch_topic_assignments on batch_topic_assignments.`id` = "+
			"talent_assignment_submissions.`batch_topic_assignment_id`"),
		repository.Filter("talent_concept_ratings.`deleted_at` IS NULL AND talent_concept_ratings.`tenant_id`=?", tenantID),
		repository.Filter("talent_assignment_submissions.`deleted_at` IS NULL AND talent_assignment_submissions.`tenant_id`=?", tenantID),
		repository.Filter("batch_topic_assignments.`deleted_at` IS NULL AND batch_topic_assignments.`tenant_id`=?", tenantID),
		repository.Filter("batch_topic_assignments.`batch_id`=?", batchID),
		repository.Filter("talent_concept_ratings.`talent_id`=?", &(*batchTalents)[i].TalentID),)
		queryProcessors = append(queryProcessors, service.addSessionTalentSearchQueries(parser.Form)...)
		
		err = service.Repository.Scan(uow, &(*batchTalents)[i].Assignments, queryProcessors...)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// addSearchQueriesForBatchbatchTalents adds all search queries if any when getAll is called
func (service *DashboardService) addSearchQueriesForBatchbatchTalents(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("batch_talents.`talent_id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addSearchQueriesForBatchSessions adds all search queries if any when getAll is called
func (service *DashboardService) addSearchQueriesForBatchSessions(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("batch_sessions.`faculty_id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}


func (service *DashboardService) addSessionTalentSearchQueries(requestForm url.Values) []repository.QueryProcessor {

	fmt.Println("================================================requestForm", requestForm)

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcessors []repository.QueryProcessor

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("faculty_talent_batch_session_feedback.`faculty_id`", "= ?", "AND",
			facultyID, &columnNames, &conditions, &operators, &values)
	}

	if questionID, ok := requestForm["questionID"]; ok {
		util.AddToSlice("faculty_talent_batch_session_feedback.`question_id`", "= ?", "AND",
			questionID, &columnNames, &conditions, &operators, &values)
	}

	if moduleConceptID, ok := requestForm["moduleConceptID"]; ok {
		util.AddToSlice("talent_concept_ratings.`module_programming_concept_id`", "= ?", "AND",
			moduleConceptID, &columnNames, &conditions, &operators, &values)
	}
	
	queryProcessors = append(queryProcessors, repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcessors
}

// doesTenantExist validates if tenant exists or not in database.
func (service *DashboardService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.db, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *DashboardService) doesTalentExist(talentID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *DashboardService) doesBatchExist(batchID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.db, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
