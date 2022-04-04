package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/admin"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TimesheetService provides methods to do different CRUD operations on time_sheet table.
type TimesheetService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewTimesheetService returns a new instance Of TimesheetService.
func NewTimesheetService(db *gorm.DB, repository repository.Repository) *TimesheetService {
	return &TimesheetService{
		DB:           db,
		Repository:   repository,
		associations: []string{"Activities.Project", "Activities.SubProject", "Activities.Batch"},
	}
}

// AddTimesheet will add new timesheet to the table
func (service *TimesheetService) AddTimesheet(timesheet *admin.Timesheet, uows ...*repository.UnitOfWork) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(timesheet, timesheet.CreatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(timesheet)
	if err != nil {
		return err
	}

	service.assignCreatedBy(timesheet)

	// Starting transaction.
	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// No activity record if on leave is set
	if timesheet.IsOnLeave {
		timesheet.Activities = nil
	}

	// check if record exist for date
	exist, err := repository.DoesRecordExistForTenant(uow.DB, timesheet.TenantID, admin.Timesheet{},
		repository.Filter("`date` = ? AND `department_id` = ? AND `credential_id` = ? ",
			timesheet.Date, timesheet.DepartmentID, timesheet.CredentialID))
	if err != nil {
		return err
	}
	// If record with same date exist, then add activities for that date.
	if exist {
		tempTimesheet := admin.Timesheet{}
		err = service.Repository.GetRecordForTenant(uow, timesheet.TenantID, &tempTimesheet,
			repository.Filter("`date` = ? AND `department_id` = ? AND `credential_id` = ? ",
				timesheet.Date, timesheet.DepartmentID, timesheet.CredentialID), repository.Select("`id`"))
		if err != nil {
			return err
		}

		// activities :=
		for index := range timesheet.Activities {
			timesheet.Activities[index].TimesheetID = tempTimesheet.ID
			err = service.Repository.Add(uow, &timesheet.Activities[index])
			if err != nil {
				// Rollback only if no transaction is passed.
				if length == 0 {
					uow.RollBack()
				}
				return err
			}
		}

		// Commit only if no transaction is passed.
		if length == 0 {
			uow.Commit()
		}
		return nil
	}

	//  else, create a new timesheet entry
	err = service.Repository.Add(uow, timesheet)
	if err != nil {
		// Rollback only if no transaction is passed.
		if length == 0 {
			uow.RollBack()
		}
		return err
	}

	// Commit only if no transaction is passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddTimesheets adds multiple timesheets to Database.
func (service *TimesheetService) AddTimesheets(timesheets []*admin.Timesheet) error {
	// If a date has a non-leave or a leave record in timesheet, a leave cannot be added to that date!
	timesheetMap := make(map[string]int8, len(timesheets))
	for index, timesheet := range timesheets {
		if timesheetMap[timesheet.Date] == -1 {
			log.NewLogger().Error("Leave is mentioned for the following date:" + timesheets[index].Date)
			return errors.NewValidationError("Leave is already mentioned for the following date:" +
				timesheets[index].Date)
		}
		timesheetMap[timesheet.Date] = timesheetMap[timesheet.Date] + 1
		if timesheet.IsOnLeave {
			if timesheetMap[timesheet.Date] != 1 {
				log.NewLogger().Error("Data exists for the following date:" + timesheets[index].Date)
				return errors.NewValidationError("Data exists for the following date:" +
					timesheets[index].Date)
			}
			timesheetMap[timesheet.Date] = -1
		}
	}

	// Add one timesheet record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, timesheet := range timesheets {
		err := service.AddTimesheet(timesheet, uow)
		if err != nil {
			uow.RollBack()
			return errors.NewValidationError("date-" + timesheet.Date + ": " + err.Error())
		}
	}

	// Commit only if all timesheets have been added.
	uow.Commit()
	return nil

}

