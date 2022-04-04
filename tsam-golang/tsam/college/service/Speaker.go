package service

import (
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// SpeakerService provide method to update, delete, add, get method for speaker.
type SpeakerService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewSpeakerService returns new instance of SpeakerService.
func NewSpeakerService(db *gorm.DB, repository repository.Repository) *SpeakerService {
	return &SpeakerService{
		DB:         db,
		Repository: repository,
	}
}

// AddSpeaker add new speaker to database.
func (service *SpeakerService) AddSpeaker(speaker *college.Speaker, uows ...*repository.UnitOfWork) error {
	// Validate tenant id.
	err := service.doesTenantExist(speaker.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// // Validate if speaker name and company exists.
	// err = service.doesSpeakerExist(speaker)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// Validate credential ID.
	err = service.doesCredentialExist(speaker.CreatedBy, speaker.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Add Speaker.
	err = service.Repository.Add(uow, speaker)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	if length == 0 {
		uow.Commit()
	}
	return nil
}

// AddSpeakers adds multiple speakers to database.
func (service *SpeakerService) AddSpeakers(speakers *[]college.Speaker, speakerIDs *[]uuid.UUID, tenantID, credentialID uuid.UUID) error {
	// // Check for same name conflict.
	// for i := 0; i < len(*speakers); i++ {
	// 	for j := 0; j < len(*speakers); j++ {
	// 		if i != j && (*speakers)[i].Name == (*speakers)[j].Name {
	// 			log.NewLogger().Error("Name:" + (*speakers)[j].Name + " exists")
	// 			return errors.NewValidationError("Name:" + (*speakers)[j].Name + " exists")
	// 		}
	// 	}
	// }

	// Add individual speaker.
	uow := repository.NewUnitOfWork(service.DB, false)
	for _, speaker := range *speakers {
		speaker.TenantID = tenantID
		speaker.CreatedBy = credentialID
		err := service.AddSpeaker(&speaker, uow)
		if err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
		*speakerIDs = append(*speakerIDs, speaker.ID)
	}

	uow.Commit()
	return nil
}

// UpdateSpeaker updates speaker in database.
func (service *SpeakerService) UpdateSpeaker(speaker *college.Speaker) error {
	// Validate tenant ID.
	err := service.doesTenantExist(speaker.TenantID)
	if err != nil {
		return err
	}

	// Validate speaker ID.
	err = service.doesSpeakerIDExist(speaker.ID, speaker.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(speaker.UpdatedBy, speaker.TenantID)
	if err != nil {
		return err
	}

	// // Validate if speaker name and company exists.
	// err = service.doesSpeakerExist(speaker)
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	return err
	// }

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update speaker.
	err = service.Repository.Update(uow, speaker)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteSpeaker deletes speaker from database.
func (service *SpeakerService) DeleteSpeaker(speaker *college.Speaker) error {
	credentialID := speaker.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(speaker.TenantID)
	if err != nil {
		return err
	}

	// Validate speaker ID.
	err = service.doesSpeakerIDExist(speaker.ID, speaker.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, speaker.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update speaker for updating deleted_by and deleted_at fields of speaker
	if err := service.Repository.UpdateWithMap(uow, speaker, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", speaker.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Speaker could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetSpeakerList returns speaker list.
func (service *SpeakerService) GetSpeakerList(speakers *[]college.Speaker, tenantID uuid.UUID) error {
	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all speakers.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, speakers, "`first_name`, `last_name`")
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetSpeaker returns all speakers.
func (service *SpeakerService) GetSpeaker(speaker *college.SpeakerDTO) error {
	// Validate tenant id.
	err := service.doesTenantExist(speaker.TenantID)
	if err != nil {
		return err
	}

	// Validate speaker id.
	err = service.doesSpeakerIDExist(speaker.ID, speaker.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one speaker by id.
	err = service.Repository.GetForTenant(uow, speaker.TenantID, speaker.ID, speaker)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetSpeakers returns specific speaker by id.
func (service *SpeakerService) GetSpeakers(speakers *[]college.SpeakerDTO, tenantID uuid.UUID, form url.Values, limit,
	offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all speakers form database.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, speakers, "`first_name`, `last_name`",
		service.addSearchQueries(form),
		repository.PreloadAssociations([]string{"Designation"}),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// addSearchQueries adds search queries.
func (service *SpeakerService) addSearchQueries(requestForm url.Values) repository.QueryProcessor {

	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	// First name.
	if _, ok := requestForm["firstName"]; ok {
		util.AddToSlice("`first_name`", "LIKE ?", "AND", "%"+requestForm.Get("firstName")+"%", &columnNames, &conditions, &operators, &values)

	}

	// Last name.
	if _, ok := requestForm["lastName"]; ok {
		util.AddToSlice("`last_name`", "LIKE ?", "AND", "%"+requestForm.Get("lastName")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Company.
	if _, ok := requestForm["company"]; ok {
		util.AddToSlice("`company`", "LIKE ?", "AND", "%"+requestForm.Get("company")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Designation.
	if designationID, ok := requestForm["designationID"]; ok {
		util.AddToSlice("`designation_id`", "=?", "AND", designationID, &columnNames, &conditions, &operators, &values)
	}

	// Minimum experience.
	if minimumExperience, ok := requestForm["minimumExperience"]; ok {
		util.AddToSlice("`experience_in_years`", ">=?", "AND", minimumExperience, &columnNames, &conditions, &operators, &values)
	}

	// Maximum experience.
	if maximumExperience, ok := requestForm["maximumExperience"]; ok {
		util.AddToSlice("`experience_in_years`", "<=?", "AND", maximumExperience, &columnNames, &conditions, &operators, &values)
	}

	return repository.FilterWithOperator(columnNames, conditions, operators, values)
}

// // doesSpeakerIDExist returns true if the speaker first name, last name and coompany
// // already exists for the speaker in database.
// func (service *SpeakerService) doesSpeakerExist(speaker *college.Speaker) error {

// 	// Check for same speaker first name, last name and company conflict.
// 	exists, err := repository.DoesRecordExistForTenant(service.DB, speaker.TenantID, &college.Speaker{},
// 		repository.Filter("first_name=? AND last_name=? AND company=? AND id!=?", speaker.FirstName,
// 			speaker.LastName, speaker.Company, speaker.ID))
// 	if err := util.HandleIfExistsError("Name:"+speaker.FirstName+" "+speaker.LastName+"exists for"+*speaker.Company, exists, err); err != nil {
// 		log.NewLogger().Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// doesTenantExist validates if tenant exists or not in database.
func (service *SpeakerService) doesTenantExist(tenantID uuid.UUID) error {
	// Check if tenant exists or not.
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesSpeakerIDExist validates if speaker exists or not in database.
func (service *SpeakerService) doesSpeakerIDExist(speakerID uuid.UUID, tenantID uuid.UUID) error {
	// Check speaker exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Speaker{},
		repository.Filter("`id` = ?", speakerID))
	if err := util.HandleError("Invalid speaker ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *SpeakerService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	// Check credential exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
