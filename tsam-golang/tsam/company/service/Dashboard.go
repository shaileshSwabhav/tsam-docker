package services

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/repository"
)

// CompanyDashboardService provides all details to be shown on company dashboard.
type CompanyDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCompanyDashboardService returns new instance of CompanyDashboardService.
func NewCompanyDashboardService(db *gorm.DB, repository repository.Repository) *CompanyDashboardService {
	return &CompanyDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetCompanyDashboardDetails gets all details required for CompanyDashboard.
func (service *CompanyDashboardService) GetCompanyDashboardDetails(companyDashboard *dashboard.CompanyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total companies.
	if err := service.Repository.GetCount(uow, company.Branch{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.TotalCompanies = uint(totalCount)

	// Get count of active companies(companies where at least 1 drive or seminar has been conducted).
	if err := service.Repository.GetCount(uow, company.Branch{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.TotalCompanies = uint(totalCount)

	if err := service.getAllExperienceRequirements(companyDashboard); err != nil {
		return err
	}

	if err := service.getAllFresherRequirements(companyDashboard); err != nil {
		return err
	}

	if err := service.getAllLiveRequirements(companyDashboard); err != nil {
		return err
	}

	return nil

}

// **************** Change from company branches to respective fields. ***********************

// Gets all fresher requirements.
func (service *CompanyDashboardService) getAllFresherRequirements(companyDashboard *dashboard.CompanyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	companyDashboard.FresherRequirements.PlacementsInProcess = 786
	companyDashboard.FresherRequirements.TotalRequirements = 786
	companyDashboard.FresherRequirements.TotalTalentsRequired = 786
	return nil
}

// Gets all live requirements.
func (service *CompanyDashboardService) getAllLiveRequirements(companyDashboard *dashboard.CompanyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	companyDashboard.LiveRequirements.PlacementsInProcess = 786
	companyDashboard.LiveRequirements.TotalRequirements = 786
	companyDashboard.LiveRequirements.TotalTalentsRequired = 786
	return nil
}

// Gets all live requirements.
func (service *CompanyDashboardService) getAllExperienceRequirements(companyDashboard *dashboard.CompanyDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	companyDashboard.ExperienceRequirements.PlacementsInProcess = 786
	companyDashboard.ExperienceRequirements.TotalRequirements = 786
	companyDashboard.ExperienceRequirements.TotalTalentsRequired = 786
	return nil
}
