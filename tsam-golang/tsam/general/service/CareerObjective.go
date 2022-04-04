package service

import (
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// CareerObjectiveService provides methods to update, delete, add, get, get all and get all by degree for career objective.
type CareerObjectiveService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCareerObjectiveService returns new instance of CareerObjectiveService.
func NewCareerObjectiveService(db *gorm.DB, repository repository.Repository) *CareerObjectiveService {
	return &CareerObjectiveService{
		DB:         db,
		Repository: repository,
	}
}

// AddCareerObjective adds new career objective in database.
func (service *CareerObjectiveService) AddCareerObjective(careerObjective *general.CareerObjective) error {
	// Get credential id from CreatedBy field of career objective(set in controller).
	credentialID := careerObjective.CreatedBy

	// Give tenant id and craeted_by to all career objectives courses of career objective.
	if careerObjective.Courses != nil && len(careerObjective.Courses) != 0 {
		for i := 0; i < len(careerObjective.Courses); i++ {
			careerObjective.Courses[i].TenantID = careerObjective.TenantID
			careerObjective.Courses[i].CreatedBy = credentialID
		}
	}

	// Validate tenant id.
	if err := service.doesTenantExist(careerObjective.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, careerObjective.TenantID); err != nil {
		return err
	}

	// Validate course id.
	if err := service.doesCourseExist(careerObjective); err != nil {
		return err
	}

	// Check if branch name exists or not.
	err := service.doesNameExist(careerObjective)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add career objective to database.
	if err := service.Repository.Add(uow, careerObjective); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Career Objective could not be added", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetCareerObjectives gets all career objectives from database.
func (service *CareerObjectiveService) GetCareerObjectives(careerObjectives *[]general.CareerObjective, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get career objectives from database.
	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, careerObjectives, "`name`",
		service.addSearchQueries(parser.Form),

		repository.PreloadAssociations([]string{"Courses"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Sort the courses by order fields.
	service.sortCoursesByOrder(careerObjectives)

	uow.Commit()
	return nil
}

// UpdateCareerObjective updates one career objective by specific career objective id in database.
func (service *CareerObjectiveService) UpdateCareerObjective(careerObjective *general.CareerObjective) error {
	// Get credential id from UpdatedBy field of career objective(set in controller).
	credentialID := careerObjective.UpdatedBy

	// Give tenant id and craeted_by to all career objectives courses of career objective.
	if careerObjective.Courses != nil && len(careerObjective.Courses) != 0 {
		for i := 0; i < len(careerObjective.Courses); i++ {
			careerObjective.Courses[i].TenantID = careerObjective.TenantID
		}
	}

	// Validate tenant id.
	if err := service.doesTenantExist(careerObjective.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, careerObjective.TenantID); err != nil {
		return err
	}

	// Validate foreign key ids.
	if err := service.doesCourseExist(careerObjective); err != nil {
		return err
	}

	// Check if career objective name exists or not.
	err := service.doesNameExist(careerObjective)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.updateCareerObjectiveCourse(uow, careerObjective.Courses, careerObjective.TenantID,
		careerObjective.UpdatedBy, careerObjective.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make courses nil so that it is not updated, added, or craeted with career objective update
	careerObjective.Courses = nil

	// Update career objective.
	if err := service.Repository.Update(uow, careerObjective); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Career Objective could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCareerObjective deletes one career objective by specific career objective id from database.
func (service *CareerObjectiveService) DeleteCareerObjective(careerObjective *general.CareerObjective) error {
	// Get credential id from DeletedBy field of talent(set in controller).
	credentialID := careerObjective.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(careerObjective.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, careerObjective.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCareerObjectiveExist(careerObjective.ID, careerObjective.TenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update career objective course for updating deleted_by and deleted_at fields of career objective course.
	if err := service.Repository.UpdateWithMap(uow, &general.CareerObjectivesCourse{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=? AND `career_objective_id`=?", careerObjective.TenantID, careerObjective.ID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Career Objective could not be deleted", http.StatusInternalServerError)
	}

	// Update career objective for updating deleted_by and deleted_at fields of career objective.
	if err := service.Repository.UpdateWithMap(uow, careerObjective, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", careerObjective.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Career Objective could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// updateCareerObjectiveCourse will update the career objective course for specified career objective.
func (service *CareerObjectiveService) updateCareerObjectiveCourse(uow *repository.UnitOfWork, courses []general.CareerObjectivesCourse,
	tenantID, credentialID, careerObjectiveID uuid.UUID) error {

	// If current courses is not present and previous courses is present.
	if courses == nil {
		err := service.Repository.UpdateWithMap(uow, general.CareerObjectivesCourse{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`career_objective_id`=?", careerObjectiveID))
		if err != nil {
			return err
		}
	}

	// Get previous courses from database.
	tempCourses := []general.CareerObjectivesCourse{}
	err := service.Repository.GetAllForTenant(uow, tenantID, &tempCourses,
		repository.Filter("`career_objective_id`=?", careerObjectiveID))
	if err != nil {
		return err
	}

	// Make course map for keeping count of occurrences of course id in previous and current courses.
	courseMap := make(map[uuid.UUID]uint)

	// Get count of course ids of previous course.
	for _, tempCourse := range tempCourses {
		courseMap[tempCourse.ID]++
	}

	// Compare with current courses for duplicate or no occurrences.
	for _, course := range courses {
		if util.IsUUIDValid(course.ID) {
			courseMap[course.ID]++
		} else { // If uuid is nil then give it created_by field.
			course.CreatedBy = credentialID
			course.TenantID = tenantID
			course.CareerObjectiveID = careerObjectiveID
			err = service.Repository.Add(uow, &course)
			if err != nil {
				return err
			}
		}
		// If duplicate occurrence then update it.
		if courseMap[course.ID] > 1 {
			course.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &course)
			if err != nil {
				return err
			}
			courseMap[course.ID] = 0
		}
	}

	// If one occurrece then delete it.
	for _, tempCourse := range tempCourses {
		if courseMap[tempCourse.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, general.CareerObjectivesCourse{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempCourse.ID))
			if err != nil {
				return err
			}
			courseMap[tempCourse.ID] = 0
		}
	}
	return nil
}

// addSearchQueries adds search queries.
func (service *CareerObjectiveService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	if _, ok := requestForm["name"]; ok {
		return repository.Filter("`name` LIKE ?", "%"+requestForm.Get("name")+"%")
	}
	return nil
}

// doesCourseExist validates if course exists or not in database.
func (service *CareerObjectiveService) doesCourseExist(careerObjective *general.CareerObjective) error {
	if careerObjective.Courses != nil && len(careerObjective.Courses) != 0 {
		for _, careerObjectiveCourse := range careerObjective.Courses {
			exists, err := repository.DoesRecordExistForTenant(service.DB, careerObjective.TenantID, general.Course{},
				repository.Filter("`id` = ?", careerObjectiveCourse.CourseID))
			if err := util.HandleError("Invalid course ID", exists, err); err != nil {
				return err
			}
		}
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CareerObjectiveService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCareerObjectiveExist if career objective exists or not in database.
func (service *CareerObjectiveService) doesCareerObjectiveExist(careerObjectiveID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.CareerObjective{},
		repository.Filter("`id` = ?", careerObjectiveID))
	if err := util.HandleError("Invalid career objective ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *CareerObjectiveService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{}, repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesNameExist checks if career objective's name already exists or not.
func (service *CareerObjectiveService) doesNameExist(careerObjective *general.CareerObjective) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, careerObjective.TenantID, &general.CareerObjective{},
		repository.Filter("`name`=? AND `id`!=?", careerObjective.Name, careerObjective.ID))
	if err := util.HandleIfExistsError("Name exists", exists, err); err != nil {
		return err
	}
	return nil
}

// sortCoursesByOrder sorts all the courses by order field.
func (service *CareerObjectiveService) sortCoursesByOrder(careerObjectives *[]general.CareerObjective) {
	if careerObjectives != nil && len(*careerObjectives) != 0 {
		for i := 0; i < len(*careerObjectives); i++ {
			courses := &(*careerObjectives)[i].Courses
			for j := 0; j < len(*courses); j++ {
				sort.Slice(*courses, func(p, q int) bool {
					return (*courses)[p].Order < (*courses)[q].Order
				})
			}
		}
	}
}
