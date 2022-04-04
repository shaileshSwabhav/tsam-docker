package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/report"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// FacultyReportService provides method to get login reports.
type FacultyReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewFacultyReportService returns a new instance of FacultyReportService.
func NewFacultyReportService(db *gorm.DB, repository repository.Repository) *FacultyReportService {
	return &FacultyReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetFacultyReport will return faculty-report for the week.
func (service *FacultyReportService) GetFacultyReport(tenantID uuid.UUID,
	facultyReport *[]report.FacultyReport, requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// get this weeks date on monday
	// weekStartDate := service.getFirstDateOfWeek()
	weekStartDate := util.GetBeginningOfWeek(time.Now())
	if err != nil {
		return err
	}

	// when date is passed.
	if _, ok := requestForm["date"]; ok {
		weekDate := requestForm.Get("date")

		// "2006-01-02"
		weekStartDate, err = time.Parse(time.RFC3339, weekDate)
		if err != nil {
			return err
		}
	}

	// for faculty_id
	// var queryProcessor repository.QueryProcessor
	// if facultyID, ok := requestForm["facultyID"]; ok {
	// 	queryProcessor = repository.Filter("`id` = ?", facultyID)
	// }

	// fmt.Println("week start date ->", weekStartDate)

	// temp variables
	faculty := []list.Faculty{}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get faculty.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, &faculty, "`first_name`",
		repository.Filter("`is_active` = '1' AND `deleted_at` IS NULL"), service.addSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range faculty {

		tempFacultyReport := report.FacultyReport{}
		tempFacultyReport.Faculty = faculty[index]
		tempFacultyReport.WorkingHours = make(map[uuid.UUID]report.WorkingHours)

		// monday
		err = service.getMondaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// tuesday
		err = service.getTuesdaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// wednesday
		err = service.getWednesdaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// thursday
		err = service.getThursdaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// friday
		err = service.getFridaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// saturday
		err = service.getSaturdaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		// sunday
		err = service.getSundaySchedule(uow, tenantID, faculty[index].ID, weekStartDate, &tempFacultyReport)
		if err != nil {
			uow.RollBack()
			return err
		}

		*facultyReport = append(*facultyReport, tempFacultyReport)
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// getMondaySchedule will load batch and timing for monday.
func (service *FacultyReportService) getMondaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Monday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ? AND batch_modules.`deleted_at` IS NULL", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`",
			weekStartDate.Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err = service.getDaySchedule(uow, tenantID, facultyID, &batches, "Monday", weekStartDate)
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {
					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Monday = append(tempFacultyReport.Monday, batches[index])
				}
			}

		}

	}

	return nil
}

// getTuesdaySchedule will load batch and timing for tuesday.
func (service *FacultyReportService) getTuesdaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Tuesday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`",
			weekStartDate.AddDate(0, 0, 1).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Tuesday", weekStartDate.AddDate(0, 0, 1))
		if err != nil {
			return err
		}
		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Tuesday = append(tempFacultyReport.Tuesday, batches[index])
				}
			}
		}
	}

	return nil
}

// getWednesdaySchedule will load batch and timing for wednesday.
func (service *FacultyReportService) getWednesdaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Wednesday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`", weekStartDate.AddDate(0, 0, 2).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Wednesday", weekStartDate.AddDate(0, 0, 2))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Wednesday = append(tempFacultyReport.Wednesday, batches[index])
				}
			}
		}
	}
	return nil
}

// getThursdaySchedule will load batch and timing for thursday.
func (service *FacultyReportService) getThursdaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Thursday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`", weekStartDate.AddDate(0, 0, 3).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Thursday", weekStartDate.AddDate(0, 0, 3))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Thursday = append(tempFacultyReport.Thursday, batches[index])
				}
			}
		}
	}
	return nil
}

// getFridaySchedule will load batch and timing for friday.
func (service *FacultyReportService) getFridaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Friday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`", weekStartDate.AddDate(0, 0, 4).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Friday", weekStartDate.AddDate(0, 0, 4))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Friday = append(tempFacultyReport.Friday, batches[index])
				}
			}
		}
	}
	return nil
}

