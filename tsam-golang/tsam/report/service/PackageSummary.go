package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// PackageSummaryService provides method to get package-summary reports.
type PackageSummaryService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	techLanguage []string
	experience   []string
	// packageDistribution []int
}

// NewPackageSummaryService returns a new instance of PackageSummarService.
func NewPackageSummaryService(db *gorm.DB, repository repository.Repository) *PackageSummaryService {
	return &PackageSummaryService{
		DB:         db,
		Repository: repository,
		techLanguage: []string{
			"Advance Java", "Dotnet", "Java", "Machine Learning", "Cloud", "Golang",
		},
		experience: []string{
			"0 to 3 Years", "3 to 6 Years", "6 to 8 Years", "8 to 10 Years", "10+ Years",
		},
		// packageDistribution: []int{
		// 	300000, 500000, 1000000, 1500000,
		// },
	}
}

// GetPackageSummary will return package summary for different pacakges.
func (service *PackageSummaryService) GetPackageSummary(tenantID uuid.UUID,
	packageSummary *[]report.PackageSummary, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	uow := repository.NewUnitOfWork(service.DB, true)

	tempPackageSummary := report.PackageSummary{}

	// Experience => 0.1 - 3 year
	tempPackageSummary.Experience = service.experience[0]
	err = service.getPackageCount(uow, tenantID, &tempPackageSummary, requestForm,
		repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 1, 36))
	if err != nil {
		uow.RollBack()
		return err
	}

	*packageSummary = append(*packageSummary, tempPackageSummary)

	// Experience => 3.1 - 6 year
	tempPackageSummary.Experience = service.experience[1]
	err = service.getPackageCount(uow, tenantID, &tempPackageSummary, requestForm,
		repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 37, 72))
	if err != nil {
		uow.RollBack()
		return err
	}

	*packageSummary = append(*packageSummary, tempPackageSummary)

	// Experience => 6.1 - 8 year
	tempPackageSummary.Experience = service.experience[2]
	err = service.getPackageCount(uow, tenantID, &tempPackageSummary, requestForm,
		repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 73, 96))
	if err != nil {
		uow.RollBack()
		return err
	}
	*packageSummary = append(*packageSummary, tempPackageSummary)

	// Experience => 8.1 - 10 year
	tempPackageSummary.Experience = service.experience[3]
	err = service.getPackageCount(uow, tenantID, &tempPackageSummary, requestForm,
		repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 97, 120))
	if err != nil {
		uow.RollBack()
		return err
	}

	*packageSummary = append(*packageSummary, tempPackageSummary)

	// Experience => 10.1+ year
	tempPackageSummary.Experience = service.experience[4]
	err = service.getPackageCount(uow, tenantID, &tempPackageSummary, requestForm,
		repository.Filter("talents.`experience_in_months` > ?", 120))
	if err != nil {
		uow.RollBack()
		return err
	}

	*packageSummary = append(*packageSummary, tempPackageSummary)

	return nil
}

