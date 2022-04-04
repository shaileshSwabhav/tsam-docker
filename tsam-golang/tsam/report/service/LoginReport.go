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

// LoginReportService provides method to get login reports.
type LoginReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewLoginReportService returns a new instance of LoginReportService.
func NewLoginReportService(db *gorm.DB, repository repository.Repository) *LoginReportService {
	return &LoginReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetLoginReports will return details for login-report.
func (service *LoginReportService) GetLoginReports(tenantID uuid.UUID, reports *[]report.LoginReport,
	limit, offset int, totalCount *int, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	credentials := new([]general.Credential)

	uow := repository.NewUnitOfWork(service.DB, true)

	// repository.Filter("roles.`is_employee` = 1"),

	searchQP := service.addSearchQueries(requestForm)

	searchQP = append(searchQP, repository.Select("credentials.`id`, `first_name`, `last_name`, credentials.`role_id`"),
		repository.Join("INNER JOIN roles ON roles.`id` = credentials.`role_id` AND roles.`tenant_id` = credentials.`tenant_id`"),
		repository.Join("INNER JOIN login_sessions ON login_sessions.`credential_id` = credentials.`id`"+
			" AND credentials.`tenant_id` = login_sessions.`tenant_id`"),
		repository.Filter("credentials.`tenant_id` = ? AND credentials.`deleted_at` IS NULL AND login_sessions.`deleted_at` IS NULL"+
			" AND roles.`deleted_at` IS NULL", tenantID), repository.OrderBy("login_sessions.`start_time` DESC"),
		repository.GroupBy("login_sessions.`credential_id`"), repository.PreloadAssociations([]string{"Role"}),
		repository.Paginate(limit, offset, totalCount))

	err = service.Repository.GetAll(uow, credentials, searchQP...)
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range *credentials {
		tempReport := new(report.LoginReport)
		var totalLoginCount int

		err = service.Repository.GetRecordForTenant(uow, tenantID, tempReport, repository.Table("login_sessions"),
			repository.Select("login_sessions.`id` AS `login_session_id`, login_sessions.`start_time` AS `login_time`,"+
				" login_sessions.`end_time` AS `logout_time`"),
			repository.Filter("`credential_id` = ?", (*credentials)[index].ID),
			repository.OrderBy("`start_time` DESC"))
		if err != nil {
			uow.RollBack()
			return err
		}

		exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.LoginSession{},
			repository.Filter("login_sessions.`end_time` IS NOT NULL"), service.addLoginSessionQueries(requestForm),
			repository.Filter("login_sessions.`credential_id` = ?", (*credentials)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}
		if exist {
			err = service.Repository.GetAllForTenant(uow, tenantID, tempReport, repository.Table("login_sessions"),
				repository.Select("SEC_TO_TIME(SUM(TIME_TO_SEC(TIMEDIFF(login_sessions.`end_time`, login_sessions.`start_time`)))) AS total_hours"),
				repository.Filter("login_sessions.`end_time` IS NOT NULL"), service.addLoginSessionQueries(requestForm),
				repository.Filter("login_sessions.`credential_id` = ?", (*credentials)[index].ID),
				repository.OrderBy("login_sessions.`created_at` DESC"))
			if err != nil {
				uow.RollBack()
				return err
			}
		}

		if tempReport.LogoutTime != nil {
			tempReport.LastLoginTime = *tempReport.LoginTime
			tempReport.LastLogoutTime = tempReport.LogoutTime
			tempReport.LoginTime = nil
			tempReport.LogoutTime = nil
		} else {
			// Check if last login_exist
			exist, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, new(general.LoginSession),
				repository.Select("login_sessions.`start_time` AS `last_login_time`, login_sessions.`end_time` AS `last_logout_time`"),
				repository.Filter("`id` != ? AND `credential_id` = ?", tempReport.LoginSessionID, (*credentials)[index].ID))
			if err != nil {
				return err
			}
			if exist {
				// Get last login and logout time.
				err = service.Repository.GetRecordForTenant(uow, tenantID, tempReport, repository.Table("login_sessions"),
					repository.Select("login_sessions.`start_time` AS `last_login_time`, login_sessions.`end_time` AS `last_logout_time`"),
					repository.Filter("`id` != ? AND `credential_id` = ?", tempReport.LoginSessionID, (*credentials)[index].ID),
					repository.OrderBy("`start_time` DESC"))
				if err != nil {
					uow.RollBack()
					return err
				}
			}
		}

		// Get login count.
		err = service.Repository.GetCountForTenant(uow, tenantID, &general.LoginSession{},
			&totalLoginCount, repository.Filter("login_sessions.`end_time` IS NOT NULL"),
			service.addLoginSessionQueries(requestForm), repository.Filter("`credential_id` = ?", (*credentials)[index].ID))
		if err != nil {
			uow.RollBack()
			return err
		}

		tempReport.CredentialID = (*credentials)[index].ID
		tempReport.LoginName = fmt.Sprintf("%s %s", (*credentials)[index].FirstName, *(*credentials)[index].LastName)
		tempReport.RoleName = (*credentials)[index].Role.RoleName
		tempReport.LoginCount = totalLoginCount

		*reports = append(*reports, *tempReport)
	}

	uow.Commit()
	return nil

	// credentialString := "credentials.`id` AS `credential_id`, CONCAT(credentials.`first_name`,' ',credentials.`last_name`) AS `login_name`"
	// loginSessionString := "login_sessions.`id` AS `login_session_id`, login_sessions.`start_time` AS `login_time`, login_sessions.`end_time` AS `logout_time`"
	// roleString := "roles.`role_name`"

	// subQuery, err := service.Repository.SubQuery(uow, &general.LoginSession{},
	// 	repository.Table("login_sessions"), repository.Select("max(login_sessions.`start_time`)"),
	// 	repository.GroupBy("login_sessions.`credential_id`"))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }
	// // var queryProcessor []repository.QueryProcessor

	// err = service.Repository.GetAll(uow, reports, repository.Table("login_sessions"),
	// 	repository.Select(credentialString+","+loginSessionString+","+roleString),
	// 	repository.Join("INNER JOIN credentials ON login_sessions.`credential_id` = credentials.`id`"+
	// 		" AND credentials.`tenant_id` = login_sessions.`tenant_id`"),
	// 	repository.Join("INNER JOIN roles ON roles.`id` = credentials.`role_id` AND roles.`tenant_id` = credentials.`tenant_id`"),
	// 	repository.Filter("login_sessions.`tenant_id` = ? AND credentials.`deleted_at` IS NULL AND roles.`deleted_at` IS NULL",
	// 		tenantID), service.addSearchQueries(requestForm),
	// 	repository.Filter("roles.`is_employee` = 1"), repository.Filter("login_sessions.`start_time` IN ? ", subQuery),
	// 	repository.OrderBy("login_sessions.`start_time` DESC"), repository.Paginate(limit, offset, totalCount))
	// if err != nil {
	// 	uow.RollBack()
	// 	return err
	// }

}

