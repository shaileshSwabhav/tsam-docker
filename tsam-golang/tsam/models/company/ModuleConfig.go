package company

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewCompanyModuleConfig Create New Talent Module Config
func NewCompanyModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (module *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&Domain{},
		&Company{},
		&Branch{},
		&Enquiry{},
		&Requirement{},
		&CallRecord{},
		&CompanyRequirementQualifications{},
		&CompanyRequirementUniversities{},
		&CompanyEnquiryTechnologies{},
		&CompanyRequirementTalents{},
	}
	for _, model := range models {
		if err := module.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}
	if err := module.DB.Model(Branch{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Branch{}).
		AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Domain{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Branch{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Branch{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(Enquiry{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Enquiry{}).
		AddForeignKey("company_branch_id", "company_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Enquiry{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Enquiry{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Requirement{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Requirement{}).
		AddForeignKey("company_id", "companies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Requirement{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Requirement{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&Requirement{}).
		AddForeignKey("designation_id", "designations(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := module.DB.Model(&CallRecord{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//talent call record's foriegn key for talent
	if err := module.DB.Model(&CallRecord{}).
		AddForeignKey("enquiry_id", "company_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Company requirement's maps' foreign keys*******************************
	// Qualification map.
	if err := module.DB.Model(&CompanyRequirementQualifications{}).
		AddForeignKey("requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&CompanyRequirementQualifications{}).
		AddForeignKey("degree_id", "degrees(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// University map.
	if err := module.DB.Model(&CompanyRequirementUniversities{}).
		AddForeignKey("requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&CompanyRequirementUniversities{}).
		AddForeignKey("university_id", "universities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Technology map.
	if err := module.DB.Model(&CompanyRequirementTechnologies{}).
		AddForeignKey("requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&CompanyRequirementTechnologies{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent map.
	if err := module.DB.Model(&CompanyRequirementTalents{}).
		AddForeignKey("requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&CompanyRequirementTalents{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Company Module Configured.")
}
