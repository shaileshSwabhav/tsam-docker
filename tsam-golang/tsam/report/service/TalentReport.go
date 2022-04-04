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
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentReportService provides method to get reports for talent.
type TalentReportService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTalentReportService returns a new instance of TalentReportService.
func NewTalentReportService(db *gorm.DB, repository repository.Repository) *TalentReportService {
	return &TalentReportService{
		DB:         db,
		Repository: repository,
	}
}

// GetTalentReport will return talent-report for the week.
func (service *TalentReportService) GetTalentReport(talentReport *report.TalentReport, talentID, tenantID uuid.UUID,
	requestForm url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(talentID, tenantID)
	if err != nil {
		return err
	}

	// Get this week's date on monday
	// weekStartDate := service.getFirstDateOfWeek()
	weekStartDate, err := service.getFirstDateOfWeek(time.Now())
	if err != nil {
		return err
	}

	// When date is passed.
	if _, ok := requestForm["date"]; ok {
		weekDate := requestForm.Get("date")
		weekStartDate, err = time.Parse(time.RFC3339, weekDate)
		if err != nil {
			return err
		}
	}

	// Create bucket for getting talent related information.
	tempTalent := list.Talent{}

	// Start enw transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent.
	err = service.Repository.GetForTenant(uow, tenantID, talentID, &tempTalent,
		repository.Filter("`is_active` = '1' AND `deleted_at` IS NULL"),
		service.addSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Give talent to talent report.
	talentReport.Talent = tempTalent
	talentReport.WorkingHours = make(map[uuid.UUID]report.WorkingHours)

	// Monday.
	err = service.getMondaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Tuesday.
	err = service.getTuesdaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Wednesday.
	err = service.getWednesdaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Thursday.
	err = service.getThursdaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Friday.
	err = service.getFridaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Saturday.
	err = service.getSaturdaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Sunday.
	err = service.getSundaySchedule(uow, tenantID, talentID, weekStartDate, talentReport)
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

// getMondaySchedule will load batch and timing for monday.
func (service *TalentReportService) getMondaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Monday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {

		batches := []report.Batch{}

		err = service.getDaySchedule(uow, tenantID, talentID, &batches, "Monday", weekStartDate)
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {
					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Monday = append(talentReport.Monday, batches[index])
				}
			}

		}

	}

	return nil
}

// getTuesdaySchedule will load batch and timing for tuesday.
func (service *TalentReportService) getTuesdaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Tuesday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 1).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Tuesday", weekStartDate.AddDate(0, 0, 1))
		if err != nil {
			return err
		}
		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Tuesday = append(talentReport.Tuesday, batches[index])
				}
			}
		}
	}

	return nil
}

// getWednesdaySchedule will load batch and timing for wednesday.
func (service *TalentReportService) getWednesdaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Wednesday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 2).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Wednesday", weekStartDate.AddDate(0, 0, 2))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Wednesday = append(talentReport.Wednesday, batches[index])
				}
			}
		}
	}
	return nil
}

// getThursdaySchedule will load batch and timing for thursday.
func (service *TalentReportService) getThursdaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Thursday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 3).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Thursday", weekStartDate.AddDate(0, 0, 3))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Thursday = append(talentReport.Thursday, batches[index])
				}
			}
		}
	}
	return nil
}

// getFridaySchedule will load batch and timing for friday.
func (service *TalentReportService) getFridaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Friday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 4).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Friday", weekStartDate.AddDate(0, 0, 4))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Friday = append(talentReport.Friday, batches[index])
				}
			}
		}
	}
	return nil
}

// getSaturdaySchedule will load batch and timing for saturday.
func (service *TalentReportService) getSaturdaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Saturday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 5).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Saturday", weekStartDate.AddDate(0, 0, 5))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {

					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Saturday = append(talentReport.Saturday, batches[index])
				}
			}
		}
	}
	return nil
}

// getSundaySchedule will load batch and timing for sunday.
func (service *TalentReportService) getSundaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	weekStartDate time.Time, talentReport *report.TalentReport) error {

	exist, err := repository.DoesRecordExist(service.DB, &report.Batch{},
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id` AND "+
			"batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Filter("days.`day` = 'Sunday' AND days.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batches.`tenant_id` = ?", tenantID),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.AddDate(0, 0, 6).Format("2006-01-02")))
	if err != nil {
		return err
	}
	if exist {
		batches := []report.Batch{}

		err := service.getDaySchedule(uow, tenantID, talentID, &batches, "Sunday", weekStartDate.AddDate(0, 0, 6))
		if err != nil {
			return err
		}

		if len(batches) != 0 {
			err = service.assignTotalWorkingHours(talentReport, &batches)
			if err != nil {
				return err
			}

			for index := range batches {
				if len(batches[index].BatchTimings) != 0 {
					err = service.assignTotalBatchWorkingHours(talentReport, batches[index])
					if err != nil {
						return err
					}
					talentReport.Sunday = append(talentReport.Sunday, batches[index])
				}
			}
		}
	}
	return nil
}

// getDaySchedule will return batch-timing for specified day.
func (service *TalentReportService) getDaySchedule(uow *repository.UnitOfWork, tenantID, talentID uuid.UUID,
	batches *[]report.Batch, day string, weekStartDate time.Time) error {

	err := service.Repository.GetAll(uow, batches,
		repository.Join("INNER JOIN batch_timing ON batches.`id` = batch_timing.`batch_id`"+
			" AND batches.`tenant_id` = batch_timing.`tenant_id`"),
		repository.Join("JOIN batch_talents ON batches.`id` = batch_talents.`batch_id`"),
		repository.Filter("`status` = 'Ongoing' AND batches.`is_active` = '1'"),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("? BETWEEN `start_date` AND `end_date`", weekStartDate.Format("2006-01-02")),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "BatchTimings",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
				repository.Filter("days.`day` = ? AND batch_timing.`tenant_id` = ? AND days.`deleted_at` IS NULL", day, tenantID),
			},
		}), repository.PreloadAssociations([]string{"BatchTimings.Day"}),
		repository.Filter("batches.`tenant_id` = ? AND batch_timing.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("batch_timing.`from_time`"), repository.GroupBy("batch_timing.`batch_id`"))
	if err != nil {
		return err
	}
	return nil
}

// assignTotalWorkingHours will assign totalhours and totalminutes to talentReport.
func (service *TalentReportService) assignTotalWorkingHours(talentReport *report.TalentReport,
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
	talentReport.TotalTrainingHours += totalTimeDiff
	return nil
}

// assignTotalBatchWorkingHours will assign totalhours and totalminutes for specified batch
func (service *TalentReportService) assignTotalBatchWorkingHours(talentReport *report.TalentReport,
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

	tempWorkingHours.TotalHours += talentReport.WorkingHours[batch.ID].TotalHours + timeDiff
	talentReport.WorkingHours[batch.ID] = tempWorkingHours

	return nil
}

// getFirstDateOfWeek will return the date of Monday this week
func (service *TalentReportService) getFirstDateOfWeek(now time.Time) (time.Time, error) {

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset), nil
}

// addSearchQueries will add search queries.
func (service *TalentReportService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	fmt.Println("================================================requestForm", requestForm)

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talents.`id`", "= ?", "AND", talentID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// returns error if there is no tenant record in table.
func (service *TalentReportService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *TalentReportService) doesTalentExist(talentID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Talent not found", exists, err); err != nil {
		return err
	}
	return nil
}
