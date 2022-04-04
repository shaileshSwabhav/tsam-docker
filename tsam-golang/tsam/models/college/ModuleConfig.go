package college

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

//ModuleConfig Automigrate Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

//NewCollegeModuleConfig Return New Module Config.
func NewCollegeModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&College{},
		&Branch{},
		&CampusDrive{},
		&Candidate{},
		&CampusTalentRegistration{},
		&Speaker{},
		&Seminar{},
		&Topic{},
		&SeminarTalentRegistration{},
		// &CampusDrive{},
		// &Talent{},
		// &DriveCandidates{},
		// &TalentTechnologies{},
		// &Seminar{},
		// &SeminarCandidates{},
	}
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}
	if err := config.DB.Model(Branch{}).
		AddForeignKey("university_id", "universities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Branch{}).
		AddForeignKey("college_id", "colleges(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(Branch{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Branch{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Branch{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&College{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Branch{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************CAMPUS DRIVE'S FOREIGN KEYS*******************************
	// Tenant.
	if err := config.DB.Model(&CampusDrive{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Campus Drive's maps' foreign keys*******************************
	// For salesperson map.
	if err := config.DB.Model(&CampusDriveSalesPerson{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CampusDriveSalesPerson{}).
		AddForeignKey("sales_person_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// For faculty map.
	if err := config.DB.Model(&CampusDriveFaculty{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CampusDriveFaculty{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// For developer map.
	if err := config.DB.Model(&CampusDriveDeveloper{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CampusDriveDeveloper{}).
		AddForeignKey("developer_id", "employees(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// For college branch map.
	if err := config.DB.Model(&CampusDriveCollegeBranch{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CampusDriveCollegeBranch{}).
		AddForeignKey("branch_id", "college_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// For company requirement map.
	if err := config.DB.Model(&CampusDriveCompanyRequirement{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CampusDriveCompanyRequirement{}).
		AddForeignKey("requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************CAMPUSTALENTREGISTRATION'S FOREIGN KEYS*******************************

	// Tenant.
	if err := config.DB.Model(&CampusTalentRegistration{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&CampusTalentRegistration{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Campus drive.
	if err := config.DB.Model(&CampusTalentRegistration{}).
		AddForeignKey("campus_drive_id", "campus_drives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************SPEAKER'S FOREIGN KEYS*******************************

	// Tenant.
	if err := config.DB.Model(&Speaker{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Designation.
	if err := config.DB.Model(&Speaker{}).
		AddForeignKey("designation_id", "designations(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************SEMINAR'S FOREIGN KEYS*******************************

	// Tenant.
	if err := config.DB.Model(&Seminar{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************SEMINAR'S MAPS' FOREIGN KEYS*******************************
	// For college branch map.
	if err := config.DB.Model(&SeminarCollegeBranch{}).
		AddForeignKey("seminar_id", "seminars(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SeminarCollegeBranch{}).
		AddForeignKey("branch_id", "college_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************SEMINAR TOPIC'S FOREIGN KEYS*******************************

	// Tenant.
	if err := config.DB.Model(&Topic{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Speaker.
	if err := config.DB.Model(&Topic{}).
		AddForeignKey("speaker_id", "speakers(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Speaker.
	if err := config.DB.Model(&Topic{}).
		AddForeignKey("seminar_id", "seminars(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************SEMINARTALENTREGISTRATION'S FOREIGN KEYS*******************************

	// Tenant.
	if err := config.DB.Model(&SeminarTalentRegistration{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&SeminarTalentRegistration{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Campus drive.
	if err := config.DB.Model(&SeminarTalentRegistration{}).
		AddForeignKey("seminar_id", "seminars(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("College Module Configured.")

}
