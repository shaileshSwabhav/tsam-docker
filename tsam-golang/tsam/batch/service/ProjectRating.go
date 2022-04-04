package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// ProjectRatingService provides methods to update, delete, add, get for project rating parameter.
type ProjectRatingService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProjectRatingService returns new instance of ProjectRatingService.
func NewProjectRatingService(db *gorm.DB, repository repository.Repository) *ProjectRatingService {
	return &ProjectRatingService{
		DB:         db,
		Repository: repository,
	}
}

// GetAllProgrammingProject will get all the project rating parameter from the table.
func (service *ProjectRatingService) GetAllProgrammingProject(projectRatingParameters *[]batch.ProgrammingProjectRatingParameter,
	tenantID uuid.UUID, totalCount *int, parser *web.Parser) error {

	// Check if tenant exists.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetAllForTenant(uow, tenantID, projectRatingParameters,
		repository.Filter("`is_active` = ?", 1))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// doesTenantExists validates tenantID.
func (ser *ProjectRatingService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}
