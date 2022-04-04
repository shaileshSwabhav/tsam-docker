package service

import (
	"net/url"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SalaryTrendService Provide method to Update, Delete, Add, Get Method For salary-trend.
type SalaryTrendService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewSalaryTrendService returns new instance of SalaryTrendService.
func NewSalaryTrendService(db *gorm.DB, repository repository.Repository) *SalaryTrendService {
	return &SalaryTrendService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Technology", "Designation",
		},
	}
}

// AddSalaryTrend will add new salary-trend to the table.
func (ser *SalaryTrendService) AddSalaryTrend(salaryTrend *general.SalaryTrend) error {

	// Checks if all foreign keys are present.
	err := ser.doesForeignKeyExist(salaryTrend, salaryTrend.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.Add(uow, salaryTrend)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateSalaryTrend will update the specifed salary-trend in the table.
func (ser *SalaryTrendService) UpdateSalaryTrend(salaryTrend *general.SalaryTrend) error {

	// Checks if all foreign keys are present.
	err := ser.doesForeignKeyExist(salaryTrend, salaryTrend.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if salary-trend exist.
	err = ser.doesSalaryTrendExist(salaryTrend.TenantID, salaryTrend.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	tempSalaryTrend := general.SalaryTrend{}

	// Get createdBy and code for specified salary-trend.
	err = ser.Repository.GetRecordForTenant(uow, salaryTrend.TenantID, &tempSalaryTrend,
		repository.Filter("`id` = ?", salaryTrend.ID), repository.Select([]string{"`created_by`"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	salaryTrend.CreatedBy = tempSalaryTrend.CreatedBy

	err = ser.Repository.Save(uow, salaryTrend)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteSalaryTrend will delete the specifed salary-trend from the table.
func (ser *SalaryTrendService) DeleteSalaryTrend(salaryTrend *general.SalaryTrend) error {

	// Check if tenant exist.
	err := ser.doesTenantExists(salaryTrend.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = ser.doesCredentialExist(salaryTrend.TenantID, salaryTrend.DeletedBy)
	if err != nil {
		return err
	}

	// Check if salary-trend exist.
	err = ser.doesSalaryTrendExist(salaryTrend.TenantID, salaryTrend.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.UpdateWithMap(uow, &general.SalaryTrend{}, map[string]interface{}{
		"DeletedBy": salaryTrend.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("`id` = ?", salaryTrend.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllSalaryTrend will get all the salary-trend from the table.
func (ser *SalaryTrendService) GetAllSalaryTrend(salaryTrends *[]general.SalaryTrendDTO, tenantID uuid.UUID,
	limit, offset int, totalCount *int, requestForm url.Values) error {

	// Check if tenant exists.
	err := ser.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, true)

	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, salaryTrends, "`date`",
		ser.addSearchQueriesParams(requestForm), repository.PreloadAssociations(ser.associations),
		repository.Paginate(limit, offset, totalCount))
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

// doesForeignKeyExist will check all the foreign keys of salary-trend.
func (ser *SalaryTrendService) doesForeignKeyExist(salaryTrend *general.SalaryTrend, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := ser.doesTenantExists(salaryTrend.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = ser.doesCredentialExist(salaryTrend.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if technology exist.
	err = ser.doesTechnologyExist(salaryTrend.TenantID, salaryTrend.Technology.ID)
	if err != nil {
		return err
	}

	// Check if designation exist.
	err = ser.doesDesignationExist(salaryTrend.TenantID, salaryTrend.Designation.ID)
	if err != nil {
		return err
	}

	// Check if experience range exist.
	err = ser.doesExperienceRangeExist(salaryTrend)
	if err != nil {
		return err
	}

	salaryTrend.TechnologyID = salaryTrend.Technology.ID
	salaryTrend.DesignationID = salaryTrend.Designation.ID

	return nil
}

// doesTenantExists validates tenantID.
func (ser *SalaryTrendService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates credentialID.
func (ser *SalaryTrendService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSalaryTrendExist validates salaryTrendID.
func (ser *SalaryTrendService) doesSalaryTrendExist(tenantID, salaryTrendID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.SalaryTrend{},
		repository.Filter("`id` = ?", salaryTrendID))
	if err := util.HandleError("Salary trend not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTechnologyExist validates technologyID.
func (ser *SalaryTrendService) doesTechnologyExist(tenantID, technologyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Technology{},
		repository.Filter("`id` = ?", technologyID))
	if err := util.HandleError("Technology not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDesignationExist validates technologyID.
func (ser *SalaryTrendService) doesDesignationExist(tenantID, designationID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Designation{},
		repository.Filter("`id` = ?", designationID))
	if err := util.HandleError("Designation not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesExperienceRangeExist will check if experience for given range exist.
func (ser *SalaryTrendService) doesExperienceRangeExist(trend *general.SalaryTrend) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, trend.TenantID, general.SalaryTrend{},
		repository.Filter("((`minimum_experience` <= ? AND `maximum_experience` >= ?) AND "+
			" `maximum_experience` = ?) AND (`technology_id` = ? AND `designation_id` = ?) AND `id` != ?",
			trend.MinimumExperience, trend.MaximumExperience, trend.MinimumExperience, trend.Technology.ID,
			trend.Designation.ID, trend.ID))
	if err := util.HandleIfExistsError(trend.Technology.Language+" has experience in range "+strconv.Itoa(int(trend.MinimumExperience))+
		"-"+strconv.Itoa(int(trend.MaximumExperience)), exists, err); err != nil {
		return err
	}
	return nil
}

// addSearchQueriesParams will search for
func (ser *SalaryTrendService) addSearchQueriesParams(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	// var queryProcessors []repository.QueryProcessor

	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(`salary_trends`.`date` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}

	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(`salary_trends`.`date` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}

	if companyRating, ok := requestForm["companyRating"]; ok {
		util.AddToSlice("`company_rating`", "= ?", "AND", companyRating, &columnNames, &conditions, &operators, &values)
	}

	if minimumExperience, ok := requestForm["minimumExperience"]; ok {
		util.AddToSlice("`minimum_experience`", ">= ?", "AND", minimumExperience, &columnNames, &conditions, &operators, &values)
	}

	if maximumExperience, ok := requestForm["maximumExperience"]; ok {
		util.AddToSlice("`maximum_experience`", "<= ?", "AND", maximumExperience, &columnNames, &conditions, &operators, &values)
	}

	if technology, ok := requestForm["technology"]; ok {
		util.AddToSlice("`technology_id`", "= ?", "AND", technology, &columnNames, &conditions, &operators, &values)
	}

	if designation, ok := requestForm["designation"]; ok {
		util.AddToSlice("`designation_id`", "= ?", "AND", designation, &columnNames, &conditions, &operators, &values)
	}

	// return queryProcessors
	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
