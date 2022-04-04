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
	general "github.com/techlabs/swabhav/tsam/models/general"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentLifetimeValueService provides method to update, delete, add, get all, get one for talent lifetime values.
type TalentLifetimeValueService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewTalentLifetimeValueService returns new instance of TalentLifetimeValueService.
func NewTalentLifetimeValueService(db *gorm.DB, repository repository.Repository) *TalentLifetimeValueService {
	return &TalentLifetimeValueService{
		DB:         db,
		Repository: repository,
	}
}

// AddTalentLifetimeValue adds one talent lifetime value to database.
func (service *TalentLifetimeValueService) AddTalentLifetimeValue(lifetimeValue *tal.LifetimeValue) error {
	// Get credential id from CreatedBy field of lifetimeValue(set in controller).
	credentialID := lifetimeValue.CreatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(lifetimeValue.TenantID, lifetimeValue.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add talent lifetime value to database.
	if err := service.Repository.Add(uow, lifetimeValue); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent lifetime value could not be added", http.StatusInternalServerError)
	}

	// Update lifetime value field of talents.
	if err := service.UpdateTalentLifetimeValueField(uow, lifetimeValue); err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTalentLifetimeValue gets one talent lifetime value form database.
func (service *TalentLifetimeValueService) GetTalentLifetimeValue(lifetimeValue *tal.LifetimeValue) error {
	// Validate tenant id.
	if err := service.doesTenantExist(lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(lifetimeValue.TenantID, lifetimeValue.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent lifetime value.
	if err := service.Repository.GetRecordForTenant(uow, lifetimeValue.TenantID, lifetimeValue,
		repository.Filter("`talent_id`=?", lifetimeValue.TalentID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetTalentLifetimeValueReports gets all talent lifetime value form database.
func (service *TalentLifetimeValueService) GetTalentLifetimeValueReports(lifetimeValues *[]tal.LifetimeValueReport, lognID uuid.UUID,
	tenantID uuid.UUID, limit int, offset int, totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, searchForm url.Values) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Query processors common for all queries.
	commonQueryProcessors := service.addSearchQueries(searchForm, tenantID)

	// Check if login id is salesperson or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.User{},
		repository.Join("left join roles ON users.`role_id` = roles.`id`"),
		repository.Filter("users.`id`=? AND roles.`role_name`=? AND users.`tenant_id`=? AND roles.`tenant_id`=?",
			lognID, "salesperson", tenantID, tenantID))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if exists {
		commonQueryProcessors = append(commonQueryProcessors, repository.Filter("`sales_person_id`=?", lognID))
	}

	// Add all filters and joins.
	commonQueryProcessors = append(commonQueryProcessors,
		repository.Join("JOIN talents ON talents.`id` = talent_lifetime_values.`talent_id`"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talent_lifetime_values.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
	)

	// Query processors for search and get.
	var queryProcessorsForSearchAndGet []repository.QueryProcessor
	queryProcessorsForSearchAndGet = append(queryProcessorsForSearchAndGet, commonQueryProcessors...)
	queryProcessorsForSearchAndGet = append(queryProcessorsForSearchAndGet, repository.Paginate(limit, offset, totalCount),
		repository.Select("*"))

	// Get talent and lifetimevalue from database.
	if err := service.Repository.GetAllInOrder(uow, lifetimeValues, "`first_name`, `last_name`", queryProcessorsForSearchAndGet...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for sum of lifetime value.
	var queryProcessorsForSum []repository.QueryProcessor
	queryProcessorsForSum = append(queryProcessorsForSum, repository.Select("sum(talents.`lifetime_value`) as total_lifetime_value"),
		repository.Table("talent_lifetime_values"))
	queryProcessorsForSum = append(queryProcessorsForSum, commonQueryProcessors...)

	// Get total lifetime value from database.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForSum...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// UpdateTalentLifetimeValue updates talent lifetime value in Database.
func (service *TalentLifetimeValueService) UpdateTalentLifetimeValue(lifetimeValue *tal.LifetimeValue) error {
	// Get credential id from UpdatedBy field of lifetimeValue(set in controller).
	credentialID := lifetimeValue.UpdatedBy

	// Validate tenant id.
	if err := service.doesTenantExist(lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(lifetimeValue.TenantID, lifetimeValue.TalentID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting talent lifetime value already present in database.
	tempTalentLifetimeValue := tal.LifetimeValue{}

	// Get talent lifetime value for getting created_by field of talent lifetime value from database.
	if err := service.Repository.GetForTenant(uow, lifetimeValue.TenantID, lifetimeValue.ID, &tempTalentLifetimeValue); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp talent lifetime value to talent lifetime value to be updated.
	lifetimeValue.CreatedBy = tempTalentLifetimeValue.CreatedBy

	// Update Talent lifetime value.
	if err := service.Repository.Save(uow, lifetimeValue); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent  lifetime value could not be updated", http.StatusInternalServerError)
	}

	// Update lifetime value field of talents.
	if err := service.UpdateTalentLifetimeValueField(uow, lifetimeValue); err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTalentLifetimeValue deletes one talent lifetime value form database.
func (service *TalentLifetimeValueService) DeleteTalentLifetimeValue(lifetimeValue *tal.LifetimeValue) error {
	// Starting new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Get credential id from DeletedBy field of lifetimeValue(set in controller).
	credentialID := lifetimeValue.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate talent lifetime value id.
	if err := service.doesTalentLifetimeValueExist(lifetimeValue.ID, lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, lifetimeValue.TenantID); err != nil {
		return err
	}

	// Validate talent id.
	if err := service.doesTalentExist(lifetimeValue.TenantID, lifetimeValue.TalentID); err != nil {
		return err
	}

	// Update talent's lifetime_value field to null.
	if err := service.Repository.UpdateWithMap(uow, &tal.Talent{}, map[string]interface{}{"LifetimeValue": nil},
		repository.Filter("`id` = ?", lifetimeValue.TalentID)); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Lifetime Value could not be updated", http.StatusInternalServerError)
	}

	// Update talent lifetime value for updating deleted_by and deleted_at field of talent lifetime value.
	if err := service.Repository.UpdateWithMap(uow, &tal.LifetimeValue{}, map[string]interface{}{
		"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`tenant_id`=? AND `id`=?", lifetimeValue.TenantID, lifetimeValue.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent lifetime value could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TalentLifetimeValueService) doesTenantExist(tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if tenant(parent tenant) exists or not.
	exists, err := repository.DoesRecordExist(uow.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *TalentLifetimeValueService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check if credential(parent credential) exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *TalentLifetimeValueService) doesTalentExist(tenantID uuid.UUID, talentID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent talent exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTalentLifetimeValueExist validates if talent lifetime value exists or not in database.
func (service *TalentLifetimeValueService) doesTalentLifetimeValueExist(lifetimeValueID uuid.UUID, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check talent lifetime value exists or not.
	exists, err := repository.DoesRecordExistForTenant(uow.DB, tenantID, tal.LifetimeValue{},
		repository.Filter("`id` = ?", lifetimeValueID))
	if err := util.HandleError("Invalid talent lifetime value ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// UpdateTalentLifetimeValueField updates the talent's lifetime value column in database.
func (service *TalentLifetimeValueService) UpdateTalentLifetimeValueField(uow *repository.UnitOfWork,
	lifetimeValue *tal.LifetimeValue) error {

	// Create bucket for sum of lifetime values.
	var lifetimeValueSum uint = 0

	// If upsell exists then add it to sum.
	if lifetimeValue.Upsell != nil {
		lifetimeValueSum = lifetimeValueSum + *lifetimeValue.Upsell
	}

	// If placement exists then add it to sum.
	if lifetimeValue.Placement != nil {
		lifetimeValueSum = lifetimeValueSum + *lifetimeValue.Placement
	}

	// If knowledge exists then add it to sum.
	if lifetimeValue.Knowledge != nil {
		lifetimeValueSum = lifetimeValueSum + *lifetimeValue.Knowledge
	}

	// If teaching exists then add it to sum.
	if lifetimeValue.Teaching != nil {
		lifetimeValueSum = lifetimeValueSum + *lifetimeValue.Teaching
	}

	// Update talent's lifetime_value field.
	if err := service.Repository.UpdateWithMap(uow, &tal.Talent{}, map[string]interface{}{"LifetimeValue": lifetimeValueSum},
		repository.Filter("`id` = ?", lifetimeValue.TalentID)); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Lifetime Value could not be updated", http.StatusInternalServerError)
	}

	return nil
}

// addSearchQueries adds all search queries if any when getAll is called.
func (service *TalentLifetimeValueService) addSearchQueries(searchForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
	fmt.Println("=========================In lifetime value search============================", searchForm)

	// Check if there is search criteria given.
	if len(searchForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	// First name search.
	if _, ok := searchForm["firstName"]; ok {
		fmt.Println("first name *********************************************************")
		util.AddToSlice("`first_name`", "LIKE ?", "AND", "%"+searchForm.Get("firstName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Last name search.
	if _, ok := searchForm["lastName"]; ok {
		util.AddToSlice("`last_name`", "LIKE ?", "AND", "%"+searchForm.Get("lastName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Email search.
	if _, ok := searchForm["email"]; ok {
		util.AddToSlice("`email`", "LIKE ?", "AND", "%"+searchForm.Get("email")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Upsell value minimum search.
	if upsellMinimum, ok := searchForm["upsellMinimum"]; ok {
		util.AddToSlice("`upsell`", ">=?", "AND", upsellMinimum, &columnNames, &conditions, &operators, &values)
	}

	// Upsell value maximum search.
	if upsellMaximum, ok := searchForm["upsellMaximum"]; ok {
		util.AddToSlice("`upsell`", "<=?", "AND", upsellMaximum, &columnNames, &conditions, &operators, &values)
	}

	// Placement value minimum search.
	if placementMinimum, ok := searchForm["placementMinimum"]; ok {
		util.AddToSlice("`placement`", ">=?", "AND", placementMinimum, &columnNames, &conditions, &operators, &values)
	}

	// Placement value maximum search.
	if placementMaximum, ok := searchForm["placementMaximum"]; ok {
		util.AddToSlice("`placement`", "<=?", "AND", placementMaximum, &columnNames, &conditions, &operators, &values)
	}

	// Knowledge value minimum search.
	if knowledgeMinimum, ok := searchForm["knowledgeMinimum"]; ok {
		util.AddToSlice("`knowledge`", ">=?", "AND", knowledgeMinimum, &columnNames, &conditions, &operators, &values)
	}

	// Knowledge value maximum search.
	if knowledgeMaximum, ok := searchForm["knowledgeMaximum"]; ok {
		util.AddToSlice("`knowledge`", "<=?", "AND", knowledgeMaximum, &columnNames, &conditions, &operators, &values)
	}

	// Teaching value minimum search.
	if teachingMinimum, ok := searchForm["teachingMinimum"]; ok {
		util.AddToSlice("`teaching`", ">=?", "AND", teachingMinimum, &columnNames, &conditions, &operators, &values)
	}

	// Teaching value maximum search.
	if teachingMaximum, ok := searchForm["teachingMaximum"]; ok {
		util.AddToSlice("`teaching`", "<=?", "AND", teachingMaximum, &columnNames, &conditions, &operators, &values)
	}

	// Total lifetime value minimum search.
	if totalMinimum, ok := searchForm["totalMinimum"]; ok {
		util.AddToSlice("`lifetime_value`", ">=?", "AND", totalMinimum, &columnNames, &conditions, &operators, &values)
	}

	// Total lifetime value maximum search.
	if totalMaximum, ok := searchForm["totalMaximum"]; ok {
		util.AddToSlice("`lifetime_value`", "<=?", "AND", totalMaximum, &columnNames, &conditions, &operators, &values)
	}

	// Group by campus drive id and add all filters.
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcesors
}
