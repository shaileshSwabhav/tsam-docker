package talentenquiry

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig used for automigration of tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewTalentEnquiryModuleConfig creates new enquiry ModuleConfig
func NewTalentEnquiryModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration update table structure with latest version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	// Table list.
	var models []interface{} = []interface{}{
		&Enquiry{},
		&Academic{},
		&Experience{},
		&CallRecord{},
	}

	// Table migrantion.
	for _, enq := range models {
		if err := config.DB.AutoMigrate(enq).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}
	//*******************************************Enquiry's foreign keys********************************
	// Country.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// State.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Source.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("source_id", "sources(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Salesperson.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("talent_id", "talents(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&Enquiry{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************Enquiry call records' foreign key******************************
	// Tenant.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Enquiry.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Purpose.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("purpose_id", "purposes(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// Outcome.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("outcome_id", "outcomes(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************Enquiry academics' foreign keys**************************
	// Enquiry.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Degree.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("degree_id", "degrees(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Specialization.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("specialization_id", "specializations(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// College branch.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("college_branch_id", "college_branches(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Enquiry experiences' foreign keys*******************************
	// Enquiry.
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Designation.
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("designation_id", "designations(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Enquiry's maps' foreign keys*******************************
	// Technology map.
	if err := config.DB.Model(&EnquiryTechnologies{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&EnquiryTechnologies{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Course map.
	if err := config.DB.Model(&EnquiryCourses{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&EnquiryCourses{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Experience technology map.
	if err := config.DB.Model(&ExperienceTechnology{}).
		AddForeignKey("experience_id", "talent_enquiry_experiences(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ExperienceTechnology{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Talent Enquiry Module Configured.")
}
