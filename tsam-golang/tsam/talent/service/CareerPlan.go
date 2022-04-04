package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// CareerPlanService provides methods to update, delete, add, get, get all and get all by degree for career plan.
type CareerPlanService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewCareerPlanService returns new instance of CareerPlanService.
func NewCareerPlanService(db *gorm.DB, repository repository.Repository) *CareerPlanService {
	return &CareerPlanService{
		DB:         db,
		Repository: repository,
	}
}

// AddCareerPlans adds new career plan in database.
func (service *CareerPlanService) AddCareerPlans(careerPlans *[]talent.CareerPlan, tenantID uuid.UUID,
	talentID uuid.UUID, credentialID uuid.UUID) error {
	//give tenant id and credential id to all career plans
	for i := 0; i < len(*careerPlans); i++ {
		(*careerPlans)[i].TenantID = tenantID
		(*careerPlans)[i].CreatedBy = credentialID
		(*careerPlans)[i].TalentID = talentID
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, tenantID); err != nil {
		return err
	}

	// Validate course id.
	if err := service.doesTalentExist(talentID, tenantID); err != nil {
		return err
	}

	// Validate foreign key of all career plans.
	for _, careerPlan := range *careerPlans {
		if err := service.doForeignKeysExist(&careerPlan); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Check career objective exists or not for the talent.
	if careerPlans != nil && len(*careerPlans) != 0 {
		careerOnjectiveID := (*careerPlans)[0].CareerObjectiveID

		exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.CareerPlan{},
			repository.Filter("`career_objective_id`=? AND `talent_id`=?", careerOnjectiveID, talentID))
		if err := util.HandleIfExistsError("Career Objective already exists", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Add career plans to database.
	for _, careerPlan := range *careerPlans {
		if err := service.Repository.Add(uow, &careerPlan); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Career Plan could not be added", http.StatusInternalServerError)
		}
	}

	uow.Commit()
	return nil
}

// GetCareerPlans gets all career plans from database.
func (service *CareerPlanService) GetCareerPlans(careerPlans *[]talent.CareerPlan, tenantID uuid.UUID, talentID uuid.UUID) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(talentID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get career objectives from database.
	err := service.Repository.GetAllForTenant(uow, tenantID, careerPlans,
		repository.Filter("`talent_id`=?", talentID))
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	uow.Commit()
	return nil
}

// UpdateCareerPlan updates one career plan by specific career plan id in database.
func (service *CareerPlanService) UpdateCareerPlan(careerPlan *talent.CareerPlan) error {
	// Get credential id from UpdatedBy field of career plan(set in controller).
	credentialID := careerPlan.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(careerPlan.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, careerPlan.TenantID); err != nil {
		return err
	}

	// Validate course id.
	if err := service.doesTalentExist(careerPlan.TalentID, careerPlan.TenantID); err != nil {
		return err
	}

	// Validate foreign key.
	if err := service.doForeignKeysExist(careerPlan); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update faculty, current rating and updated_by fields of career plan.
	if err := service.Repository.UpdateWithMap(uow, &talent.CareerPlan{}, map[string]interface{}{
		"UpdatedBy":     credentialID,
		"FacultyID":     careerPlan.FacultyID,
		"CurrentRating": careerPlan.CurrentRating,
	},
		repository.Filter("`id`=? AND `tenant_id`=?", careerPlan.ID, careerPlan.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Career Plan could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteCareerPlan deletes one career plan by specific career plan id from database.
func (service *CareerPlanService) DeleteCareerPlan(careerObjectiveID uuid.UUID, tenantID uuid.UUID, credentialID uuid.UUID) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update career plan for updating deleted_by and deleted_at fields of career plan.
	if err := service.Repository.UpdateWithMap(uow, &talent.CareerPlan{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=? AND `career_objective_id`=?", tenantID, careerObjectiveID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Career Plan could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *CareerPlanService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *CareerPlanService) doesTalentExist(talentID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tal exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, talent.Talent{}, repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credental exists or not in database.
func (service *CareerPlanService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Credential{}, repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates if all foreign jey ids exists or not in database.
func (service *CareerPlanService) doForeignKeysExist(careerPlan *talent.CareerPlan) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check career objective exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, careerPlan.TenantID, general.CareerObjective{},
		repository.Filter("`id` = ?", careerPlan.CareerObjectiveID))
	if err := util.HandleError("Invalid career objective ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check career objectives courses exists or not.
	exists, err = repository.DoesRecordExistForTenant(uow.DB, careerPlan.TenantID, general.CareerObjectivesCourse{},
		repository.Filter("`id` = ?", careerPlan.CareerObjectivesCoursesID))
	if err := util.HandleError("Invalid career objectives courses ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check faculty exists or not.
	exists, err = repository.DoesRecordExistForTenant(uow.DB, careerPlan.TenantID, faculty.Faculty{},
		repository.Filter("`id` = ?", careerPlan.FacultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}
