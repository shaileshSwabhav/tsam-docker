package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/repository"
)

// SalesPersonDashboardService provides all details to be shown on salesPerson dashboard.
type SalesPersonDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewSalesPersonDashboardService returns new instance of SalesPersonDashboardService.
func NewSalesPersonDashboardService(db *gorm.DB, repository repository.Repository) *SalesPersonDashboardService {
	return &SalesPersonDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetSalesPersonDashboardDetails gets all details required for SalesPersonDashboard.
func (service *SalesPersonDashboardService) GetSalesPersonDashboardDetails(tenantID uuid.UUID, salesPersonDashboard *dashboard.SalesPersonDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(service.DB, true)

	// // Get count of total salesPeople.
	// if err := service.Repository.GetCount(uow, model.User{}, &totalCount,
	// 	repository.Filter("")); err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	// }
	salesPersonDashboard.TotalSalesPeople = 786
	salesPersonDashboard.EnquiriesAssigned = 786
	salesPersonDashboard.EnquiriesConverted = 786
	// salesPersonDashboard.EnquiriesNotAssigned = 786
	salesPersonDashboard.EnquiriesNotHandled = 786
	salesPersonDashboard.TalentsApproached = 786
	salesPersonDashboard.JoinedBatches = 786
	salesPersonDashboard.CampusDrivesCompleted = 786
	salesPersonDashboard.SeminarsCompleted = 786

	return nil

}
