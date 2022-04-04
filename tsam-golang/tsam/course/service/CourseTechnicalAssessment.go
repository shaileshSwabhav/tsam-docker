package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// AssessmentService Provide method to Update, Delete, Add, Get Method For CourseTechnicalAssessment.
type AssessmentService struct {
	DB           *gorm.DB
	Repository   repository.Repository
	associations []string
}

// NewCourseTechnicalAssessmentService returns new instance of CourseTechnicalAssessmentService.
func NewCourseTechnicalAssessmentService(db *gorm.DB, repository repository.Repository) *AssessmentService {
	return &AssessmentService{
		DB:         db,
		Repository: repository,
		associations: []string{
			"Course", "Faculty",
		},
	}
}

// AddAssessment will add single assessment to the table.
func (service *AssessmentService) AddAssessment(uow *repository.UnitOfWork,
	assessment *course.CourseTechnicalAssessment) error {

	// extract ID's from the struct.
	// service.extractCourseTechnicalAssessmentID(assessment)

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(assessment, assessment.CreatedBy)
	if err != nil {
		return err
	}

	err = service.Repository.Add(uow, assessment)
	if err != nil {
		return err
	}

	return nil
}

// AddAssessments will add multiple course technical assessments to the table.
func (service *AssessmentService) AddAssessments(assessments *[]course.CourseTechnicalAssessment,
	tenantID, credentialID, facultyID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, false)

	for index := range *assessments {

		(*assessments)[index].TenantID = tenantID
		(*assessments)[index].CreatedBy = credentialID
		(*assessments)[index].FacultyID = facultyID

		err := service.AddAssessment(uow, &(*assessments)[index])
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// UpdateAssessment will update specified assessment in the table.
func (service *AssessmentService) UpdateAssessment(assessment *course.CourseTechnicalAssessment) error {

	// extract ID's from the struct.
	// service.extractCourseTechnicalAssessmentID(assessment)

	// Check if foreign keys exist.
	err := service.doesForeignKeyExist(assessment, assessment.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if assessment exist.
	err = service.doesAssessmentExist(assessment.TenantID, assessment.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempAssessment := course.CourseTechnicalAssessment{}
	err = service.Repository.GetRecordForTenant(uow, assessment.TenantID, &tempAssessment,
		repository.Filter("`id` = ?", assessment.ID), repository.Select("`created_by`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	assessment.CreatedBy = tempAssessment.CreatedBy

	err = service.Repository.Save(uow, assessment)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteAssessment will delete specified assessment from in the table.
func (service *AssessmentService) DeleteAssessment(assessment *course.CourseTechnicalAssessment) error {

	// Check if tenant exist.
	err := service.doesTenantExist(assessment.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(assessment.TenantID, assessment.DeletedBy)
	if err != nil {
		return err
	}

	// Check if assessment exist.
	err = service.doesAssessmentExist(assessment.TenantID, assessment.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, course.CourseTechnicalAssessment{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": assessment.DeletedBy,
	}, repository.Filter("`id` = ?", assessment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllAssessments will return all the technical assessments.
func (service *AssessmentService) GetAllAssessments(assessments *[]course.CourseTechnicalAssessmentDTO,
	tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAll(uow, assessments,
		repository.PreloadAssociations(service.associations),
		repository.Join("INNER JOIN courses ON courses.`id` = course_technical_assessments.`course_id` AND "+
			"course_technical_assessments.`tenant_id` = courses.`tenant_id`"),
		repository.Filter("course_technical_assessments.`tenant_id`=? AND courses.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("courses.`name`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAssessmentsForFaculty will return all the technical assessments for the specified faculty.
func (service *AssessmentService) GetAssessmentsForFaculty(assessments *[]course.CourseTechnicalAssessmentDTO,
	tenantID, facultyID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if faculty exist.
	err = service.doesFacultyExist(tenantID, facultyID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAll(uow, assessments,
		repository.Filter("course_technical_assessments.`faculty_id` = ?", facultyID), repository.PreloadAssociations([]string{"Course"}),
		repository.Join("INNER JOIN courses ON courses.`id` = course_technical_assessments.`course_id` AND "+
			"course_technical_assessments.`tenant_id` = courses.`tenant_id`"),
		repository.Filter("course_technical_assessments.`tenant_id`=? AND courses.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("courses.`name`"))
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

// doesForeignKeyExist will check if all the foreign keys are valid.
func (service *AssessmentService) doesForeignKeyExist(assessment *course.CourseTechnicalAssessment,
	credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(assessment.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(assessment.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if faculty exist.
	err = service.doesFacultyExist(assessment.TenantID, assessment.FacultyID)
	if err != nil {
		return err
	}

	// Check if course exist.
	err = service.doesCourseExist(assessment.TenantID, assessment.CourseID)
	if err != nil {
		return err
	}

	// Check if course rating is repeated for faculty.
	err = service.doesCourseAssessmentExist(assessment.TenantID, assessment.CourseID,
		assessment.FacultyID, assessment.ID)
	if err != nil {
		return err
	}

	return nil
}

// extractCourseTechnicalAssessmentID will extract ID's from struct.
// func (service *AssessmentService) extractCourseTechnicalAssessmentID(assessment *course.CourseTechnicalAssessment) {

// 	assessment.CourseID = assessment.CourseID
// }

// Returns error if there is no tenant record in table.
func (service *AssessmentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))

	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no credential record in table for the given tenant.
func (service *AssessmentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))

	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no faculty record in table for the given tenant.
func (service *AssessmentService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, faculty.Faculty{},
		repository.Filter("`id` = ?", facultyID))

	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no course record in table for the given tenant.
func (service *AssessmentService) doesCourseExist(tenantID, courseID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.Course{},
		repository.Filter("`id` = ?", courseID))

	if err := util.HandleError("Invalid course ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no course technical assessment record in table for the given tenant.
func (service *AssessmentService) doesAssessmentExist(tenantID, assessmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.CourseTechnicalAssessment{},
		repository.Filter("`id` = ?", assessmentID))

	if err := util.HandleError("Invalid assessment ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no course technical assessment record for course and faculty exist in table for the given tenant.
func (service *AssessmentService) doesCourseAssessmentExist(tenantID, courseID, facultyID, assessmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, course.CourseTechnicalAssessment{},
		repository.Filter("`course_id` = ? AND `faculty_id` = ? AND `id` != ?", courseID, facultyID, assessmentID))

	if err := util.HandleIfExistsError("Rating already assigned to course", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}
