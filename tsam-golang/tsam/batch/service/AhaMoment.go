package service

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// AhaMomentService provides methods to do different CRUD operations on feeling table.
type AhaMomentService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewAhaMomentService returns a new instance Of AhaMomentService.
func NewAhaMomentService(db *gorm.DB, repository repository.Repository) *AhaMomentService {
	return &AhaMomentService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Batch", "BatchTopic", "Faculty",
			"Feeling", "FeelingLevel",
			"Talent",
			"AhaMomentResponse", "AhaMomentResponse.Question",
		},
	}
}

// AddAhaMoment will add ahaMoment and responses for the ahaMoment.
func (ser *AhaMomentService) AddAhaMoment(ahaMoment *batch.AhaMoment, uows ...*repository.UnitOfWork) error {

	// Extracts all the ID from the objects
	ser.extractAhaMomentID(ahaMoment)

	// Create new unit of work, if no transaction has been passed to the function.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {

		// Check if tenant exist.
		err := ser.doesTenantExist(ahaMoment.TenantID)
		if err != nil {
			return err
		}

		// Check if credential exist.
		err = ser.doesCredentialExist(ahaMoment.TenantID, ahaMoment.CreatedBy)
		if err != nil {
			return err
		}

		uow = repository.NewUnitOfWork(ser.DB, false)
	} else {
		uow = uows[0]

	}

	// Check if all foreign keys exist.
	err := ser.doesForeignKeyExist(ahaMoment)
	if err != nil {
		return nil
	}

	for index := range ahaMoment.AhaMomentResponse {

		ser.assignAhaMomentResponse(&ahaMoment.AhaMomentResponse[index], ahaMoment)

		err = ser.doesFeedbackQuestionExist(ahaMoment.TenantID, ahaMoment.AhaMomentResponse[index].QuestionID)
		if err != nil {
			return err
		}
	}

	err = ser.Repository.Add(uow, ahaMoment)
	if err != nil {
		if length == 0 {
			uow.RollBack()
		}
		return err
	}
	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddAhaMoments will add multiple ahaMoments and its response to the table.
