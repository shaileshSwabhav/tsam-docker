package admin

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewAdministrationModuleConfig Return New Module Config.
func NewAdministrationModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// extract all to a general module config, service & controller as well
// see to it these tables are created before everyone.

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		// &Timesheet{},
		&Timesheet{},
		&TimesheetActivity{},
		&Employee{},
		&TargetCommunity{},
		&SwabhavEvent{},
	}

	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	//**********************************Timesheet's foreign keys*******************************
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("department_id", "departments(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("credential_id", "credentials(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("project_id", "projects(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("sub_project_id", "projects(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("batch_id", "batches(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&Timesheet{}).
	// 	AddForeignKey("session_id", "batch_sessions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//**********************************Supervisor foreign keys*******************************
	if err := config.DB.Model(&EmployeeSupervisor{}).
		AddForeignKey("employee_credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&EmployeeSupervisor{}).
		AddForeignKey("supervisor_credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Target community's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&TargetCommunity{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Department.
	if err := config.DB.Model(&TargetCommunity{}).
		AddForeignKey("department_id", "departments(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Credential.
	if err := config.DB.Model(&TargetCommunity{}).
		AddForeignKey("credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Function.
	if err := config.DB.Model(&TargetCommunity{}).
		AddForeignKey("function_id", "target_community_functions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Faculty.
	if err := config.DB.Model(&TargetCommunity{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Target community's maps' foreign keys*******************************
	// College branch map.
	if err := config.DB.Model(&TargetCommunityColleges{}).
		AddForeignKey("target_community_id", "target_communities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TargetCommunityColleges{}).
		AddForeignKey("branch_id", "college_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Company branch map.
	if err := config.DB.Model(&TargetCommunityCompanies{}).
		AddForeignKey("target_community_id", "target_communities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TargetCommunityCompanies{}).
		AddForeignKey("company_branch_id", "company_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Course map.
	if err := config.DB.Model(&TargetCommunityCourses{}).
		AddForeignKey("target_community_id", "target_communities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TargetCommunityCourses{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Timesheet's foreign keys*******************************
	if err := config.DB.Model(&Timesheet{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Timesheet{}).
		AddForeignKey("department_id", "departments(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Timesheet{}).
		AddForeignKey("credential_id", "credentials(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TimesheetActivity{}).
		AddForeignKey("project_id", "projects(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TimesheetActivity{}).
		AddForeignKey("sub_project_id", "projects(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TimesheetActivity{}).
		AddForeignKey("batch_id", "batches(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TimesheetActivity{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TimesheetActivity{}).
		AddForeignKey("timesheet_id", "timesheets(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Event's foreign keys*******************************
	if err := config.DB.Model(&SwabhavEvent{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Administration Module Configured.")
}