// GetCredentialLoginReports returns details of login and logout
func (service *LoginReportService) GetCredentialLoginReports(tenantID, credentialID uuid.UUID,
	loginReports *[]report.CredentialLoginReport, limit, offset int, totalCount *int, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAll(uow, loginReports, repository.Table("login_sessions"),
		repository.Select("login_sessions.`start_time` AS login_time, login_sessions.`end_time` AS logout_time, "+
			"TIMEDIFF(login_sessions.`end_time`, login_sessions.`start_time`) AS total_hours"),
		repository.Filter("login_sessions.`tenant_id` = ?", tenantID), repository.Filter("login_sessions.`end_time` IS NOT NULL"),
		repository.Filter("login_sessions.`credential_id` = ?", credentialID), service.addLoginSessionQueries(requestForm),
		repository.OrderBy("login_sessions.`created_at` DESC"), repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil

	// credentialString := "CONCAT(credentials.`first_name`,' ',credentials.`last_name`) AS `login_name`"
	// loginSessionString := "login_sessions.`start_time` AS login_time, login_sessions.`end_time` AS logout_time, " +
	// 	"TIMEDIFF(login_sessions.`end_time`, login_sessions.`start_time`) AS total_hours"
	// roleString := "roles.`role_name`"	+","+roleString

	// searchQP, err := service.addSearchQueries(requestForm)
	// if err != nil {
	// 	return err
	// }

	// repository.Join("INNER JOIN credentials ON credentials.`id` = login_sessions.`credential_id` AND "+
	// 	"login_sessions.`tenant_id` = credentials.`tenant_id`"),
	// repository.Join("INNER JOIN roles ON roles.`id` = credentials.`role_id` AND "+
	// 	"credentials.`tenant_id` = roles.`tenant_id`"),AND roles.`deleted_at` IS NULL
	// repository.Filter("credentials.`deleted_at` IS NULL"),

}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// returns error if there is no tenant record in table.
func (service *LoginReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table.
func (service *LoginReportService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *LoginReportService) addSearchQueries(requestForm url.Values) []repository.QueryProcessor {

	// fmt.Println("=========================In addSearchQueries============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	var queryProcessors []repository.QueryProcessor

	if isActive, ok := requestForm["isActive"]; ok {
		util.AddToSlice("credentials.`is_active`", "= ?", "AND", isActive, &columnNames, &conditions, &operators, &values)
	}

	if roleID, ok := requestForm["roleID"]; ok {
		util.AddToSlice("roles.`id`", "= ?", "AND", roleID, &columnNames, &conditions, &operators, &values)
	}

	if batchID, ok := requestForm["batchID"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = credentials.`talent_id` "+
				"AND batch_talents.`tenant_id` = credentials.`tenant_id`"),
			repository.Filter("batch_talents.`batch_id` = ?", batchID))
	}

	if courseID, ok := requestForm["courseID"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = credentials.`talent_id` "+
				"AND batch_talents.`tenant_id` = credentials.`tenant_id`"),
			repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id` AND batches.`tenant_id` = batch_talents.`tenant_id`"),
			repository.Filter("batches.`course_id` = ?", courseID))
	}

	// if fromDate, ok := requestForm["fromDate"]; ok {
	// 	util.AddToSlice("login_sessions.`start_time`", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	// }

	// if toDate, ok := requestForm["toDate"]; ok {
	// 	util.AddToSlice("login_sessions.`start_time`", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	// }

	// if durationString, ok := requestForm["duration"]; ok {
	// 	duration, err := strconv.Atoi(durationString[0])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	startDate := time.Now()
	// 	endDate := startDate.AddDate(0, 0, -duration)
	// 	util.AddToSlice("CAST(login_sessions.`start_time` AS DATE)", "<= ?", "AND",
	// 		startDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
	// 	util.AddToSlice("CAST(login_sessions.`end_time` AS DATE)", ">= ?", "AND",
	// 		endDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
	// }

	queryProcessors = append(queryProcessors,
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("credentials.`id`")) // <- what should be used for group_by

	// return repository.FilterWithOperator(columnNames, conditions, operators, values)
	return queryProcessors
}

func (service *LoginReportService) addLoginSessionQueries(requestForm url.Values) repository.QueryProcessor {

	// fmt.Println("=========================In addLoginSessionQueries============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	// var queryProcessors []repository.QueryProcessor

	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(login_sessions.`start_time` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(login_sessions.`start_time` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
