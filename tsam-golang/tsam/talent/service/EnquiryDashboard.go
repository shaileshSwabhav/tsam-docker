package service

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	talentenquiry "github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/repository"
)

// EnquiryDashboardService provides all details to be shown on enquiry dashboard.
type EnquiryDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewEnquiryDashboardService returns new instance of EnquiryDashboardService.
func NewEnquiryDashboardService(db *gorm.DB, repository repository.Repository) *EnquiryDashboardService {
	return &EnquiryDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetEnquiryDashboardDetails gets all details required for EnquiryDashboard.
func (service *EnquiryDashboardService) GetEnquiryDashboardDetails(enquiryDashboard *dashboard.EnquiryDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	var totalCount int = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total enquiries.
	if err := service.Repository.GetCount(uow, talentenquiry.Enquiry{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.TotalEnquiries = uint(totalCount)

	enquiryDashboard.NewEnquiries = 786
	enquiryDashboard.EnquiriesAssigned = 786
	enquiryDashboard.EnquiriesNotHandled = 786
	enquiryDashboard.EnquiriesNotAssigned = 786
	enquiryDashboard.EnquiriesConverted = 786

	return nil

}