// getSaturdaySchedule will load batch and timing for saturday.
func (service *FacultyReportService) getSaturdaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Saturday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`", weekStartDate.AddDate(0, 0, 5).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Saturday", weekStartDate.AddDate(0, 0, 5))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Saturday = append(tempFacultyReport.Saturday, batches[index])
				}
			}
		}
	}
	return nil
}

// getSundaySchedule will load batch and timing for sunday.
func (service *FacultyReportService) getSundaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	weekStartDate time.Time, tempFacultyReport *report.FacultyReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("days.`day` = 'Sunday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`",
			weekStartDate.AddDate(0, 0, 6).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, facultyID, &batches, "Sunday", weekStartDate.AddDate(0, 0, 6))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(tempFacultyReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {
					err = service.assignTotalBatchWorkingHours(tempFacultyReport, batches[index])
					if err != nil {
						return err
					}
					tempFacultyReport.Sunday = append(tempFacultyReport.Sunday, batches[index])
				}
			}
		}
	}
	return nil
}

// getDaySchedule will return batch-timing for specified day. #shailesh -> needs refactoring
func (service *FacultyReportService) getDaySchedule(uow *repository.UnitOfWork, tenantID, facultyID uuid.UUID,
	batches *[]report.Batch, day string, weekStartDate time.Time) error {

	err := service.Repository.GetAll(uow, batches,
		repository.Join("INNER JOIN batch_modules ON batch_modules.`batch_id` = batches.`id` AND"+
			" batch_modules.`tenant_id` = batches.`tenant_id`"),
		repository.Join("INNER JOIN batch_module_timings ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND "+
			"batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
		repository.Filter("batch_modules.`faculty_id` = ? AND batches.`status` = 'Ongoing' AND batches.`is_active` = '1'", facultyID),
		repository.Filter("? BETWEEN batch_modules.`start_date` AND batch_modules.`estimated_end_date`", weekStartDate.Format("2006-01-02")),
		repository.Filter("batches.`tenant_id` = ? AND batch_module_timings.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("batch_module_timings.`from_time`"), repository.GroupBy("batch_module_timings.`batch_id`"))
	if err != nil {
		return err
	}

	for index := range *batches {
		err = service.Repository.GetAll(uow, &(*batches)[index].BatchTimings,
			repository.Join("INNER JOIN batch_modules ON batch_modules.`id` = batch_module_timings.`batch_module_id` AND"+
				" batch_modules.`tenant_id` = batch_module_timings.`tenant_id`"),
			repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND"+
				" days.`tenant_id` = batch_module_timings.`tenant_id`"),
			repository.Filter("days.`day` = ? AND batch_module_timings.`tenant_id` = ? AND days.`deleted_at` IS NULL AND"+
				" batch_modules.`deleted_at` IS NULL", day, tenantID),
			repository.Filter("batch_modules.`batch_id` = ?", (*batches)[index].ID), repository.PreloadAssociations([]string{"Day"}),
			repository.GroupBy("batch_modules.`batch_id`"))
		if err != nil {
			return err
		}
	}

	// repository.PreloadWithCustomCondition(repository.Preload{
	// 	Schema: "BatchTimings",
	// 	Queryprocessors: []repository.QueryProcessor{
	// 		repository.Join("INNER JOIN days ON days.`id` = batch_module_timings.`day_id` AND days.`tenant_id` = batch_module_timings.`tenant_id`"),
	// 		repository.Filter("days.`day` = ? AND batch_module_timings.`tenant_id` = ? AND days.`deleted_at` IS NULL", day, tenantID),
	// 	},
	// }), repository.PreloadAssociations([]string{"BatchModuleTimings.Day"}),

	return nil
}

// assignTotalWorkingHours will assign totalhours and totalminutes to tempFacultyReport.
func (service *FacultyReportService) assignTotalWorkingHours(tempFacultyReport *report.FacultyReport,
	batches *[]report.Batch) error {

	var fromTime, toTime time.Time
	var totalTimeDiff float64
	var err error

	for index, batch := range *batches {
		var timeDiff float64
		for _, timing := range batch.BatchTimings {
			if timing.FromTime != nil {
				fromTime, err = time.Parse("15:04:05", *timing.FromTime)
				if err != nil {
					return err
				}
			}
			if timing.ToTime != nil {
				toTime, err = time.Parse("15:04:05", *timing.ToTime)
				if err != nil {
					return err
				}
			}
			totalTimeDiff += (toTime.Sub(fromTime).Hours())
			timeDiff += (toTime.Sub(fromTime).Hours())
		}
		(*batches)[index].TotalDailyHours = timeDiff
	}
	tempFacultyReport.TotalTrainingHours += totalTimeDiff
	return nil
}

// assignTotalBatchWorkingHours will assign totalhours and totalminutes for specified batch
func (service *FacultyReportService) assignTotalBatchWorkingHours(tempFacultyReport *report.FacultyReport,
	batch report.Batch) error {

	var fromTime, toTime time.Time
	var timeDiff float64
	var err error

	for _, timing := range batch.BatchTimings {
		if timing.FromTime != nil {
			fromTime, err = time.Parse("15:04:05", *timing.FromTime)
			if err != nil {
				return err
			}
		}
		if timing.ToTime != nil {
			toTime, err = time.Parse("15:04:05", *timing.ToTime)
			if err != nil {
				return err
			}
		}
		timeDiff += toTime.Sub(fromTime).Hours()
	}

	tempWorkingHours := report.WorkingHours{}
	tempWorkingHours.BatchName = batch.BatchName

	tempWorkingHours.TotalHours += tempFacultyReport.WorkingHours[batch.ID].TotalHours + timeDiff
	tempFacultyReport.WorkingHours[batch.ID] = tempWorkingHours

	return nil
}

// func (service *FacultyReportService) assignTotalDailyHours(tempFacultyReport *report.FacultyReport,
// 	batches *[]report.Batch) error {

// 	return nil
// }

// // calculateTotalTimeDifference will calculate difference between two times for all batches of specified day.
// func (service *FacultyReportService) calculateTotalTimeDifference(batches []report.Batch) (int, int, error) {

// 	var toTime, fromTime int = 0, 0

// 	for _, batch := range batches {
// 		for _, timing := range batch.BatchTimings {
// 			time1, err := service.removeColon(*timing.FromTime)
// 			if err != nil {
// 				return 0, 0, err
// 			}
// 			fromTime += time1

// 			time2, err := service.removeColon(*timing.ToTime)
// 			if err != nil {
// 				return 0, 0, err
// 			}
// 			toTime += time2
// 		}
// 	}

// 	hourDiff := toTime/100 - fromTime/100 - 1
// 	minDiff := toTime%100 + (60 - fromTime%100)

// 	for minDiff >= 60 {
// 		hourDiff++
// 		minDiff = minDiff - 60
// 	}

// 	return hourDiff, minDiff, nil
// }

// // calculateTimeDifference will calculate difference between two times for single batch.
// func (service *FacultyReportService) calculateTimeDifference(batch report.Batch) (int, int, error) {

// 	var toTime, fromTime int = 0, 0

// 	for _, timing := range batch.BatchTimings {
// 		time1, err := service.removeColon(*timing.FromTime)
// 		if err != nil {
// 			return 0, 0, err
// 		}
// 		fromTime += time1

// 		time2, err := service.removeColon(*timing.ToTime)
// 		if err != nil {
// 			return 0, 0, err
// 		}
// 		toTime += time2
// 	}

// 	hourDiff := toTime/100 - fromTime/100 - 1
// 	minDiff := toTime%100 + (60 - fromTime%100)

// 	for minDiff >= 60 {
// 		hourDiff++
// 		minDiff = minDiff - 60
// 	}

// 	return hourDiff, minDiff, nil
// }

// // remove ':' and convert it into an integer
// func (service *FacultyReportService) removeColon(time string) (int, error) {
// 	newTime := strings.Split(time, ":")
// 	time = newTime[0] + newTime[1]

// 	time = strings.Replace(time, ":", "", 1)
// 	return strconv.Atoi(time)

// }

// getFirstDateOfWeek will return the date of Monday this week
// func (service *FacultyReportService) getFirstDateOfWeek(now time.Time) time.Time {

// 	offset := int(time.Monday - now.Weekday())
// 	if offset > 0 {
// 		offset = -6
// 	}

// 	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
// }

// addSearchQueries will add search queries.
func (service *FacultyReportService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	fmt.Println("================================================requestForm", requestForm)

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if facultyID, ok := requestForm["facultyID"]; ok {
		util.AddToSlice("faculties.`id`", "= ?", "AND", facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *FacultyReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
