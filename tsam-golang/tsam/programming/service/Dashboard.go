package service

import (
	"net/url"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ConceptDashboardService provides methods to get for concept dashboard.
type ConceptDashboardService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewConceptDashboardService returns new instance of ConceptDashboardService.
func NewConceptDashboardService(db *gorm.DB, repository repository.Repository) *ConceptDashboardService {
	return &ConceptDashboardService{
		DB:         db,
		Repository: repository,
		association: []string{
		},
	}
}

// GetComplexConcepts all complex concepts.
func (service *ConceptDashboardService) GetComplexConcepts(concepts *[]programming.ComplexConcept, tenantID uuid.UUID,
	form url.Values) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create query precessors, add search quesries.
	var queryProcessors []repository.QueryProcessor

	// Add all necessary queries for sub query one.
	queryProcessors = append(queryProcessors,
		repository.Table("talent_concept_ratings"),
		repository.Select("AVG(talent_concept_ratings.`score`) AS score, talent_concept_ratings.`talent_id` AS talent_id,"+
			"programming_concepts.`name` AS concept_name, AVG(score) < 5 as less, programming_concepts.`complexity` AS complexity,"+
			"modules_programming_concepts.`level` AS level, programming_concepts.`description` AS description"),
		repository.Join("INNER JOIN modules_programming_concepts ON modules_programming_concepts.`id` = talent_concept_ratings.`module_programming_concept_id`"),
		repository.Join("INNER JOIN programming_concepts on programming_concepts.`id` = modules_programming_concepts.`programming_concept_id`"),
		repository.Join("INNER JOIN batch_modules on batch_modules.`module_id` = modules_programming_concepts.`module_id`"),
		repository.Filter("talent_concept_ratings.`deleted_at` IS NULL AND talent_concept_ratings.`tenant_id`=?", tenantID),
		repository.Filter("modules_programming_concepts.`deleted_at` IS NULL AND modules_programming_concepts.`tenant_id`=?", tenantID),
		repository.Filter("programming_concepts.`deleted_at` IS NULL AND programming_concepts.`tenant_id`=?", tenantID),
		repository.Filter("batch_modules.`deleted_at` IS NULL AND batch_modules.`tenant_id`=?", tenantID),
		service.addSearchQueries(form),
		repository.GroupBy("talent_concept_ratings.`module_programming_concept_id`, talent_concept_ratings.`talent_id`"))

	// Create query expression for sub query one.
	subQueryOne, err := service.Repository.SubQuery(uow, programming.ProgrammingQuestionSolution{}, queryProcessors...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	// Get all complex concepts.
	if err := service.Repository.GetAll(uow, &concepts,
		repository.RawQuery("SELECT * FROM ? as sub WHERE sub.less = 1", subQueryOne),
	); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return err
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search criteria.
func (service *ConceptDashboardService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if moduleID, ok := requestForm["moduleID"]; ok {
		util.AddToSlice("modules_programming_concepts.`module_id`", "=?", "AND", moduleID,
			&columnNames, &conditions, &operators, &values)
	}

	if talentID, ok := requestForm["talentID"]; ok {
		util.AddToSlice("talent_concept_ratings.`talent_id`", "=?", "AND", talentID,
			&columnNames, &conditions, &operators, &values)
	}

	if batchID, ok := requestForm["batchID"]; ok {
		util.AddToSlice("batch_modules.`batch_id`", "=?", "AND", batchID,
			&columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// doesTenantExist returns error if there is no tenant record in table.
func (service *ConceptDashboardService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id`=?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
