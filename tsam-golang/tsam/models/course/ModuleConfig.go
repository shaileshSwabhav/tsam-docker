package course

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewCourseModuleConfig Return New Module Config.
func NewCourseModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&CourseSession{},
		&Course{},
		&CourseTechnicalAssessment{},
		&TopicProgrammingQuestion{},
		&TopicProgrammingConcept{},
		&CourseModule{},
		&Module{},
		&ModuleTopic{},
		&ModuleResource{},
	}

	// Migrate Table
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	// Add Foreign Key

	if err := config.DB.Model(&Course{}).
		AddForeignKey("eligibility_id", "eligibilities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Course{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseSession{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseSession{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseSession{}).
		AddForeignKey("session_id", "course_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseSession{}).
		AddForeignKey("course_module_id", "course_modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleResource{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&ModuleResource{}).
		AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Course Technical Assessment foreign keys *******************************

	if err := config.DB.Model(&CourseTechnicalAssessment{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseTechnicalAssessment{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseTechnicalAssessment{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Course Sessions Resource foreign keys *******************************

	// if err := config.DB.Model(&Resource{}).
	// 	AddForeignKey("course_session_id", "course_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&Resource{}).
	// 	AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//********************************** COURSE PROGRAMMING ASSIGNMENT FOREIGN KEYS *******************************

	if err := config.DB.Model(&TopicProgrammingQuestion{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingQuestion{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingQuestion{}).
		AddForeignKey("topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingQuestion{}).
		AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** COURSE PROGRAMMING CONCEPT FOREIGN KEYS *******************************

	if err := config.DB.Model(&TopicProgrammingConcept{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingConcept{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingConcept{}).
		AddForeignKey("sub_topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicProgrammingConcept{}).
		AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** COURSE MODULE FOREIGN KEYS *******************************

	if err := config.DB.Model(&CourseModule{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseModule{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CourseModule{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// ********************************** MODULE **********************************

	if err := config.DB.Model(&Module{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// ********************************** TOPIC **********************************

	if err := config.DB.Model(&ModuleTopic{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTopic{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// if err := config.DB.Model(&ModuleTopic{}).
	// 	AddForeignKey("programming_concept_id", "programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	if err := config.DB.Model(&ModuleTopic{}).
		AddForeignKey("topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Course Module Configured.")

}
