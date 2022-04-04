package talent

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig used for automifrating tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewTalentModuleConfig create new Talent ModuleConfig.
func NewTalentModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration update table structure with latest version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	// Table list.
	var models []interface{} = []interface{}{
		&Talent{},
		&Academic{},
		&Experience{},
		&CallRecord{},
		&LifetimeValue{},
		&Interview{},
		&InterviewSchedule{},
		&InterviewRound{},
		&NextAction{},
		&NextActionType{},
		&CareerPlan{},
		&InterviewTakenBy{},
		&WaitingList{},
		&TalentEventRegistration{},
		&AssignmentSubmission{},
		&AssignmentSubmissionUpload{},
		&TalentConceptRating{},
		&ProjectSubmission{},
		&ProjectSubmissionUpload{},
	}

	// Table migrantion.
	for _, tal := range models {
		if err := config.DB.AutoMigrate(tal).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}
	//*******************************************Talent's foreign keys********************************
	// Country.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("country_id", "countries(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// State.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("state_id", "states(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Source.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("source_id", "sources(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Salesperson.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("sales_person_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("referral_id", "talents(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&Talent{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
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

	//*************************************Talent call records' foreign key******************************
	// Tenant.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&CallRecord{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//*************************************Talent academics' foreign keys**************************
	// Talent.
	if err := config.DB.Model(&Academic{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
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

	//**********************************Talent experiences' foreign keys*******************************
	// Talent.
	if err := config.DB.Model(&Experience{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
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

	//**********************************Talent Lifetime Value's foreign keys*******************************
	// Talent.
	if err := config.DB.Model(&LifetimeValue{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&LifetimeValue{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Interview's foreign keys*******************************
	// Talent.
	if err := config.DB.Model(&Interview{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&Interview{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Interview schedule.
	if err := config.DB.Model(&Interview{}).
		AddForeignKey("schedule_id", "talent_interview_schedules(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Interview round.
	if err := config.DB.Model(&Interview{}).
		AddForeignKey("round_id", "talent_interview_rounds(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Interview schedule's foreign keys*******************************
	// Talent.
	if err := config.DB.Model(&InterviewSchedule{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Tenant.
	if err := config.DB.Model(&InterviewSchedule{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Interview round's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&InterviewRound{}).
		AddForeignKey("tenant_id", "tenants(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Next Action Type's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&NextActionType{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Next Action's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&NextAction{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&NextAction{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Next action type.
	if err := config.DB.Model(&NextAction{}).
		AddForeignKey("action_type_id", "talent_next_action_types(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent Next Action's maps' foreign keys*******************************
	// Course map.
	if err := config.DB.Model(&NextActionCourse{}).
		AddForeignKey("next_action_id", "talent_next_actions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&NextActionCourse{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Company map.
	if err := config.DB.Model(&NextActionCompanyBranch{}).
		AddForeignKey("next_action_id", "talent_next_actions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&NextActionCompanyBranch{}).
		AddForeignKey("company_branch_id", "company_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Technology map.
	if err := config.DB.Model(&NextActionTechnology{}).
		AddForeignKey("next_action_id", "talent_next_actions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&NextActionTechnology{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent's maps' foreign keys*******************************
	// Technology map.
	if err := config.DB.Model(&TalentTechnologies{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentTechnologies{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Experience technology map.
	if err := config.DB.Model(&ExperienceTechnology{}).
		AddForeignKey("experience_id", "talent_experiences(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ExperienceTechnology{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Talent's Career Plan's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&CareerPlan{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&CareerPlan{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Career objective.
	if err := config.DB.Model(&CareerPlan{}).
		AddForeignKey("career_objective_id", "career_objectives(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Career objectives courses.
	if err := config.DB.Model(&CareerPlan{}).
		AddForeignKey("career_objectives_courses_id", "career_objectives_courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Faculty.
	if err := config.DB.Model(&CareerPlan{}).
		AddForeignKey("faculty_id", "faculties(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Interview's maps' foreign keys*******************************
	// Taken by map.
	if err := config.DB.Model(&InterviewTakenBy{}).
		AddForeignKey("interview_id", "talent_interviews(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&InterviewTakenBy{}).
		AddForeignKey("credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Waiting list's foreign keys*******************************
	// Tenant.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Talent.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Company branch.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("company_branch_id", "company_branches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Company requirement.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("company_requirement_id", "company_requirements(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Course.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("course_id", "courses(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Batch.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Source.
	if err := config.DB.Model(&WaitingList{}).
		AddForeignKey("source_id", "sources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Talent event registration's foreign keys *******************************
	if err := config.DB.Model(&TalentEventRegistration{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentEventRegistration{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentEventRegistration{}).
		AddForeignKey("event_id", "events(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Talent assignment submissions foreign keys *******************************
	if err := config.DB.Model(&AssignmentSubmission{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&AssignmentSubmission{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&AssignmentSubmission{}).
		AddForeignKey("faculty_id", "faculties(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&AssignmentSubmission{}).
		AddForeignKey("batch_topic_assignment_id", "batch_topic_assignments(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// if err := config.DB.Model(&AssignmentSubmission{}).
	// 	AddForeignKey("previous_submission_id", "talent_assignment_submissions(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	//********************************** Talent submission uploads foreign keys *******************************
	if err := config.DB.Model(&AssignmentSubmissionUpload{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&AssignmentSubmissionUpload{}).
		AddForeignKey("talent_assignment_submission_id", "talent_assignment_submissions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentConceptRating{}).
		AddForeignKey("module_programming_concept_id", "modules_programming_concepts(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentConceptRating{}).
		AddForeignKey("talent_submission_id", "talent_assignment_submissions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentConceptRating{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&TalentConceptRating{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Talent Project submissions foreign keys *******************************
	if err := config.DB.Model(&ProjectSubmission{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProjectSubmission{}).
		AddForeignKey("talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProjectSubmission{}).
		AddForeignKey("faculty_id", "faculties(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProjectSubmission{}).
		AddForeignKey("batch_id", "batches(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProjectSubmission{}).
		AddForeignKey("batch_project_id", "batch_projects(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** Talent Project submission uploads foreign keys *******************************
	if err := config.DB.Model(&ProjectSubmissionUpload{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&ProjectSubmissionUpload{}).
		AddForeignKey("project_submission_id", "talent_project_submissions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Talent Module Configured.")

}
