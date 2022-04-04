package service

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentAssignmentSubmissionService provides method to Update, Delete, Add, Get Method For talent_assignment_submission.
type TalentAssignmentSubmissionService struct {
	DB           *gorm.DB
	repository   repository.Repository
	associations []string
}

// TalentAssignmentSubmissionService returns a new instance of TalentAssignmentSubmissionService.
func NewTalentAssignmentSubmissionService(db *gorm.DB, repository repository.Repository) *TalentAssignmentSubmissionService {
	return &TalentAssignmentSubmissionService{
		DB:         db,
		repository: repository,
		associations: []string{
			"Talent", "AssignmentSubmissionUploads", "TalentConceptRatings",
			"TalentConceptRatings.ModuleProgrammingConcept",
			// "BatchSessionProgrammingAssignment", "BatchSessionProgrammingAssignment.ProgrammingAssignment",
			// "BatchSessionProgrammingAssignment.ProgrammingAssignment.ProgrammingQuestion",
			// "BatchSessionProgrammingAssignment.ProgrammingAssignment.ProgrammingQuestion.ProgrammingQuestionTypes",
			// "BatchSessionProgrammingAssignment.ProgrammingAssignment.ProgrammingAssignmentSubTask",
		},
	}
}

// AddTalentSubmission will add assignment submission of talent.
func (service *TalentAssignmentSubmissionService) AddTalentSubmission(submission *talent.AssignmentSubmission) error {

	// Check if foreign keys exist.
	err := service.doForeignKeysExist(submission, submission.CreatedBy)
	if err != nil {
		return err
	}

	// assign tenantID and createdBy to uploads slice if it exist.
	if submission.AssignmentSubmissionUploads != nil {
		for index := range submission.AssignmentSubmissionUploads {
			submission.AssignmentSubmissionUploads[index].TenantID = submission.TenantID
			submission.AssignmentSubmissionUploads[index].CreatedBy = submission.CreatedBy
		}
	}

	// assign tenantID and createdBy to uploads slice if it exist.
	if submission.TalentConceptRatings != nil {
		for index := range submission.TalentConceptRatings {
			submission.TalentConceptRatings[index].TenantID = submission.TenantID
			submission.TalentConceptRatings[index].CreatedBy = submission.CreatedBy
		}
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	submission.SubmittedOn = time.Now()

	err = service.repository.Add(uow, submission)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateTalentSubmission will update submitted assignment of talent.
func (service *TalentAssignmentSubmissionService) UpdateTalentSubmission(submission *talent.AssignmentSubmission) error {

	// Check if foreign keys exist.
	err := service.doForeignKeysExist(submission, submission.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if submission exist.
	err = service.doesSubmissionExist(submission.TenantID, submission.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempSubmission := talent.AssignmentSubmission{}

	err = service.repository.GetRecordForTenant(uow, submission.TenantID, &tempSubmission,
		repository.Filter("`id` = ?", submission.ID), repository.Select("`created_by`,`created_at`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	submission.CreatedBy = tempSubmission.CreatedBy
	submission.SubmittedOn = tempSubmission.CreatedAt

	err = service.repository.Save(uow, submission)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ScoreTalentSubmission will allow faculty to give scores to the latest submission.
func (service *TalentAssignmentSubmissionService) ScoreTalentSubmission(submission *talent.AssignmentSubmission) error {

	// Check if foreign keys exist.
	err := service.doForeignKeysExist(submission, submission.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if submission exist.
	err = service.doesSubmissionExist(submission.TenantID, submission.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)
	defer uow.RollBack()

	tempSubmission := &talent.AssignmentSubmission{}

	err = service.repository.GetRecordForTenant(uow, submission.TenantID, tempSubmission,
		repository.Filter("`id` = ?", submission.ID),
		repository.Select("`batch_topic_assignment_id`,`talent_id`,`submitted_on`,`tenant_id`"))
	if err != nil {
		return err
	}

	err = service.validateLatestSubmission(tempSubmission)
	if err != nil {
		return err
	}

	if submission.IsAccepted != nil && *submission.IsAccepted {
		err = service.addConceptRating(uow, submission)
		if err != nil {
			return err
		}
	}

	err = service.repository.UpdateWithMap(uow, submission, map[string]interface{}{
		"FacultyID":        submission.FacultyID,
		"Score":            submission.Score,
		"IsChecked":        true,
		"IsAccepted":       submission.IsAccepted,
		"FacultyRemarks":   submission.FacultyRemarks,
		"FacultyVoiceNote": submission.FacultyVoiceNote,
		"AcceptanceDate": func() *time.Time {
			if submission.IsAccepted != nil && *submission.IsAccepted {
				now := time.Now()
				submission.AcceptanceDate = &now
			}
			return submission.AcceptanceDate
		}(),
	})
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (service *TalentAssignmentSubmissionService) addConceptRating(uow *repository.UnitOfWork, submission *talent.AssignmentSubmission) error {
	for _, conceptRating := range submission.TalentConceptRatings {
		err := service.doesConceptModuleExist(submission.TenantID, conceptRating.ModuleProgrammingConceptID)
		if err != nil {
			return err
		}
		conceptRating.TalentID = submission.TalentID
		conceptRating.TalentSubmissionID = submission.ID
		conceptRating.TenantID = submission.TenantID

		err = conceptRating.Validate()
		if err != nil {
			return err
		}

		err = service.validateConceptRating(submission)
		if err != nil {
			return err
		}

		// Adding ratings
		err = service.repository.Add(uow, conceptRating)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetTalentSubmissions will return talent submissions.
func (service *TalentAssignmentSubmissionService) GetTalentSubmissions(tenantID, sessionAssignmentID, talentID uuid.UUID,
	talentSubmissions *[]talent.AssignmentSubmissionDTO, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesSessionAssignmentExist(tenantID, sessionAssignmentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.repository.GetAllForTenant(uow, tenantID, talentSubmissions,
		repository.Filter("`batch_topic_assignment_id` = ?", sessionAssignmentID),
		repository.Filter("`talent_id` = ?", talentID),
		repository.OrderBy("`submitted_on` DESC"),
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetTalentScores will return scores of all talents present in batch.
// works for now but needs optimization.
func (service *TalentAssignmentSubmissionService) GetTalentScores(tenantID, batchID uuid.UUID,
	talentAssignmentScores *[]talent.TalentAssignmentScoreDTO, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	var tempTalents []list.Talent

	err = service.repository.GetAll(uow, &tempTalents, repository.Filter("batch_talents.`batch_id` = ?", batchID),
		repository.Join("INNER JOIN `batch_talents` ON `talents`.`id` = `batch_talents`.`talent_id` AND "+
			"`talents`.`tenant_id` = `batch_talents`.`tenant_id`"), repository.Filter("`talents`.`tenant_id` = ?", tenantID),
		repository.Filter("`batch_talents`.`deleted_at` IS NULL AND `talents`.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`suspension_date` IS NULL AND batch_talents.`is_active` = ?", true))
	if err != nil {
		uow.RollBack()
		return err
	}

	for index := range tempTalents {
		var tempTalentScore talent.TalentAssignmentScoreDTO
		tempTalentScore.Talent = &tempTalents[index]

		err = service.repository.GetAllForTenant(uow, tenantID, &tempTalentScore.TalentSubmission,
			repository.Filter("`talent_id` = ?", tempTalents[index].ID),
			repository.PreloadWithCustomCondition(repository.Preload{
				Schema: "BatchTopicAssignment",
				Queryprocessors: []repository.QueryProcessor{
					repository.Join("LEFT JOIN `talent_assignment_submissions` ON " +
						"`talent_assignment_submissions`.`batch_topic_assignment_id` = batch_topic_assignments.`id` AND" +
						" `talent_assignment_submissions`.`tenant_id` = batch_topic_assignments.`tenant_id`"),
					service.addSearchQueries(parser.Form),
				}}),
			repository.PreloadAssociations([]string{"BatchTopicAssignment.ProgrammingQuestion"}))
		if err != nil {
			uow.RollBack()
			return err
		}

		*talentAssignmentScores = append(*talentAssignmentScores, tempTalentScore)
	}

	uow.Commit()
	return nil
}

// GetTalentScores will return scores of all talents present in batch.
// works for now but needs optimization.
func (service *TalentAssignmentSubmissionService) GetLatestAssignmentSubmissions(tenantID, batchID uuid.UUID,
	talentSubmissions *[]talent.AssignmentSubmissionDTO, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.repository.GetAllForTenant(uow, tenantID, talentSubmissions,
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.associations),
		repository.OrderBy("`created_at`"),
		repository.GroupBy("`talent_id`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	var tempTalents []list.Talent

	err = service.repository.GetAll(uow, &tempTalents, repository.Filter("batch_talents.`batch_id` = ?", batchID),
		repository.Join("INNER JOIN `batch_talents` ON `talents`.`id` = `batch_talents`.`talent_id` AND "+
			"`talents`.`tenant_id` = `batch_talents`.`tenant_id`"), repository.Filter("`talents`.`tenant_id` = ?", tenantID),
		repository.Filter("`batch_talents`.`deleted_at` IS NULL AND `talents`.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`suspension_date` IS NULL AND batch_talents.`is_active` = ?", true))
	if err != nil {
		uow.RollBack()
		return err
	}

	// for index := range tempTalents {
	// 	var tempTalentScore talent.TalentAssignmentScoreDTO
	// 	tempTalentScore.Talent = &tempTalents[index]

	// 	err = service.Repository.GetAllForTenant(uow, tenantID, &tempTalentScore.TalentSubmission,
	// 		repository.Filter("`talent_id` = ?", tempTalents[index].ID),
	// 		repository.PreloadWithCustomCondition(repository.Preload{
	// 			Schema: "BatchTopicAssignment",
	// 			Queryprocessors: []repository.QueryProcessor{
	// 				repository.Join("LEFT JOIN `talent_assignment_submissions` assignment ON " +
	// 					"assignment.`batch_topic_assignment_id` = batch_topic_assignments.`id` AND" +
	// 					" assignment.`tenant_id` = batch_topic_assignments.`tenant_id`"),
	// 				service.addSearchQueries(parser.Form),
	// 			}}),
	// 		repository.PreloadAssociations([]string{"BatchTopicAssignment.ProgrammingQuestion"}))
	// 	if err != nil {
	// 		uow.RollBack()
	// 		return err
	// 	}

	// 	*talentAssignmentScores = append(*talentAssignmentScores, tempTalentScore)
	// }

	uow.Commit()
	return nil
}

// GetAllAssignmentsSubmissionsOfTalent will return all submissions of a particular assignment
// of the given talent.
func (service *TalentAssignmentSubmissionService) GetAllAssignmentsSubmissionsOfTalent(tenantID, batchID, talentID uuid.UUID,
	topicAssignment *[]batch.TopicAssignmentDTO, parser *web.Parser) error {

	// check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// get all batch topic assignments.
	err = service.repository.GetAll(uow, topicAssignment,
		repository.Join("INNER JOIN `batch_modules` ON "+
		"`batch_modules`.`batch_id` = `batch_topic_assignments`.`batch_id` AND "+
		"`batch_modules`.`module_id` = `batch_topic_assignments`.`module_id` AND "+
		"`batch_modules`.`tenant_id` = `batch_topic_assignments`.`tenant_id`"),
		repository.Filter("`batch_topic_assignments`.`batch_id` = ?", batchID),
		repository.Filter("`batch_topic_assignments`.`tenant_id` = ? AND batch_topic_assignments.`deleted_at` IS NULL", tenantID),
		repository.Filter("`batch_topic_assignments`.`assigned_date` IS NOT NULL"),
		repository.OrderBy("`batch_topic_assignments`.`assigned_date`"),
		service.addSearchQueries(parser.Form),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ProgrammingQuestion"},
			repository.Preload{Schema: "Submissions",
				Queryprocessors: []repository.QueryProcessor{
					repository.Filter("`talent_assignment_submissions`.`talent_id` = ?", talentID),
					repository.OrderBy("`talent_assignment_submissions`.`submitted_on` DESC"),
					repository.PreloadAssociations([]string{"AssignmentSubmissionUploads", "TalentConceptRatings",
						"TalentConceptRatings.ModuleProgrammingConcept.ProgrammingConcept"}),
				},
			}),
	)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doForeignKeysExist will check if all foreign keys are valid.
func (service *TalentAssignmentSubmissionService) doForeignKeysExist(submission *talent.AssignmentSubmission,
	credentialID uuid.UUID) error {

	// check if tenant exist.
	err := service.doesTenantExist(submission.TenantID)
	if err != nil {
		return err
	}

	// check if credential exist.
	err = service.doesCredentialExist(submission.TenantID, credentialID)
	if err != nil {
		return err
	}

	// check if talent exist.
	err = service.doesTalentExist(submission.TenantID, submission.TalentID)
	if err != nil {
		return err
	}

	// check if faculty exist.
	if submission.FacultyID != nil {
		err = service.doesFacultyExist(submission.TenantID, *submission.FacultyID)
		if err != nil {
			return err
		}
	}

	// check if sessionAssignmentID exist.
	err = service.doesBatchSessionAssignmentExist(submission.TenantID, submission.BatchTopicAssignmentID)
	if err != nil {
		return err
	}

	return nil
}

// returns error if given submission is not the latest submission.
func (service *TalentAssignmentSubmissionService) validateLatestSubmission(submission *talent.AssignmentSubmission) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, submission.TenantID, submission,
		repository.Filter("`talent_id` = ?", submission.TalentID),
		repository.Filter("`batch_topic_assignment_id` = ?", submission.BatchTopicAssignmentID),
		repository.Filter("`submitted_on` > ?", submission.SubmittedOn))
	if err := util.HandleIfExistsError("Only latest submission can be updated", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if given submission's concepts are already rated.
func (service *TalentAssignmentSubmissionService) validateConceptRating(submission *talent.AssignmentSubmission) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, submission.TenantID, talent.TalentConceptRating{},
		repository.Filter("`talent_id` = ?", submission.TalentID),
		repository.Filter("`talent_submission_id` = ?", submission.ID))
	if err := util.HandleIfExistsError("Already rated for this submission", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no tenant record in table.
func (service *TalentAssignmentSubmissionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch_session record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesSessionAssignmentExist(tenantID, sessionAssignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.TopicAssignment{},
		repository.Filter("`id` = ?", sessionAssignmentID))
	if err := util.HandleError("Invalid session programming assignment ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no faculty record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.FacultyDetails{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no concept record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesConceptModuleExist(tenantID, conceptModuleID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ModuleProgrammingConcepts{},
		repository.Filter("`id` = ?", conceptModuleID))
	if err := util.HandleError("Invalid programming_concept_module ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no concept record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesProgrammingConceptExist(tenantID, conceptID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, programming.ProgrammingConcept{},
		repository.Filter("`id` = ?", conceptID))
	if err := util.HandleError("Invalid programming concept ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent submission record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesSubmissionExist(tenantID, submissionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.AssignmentSubmission{},
		repository.Filter("`id` = ?", submissionID))
	if err := util.HandleError("Invalid submission ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch_assignment record in table for the given tenant.
func (service *TalentAssignmentSubmissionService) doesBatchSessionAssignmentExist(tenantID, sessionAssignmentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.TopicAssignment{},
		repository.Filter("`id` = ?", sessionAssignmentID))
	if err := util.HandleError("Invalid batch session assignment ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *TalentAssignmentSubmissionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if dueDate, ok := requestForm["dueDate"]; ok {
		util.AddToSlice("`batch_topic_assignments`.`due_date`", "<= ?", "AND", dueDate, &columnNames, &conditions, &operators, &values)
	}

	facultyID := requestForm.Get("facultyID")
	if !util.IsEmpty(facultyID) {
		util.AddToSlice("`batch_modules`.`faculty_id`", "= ?", "AND",
			facultyID, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}
