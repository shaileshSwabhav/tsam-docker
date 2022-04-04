package batch

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewBatchModuleConfig Return New Module Config.
func NewBatchModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&AhaMoment{},
		&AhaMomentResponse{},
		&Batch{},
		&Timing{},
		&Eligibility{},
		&MappedTalent{},
		&FacultyTalentFeedback{},
		&FacultyTalentBatchSessionFeedback{},
		&Module{},
		&ModuleTiming{},
		&TalentBatchSessionFeedback{},
		&TalentFeedback{},
		&TopicAssignment{},
		&BatchSessionTalent{},
		&Session{},
		&SessionTopic{},
		&BatchSessionPrerequisite{},
		&Project{},
		&ProgrammingProjectRatingParameter{},
		&ProgrammingProjectRating{},
		// &BatchTopic{},
		// &MappedSession{},
	}

	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}
	// Add Foreign Key
	//*************************************Batch foreign key******************************
	if err := config.DB.Model(&Batch{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Batch{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&Batch{}).
		AddForeignKey("eligibility_id", "batch_eligibilities(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// if err := config.DB.Model(&Batch{}).
	// 	AddForeignKey("faculty_id", "faculties(id)", "RESTRICT", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	if err := config.DB.Model(&Batch{}).
		AddForeignKey("course_id", "courses(id)", "RESTRICT", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	//*************************************BatchTalents foreign key******************************
	if err := config.DB.Model(&MappedTalent{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&MappedTalent{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&MappedTalent{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	//*************************************Eligibility foreign key******************************
	if err := config.DB.Model(&Eligibility{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// if err := config.DB.Model(&Eligibility{}).
	// 	AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	//*************************************Session foreign key******************************
	// if err := config.DB.Model(&MappedSession{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&MappedSession{}).
	// 	AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	// if err := config.DB.Model(&MappedSession{}).
	// 	AddForeignKey("course_session_id", "course_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }
	//*************************************Batch Timing foreign key******************************
	if err := config.DB.Model(&Timing{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Timing{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// if err := config.DB.Model(&Timing{}).
	// 	AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&Timing{}).
	// 	AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//*************************************FacultyBatchFeedback foreign key******************************
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("question_id", "feedback_questions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentFeedback{}).
		AddForeignKey("option_id", "feedback_options(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************FacultyBatchSessionFeedback foreign key******************************
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("question_id", "feedback_questions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("option_id", "feedback_options(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&FacultyTalentBatchSessionFeedback{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************TalentBatchSessionFeedback foreign key******************************
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("question_id", "feedback_questions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentBatchSessionFeedback{}).
		AddForeignKey("option_id", "feedback_options(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************TalentBatchFeedback foreign key******************************
	if err := config.DB.Model(&TalentFeedback{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentFeedback{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentFeedback{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentFeedback{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&TalentFeedback{}).
		AddForeignKey("question_id", "feedback_questions(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************AhaMoment foreign key******************************
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("feeling_id", "feelings(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMoment{}).
		AddForeignKey("feeling_level_id", "feeling_levels(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************AhaMomentResponse foreign key******************************
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("feeling_id", "feelings(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("feeling_level_id", "feeling_levels(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("question_id", "feedback_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	if err := config.DB.Model(&AhaMomentResponse{}).
		AddForeignKey("aha_moment_id", "aha_moments(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Batch Session Programming Assignment foreign keys*******************************
	if err := config.DB.Model(&TopicAssignment{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicAssignment{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicAssignment{}).
		AddForeignKey("topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicAssignment{}).
		AddForeignKey("programming_question_id", "programming_questions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TopicAssignment{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Batch Sessions Talents foreign keys*******************************
	if err := config.DB.Model(&BatchSessionTalent{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&BatchSessionTalent{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&BatchSessionTalent{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&BatchSessionTalent{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH TOPIC *******************************
	// if err := config.DB.Model(&BatchTopic{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&BatchTopic{}).
	// 	AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&BatchTopic{}).
	// 	AddForeignKey("module_topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&BatchTopic{}).
	// 	AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//********************************** BATCH SESSION *******************************
	if err := config.DB.Model(&Session{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Session{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH SESSION TOPIC *******************************
	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("sub_topic_id", "module_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&SessionTopic{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH SESSION PREREQUISITE *******************************
	if err := config.DB.Model(&BatchSessionPrerequisite{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&BatchSessionPrerequisite{}).
		AddForeignKey("batch_session_id", "batch_sessions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH PROJECTS *******************************
	if err := config.DB.Model(&Project{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Project{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Project{}).
		AddForeignKey("programming_project_id", "programming_projects(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH MODULE *******************************
	if err := config.DB.Model(&Module{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Module{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Module{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Module{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BATCH MODULE TIMING *******************************
	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("module_id", "modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("batch_module_id", "batch_modules(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ModuleTiming{}).
		AddForeignKey("day_id", "days(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Programming Project Ratings *******************************
	if err := config.DB.Model(&ProgrammingProjectRatingParameter{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProgrammingProjectRating{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProgrammingProjectRating{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProgrammingProjectRating{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProgrammingProjectRating{}).
		AddForeignKey("talent_submission_id", "talent_project_submissions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProgrammingProjectRating{}).
		AddForeignKey("programming_project_rating_parameter_id", "programming_project_rating_parameter(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Batch Module Configured.")

}
