package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FresherSummaryService provides method to get fresher-summary reports.
type FresherSummaryService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	techLanguage []string
	academicYear []string
	jobSwitch    string
}

// NewFresherSummaryService returns a new instance of FresherSummaryService.
func NewFresherSummaryService(db *gorm.DB, repository repository.Repository) *FresherSummaryService {
	return &FresherSummaryService{
		DB:         db,
		Repository: repository,
		techLanguage: []string{
			"Advance Java", "Dotnet", "Java", "Machine Learning", "Cloud", "Golang",
		},
		academicYear: []string{
			"Professional", "Looking to switch", "Current Requirements",
		},
		jobSwitch: "Job switch",
	}
}

// GetAllFresherSummary will return fresher-summary
func (service *FresherSummaryService) GetAllFresherSummary(tenantID uuid.UUID,
	fresherSummary *[]report.FresherSummary, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	uow := repository.NewUnitOfWork(service.DB, false)

	academics := []general.CommonType{}
	// var totalCount uint

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &academics, "`key`",
		repository.Table("general_types"), repository.Filter("`type` = 'academic_year'"))
	if err != nil {
		uow.RollBack()
		return err
	}

	var queryProcessor []repository.QueryProcessor
	var subQueryProcessor []repository.QueryProcessor
	tempSummary := report.FresherSummary{}

	for _, academic := range academics {
		tempSummary.ColumnName = academic.Value
		tempSummary.Academic = academic
		queryProcessor = nil

		queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
			repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = ?", academic.Key))

		if academic.Key == 5 {
			queryProcessor = append(queryProcessor, repository.Filter("`talents`.`is_experience` = '0'"))
		}

		err = service.getTalentCount(uow, tenantID, queryProcessor, &tempSummary)
		if err != nil {
			uow.RollBack()
			return err
		}

		*fresherSummary = append(*fresherSummary, tempSummary)
	}

	// Calculate Professional.
	tempSummary.ColumnName = service.academicYear[0]
	queryProcessor = nil
	queryProcessor = append(queryProcessor, repository.Table("talents"),
		repository.Filter("`talents`.`is_experience` = '1'"), repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talents.`academic_year` = 5"))

	err = service.getTalentCount(uow, tenantID, queryProcessor, &tempSummary)
	if err != nil {
		uow.RollBack()
		return err
	}

	*fresherSummary = append(*fresherSummary, tempSummary)

	// Calculate for job seekers.
	// tempSummary.AcademicYear = "Job Seekers"
	tempSummary.ColumnName = service.academicYear[1]
	queryProcessor = nil
	subQueryProcessor = nil

	subQueryProcessor = append(subQueryProcessor, repository.Table("talent_call_records"),
		repository.Select("max(`date_time`)"), repository.GroupBy("talent_call_records.`talent_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, talent.CallRecord{}, subQueryProcessor...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	queryProcessor = append(queryProcessor, repository.Table("talents"),
		repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
			"talents.`tenant_id` = talent_call_records.`tenant_id`"),
		repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id` = outcomes.`id` AND "+
			"outcomes.`tenant_id` = talent_call_records.`tenant_id`"), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Filter("`talents`.`is_experience` = '1'"), repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talent_call_records.`deleted_at` IS NULL AND outcomes.`deleted_at` IS NULL"),
		repository.Filter("talent_call_records.`date_time` IN ? AND outcomes.`outcome` = ?", subQueryOne, service.jobSwitch))

	err = service.getTalentCount(uow, tenantID, queryProcessor, &tempSummary)
	if err != nil {
		uow.RollBack()
		return err
	}

	*fresherSummary = append(*fresherSummary, tempSummary)

	// Calculate for requirement available.
	tempSummary.ColumnName = service.academicYear[2]
	queryProcessor = nil

	queryProcessor = append(queryProcessor, repository.Table("company_requirements"),
		repository.Select(
			"COUNT( CASE WHEN company_requirements.`talent_rating` >= 5 THEN 1 ELSE NULL END) AS `outstanding_count`, "+
				"COUNT( CASE WHEN (company_requirements.`talent_rating` BETWEEN 3 AND 4) THEN 1 ELSE NULL END) AS `excellent_count`, "+
				"COUNT( CASE WHEN company_requirements.`talent_rating` <= 2 THEN 1 ELSE NULL END) AS `average_count`, "+
				"COUNT( CASE WHEN company_requirements.`talent_rating` IS NULL THEN 1 ELSE NULL END) AS `unranked_count`"),
		repository.Filter("company_requirements.`tenant_id` = ? AND company_requirements.`deleted_at` IS NULL", tenantID))

	err = service.Repository.Scan(uow, &tempSummary, queryProcessor...)
	if err != nil {
		uow.RollBack()
		return err
	}

	*fresherSummary = append(*fresherSummary, tempSummary)

	uow.Commit()
	return nil
}

// GetSummaryForAcademicTechnology will get technology wise fresher-summary for all academic year.
func (service *FresherSummaryService) GetSummaryForAcademicTechnology(tenantID uuid.UUID,
	technologySummary *[]report.AcademicTechnologySummary, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	academics := []general.CommonType{}
	tempTechnology := []general.Technology{}

	// temp variables
	tempAcademicSummary := report.AcademicTechnologySummary{}

	var queryProcessor []repository.QueryProcessor
	// var totalCount uint

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &tempTechnology,
		"`language`", repository.Filter("`language` IN (?)", service.techLanguage))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &academics, "`key`",
		repository.Filter("`type` = 'academic_year'"))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range academics {
		// Set academic value.
		tempAcademicSummary.Academic = &(academics)[index]
		tempAcademicSummary.ColumnName = (academics)[index].Value

		// creating temp tech summary slice.
		techSummary := []report.TechnologySummary{}

		for j := range tempTechnology {
			tempTechSummary := report.TechnologySummary{}
			tempTechSummary.Technology = &(tempTechnology)[j]
			tempTechSummary.TechnologyLangugage = (tempTechnology)[j].Language

			queryProcessor = nil

			if (academics)[index].Key == 5 {
				queryProcessor = append(queryProcessor, repository.Filter("`talents`.`is_experience` = '0'"))
			}
			queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
				repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = ?", (academics)[index].Key),
				repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
				repository.Filter("talents_technologies.`technology_id` = ? ", (tempTechnology)[j].ID), service.addSearchQueries(requestForm))

			err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
			if err != nil {
				uow.RollBack()
				return err
			}

			techSummary = append(techSummary, tempTechSummary)
		}
		// Other techs which are not specified to be fetched.
		tempTechSummary := report.TechnologySummary{
			TechnologyLangugage: "Other",
			Technology:          nil,
			Academic:            &academics[index],
			ColumnName:          academics[index].Value,
			TotalCount:          0,
		}
		err = service.getOtherTechnologyCount(uow, tenantID, &tempTechSummary,
			int(academics[index].Key), requestForm)
		if err != nil {
			uow.RollBack()
			return err
		}
		techSummary = append(techSummary, tempTechSummary)

		tempAcademicSummary.TechnologySummary = techSummary
		*technologySummary = append(*technologySummary, tempAcademicSummary)
	}

	// ======================================================== PROFESSIONALS ========================================================

	// Calculate Professional.
	tempAcademicSummary.ColumnName = service.academicYear[0]
	techSummary := []report.TechnologySummary{}

	for index := range tempTechnology {

		tempTechSummary := report.TechnologySummary{}
		tempTechSummary.Technology = &(tempTechnology)[index]
		tempTechSummary.TechnologyLangugage = (tempTechnology)[index].Language

		queryProcessor = nil
		queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
			repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = 5"),
			repository.Filter("`talents`.`is_experience` = '1'"), service.addSearchQueries(requestForm),
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Filter("talents_technologies.`technology_id` = ? ", tempTechnology[index].ID))

		err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
		if err != nil {
			uow.RollBack()
			return err
		}
		techSummary = append(techSummary, tempTechSummary)
	}

	// Caculate other for professional.
	tempTechSummary := report.TechnologySummary{
		ColumnName:          service.academicYear[0],
		TechnologyLangugage: "Other",
		Academic:            nil,
		Technology:          nil,
		TotalCount:          0,
	}
	err = service.getOtherTechnologyCount(uow, tenantID, &tempTechSummary, 5, requestForm)
	if err != nil {
		uow.RollBack()
		return err
	}
	techSummary = append(techSummary, tempTechSummary)
	tempAcademicSummary.TechnologySummary = techSummary
	*technologySummary = append(*technologySummary, tempAcademicSummary)

	// ======================================================== JOB SEEKERS ========================================================
	// Calculate for job seekers.
	var subQueryProcessor []repository.QueryProcessor
	tempAcademicSummary.ColumnName = service.academicYear[1]
	techSummary = []report.TechnologySummary{}

	subQueryProcessor = append(subQueryProcessor, repository.Table("talent_call_records"),
		repository.Select("max(talent_call_records.`date_time`)"), repository.GroupBy("talent_call_records.`talent_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, talent.CallRecord{}, subQueryProcessor...)
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range tempTechnology {

		tempTechSummary := report.TechnologySummary{}
		tempTechSummary.Technology = &(tempTechnology)[index]
		tempTechSummary.TechnologyLangugage = (tempTechnology)[index].Language

		queryProcessor = nil
		queryProcessor = append(queryProcessor, repository.Table("talents"),
			repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
				"talents.`tenant_id` = talent_call_records.`tenant_id`"),
			repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id` = outcomes.`id` AND "+
				"outcomes.`tenant_id` = talent_call_records.`tenant_id`"),
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Filter("talents_technologies.`technology_id` = ? ", tempTechnology[index].ID),
			repository.Filter("talents.`tenant_id` = ?", tenantID), repository.Filter("`talents`.`is_experience` = '1'"),
			repository.Filter("talents.`deleted_at` IS NULL"), service.addSearchQueries(requestForm),
			repository.Filter("talent_call_records.`deleted_at` IS NULL AND outcomes.`deleted_at` IS NULL"),
			repository.Filter("talent_call_records.`date_time` IN ? AND outcomes.`outcome` = ?", subQueryOne, service.jobSwitch))

		err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
		if err != nil {
			uow.RollBack()
			return err
		}
		techSummary = append(techSummary, tempTechSummary)
	}

	// Caculate other for job seekers.
	tempTechSummary = report.TechnologySummary{
		ColumnName:          service.academicYear[1],
		TechnologyLangugage: "Other",
		Academic:            nil,
		Technology:          nil,
		TotalCount:          0,
	}
	queryProcessor = nil
	queryProcessor = append(queryProcessor, repository.Table("talents"),
		repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
			"talents.`tenant_id` = talent_call_records.`tenant_id`"),
		repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id` = outcomes.`id` AND "+
			"outcomes.`tenant_id` = talent_call_records.`tenant_id`"),
		repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
		repository.Join("INNER JOIN technologies ON technologies.`id` = talents_technologies.`technology_id`"),
		repository.Filter("talents.`tenant_id` = ?", tenantID), repository.Filter("`talents`.`is_experience` = '1'"),
		repository.Filter("talents.`deleted_at` IS NULL"), service.addSearchQueries(requestForm),
		repository.Filter("technologies.`tenant_id` = ? AND technologies.`deleted_at` IS NULL", tenantID),
		repository.Filter("technologies.`language` NOT IN (?)", service.techLanguage),
		repository.Filter("talent_call_records.`deleted_at` IS NULL AND outcomes.`deleted_at` IS NULL"),
		repository.Filter("talent_call_records.`date_time` IN ? AND outcomes.`outcome` = ?", subQueryOne, service.jobSwitch))

	err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
	if err != nil {
		return err
	}
	techSummary = append(techSummary, tempTechSummary)
	tempAcademicSummary.TechnologySummary = techSummary
	*technologySummary = append(*technologySummary, tempAcademicSummary)

	// ======================================================== REQUIREMENT ========================================================

	// Calculate Requirement.
	tempAcademicSummary.ColumnName = service.academicYear[2]
	techSummary = []report.TechnologySummary{}

	for index := range tempTechnology {

		tempTechSummary := report.TechnologySummary{}
		tempTechSummary.Technology = &(tempTechnology)[index]
		tempTechSummary.TechnologyLangugage = (tempTechnology)[index].Language

		queryProcessor = nil
		queryProcessor = append(queryProcessor, repository.Table("company_requirements"),
			repository.Filter("company_requirements.`tenant_id` = ? AND company_requirements.`deleted_at` IS NULL", tenantID),
			service.addSearchQueriesForRequirement(requestForm),
			repository.Join("INNER JOIN company_requirements_technologies ON company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
			repository.Filter("company_requirements_technologies.`technology_id` = ? ", tempTechnology[index].ID))

		err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
		if err != nil {
			uow.RollBack()
			return err
		}
		techSummary = append(techSummary, tempTechSummary)
	}

	// Caculate other for requirement.
	tempTechSummary = report.TechnologySummary{
		ColumnName:          service.academicYear[2],
		TechnologyLangugage: "Other",
		Academic:            nil,
		Technology:          nil,
		TotalCount:          0,
	}
	queryProcessor = nil
	queryProcessor = append(queryProcessor, repository.Table("company_requirements"),
		repository.Filter("company_requirements.`tenant_id` = ? AND company_requirements.`deleted_at` IS NULL", tenantID),
		service.addSearchQueriesForRequirement(requestForm),
		repository.Join("INNER JOIN company_requirements_technologies ON company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
		repository.Join("INNER JOIN technologies ON technologies.`id` = company_requirements_technologies.`technology_id`"),
		repository.Filter("technologies.`tenant_id` = ? AND technologies.`deleted_at` IS NULL", tenantID),
		repository.Filter("technologies.`language` NOT IN (?)", service.techLanguage))

	err = service.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
	if err != nil {
		return err
	}

	techSummary = append(techSummary, tempTechSummary)
	tempAcademicSummary.TechnologySummary = techSummary
	*technologySummary = append(*technologySummary, tempAcademicSummary)

	uow.Commit()
	return nil
}

// GetSummaryTechnologyList will return list of technologies used in fresher summary
func (service *FresherSummaryService) GetSummaryTechnologyList(tenantID uuid.UUID,
	technologies *[]general.Technology) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, technologies,
		"`language`", repository.Filter("`language` IN (?)", service.techLanguage))
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

