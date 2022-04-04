package callreportsvc

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/callreport"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CallingReportService provides method to Get calling reports.
type CallingReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// New returns a new instance of CallingReportService.
func New(db *gorm.DB, repository repository.Repository) *CallingReportService {
	return &CallingReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetLoginwiseTalentCallingReports returns loginwise calling report details related to talent.
func (service *CallingReportService) GetLoginwiseTalentCallingReports(reports *[]callreport.LoginwiseCallingReport,
	tenantID uuid.UUID, form url.Values) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	credentials := []general.Credential{}
	err = service.Repository.GetAll(uow, &credentials,
		repository.Join("INNER JOIN roles on credentials.`role_id` = roles.`id`"),
		// Remove select if full credential is needed.
		repository.Select([]string{"credentials.`id`", "credentials.`first_name`", "credentials.`last_name`"}),
		// Remove admin if not needed.
		repository.Filter("(roles.`role_name`=? OR roles.`role_name`=?) AND credentials.`tenant_id`=? AND roles.`tenant_id`=?",
			"salesperson", "admin", tenantID, tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	for _, credential := range credentials {
		var totalTalentCount uint = 0
		var totalCallingCount uint = 0
		report := callreport.LoginwiseCallingReport{}
		searchQP, err := service.addSearchQueriesForLoginwise(form)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInvalidFields, http.StatusBadRequest)
		}

		report.CredentialID = credential.ID
		report.FirstName = credential.FirstName
		if credential.LastName != nil {
			report.LastName = *credential.LastName
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talent.CallRecord{}, &totalCallingCount,
			repository.Filter("`created_by`=?", credential.ID), searchQP)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talent.CallRecord{}, &totalTalentCount,
			repository.Filter("`created_by`=?", credential.ID), searchQP,
			repository.GroupBy("talent_id"))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}
		report.TotalCallingCount = totalCallingCount
		report.TotalTalentCount = totalTalentCount
		*reports = append(*reports, report)
	}
	return nil
}

// GetDaywiseTalentCallingReports returns daywise calling report details related to talent.
func (service *CallingReportService) GetDaywiseTalentCallingReports(reports *[]callreport.DaywiseCallingReport,
	tenantID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)

	today := time.Now()
	noOfDays := 10
	for i := 0; i < noOfDays; i++ {
		var totalCallingCount uint = 0
		var totalTalentCount uint = 0
		report := callreport.DaywiseCallingReport{}
		report.Date = today.AddDate(0, 0, -i).Format("2006-01-02")
		err = service.Repository.GetCountForTenant(uow, tenantID, &talent.CallRecord{}, &totalCallingCount,
			repository.Filter("CAST(`date_time` AS DATE)=?", report.Date))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talent.CallRecord{}, &totalTalentCount,
			repository.Filter("CAST(`date_time` AS DATE)=?", report.Date),
			repository.GroupBy("talent_id"))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		report.TotalCallingCount = totalCallingCount
		report.TotalTalentCount = totalTalentCount
		*reports = append(*reports, report)
	}
	fmt.Println("==========daywise reports are============= ", reports)
	return nil
}

