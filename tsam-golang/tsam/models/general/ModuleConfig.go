package general

import (
	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewGeneralModuleConfig Return New Module Config.
func NewGeneralModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration() {
	var models []interface{} = []interface{}{
		&City{},
		&CommonType{},
		&Country{},
		&Credential{},
		&Day{},
		&Degree{},
		&Department{},
		&Designation{},
		&Eligibility{},
		&Employee{},
		&Examination{},
		&FeedbackOption{},
		&FeedbackQuestion{},
		&LoginSession{},
		&MastersAbroad{},
		&Menu{},
		&Outcome{},
		&Purpose{},
		&Role{},
		&Source{},
		&Specialization{},
		&State{},
		&SwabhavProject{},
		&Technology{},
		&Tenant{},
		&User{},
		&University{},
		&Score{},
		&CareerObjective{},
		&CareerObjectivesCourse{},
		&Feeling{},
		&FeelingLevel{},
		&FeedbackQuestionGroup{},
		&TargetCommunityFunction{},
		&SalaryTrend{},
		&ProgrammingLanguage{},

		//android
		&AndroidUser{},

		//notification
		&Notification_Test{},
	}

	for _, model := range models {
		if err := config.DB.Debug().AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	//foreign keys for notification
	if err := config.DB.Model(&Notification_Test{}).
		AddForeignKey("blog_notification_id", "blog_notifications(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// Tenant.
	if err := config.DB.Model(&Notification_Test{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// notifierID.
	if err := config.DB.Model(&Notification_Test{}).
		AddForeignKey("notifier_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// notifiediD.
	if err := config.DB.Model(&Notification_Test{}).
		AddForeignKey("notified_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// typeID.
	if err := config.DB.Model(&Notification_Test{}).
		AddForeignKey("type_id", "notification_types(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// foreign key for state.
	if err := config.DB.Model(&State{}).
		AddForeignKey("country_id", "countries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&State{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// foreign key for city.
	if err := config.DB.Model(&City{}).
		AddForeignKey("state_id", "states(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&City{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// foreign key for role.
	if err := config.DB.Debug().Model(&Role{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// foreign key for menu.
	if err := config.DB.Debug().Model(&Menu{}).
		AddForeignKey("menu_id", "menus(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Menu{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Menu{}).
		AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// foreign keys for login.

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("college_id", "colleges(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("sales_person_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("company_id", "company_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Debug().Model(&Credential{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Foreign keys for login sessions.
	// if err := config.DB.Debug().Model(&LoginSession{}).
	// 	AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	if err := config.DB.Debug().Model(&LoginSession{}).
		AddForeignKey("login_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Debug().Model(&LoginSession{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Debug().Model(&Country{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Debug().Model(&Eligibility{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Debug().Model(&Specialization{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Debug().Model(&Specialization{}).
		AddForeignKey("degree_id", "degrees(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Degree{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Designation{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&CommonType{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Outcome{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Outcome{}).
		AddForeignKey("purpose_id", "purposes(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Purpose{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Source{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Technology{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&User{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&User{}).
		AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&University{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&University{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Examination{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Masters abroad foreign keys*******************************
	// Adding tenant_id foreign key in masters_abroad.
	if err := config.DB.Model(&MastersAbroad{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Adding degree_id foreign key in masters_abroad.
	if err := config.DB.Model(&MastersAbroad{}).
		AddForeignKey("degree_id", "degrees(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Adding talent_id foreign key in masters_abroad.
	if err := config.DB.Model(&MastersAbroad{}).
		AddForeignKey("talent_id", "talents(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Adding enquiry_id foreign key in masters_abroad.
	if err := config.DB.Model(&MastersAbroad{}).
		AddForeignKey("enquiry_id", "talent_enquiries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Masters Abroad's maps' foreign keys*******************************
	// For country map.
	if err := config.DB.Model(&MastersAbroadCountries{}).
		AddForeignKey("masters_abroad_id", "masters_abroad(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&MastersAbroadCountries{}).
		AddForeignKey("country_id", "countries(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// For university map.
	if err := config.DB.Model(&MastersAbroadUniversities{}).
		AddForeignKey("masters_abroad_id", "masters_abroad(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&MastersAbroadUniversities{}).
		AddForeignKey("university_id", "universities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Score's foreign keys*******************************
	// Adding tenant_id foreign key in score.
	if err := config.DB.Model(&Score{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Adding examination_id foreign key in score.
	if err := config.DB.Model(&Score{}).
		AddForeignKey("examination_id", "examinations(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Adding masters_abroad_id foreign key in score.
	if err := config.DB.Model(&Score{}).
		AddForeignKey("masters_abroad_id", "masters_abroad(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Day's foreign keys*******************************
	if err := config.DB.Model(&Day{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************FeedbackQuestion's foreign keys*******************************
	if err := config.DB.Model(&FeedbackQuestion{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FeedbackQuestion{}).
		AddForeignKey("group_id", "feedback_question_groups(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************FeedbackOption's foreign keys*******************************
	if err := config.DB.Model(&FeedbackOption{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FeedbackOption{}).
		AddForeignKey("question_id", "feedback_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Department's foreign keys*******************************
	if err := config.DB.Model(&Department{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Project's foreign keys*******************************
	if err := config.DB.Model(&SwabhavProject{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// if err := config.DB.Model(&SwabhavProject{}).
	// 	AddForeignKey("project_id", "swabhav_projects(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//**********************************Employee's foreign keys*******************************
	if err := config.DB.Model(&Employee{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Mapped table
	// if err := config.DB.Model(&EmployeeTechnologies{}).
	// 	AddForeignKey("employee_id", "employees(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	if err := config.DB.Table("employees_technologies").
		AddForeignKey("employee_id", "employees(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Table("employees_technologies").
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Career Objective's foreign keys*******************************
	if err := config.DB.Model(&CareerObjective{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Career Objective Course's foreign keys*******************************
	if err := config.DB.Model(&CareerObjectivesCourse{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CareerObjectivesCourse{}).
		AddForeignKey("career_objective_id", "career_objectives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&CareerObjectivesCourse{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	//**********************************Feelings foreign keys*******************************
	if err := config.DB.Model(&Feeling{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	//**********************************Feelings foreign keys*******************************
	if err := config.DB.Model(&FeelingLevel{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FeelingLevel{}).
		AddForeignKey("feeling_id", "feelings(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Feedback Question Group's foreign keys*******************************
	if err := config.DB.Model(&FeedbackQuestionGroup{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Target Community Function's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&TargetCommunityFunction{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Department.
	if err := config.DB.Model(&TargetCommunityFunction{}).
		AddForeignKey("department_id", "departments(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Resources foreign keys*******************************

	// if err := config.DB.Model(&Resource{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//**********************************Salary-trend foreign keys*******************************

	if err := config.DB.Model(&SalaryTrend{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SalaryTrend{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SalaryTrend{}).
		AddForeignKey("designation_id", "designations(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************PRGRAMMING LANGUAGE*******************************

	// Tenant.
	if err := config.DB.Model(&ProgrammingLanguage{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("General Module Configured.")
}