func (service *FresherSummaryService) getTalentCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
	queryProcessors []repository.QueryProcessor, tempSummary *report.FresherSummary) error {

	queryProcessors = append(queryProcessors, repository.Select(
		"COUNT( CASE WHEN talents.`talent_type` >= 5 THEN 1 ELSE NULL END) AS `outstanding_count`, "+
			"COUNT( CASE WHEN (talents.`talent_type` BETWEEN 3 AND 4) THEN 1 ELSE NULL END) AS `excellent_count`, "+
			"COUNT( CASE WHEN talents.`talent_type` <= 2 THEN 1 ELSE NULL END) AS `average_count`, "+
			"COUNT( CASE WHEN talents.`talent_type` IS NULL THEN 1 ELSE NULL END) AS `unranked_count`"))

	err := service.Repository.Scan(uow, &tempSummary, queryProcessors...)
	if err != nil {
		return err
	}

	return nil
}

func (service *FresherSummaryService) getOtherTechnologyCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
	tempTechSummary *report.TechnologySummary, academicYear int, requestForm url.Values) error {

	var queryProcessor []repository.QueryProcessor
	if tempTechSummary.ColumnName == service.academicYear[0] {
		queryProcessor = append(queryProcessor, repository.Filter("`talents`.`is_experience` = '1'"))
	} else {
		queryProcessor = append(queryProcessor, repository.Filter("`talents`.`is_experience` = '0'"))
	}

	queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = ?", academicYear),
		repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
		repository.Join("INNER JOIN technologies ON technologies.`id` = talents_technologies.`technology_id`"),
		repository.Filter("technologies.`tenant_id` = ? AND technologies.`deleted_at` IS NULL", tenantID),
		repository.Filter("technologies.`language` NOT IN (?)", service.techLanguage))

	queryProcessor = append(queryProcessor, service.addSearchQueries(requestForm))
	err := service.Repository.Scan(uow, tempTechSummary, queryProcessor...)
	if err != nil {
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (service *FresherSummaryService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// addSearchQueries will add search queries to queryprocessor.
func (service *FresherSummaryService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	// fmt.Println("=========================In addSearchQueries============================", requestForm)
	if len(requestForm) == 0 {
		return repository.Select("COUNT(*) AS `total_count`")
	}

	talentTypeArray := requestForm["talentType"]
	if len(talentTypeArray) != 0 {
		if requestForm.Get("talentType") == "Outstanding" {
			return repository.Select("COUNT( CASE WHEN talents.`talent_type` >= 5 THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Excellent" {
			return repository.Select("COUNT( CASE WHEN (talents.`talent_type` BETWEEN 3 AND 4) THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Average" {
			return repository.Select("COUNT( CASE WHEN talents.`talent_type` <= 2 THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Unranked" {
			return repository.Select("COUNT( CASE WHEN talents.`talent_type` IS NULL THEN 1 ELSE NULL END) AS `total_count`")
		}
	}

	return nil
}

// addSearchQueriesForRequirement will add search queries to queryprocessor.
func (service *FresherSummaryService) addSearchQueriesForRequirement(requestForm url.Values) repository.QueryProcessor {

	// fmt.Println("=========================In addSearchQueries============================", requestForm)
	if len(requestForm) == 0 {
		return repository.Select("COUNT(DISTINCT company_requirements_technologies.`requirement_id`) AS `total_count`")
	}

	talentTypeArray := requestForm["talentType"]
	if len(talentTypeArray) != 0 {
		if requestForm.Get("talentType") == "Outstanding" {
			return repository.Select("COUNT( CASE WHEN company_requirements.`talent_rating` >= 5 THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Excellent" {
			return repository.Select("COUNT( CASE WHEN (company_requirements.`talent_rating` BETWEEN 3 AND 4) THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Average" {
			return repository.Select("COUNT( CASE WHEN company_requirements.`talent_rating` <= 2 THEN 1 ELSE NULL END) AS `total_count`")
		}
		if requestForm.Get("talentType") == "Unranked" {
			return repository.Select("COUNT( CASE WHEN company_requirements.`talent_rating` IS NULL THEN 1 ELSE NULL END) AS `total_count`")
		}
	}

	return nil
}

// func (ser *FresherSummaryService) getAllCounts(uow *repository.UnitOfWork, tenantID uuid.UUID,
// 	tempSummary *report.FresherSummary, academicYear int, queryProcessor []repository.QueryProcessor) error {

// 	var totalCount uint

// 	err := ser.getOutstandingCount(uow, tenantID, &totalCount, academicYear, queryProcessor)
// 	if err != nil {
// 		return err
// 	}
// 	tempSummary.OutstandingCount = totalCount

// 	err = ser.getExcellentCount(uow, tenantID, &totalCount, academicYear, queryProcessor)
// 	if err != nil {
// 		return err
// 	}
// 	tempSummary.ExcellentCount = totalCount

// 	err = ser.getAverageCount(uow, tenantID, &totalCount, academicYear, queryProcessor)
// 	if err != nil {
// 		return err
// 	}
// 	tempSummary.AverageCount = totalCount
// 	return nil
// }

// func (ser *FresherSummaryService) getTechnologyCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
// 	tempSummary *report.TechnologySummary, academicYear int, queryProcessor []repository.QueryProcessor,
// 	requestForm url.Values) error {

// 	if requestForm.Get("talentType") == "Outstanding" {
// 		err := ser.getOutstandingCount(uow, tenantID, &tempSummary.TotalCount,
// 			academicYear, queryProcessor)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}
// 	if requestForm.Get("talentType") == "Excellent" {
// 		err := ser.getExcellentCount(uow, tenantID, &tempSummary.TotalCount,
// 			academicYear, queryProcessor)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}
// 	if requestForm.Get("talentType") == "Average" {
// 		err := ser.getAverageCount(uow, tenantID, &tempSummary.TotalCount,
// 			academicYear, queryProcessor)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 	}
// 	return nil
// }

// getOutstandingCount calucaltes count of talents whose feedback score is greater than 8.
// func (ser *FresherSummaryService) getOutstandingCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
// 	totalCount *uint, academicYear int, queryProcessor []repository.QueryProcessor) error {

// 	queryProcessor = append(queryProcessor,
// 		repository.Filter("talents.`talent_type` >= 5 AND talents.`academic_year` = ?", academicYear),
// 		repository.Filter("talents.`tenant_id` = ?", tenantID))
// 	err := ser.Repository.GetCount(uow, &talent.Talent{}, totalCount, queryProcessor...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // getExcellentCount calucaltes count of talents whose feedback score is less than 8 and greater than 5.
// func (ser *FresherSummaryService) getExcellentCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
// 	totalCount *uint, academicYear int, queryProcessor []repository.QueryProcessor) error {

// 	queryProcessor = append(queryProcessor,
// 		repository.Filter("talents.`talent_type` BETWEEN 3 AND 4 AND talents.`academic_year` = ?", academicYear),
// 		repository.Filter("talents.`tenant_id` = ?", tenantID))
// 	err := ser.Repository.GetCount(uow, &talent.Talent{}, totalCount, queryProcessor...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // getAverageCount calucaltes count of talents whose feedback score is less than 5.
// func (ser *FresherSummaryService) getAverageCount(uow *repository.UnitOfWork, tenantID uuid.UUID,
// 	totalCount *uint, academicYear int, queryProcessor []repository.QueryProcessor) error {

// 	queryProcessor = append(queryProcessor,
// 		repository.Filter("talents.`talent_type` <= 2 AND talents.`academic_year` = ?", academicYear),
// 		repository.Filter("talents.`tenant_id` = ?", tenantID))
// 	err := ser.Repository.GetCount(uow, &talent.Talent{}, totalCount, queryProcessor...)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // GetTechnologySummary will get technology wise fresher-summary for specified academic year.
// func (ser *FresherSummaryService) GetTechnologySummary(tenantID uuid.UUID,
// 	technologySummary *[]report.TechnologySummary, requestForm url.Values) error {

// 	// Check if tenant exist.
// 	err := ser.doesTenantExist(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, true)

// 	tempTechnology := []general.Technology{}
// 	var totalCount uint

// 	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempTechnology,
// 		"`language`", repository.Filter("`language` IN ('Dotnet', 'Java', 'Machine Learning' , 'Cloud', 'Golang')"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	for _, tech := range tempTechnology {
// 		tempTechSummary := report.TechnologySummary{
// 			Technology: tech,
// 		}

// 		err = ser.Repository.GetCount(uow, &talent.Talent{}, &totalCount,
// 			ser.addSearchQueries(requestForm), repository.Filter("talents.`tenant_id` = ?", tenantID),
// 			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
// 			repository.Filter("talents_technologies.`technology_id` = ? ", tech.ID))
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 		tempTechSummary.TotalCount = totalCount
// 		*technologySummary = append(*technologySummary, tempTechSummary)
// 	}

// 	uow.Commit()
// 	return nil
// }

// // GetSummaryForTechnology will get technology wise fresher-summary for all academic year.
// func (ser *FresherSummaryService) GetSummaryForTechnology(tenantID uuid.UUID,
// 	technologySummary *[]report.TechnologySummary, requestForm url.Values) error {

// 	// Check if tenant exist.
// 	err := ser.doesTenantExist(tenantID)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(ser.DB, true)

// 	academics := []general.CommonType{}
// 	tempTechnology := []general.Technology{}
// 	var queryProcessor []repository.QueryProcessor
// 	// var totalCount uint

// 	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &tempTechnology,
// 		"`language`", repository.Filter("`language` IN (?)", ser.techLanguage))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	err = ser.Repository.GetAllInOrderForTenant(uow, tenantID, &academics, "`key`",
// 		repository.Filter("`type` = 'academic_year'"))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	for j, academic := range academics {
// 		tempTechSummary := report.TechnologySummary{
// 			AcademicYear:        academic.Value,
// 			Academic:            &academic,
// 			TechnologyLangugage: "",
// 			Technology:          nil,
// 			TotalCount:          0,
// 		}
// 		for index, tech := range tempTechnology {
// 			tempTechSummary.Technology = &tempTechnology[index]
// 			tempTechSummary.TechnologyLangugage = tech.Language
// 			queryProcessor = nil

// 			if academic.Key == 5 {
// 				queryProcessor = append(queryProcessor, repository.Filter("`talents`.`is_experience` = '0'"))
// 			}

// 			queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
// 				repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = ?", academic.Key),
// 				repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
// 				repository.Filter("talents_technologies.`technology_id` = ? ", tech.ID), ser.addSearchQueries(requestForm))

// 			// if _, ok := requestForm["talentType"]; ok {

// 			err = ser.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
// 			if err != nil {
// 				uow.RollBack()
// 				return err
// 			}
// 			// }
// 			*technologySummary = append(*technologySummary, tempTechSummary)

// 		}
// 		// Other techs which are not specified to be fetch.
// 		tempTechSummary = report.TechnologySummary{
// 			TechnologyLangugage: "Other",
// 			Technology:          nil,
// 			Academic:            &academics[j],
// 			AcademicYear:        academic.Value,
// 			TotalCount:          0,
// 		}
// 		err = ser.getOtherTechnologyCount(uow, tenantID, &tempTechSummary,
// 			int(academics[j].Key), requestForm)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 		*technologySummary = append(*technologySummary, tempTechSummary)
// 	}

// 	for index, tech := range tempTechnology {
// 		// Calculate for professional.
// 		tempTechSummary := report.TechnologySummary{
// 			AcademicYear:        ser.academicYear[0],
// 			Technology:          &tempTechnology[index],
// 			TechnologyLangugage: tech.Language,
// 			TotalCount:          0,
// 			Academic:            nil,
// 		}
// 		queryProcessor = nil
// 		queryProcessor = append(queryProcessor, repository.Table("talents"), repository.Filter("talents.`tenant_id` = ?", tenantID),
// 			repository.Filter("talents.`deleted_at` IS NULL"), repository.Filter("talents.`academic_year` = 5"),
// 			repository.Filter("`talents`.`is_experience` = '1'"), ser.addSearchQueries(requestForm),
// 			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
// 			repository.Filter("talents_technologies.`technology_id` = ? ", tech.ID))

// 		err = ser.Repository.Scan(uow, &tempTechSummary, queryProcessor...)
// 		if err != nil {
// 			uow.RollBack()
// 			return err
// 		}
// 		// tempTechSummary.TotalCount = totalCount
// 		*technologySummary = append(*technologySummary, tempTechSummary)
// 	}
// 	tempTechSummary := report.TechnologySummary{
// 		AcademicYear:        ser.academicYear[0],
// 		TechnologyLangugage: "Other",
// 		Academic:            nil,
// 		Technology:          nil,
// 		TotalCount:          0,
// 	}
// 	err = ser.getOtherTechnologyCount(uow, tenantID, &tempTechSummary, 5, requestForm)
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	*technologySummary = append(*technologySummary, tempTechSummary)

// 	uow.Commit()
// 	return nil
// }
