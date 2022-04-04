package service

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// NextActionReportService provides method to Get next action reports.
type NextActionReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// New returns a new instance of NextActionReportService.
func New(db *gorm.DB, repository repository.Repository) *NextActionReportService {
	return &NextActionReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetTalentNextActionReports returns talent next action report details with limit and offset.
func (ser *NextActionReportService) GetTalentNextActionReports(tenantID uuid.UUID,
	reports *[]report.TalentNextActionReportDTO, form url.Values, limit, offset int, totalCount *int) error {

	// check if tenant exists
	err := ser.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(ser.DB, true)
	selectCredential := "cred.`id` AS `login_id`,CONCAT(cred.`first_name`,' ',cred.`last_name`) AS `login_name`"
	selectTalent := "talents.`first_name`,talents.`last_name`,talents.`contact`,talents.`email`"
	selectNextAction := "talent_next_actions.`id`, talent_next_actions.`talent_id`, talent_next_actions.`stipend`, talent_next_actions.`referral_count`," +
		"talent_next_actions.`from_date`, talent_next_actions.`to_date`, talent_next_actions.`target_date`, talent_next_actions.`comment`"
	nextActionType := "talent_next_action_types.`type` AS action_type"

	var queryProcessors []repository.QueryProcessor
	queryProcessors = append(queryProcessors, repository.Select(selectCredential+","+selectTalent+","+selectNextAction+","+nextActionType),
		repository.Table("talent_next_actions"),
		repository.Join("INNER JOIN credentials cred ON talent_next_actions.`created_by`=cred.`id`"),
		repository.Join("INNER JOIN talents ON talent_next_actions.`talent_id`=talents.`id`"),
		repository.Join("INNER JOIN talent_next_action_types ON talent_next_actions.`action_type_id`=talent_next_action_types.`id`"),
		repository.Filter("talent_next_actions.`deleted_at` IS NULL AND talent_next_actions.`tenant_id`=?", tenantID),
		repository.Filter("cred.`deleted_at` IS NULL AND cred.`tenant_id`=?", tenantID),
		repository.Filter("talent_next_action_types.`deleted_at` IS NULL AND talent_next_action_types.`tenant_id`=?", tenantID))

	queryProcessors = append(queryProcessors, ser.addSearchQueriesForNextAction(form)...)
	queryProcessors = append(queryProcessors, repository.Paginate(limit, offset, totalCount))

	err = ser.Repository.GetAll(uow, reports, queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	err = ser.getCourses(tenantID, reports, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = ser.getCompanyBranches(tenantID, reports, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = ser.getTechnologies(tenantID, reports, uow)
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

func (ser *NextActionReportService) getCourses(tenantID uuid.UUID, reports *[]report.TalentNextActionReportDTO,
	uow *repository.UnitOfWork) error {

	for index := range *reports {

		err := ser.Repository.GetAll(uow, &(*reports)[index].Courses,
			repository.Join("INNER JOIN talent_next_actions_courses ON talent_next_actions_courses.`course_id`=courses.`id`"),
			repository.Join("INNER JOIN talent_next_actions ON talent_next_actions_courses.`next_action_id`=talent_next_actions.`id`"),
			repository.Filter("courses.`deleted_at` IS NULL AND courses.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`deleted_at` IS NULL AND talent_next_actions.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`id`=?", (*reports)[index].ID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ser *NextActionReportService) getCompanyBranches(tenantID uuid.UUID, reports *[]report.TalentNextActionReportDTO,
	uow *repository.UnitOfWork) error {

	for index := range *reports {

		err := ser.Repository.GetAll(uow, &(*reports)[index].Companies,
			repository.Join("INNER JOIN talent_next_actions_company_branches ON talent_next_actions_company_branches.`company_branch_id`=company_branches.`id`"),
			repository.Join("INNER JOIN talent_next_actions ON talent_next_actions_company_branches.`next_action_id`=talent_next_actions.`id`"),
			repository.Filter("company_branches.`deleted_at` IS NULL AND company_branches.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`deleted_at` IS NULL AND talent_next_actions.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`id`=?", (*reports)[index].ID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ser *NextActionReportService) getTechnologies(tenantID uuid.UUID, reports *[]report.TalentNextActionReportDTO,
	uow *repository.UnitOfWork) error {

	for index := range *reports {

		err := ser.Repository.GetAll(uow, &(*reports)[index].Technologies,
			repository.Join("INNER JOIN talent_next_actions_technologies ON talent_next_actions_technologies.`technology_id`=technologies.`id`"),
			repository.Join("INNER JOIN talent_next_actions ON talent_next_actions_technologies.`next_action_id`=talent_next_actions.`id`"),
			repository.Filter("technologies.`deleted_at` IS NULL AND technologies.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`deleted_at` IS NULL AND talent_next_actions.`tenant_id`=?", tenantID),
			repository.Filter("talent_next_actions.`id`=?", (*reports)[index].ID))
		if err != nil {
			return err
		}
	}
	return nil
}

// returns error if there is no tenant record in table.
func (ser *NextActionReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (ser *NextActionReportService) addSearchQueriesForNextAction(requestForm url.Values) []repository.QueryProcessor {
	fmt.Println("=========================In addSearchQueriesForNextAction============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if actionType, ok := requestForm["actionType"]; ok {
		util.AddToSlice("talent_next_actions.`action_type_id`", "= ?", "AND", actionType, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(talent_next_actions.`from_date` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}
	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(talent_next_actions.`to_date` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}
	if targetDate, ok := requestForm["targetDate"]; ok {
		util.AddToSlice("CAST(talent_next_actions.`target_date` AS DATE)", "= ?", "AND", targetDate, &columnNames, &conditions, &operators, &values)
	}
	if referralCount, ok := requestForm["referralCount"]; ok {
		util.AddToSlice("talent_next_actions.`referral_count`", "= ?", "AND", referralCount, &columnNames, &conditions, &operators, &values)
	}
	if stipend, ok := requestForm["stipend"]; ok {
		util.AddToSlice("talent_next_actions.`stipend`", "= ?", "AND", stipend, &columnNames, &conditions, &operators, &values)
	}
	if loginID, ok := requestForm["loginID"]; ok {
		util.AddToSlice("talent_next_actions.`created_by`", "= ?", "AND", loginID, &columnNames, &conditions, &operators, &values)
	}
	//if technologies is present then join talent_next_actions_technologies and technolgies table
	if technologies, ok := requestForm["technologies"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN talent_next_actions_technologies ON talent_next_actions.`id` = talent_next_actions_technologies.`next_action_id`"))
		if len(technologies) > 0 {
			util.AddToSlice("talent_next_actions_technologies.`technology_id`", "IN(?)", "AND", technologies, &columnNames, &conditions, &operators, &values)
		}
	}
	//if technologies is present then join talent_next_actions_courses and technolgies table
	if courses, ok := requestForm["courses"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN talent_next_actions_courses ON talent_next_actions.`id` = talent_next_actions_courses.`next_action_id`"))
		if len(courses) > 0 {
			util.AddToSlice("talent_next_actions_courses.`course_id`", "IN(?)", "AND", courses, &columnNames, &conditions, &operators, &values)
		}
	}
	//if technologies is present then join talent_next_actions_company_branches and technolgies table
	if companies, ok := requestForm["companies"]; ok {
		queryProcessors = append(queryProcessors, repository.Join("INNER JOIN talent_next_actions_company_branches ON talent_next_actions.`id` = talent_next_actions_company_branches.`next_action_id`"))
		if len(companies) > 0 {
			util.AddToSlice("talent_next_actions_company_branches.`company_branch_id`", "IN(?)", "AND", companies, &columnNames, &conditions, &operators, &values)
		}
	}

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("talent_next_actions.`id`"))

	return queryProcessors
	// return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