// GetTalentCallingReports returns talent calling report details with limit and offset.
func (service *CallingReportService) GetTalentCallingReports(tenantID uuid.UUID,
	reports *[]callreport.TalentCallingReportDTO, form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	// get all colleges ordered by name order
	selectCredential := "cred.`id` AS `login_id`,CONCAT(cred.`first_name`,' ',cred.`last_name`) AS `login_name`"
	selectTalent := "talents.`first_name`,talents.`last_name`,talents.`contact`,talents.`email`"
	selectCallRecord := "talent_call_records.`date_time`,talent_call_records.`comment`,talent_call_records.`expected_ctc`,talent_call_records.`notice_period`,talent_call_records.`target_date`"
	selectPurposeOutcome := "purposes.`purpose`,outcomes.`outcome`"
	err = service.Repository.GetAll(uow, reports,
		repository.Select(selectCredential+","+selectTalent+","+selectCallRecord+","+selectPurposeOutcome),
		repository.Table("talent_call_records"),
		repository.Join("INNER JOIN credentials cred ON talent_call_records.`created_by`=cred.`id`"),
		repository.Join("INNER JOIN talents ON talent_call_records.`talent_id`=talents.`id`"),
		repository.Join("INNER JOIN purposes ON talent_call_records.`purpose_id`=purposes.`id`"),
		repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id`=outcomes.`id`"),
		repository.Filter("talent_call_records.`deleted_at` IS NULL AND talent_call_records.`tenant_id`=?", tenantID),
		repository.Filter("cred.`deleted_at` IS NULL AND cred.`tenant_id`=?", tenantID),
		repository.Filter("purposes.`deleted_at` IS NULL AND purposes.`tenant_id`=?", tenantID),
		repository.Filter("outcomes.`deleted_at` IS NULL AND outcomes.`tenant_id`=?", tenantID),
		service.addSearchQueriesForCallingReport(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	fmt.Println("=========miraclee===============", reports)
	return nil
}

// =======================================================TALENT-ENQUIRY=======================================================

// GetLoginwiseTalentEnquiryCallingReports returns loginwise calling report details related to talent-enquiry.
func (service *CallingReportService) GetLoginwiseTalentEnquiryCallingReports(reports *[]callreport.LoginwiseTalentEnquiryCallingReport,
	tenantID uuid.UUID, form url.Values) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	credentials := []general.Credential{}
	err = service.Repository.GetAll(uow, &credentials,
		repository.Join("INNER JOIN roles on credentials.`role_id` = roles.`id`"),
		// Remove select if full credential is needed.
		repository.Select([]string{"credentials.`id`", "credentials.`first_name`", "credentials.`last_name`"}),
		// Remove admin if not needed.
		repository.Filter("(roles.`role_name`=? OR roles.`role_name`=?) AND credentials.`tenant_id`=? AND roles.`tenant_id`=?",
			"salesperson", "admin", tenantID, tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	for _, credential := range credentials {
		var totalEnquiryCount uint = 0
		var totalCallingCount uint = 0
		report := callreport.LoginwiseTalentEnquiryCallingReport{}
		searchQP, err := service.addSearchQueriesForEnquiryLoginwise(form)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInvalidFields, http.StatusBadRequest)
		}

		report.CredentialID = credential.ID
		report.FirstName = credential.FirstName
		if credential.LastName != nil {
			report.LastName = *credential.LastName
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talentenquiry.CallRecord{}, &totalCallingCount,
			repository.Filter("`created_by`=?", credential.ID), searchQP)
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talentenquiry.CallRecord{}, &totalEnquiryCount,
			repository.Filter("`created_by`=?", credential.ID), searchQP,
			repository.GroupBy("enquiry_id"))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}
		report.TotalCallingCount = totalCallingCount
		report.TotalEnquiryCount = totalEnquiryCount
		*reports = append(*reports, report)
	}
	return nil
}

// GetDaywiseTalentEnquiryCallingReports returns daywise calling report details related to talent-enquiry.
func (service *CallingReportService) GetDaywiseTalentEnquiryCallingReports(reports *[]callreport.DaywiseTalentEnquiryCallingReport,
	tenantID uuid.UUID) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)

	today := time.Now()
	noOfDays := 10
	for i := 0; i < noOfDays; i++ {
		var totalCallingCount uint = 0
		var totalEnquiryCount uint = 0
		report := callreport.DaywiseTalentEnquiryCallingReport{}
		report.Date = today.AddDate(0, 0, -i).Format("2006-01-02")
		err = service.Repository.GetCountForTenant(uow, tenantID, &talentenquiry.CallRecord{}, &totalCallingCount,
			repository.Filter("CAST(`date_time` AS DATE)=?", report.Date))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		err = service.Repository.GetCountForTenant(uow, tenantID, &talentenquiry.CallRecord{}, &totalEnquiryCount,
			repository.Filter("CAST(`date_time` AS DATE)=?", report.Date),
			repository.GroupBy("enquiry_id"))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
		}

		report.TotalCallingCount = totalCallingCount
		report.TotalEnquiryCount = totalEnquiryCount
		*reports = append(*reports, report)
	}
	return nil
}

// GetTalentEnquiryCallingReports returns talent-enquiry calling report details with limit and offset.
func (service *CallingReportService) GetTalentEnquiryCallingReports(tenantID uuid.UUID,
	reports *[]callreport.TalentEnquiryCallingReportDTO, form url.Values, limit, offset int, totalCount *int) error {

	// Check if tenant exists
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	uow := repository.NewUnitOfWork(service.DB, true)
	// get all colleges ordered by name order
	selectCredential := "cred.`id` AS `login_id`,CONCAT(cred.`first_name`,' ',cred.`last_name`) AS `login_name`"
	selectTalent := "talent_enquiries.`first_name`,talent_enquiries.`last_name`,talent_enquiries.`contact`,talent_enquiries.`email`"
	selectCallRecord := "talent_enquiry_call_records.`date_time`,talent_enquiry_call_records.`comment`," +
		"talent_enquiry_call_records.`expected_ctc`,talent_enquiry_call_records.`notice_period`,talent_enquiry_call_records.`target_date`"
	selectPurposeOutcome := "purposes.`purpose`,outcomes.`outcome`"
	err = service.Repository.GetAll(uow, reports,
		repository.Select(selectCredential+","+selectTalent+","+selectCallRecord+","+selectPurposeOutcome),
		repository.Table("talent_enquiry_call_records"),
		repository.Join("INNER JOIN credentials cred ON talent_enquiry_call_records.`created_by`=cred.`id`"),
		repository.Join("INNER JOIN talent_enquiries ON talent_enquiry_call_records.`enquiry_id`=talent_enquiries.`id`"),
		repository.Join("INNER JOIN purposes ON talent_enquiry_call_records.`purpose_id`=purposes.`id`"),
		repository.Join("INNER JOIN outcomes ON talent_enquiry_call_records.`outcome_id`=outcomes.`id`"),
		repository.Filter("talent_enquiry_call_records.`deleted_at` IS NULL AND talent_enquiry_call_records.`tenant_id`=?", tenantID),
		repository.Filter("cred.`deleted_at` IS NULL AND cred.`tenant_id`=?", tenantID),
		repository.Filter("purposes.`deleted_at` IS NULL AND purposes.`tenant_id`=?", tenantID),
		repository.Filter("outcomes.`deleted_at` IS NULL AND outcomes.`tenant_id`=?", tenantID),
		service.addSearchQueriesForTalentEnquiryCallingReport(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// returns error if there is no tenant record in table.
func (service *CallingReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *CallingReportService) addSearchQueriesForCallingReport(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if credentialID, ok := requestForm["credentialID"]; ok {
		util.AddToSlice("cred.`id`", "= ?", "AND", credentialID, &columnNames, &conditions, &operators, &values)
	}
	if purposeID, ok := requestForm["purposeID"]; ok {
		util.AddToSlice("purposes.`id`", "= ?", "AND", purposeID, &columnNames, &conditions, &operators, &values)
	}
	if outcomeID, ok := requestForm["outcomeID"]; ok {
		util.AddToSlice("outcomes.`id`", "= ?", "AND", outcomeID, &columnNames, &conditions, &operators, &values)
	}
	if date, ok := requestForm["date"]; ok {
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", "= ?", "AND", date, &columnNames, &conditions, &operators, &values)
	}
	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}
	if targetToDate, ok := requestForm["targetToDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`target_date` AS DATE)", "<= ?", "AND", targetToDate, &columnNames, &conditions, &operators, &values)
	}
	if targetFromDate, ok := requestForm["targetFromDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`target_date` AS DATE)", ">= ?", "AND", targetFromDate, &columnNames, &conditions, &operators, &values)
	}
	if noticePeriod, ok := requestForm["noticePeriod"]; ok {
		util.AddToSlice("talent_call_records.`notice_period`", "<= ?", "AND", noticePeriod, &columnNames, &conditions, &operators, &values)
	}
	if expectedCTC, ok := requestForm["expectedCTC"]; ok {
		util.AddToSlice("talent_call_records.`expected_ctc`", "<= ?", "AND", expectedCTC, &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

func (service *CallingReportService) addSearchQueriesForLoginwise(requestForm url.Values) (repository.QueryProcessor, error) {
	fmt.Println("=========================In loginwise search============================", requestForm)
	if len(requestForm) == 0 {
		return nil, nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if durationString, ok := requestForm["duration"]; ok {
		duration, err := strconv.Atoi(durationString[0])
		if err != nil {
			return nil, err
		}
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, -duration)
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", "<= ?", "AND",
			startDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", ">= ?", "AND",
			endDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
	}
	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(talent_call_records.`date_time` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values), nil
}

// addSearchQueriesForTalentEnquiryCallingReport will add search queries for talent enquiry
func (service *CallingReportService) addSearchQueriesForTalentEnquiryCallingReport(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if credentialID, ok := requestForm["credentialID"]; ok {
		util.AddToSlice("cred.`id`", "= ?", "AND", credentialID, &columnNames, &conditions, &operators, &values)
	}
	if purposeID, ok := requestForm["purposeID"]; ok {
		util.AddToSlice("purposes.`id`", "= ?", "AND", purposeID, &columnNames, &conditions, &operators, &values)
	}
	if outcomeID, ok := requestForm["outcomeID"]; ok {
		util.AddToSlice("outcomes.`id`", "= ?", "AND", outcomeID, &columnNames, &conditions, &operators, &values)
	}
	if date, ok := requestForm["date"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", "= ?", "AND",
			date, &columnNames, &conditions, &operators, &values)
	}
	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", "<= ?", "AND",
			toDate, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", ">= ?", "AND",
			fromDate, &columnNames, &conditions, &operators, &values)
	}
	if targetToDate, ok := requestForm["targetToDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`target_date` AS DATE)", "<= ?", "AND",
			targetToDate, &columnNames, &conditions, &operators, &values)
	}
	if targetFromDate, ok := requestForm["targetFromDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`target_date` AS DATE)", ">= ?", "AND",
			targetFromDate, &columnNames, &conditions, &operators, &values)
	}
	if noticePeriod, ok := requestForm["noticePeriod"]; ok {
		util.AddToSlice("talent_enquiry_call_records.`notice_period`", "<= ?", "AND", noticePeriod,
			&columnNames, &conditions, &operators, &values)
	}
	if expectedCTC, ok := requestForm["expectedCTC"]; ok {
		util.AddToSlice("talent_enquiry_call_records.`expected_ctc`", "<= ?", "AND", expectedCTC,
			&columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

func (service *CallingReportService) addSearchQueriesForEnquiryLoginwise(requestForm url.Values) (repository.QueryProcessor, error) {
	fmt.Println("=========================In loginwise search============================", requestForm)
	if len(requestForm) == 0 {
		return nil, nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if durationString, ok := requestForm["duration"]; ok {
		duration, err := strconv.Atoi(durationString[0])
		if err != nil {
			return nil, err
		}
		startDate := time.Now()
		endDate := startDate.AddDate(0, 0, -duration)
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", "<= ?", "AND",
			startDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", ">= ?", "AND",
			endDate.Format("2006-01-02"), &columnNames, &conditions, &operators, &values)
	}
	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", "<= ?", "AND",
			toDate, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(talent_enquiry_call_records.`date_time` AS DATE)", ">= ?", "AND",
			fromDate, &columnNames, &conditions, &operators, &values)
	}
	return repository.FilterWithOperator(columnNames, conditions, operators, values), nil
}
