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

// TechnologyService provide method to update, delete, add, get method for Technology.
type TechnologyService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTechnologyService creates new instance of TechnologyService.
func NewTechnologyService(db *gorm.DB, repository repository.Repository) *TechnologyService {
	return &TechnologyService{
		DB:         db,
		Repository: repository,
	}
}

// AddTechnology adds new technology to database.
func (service *TechnologyService) AddTechnology(technology *general.Technology, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(technology.TenantID)
	if err != nil {
		return err
	}

	// Validate if same language exists.
	err = service.doesLanguageExist(technology)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(technology.CreatedBy, technology.TenantID)
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

	// Add technology.
	err = service.Repository.Add(uow, technology)
	if err != nil {
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddTechnologies adds multiple technologies to database.
func (service *TechnologyService) AddTechnologies(technologies *[]general.Technology, technologyIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// Check for same language conflict.
	for i := 0; i < len(*technologies); i++ {
		for j := 0; j < len(*technologies); j++ {
			if i != j && (*technologies)[i].Language == (*technologies)[j].Language {
				log.NewLogger().Error("Language:" + (*technologies)[j].Language + " exists")
				return errors.NewValidationError("Language:" + (*technologies)[j].Language + " exists")
			}
		}
	}

	// Add individual technology.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, technology := range *technologies {
		technology.TenantID = tenantID
		technology.CreatedBy = credentialID
		err := service.AddTechnology(&technology, uow)
		if err != nil {
			return err
		}
		*technologyIDs = append(*technologyIDs, technology.ID)
	}

	uow.Commit()
	return nil
}

// UpdateTechnology updates technology in database.
func (service *TechnologyService) UpdateTechnology(technology *general.Technology) error {
	// Validate tenant ID.
	err := service.doesTenantExist(technology.TenantID)
	if err != nil {
		return err
	}

	// Validate technology ID.
	err = service.doesTechnologyExist(technology.ID, technology.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(technology.UpdatedBy, technology.TenantID)
	if err != nil {
		return err
	}

	// Validate if same language exists.
	err = service.doesLanguageExist(technology)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Update(uow, technology)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTechnology delete technology from database.
func (service *TechnologyService) DeleteTechnology(technology *general.Technology) error {
	credentialID := technology.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(technology.TenantID)
	if err != nil {
		return err
	}

	// Validate technology ID.
	err = service.doesTechnologyExist(technology.ID, technology.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, technology.TenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update technology for updating deleted_by and deleted_at fields of technology
	if err := service.Repository.UpdateWithMap(uow, technology, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", technology.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Technology could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetTechnology returns one technology by id.
func (service *TechnologyService) GetTechnology(technology *general.Technology, tenantID, technologyID uuid.UUID) error {

	// reason??
	// *technology = append(*technology, general.Technology{TenantBase: general.TenantBase{ID: technologyID}})

	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}
	technology.ID = technologyID

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetForTenant(uow, tenantID, technologyID, technology)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetTechnologies returns all technologies.
func (service *TechnologyService) GetTechnologies(technologies *[]general.Technology, parser *web.Parser, tenantID uuid.UUID, totalCount *int) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get all technologies.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, technologies, "`language`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTechnologyList returns all the technologies(without pagination).
func (service *TechnologyService) GetTechnologyList(technologies *[]general.Technology, tenantID uuid.UUID,
	requestForm url.Values) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all technologies.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, technologies, "`language`",
		service.addSearchQueries(requestForm))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries by comparing with the technology data.
func (service *TechnologyService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if _, ok := requestForm["language"]; ok {
		util.AddToSlice("`language`", "LIKE ?", "AND", "%"+requestForm.Get("language")+"%", &columnNames, &conditions, &operators, &values)
	}
	if rating, ok := requestForm["rating"]; ok {
		util.AddToSlice("`rating`", "= ?", "AND", rating, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesLanguageExist returns true if the language already exists for the technology in database.
func (service *TechnologyService) doesLanguageExist(technology *general.Technology) error {
	// Check for same language conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, technology.TenantID, &general.Technology{},
		repository.Filter("`language`=? AND `id`!=?", technology.Language, technology.ID))
	if err := util.HandleIfExistsError("Language:"+technology.Language+" exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TechnologyService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTechnologyExist validates if technology exists or not in database.
func (service *TechnologyService) doesTechnologyExist(technologyID uuid.UUID, tenantID uuid.UUID) error {
	// Check technology exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id` = ?", technologyID))
	if err := util.HandleError("Invalid technology ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *TechnologyService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}
