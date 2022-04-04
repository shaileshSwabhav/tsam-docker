package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/models/talentenquiry"
	"github.com/techlabs/swabhav/tsam/util"

	//"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
)

// AdminDashboardService provides all details to be shown on admin dashboard.
type AdminDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewAdminDashboardService returns new instance of AdminDashboardService.
func NewAdminDashboardService(db *gorm.DB, repository repository.Repository) *AdminDashboardService {
	return &AdminDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetTalentEnquiryDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetTalentEnquiryDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.EnquiryDashboard,
	requestForm url.Values) error {
	if err := service.getTalentEnquiryDashboardDetails(tenantID, dashboard, requestForm); err != nil {
		return err
	}
	return nil
}

// GetSalesPeopleDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetSalesPeopleDashboardDetails(tenantID uuid.UUID, form url.Values,
	dashboard *dashboard.SalesPersonDashboard) error {
	if err := service.getSalesPeopleDashboardDetails(tenantID, form, dashboard); err != nil {
		return err
	}
	return nil
}

// GetTalentDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetTalentDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.TalentDashboard) error {
	if err := service.getTalentDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetFacultyDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetFacultyDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.FacultyDashboard) error {
	if err := service.getFacultyDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetCollegeDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetCollegeDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.CollegeDashboard) error {
	if err := service.getCollegeDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetCompanyDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetCompanyDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.CompanyDashboard) error {
	if err := service.getCompanyDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetCourseDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetCourseDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.CourseDashboard) error {
	if err := service.getCourseDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetTechnologyDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetTechnologyDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.TechnologyDashboard) error {
	if err := service.getTechnologyDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetBatchDashboardDetails gets all details required for AdminDashboard.
func (service *AdminDashboardService) GetBatchDashboardDetails(tenantID uuid.UUID, dashboard *dashboard.BatchDashboard) error {
	if err := service.getBatchDashboardDetails(tenantID, dashboard); err != nil {
		return err
	}
	return nil
}

// GetAdminDashboardDetails gets total values of specific entities.
func (service *AdminDashboardService) GetAdminDashboardDetails(tenantID uuid.UUID, adminDashboard *dashboard.AdminDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Total courses
	err := service.Repository.GetCountForTenant(uow, tenantID, course.Course{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	adminDashboard.TotalCourses = totalCount

	// Total technologies
	err = service.Repository.GetCountForTenant(uow, tenantID, general.Technology{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	adminDashboard.TotalTechnologies = totalCount

	return nil

}

// GetBatchPerformances will return talents having outstanding, good or average score
func (service *AdminDashboardService) GetBatchPerformances(batchPerformace *dashboard.BatchPerformance,
	tenantID uuid.UUID, form url.Values) error {

	err := service.getAllBatchPerformance(batchPerformace, tenantID, form)
	if err != nil {
		return err
	}

	return nil
}

// GetSessionWiseTalentScore returns specified talent's session-wise feedback score
func (service *AdminDashboardService) GetSessionWiseTalentScore(talentSessionFeedbackScore *dashboard.TalentSessionFeedbackScore,
	tenantID, talentID, batchID uuid.UUID, form url.Values) error {

	err := service.getTalentSessionWiseScore(talentSessionFeedbackScore, tenantID, talentID, batchID, form)
	if err != nil {
		return err
	}
	return nil
}

// GetBatchStatusDetails returns details of batch for specified batch status
func (service *AdminDashboardService) GetBatchStatusDetails(batchDetails *[]dashboard.BatchDetails,
	tenantID uuid.UUID, form url.Values) error {

	err := service.getBatchStatusDetails(batchDetails, tenantID, form)
	if err != nil {
		return err
	}
	return nil
}

// GetBatchTalentFeedbackScore returns talents feedback score for specified batch
func (service *AdminDashboardService) GetBatchTalentFeedbackScore(talentFeedbackScore *[]dashboard.TalentFeedbackScore,
	tenantID, batchID uuid.UUID) error {

	err := service.getBatchTalentFeedbackScore(talentFeedbackScore, tenantID, batchID)
	if err != nil {
		return err
	}
	return nil
}

// Sales people dashboard for all salesperson.
func (service *AdminDashboardService) getSalesPeopleDashboardDetails(tenantID uuid.UUID, form url.Values,
	salesPersonDashboard *dashboard.SalesPersonDashboard) error {
	var totalCount uint = 0
	// contains date of 30 days ago.
	date := time.Now().AddDate(0, 0, -30).UTC()
	// last30DaysQuery := repository.Filter("`enquiry_date`>= ?", date)

	var talentSalesPersonQuery, talentEnquirySalesPersonQuery repository.QueryProcessor

	if salesPersonID, ok := form["salesPersonID"]; ok {
		talentSalesPersonQuery = repository.Filter("talents.`sales_person_id` = ?", salesPersonID)
		talentEnquirySalesPersonQuery = repository.Filter("talent_enquiries.`sales_person_id` = ?", salesPersonID)
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Need to check with roles in user table.
	// Active deactive fields( only who are working right now)
	// Also make role names constants #niranjan
	err := service.Repository.GetCount(uow, general.User{}, &totalCount,
		repository.Join("INNER JOIN roles ON users.`role_id` = roles.`id` AND roles.`tenant_id` = users.`tenant_id`"),
		repository.Filter("users.`tenant_id` = ? AND roles.`deleted_at` IS NULL AND "+
			"roles.`role_name`='SalesPerson'", tenantID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.TotalSalesPeople = totalCount

	// Need to add if enquiry is closed condition.( ADD FLAG)
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`enquiry_date`>= ? AND `sales_person_id` IS NOT NULL", date),
		talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	salesPersonDashboard.EnquiriesAssigned = totalCount

	// Need to add enquiry closed flag. #niranjan
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`enquiry_date`>= ? AND `talent_id` IS NOT NULL", date),
		talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.EnquiriesConverted = totalCount

	// // No salesperson is assigned.
	// err = service.Repository.GetCount(uow, talentenquiry.Enquiry{}, &totalCount,
	// 	repository.Filter("`sales_person_id` IS NULL AND `talent_id` IS NULL AND "+
	// 		"talent_enquiries.`enquiry_date`>= ?", date))
	// if err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	// }
	// salesPersonDashboard.EnquiriesNotAssigned = totalCount

	// Salesperson yet to call
	err = service.Repository.GetCount(uow, talentenquiry.Enquiry{}, &totalCount,
		repository.Join("LEFT JOIN talent_enquiry_call_records ON "+
			"talent_enquiries.`id` = talent_enquiry_call_records.`enquiry_id` AND "+
			"talent_enquiries.`tenant_id` = talent_enquiry_call_records.`tenant_id`"),
		repository.Filter("talent_enquiry_call_records.`deleted_at` IS NULL AND "+
			"talent_enquiry_call_records.`id` IS NULL AND talent_enquiries.`tenant_id` = ? AND "+
			"talent_enquiries.`enquiry_date`>= ? AND talent_enquiries.`talent_id` IS NULL", tenantID, date),
		talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.EnquiriesNotHandled = totalCount

	// Question( CALL records any entry) also last 30 days
	err = service.Repository.GetCount(uow, talent.Talent{}, &totalCount,
		repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
			"talents.`tenant_id` = talent_call_records.`tenant_id`"),
		repository.Filter("talents.`tenant_id` = ? AND talent_call_records.`deleted_at` IS NULL AND "+
			"CAST(talent_call_records.`date_time` AS DATE) >= ?", tenantID, date),
		talentSalesPersonQuery, repository.GroupBy("talents.`id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.TalentsApproached = totalCount

	// Question (when a talent joins a batch, assigned SP will have closed+1) LAST 30 days
	// active salesperson (role sp check if mitali case)
	err = service.Repository.GetCount(uow, talent.Talent{}, &totalCount,
		repository.Join("INNER JOIN batch_talents ON talents.`id`=batch_talents.`talent_id` AND "+
			"talents.`tenant_id` = batch_talents.`tenant_id`"),
		// change to date_of_joining #niranjan
		repository.Filter("talents.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND "+
			"batch_talents.`created_at` >= ?", tenantID, date),
		talentSalesPersonQuery, repository.GroupBy("talents.`id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.JoinedBatches = totalCount

	// Define constants.
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("talent_enquiries.`enquiry_type` = ? AND talent_enquiries.`enquiry_date` >= ?",
			"Training", date), talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.TrainingEnquiries = totalCount

	// Define constants.
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("talent_enquiries.`enquiry_type` = ? AND talent_enquiries.`enquiry_date` >= ?",
			"Placement", date), talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.PlacementEnquiries = totalCount

	// Define constants.
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("talent_enquiries.`enquiry_type` = ? AND talent_enquiries.`enquiry_date` >= ?",
			"Training And Placement", date), talentEnquirySalesPersonQuery)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	salesPersonDashboard.TrainingAndPlacementEnquiries = totalCount

	// campus cancelled handle (Add current time on check)**********

	salesPersonDashboard.CampusDrivesCompleted = 786

	// Seminar is pending(same as campus)
	salesPersonDashboard.SeminarsCompleted = 786

	return nil

}

// *********************************** Talent Dashboard ***********************************

// GetTalentDashboardDetails gets all details required for TalentDashboard.
func (service *AdminDashboardService) getTalentDashboardDetails(tenantID uuid.UUID, talentDashboard *dashboard.TalentDashboard) error {

	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total talents.
	err := service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.TotalTalents = totalCount

	// Completed = 5 should be defined #niranjan
	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`academic_year` != ? AND `is_experience`= ?", 5, false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.B2CTalents = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`academic_year` = ? AND `is_experience` = ?", 5, false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.TotalFreshers = totalCount

	// Need to check talentDashboard.TotalExperienced = talentDashboard.TotalTalents- talentDashboard.TotalFreshers
	// Would be a problem with null values
	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`is_experience` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.TotalExperienced = totalCount

	// Need to add swabhav talent flag. How is he/she swabhav talents? (add batch flag change all existing true)
	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`is_swabhav_talent` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.SwabhavTalents = totalCount

	// Need to add swabhav talent flag. How is he/she swabhav talents? (add batch flag change all existing true)
	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`is_swabhav_talent` = ?", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.NonSwabhavTalents = totalCount

	// Need to add interested in foreign flag (add flag)
	err = service.Repository.GetCountForTenant(uow, tenantID, talent.Talent{}, &totalCount,
		repository.Filter("`is_masters_abroad` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	talentDashboard.TotalInterestedInForeign = totalCount

	err = service.Repository.Scan(uow, talentDashboard,
		repository.Model(&talent.Talent{}),
		repository.Select("SUM(talents.`lifetime_value`) as total_lifetime_value"),
		repository.Filter("talents.`tenant_id` = ?", tenantID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	return nil

}

// *********************************** Enquiry Dashboard ***********************************

// GetEnquiryDashboardDetails gets all details required for EnquiryDashboard.
func (service *AdminDashboardService) getTalentEnquiryDashboardDetails(tenantID uuid.UUID, enquiryDashboard *dashboard.EnquiryDashboard,
	requestForm url.Values) error {
	var totalCount uint = 0
	// contains date of 30 days ago.
	date := time.Now().AddDate(0, 0, -30).UTC()

	uow := repository.NewUnitOfWork(service.DB, true)

	var queryProcessor repository.QueryProcessor

	if sourceID, ok := requestForm["sourceID"]; ok {
		queryProcessor = repository.Filter("talent_enquiries.`source_id` = ?", sourceID)
	}

	// Get count of total enquiries.
	err := service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`enquiry_date`>= ?", date), queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.TotalEnquiries = totalCount

	// Date criteria for new enquiries (7 days)
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`enquiry_date` >= ?", time.Now().AddDate(0, 0, -7).UTC()), queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.NewEnquiries = totalCount

	// Need to add if enquiry is closed condition.
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`sales_person_id` IS NOT NULL AND `enquiry_date` >= ?", date), queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.EnquiriesAssigned = totalCount

	// Need to add if enquiry is closed condition. OR NOT
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`sales_person_id` IS NULL AND `talent_id` IS NULL AND `enquiry_date`>= ?", date),
		queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)

	}
	enquiryDashboard.EnquiriesNotAssigned = totalCount

	// Clarification needed here (number of days with no call record?)
	err = service.Repository.GetCount(uow, talentenquiry.Enquiry{}, &totalCount,
		repository.Join("LEFT JOIN talent_enquiry_call_records ON "+
			"talent_enquiries.`id` = talent_enquiry_call_records.`enquiry_id` AND "+
			"talent_enquiries.`tenant_id` = talent_enquiry_call_records.`tenant_id`"),
		repository.Filter("talent_enquiries.`tenant_id` = ? AND talent_enquiry_call_records.`deleted_at` IS NULL AND "+
			"talent_enquiry_call_records.`id` IS NULL AND talent_enquiries.`enquiry_date`>= ? "+
			"AND talent_enquiries.`talent_id` IS NULL", tenantID, date), queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.EnquiriesNotHandled = totalCount

	// Need to add flag
	err = service.Repository.GetCountForTenant(uow, tenantID, talentenquiry.Enquiry{}, &totalCount,
		repository.Filter("`enquiry_date`>= ? AND `talent_id` IS NOT NULL", date), queryProcessor)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	enquiryDashboard.EnquiriesConverted = totalCount

	return nil

}

// GetEnquirySourceCount will return count of all the enquiries from a specific source
func (service *AdminDashboardService) GetEnquirySourceCount(tenantID uuid.UUID, sources *[]dashboard.EnquirySource,
	requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return nil
	}

	sourcesRequired := []string{"inst", "fb", "link", "nak", "itshal", "ind", "web"}

	tempSources := []general.Source{}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, &tempSources,
		repository.Filter("`name` IN (?)", sourcesRequired), repository.OrderBy("`name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	var queryProcessor repository.QueryProcessor

	if salesPersonID, ok := requestForm["salesPersonID"]; ok {
		queryProcessor = repository.Filter("talent_enquiries.`sales_person_id` = ?", salesPersonID)
	}

	for index := range tempSources {

		tempEnquirySource := dashboard.EnquirySource{}
		tempEnquirySource.SourceID = tempSources[index].ID
		tempEnquirySource.Description = tempSources[index].Description

		err = service.Repository.GetCount(uow, &tempEnquirySource, &tempEnquirySource.EnquiryCount,
			repository.Table("sources"),
			repository.Join("INNER JOIN talent_enquiries ON sources.`id` = talent_enquiries.`source_id` AND "+
				"sources.`tenant_id` = talent_enquiries.`tenant_id`"), repository.Filter("sources.`id` = ?", tempEnquirySource.SourceID),
			repository.Filter("sources.`tenant_id` = ?", tenantID), repository.Filter("talent_enquiries.`deleted_at` IS NULL"),
			queryProcessor)
		if err != nil {
			uow.RollBack()
			return err
		}

		*sources = append(*sources, tempEnquirySource)
	}

	return nil
}

// *********************************** Faculty Dashboard ***********************************

// GetFacultyDashboardDetails gets all details required for FacultyDashboard.
func (service *AdminDashboardService) getFacultyDashboardDetails(tenantID uuid.UUID, facultyDashboard *dashboard.FacultyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total faculty.
	err := service.Repository.GetCountForTenant(uow, tenantID, faculty.Faculty{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.TotalFaculty = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, faculty.Faculty{}, &totalCount,
		repository.Filter("`is_full_time` = ?", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.PartTime = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, faculty.Faculty{}, &totalCount,
		repository.Filter("`is_full_time` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.FullTime = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, faculty.Faculty{}, &totalCount,
		repository.Filter("`is_active` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.Active = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, faculty.Faculty{}, &totalCount,
		repository.Filter("`is_active` = ?", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.InActive = totalCount

	if err := service.getAllCompletedBatches(tenantID, facultyDashboard); err != nil {
		return err
	}

	if err := service.getAllLiveBatches(tenantID, facultyDashboard); err != nil {
		return err
	}

	if err := service.getAllUpcomingBatches(tenantID, facultyDashboard); err != nil {
		return err
	}

	return nil

}

// Gets all fresher requirements.
func (service *AdminDashboardService) getAllLiveBatches(tenantID uuid.UUID, facultyDashboard *dashboard.FacultyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Can a batch have a end date and be active? there are a few batches.
	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ?", "Ongoing"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.LiveBatches.TotalBatches = totalCount

	// Can 1 talent join two live batches? $niranjan
	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN `batch_talents` ON batches.`id` = `batch_talents`.`batch_id` AND "+
			"batches.`tenant_id` = `batch_talents`.`tenant_id`"),
		repository.Filter("batches.`tenant_id` = ? AND `batch_talents`.`deleted_at` IS NULL AND "+
			"batches.`status` = ?", tenantID, "Ongoing"),
		repository.GroupBy("`batch_talents`.`talent_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	facultyDashboard.LiveBatches.TotalTalentsJoined = totalCount
	return nil
}

// Gets all live requirements.
func (service *AdminDashboardService) getAllUpcomingBatches(tenantID uuid.UUID, facultyDashboard *dashboard.FacultyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		// time.Now().Format("2006-01-02")
		repository.Filter("`batches`.`status` = ?", "Upcoming"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.UpcomingBatches.TotalBatches = totalCount

	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN `batch_talents` ON batches.`id` = `batch_talents`.`batch_id` AND "+
			"batches.`tenant_id` = `batch_talents`.`tenant_id`"),
		repository.Filter("batches.`tenant_id` = ? AND `batch_talents`.`deleted_at` IS NULL AND "+
			"batches.`status` = ?", tenantID, "Upcoming"),
		repository.GroupBy("`batch_talents`.`talent_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.UpcomingBatches.TotalTalentsJoined = totalCount
	return nil
}

// Gets all live requirements.
func (service *AdminDashboardService) getAllCompletedBatches(tenantID uuid.UUID, facultyDashboard *dashboard.FacultyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ?", "Finished"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.BatchesCompleted.TotalBatches = totalCount

	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id` AND "+
			"batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Filter("batches.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND "+
			"batches.`status` = ?", tenantID, "Finished"),
		repository.GroupBy("`batch_talents`.`talent_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	facultyDashboard.BatchesCompleted.TotalTalentsJoined = totalCount
	return nil
}

// *********************************** College Dashboard ***********************************

// GetCollegeDashboardDetails gets all details required for CollegeDashboard.
func (service *AdminDashboardService) getCollegeDashboardDetails(tenantID uuid.UUID, collegeDashboard *dashboard.CollegeDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total colleges.
	if err := service.Repository.GetCountForTenant(uow, tenantID, college.Branch{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.TotalColleges = totalCount

	// Optional
	collegeDashboard.TotalActiveColleges = 786

	if err := service.getAllCampusDetails(tenantID, collegeDashboard); err != nil {
		return err
	}

	if err := service.getAllSeminarDetails(tenantID, collegeDashboard); err != nil {
		return err
	}

	return nil

}

// **************** Change from college branches to respective fields. ***********************

// Gets all campus details.
func (service *AdminDashboardService) getAllCampusDetails(tenantID uuid.UUID, collegeDashboard *dashboard.CollegeDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of campus drives.
	if err := service.Repository.GetCountForTenant(uow, tenantID, college.College{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.CampusDetails.TotalCampusDrives = 786

	// Pending
	collegeDashboard.CampusDetails.AllTimeRegisteredTalents = 786
	collegeDashboard.CampusDetails.Ongoing = 786
	collegeDashboard.CampusDetails.Upcoming = 786

	return nil
}

// Gets all seminar details
func (service *AdminDashboardService) getAllSeminarDetails(tenantID uuid.UUID, collegeDashboard *dashboard.CollegeDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of seminars.
	if err := service.Repository.GetCountForTenant(uow, tenantID, college.College{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	collegeDashboard.SeminarDetails.TotalSeminars = 786

	// Pending
	collegeDashboard.SeminarDetails.Upcoming = 786
	collegeDashboard.SeminarDetails.AllTimeRegisteredTalents = 786
	collegeDashboard.SeminarDetails.Ongoing = 786

	return nil
}

// *********************************** Company Dashboard ***********************************

// GetCompanyDashboardDetails gets all details required for CompanyDashboard.
func (service *AdminDashboardService) getCompanyDashboardDetails(tenantID uuid.UUID, companyDashboard *dashboard.CompanyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total companies.
	if err := service.Repository.GetCountForTenant(uow, tenantID, company.Branch{}, &totalCount); err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.TotalCompanies = totalCount

	if err := service.getAllExperienceRequirements(tenantID, companyDashboard); err != nil {
		return err
	}

	if err := service.getAllFresherRequirements(tenantID, companyDashboard); err != nil {
		return err
	}

	if err := service.getAllLiveRequirements(tenantID, companyDashboard); err != nil {
		return err
	}

	return nil

}

// **************** Change from company branches to respective fields. ***********************

// Gets all fresher requirements.
func (service *AdminDashboardService) getAllFresherRequirements(tenantID uuid.UUID, companyDashboard *dashboard.CompanyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total fresher requirements.  ***************ACTIVE?***************
	err := service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Filter("(`minimum_experience` IS NULL OR `minimum_experience` = ?)", 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.FresherRequirements.TotalRequirements = totalCount

	// count of talents from allocated talents?
	err = service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Join("INNER JOIN company_requirements_talents ON "+
			"company_requirements.`id` = company_requirements_talents.`requirement_id`"),
		repository.Filter("(company_requirements.`minimum_experience` IS NULL OR "+
			"company_requirements.`minimum_experience` = ?)", 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.FresherRequirements.PlacementsInProcess = totalCount

	// count of talets from allocated talents- total requirement?
	err = service.Repository.Scan(uow, &companyDashboard.FresherRequirements,
		repository.Model(&company.Requirement{}),
		repository.Select("SUM(company_requirements.`vacancy`) AS total_talents_required"),
		repository.Filter("company_requirements.`tenant_id` = ? AND "+
			"(company_requirements.`minimum_experience` IS NULL OR company_requirements.`minimum_experience` = ?)",
			tenantID, 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	totalTalentsRequried := companyDashboard.FresherRequirements.TotalTalentsRequired
	placementInProcess := companyDashboard.FresherRequirements.PlacementsInProcess

	companyDashboard.FresherRequirements.TalentsNeeded = totalTalentsRequried - placementInProcess
	return nil
}

// Gets all live requirements.
func (service *AdminDashboardService) getAllLiveRequirements(tenantID uuid.UUID, companyDashboard *dashboard.CompanyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Filter("`is_active` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.LiveRequirements.TotalRequirements = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Join("INNER JOIN company_requirements_talents ON "+
			"company_requirements.`id` = company_requirements_talents.`requirement_id`"),
		repository.Filter("company_requirements.`is_active` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.LiveRequirements.PlacementsInProcess = totalCount

	err = service.Repository.Scan(uow, &companyDashboard.LiveRequirements,
		repository.Model(&company.Requirement{}),
		repository.Select("SUM(company_requirements.`vacancy`) AS total_talents_required"),
		repository.Filter("company_requirements.`tenant_id` = ? AND (`is_active` = ?)", tenantID, true))
	if err != nil {
		return err
	}
	totalTalentsRequried := companyDashboard.LiveRequirements.TotalTalentsRequired
	placementInProcess := companyDashboard.LiveRequirements.PlacementsInProcess

	companyDashboard.LiveRequirements.TalentsNeeded = totalTalentsRequried - placementInProcess
	return nil
}

// Gets all live requirements.
func (service *AdminDashboardService) getAllExperienceRequirements(tenantID uuid.UUID, companyDashboard *dashboard.CompanyDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Filter("(`minimum_experience` IS NOT NULL AND `minimum_experience` != ?)", 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.ExperienceRequirements.TotalRequirements = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, company.Requirement{}, &totalCount,
		repository.Join("INNER JOIN company_requirements_talents ON "+
			"company_requirements.`id` = company_requirements_talents.`requirement_id`"),
		repository.Filter("(company_requirements.`minimum_experience` IS NOT NULL AND "+
			"company_requirements.`minimum_experience` != ?)", 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	companyDashboard.ExperienceRequirements.PlacementsInProcess = totalCount

	err = service.Repository.Scan(uow, &companyDashboard.ExperienceRequirements,
		repository.Model(&company.Requirement{}),
		repository.Select("SUM(company_requirements.`vacancy`) AS total_talents_required"),
		repository.Filter("company_requirements.`tenant_id` = ? AND "+
			"(company_requirements.`minimum_experience` IS NOT NULL AND company_requirements.`minimum_experience` != ?)",
			tenantID, 0))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	totalTalentsRequried := companyDashboard.ExperienceRequirements.TotalTalentsRequired
	placementInProcess := companyDashboard.ExperienceRequirements.PlacementsInProcess
	companyDashboard.ExperienceRequirements.TalentsNeeded = totalTalentsRequried - placementInProcess
	return nil
}

// *********************************** Batch-Course Dashboard ***********************************
func (service *AdminDashboardService) getBatchDashboardDetails(tenantID uuid.UUID, batchDashboard *dashboard.BatchDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Live category
	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Ongoing", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Live.B2BBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Ongoing", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Live.B2CBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ?", "Ongoing"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Live.AllBatches = totalCount

	err = service.Repository.GetCount(uow, batch.MappedTalent{}, &totalCount,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id`"),
		repository.Filter("batches.`status` = ?", "Ongoing"), repository.Filter("batch_talents.`tenant_id`=?", tenantID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Live.TotalStudents = totalCount

	// Completed category
	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Finished", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Completed.B2BBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Finished", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Completed.B2CBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ?", "Finished"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Completed.AllBatches = totalCount

	err = service.Repository.GetCount(uow, batch.MappedTalent{}, &totalCount,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id`"),
		repository.Filter("batches.`status` = ?", "Finished"), repository.Filter("batch_talents.`tenant_id`=?", tenantID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Completed.TotalStudents = totalCount

	// Upcoming category
	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Upcoming", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Upcoming.B2BBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `is_b2b` = ?", "Upcoming", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Upcoming.B2CBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ?", "Upcoming"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Upcoming.AllBatches = totalCount

	err = service.Repository.GetCount(uow, batch.MappedTalent{}, &totalCount,
		repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id`"),
		repository.Filter("batches.`status` = ?", "Upcoming"), repository.Filter("batch_talents.`tenant_id`=?", tenantID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Upcoming.TotalStudents = totalCount

	// Total batches
	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`is_b2b` = ?", true))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Total.B2BBatches = totalCount

	// total category
	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`is_b2b` = ?", false))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Total.B2CBatches = totalCount

	err = service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Total.AllBatches = totalCount

	err = service.Repository.GetCount(uow, batch.MappedTalent{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	batchDashboard.Total.TotalStudents = totalCount

	return nil
}

// GetCourseBatchDashboardDetails gets all details required for CourseBatchDashboard.
func (service *AdminDashboardService) getCourseDashboardDetails(tenantID uuid.UUID, courseDashboard *dashboard.CourseDashboard) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get count of total courses.
	err := service.Repository.GetCountForTenant(uow, tenantID, general.Course{}, &totalCount)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDashboard.TotalCourses = totalCount

	// Get courses data of group A.
	err = service.getCourseGroupsData("GroupA", "Normal", courseDashboard, tenantID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	// Get courses data of group B.
	err = service.getCourseGroupsData("GroupB", "Hot", courseDashboard, tenantID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	// Get courses data of group C.
	err = service.getCourseGroupsData("GroupC", "Very Hot", courseDashboard, tenantID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	// Get courses data of group D.
	err = service.getCourseGroupsData("GroupD", "Hot & Rare", courseDashboard, tenantID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	return nil

}

// Gets all fresher requirements.
func (service *AdminDashboardService) getCourseGroupsData(groupName, courseType string,
	courseDashboard *dashboard.CourseDashboard, tenantID uuid.UUID) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)
	var courseGroup dashboard.CourseGroup

	courseGroup.GroupName = groupName

	err := service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id`"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"courses.`course_type` = ?", tenantID, courseType))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.AllBatches.TotalBatches = totalCount

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id` "+
			"INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id` AND "+
			"batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"batch_talents.`deleted_at` IS NULL AND courses.`course_type` = ?", tenantID, courseType))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.AllBatches.TotalTalentsJoined = totalCount

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id`"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"courses.`course_type` = ? AND batches.`status` = ? ", tenantID, courseType, "Ongoing"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.LiveBatches = totalCount

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id`"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"courses.`course_type` = ? AND batches.`status` = ? ", tenantID, courseType, "Upcoming"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.UpcomingBatches.TotalBatches = totalCount

	err = service.Repository.Scan(uow, &courseGroup.UpcomingBatches,
		repository.Model(&course.Course{}),
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id`"),
		repository.Select("SUM(batches.`total_intake`) AS total_talent_intake"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL "+
			"AND batches.`status` = ? AND courses.`course_type` = ? ", tenantID, "Upcoming", courseType))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id` "+
			"INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id` AND "+
			"courses.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"batch_talents.`deleted_at` IS NULL AND batches.`status` = ? AND courses.`course_type` = ?",
			tenantID, "Upcoming", courseType))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.UpcomingBatches.TotalTalentsJoined = totalCount
	intake := courseGroup.UpcomingBatches.TotalTalentIntake
	courseGroup.UpcomingBatches.RequiredTalents = intake - totalCount
	// Not sure if this is needed.
	if intake == 0 {
		courseGroup.UpcomingBatches.RequiredTalents = 0
	}

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id` "+
			"INNER JOIN batch_session_topic ON batches.`id` = batch_session_topic.`batch_id` AND "+
			"batches.`tenant_id` = batch_session_topic.`tenant_id` "),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"batch_session_topic.`deleted_at` IS NULL AND "+
			"courses.`course_type` = ?", tenantID, courseType),
		repository.GroupBy("batch_session_topics.`batch_session_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.Sessions.TotalSessions = totalCount

	err = service.Repository.GetCount(uow, course.Course{}, &totalCount,
		repository.Join("INNER JOIN batches ON courses.`id` = batches.`course_id` AND "+
			"courses.`tenant_id` = batches.`tenant_id` "+
			"INNER JOIN batch_session_topic ON batches.`id` = batch_session_topic.`batch_id` AND "+
			"batches.`tenant_id` = batch_session_topic.`tenant_id` "),
		repository.Filter("courses.`tenant_id` = ? AND batches.`deleted_at` IS NULL AND "+
			"batch_session_topic.`deleted_at` IS NULL AND courses.`course_type` = ? AND "+
			"batch_session_topic.`is_completed` = ?", tenantID, courseType, false),
		repository.GroupBy("batch_session_topics.`batch_session_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseGroup.Sessions.RemainingSessions = totalCount

	err = service.getAllCourseDetails(courseType, &courseGroup, tenantID)
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	courseDashboard.CourseGroups = append(courseDashboard.CourseGroups, courseGroup)
	return nil
}

// Gets all fresher requirements.
func (service *AdminDashboardService) getAllCourseDetails(courseType string,
	courseGroup *dashboard.CourseGroup, tenantID uuid.UUID) error {
	allCourses := []course.Course{}

	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, &allCourses, "`name`",
		repository.Select("`id`, `name`, `course_level`"),
		repository.Filter("`course_type` = ?", courseType))
	if err != nil {
		return err
	}

	for _, course := range allCourses {
		courseData := &dashboard.CourseData{}
		courseData.CourseName = course.Name
		courseData.CourseLevel = course.CourseLevel
		err = service.getAllBatchesForCourse(course.ID, tenantID, courseData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		err = service.getAllLiveBatchesForCourse(course.ID, tenantID, courseData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		err = service.getAllUpcomingBatchesForCourse(course.ID, tenantID, courseData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		err = service.getAllSessionForCourse(course.ID, tenantID, courseData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}

		courseGroup.CoursesData = append(courseGroup.CoursesData, courseData)
	}
	return nil
}

func (service *AdminDashboardService) getAllBatchesForCourse(courseID, tenantID uuid.UUID,
	courseDetail *dashboard.CourseData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	// all batches for course.
	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`course_id` = ?", courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.AllBatches.TotalBatches = totalCount

	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id` AND "+
			"batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Filter("batches.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND "+
			"batches.`course_id` = ?", tenantID, courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.AllBatches.TotalTalentsJoined = totalCount
	return nil
}

func (service *AdminDashboardService) getAllLiveBatchesForCourse(courseID, tenantID uuid.UUID,
	courseDetail *dashboard.CourseData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `course_id` = ?", "Ongoing", courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.LiveBatches = totalCount
	return nil
}

func (service *AdminDashboardService) getAllUpcomingBatchesForCourse(courseID, tenantID uuid.UUID,
	courseDetail *dashboard.CourseData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, &totalCount,
		repository.Filter("`status` = ? AND `course_id` = ?", "Upcoming", courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.UpcomingBatches.TotalBatches = totalCount

	err = service.Repository.Scan(uow, &courseDetail.UpcomingBatches,
		repository.Model(&batch.Batch{}),
		repository.Select("SUM(`total_intake`) AS total_talent_intake"),
		repository.Filter("`tenant_id` = ? AND `status` = ? AND `course_id` = ?", tenantID, "Upcoming", courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}

	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN batch_talents ON batches.`id` = batch_talents.`batch_id` AND "+
			"batches.`tenant_id` = batch_talents.`tenant_id`"),
		repository.Filter("batches.`tenant_id` = ? AND batch_talents.`deleted_at` IS NULL AND "+
			"batches.`status` = ? AND batches.`course_id` = ?", tenantID, "Upcoming", courseID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.UpcomingBatches.TotalTalentsJoined = totalCount
	courseDetail.UpcomingBatches.RequiredTalents = courseDetail.UpcomingBatches.TotalTalentIntake - totalCount
	return nil
}

func (service *AdminDashboardService) getAllSessionForCourse(courseID, tenantID uuid.UUID,
	courseDetail *dashboard.CourseData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN batch_session_topic ON batches.`id` = batch_session_topic.`batch_id` AND "+
			"batches.`tenant_id` = batch_session_topic.`tenant_id` "),
		repository.Filter("batches.`tenant_id` = ? AND batch_session_topic.`deleted_at` IS NULL AND "+
			"batches.`course_id` = ?", tenantID, courseID), repository.GroupBy("batch_session_topic.`batch_session_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.Sessions.TotalSessions = totalCount

	err = service.Repository.GetCount(uow, batch.Batch{}, &totalCount,
		repository.Join("INNER JOIN batch_session_topic ON batches.`id` = batch_session_topic.`batch_id` AND "+
			"batches.`tenant_id` = batch_session_topic.`tenant_id` "),
		repository.Filter("batches.`tenant_id` = ? AND batch_session_topic.`deleted_at` IS NULL AND "+
			"batches.`course_id` = ? AND batch_session_topic.`is_completed` = ?", tenantID, courseID, false),
		repository.GroupBy("batch_session_topic.`batch_session_id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	courseDetail.Sessions.RemainingSessions = totalCount
	return nil

}

// GetTechnologyDashboardDetails gets all details required for TechnologyDashboard.
func (service *AdminDashboardService) getTechnologyDashboardDetails(tenantID uuid.UUID, technologyDashboard *dashboard.TechnologyDashboard) error {
	var allTechnologies []general.Technology

	uow := repository.NewUnitOfWork(service.DB, true)
	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, &allTechnologies, "`language`")
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	technologyDashboard.TotalCount = uint(len(allTechnologies))

	for _, technology := range allTechnologies {
		technologyData := &dashboard.TechnologyData{}
		technologyData.Name = technology.Language
		err = service.getFresherCountForTechnology(tenantID, technology.ID, technologyData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		err = service.getExperiencedCountForTechnology(tenantID, technology.ID, technologyData)
		if err != nil {
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		technologyData.TotalTalents = technologyData.TotalFreshers + technologyData.TotalExperienced

		// Update other total fields of the main technology dashboard.
		technologyDashboard.TotalFreshers = technologyDashboard.TotalFreshers + technologyData.TotalFreshers
		technologyDashboard.TotalExperienced = technologyDashboard.TotalExperienced + technologyData.TotalExperienced
		technologyDashboard.TotalTalents = technologyDashboard.TotalTalents + technologyData.TotalTalents

		if technologyData.TotalTalents != 0 {
			technologyDashboard.TechnologiesData = append(technologyDashboard.TechnologiesData, technologyData)
		}
	}

	return nil
}

func (service *AdminDashboardService) getFresherCountForTechnology(tenantID, technologyID uuid.UUID,
	technologyData *dashboard.TechnologyData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCount(uow, talent.Talent{}, &totalCount,
		repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
		repository.Filter("talents.`tenant_id` = ? AND talents.`is_experience` = ? AND "+
			"talents_technologies.`technology_id` = ?", tenantID, false, technologyID))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	technologyData.TotalFreshers = totalCount
	return nil
}

func (service *AdminDashboardService) getExperiencedCountForTechnology(tenantID, technologyID uuid.UUID,
	technologyData *dashboard.TechnologyData) error {
	var totalCount uint = 0
	uow := repository.NewUnitOfWork(service.DB, true)

	err := service.Repository.GetCount(uow, talent.Talent{}, &totalCount,
		repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` AND "+
			"talents.`tenant_id` = talent_experiences.`tenant_id`"+
			"INNER JOIN talent_experiences_technologies ON "+
			"talent_experiences.`id`= talent_experiences_technologies.`experience_id`"),
		repository.Filter("talents.`tenant_id` = ? AND talent_experiences.`deleted_at` IS NULL AND "+
			"talents.`is_experience` = ? AND talent_experiences_technologies.`technology_id` = ? ",
			tenantID, true, technologyID),
		repository.GroupBy("talents.`id`"))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	}
	technologyData.TotalExperienced = totalCount

	return nil
}

// *********************************** Faculty-Talent-Batch-Feedback Dashboard ***********************************

// getAllBatchPerformance will return talents having outstanding, good or average score
func (service *AdminDashboardService) getAllBatchPerformance(batchPerformance *dashboard.BatchPerformance,
	tenantID uuid.UUID, form url.Values) error {

	batchFeedbackScore := []dashboard.TalentFeedbackScore{}
	feedbackQuestionGroups := []general.FeedbackQuestionGroup{}
	// batchFeedbackScore := []dashboard.BatchFeedbackScore{}

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// get all talents for whom feedback is given
	selectTalent := "talents.`first_name`, talents.`last_name`, talents.`id` AS talent_id"
	selectBatch := "batches.`batch_name`, batches.`id` AS batch_id"
	selectScore := "((SUM(feedback_options.`key`) / SUM(feedback_questions.`max_score`)) * 10) AS score"

	err = service.Repository.GetAll(uow, &batchFeedbackScore,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select(selectTalent+","+selectBatch+","+selectScore),
		repository.Join("INNER JOIN batches ON faculty_talent_batch_session_feedback.`batch_id` = batches.`id`"+
			" AND batches.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN talents ON faculty_talent_batch_session_feedback.`talent_id` = talents.`id`"+
			" AND talents.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id`"+
			" AND feedback_options.`tenant_id` = faculty_talent_batch_session_feedback.`tenant_id`"),
		service.addBatchFeedbackSearchQueries(form),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=? AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
		repository.Filter("feedback_questions.`has_options` = true"),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("batches.`deleted_at` IS NULL AND talents.`deleted_at` IS NULL"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"),
		repository.GroupBy("faculty_talent_batch_session_feedback.`talent_id`"),
		repository.OrderBy("score DESC"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Get feedback question keyword group wise
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &feedbackQuestionGroups, "`order`",
		repository.Table("feedback_question_groups"), repository.Select("`id`, `group_name`"),
		repository.Filter("`type` = ?", "Faculty_Session_Feedback"))
	if err != nil {
		uow.RollBack()
		return err
	}

	tempGroupwiseKeywords := make([]dashboard.GroupWiseKeywordName, len(feedbackQuestionGroups))
	for index := range feedbackQuestionGroups {

		tempGroupwiseKeywords[index].GroupName = (feedbackQuestionGroups)[index].GroupName

		err = service.getKeywords(uow, &tempGroupwiseKeywords[index].Keywords,
			tenantID, (feedbackQuestionGroups)[index].ID)
		if err != nil {
			return err
		}
	}
	batchPerformance.KeywordNames = tempGroupwiseKeywords

	// get group wise batch feedback score
	for i := range batchFeedbackScore {
		tempFeedbackGroup := make([]dashboard.GroupScore, len(feedbackQuestionGroups))

		for j := range feedbackQuestionGroups {

			tempFeedbackGroup[j].GroupID = feedbackQuestionGroups[j].ID
			tempFeedbackGroup[j].GroupName = feedbackQuestionGroups[j].GroupName

			err = service.getGroupWiseFeedbackScore(&tempFeedbackGroup[j].FeedbackScore, tenantID,
				batchFeedbackScore[i].TalentID, batchFeedbackScore[i].BatchID, feedbackQuestionGroups[j].ID, uow)
			if err != nil {
				uow.RollBack()
				return err
			}

		}
		// err = service.getBatchFeedbackAverageScore(uow, tenantID, batchFeedbackScore[i].BatchID,
		// 	(batchFeedbackScore)[i].TalentID, &(batchFeedbackScore)[i])
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }
		batchFeedbackScore[i].FeedbackGroup = append(batchFeedbackScore[i].FeedbackGroup, tempFeedbackGroup...)
	}

	err = service.getInterviewRating(uow, &batchFeedbackScore, tenantID)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.getTotalBatchCount(uow, &batchPerformance.TotalBatches, tenantID, form)
	if err != nil {
		uow.RollBack()
		return err
	}

	for _, feedback := range batchFeedbackScore {

		if feedback.Score > 8.0 {
			batchPerformance.Outstanding++
			batchPerformance.OutstandingTalent = append(batchPerformance.OutstandingTalent, feedback)
			continue
		}
		if feedback.Score < 8.0 && feedback.Score > 5.0 {
			batchPerformance.Good++
			batchPerformance.GoodTalent = append(batchPerformance.GoodTalent, feedback)
			continue
		}
		if feedback.Score < 5.0 {
			batchPerformance.Average++
			batchPerformance.AverageTalent = append(batchPerformance.AverageTalent, feedback)
			continue
		}
	}

	uow.Commit()
	return nil
}

// getTalentSessionWiseScore returns specified talent's session-wise feedback score
func (service *AdminDashboardService) getTalentSessionWiseScore(talentSessionFeedbackScore *dashboard.TalentSessionFeedbackScore,
	tenantID, talentID, batchID uuid.UUID, requestForm url.Values) error {

	sessionFeedbackScores := []dashboard.SessionKeywordFeedback{}
	feedbackQuestionGroups := []general.FeedbackQuestionGroup{}
	tempTalent := talent.Talent{}
	tempBatch := batch.Batch{}

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetRecordForTenant(uow, tenantID, &tempTalent,
		repository.Filter("`id`=?", talentID),
		repository.Select([]string{"`id`, `first_name`", "`last_name`", "`personality_type`", "`talent_type`"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.GetRecordForTenant(uow, tenantID, &tempBatch,
		repository.Filter("`id`=?", batchID), repository.Select([]string{"`batch_name`", "`id`"}))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &feedbackQuestionGroups, "`order`",
		repository.Table("feedback_question_groups"), repository.Select("`id`, `group_name`"),
		repository.Filter("`type` = ?", "Faculty_Session_Feedback"))
	if err != nil {
		uow.RollBack()
		return err
	}

	selectCourse := "module_topics.`name` AS session_name, module_topics.`order`"
	selectBatchSession := "faculty_talent_batch_session_feedback.`batch_session_topic_id` AS batch_topic_id"
	selectSession := "batch_session_topics.`initial_date` as session_date"

	// get all sessions of specified batch and talent
	err = service.Repository.GetAll(uow, &sessionFeedbackScores,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select(selectCourse+","+selectBatchSession+","+selectSession),
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN batch_sessions ON batch_sessions.`id` = faculty_talent_batch_session_feedback.`batch_session_id`"+
			" AND faculty_talent_batch_session_feedback.`tenant_id`=batch_sessions.`tenant_id`"),
		repository.Join("INNER JOIN batch_session_topics ON batch_sessions.`id` = batch_session_topics.`batch_session_id`"+
			" AND batch_session_topics.`tenant_id`=batch_sessions.`tenant_id`"),
		repository.Filter("feedback_questions.`has_options` = true"),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("faculty_talent_batch_session_feedback.`batch_id`=? AND "+
			"faculty_talent_batch_session_feedback.`talent_id`=?", batchID, talentID),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=? AND"+
			" faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
		repository.Filter("feedback_questions.`deleted_at` IS NULL AND batch_session_topics.`deleted_at` IS NULL"),
		repository.Filter("module_topics.`deleted_at` IS NULL"),
		repository.GroupBy("faculty_talent_batch_session_feedback.`talent_id`, "+
			"faculty_talent_batch_session_feedback.`batch_session_topic_id`", "batch_session_topics.`batch_session_id`"),
		repository.OrderBy("batch_session_topics.`order`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// get group wise keywords
	tempGroupwiseKeywords := make([]dashboard.GroupWiseKeywordName, len(feedbackQuestionGroups))

	for index := range feedbackQuestionGroups {

		// get groupwise keywords
		err = service.getKeywords(uow, &tempGroupwiseKeywords[index].Keywords, tenantID, feedbackQuestionGroups[index].ID)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	for i := range sessionFeedbackScores {

		sessionFeedbackScores[i].Score = 0
		tempFeedbackGroup := make([]dashboard.GroupScore, len(feedbackQuestionGroups))

		for j := range feedbackQuestionGroups {
			tempFeedbackGroup[j].GroupID = feedbackQuestionGroups[j].ID
			tempFeedbackGroup[j].GroupName = feedbackQuestionGroups[j].GroupName
			tempGroupwiseKeywords[j].GroupName = feedbackQuestionGroups[j].GroupName

			err = service.getGroupWiseSessionScore(uow, tenantID, batchID, talentID,
				(sessionFeedbackScores)[i].BatchSessionID, feedbackQuestionGroups[j].ID, &tempFeedbackGroup[j].FeedbackScore)
			if err != nil {
				uow.RollBack()
				return err
			}

		}

		err = service.getTalentFeelingDetails(uow, tenantID, talentID, batchID, (sessionFeedbackScores)[i].BatchSessionID,
			&sessionFeedbackScores[i])
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.getAverageSessionScore(uow, tenantID, batchID, talentID, (sessionFeedbackScores)[i].BatchSessionID,
			&(sessionFeedbackScores)[i])
		if err != nil {
			uow.RollBack()
			return err
		}

		sessionFeedbackScores[i].FeedbackGroup = append(sessionFeedbackScores[i].FeedbackGroup, tempFeedbackGroup...)
	}

	talentSessionFeedbackScore.TalentID = tempTalent.ID
	talentSessionFeedbackScore.FirstName = tempTalent.FirstName
	talentSessionFeedbackScore.LastName = tempTalent.LastName
	talentSessionFeedbackScore.PersonalityType = tempTalent.PersonalityType
	talentSessionFeedbackScore.TalentType = tempTalent.TalentType
	talentSessionFeedbackScore.BatchName = tempBatch.BatchName
	talentSessionFeedbackScore.BatchID = tempBatch.ID
	talentSessionFeedbackScore.KeywordNames = tempGroupwiseKeywords
	talentSessionFeedbackScore.SessionFeedback = sessionFeedbackScores

	return nil
}

// getBatchTalentFeedbackScore returns talents feedback score for specified batch
func (service *AdminDashboardService) getBatchTalentFeedbackScore(talentFeedbackScore *[]dashboard.TalentFeedbackScore,
	tenantID, batchID uuid.UUID) error {

	feedbackQuestionGroups := []general.FeedbackQuestionGroup{}

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	selectTalent := "talents.`first_name`, talents.`last_name`, talents.`personality_type`, talents.`talent_type`"
	selectBatch := "batches.`batch_name`, batches.`id` AS batch_id"
	selectScore := "talents.`id` AS talent_id"

	err = service.Repository.GetAll(uow, &talentFeedbackScore,
		repository.Table("talents"),
		repository.Select(selectTalent+","+selectBatch+","+selectScore),
		repository.Join("LEFT JOIN batch_talents ON batch_talents.`talent_id` = talents.`id` AND batch_talents.`tenant_id` = talents.`tenant_id`"),
		repository.Join("LEFT JOIN batches ON batches.`id` = batch_talents.`batch_id` AND batch_talents.`tenant_id` = batches.`tenant_id`"),
		repository.Filter("batches.`id`=? ", batchID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL AND batch_talents.`deleted_at` IS NULL"),
		repository.GroupBy("talents.`id`"),
		repository.OrderBy("talents.`first_name`"))

	if err != nil {
		uow.RollBack()
		return err
	}

	// Get feedback question keyword group wise
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &feedbackQuestionGroups, "`order`",
		repository.Table("feedback_question_groups"), repository.Select("`id`, `group_name`"),
		repository.Filter("`type` = ?", "Faculty_Session_Feedback"))
	if err != nil {
		uow.RollBack()
		return err
	}

	// get group wise batch feedback score
	for i := range *talentFeedbackScore {
		tempFeedbackGroup := make([]dashboard.GroupScore, len(feedbackQuestionGroups))

		for j := range feedbackQuestionGroups {
			tempFeedbackGroup[j].GroupID = feedbackQuestionGroups[j].ID
			tempFeedbackGroup[j].GroupName = feedbackQuestionGroups[j].GroupName

			err = service.getGroupWiseFeedbackScore(&tempFeedbackGroup[j].FeedbackScore, tenantID,
				(*talentFeedbackScore)[i].TalentID, batchID, feedbackQuestionGroups[j].ID, uow)
			if err != nil {
				uow.RollBack()
				return err
			}
		}
		err = service.getBatchFeedbackAverageScore(uow, tenantID, batchID,
			(*talentFeedbackScore)[i].TalentID, &(*talentFeedbackScore)[i])
		if err != nil {
			uow.RollBack()
			return err
		}
		(*talentFeedbackScore)[i].FeedbackGroup = append((*talentFeedbackScore)[i].FeedbackGroup, tempFeedbackGroup...)
	}

	err = service.getInterviewRating(uow, talentFeedbackScore, tenantID)
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// getGroupWiseFeedbackScore will calculate average score for specified group
func (service *AdminDashboardService) getGroupWiseFeedbackScore(feedbackKeyword *[]dashboard.FeedbackKeywords,
	tenantID, talentID, batchID, groupID uuid.UUID, uow *repository.UnitOfWork) error {

	exist, err := repository.DoesRecordExist(service.DB, batch.FacultyTalentBatchSessionFeedback{},
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("feedback_questions.`group_id` = ?", groupID),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id` = ? AND "+
			" faculty_talent_batch_session_feedback.`batch_id` = ?", talentID, batchID),
		repository.Filter("feedback_questions.`has_options` = true"),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL"))
	if err != nil {
		return err
	}
	if exist {
		err = service.Repository.GetAll(uow, feedbackKeyword,
			repository.Table("feedback_questions"),
			repository.Select([]string{"feedback_questions.`keyword`," +
				"((SUM(feedback_options.`key`) / SUM(feedback_questions.`max_score`)) * 10) AS keyword_score"}),
			repository.Join("INNER JOIN faculty_talent_batch_session_feedback ON"+
				" faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
				" AND faculty_talent_batch_session_feedback.`tenant_id` = feedback_questions.`tenant_id`"),
			repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id`"+
				" AND faculty_talent_batch_session_feedback.`tenant_id` = feedback_options.`tenant_id`"),
			repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=?", talentID),
			repository.Filter("faculty_talent_batch_session_feedback.`batch_id`=?", batchID),
			repository.Filter("feedback_questions.`has_options` = true"),
			// repository.Filter("feedback_questions.`is_active` = true"),
			repository.Filter("feedback_questions.`group_id` = ?", groupID),
			repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=? AND "+
				"faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
			repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"),
			repository.GroupBy("feedback_questions.`id`"),
			repository.OrderBy("feedback_questions.`order`"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *AdminDashboardService) getBatchFeedbackAverageScore(uow *repository.UnitOfWork, tenantID,
	batchID, talentID uuid.UUID, talentFeedbackScore *dashboard.TalentFeedbackScore) error {

	// check does record exist
	exist, err := repository.DoesRecordExist(service.DB, batch.FacultyTalentBatchSessionFeedback{},
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("feedback_questions.`has_options` = true"),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL"),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id` = ? AND "+
			" faculty_talent_batch_session_feedback.`batch_id` = ?", talentID, batchID))
	if err != nil {
		return err
	}
	if exist {
		err = service.Repository.GetAll(uow, talentFeedbackScore, repository.Table("faculty_talent_batch_session_feedback"),
			repository.Select([]string{"((SUM(feedback_options.`key`) / SUM(feedback_questions.`max_score`)) * 10) AS score"}),
			repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
				" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id`"+
				" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
			repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=? AND faculty_talent_batch_session_feedback.`batch_id`=?",
				talentID, batchID),
			repository.Filter("feedback_questions.`has_options` = true"),
			// repository.Filter("feedback_questions.`is_active` = true"),
			repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?"+
				" AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
			repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"),
			repository.GroupBy("faculty_talent_batch_session_feedback.`talent_id`"))
		if err != nil {
			return err
		}
	}

	return nil
}

// getAverageSessionScore returns overall score for specified score
func (service *AdminDashboardService) getAverageSessionScore(uow *repository.UnitOfWork, tenantID,
	batchID, talentID, sessionID uuid.UUID, sessionFeedbackScore *dashboard.SessionKeywordFeedback) error {

	err := service.Repository.GetAll(uow, sessionFeedbackScore, repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select([]string{"AVG(feedback_options.`key` * 10 / feedback_questions.`max_score`) AS score"}),
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=? AND faculty_talent_batch_session_feedback.`batch_id`=?",
			talentID, batchID),
		repository.Filter("feedback_questions.`has_options` = ? AND "+
			"faculty_talent_batch_session_feedback.`batch_session_id`=?", sessionID),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?"+
			" AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"),
		repository.OrderBy("feedback_questions.`order`"))
	if err != nil {
		return err
	}

	return nil
}

// getGroupWiseSessionScore will retrun session score according to the specified group
func (service *AdminDashboardService) getGroupWiseSessionScore(uow *repository.UnitOfWork, tenantID,
	batchID, talentID, sessionID, groupID uuid.UUID, feedbackScore *[]dashboard.FeedbackKeywords) error {

	// get feedback score for all the feedback given for sessions
	err := service.Repository.GetAll(uow, feedbackScore,
		repository.Table("faculty_talent_batch_session_feedback"),
		repository.Select([]string{"feedback_questions.`keyword`, (feedback_options.`key` * 10 / feedback_questions.`max_score`) AS keyword_score"}),
		repository.Join("INNER JOIN feedback_questions ON faculty_talent_batch_session_feedback.`question_id` = feedback_questions.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Join("INNER JOIN feedback_options ON faculty_talent_batch_session_feedback.`option_id` = feedback_options.`id`"+
			" AND feedback_questions.`tenant_id`=faculty_talent_batch_session_feedback.`tenant_id`"),
		repository.Filter("faculty_talent_batch_session_feedback.`talent_id`=? AND faculty_talent_batch_session_feedback.`batch_id`=?",
			talentID, batchID), repository.Filter("feedback_questions.`group_id` = ?", groupID),
		repository.Filter("feedback_questions.`has_options` = true AND"+
			" faculty_talent_batch_session_feedback.`batch_session_id`=?", sessionID),
		repository.Filter("faculty_talent_batch_session_feedback.`tenant_id`=?"+
			" AND faculty_talent_batch_session_feedback.`deleted_at` IS NULL", tenantID),
		// repository.Filter("feedback_questions.`is_active` = true"),
		repository.Filter("feedback_questions.`deleted_at` IS NULL AND feedback_options.`deleted_at` IS NULL"),
		repository.OrderBy("feedback_questions.`order`"))
	if err != nil {
		return err
	}
	return nil
}

// getBatchStatusDetails returns details of batch for specified batch status
func (service *AdminDashboardService) getBatchStatusDetails(batchDetails *[]dashboard.BatchDetails,
	tenantID uuid.UUID, requestForm url.Values) error {

	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batchDetails, "`batch_name`",
		repository.Table("batches"), repository.Select("`id` AS batch_id, `batch_name`, `is_b2b` AS batch_type, `total_students`"),
		repository.Filter("`deleted_at` IS NULL"), service.addBatchFeedbackSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// getTotalBatchCount returns total batch count for specified search queries
func (service *AdminDashboardService) getTotalBatchCount(uow *repository.UnitOfWork, totalCount *uint,
	tenantID uuid.UUID, requestForm url.Values) error {

	err := service.Repository.GetCountForTenant(uow, tenantID, batch.Batch{}, totalCount,
		service.addBatchFeedbackSearchQueries(requestForm))
	if err != nil {
		return err
	}
	return nil
}

// getInterviewRating will return interview rating of all the talents
func (service *AdminDashboardService) getInterviewRating(uow *repository.UnitOfWork, batchSessionFeedbackScore *[]dashboard.TalentFeedbackScore,
	tenantID uuid.UUID) error {

	for index := range *batchSessionFeedbackScore {

		err := service.Repository.Scan(uow, &(*batchSessionFeedbackScore)[index],
			repository.Model(&talent.Interview{}),
			repository.Select("AVG(`talent_interviews`.`rating`) AS interview_rating"),
			repository.Join("INNER JOIN `talent_interview_schedules` ON "+
				"`talent_interviews`.`schedule_id` = `talent_interview_schedules`.`id`"),
			repository.Filter("`talent_interviews`.`talent_id` = ?", (*batchSessionFeedbackScore)[index].TalentID),
			repository.Filter("talent_interviews.`tenant_id`=? AND talent_interviews.`deleted_at` IS NULL", tenantID),
			repository.Filter("talent_interview_schedules.`tenant_id`=? AND talent_interview_schedules.`deleted_at` IS NULL", tenantID))
		if err != nil {
			return err
		}
	}
	return nil
}

// getKeywords will return keyword names from feedback_questions table
func (service *AdminDashboardService) getKeywords(uow *repository.UnitOfWork,
	keywords *[]dashboard.KeywordName, tenantID uuid.UUID, groupID ...uuid.UUID) error {

	var queryProcessor repository.QueryProcessor

	if groupID != nil {
		queryProcessor = repository.Filter("`group_id` = ?", groupID)
	}

	err := service.Repository.GetAllForTenant(uow, tenantID, keywords,
		repository.Table("feedback_questions"), repository.Select("feedback_questions.`keyword` AS name"),
		repository.Filter("feedback_questions.`has_options` = true AND feedback_questions.`type`='Faculty_Session_Feedback'"),
		repository.Filter("feedback_questions.`is_active` = true"),
		queryProcessor, repository.OrderBy("feedback_questions.`order`"))
	if err != nil {
		return err
	}

	return nil
}

// getTalentFeelingDetails will return feeling details for specified talent and session
func (service *AdminDashboardService) getTalentFeelingDetails(uow *repository.UnitOfWork, tenantID,
	talentID, batchID, sessionID uuid.UUID, sessionFeedbackScore *dashboard.SessionKeywordFeedback) error {

	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.AhaMoment{},
		repository.Filter("`talent_id`=? AND `batch_id`=?",
			talentID, batchID))
	if err != nil {
		return err

	}
	if exist {
		err = service.Repository.GetRecord(uow, sessionFeedbackScore,
			repository.Table("aha_moments"),
			repository.Select("feelings.`feeling_name`, feeling_levels.`level_number`, feeling_levels.`description`"),
			repository.Join("INNER JOIN feelings ON feelings.`id` = aha_moments.`feeling_id` AND "+
				"feelings.`tenant_id` = aha_moments.`tenant_id`"),
			repository.Join("INNER JOIN feeling_levels ON feeling_levels.`id` = aha_moments.`feeling_level_id` AND "+
				"feeling_levels.`tenant_id` = aha_moments.`tenant_id`"),
			repository.Filter("feelings.`deleted_at` IS NULL AND feeling_levels.`deleted_at` IS NULL"),
			repository.Filter("aha_moments.`batch_id` = ?",
				batchID),
			repository.Filter("aha_moments.`talent_id` = ? AND aha_moments.`tenant_id` = ?",
				talentID, tenantID))
		if err != nil {
			return err
		}

	}

	return nil
}

// func (service *AdminDashboardService) getFeedbackQuestionGroupsByType(uow *repository.UnitOfWork,
// 	questionType string, tenantID uuid.UUID, groups *[]general.FeedbackQuestionGroup) error {

// 	// Get feedbackQuestionGroup list.
// 	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, groups, "`order`",
// 		repository.Filter("`type` = ?", questionType))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// addBatchFeedbackSearchQueries will add search queries from query params
func (service *AdminDashboardService) addBatchFeedbackSearchQueries(requestForm url.Values) repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if batchStatus, ok := requestForm["batchStatus"]; ok {
		util.AddToSlice("batches.`status`", "=?", "AND", batchStatus, &columnNames, &conditions, &operators, &values)
	}
	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("batches.`id`", "=?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}
	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("batches.`faculty_id`", "=?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)

}

// doesTenantExists validates tenantID
func (service *AdminDashboardService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}
