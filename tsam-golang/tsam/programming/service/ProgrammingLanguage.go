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
)

// ProgrammingLanguageService provide method to update, delete, add, get method for ProgrammingLanguage.
type ProgrammingLanguageService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewProgrammingLanguageService creates new instance of ProgrammingLanguageService.
func NewProgrammingLanguageService(db *gorm.DB, repository repository.Repository) *ProgrammingLanguageService {
	return &ProgrammingLanguageService{
		DB:         db,
		Repository: repository,
	}
}

// AddProgrammingLanguage adds new programming language to database.
func (service *ProgrammingLanguageService) AddProgrammingLanguage(language *general.ProgrammingLanguage) error {

	// Validate tenant id.
	err := service.doesTenantExist(language.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same language name exists.
	err = service.doesNameExist(language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(language.CreatedBy, language.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add language.
	err = service.Repository.Add(uow, language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateProgrammingLanguage updates language in database.
func (service *ProgrammingLanguageService) UpdateProgrammingLanguage(language *general.ProgrammingLanguage) error {

	// Validate tenant ID.
	err := service.doesTenantExist(language.TenantID)
	if err != nil {
		return err
	}

	// Validate language ID.
	err = service.doesProgrammingLanguageExist(language.ID, language.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(language.UpdatedBy, language.TenantID)
	if err != nil {
		return err
	}

	// Validate if same language name exists.
	err = service.doesNameExist(language)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update language.
	err = service.Repository.Update(uow, language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteProgrammingLanguage delete language from database.
func (service *ProgrammingLanguageService) DeleteProgrammingLanguage(language *general.ProgrammingLanguage) error {
	credentialID := language.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(language.TenantID)
	if err != nil {
		return err
	}

	// Validate language ID.
	err = service.doesProgrammingLanguageExist(language.ID, language.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, language.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update language for updating deleted_by and deleted_at fields of language
	if err := service.Repository.UpdateWithMap(uow, language, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", language.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Programming Language could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetProgrammingLanguage returns one language by id.
func (service *ProgrammingLanguageService) GetProgrammingLanguage(language *general.ProgrammingLanguage) error {

	// Validate tenant ID.
	err := service.doesTenantExist(language.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get language by id from database.
	err = service.Repository.GetForTenant(uow, language.TenantID, language.ID, language)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingLanguages returns all programming languages.
func (service *ProgrammingLanguageService) GetProgrammingLanguages(languages *[]general.ProgrammingLanguage,
	form url.Values, tenantID uuid.UUID, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all programming language.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, languages, "`name`",
		service.addSearchQueries(form),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetProgrammingLanguageList returns all the programming language(without pagination).
func (service *ProgrammingLanguageService) GetProgrammingLanguageList(languages *[]general.ProgrammingLanguage,
	tenantID uuid.UUID) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all programming language.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, languages, "`name`")
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds all search queries by comparing with the language data.
func (service *ProgrammingLanguageService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// Name.
	if _, ok := requestForm["name"]; ok {
		util.AddToSlice("`name`", "LIKE ?", "AND", "%"+requestForm.Get("name")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Rating.
	if rating, ok := requestForm["rating"]; ok {
		util.AddToSlice("`rating`", "= ?", "AND", rating, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesNameExist returns true if the name already exists for the language in database.
func (service *ProgrammingLanguageService) doesNameExist(language *general.ProgrammingLanguage) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, language.TenantID, &general.ProgrammingLanguage{},
		repository.Filter("`name`=? AND `id`!=?", language.Name, language.ID))
	if err := util.HandleIfExistsError("names:"+language.Name+" exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *ProgrammingLanguageService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesProgrammingLanguageExist validates if language exists or not in database.
func (service *ProgrammingLanguageService) doesProgrammingLanguageExist(languageID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.ProgrammingLanguage{},
		repository.Filter("`id` = ?", languageID))
	if err := util.HandleError("Invalid programming language ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *ProgrammingLanguageService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
