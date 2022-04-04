package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// DegreeService provide method to update, delete, add, get method for Degree.
type DegreeService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewDegreeService returns new instance of DegreeService.
func NewDegreeService(db *gorm.DB, repository repository.Repository) *DegreeService {
	return &DegreeService{
		DB:         db,
		Repository: repository,
	}
}

// AddDegree add new degree to database.
func (service *DegreeService) AddDegree(degree *general.Degree, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(degree.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesDegreeNameExist(degree)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(degree.CreatedBy, degree.TenantID)
	if err != nil {
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add Degree.
	err = service.Repository.Add(uow, degree)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddDegrees adds multiple degrees to database.
func (service *DegreeService) AddDegrees(degrees *[]general.Degree, degreeIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same name conflict.
	for i := 0; i < len(*degrees); i++ {
		for j := 0; j < len(*degrees); j++ {
			if i != j && (*degrees)[i].Name == (*degrees)[j].Name {
				log.NewLogger().Error("Name:" + (*degrees)[j].Name + " exists")
				return errors.NewValidationError("Name:" + (*degrees)[j].Name + " exists")
			}
		}
	}

	// Add individual degree.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, degree := range *degrees {
		degree.TenantID = tenantID
		degree.CreatedBy = credentialID
		err := service.AddDegree(&degree, uow)
		if err != nil {
			return err
		}
		*degreeIDs = append(*degreeIDs, degree.ID)
	}

	uow.Commit()
	return nil
}

// UpdateDegree update degree to database.
func (service *DegreeService) UpdateDegree(degree *general.Degree) error {
	// Validate tenant ID.
	err := service.doesTenantExist(degree.TenantID)
	if err != nil {
		return err
	}

	// Validate degree ID.
	err = service.doesDegreeExist(degree.ID, degree.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(degree.UpdatedBy, degree.TenantID)
	if err != nil {
		return err
	}

	// Validate if same city name exists.
	err = service.doesDegreeNameExist(degree)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update degree
	err = service.Repository.Update(uow, degree)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteDegree deletes degree from database.
func (service *DegreeService) DeleteDegree(degree *general.Degree) error {
	credentialID := degree.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(degree.TenantID)
	if err != nil {
		return err
	}

	// Validate degree ID.
	err = service.doesDegreeExist(degree.ID, degree.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, degree.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// ******************************************Delete specializations****************************************
	// Update specializations for updating deleted_by and deleted_at fields of specializations
	if err := service.Repository.UpdateWithMap(uow, &general.Specialization{}, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=? AND `degree_id`=?", degree.TenantID, degree.ID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Degree could not be deleted", http.StatusInternalServerError)
	}

	// Update degree for updating deleted_by and deleted_at fields of degree
	if err := service.Repository.UpdateWithMap(uow, degree, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", degree.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Degree could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetDegreeList returns all degree list.
func (service *DegreeService) GetDegreeList(degrees *[]general.Degree, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all degrees
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, degrees, "`name`")
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetDegree returns all degrees.
func (service *DegreeService) GetDegree(degree *general.Degree) error {
	// Validate tenant id.
	err := service.doesTenantExist(degree.TenantID)
	if err != nil {
		return err
	}

	// Validate degree id.
	err = service.doesDegreeExist(degree.ID, degree.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one degree by id.
	err = service.Repository.GetForTenant(uow, degree.TenantID, degree.ID, degree)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetDegrees returns specific degree by id.
func (service *DegreeService) GetDegrees(degrees *[]general.Degree, tenantID uuid.UUID, parser *web.Parser, totalCount *int) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get all degrees form database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, degrees, "`name`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search queries.
func (service *DegreeService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	if _, ok := requestForm["name"]; ok {
		return repository.Filter("`name` LIKE ?", "%"+requestForm.Get("name")+"%")
	}
	return nil
}

// doesDegreeNameExist returns true if the degree name already exists for the state in database.
func (service *DegreeService) doesDegreeNameExist(degree *general.Degree) error {
	// Check for same degree name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, degree.TenantID, &general.Degree{},
		repository.Filter("`name`=? AND `id`!=?", degree.Name, degree.ID))
	if err := util.HandleIfExistsError("Name:"+degree.Name+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *DegreeService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesDegreeExist validates if degree exists or not in database.
func (service *DegreeService) doesDegreeExist(degreeID uuid.UUID, tenantID uuid.UUID) error {
	// Check degree exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Degree{},
		repository.Filter("`id` = ?", degreeID))
	if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *DegreeService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