// GetTechnologyPackageSummary will give package-summary based on technology.
func (service *PackageSummaryService) GetTechnologyPackageSummary(tenantID uuid.UUID,
	technologySummary *[]report.ExperienceTechnologySummary, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	uow := repository.NewUnitOfWork(service.DB, true)

	technologies := []general.Technology{}

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &technologies,
		"`language`", repository.Filter("`language` IN (?)", service.techLanguage))
	if err != nil {
		uow.RollBack()
		return err
	}

	experienceTechSummary := report.ExperienceTechnologySummary{}

	// Experience => 0.1 - 3 year
	experienceTechSummary.Experience = service.experience[0]
	err = service.getExperienceSummary(uow, tenantID, technologies, &experienceTechSummary.TechnologySummary,
		requestForm, repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 1, 36))
	if err != nil {
		uow.RollBack()
		return err
	}

	*technologySummary = append(*technologySummary, experienceTechSummary)

	experienceTechSummary.Experience = service.experience[1]
	experienceTechSummary.TechnologySummary = nil
	err = service.getExperienceSummary(uow, tenantID, technologies, &experienceTechSummary.TechnologySummary,
		requestForm, repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 37, 72))
	if err != nil {
		uow.RollBack()
		return err
	}
	*technologySummary = append(*technologySummary, experienceTechSummary)

	experienceTechSummary.Experience = service.experience[2]
	experienceTechSummary.TechnologySummary = nil
	err = service.getExperienceSummary(uow, tenantID, technologies, &experienceTechSummary.TechnologySummary,
		requestForm, repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 73, 97))
	if err != nil {
		uow.RollBack()
		return err
	}
	*technologySummary = append(*technologySummary, experienceTechSummary)

	experienceTechSummary.Experience = service.experience[3]
	experienceTechSummary.TechnologySummary = nil
	err = service.getExperienceSummary(uow, tenantID, technologies, &experienceTechSummary.TechnologySummary,
		requestForm, repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 97, 120))
	if err != nil {
		uow.RollBack()
		return err
	}
	*technologySummary = append(*technologySummary, experienceTechSummary)

	experienceTechSummary.Experience = service.experience[4]
	experienceTechSummary.TechnologySummary = nil
	err = service.getExperienceSummary(uow, tenantID, technologies, &experienceTechSummary.TechnologySummary,
		requestForm, repository.Filter("talents.`experience_in_months` > ?", 120))
	if err != nil {
		uow.RollBack()
		return err
	}
	*technologySummary = append(*technologySummary, experienceTechSummary)

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// getPackageCount will return count for various package ranges based on condition.
func (service *PackageSummaryService) getPackageCount(uow *repository.UnitOfWork,
	tenantID uuid.UUID, packageSummary *report.PackageSummary, requestForm url.Values,
	queryProcessor repository.QueryProcessor) error {
	err := service.Repository.Scan(uow, packageSummary, repository.Table("talents"),
		repository.Select("COUNT( IF(talent_experiences.`package` <= ?, 1, NULL) ) AS less_than_three, "+
			"COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS three_to_five, "+
			"COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS five_to_ten, "+
			"COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS ten_to_fifteen, "+
			"COUNT( IF(talent_experiences.`package` > ?, 1, NULL) ) AS greater_than_fifteen",
			300000, 300001, 500000, 500001, 1000000, 1000001, 1500000, 1500000),
		repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` "+
			" AND talents.`tenant_id` = talent_experiences.`tenant_id`"), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL AND talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talent_experiences.`package` IS NOT NULL"),
		repository.Filter("talent_experiences.`to_date` IS NULL AND talent_experiences.`from_date` IS NOT NULL"), queryProcessor)
	if err != nil {
		return err
	}
	return nil
}

func (service *PackageSummaryService) getExperienceSummary(uow *repository.UnitOfWork,
	tenantID uuid.UUID, technologies []general.Technology, techSummary *[]report.TechnologyPackageSummary,
	requestForm url.Values, queryProcessor repository.QueryProcessor) error {

	tempTechSummary := report.TechnologyPackageSummary{}

	for index := range technologies {

		tempTechSummary.TechLanguage = technologies[index].Language
		tempTechSummary.Technology = &technologies[index]
		err := service.getTechnologyPackageSummary(uow, tenantID, technologies[index].ID, &tempTechSummary,
			requestForm, queryProcessor)
		if err != nil {
			uow.RollBack()
			return err
		}

		*techSummary = append(*techSummary, tempTechSummary)
	}

	// FOR OTHER TECHNOLOGIES
	// tempTechSummary.Technology = nil
	// tempTechSummary.TechLanguage = "Other"

	// err := service.Repository.Scan(uow, &tempTechSummary, repository.Table("talents"),
	// 	service.addQueries(requestForm), repository.Filter("talents.`tenant_id` = ?", tenantID),
	// 	repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` "+
	// 		" AND talents.`tenant_id` = talent_experiences.`tenant_id`"),
	// 	repository.Join("INNER JOIN talent_experiences_technologies ON "+
	// 		"talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
	// 	repository.Join("INNER JOIN technologies ON talent_experiences_technologies.`technology_id` = technologies.`id`"),
	// 	repository.Filter("technologies.`id` NOT IN (?)", service.techLanguage), repository.Filter("talent_experiences.`package` IS NOT NULL"),
	// 	repository.Filter("talents.`deleted_at` IS NULL AND talent_experiences.`deleted_at` IS NULL"),
	// 	repository.Filter("talent_experiences.`to_date` IS NULL AND talent_experiences.`from_date` IS NOT NULL"), queryProcessor)
	// if err != nil {
	// 	return err
	// }

	// *techSummary = append(*techSummary, tempTechSummary)

	return nil
}

// getTechnologyPackageSummary will get total count for specified technology.
func (service *PackageSummaryService) getTechnologyPackageSummary(uow *repository.UnitOfWork,
	tenantID, techID uuid.UUID, tempTechSummary *report.TechnologyPackageSummary,
	requestForm url.Values, queryProcessor repository.QueryProcessor) error {

	err := service.Repository.Scan(uow, tempTechSummary, repository.Table("talents"),
		service.addQueries(requestForm), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` "+
			"AND talents.`tenant_id` = talent_experiences.`tenant_id`"),
		repository.Join("INNER JOIN talent_experiences_technologies ON "+
			"talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
		repository.Filter("talent_experiences_technologies.`technology_id` = ?", techID), repository.Filter("talent_experiences.`package` IS NOT NULL"),
		repository.Filter("talents.`deleted_at` IS NULL AND talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talent_experiences.`to_date` IS NULL AND talent_experiences.`from_date` IS NOT NULL"), queryProcessor)
	if err != nil {
		return err
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *PackageSummaryService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *PackageSummaryService) addQueries(requestForm url.Values) repository.QueryProcessor {

	// fmt.Println("=========================In addSearchQueries============================", requestForm)
	if len(requestForm) == 0 {
		// DISTINCT talent_experiences_technologies.`experience_id`
		return repository.Select("COUNT(*) AS `total_count`")
		// return repository.Select("COUNT(DISTINCT talent_experiences.`talent_id`) AS `total_count`")
	}

	packageTypeArray := requestForm["packageType"]
	if len(packageTypeArray) != 0 {
		packageType := requestForm.Get("packageType")
		if packageType == "LessThanThree" {
			return repository.Select("COUNT( IF(talent_experiences.`package` <= ?, 1, NULL) ) AS `total_count`", 300000)
		}
		if packageType == "ThreeToFive" {
			return repository.Select("COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS `total_count`", 300001, 500000)
		}
		if packageType == "FiveToTen" {
			return repository.Select("COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS `total_count`", 500001, 1000000)
		}
		if packageType == "TenToFifteen" {
			return repository.Select("COUNT( IF(talent_experiences.`package` BETWEEN ? AND ?, 1, NULL) ) AS `total_count`", 1000001, 1500000)
		}
		if packageType == "GreaterThanFifteen" {
			return repository.Select("COUNT( IF(talent_experiences.`package` > ?, 1, NULL) ) AS `total_count`", 1500000)
		}
	}

	return nil

}
