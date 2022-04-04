package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProfessionalSummaryReportService provides methods to get professional summary report of talents.
type ProfessionalSummaryReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProfessionalSummaryReportService returns new instance of ProfessionalSummaryReportService.
func NewProfessionalSummaryReportService(db *gorm.DB, repository repository.Repository) *ProfessionalSummaryReportService {
	return &ProfessionalSummaryReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetProfessionalSummaryReport gets professional summary report of talents.
func (service *ProfessionalSummaryReportService) GetProfessionalSummaryReport(proSummaryReport *[]tal.ProfessionalSummaryReport, tenantID uuid.UUID,
	form url.Values, limit int, offset int, proSummaryReportCounts *tal.ProfessionalSummaryReportCounts) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors for sub query number one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Select("IF(`experience_in_months`>=12 AND `experience_in_months`<=24,1,0) as first,"+
			"IF(`experience_in_months`>24 AND `experience_in_months`<=60,1,0) as second,"+
			"IF(`experience_in_months`>60 AND `experience_in_months`<=84,1,0) as third,"+
			"IF(`experience_in_months`>84,1,0) as fourth,talent_experiences.*"),
		repository.Filter("`from_date` IS NOT NULL and `to_date` IS NULL AND talent_experiences.`deleted_at` IS NULL "+
			"AND talent_experiences.`tenant_id`=? AND talents.`deleted_at` IS NULL AND talents.`tenant_id`=?", tenantID, tenantID),
		repository.Table("talent_experiences"),
		repository.Join("LEFT JOIN talents on talents.`id` = talent_experiences.`talent_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, tal.ProfessionalSummaryReport{}, queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, proSummaryReport,
		repository.RawQuery("SELECT SUM(first) as first_count, SUM(second) as second_count, SUM(third) as third_count,"+
			"SUM(fourth) as fourth_count, temp.* FROM  ? as temp GROUP BY `company` order by `company`", subQueryOne),
	); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create subQuery for total count.
	subQueryForCount, err := service.Repository.SubQuery(uow, tal.ProfessionalSummaryReport{},
		repository.RawQuery("SELECT SUM(first) as first_count, SUM(second) as second_count, SUM(third) as third_count,"+
			"SUM(fourth) as fourth_count, temp.* FROM  ? as temp GROUP BY `company`", subQueryOne),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// create bucket for total count.
	total := tal.TotalCount{}

	// Get total count of entries.
	if err := service.Repository.Scan(uow, &total,
		repository.RawQuery("SELECT COUNT(*) as total_count FROM ? as temptwo", subQueryForCount)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give total to total count field of pro summar report counts.
	proSummaryReportCounts.TotalCount = total.TotalCount

	// Get total counts of all categories.
	if err := service.Repository.Scan(uow, &proSummaryReportCounts,
		repository.RawQuery("SELECT SUM(first_count) as first_count_total, SUM(second_count) as second_count_total,"+
			"SUM(third_count) as third_count_total,SUM(fourth_count) as fourth_count_total, temptwo.* FROM  ? as temptwo", subQueryForCount)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetProfessionalSummaryReportByTechCount gets talent counts by company name and technology.
func (service *ProfessionalSummaryReportService) GetProfessionalSummaryReportByTechCount(companyTechTalents *[]tal.CompanyTechnologyTalent, tenantID uuid.UUID,
	form url.Values) error {

	//********************************************Category filter***************************************************
	// Variables for company name and category.
	category := ""

	// Get query params for has appeared.
	categoryArray, _ := form["category"]
	if categoryArray != nil && len(categoryArray) != 0 {
		category = categoryArray[0]
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors for sub query number one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Filter("`from_date` IS NOT NULL and `to_date` IS NULL AND talent_experiences.`deleted_at` IS NULL "+
			"AND talent_experiences.`tenant_id`=? AND talents.`deleted_at` IS NULL AND talents.`tenant_id`=?", tenantID, tenantID),
		repository.Table("talent_experiences"),
		repository.Join("JOIN talents on talents.`id` = talent_experiences.`talent_id`"),
		repository.Join("JOIN talent_experiences_technologies on talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
		repository.Join("JOIN technologies on talent_experiences_technologies.`technology_id` = technologies.`id`"))

	// If category is first then find talents with 12-24 months of experience.
	if category == "first" {
		queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
			repository.Select("IF(`experience_in_months`>=12 AND `experience_in_months`<=24,1,0) as talents,"+
				"talent_experiences.`company`,technologies.`id` as tech_id, technologies.`language` as tech_name"))
	}

	// If category is second then find talents with 24-60 months of experience.
	if category == "second" {
		queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
			repository.Select("IF(`experience_in_months`>=24 AND `experience_in_months`<=60,1,0) as talents,"+
				"talent_experiences.`company`,technologies.`id` as tech_id, technologies.`language` as tech_name"))
	}

	// If category is third then find talents with 60-84 months of experience.
	if category == "third" {
		queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
			repository.Select("IF(`experience_in_months`>=60 AND `experience_in_months`<=84,1,0) as talents,"+
				"talent_experiences.`company`,technologies.`id` as tech_id, technologies.`language` as tech_name"))
	}

	// If category is fourth then find talents with above 84 months of experience.
	if category == "fourth" {
		queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
			repository.Select("IF(`experience_in_months`>84,1,0) as talents,"+
				"talent_experiences.`company`,technologies.`id` as tech_id, technologies.`language` as tech_name"))
	}

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, tal.ProfessionalSummaryReport{}, queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, companyTechTalents,
		repository.RawQuery("SELECT SUM(talents) as talent_count,temp.* FROM  ? as temp GROUP BY `company`, tech_id HAVING talent_count > 0", subQueryOne),
	); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProfessionalSummaryReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
