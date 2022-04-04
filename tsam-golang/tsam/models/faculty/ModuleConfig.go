package faculty

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewFacultyModuleConfig Create New Faculty Module Config
func NewFacultyModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&Faculty{},
		&Academic{},
		&Experience{},
		&FacultyAssessment{},
		// &Technology{},
	}
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	if err := config.DB.Model(&Faculty{}).
		AddForeignKey("country_id", "countries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Faculty{}).
		AddForeignKey("state_id", "states(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Faculty{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("degree_id", "degrees(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("specialization_id", "specializations(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("designation_id", "designations(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// *********************************M2M Faculty Technoloiges*********************************
	if err := config.DB.Table("faculties_technologies").
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Table("faculties_technologies").
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("created_by", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("question_id", "feedback_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("option_id", "feedback_options(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&FacultyAssessment{}).
		AddForeignKey("group_id", "feedback_question_groups(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Faculty Module Configured.")

}
