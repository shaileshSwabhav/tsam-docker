package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	colg "github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/repository"
)

// CollegeDashboardService provides all details to be shown on college dashboard.
type CollegeDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCollegeDashboardService returns new instance of CollegeDashboardService.
func NewCollegeDashboardService(db *gorm.DB, repository repository.Repository) *CollegeDashboardService {
	return &CollegeDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetCollegeDashboardDetails gets all details required for CollegeDashboard.
func (service *CollegeDashboardService) GetCollegeDashboardDetails(collegeDashboard *dashboard.CollegeDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total colleges.
	if err := service.Repository.GetCount(uow, colg.Branch{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.TotalColleges = uint(totalCount)

	// Get count of active colleges(colleges where at least 1 drive or seminar has been conducted).
	if err := service.Repository.GetCount(uow, colg.College{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.TotalActiveColleges = uint(totalCount)

	if err := service.getAllCampusDetails(collegeDashboard); err != nil {
		return err
	}

	if err := service.getAllSeminarDetails(collegeDashboard); err != nil {
		return err
	}

	return nil

}

// **************** Change from college branches to respective fields. ***********************

// Gets all campus details.
func (service *CollegeDashboardService) getAllCampusDetails(collegeDashboard *dashboard.CollegeDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of campus drives.
	if err := service.Repository.GetCount(uow, colg.College{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.CampusDetails.TotalCampusDrives = uint(totalCount)
	return nil
}

// Gets all seminar details
func (service *CollegeDashboardService) getAllSeminarDetails(collegeDashboard *dashboard.CollegeDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of seminars.
	if err := service.Repository.GetCount(uow, colg.College{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.SeminarDetails.TotalSeminars = uint(totalCount)
	return nil
}
