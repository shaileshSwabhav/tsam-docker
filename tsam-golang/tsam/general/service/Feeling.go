package service

import (
	"net/url"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// FeelingService provides methods to do different CRUD operations on feeling table.
type FeelingService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewFeelingService returns a new instance Of FeelingService.
func NewFeelingService(db *gorm.DB, repository repository.Repository) *FeelingService {
	return &FeelingService{
		DB:         db,
		Repository: repository,
		association: []string{
			"FeelingLevels",
		},
	}
}

// AddFeeling will add the feeling and its feeling level.
func (service *FeelingService) AddFeeling(feeling *general.Feeling) error {

	// Check if tenant exist.
	err := service.doesTenantExist(feeling.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(feeling.TenantID, feeling.CreatedBy)
	if err != nil {
		return err
	}

	// Check if feeling name is unique.
	err = service.doesFeelingNameExist(feeling.TenantID, feeling.ID, feeling.FeelingName)
	if err != nil {
		return err
	}

	// Set tenant id and created_by field of all feeling levels of feeling.
	if feeling.FeelingLevels != nil {
		for index := range feeling.FeelingLevels {
			feeling.FeelingLevels[index].TenantID = feeling.TenantID
			feeling.FeelingLevels[index].CreatedBy = feeling.CreatedBy
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, feeling)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateFeeling will update the existing feeling.
func (service *FeelingService) UpdateFeeling(feeling *general.Feeling) error {

	// Check if tenant exist.
	err := service.doesTenantExist(feeling.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(feeling.TenantID, feeling.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if feeling exist.
	err = service.doesFeelingExist(feeling.TenantID, feeling.ID)
	if err != nil {
		return err
	}

	// Check if feeling name is unique.
	err = service.doesFeelingNameExist(feeling.TenantID, feeling.ID, feeling.FeelingName)
	if err != nil {
		return err
	}

	// Set tenant id and created_by field of all feeling levels of feeling.
	if feeling.FeelingLevels != nil {
		for index := range feeling.FeelingLevels {
			feeling.FeelingLevels[index].TenantID = feeling.TenantID
			feeling.FeelingLevels[index].CreatedBy = feeling.CreatedBy
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.updateFeelingLevel(uow, feeling.FeelingLevels, feeling.TenantID,
		feeling.UpdatedBy, feeling.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	feeling.FeelingLevels = nil

	err = service.Repository.Update(uow, feeling)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteFeeling will delete feeling and its feeling levels.
func (service *FeelingService) DeleteFeeling(feeling *general.Feeling) error {

	// Check if tenant exist.
	err := service.doesTenantExist(feeling.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExist(feeling.TenantID, feeling.DeletedBy)
	if err != nil {
		return err
	}

	// Check if feeling exist.
	err = service.doesFeelingExist(feeling.TenantID, feeling.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.UpdateWithMap(uow, general.FeelingLevel{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": feeling.DeletedBy,
	}, repository.Filter("`feeling_id`=?", feeling.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.UpdateWithMap(uow, general.Feeling{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": feeling.DeletedBy,
	}, repository.Filter("`id`=?", feeling.ID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetFeelingsList will return list of feelings from feelings table.
func (service *FeelingService) GetFeelingsList(feelings *[]general.Feeling, tenantID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feelings, "`feeling_name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetFeelingLevels will return all the feeling levels for the specified feeling.
func (service *FeelingService) GetFeelingLevels(feelingLevels *[]general.FeelingLevel,
	tenantID, feelingID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if feeling exist.
	err = service.doesFeelingExist(tenantID, feelingID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, feelingLevels, "`level_number`",
		repository.Filter("`feeling_id` = ?", feelingID))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllFeelings gets all feelings from database.
func (service *FeelingService) GetAllFeelings(feelings *[]general.Feeling, tenantID uuid.UUID,
	parser *web.Parser, totalCount *int) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get feelings from database.
	err := service.Repository.GetAllInOrderForTenant(uow, tenantID, feelings, "`feeling_name`",
		service.addSearchQueries(parser.Form),
		repository.PreloadAssociations([]string{"FeelingLevels"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		uow.RollBack()
		return err
	}

	// Sort the feeling levels by level number field.
	service.sortFeelingLevelsByLevelNumber(feelings)

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// updateFeelingLevel will update the feeling_level for specified feeling
func (service *FeelingService) updateFeelingLevel(uow *repository.UnitOfWork, feelingLevels []general.FeelingLevel,
	tenantID, credentialID, feelingID uuid.UUID) error {

	if feelingLevels == nil {
		err := service.Repository.UpdateWithMap(uow, general.FeelingLevel{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`feeling_id`=?", feelingID))
		if err != nil {
			return err
		}
	}

	tempFeelingLevels := []general.FeelingLevel{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempFeelingLevels,
		repository.Filter("`feeling_id`=?", feelingID))
	if err != nil {
		return err
	}

	feelingLevelMap := make(map[uuid.UUID]uint)

	for _, tempFeelingLevel := range tempFeelingLevels {
		feelingLevelMap[tempFeelingLevel.ID]++
	}

	for _, feelingLevel := range feelingLevels {

		if util.IsUUIDValid(feelingLevel.ID) {
			feelingLevelMap[feelingLevel.ID]++
		} else {
			feelingLevel.CreatedBy = credentialID
			feelingLevel.TenantID = tenantID
			feelingLevel.FeelingID = feelingID
			err = service.Repository.Add(uow, &feelingLevel)
			if err != nil {
				return err
			}
		}

		if feelingLevelMap[feelingLevel.ID] > 1 {
			feelingLevel.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &feelingLevel)
			if err != nil {
				return err
			}
			feelingLevelMap[feelingLevel.ID] = 0
		}
	}

	for _, tempFeelingLevel := range tempFeelingLevels {
		if feelingLevelMap[tempFeelingLevel.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, general.FeelingLevel{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempFeelingLevel.ID))
			if err != nil {
				return err
			}
			feelingLevelMap[tempFeelingLevel.ID] = 0
		}
	}
	return nil
}

// addSearchQueries adds search criteria.
func (service *FeelingService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if _, ok := requestForm["feelingName"]; ok {
		util.AddToSlice("`feeling_name`", "LIKE ?", "AND", "%"+requestForm.Get("feelingName")+"%", &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *FeelingService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *FeelingService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFeelingExist returns error if there is no feeling record in table.
func (service *FeelingService) doesFeelingExist(tenantID, feelingID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Feeling{},
		repository.Filter("`id` = ?", feelingID))
	if err := util.HandleError("Invalid feeling ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesFeelingNameExist returns error if same feeling name exists in table.
func (service *FeelingService) doesFeelingNameExist(tenantID, feelingID uuid.UUID, feelingName string) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Feeling{},
		repository.Filter("`feeling_name` = ? AND `id`!=?", feelingName, feelingID))
	if err := util.HandleIfExistsError("Feeling name already exists", exists, err); err != nil {
		return err
	}
	return nil
}

// sortFeelingLevelsByLevelNumber sorts all the feeling levels by level number field.
func (service *FeelingService) sortFeelingLevelsByLevelNumber(feelings *[]general.Feeling) {
	if feelings != nil && len(*feelings) != 0 {
		for i := 0; i < len(*feelings); i++ {
			feelingLevels := &(*feelings)[i].FeelingLevels
			for j := 0; j < len(*feelingLevels); j++ {
				sort.Slice(*feelingLevels, func(p, q int) bool {
					return (*feelingLevels)[p].LevelNumber < (*feelingLevels)[q].LevelNumber
				})
			}
		}
	}
}