func (ser *AhaMomentService) AddAhaMoments(ahaMoments *[]batch.AhaMoment,
	tenantID, batchID, batchTopicID, credentialID uuid.UUID) error {

	// Check if tenant exist
	err := ser.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential exist
	err = ser.doesCredentialExist(tenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if credential is of faculty
	exist, err := repository.DoesRecordExist(ser.DB, general.Credential{},
		repository.Join("LEFT JOIN faculties ON credentials.`faculty_id` = faculties.`id` AND "+
			"credentials.`tenant_id` = faculties.`tenant_id`"), repository.Filter("faculties.`deleted_at` IS NULL"),
		repository.Filter("credentials.`id` = ? AND credentials.`tenant_id` = ?", credentialID, tenantID))
	if err := util.HandleError("Aha moments can only be added by faculty", exist, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	for index := range *ahaMoments {

		(*ahaMoments)[index].BatchID = batchID
		(*ahaMoments)[index].BatchSessionID = batchTopicID
		(*ahaMoments)[index].TenantID = tenantID
		(*ahaMoments)[index].CreatedBy = credentialID
		err := ser.AddAhaMoment(&(*ahaMoments)[index], uow)
		if err != nil {
			uow.RollBack()
			return err
		}
	}
	uow.Commit()
	return nil
}

// DeleteAhaMoment will delete aha moment and its responses.
func (ser *AhaMomentService) DeleteAhaMoment(ahaMoment *batch.AhaMoment) error {

	// Check if tenant exist
	err := ser.doesTenantExist(ahaMoment.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist
	err = ser.doesCredentialExist(ahaMoment.TenantID, ahaMoment.DeletedBy)
	if err != nil {
		return err
	}

	// Check if ahaMoment exist
	err = ser.doesAhaMomentExist(ahaMoment.TenantID, ahaMoment.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.DB, false)

	err = ser.Repository.UpdateWithMap(uow, batch.AhaMomentResponse{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": ahaMoment.DeletedBy,
	}, repository.Filter("`aha_moment_id` = ?", ahaMoment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = ser.Repository.UpdateWithMap(uow, batch.AhaMoment{}, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": ahaMoment.DeletedBy,
	}, repository.Filter("`id` = ?", ahaMoment.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllAhaMoments will return all the aha moments and its response for specified batch and session
func (ser *AhaMomentService) GetAllAhaMoments(ahaMoments *[]batch.AhaMomentDTO,
	tenantID, batchID, sessionID uuid.UUID) error {

	// Check if tenant exist.
	err := ser.doesTenantExist(tenantID)
	if err != nil {
		return nil
	}

	uow := repository.NewUnitOfWork(ser.DB, true)

	err = ser.Repository.GetAll(uow, ahaMoments,
		repository.Filter("aha_moments.`batch_id` = ? AND aha_moments.`batch_topic_id` = ?", batchID, sessionID),
		repository.PreloadAssociations(ser.association),
		repository.Join("INNER JOIN talents ON aha_moments.`talent_id`=talents.`id`"+
			" AND aha_moments.`tenant_id` = talents.`tenant_id`"),
		repository.Filter("aha_moments.`tenant_id` = ? AND talents.`deleted_at` IS NULL", tenantID),
		repository.OrderBy("talents.`first_name`"))
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

// doesForeignKeyExist will check if all the foreign keys are valid.
func (ser *AhaMomentService) doesForeignKeyExist(ahaMoment *batch.AhaMoment) error {

	// Check if batch exist.
	err := ser.doesBatchExist(ahaMoment.TenantID, ahaMoment.BatchID)
	if err != nil {
		return err
	}

	// Check if batch-session exist.
	// err = ser.doesBatchTopicExist(ahaMoment.TenantID, ahaMoment.BatchID, ahaMoment.BatchTopicID)
	// if err != nil {
	// 	return err
	// }

	// Check if faculty exist.
	err = ser.doesFacultyExist(ahaMoment.TenantID, ahaMoment.FacultyID)
	if err != nil {
		return err
	}

	// Check if talent exist.
	err = ser.doesTalentExist(ahaMoment.TenantID, ahaMoment.TalentID)
	if err != nil {
		return err
	}

	// Check if feeling exist.
	err = ser.doesFeelingExist(ahaMoment.TenantID, ahaMoment.FeelingID)
	if err != nil {
		return err
	}

	// Check if feeling level exist.
	err = ser.doesFeelingLevelExist(ahaMoment.TenantID, ahaMoment.FeelingID, ahaMoment.FeelingLevelID)
	if err != nil {
		return err
	}

	return nil
}

// extractAhaMomentID will extract ID's from ahaMoment.
func (ser *AhaMomentService) extractAhaMomentID(ahaMoment *batch.AhaMoment) {

	ahaMoment.FeelingID = ahaMoment.Feeling.ID
	ahaMoment.FeelingLevelID = ahaMoment.FeelingLevel.ID
}

// assignAhaMomentResponse will extract ID's from ahaMomentResponse.
func (ser *AhaMomentService) assignAhaMomentResponse(response *batch.AhaMomentResponse, ahaMoment *batch.AhaMoment) {

	response.QuestionID = response.Question.ID
	response.BatchID = ahaMoment.BatchID
	response.BatchSessionID = ahaMoment.BatchSessionID
	response.TalentID = ahaMoment.TalentID
	response.FacultyID = ahaMoment.FacultyID
	response.TenantID = ahaMoment.TenantID
	response.CreatedBy = ahaMoment.CreatedBy
	response.FeelingID = ahaMoment.FeelingID
	response.FeelingLevelID = ahaMoment.FeelingLevelID
}

// Returns error if there is no tenant record in table.
func (ser *AhaMomentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(ser.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no credential record in table for the given tenant.
func (ser *AhaMomentService) doesCredentialExist(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no aha-moment record in table for the given tenant.
func (ser *AhaMomentService) doesAhaMomentExist(tenantID, momentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.AhaMoment{},
		repository.Filter("`id` = ?", momentID))
	if err := util.HandleError("Invalid Aha moment ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no batch record in table for the given tenant.
func (ser *AhaMomentService) doesBatchExist(tenantID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Invalid batch ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no batch session record in table for the given tenant.
// func (ser *AhaMomentService) doesBatchTopicExist(tenantID, batchID, batchTopicID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.BatchTopic{},
// 		repository.Filter("`id` = ? AND `batch_id` =? ", batchTopicID, batchID))
// 	if err := util.HandleError("Invalid batch-session ID", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// Returns error if there is no faculty record in table for the given tenant.
func (ser *AhaMomentService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, faculty.Faculty{},
		repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Invalid faculty ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no talent record in table for the given tenant.
func (ser *AhaMomentService) doesTalentExist(tenantID, talentID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, talent.Talent{},
		repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Invalid talent ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no feedback question record in table for the given tenant.
func (ser *AhaMomentService) doesFeedbackQuestionExist(tenantID, feedbackQuestionID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.FeedbackQuestion{},
		repository.Filter("`id` = ?", feedbackQuestionID))
	if err := util.HandleError("Invalid feedback question ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no feeling record in table for the given tenant.
func (ser *AhaMomentService) doesFeelingExist(tenantID, feelingID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.Feeling{},
		repository.Filter("`id` = ?", feelingID))
	if err := util.HandleError("Invalid feeling ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// Returns error if there is no feeling level record in table for the given tenant.
func (ser *AhaMomentService) doesFeelingLevelExist(tenantID, feelingID, feelingLevelID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, general.FeelingLevel{},
		repository.Filter("`id` = ? AND `feeling_id` = ?", feelingLevelID, feelingID))
	if err := util.HandleError("Invalid feeling level ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// // Returns error if there is no talent and faculty record in table for the given tenant.
// func (ser *AhaMomentService) doesTalentAndFacultyExist(tenantID, talentID, facultyID uuid.UUID) error {
// 	exists, err := repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.AhaMoment{},
// 		repository.Filter("`talent_id` = ? AND `faculty_id` = ?", talentID, facultyID))
// 	if err := util.HandleError("Talent and faculty does not exist in aha moment", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}

// 	exists, err = repository.DoesRecordExistForTenant(ser.DB, tenantID, batch.AhaMomentResponse{},
// 		repository.Filter("`talent_id` = ? AND `faculty_id` = ?", talentID, facultyID))
// 	if err := util.HandleError("Talent and faculty does not exist in aha moment response", exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}

// 	return nil
// }
