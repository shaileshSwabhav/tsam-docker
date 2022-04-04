package programming

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewProgrammingModuleConfig Return New programming Module Config.
func NewProgrammingModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (module *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&ProgrammingQuestion{},
		&ProgrammingQuestionType{},
		&ProgrammingQuestionOption{},
		&ProblemOfTheDay{},
		&ProgrammingQuestionTalentAnswer{},
		&ProgrammingConcept{},
		&ProgrammingQuestionSolution{},
		&ProgrammingQuestionSolutionIsViewed{},
		&ProgrammingQuestionTestCase{},
		&ProgrammingAssignment{},
		&ProgrammingAssignmentSubTask{},
		&ProgrammingProject{},
		&ModuleProgrammingConcepts{},
		// &ReadingProgrammingAssignment{},
	}

	for _, model := range models {
		if err := module.DB.Debug().AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	//********************************** PROGRAMMING QUESTION FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestion{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// // Programming question type.
	// if err := module.DB.Model(&ProgrammingQuestion{}).
	// 	AddForeignKey("programming_question_type_id", "programming_question_types(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//********************************** PROGRAMMING QUESTION TYPE FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionType{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING QUESTION OPTION FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionOption{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionOption{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROBLEM OF THE DAY FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProblemOfTheDay{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProblemOfTheDay{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING SOLUTION FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionTalentAnswer{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionTalentAnswer{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionTalentAnswer{}).
		AddForeignKey("programming_question_option_id", "programming_question_options(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionTalentAnswer{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING CONCEPT FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingConcept{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming concept.
	if err := module.DB.Model(&ProgrammingConcept{}).
		AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** CONCEPTMODULE FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ModuleProgrammingConcepts{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming concept.
	if err := module.DB.Model(&ModuleProgrammingConcepts{}).
		AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Module.
	if err := module.DB.Model(&ModuleProgrammingConcepts{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Parent programming concept module foreign keys *******************************

	if err := module.DB.Model(&ModuleProgrammingConceptsParentModuleProgrammingConcepts{}).
		AddForeignKey("parent_module_programming_concept_id", "modules_programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&ModuleProgrammingConceptsParentModuleProgrammingConcepts{}).
		AddForeignKey("module_programming_concept_id", "modules_programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************PROGRAMMING CONCEPT'S MAPS' FOREIGN KEYS*******************************

	// Programming question map.
	if err := module.DB.Model(&ProgrammingConceptsProgrammingQuestions{}).
		AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&ProgrammingConceptsProgrammingQuestions{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING QUESTION SOLUTION FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionSolution{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionSolution{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming language.
	if err := module.DB.Model(&ProgrammingQuestionSolution{}).
		AddForeignKey("programming_language_id", "programming_languages(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING SOLUTION IS VIEWED FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionSolutionIsViewed{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionSolutionIsViewed{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming language.
	if err := module.DB.Model(&ProgrammingQuestionSolutionIsViewed{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING QUESTION TEST CASE FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingQuestionTestCase{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Programming question.
	if err := module.DB.Model(&ProgrammingQuestionTestCase{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** PROGRAMMING ASSIGNMENT CASE FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingAssignment{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	// if err := module.DB.Model(&ReadingProgrammingAssignment{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//********************************** PROGRAMMING ASSIGNMENT SUB TASK CASE FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&ProgrammingAssignmentSubTask{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&ProgrammingAssignmentSubTask{}).
		AddForeignKey("programming_assignment_id", "programming_assignments(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&ProgrammingAssignmentSubTask{}).
		AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Batch Project foreign keys*******************************
	if err := module.DB.Model(&ProgrammingProject{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Table("programming_projects_technologies").
		AddForeignKey("programming_project_id", "programming_projects(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Table("programming_projects_technologies").
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Table("programming_projects_resources").
		AddForeignKey("programming_project_id", "programming_projects(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Table("programming_projects_resources").
		AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	log.NewLogger().Info("Programming Module Configured.")
}
