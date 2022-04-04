package service

import (
	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/models/dashboard"
	"github.com/techlabs/swabhav/tsam/repository"
)

// CourseBatchDashboardService provides all details to be shown on course dashboard.
type CourseBatchDashboardService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCourseBatchDashboardService returns new instance of CourseBatchDashboardService.
func NewCourseBatchDashboardService(db *gorm.DB, repository repository.Repository) *CourseBatchDashboardService {
	return &CourseBatchDashboardService{
		DB:         db,
		Repository: repository,
	}
}

// GetCourseBatchDashboardDetails gets all details required for CourseBatchDashboard.
func (ser *CourseBatchDashboardService) GetCourseBatchDashboardDetails(courseDashboard *dashboard.CourseDashboard,
	queryProcessors ...repository.QueryProcessor) error {
	// var totalCount int = 0
	// uow := repository.NewUnitOfWork(ser.DB, true)

	// // Get count of total courses.
	// if err := ser.Repository.GetCount(uow, common.Course{}, &totalCount); err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	// }
	// courseDashboard.TotalCourses = uint(totalCount)

	// // Get count of total batches.
	// if err := ser.Repository.GetCount(uow, batch.Batch{}, &totalCount); err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	// }
	// courseDashboard.TotalBatches = uint(totalCount)

	// // Get count of live batches.
	// if err := ser.Repository.GetCount(uow, batch.Batch{}, &totalCount,
	// 	repository.Filter("end_time IS NOT NULL")); err != nil {
	// 	return errors.NewHTTPError(err.Error(), http.StatusNotFound)
	// }
	// courseDashboard.AllLiveBatches = uint(totalCount)

	// if err := ser.getAllCourseDetails(courseDashboard); err != nil {
	// 	return err
	// }

	return nil

}

// Gets all fresher requirements.
// func (ser *CourseBatchDashboardService) getAllCourseDetails(courseDashboard *dashboard.CourseDashboard,
// 	queryProcessors ...repository.QueryProcessor) error {
// 	// var totalCount int = 0
// 	// uow := repository.NewUnitOfWork(ser.DB, true)
// 	course := &dashboard.CourseData{}
// 	course.CourseName = "786"
// 	// courseDashboard.Courses = append(courseDashboard.Courses, course)
// 	// for _, course := range courseDashboard.Courses {
// 	// 	ser.getAllBatchesForCourse(course)
// 	// 	ser.getAllSessionForCourse(course)
// 	// 	ser.getAllLiveBatchesForCourse(course)
// 	// }
// 	return nil
// }

// func (ser *CourseBatchDashboardService) getAllLiveBatchesForCourse(course *dashboard.CourseData,
// 	queryProcessors ...repository.QueryProcessor) error {
// 	// var totalCount int = 0
// 	// uow := repository.NewUnitOfWork(ser.DB, true)

// 	// course.LiveBatches.RequiredTalents = 786
// 	// course.LiveBatches.TotalBatches = 786
// 	// course.LiveBatches.TotalTalentsJoined = 786
// 	return nil
// }

// func (ser *CourseBatchDashboardService) getAllBatchesForCourse(course *dashboard.CourseData,
// 	queryProcessors ...repository.QueryProcessor) error {
// 	// var totalCount int = 0
// 	// uow := repository.NewUnitOfWork(ser.DB, true)

// 	course.AllBatches.TotalBatches = 786
// 	course.AllBatches.TotalTalentsJoined = 786
// 	return nil
// }

// func (ser *CourseBatchDashboardService) getAllSessionForCourse(courseDashboard *dashboard.CourseData,
// 	queryProcessors ...repository.QueryProcessor) error {
// 	// var totalCount int = 0
// 	// uow := repository.NewUnitOfWork(ser.DB, true)

// 	courseDashboard.Sessions.RemainingSessions = 786
// 	courseDashboard.Sessions.TotalSessions = 786
// 	return nil
// }
