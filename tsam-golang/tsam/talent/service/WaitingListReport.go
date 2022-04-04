package service

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/company"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// WaitingListReportService provides methods to get waiting list of talents and enquiries in a report format.
type WaitingListReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewWaitingListReportService returns new instance of WaitingListReportService.
func NewWaitingListReportService(db *gorm.DB, repository repository.Repository) *WaitingListReportService {
	return &WaitingListReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetCompanyBranchWaitingListReport gets waiting list report for company branch.
func (service *WaitingListReportService) GetCompanyBranchWaitingListReport(waitingListReport *[]tal.WaitingListCompanyBranchDTO, tenantID uuid.UUID,
	form url.Values, limit int, offset int, totalCount *tal.TotalCount) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	var queryProcessorsForSubQuery []repository.QueryProcessor
	queryProcessorsForSubQuery = append(queryProcessorsForSubQuery,
		repository.Select("IF(`talent_id` is not null,1,0) talents, IF(`talent_id` is null and `enquiry_id` is not null, 1, 0) enquiries,waiting_list.*"),
		repository.Filter("waiting_list.`deleted_at` IS NULL AND `is_active`=? AND `tenant_id`=? AND `company_branch_id` IS NOT NULL", 1, tenantID),
		repository.Table("waiting_list"),
		repository.GroupBy("`company_branch_id`, `email`"))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, tal.WaitingList{}, queryProcessorsForSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, waitingListReport, repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `company_branch_id` limit ? offset ?", subQuery, limit, limit*offset),
		repository.PreloadAssociations([]string{"CompanyBranch"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create subQuery for total count.
	subQueryForCount, err := service.Repository.SubQuery(uow, tal.WaitingList{},
		repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `company_branch_id`", subQuery),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total count of entries.
	if err := service.Repository.Scan(uow, totalCount,
		repository.RawQuery("select count(*) as total_count from ? as temptwo", subQueryForCount)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetCourseWaitingListReport gets waiting list report for course.
func (service *WaitingListReportService) GetCourseWaitingListReport(waitingListReport *[]tal.WaitingListCourseDTO, tenantID uuid.UUID,
	form url.Values, limit int, offset int, totalCount *tal.TotalCount) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	var queryProcessorsForSubQuery []repository.QueryProcessor
	queryProcessorsForSubQuery = append(queryProcessorsForSubQuery,
		repository.Select("IF(`talent_id` is not null,1,0) talents, IF(`talent_id` is null and `enquiry_id` is not null, 1, 0) enquiries,waiting_list.*"),
		repository.Filter("waiting_list.`deleted_at` IS NULL AND `is_active`=? AND `tenant_id`=? AND `course_id` IS NOT NULL", 1, tenantID),
		repository.Table("waiting_list"),
		repository.GroupBy("`course_id`, `email`"))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, tal.WaitingList{}, queryProcessorsForSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, waitingListReport, repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `course_id` limit ? offset ?", subQuery, limit, limit*offset),
		repository.PreloadAssociations([]string{"Course"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create subQuery for total count.
	subQueryForCount, err := service.Repository.SubQuery(uow, tal.WaitingList{},
		repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `course_id`", subQuery),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total count of entries.
	if err := service.Repository.Scan(uow, totalCount,
		repository.RawQuery("select count(*) as total_count from ? as temptwo", subQueryForCount)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetRequirementWaitingListReport gets waiting list report for requirement.
func (service *WaitingListReportService) GetRequirementWaitingListReport(waitingListReport *[]tal.WaitingListRequirementDTO,
	tenantID, companyBranchID uuid.UUID, form url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate company branch id.
	if err := service.doesCompanyBranchExist(companyBranchID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	var queryProcessorsForSubQuery []repository.QueryProcessor
	queryProcessorsForSubQuery = append(queryProcessorsForSubQuery,
		repository.Select("IF(`talent_id` is not null,1,0) talents, IF(`talent_id` is null and `enquiry_id` is not null, 1, 0) enquiries,waiting_list.*"),
		repository.Filter("waiting_list.`deleted_at` IS NULL AND `is_active`=? AND `tenant_id`=? AND `company_requirement_id` IS NOT NULL AND `company_branch_id`=?", 1, tenantID, companyBranchID),
		repository.Table("waiting_list"))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, tal.WaitingList{}, queryProcessorsForSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, waitingListReport, repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `company_requirement_id`", subQuery),
		repository.PreloadAssociations([]string{"CompanyRequirement"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetBatchWaitingListReport gets waiting list report for batch.
func (service *WaitingListReportService) GetBatchWaitingListReport(waitingListReport *[]tal.WaitingListBatchDTO,
	tenantID, courseID uuid.UUID, form url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate course id.
	if err := service.doesCourseExist(courseID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors.
	var queryProcessorsForSubQuery []repository.QueryProcessor
	queryProcessorsForSubQuery = append(queryProcessorsForSubQuery,
		repository.Select("IF(`talent_id` is not null,1,0) talents, IF(`talent_id` is null and `enquiry_id` is not null, 1, 0) enquiries,waiting_list.*"),
		repository.Filter("waiting_list.`deleted_at` IS NULL AND `is_active`=? AND `tenant_id`=? AND `batch_id` IS NOT NULL AND `course_id`=?", 1, tenantID, courseID),
		repository.Table("waiting_list"))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, tal.WaitingList{}, queryProcessorsForSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, waitingListReport, repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temp.* from ? as temp group by `batch_id`", subQuery),
		repository.PreloadAssociations([]string{"Batch"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetTechnologyWaitingListReport gets waiting list report by technology.
func (service *WaitingListReportService) GetTechnologyWaitingListReport(waitingListReport *[]tal.WaitingListTechnologyDTO, tenantID uuid.UUID,
	form url.Values, limit int, offset int, totalCount *tal.TotalCount) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors for query number one.
	var queryProcessorsForSubQueryOne []repository.QueryProcessor
	queryProcessorsForSubQueryOne = append(queryProcessorsForSubQueryOne,
		repository.Select("IF(company_requirements_technologies.`technology_id` IS NULL,courses_technologies.`technology_id`,company_requirements_technologies.`technology_id`) as technology_id , waiting_list.*"),
		repository.Filter("waiting_list.`deleted_at` IS NULL AND waiting_list.`is_active`=? AND waiting_list.`tenant_id`=?", 1, tenantID),
		repository.Table("waiting_list"),
		repository.Join("LEFT JOIN company_requirements on waiting_list.`company_requirement_id` = company_requirements.`id`"),
		repository.Join("LEFT JOIN company_requirements_technologies on company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
		repository.Join("LEFT JOIN courses on waiting_list.`course_id` = courses.`id`"),
		repository.Join("LEFT JOIN courses_technologies on courses.`id` = courses_technologies.`course_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, tal.WaitingList{}, queryProcessorsForSubQueryOne...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create subQuery for query number two.
	subQueryTwo, err := service.Repository.SubQuery(uow, tal.WaitingList{},
		repository.RawQuery("SELECT IF(`talent_id` is not null,1,0) talents, IF(`talent_id` is null and `enquiry_id` is not null, 1, 0) enquiries,temp.* FROM ? as temp GROUP BY `technology_id`, `email`", subQueryOne),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total entries.
	if err := service.Repository.GetAll(uow, waitingListReport, repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temptwo.* FROM ? as temptwo GROUP BY `technology_id` limit ? offset ?", subQueryTwo, limit, limit*offset),
		repository.PreloadAssociations([]string{"Technology"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create subQuery for total count.
	subQueryForCount, err := service.Repository.SubQuery(uow, tal.WaitingList{},
		repository.RawQuery("select sum(talents) as talent_count, sum(enquiries) as enquiry_count, temptwo.* from ? as temptwo GROUP BY `technology_id`", subQueryTwo),
	)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Get total count of entries.
	if err := service.Repository.Scan(uow, totalCount,
		repository.RawQuery("select count(*) as total_count from ? as tempthree", subQueryForCount)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *WaitingListReportService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID, orderBy *string) []repository.QueryProcessor {
	fmt.Println("=========================In waiting list report search============================", searchForm)

	// Check if there is search criteria given.
	if len(searchForm) == 0 {
		return nil
	}

	// Create query precessors.
	var queryProcessors []repository.QueryProcessor

	// Company branch id filter.
	if _, ok := searchForm["companyBranch"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Filter("`company_branch_id` IS NOT NULL"),
			repository.GroupBy("`company_branch_id`, `email`"),
			repository.PreloadAssociations([]string{"CompanyBranch"}))

		orderByInString := "company_branches.`company_name`"
		orderBy = &orderByInString
	}

	// Course id filter.
	if courseID, ok := searchForm["courseID"]; ok {
		queryProcessors = append(queryProcessors,
			repository.Filter("`course_id`=?", courseID),
			repository.GroupBy("`course_id`, `email`"),
			repository.PreloadAssociations([]string{"Course"}))

		orderByInString := "courses.`name`"
		orderBy = &orderByInString
	}

	return queryProcessors
}

// doesTenantExist validates if tenant exists or not in database.
func (service *WaitingListReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCompanyBranchExist validates if company branch exists or not in database.
func (service *WaitingListReportService) doesCompanyBranchExist(companyBranchID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Branch{},
		repository.Filter("`id` = ?", companyBranchID))
	if err := util.HandleError("Invalid company branch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCourseExist validates if batch exists or not in database.
func (service *WaitingListReportService) doesCourseExist(courseID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
		repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// // sortTalentAcademicsByOrder sorts all the talent academics by order of Passout field.
// func (service *WaitingListReportService) sortTalentAcademicsByOrder(talents *[]tal.DTO) {
// 	if talents != nil && len(*talents) != 0 {
// 		for i := 0; i < len(*talents); i++ {
// 			academics := &(*talents)[i].Academics
// 			for j := 0; j < len(*academics); j++ {
// 				if (*academics)[j].Passout == 0 {
// 					return
// 				}
// 			}
// 			for j := 0; j < len(*academics); j++ {
// 				sort.Slice(*academics, func(p, q int) bool {
// 					return (*academics)[p].Passout < (*academics)[q].Passout
// 				})
// 			}
// 		}
// 	}
// }

// // sortEnquiryAcademicsByOrder sorts all the enquiry academics by order of Passout field.
// func (service *WaitingListReportService) sortEnquiryAcademicsByOrder(enquiries *[]talenq.Enquiry) {
// 	if enquiries != nil && len(*enquiries) != 0 {
// 		for i := 0; i < len(*enquiries); i++ {
// 			academics := &(*enquiries)[i].Academics
// 			for j := 0; j < len(*academics); j++ {
// 				if (*academics)[j].Passout == 0 {
// 					return
// 				}
// 			}
// 			for j := 0; j < len(*academics); j++ {
// 				sort.Slice(*academics, func(p, q int) bool {
// 					return (*academics)[p].Passout < (*academics)[q].Passout
// 				})
// 			}
// 		}
// 	}
// }

// // createWaitingListPartOfTalentWaitingList creates waiting list part of talent waiting list.
// func (service *WaitingListReportService) createWaitingListPartOfTalentWaitingList(waitingList *tal.WaitingListDTO,
// 	talentWaitingList *tal.TalentWaitingList) {
// 	talentWaitingList.TalentID = *waitingList.TalentID
// 	if waitingList.CompanyBranch != nil {
// 		talentWaitingList.CompanyBranch = &waitingList.CompanyBranch.BranchName
// 	}
// 	if waitingList.CompanyRequirement != nil {
// 		talentWaitingList.CompanyRequirement = &waitingList.CompanyRequirement.Branch.BranchName
// 		talentWaitingList.CompanyRequirementCode = &waitingList.CompanyRequirement.Code
// 	}
// 	if waitingList.Course != nil {
// 		talentWaitingList.Course = &waitingList.Course.Name
// 	}
// 	if waitingList.Course != nil {
// 		talentWaitingList.Batch = &waitingList.Batch.BatchName
// 	}
// }

// // createTalentPartOfTalentWaitingList creates talent part of talent waiting list.
// func (service *WaitingListReportService) createTalentPartOfTalentWaitingList(talent *tal.DTO, talentWaitingList *tal.TalentWaitingList) {
// 	talentWaitingList.FirstName = talent.FirstName
// 	talentWaitingList.LastName = talent.LastName
// 	talentWaitingList.Contact = talent.Contact
// 	talentWaitingList.AcademicYear = talent.AcademicYear
// 	if talent.SalesPerson != nil {
// 		talentWaitingList.SalesPersonFirstName = &talent.SalesPerson.FirstName
// 		talentWaitingList.SalesPersonLastName = &talent.SalesPerson.LastName
// 	}

// 	if len(talent.Academics) != 0 {
// 		talentAcademic := talent.Academics[len(talent.Academics)-1]
// 		talentWaitingList.College = &talentAcademic.College
// 	}
// }

// // createWaitingListPartOfEnquiryWaitingList creates waiting list part of enquiry waiting list.
// func (service *WaitingListReportService) createWaitingListPartOfEnquiryWaitingList(waitingList *tal.WaitingListDTO,
// 	enquiryWaitingList *talenq.EnquiryWaitingList) {
// 	enquiryWaitingList.EnquiryID = *waitingList.EnquiryID
// 	if waitingList.CompanyBranch != nil {
// 		enquiryWaitingList.CompanyBranch = &waitingList.CompanyBranch.BranchName
// 	}
// 	if waitingList.CompanyRequirement != nil {
// 		enquiryWaitingList.CompanyRequirement = &waitingList.CompanyRequirement.Branch.BranchName
// 		enquiryWaitingList.CompanyRequirementCode = &waitingList.CompanyRequirement.Code
// 	}
// 	if waitingList.Course != nil {
// 		enquiryWaitingList.Course = &waitingList.Course.Name
// 	}
// 	if waitingList.Course != nil {
// 		enquiryWaitingList.Batch = &waitingList.Batch.BatchName
// 	}
// }

// // createEnquiryPartOfEnquiryWaitingList creates enquiry part of enquiry waiting list.
// func (service *WaitingListReportService) createEnquiryPartOfEnquiryWaitingList(enquiry *talenq.DTO, enquiryWaitingList *talenq.EnquiryWaitingList) {
// 	enquiryWaitingList.FirstName = enquiry.FirstName
// 	enquiryWaitingList.LastName = enquiry.LastName
// 	enquiryWaitingList.Contact = enquiry.Contact
// 	enquiryWaitingList.AcademicYear = enquiry.AcademicYear
// 	if enquiry.SalesPerson != nil {
// 		enquiryWaitingList.SalesPersonFirstName = &enquiry.SalesPerson.FirstName
// 		enquiryWaitingList.SalesPersonLastName = &enquiry.SalesPerson.LastName
// 	}

// 	if len(enquiry.Academics) != 0 {
// 		enquiryAcademic := enquiry.Academics[len(enquiry.Academics)-1]
// 		enquiryWaitingList.College = enquiryAcademic.College
// 	}
// }

// // GetTalentWaitingListByCompanyRequirement returns all talents and its corresponding waiting list entries by
// // company requirement id.
// func (service *WaitingListReportService) GetTalentWaitingListByCompanyRequirement(talentWaitingLists *[]tal.TalentWaitingList,
// 	tenantID uuid.UUID, limit, offset int, totalCount *int) error {

// 	// Validate tenant id.
// 	if err := service.doesTenantExist(tenantID); err != nil {
// 		return err
// 	}

// 	// Start new transaction.
// 	uow := repository.NewUnitOfWork(service.DB, true)

// 	// Create bucket for talents and waiting list entries.
// 	talents := []tal.DTO{}
// 	waitingLists := []tal.WaitingListDTO{}

// 	// Get all waiting list entries from database.
// 	if err := service.Repository.GetAllInOrder(uow, &waitingLists, "first_name, last_name",
// 		repository.Join("JOIN talents ON talents.`id` = waiting_list.`talent_id`"),
// 		repository.Filter("talents.`deleted_at` IS NULL"),
// 		repository.Filter("waiting_list.`tenant_id`=?", tenantID),
// 		repository.Filter("talents.`tenant_id`=?", tenantID),
// 		repository.Filter("waiting_list.`company_branch_id` IS NOT NULL"),
// 		repository.PreloadAssociations([]string{"CompanyBranch", "CompanyRequirement", "CompanyRequirement.Branch", "Course", "Batch"}),
// 		repository.Paginate(limit, offset, totalCount)); err != nil {
// 		uow.RollBack()
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError("Record not found")
// 	}

// 	// Get talents from database.
// 	if err := service.Repository.GetAllInOrder(uow, &talents, "first_name, last_name",
// 		repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
// 		repository.Filter("waiting_list.`deleted_at` IS NULL"),
// 		repository.Filter("waiting_list.`tenant_id`=?", tenantID),
// 		repository.Filter("talents.`tenant_id`=?", tenantID),
// 		repository.Filter("waiting_list.`company_branch_id` IS NOT NULL"),
// 		repository.PreloadAssociations([]string{"Academics", "SalesPerson"}),
// 		repository.Paginate(limit, offset, totalCount)); err != nil {
// 		uow.RollBack()
// 		log.NewLogger().Error(err.Error())
// 		return errors.NewValidationError("Record not found")
// 	}

// 	// Give waiting list entries to talent waiting lists.
// 	for _, waitingList := range waitingLists {
// 		talentWaitingList := tal.TalentWaitingList{}
// 		service.createWaitingListPartOfTalentWaitingList(&waitingList, &talentWaitingList)
// 		*talentWaitingLists = append(*talentWaitingLists, talentWaitingList)
// 	}

// 	// Sort the child tables.
// 	service.sortTalentAcademicsByOrder(&talents)

// 	// Give talents to talent waiting lists.
// 	for _, talent := range talents {
// 		for i := 0; i < len(*talentWaitingLists); i++ {
// 			if talent.ID == (*talentWaitingLists)[i].TalentID {
// 				service.createTalentPartOfTalentWaitingList(&talent, &(*talentWaitingLists)[i])
// 			}
// 		}
// 	}

// 	// fmt.Println("talents waiting list*************************************************************")
// 	// fmt.Println("talent id :", (*talentWaitingLists)[0].TalentID)
// 	// fmt.Println("first name :", (*talentWaitingLists)[0].FirstName)
// 	// fmt.Println("last name :", (*talentWaitingLists)[0].LastName)
// 	// fmt.Println("CONTACT :", (*talentWaitingLists)[0].Contact)
// 	// fmt.Println("academic year :", (*talentWaitingLists)[0].AcademicYear)
// 	// fmt.Println("college :", (*talentWaitingLists)[0].College)
// 	// fmt.Println("salesperson :", (*talentWaitingLists)[0].SalesPersonFirstName)
// 	// fmt.Println("branch :", (*talentWaitingLists)[0].SalesPersonLastName)
// 	// fmt.Println("req :", (*talentWaitingLists)[0].CompanyBranch)
// 	// fmt.Println("course :", (*talentWaitingLists)[0].CompanyRequirement)
// 	// fmt.Println("batch :", (*talentWaitingLists)[0].Course)
// 	// fmt.Println("batch :", (*talentWaitingLists)[0].Batch)

// 	// fmt.Println("talent id :", (*talentWaitingLists)[1].TalentID)
// 	// fmt.Println("first name :", (*talentWaitingLists)[1].FirstName)
// 	// fmt.Println("last name :", (*talentWaitingLists)[1].LastName)
// 	// fmt.Println("CONTACT :", (*talentWaitingLists)[1].Contact)
// 	// fmt.Println("academic year :", (*talentWaitingLists)[1].AcademicYear)
// 	// fmt.Println("college :", (*talentWaitingLists)[1].College)
// 	// fmt.Println("salesperson :", (*talentWaitingLists)[1].SalesPersonFirstName)
// 	// fmt.Println("branch :", (*talentWaitingLists)[1].SalesPersonLastName)
// 	// fmt.Println("req :", (*talentWaitingLists)[1].CompanyBranch)
// 	// fmt.Println("course :", (*talentWaitingLists)[1].CompanyRequirement)
// 	// fmt.Println("batch :", (*talentWaitingLists)[1].Course)
// 	// fmt.Println("batch :", (*talentWaitingLists)[1].Batch)

// 	uow.Commit()
// 	return nil
// }
