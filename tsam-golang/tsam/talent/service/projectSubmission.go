package service

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// TalentProjectSubmissionService provides method to Update, Delete, Add, Get Method For talent_project_submission.
type TalentProjectSubmissionService struct {
	DB           *gorm.DB
	repository   repository.Repository
	associations []string
}

// TalentProjectSubmissionService returns a new instance of TalentProjectSubmissionService.
func NewTalentProjectSubmissionService(db *gorm.DB, repository repository.Repository) *TalentProjectSubmissionService {
	return &TalentProjectSubmissionService{
		DB:         db,
		repository: repository,
		associations: []string{
			"Talent", "ProjectSubmissionUpload",
		},
	}
}

// AddTalentProjectSubmission will add project submission of talent.
func (service *TalentProjectSubmissionService) AddTalentProjectSubmission(submission *talent.ProjectSubmission) error {

	// Check if foreign keys exist.
	fmt.Println("=========================", submission)
	err := service.doesForeignKeysExist(submission, submission.CreatedBy)
	if err != nil {
		return err
	}

	// assign tenantID and createdBy to uploads slice if it exist.
	if submission.ProjectSubmissionUploads != nil {
		for index := range submission.ProjectSubmissionUploads {
			submission.ProjectSubmissionUploads[index].TenantID = submission.TenantID
			submission.ProjectSubmissionUploads[index].CreatedBy = submission.CreatedBy
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

// UpdateTalentProjectSubmission will update submitted project of talent.
func (service *TalentProjectSubmissionService) UpdateTalentProjectSubmission(submission *talent.ProjectSubmission) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeysExist(submission, submission.UpdatedBy)
	if err != nil {
		return err
	}

	// Check if submission exist.
	err = service.doesSubmissionExist(submission.TenantID, submission.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	tempSubmission := talent.ProjectSubmission{}

	err = service.repository.GetRecordForTenant(uow, submission.TenantID, &tempSubmission,
		repository.Filter("`id` = ?", submission.ID), repository.Select("`created_by`,`created_at`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	submission.CreatedBy = tempSubmission.CreatedBy
	submission.SubmittedOn = tempSubmission.CreatedAt
	// submission.ProjectSubmissionUploads = nil

	// err = service.updateProjectSubmissionUpload(uow, submission, submission.TenantID, submission.UpdatedBy)
	// if err != nil {
	// 	return err
	// }

	err = service.repository.Save(uow, submission)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ScoreTalentSubmission will allow faculty to give scores to the latest submission.
func (service *TalentProjectSubmissionService) ScoreTalentProjectSubmission(submission *talent.ProjectSubmission) error {

	// Check if foreign keys exist.
	err := service.doesForeignKeysExist(submission, submission.UpdatedBy)
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

	tempSubmission := &talent.ProjectSubmission{}

	err = service.repository.GetRecordForTenant(uow, submission.TenantID, tempSubmission,
		repository.Filter("`id` = ?", submission.ID),
		repository.Select("`batch_project_id`,`talent_id`,`submitted_on`,`tenant_id`"))
	if err != nil {
		return err
	}

	err = service.validateLatestSubmission(tempSubmission)
	if err != nil {
		return err
	}

	if submission.IsAccepted != nil && *submission.IsAccepted {
		err = service.addProjectRating(uow, submission)
		if err != nil {
			return err
		}
	}

	err = service.repository.UpdateWithMap(uow, submission, map[string]interface{}{
		"BatchID":          submission.BatchID,
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

// GetTalentProjectSubmissions will return talent project submissions.
func (service *TalentProjectSubmissionService) GetTalentProjectSubmissions(tenantID, batchID, talentID, projectID uuid.UUID,
	talentSubmissions *[]talent.ProjectSubmissionDTO, parser *web.Parser) error {

	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchProjectExist(tenantID, projectID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesBatchExist(tenantID, batchID)
	if err != nil {
		return err
	}

	// Check if batch session exist.
	err = service.doesTalentExist(tenantID, talentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.repository.GetAllForTenant(uow, tenantID, talentSubmissions,
		repository.Filter("`batch_project_id` = ?", projectID),
		repository.Filter("`talent_id` = ?", talentID),
		repository.Filter("`batch_id` = ?", batchID),
		repository.OrderBy("`submitted_on` DESC"),
		service.addSearchQueries(parser.Form), repository.PreloadAssociations(service.associations))
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// GetAllAssignmentsSubmissionsOfTalent will return batch topic assignments for one talent.
func (service *TalentProjectSubmissionService) GetAllProjectsSubmissionsOfTalent(tenantID, batchID, talentID uuid.UUID,
	project *[]batch.ProjectDTO, parser *web.Parser) error {

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

	// Check if batch session exist.
	err = service.doesTalentExist(tenantID, talentID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// get all batch project.
	err = service.repository.GetAllForTenant(uow, tenantID, project,
		repository.Filter("`batch_projects`.`batch_id` = ?", batchID),
		repository.Filter("`batch_projects`.`assigned_date` IS NOT NULL"),
		repository.OrderBy("`batch_projects`.`assigned_date`"),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ProgrammingProject"},
			repository.Preload{Schema: "Submissions",
				Queryprocessors: []repository.QueryProcessor{
					repository.Filter("`talent_project_submissions`.`talent_id` = ?", talentID),
					repository.OrderBy("`talent_project_submissions`.`submitted_on` DESC"),
					repository.PreloadAssociations([]string{"ProjectSubmissionUpload", "ProgrammingProjectRatings",
						"ProgrammingProjectRatings.ProgrammingProjectRatingParameter"}),
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

// GetAllProjectsWithSubmissions will return batch_project with scores.
func (service *TalentProjectSubmissionService) GetAllProjectsWithSubmissions(tenantID, batchID uuid.UUID,
	topicProject *[]batch.ProjectDTO, parser *web.Parser) error {

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

	err = service.repository.GetAllForTenant(uow, tenantID, topicProject,
		repository.Filter("`batch_projects`.`batch_id` = ?", batchID),
		// Need to move below join & optimize code #Niranjan.
		service.addSearchQueries(parser.Form),
		repository.PreloadWithCustomCondition(repository.Preload{Schema: "ProgrammingProject"},
			repository.Preload{Schema: "Submissions",
				Queryprocessors: []repository.QueryProcessor{
					repository.Join("INNER JOIN `talents` ON `talents`.`id` = `talent_project_submissions`.`talent_id`"),
					repository.OrderBy("`talents`.`first_name`,`talent_project_submissions`.`submitted_on` DESC"),
					repository.PreloadAssociations([]string{"Talent"}),
				},
			},
		),
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
func (service *TalentProjectSubmissionService) doesForeignKeysExist(submission *talent.ProjectSubmission,
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

	//check if batchID exit
	err = service.doesBatchExist(submission.TenantID, submission.BatchID)
	if err != nil {
		return err
	}

	//check if projectID exit
	err = service.doesBatchProjectExist(submission.TenantID, submission.BatchProjectID)
	if err != nil {
		return err
	}

	return nil
}

// returns error if there is no tenant record in table.
func (service *TalentProjectSubmissionService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no credential record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no batch_assignment record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesBatchProjectExist(tenantID, projectID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.Project{},
		repository.Filter("`id` = ?", projectID))
	if err := util.HandleError("Invalid batch session project ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no faculty record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.FacultyDetails{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if there is no talent project submission record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesSubmissionExist(tenantID, submissionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, talent.ProjectSubmission{},
		repository.Filter("`id` = ?", submissionID))
	if err := util.HandleError("Invalid submission ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *TalentProjectSubmissionService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}
	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	if dueDate, ok := requestForm["dueDate"]; ok {
		util.AddToSlice("`batch_projects`.`due_date`", "<= ?", "AND", dueDate, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

func (service *TalentProjectSubmissionService) updateProjectSubmissionUpload(uow *repository.UnitOfWork, submission *talent.ProjectSubmission,
	tenantID, credentialID uuid.UUID) error {

	projectSubmissionUploadeds := []talent.ProjectSubmissionUpload{}
	projectSubmissionUploadMap := make(map[uuid.UUID]uint)

	err := service.repository.GetAllForTenant(uow, tenantID, &projectSubmissionUploadeds,
		repository.Filter("`project_submission_id` = ?", submission.ID))
	if err != nil {
		return err
	}

	for _, tempProjectSubmissionUploaded := range projectSubmissionUploadeds {
		projectSubmissionUploadMap[tempProjectSubmissionUploaded.ID] = 1
	}

	for _, projectSubmissionUpload := range submission.ProjectSubmissionUploads {
		// if it is not an existing programmingConcept
		if util.IsUUIDValid(projectSubmissionUpload.ID) {
			projectSubmissionUploadMap[projectSubmissionUpload.ID]++
		}

		// new entry of programmingConcept
		if !util.IsUUIDValid(projectSubmissionUpload.ID) {

			projectSubmissionUpload.CreatedBy = credentialID
			projectSubmissionUpload.ProjectSubmissionID = submission.ID
			projectSubmissionUpload.TenantID = tenantID

			// adding new record to the DB
			err := service.repository.Add(uow, &projectSubmissionUpload)
			if err != nil {
				return err
			}
		}

		if projectSubmissionUploadMap[projectSubmissionUpload.ID] > 1 {

			// check foreign keys of programmingConcept
			projectSubmissionUpload.TenantID = tenantID
			projectSubmissionUpload.ProjectSubmissionID = submission.ID
			projectSubmissionUpload.UpdatedBy = credentialID

			err := service.repository.Update(uow, &projectSubmissionUpload)
			if err != nil {
				return nil
			}

			projectSubmissionUploadMap[projectSubmissionUpload.ID] = 0
			continue
		}

	}
	submission.ProjectSubmissionUploads = nil
	return nil
}

// returns error if given submission is not the latest submission.
func (service *TalentProjectSubmissionService) validateLatestSubmission(submission *talent.ProjectSubmission) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, submission.TenantID, submission,
		repository.Filter("`talent_id` = ?", submission.TalentID),
		repository.Filter("`batch_project_id` = ?", submission.BatchProjectID),
		repository.Filter("`submitted_on` > ?", submission.SubmittedOn))
	if err := util.HandleIfExistsError("Only latest submission can be updated", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

func (service *TalentProjectSubmissionService) addProjectRating(uow *repository.UnitOfWork, submission *talent.ProjectSubmission) error {
	for _, projectRating := range submission.ProgrammingProjectRatings {
		err := service.doesProjectRatingParameterExist(submission.TenantID, projectRating.ProgrammingProjectRatingParameterID)
		if err != nil {
			return err
		}
		projectRating.TalentID = submission.TalentID
		projectRating.TalentSubmissionID = submission.ID
		projectRating.TenantID = submission.TenantID
		projectRating.BatchID = submission.BatchID

		err = projectRating.Validate()
		if err != nil {
			return err
		}

		err = service.validateProjectRating(submission)
		if err != nil {
			return err
		}

		// Adding ratings
		err = service.repository.Add(uow, projectRating)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns error if there is no project rating record in table for the given tenant.
func (service *TalentProjectSubmissionService) doesProjectRatingParameterExist(tenantID, ProgrammingProjectRatingParameterID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, batch.ProgrammingProjectRatingParameter{},
		repository.Filter("`id` = ?", ProgrammingProjectRatingParameterID))
	if err := util.HandleError("Invalid project_parameter_id ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// returns error if given submission's project rating are already rated.
func (service *TalentProjectSubmissionService) validateProjectRating(submission *talent.ProjectSubmission) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, submission.TenantID, batch.ProgrammingProjectRating{},
		repository.Filter("`talent_id` = ?", submission.TalentID),
		repository.Filter("`talent_submission_id` = ?", submission.ID))
	if err := util.HandleIfExistsError("Already rated for this submission", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