// UpdateTimesheet will update the specified timesheet
func (service *TimesheetService) UpdateTimesheet(timesheet *admin.Timesheet) error {

	// check if all foreign keys exist.
	err := service.doForeignKeysExist(timesheet, timesheet.UpdatedBy)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(timesheet)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// if using save then created_by field is required to be fetched from DB else it will be marked as null
	tempTimesheet := admin.Timesheet{}
	err = service.Repository.GetRecordForTenant(uow, timesheet.TenantID, &tempTimesheet,
		repository.Filter("`id` = ?", timesheet.ID), repository.Select([]string{"`created_by`,`credential_id`"}))
	if err != nil {
		uow.RollBack()
		return err
	}
	timesheet.CreatedBy = tempTimesheet.CreatedBy
	if timesheet.UpdatedBy != tempTimesheet.CredentialID {
		log.NewLogger().Error("Not authorized to update timesheet.")
		return errors.NewValidationError("Not authorized to update timesheet.")
	}

	exists, err := repository.DoesRecordExistForTenant(uow.DB, timesheet.TenantID, admin.Timesheet{},
		repository.Filter("`date` = ? AND `department_id` = ? AND `credential_id` = ? AND `id` != ?",
			timesheet.Date, timesheet.DepartmentID, timesheet.CredentialID, timesheet.ID))
	err = util.HandleIfExistsError("Record already exists for the date: "+timesheet.Date, exists, err)
	if err != nil {
		return err
	}

	err = service.updateTimesheetActivities(uow, timesheet)
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.Save(uow, timesheet)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// DeleteTimesheet will delete the specified timesheet from the table
func (service *TimesheetService) DeleteTimesheet(timesheet *admin.Timesheet) error {

	// check if tenant exist
	err := service.doesTenantExist(timesheet.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist
	err = service.doesCredentialExist(timesheet.TenantID, timesheet.DeletedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempTimesheet := admin.Timesheet{}
	err = service.Repository.GetRecordForTenant(uow, timesheet.TenantID, &tempTimesheet,
		repository.Filter("`id` = ?", timesheet.ID), repository.Select([]string{"`credential_id`"}))
	if err != nil {
		uow.RollBack()
		return err
	}
	if timesheet.DeletedBy != tempTimesheet.CredentialID {
		log.NewLogger().Error("Not authorized to delete timesheet.")
		return errors.NewValidationError("Not authorized to delete timesheet.")
	}

	// soft-delete timesheet activities
	err = service.Repository.UpdateWithMap(uow, &admin.TimesheetActivity{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": timesheet.DeletedBy,
	}, repository.Filter("`timesheet_id` = ?", timesheet.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	// soft-delete timesheet
	err = service.Repository.UpdateWithMap(uow, &admin.Timesheet{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": timesheet.DeletedBy,
	}, repository.Filter("`id` = ?", timesheet.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateTimesheetActivity will update specific timesheet activity
func (service *TimesheetService) UpdateTimesheetActivity(timesheetActivity *admin.TimesheetActivity) error {

	tenantID := timesheetActivity.TenantID

	// check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exists for master fields.
	err = service.doesCredentialExist(tenantID, timesheetActivity.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if timesheet foreign key exists.
	err = service.doesActivitiesForeignKeyExist(tenantID, timesheetActivity)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// while using save, created_by field is required to be fetched from DB else it will be marked as null
	timesheetInDB := &admin.Timesheet{}
	err = service.Repository.GetRecordForTenant(uow, tenantID, &timesheetInDB,
		repository.Filter("`id` = ?", timesheetActivity.TimesheetID),
		repository.Select([]string{"`credential_id`", "`department_id`", "`date`", "`tenant_id`"}))
	if err != nil {
		uow.RollBack()
		return err
	}
	if timesheetActivity.UpdatedBy != timesheetInDB.CredentialID {
		log.NewLogger().Error("Not authorized to update timesheet.")
		return errors.NewValidationError("Not authorized to update timesheet.")
	}

	timesheetActivityInDB := admin.TimesheetActivity{}
	err = service.Repository.GetRecordForTenant(uow, tenantID, &timesheetActivityInDB,
		repository.Filter("`id` = ?", timesheetActivity.ID),
		repository.Select([]string{"`created_by`", "`next_estimated_date`"}))
	if err != nil {
		uow.RollBack()
		return err
	}
	if timesheetActivityInDB.NextEstimatedDate != nil {
		log.NewLogger().Error("Already has a record on date: " + *timesheetActivityInDB.NextEstimatedDate)
		return errors.NewValidationError("Can't update activity as this activity has been transferred to: " +
			*timesheetActivityInDB.NextEstimatedDate)
	}

	timesheetActivity.CreatedBy = timesheetActivityInDB.CreatedBy
	if timesheetActivity.NextEstimatedDate != nil {
		date, err := time.Parse(time.RFC3339, timesheetInDB.Date)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		timesheetInDB.Date = date.Format("2006-01-02")
		err = timesheetInDB.ValidateNextEstimateDate(*timesheetActivity.NextEstimatedDate)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		err = service.addNextEstimateTimesheet(uow, timesheetInDB, *timesheetActivity)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	err = service.Repository.Save(uow, timesheetActivity)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetTimesheets returns timesheets based on limit and offset
func (service *TimesheetService) GetTimesheets(tenantID uuid.UUID,
	timesheets *[]*admin.TimesheetDTO, form url.Values, header *admin.TimesheetHeader) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Orderby query should start with ASC if from date is searched.
	orderBy := "`timesheets`.`date` DESC"
	if _, ok := form["fromDate"]; ok {
		orderBy = "`timesheets`.`date`"
	}
	// Get all timesheets and orderby timesheet_name.
	err = service.Repository.GetAllInOrder(uow, timesheets, orderBy,
		// repository.Filter("`timesheets`.`deleted_at` IS NULL"),
		service.addSearchQueries(form),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Activities",
			Queryprocessors: []repository.QueryProcessor{
				service.addIsCompletedSearchQueries(form),
			},
		}), repository.PreloadAssociations(service.associations),
		repository.Join("LEFT JOIN timesheet_activities ON timesheets.`id` = timesheet_activities.`timesheet_id` AND "+
			"timesheets.`tenant_id` = timesheet_activities.`tenant_id`"),
		repository.Filter("`timesheets`.`tenant_id` = ?", tenantID), repository.GroupBy("timesheets.`id`"),
		repository.Paginate(header.Limit, header.Offset, &header.TotalCount))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Sum all hours_needed.
	err = service.Repository.GetRecord(uow, header,
		repository.Table("timesheets"),
		repository.Filter("`timesheets`.`deleted_at` IS NULL AND `timesheet_activities`.`deleted_at` IS NULL"),
		repository.Select("SUM(`timesheet_activities`.`hours_needed`) as total_hours"),
		repository.Join("LEFT JOIN timesheet_activities ON timesheets.`id` = timesheet_activities.`timesheet_id` AND "+
			"timesheets.`tenant_id` = timesheet_activities.`tenant_id`"),
		repository.Filter("`timesheets`.`tenant_id` = ?", tenantID),
		service.addSearchQueries(form))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}

	// =====================================Below code just for faculty=====================================
	var credential list.Credential
	credentialID := form.Get("credentialID")
	if util.IsEmpty(credentialID) {
		return nil
	}
	credential.ID, err = util.ParseUUID(credentialID)
	if err != nil {
		return nil
	}

	err = service.Repository.GetRecordForTenant(uow, tenantID, &credential,
		repository.Filter("`credentials`.`deleted_at` IS NULL AND `id` = ?", credential.ID),
		repository.PreloadAssociations([]string{"Role"}))
	if err != nil {
		return nil
	}

	if credential.Role.RoleName == "Faculty" {
		// Get free hours.
		err = service.Repository.GetRecord(uow, header,
			repository.Table("timesheets"),
			repository.Join("LEFT JOIN timesheet_activities ON timesheets.`id` = timesheet_activities.`timesheet_id` AND "+
				"timesheets.`tenant_id` = timesheet_activities.`tenant_id`"),
			repository.Filter("`timesheet_activities`.`deleted_at` IS NULL AND `timesheet_activities`.`batch_id` is NULL AND "+
				"`timesheets`.`deleted_at` IS NULL AND `timesheets`.`is_on_leave` = ?", false),
			repository.Filter("`timesheets`.`tenant_id` = ?", tenantID),
			repository.Select("SUM(`timesheet_activities`.`hours_needed`) as free_hours"),
			service.addSearchQueries(form))
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
	}

	// =====================================Above code just for faculty=====================================
	// jump here to skip the free hours calculation which is just meant for faculty.
	// Old batch sessions
	// for _, timesheet := range *timesheets {
	// 	if timesheet.Activities != nil {
	// 		for _, activity := range timesheet.Activities {
	// 			// if activity.BatchSessionID != nil {
	// 			if activity.BatchSessionID != nil {

	// 				err = service.Repository.GetRecordForTenant(uow, tenantID, activity.BatchSession,
	// 					repository.Filter("`batch_session_id` = ?", activity.BatchSessionID),
	// 					repository.PreloadAssociations([]string{"BatchSessionTopic"}))
	// 				if err != nil {
	// 					uow.RollBack()
	// 					return err
	// 				}

	// 				// batchSession := list.BatchSession{}
	// 				// // Get session details for each timesheet.
	// 				// err = service.Repository.GetRecord(uow, &batchSession,
	// 				// 	repository.Select([]string{"`batch_sessions`.`id`,`course_sessions`.`name`"}),
	// 				// 	repository.Join("INNER JOIN `course_sessions` ON "+
	// 				// 		"`batch_sessions`.`course_session_id` = `course_sessions`.`id` AND "+
	// 				// 		"`batch_sessions`.`tenant_id` = `course_sessions`.`tenant_id`"),
	// 				// 	repository.Filter("`batch_sessions`.`id` = ?", activity.BatchSessionID))
	// 				// if err != nil {
	// 				// 	uow.RollBack()
	// 				// 	log.NewLogger().Error(err.Error())
	// 				// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// 				// }
	// 				// activity.BatchSession = &batchSession
	// 			}
	// 		}
	// 	}
	// }

	// for _, timesheet := range *timesheets {
	// 	for _, activity := range timesheet.Activities {
	// 		// if activity.BatchSessionID != nil {
	// 		if activity.BatchSessionID != nil {

	// 			err = service.Repository.GetRecordForTenant(uow, tenantID, &activity.BatchSession,
	// 				repository.Filter("`batch_session_id` = ?", activity.BatchSessionID),
	// 				repository.PreloadAssociations([]string{"BatchSession", "BatchSession.BatchSessionTopic"}))
	// 			if err != nil {
	// 				uow.RollBack()
	// 				return err
	// 			}
	// 			// batchSession := list.BatchSession{}
	// 			// // Get session details for each activity.
	// 			// err = service.Repository.GetRecord(uow, &batchSession,
	// 			// 	repository.Select([]string{"`batch_sessions`.`id`,`course_sessions`.`name`"}),
	// 			// 	repository.Join("INNER JOIN `course_sessions` ON "+
	// 			// 		"`batch_sessions`.`course_session_id` = `course_sessions`.`id` AND "+
	// 			// 		"`batch_sessions`.`tenant_id` = `course_sessions`.`tenant_id`"),
	// 			// 	repository.Filter("`batch_sessions`.`id` = ?", activity.BatchSessionID))
	// 			// if err != nil {
	// 			// 	uow.RollBack()
	// 			// 	log.NewLogger().Error(err.Error())
	// 			// 	return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	// 			// }
	// 			// activity.BatchSession = &batchSession
	// 		}
	// 	}
	// }

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// validateFieldUniqueness will check the table for any repitition for unique fields
func (service *TimesheetService) validateFieldUniqueness(timesheet *admin.Timesheet) error {
	// If on leave flag is set, we need to check if there are previous records or data on that date.
	queryProcessors := []repository.QueryProcessor{
		repository.Filter("`date` = ? AND `credential_id` = ?  AND `id`!= ?", timesheet.Date,
			timesheet.CredentialID, timesheet.ID)}
	if !timesheet.IsOnLeave {
		queryProcessors = append(queryProcessors, repository.Filter("`is_on_leave` = ?", true))
	}
	errMsg := "Can't add data as there is a leave mentioned on "
	if timesheet.IsOnLeave {
		errMsg = "Can't add leave as there is data or a leave mentioned on "
	}
	errMsg += timesheet.Date
	exists, err := repository.DoesRecordExistForTenant(service.DB, timesheet.TenantID, admin.Timesheet{},
		queryProcessors...)
	err = util.HandleIfExistsError(errMsg, exists, err)
	if err != nil {
		return err
	}
	return nil
}

// doForeignKeysExist checks if all foregin key exist and if not returns error
func (service *TimesheetService) doForeignKeysExist(timesheet *admin.Timesheet, credentialID uuid.UUID) error {

	// check if tenant exists.
	err := service.doesTenantExist(timesheet.TenantID)
	if err != nil {
		return err
	}

	// check if credential exists for master fields.
	err = service.doesCredentialExist(timesheet.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if credential exists.
	err = service.doesCredentialExist(timesheet.TenantID, timesheet.CredentialID)
	if err != nil {
		return err
	}

	err = service.doesDepartmentExist(timesheet.TenantID, timesheet.DepartmentID)
	if err != nil {
		return err
	}

	if !timesheet.IsOnLeave {
		for _, activity := range timesheet.Activities {
			err = service.doesActivitiesForeignKeyExist(timesheet.TenantID, &activity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (service *TimesheetService) doesActivitiesForeignKeyExist(tenantID uuid.UUID, timesheetActivity *admin.TimesheetActivity) error {

	if timesheetActivity.ProjectID != nil {
		err := service.doesProjectExist(tenantID, timesheetActivity.ProjectID)
		if err != nil {
			return err
		}
	}

	if timesheetActivity.SubProjectID != nil {
		err := service.doesProjectExist(tenantID, timesheetActivity.SubProjectID)
		if err != nil {
			return err
		}

		err = service.doesSubProjectExist(tenantID, timesheetActivity.SubProjectID, timesheetActivity.ProjectID)
		if err != nil {
			return err
		}
	}

	if timesheetActivity.BatchID != nil {
		err := service.doesBatchExist(tenantID, timesheetActivity.BatchID)
		if err != nil {
			return err
		}
	}

	if timesheetActivity.BatchSessionID != nil {
		err := service.doesSessionExist(tenantID, timesheetActivity.BatchSessionID)
		if err != nil {
			return err
		}
	}

	// if timesheetActivity.BatchSessionID != nil {
	// 	err := service.doesBatchTopicExist(tenantID, timesheetActivity.BatchSessionID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// updateTimesheetActivities will update timesheet activities
func (service *TimesheetService) updateTimesheetActivities(uow *repository.UnitOfWork, timesheet *admin.Timesheet) error {

	tempTimesheetActivities := []admin.TimesheetActivity{}
	activityMap := make(map[uuid.UUID]uint)

	// get all activities for current timesheet
	err := service.Repository.GetAllForTenant(uow, timesheet.TenantID, &tempTimesheetActivities,
		repository.Filter("`timesheet_id` = ?", timesheet.ID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// populating activityMap
	for _, tempActivity := range tempTimesheetActivities {
		activityMap[tempActivity.ID] = 1
	}

	for _, activity := range timesheet.Activities {
		if util.IsUUIDValid(activity.ID) {
			activityMap[activity.ID]++
		}
		activity.TimesheetID = timesheet.ID
		activity.TenantID = timesheet.TenantID

		// check if activity already exists in the DB
		if activityMap[activity.ID] > 1 {
			timesheetActivityInDB := admin.TimesheetActivity{}
			err = service.Repository.GetRecordForTenant(uow, timesheet.TenantID, &timesheetActivityInDB,
				repository.Filter("`id` = ?", activity.ID), repository.Select([]string{"`created_by`"}))
			if err != nil {
				uow.RollBack()
				return err
			}

			activity.UpdatedBy = timesheet.UpdatedBy
			activity.CreatedBy = timesheetActivityInDB.CreatedBy
			if activity.NextEstimatedDate != nil {
				err = service.addNextEstimateTimesheet(uow, timesheet, activity)
				if err != nil {
					log.NewLogger().Error(err.Error())
					return err
				}
			}
			err = service.Repository.Save(uow, &activity)
			if err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
			activityMap[activity.ID] = 0
		}

		// add activity when uuid is nil
		if !util.IsUUIDValid(activity.ID) {
			activity.CreatedBy = timesheet.UpdatedBy

			// // should not be allowed in frontend (uncomment if allowed) #Niranjan.
			// if activity.NextEstimatedDate != nil {
			// 	err = service.addNextEstimateTimesheet(uow, timesheet, activity)
			// 	if err != nil {
			// 		log.NewLogger().Error(err.Error())
			// 		return err
			// 	}
			// }
			err := service.Repository.Add(uow, &activity)
			if err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
		}
	}

	// deleting all records where count is 1 as they have been removed from the activites
	for _, activity := range tempTimesheetActivities {
		if activityMap[activity.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, &activity, map[string]interface{}{
				"DeletedBy": timesheet.UpdatedBy,
				"DeletedAt": time.Now(),
			})
			if err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
		}
		activityMap[activity.ID] = 0
	}

	timesheet.Activities = nil
	return nil
}

// assignCreatedBy will assign createdby to activites
func (service *TimesheetService) assignCreatedBy(timesheet *admin.Timesheet) {
	for index := range timesheet.Activities {
		timesheet.Activities[index].CreatedBy = timesheet.CreatedBy
		timesheet.Activities[index].TenantID = timesheet.TenantID
	}
}

// returns error if there is no tenant record in table.
func (service *TimesheetService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *TimesheetService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no department record in table for the given tenant.
func (service *TimesheetService) doesDepartmentExist(tenantID, departmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Department{},
		repository.Filter("`id` = ?", departmentID))
	if err := util.HandleError("Invalid department ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no project record in table for the given tenant.
func (service *TimesheetService) doesProjectExist(tenantID uuid.UUID, projectID *uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.SwabhavProject{},
		repository.Filter("`id` = ?", projectID))
	if err := util.HandleError("Invalid project ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no project record in table for the given tenant.
func (service *TimesheetService) doesSubProjectExist(tenantID uuid.UUID, subProjectID, projectID *uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.SwabhavProject{},
		repository.Filter("`id` = ? AND `project_id` = ?", subProjectID, projectID))
	if err := util.HandleError("Project does not have the specified sub-project", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *TimesheetService) doesBatchExist(tenantID uuid.UUID, batchID *uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch topic record in table for the given tenant.
// func (service *TimesheetService) doesBatchTopicExist(tenantID uuid.UUID, batchTopicID *uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.BatchTopic{},
// 		repository.Filter("`id` = ?", batchTopicID))
// 	if err := util.HandleError("Invalid batch topic ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// returns error if there is no session record in table for the given tenant.
func (service *TimesheetService) doesSessionExist(tenantID uuid.UUID, batchSessionID *uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Session{},
		repository.Filter("`id` = ?", batchSessionID))
	if err := util.HandleError("Invalid session ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// // returns error if there is no time_sheet record in table for the given tenant.
// func (service *TimesheetService) doesTimesheetExist(tenantID, timesheetID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.Timesheet{},
// 		repository.Filter("`id` = ?", timesheetID))
// 	if err := util.HandleError("Invalid Timesheet ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// // returns error if there is no timesheet_activity record in table for the given tenant.
// func (service *TimesheetService) doesTimesheetActivityExist(tenantID, timesheetActivityID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, admin.TimesheetActivity{},
// 		repository.Filter("`id` = ?", timesheetActivityID))
// 	if err := util.HandleError("Invalid Timesheet Activity ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

func (service *TimesheetService) addNextEstimateTimesheet(uow *repository.UnitOfWork,
	timesheet *admin.Timesheet, timesheetActivity admin.TimesheetActivity) error {
	newTimesheet := &admin.Timesheet{}

	newTimesheet.TenantID = timesheet.TenantID
	newTimesheet.DepartmentID = timesheet.DepartmentID
	newTimesheet.CredentialID = timesheet.CredentialID
	newTimesheet.CreatedBy = timesheetActivity.UpdatedBy
	newTimesheet.Date = *timesheetActivity.NextEstimatedDate

	newTimesheetActivity := timesheetActivity
	// Make uuid nil because the add operation won't work on child if ID is passed.
	newTimesheetActivity.ID = uuid.Nil
	newTimesheetActivity.IsCompleted = nil
	newTimesheetActivity.NextEstimatedDate = nil
	newTimesheetActivity.WorkDone = nil

	newTimesheet.Activities = append(newTimesheet.Activities, newTimesheetActivity)

	err := service.AddTimesheet(newTimesheet, uow)
	if err != nil {
		uow.RollBack()
		return err
	}
	// #weird
	// uow.Commit()
	// return nil
	return nil
}

// addSearchQueries will form queries on the basis of search specified in form.
func (service *TimesheetService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if toDate, ok := requestForm["toDate"]; ok {
		util.AddToSlice("CAST(`timesheets`.`date` AS DATE)", "<= ?", "AND", toDate, &columnNames, &conditions, &operators, &values)
	}
	if fromDate, ok := requestForm["fromDate"]; ok {
		util.AddToSlice("CAST(`timesheets`.`date` AS DATE)", ">= ?", "AND", fromDate, &columnNames, &conditions, &operators, &values)
	}
	if departmentID, ok := requestForm["departmentID"]; ok {
		util.AddToSlice("`timesheets`.`department_id`", "= ?", "AND", departmentID, &columnNames, &conditions, &operators, &values)
	}
	if credentialID, ok := requestForm["credentialID"]; ok {
		util.AddToSlice("`timesheets`.`credential_id`", "= ?", "AND", credentialID, &columnNames, &conditions, &operators, &values)
	}
	if projectID, ok := requestForm["projectID"]; ok {
		util.AddToSlice("`timesheet_activities`.`project_id`", "= ?", "AND", projectID, &columnNames, &conditions, &operators, &values)
	}
	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("`timesheet_activities`.`batch_id`", "= ?", "AND", batchID, &columnNames, &conditions, &operators, &values)
	}
	if batchSessionID, ok := requestForm["batchSessionID"]; ok {
		util.AddToSlice("`timesheet_activities`.`batch_session_id`", "= ?", "AND", batchSessionID, &columnNames, &conditions, &operators, &values)
	}
	// if BatchTopicID, ok := requestForm["BatchTopicID"]; ok {
	// 	util.AddToSlice("`timesheet_activities`.`batch_topic_id`", "= ?", "AND", BatchTopicID, &columnNames, &conditions, &operators, &values)
	// }
	if isCompleted, ok := requestForm["isCompleted"]; ok {
		// when isCompleted is null
		if isCompleted[0] == "null" {
			util.AddToSlice("`timesheet_activities`.`is_completed`", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		} else {
			// convert to bool (Pass 1 or 0)
			util.AddToSlice("`timesheet_activities`.`is_completed`", "= ?", "AND", isCompleted, &columnNames, &conditions, &operators, &values)
		}
	}
	if isOnLeave, ok := requestForm["isOnLeave"]; ok {
		util.AddToSlice("`timesheets`.`is_on_leave`", "= ?", "AND", isOnLeave, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// addIsCompletedSearchQueries will form query for isCompleted field.
func (service *TimesheetService) addIsCompletedSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if isCompleted, ok := requestForm["isCompleted"]; ok {
		// when isCompleted is null
		if isCompleted[0] == "null" {
			util.AddToSlice("`timesheet_activities`.`is_completed`", "IS NULL", "AND", nil, &columnNames, &conditions, &operators, &values)
		} else {
			// convert to bool (Pass 1 or 0)
			util.AddToSlice("`timesheet_activities`.`is_completed`", "= ?", "AND", isCompleted, &columnNames, &conditions, &operators, &values)
		}
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// DeleteTimesheetActivity will delete specific timesheet activity
// func (service *TimesheetService) DeleteTimesheetActivity(timesheetActivity *admin.TimesheetActivity) error {

// 	// check if tenant exists.
// 	err := service.doesTenantExist(timesheetActivity.TenantID)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if credential exists for master fields.
// 	err = service.doesCredentialExist(timesheetActivity.TenantID, timesheetActivity.UpdatedBy)
// 	if err != nil {
// 		return err
// 	}

// 	uow := repository.NewUnitOfWork(service.DB, false)

// 	tempTimesheet := admin.Timesheet{}
// 	err = service.Repository.GetRecordForTenant(uow, timesheetActivity.TenantID, &tempTimesheet,
// 		repository.Filter("`id` = ?", timesheetActivity.TimesheetID), repository.Select([]string{"`credential_id`"}))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}
// 	if timesheetActivity.UpdatedBy != tempTimesheet.CredentialID {
// 		log.NewLogger().Error("Not authorized to delete timesheet.")
// 		return errors.NewValidationError("Not authorized to delete timesheet.")
// 	}

// 	// soft-delete timesheet_activities
// 	err = service.Repository.UpdateWithMap(uow, &admin.TimesheetActivity{}, map[string]interface{}{
// 		"DeletedAt": time.Now(),
// 		"DeletedBy": timesheetActivity.DeletedBy,
// 	}, repository.Filter("`id` = ?", timesheetActivity.ID))
// 	if err != nil {
// 		uow.RollBack()
// 		return err
// 	}

// 	return nil
// }
