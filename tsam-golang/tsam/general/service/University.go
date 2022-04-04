package service

import (
	"fmt"
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

// UniversityService provides method like add, update, delete, get by ID, get all for University.
type UniversityService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewUniversityService returns the new instance of UniversityService.
func NewUniversityService(db *gorm.DB, repo repository.Repository) *UniversityService {
	return &UniversityService{
		DB:         db,
		Repository: repo,
	}
}

// AddUniversity adds new University to database.
func (service *UniversityService) AddUniversity(university *general.University,
	uows ...*repository.UnitOfWork) error {

	// Assign countryID if no ID is given.
	if university.CountryID == uuid.Nil {
		err := service.assignCountryIDByName(university)
		if err != nil {
			return err
		}
	}

	// Extract all foreign key IDs assign to entityID field.
	err := service.extractID(university)
	if err != nil {
		return err
	}

	// Check all foreign key records.
	err = service.doForeignKeysExist(university.TenantID, university.CreatedBy, university.CountryID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(university)
	if err != nil {
		return err
	}

	// Create new unit of work, if no transaction has been passed to the function.
	length := len(uows)
	if length == 0 {
		uows = append(uows, repository.NewUnitOfWork(service.DB, false))
	}
	uow := uows[0]

	// Add repo call.
	err = service.Repository.Add(uow, university)
	if err != nil {
		// Rollback only if no transaction is passed.
		if length == 0 {
			uow.RollBack()
		}
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	// Commit only if no transaction is passed.
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddUniversities adds multiple universities to Database.
func (service *UniversityService) AddUniversities(universities []*general.University) error {
	// Check for same university name conflict.
	for i := 0; i < len(universities); i++ {
		for j := i + 1; j < len(universities); j++ {
			if (universities)[i].UniversityName == (universities)[j].UniversityName && (universities)[i].CountryID == (universities)[j].CountryID {
				log.NewLogger().Error("University Name:" + (universities)[j].UniversityName + " exists")
				return errors.NewValidationError("University Name:" + (universities)[j].UniversityName + " exists")
			}
		}
	}

	// Add one university record at a time.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, university := range universities {
		err := service.AddUniversity(university, uow)
		if err != nil {
			uow.RollBack()
			return errors.NewValidationError("name-" + university.UniversityName + ": " + err.Error())
		}
	}

	// Commit only if all universities have been added.
	uow.Commit()
	return nil

}

// GetUniversityList returns all University from database.
func (service *UniversityService) GetUniversityList(tenantID uuid.UUID,
	universities *[]general.University) error {

	// Check if tenant exists..
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all universities and orderby university_name.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, universities, "`university_name`",
		repository.PreloadAssociations([]string{"Country"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetUniversities returns universities based on limit and offset.
func (service *UniversityService) GetUniversities(tenantID uuid.UUID,
	universities *[]general.University, parser *web.Parser, totalCount *int) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	limit, offset := parser.ParseLimitAndOffset()
	// Get all universities and orderby university_name.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, universities, "`university_name`",
		service.addSearchQueries(parser.Form),
		repository.Paginate(limit, offset, totalCount),
		repository.PreloadAssociations([]string{"Country"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetUniversitiesByCountryList returns all universities from database.
func (service *UniversityService) GetUniversitiesByCountryList(tenantID, countryID uuid.UUID,
	universities *[]general.University) error {

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all universities where country_id matches and orderby university_name.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, universities, "`university_name`",
		repository.Filter("`country_id`=?", countryID),
		repository.PreloadAssociations([]string{"Country"}))
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	uow.Commit()
	return nil
}

// GetUniversity returns particular University by ID.
func (service *UniversityService) GetUniversity(university *general.University) error {

	tenantID := university.TenantID

	// Check if tenant exists.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	err = service.doesUniversityExist(tenantID, university.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)
	err = service.Repository.GetRecordForTenant(uow, tenantID, university,
		repository.PreloadAssociations([]string{"Country"}))
	if err != nil {
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

// UpdateUniversity updates the data of University.
func (service *UniversityService) UpdateUniversity(university *general.University) error {

	tenantID := university.TenantID

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	err := service.extractID(university)
	if err != nil {
		return err
	}

	// Check all foreign key records.
	err = service.doForeignKeysExist(tenantID, university.UpdatedBy, university.CountryID)
	if err != nil {
		return err
	}

	// Check general type record.
	err = service.doesUniversityExist(tenantID, university.ID)
	if err != nil {
		return err
	}

	// Validate if fields that should have unique value are having unique values.
	err = service.validateFieldUniqueness(university)
	if err != nil {
		return err
	}

	// Update method of repo to update.
	uow := repository.NewUnitOfWork(service.DB, false)
	err = service.Repository.Update(uow, university)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteUniversity deletes the data of University.
func (service *UniversityService) DeleteUniversity(university *general.University) error {
	tenantID := university.TenantID

	// Check general type record.
	err := service.doesUniversityExist(tenantID, university.ID)
	if err != nil {
		return err
	}

	// Repository deleted_at update call.
	uow := repository.NewUnitOfWork(service.DB, false)

	// First update the deleted_by field of record and then soft delete.
	err = service.Repository.UpdateWithMap(uow, university, map[string]interface{}{
		"DeletedBy": university.DeletedBy,
		"DeletedAt": time.Now(),
	})
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

// validateFieldUniqueness validates id university name is already present or not.
func (service *UniversityService) validateFieldUniqueness(university *general.University) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, university.TenantID, general.University{},
		repository.Filter("`university_name`=? AND `id`!= ? AND `country_id`=?", university.UniversityName,
			university.ID, university.CountryID))
	if err := util.HandleIfExistsError("Record already exists with the name "+university.UniversityName, exists, err); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

// doForeignKeysExist will check the DB whether all foreign-keys are present in the table
// it will return error if no record is found in table.
func (service *UniversityService) doForeignKeysExist(tenantID, credentialID, countryID uuid.UUID) error {

	// Check if tenant exists.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Check if credential exists.
	if err := service.doesCredentialExist(tenantID, credentialID); err != nil {
		return err
	}

	// Check if country exists.
	if err := service.doesCountryExist(tenantID, countryID); err != nil {
		return err
	}
	return nil
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *UniversityService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExist returns error if there is no credential record in table for the given tenant.
func (service *UniversityService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id`=?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCountryExist returns error if there is no country record for the given tenant.
func (service *UniversityService) doesCountryExist(tenantID, countryID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Country{},
		repository.Filter("`id`=?", countryID))
	if err := util.HandleError("Invalid country ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesUniversityExist returns error if there is no university record for the given tenant.
func (service *UniversityService) doesUniversityExist(tenantID, universityID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.University{},
		repository.Filter("`id` = ?", universityID))
	if err := util.HandleError("Invalid university ID", exists, err); err != nil {
		return err
	}
	return nil
}

// assignCountryIDByName returns error if there is no university record for the given tenant.
func (service *UniversityService) assignCountryIDByName(university *general.University) error {
	tenantID := university.TenantID
	country := &general.Country{}
	countryName := university.Country.Name

	uow := repository.NewUnitOfWork(service.DB, true)

	fmt.Println("=========================name================================", countryName)
	// If only country name has been given, it will get record of that country and assign the ID
	err := service.Repository.GetRecordForTenant(uow, tenantID, country,
		repository.Filter("`name`=?", countryName))
	if err != nil {
		return errors.NewValidationError("Invalid country name")
	}
	fmt.Println(country.ID)
	university.Country.ID = country.ID
	return nil
}

// extractID extracts ids from forrign key objects.
func (service *UniversityService) extractID(university *general.University) error {
	if university.Country != nil {
		university.CountryID = university.Country.ID
	}
	return nil
}

// addSearchQueries adds search criteria to get all universities.
func (service *UniversityService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	fmt.Println("=========================In uni search============================", requestForm)
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}
	if universityName, ok := requestForm["universityName"]; ok {
		util.AddToSlice("`university_name`", "LIKE ?", "AND", "%"+universityName[0]+"%", &columnNames, &conditions, &operators, &values)
	}
	if countryID, ok := requestForm["countryID"]; ok {
		util.AddToSlice("`country_id`", "= ?", "AND", countryID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
