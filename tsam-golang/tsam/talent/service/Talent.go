package service

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"

	"github.com/techlabs/swabhav/tsam/errors"
	genService "github.com/techlabs/swabhav/tsam/general/service"
	"github.com/techlabs/swabhav/tsam/log"
	bat "github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/models/college"
	"github.com/techlabs/swabhav/tsam/models/company"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	fclt "github.com/techlabs/swabhav/tsam/models/faculty"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/models/talent"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentService provides methods to update, delete, add, get all, get all by salesperson,
// get one by id and get all by company requiremnt for talent.
type TalentService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// TalentAssociationNames provides preload associations array for talent.
var TalentAssociationNames []string = []string{
	"Academics", "Experiences", "Technologies", "State", "Country", "Academics.Degree",
	"Experiences.Technologies", "Experiences.Designation", "Academics.Specialization", "SalesPerson",
	"TalentSource", "MastersAbroad", "MastersAbroad.Scores", "MastersAbroad.Scores.Examination",
	"MastersAbroad.Degree", "MastersAbroad.Universities", "MastersAbroad.Countries",
}

// BatchAssociations provides preload associations array for batch
var BatchAssociations []string = []string{
	"Course", "Course.Eligibility",
	"Course.Eligibility.Technologies",
	"Eligibility", "Eligibility.Technologies",
	// "Faculty",
	// "Timing", "Timing.Day",
	"SalesPerson",
}

// NewTalentService returns new instance of TalentServcie.
func NewTalentService(db *gorm.DB, repository repository.Repository) *TalentService {
	return &TalentService{
		DB:         db,
		Repository: repository,
	}
}

// AddTalent adds new talent to database.
func (service *TalentService) AddTalent(talent *tal.Talent, isForm bool, uows ...*repository.UnitOfWork) (bool, error) {

	// Validate compulsary fields.
	if err := talent.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		return false, err
	}

	// Get credential id from CreatedBy field of talent(set in controller).
	credentialID := talent.CreatedBy

	// Starting new transaction.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// If talent is getting added from campus regsitration form then skip credential id validation and setting.
	if !isForm {
		// Validate credential id.
		if err := service.doesCredentialExist(credentialID, talent.TenantID); err != nil {
			return false, err
		}

		// Check if login id's role name is salesperson or not.
		isRecordNotFound := false
		tempUser := general.User{}
		err := service.Repository.GetRecord(uow, &tempUser,
			repository.Select("users.`id`"),
			repository.Join("join credentials on credentials.`sales_person_id` = users.`id`"),
			repository.Filter("credentials.`id`=? AND credentials.`tenant_id`=? AND users.`tenant_id`=?",
				credentialID, talent.TenantID, talent.TenantID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				isRecordNotFound = true
			} else {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}
		}
		if !isRecordNotFound {
			talent.SalesPersonID = &tempUser.ID
		}

		// Set createdBy field to all academics.
		if talent.Academics != nil && len(talent.Academics) != 0 {
			for i := 0; i < len(talent.Academics); i++ {
				//set created_by field of academics
				talent.Academics[i].CreatedBy = credentialID
			}
		}

		// Give 01 to all talent's experiences' fromDate and toDate field and set createdBy field.
		if talent.Experiences != nil && len(talent.Experiences) != 0 {
			for i := 0; i < len(talent.Experiences); i++ {
				//set created_by field of experiences
				talent.Experiences[i].CreatedBy = credentialID
			}
		}

		// Give masters abroad and its score arrays created by field.
		if talent.MastersAbroad != nil {
			talent.MastersAbroad.CreatedBy = credentialID
			for i := 0; i < len(talent.MastersAbroad.Scores); i++ {
				talent.MastersAbroad.Scores[i].CreatedBy = credentialID
			}
		}
	}

	// Extract foreign key IDs and remove the object.
	service.extractID(talent)

	// Give tenant id to academics and experiences.
	service.setTenantID(talent, talent.TenantID)

	// Validate tenant id.
	if err := service.doesTenantExist(talent.TenantID); err != nil {
		return false, err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(talent); err != nil {
		return false, err
	}

	// Assign Talent Code.
	var codeError error
	talent.Code, codeError = util.GenerateUniqueCode(uow.DB, talent.FirstName, "`code` = ?", &tal.Talent{})
	if codeError != nil {
		log.NewLogger().Error(codeError.Error())
		uow.RollBack()
		return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, talent); err != nil {
		uow.RollBack()
		return false, err
	}

	// Give 01 to all talent's experiences' fromDate and toDate field and set createdBy field.
	if talent.Experiences != nil && len(talent.Experiences) != 0 {
		for i := 0; i < len(talent.Experiences); i++ {
			//set created_by field of experiences
			if len(talent.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(talent.Experiences[i].FromDate))
			}
			if talent.Experiences[i].ToDate != nil && len(*talent.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((talent.Experiences[i].ToDate))
			}
		}
	}

	// Call function to add talent and credential.
	err := service.AddTalentAndCredential(uow, talent)
	if err != nil {
		// If email already exists then send true value.
		if err.Error()[:6] == "Email:" && err.Error()[len(err.Error())-14:] == "already exists" || err.Error() == "Email already exists" {
			return true, err
		}
		return false, err
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}
	// uow.RollBack()

	return false, nil
}

// AddTalents adds multiple talents to database.
func (service *TalentService) AddTalents(talents *[]tal.TalentExcel, talentAddedCount *int,
	tenantID uuid.UUID, errorList *[]error, credentialID uuid.UUID) error {

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get IDs for all values and add individual talent to database.
	if talents != nil && len(*talents) != 0 {

		// ExcelLoop is a label for outer loop
	ExcelLoop:

		// Loop the talents in talents from excel.
		for _, talentExcel := range *talents {

			// Create bucket for talent.
			talent := tal.Talent{}

			// Get country ID.
			if talentExcel.CountryName != nil {

				// Create bucket for country ID.
				countryID := tal.IDModel{}

				// Get country ID from database.
				err := service.Repository.GetRecordForTenant(uow, tenantID, &countryID,
					repository.Select("`id`"),
					repository.Filter("`name`=?", talentExcel.CountryName),
					repository.Filter("`tenant_id`=?", tenantID),
					repository.Table("countries"))

				if err != nil {
					// More efficiency needed #Niranjan
					log.NewLogger().Error(err.Error())
					errorString := ""
					if err == gorm.ErrRecordNotFound {
						errorString = talentExcel.Email + " : Invalid country name"
						err = errors.NewHTTPError(errorString, http.StatusBadRequest)
					} else {
						errorString = talentExcel.Email + " : Internal Server Error"
						err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
					}
					*errorList = append(*errorList, err)
					continue
				}

				// Give country ID to talent.
				talent.CountryID = &countryID.ID
			}

			// Get state ID.
			if talentExcel.StateName != nil {

				// Create bucket for state ID.
				stateID := tal.IDModel{}

				// Get state ID from database.
				err := service.Repository.GetRecordForTenant(uow, tenantID, &stateID,
					repository.Select("`id`"),
					repository.Filter("`name`=?", talentExcel.StateName),
					repository.Filter("`tenant_id`=?", tenantID),
					repository.Table("states"))

				if err != nil {
					log.NewLogger().Error(err.Error())
					errorString := ""
					if err == gorm.ErrRecordNotFound {
						errorString = talentExcel.Email + " : Invalid state name"
						err = errors.NewHTTPError(errorString, http.StatusBadRequest)
					} else {
						errorString = talentExcel.Email + " : Internal Server Error"
						err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
					}
					*errorList = append(*errorList, err)
					continue
				}

				// Give state ID to talent.
				talent.StateID = &stateID.ID
			}

			// Check if academics exist or not.
			if talentExcel.Academics != nil && len(talentExcel.Academics) > 0 {

				// Create bucket for academics.
				academics := []*tal.Academic{}

				// Loop the academics of talent excel.
				for _, academicExcel := range talentExcel.Academics {

					// Create bucket for academic.
					academic := tal.Academic{}

					// Create bucket for degree ID.
					degreeID := tal.IDModel{}

					// Get degree ID from database.
					err := service.Repository.GetRecordForTenant(uow, tenantID, &degreeID,
						repository.Select("`id`"),
						repository.Filter("`name`=?", academicExcel.DegreeName),
						repository.Filter("`tenant_id`=?", tenantID),
						repository.Table("degrees"))

					if err != nil {
						log.NewLogger().Error(err.Error())
						errorString := ""
						if err == gorm.ErrRecordNotFound {
							errorString = talentExcel.Email + " : Invalid degree name"
							err = errors.NewHTTPError(errorString, http.StatusBadRequest)
						} else {
							errorString = talentExcel.Email + " : Internal Server Error"
							err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
						}
						*errorList = append(*errorList, err)
						continue ExcelLoop
					}

					// Give degree ID to academic.
					academic.DegreeID = degreeID.ID

					// Create bucket for specialization ID.
					specializationID := tal.IDModel{}

					// Get specialization ID from database.
					err = service.Repository.GetRecordForTenant(uow, tenantID, &specializationID,
						repository.Select("`id`"),
						repository.Filter("`branch_name`=?", academicExcel.SpecializationName),
						repository.Filter("`degree_id`=?", degreeID.ID),
						repository.Filter("`tenant_id`=?", tenantID),
						repository.Table("specializations"))

					if err != nil {
						log.NewLogger().Error(err.Error())
						errorString := ""
						if err == gorm.ErrRecordNotFound {
							errorString = talentExcel.Email + " : Invalid specialization name"
							err = errors.NewHTTPError(errorString, http.StatusBadRequest)
						} else {
							errorString = talentExcel.Email + " : Internal Server Error"
							err = errors.NewHTTPError(errorString, http.StatusInternalServerError)
						}
						*errorList = append(*errorList, err)
						continue ExcelLoop
					}

					// Give specialization ID to academic.
					academic.SpecializationID = specializationID.ID

					// Give college name to academic.
					academic.College = academicExcel.CollegeName

					// Give percentage to academic.
					academic.Percentage = academicExcel.Percentage

					// Give year of passout to academic.
					academic.Passout = academicExcel.YearOfPassout

					// Push academic into academics.
					academics = append(academics, &academic)
				}

				// Give academics to talent.
				talent.Academics = academics
			}

			// Give first name to talent.
			talent.FirstName = talentExcel.FirstName

			// Give last name to talent.
			talent.LastName = talentExcel.LastName

			// Give email to talent.
			talent.Email = talentExcel.Email

			// Give contact to talent.
			talent.Contact = talentExcel.Contact

			// Give academic year to talent.
			talent.AcademicYear = talentExcel.AcademicYear

			// Give is active to talent.
			talent.IsActive = talentExcel.IsActive

			// Give is swabhav talent to talent.
			talent.IsSwabhavTalent = talentExcel.IsSwabhavTalent

			// Give address to talent.
			talent.Address.Address = talentExcel.Address

			// Give city to talent.
			talent.City = talentExcel.City

			// Give Pin code to talent.
			talent.PINCode = talentExcel.PINCode

			// Give tenant ID to talent.
			talent.TenantID = tenantID

			// Give created by to talent.
			talent.CreatedBy = credentialID

			// Add talent individually.
			if _, err := service.AddTalent(&talent, false); err != nil {

				// If error then push it in error list.
				errorString := talentExcel.Email + " : " + err.Error()
				er := errors.NewHTTPError(errorString, http.StatusBadRequest)
				*errorList = append(*errorList, er)
				continue
			}

			// Increment count of talents added successfully.
			*talentAddedCount++
		}
	}
	return nil
}

// AddTalentFromExcel adds one talent from excel to database.
func (service *TalentService) AddTalentFromExcel(talentExcel *tal.TalentExcel, tenantID, credentialID uuid.UUID) error {

	// Start new transaction
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get IDs for all values and add talent to database.
	// Create bucket for talent.
	talent := tal.Talent{}

	// Get country ID.
	if talentExcel.CountryName != nil {

		// Create bucket for country ID.
		countryID := tal.IDModel{}

		// Get country ID from database.
		err := service.Repository.GetRecordForTenant(uow, tenantID, &countryID,
			repository.Select("`id`"),
			repository.Filter("`name`=?", talentExcel.CountryName),
			repository.Filter("`tenant_id`=?", tenantID),
			repository.Table("countries"))

		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.Commit()
			if err == gorm.ErrRecordNotFound {
				return errors.NewValidationError("Invalid country name")
			} else {
				return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}
		}

		// Give country ID to talent.
		talent.CountryID = &countryID.ID
	}

	// Get state ID.
	if talentExcel.StateName != nil {

		// Create bucket for state ID.
		stateID := tal.IDModel{}

		// Get state ID from database.
		err := service.Repository.GetRecordForTenant(uow, tenantID, &stateID,
			repository.Select("`id`"),
			repository.Filter("`name`=?", talentExcel.StateName),
			repository.Filter("`tenant_id`=?", tenantID),
			repository.Table("states"))

		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.Commit()
			if err == gorm.ErrRecordNotFound {
				return errors.NewValidationError("Invalid state name")
			} else {
				return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
			}
		}

		// Give state ID to talent.
		talent.StateID = &stateID.ID
	}

	// Check if academics exist or not.
	if talentExcel.Academics != nil && len(talentExcel.Academics) > 0 {

		// Create bucket for academics.
		academics := []*tal.Academic{}

		// Loop the academics of talent excel.
		for _, academicExcel := range talentExcel.Academics {

			// Create bucket for academic.
			academic := tal.Academic{}

			// Create bucket for degree ID.
			degreeID := tal.IDModel{}

			// Get degree ID from database.
			err := service.Repository.GetRecordForTenant(uow, tenantID, &degreeID,
				repository.Select("`id`"),
				repository.Filter("`name`=?", academicExcel.DegreeName),
				repository.Filter("`tenant_id`=?", tenantID),
				repository.Table("degrees"))

			if err != nil {
				log.NewLogger().Error(err.Error())
				uow.Commit()
				if err == gorm.ErrRecordNotFound {
					return errors.NewValidationError("Invalid degree name")
				} else {
					return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
				}
			}

			// Give degree ID to academic.
			academic.DegreeID = degreeID.ID

			// Create bucket for specialization ID.
			specializationID := tal.IDModel{}

			// Get specialization ID from database.
			err = service.Repository.GetRecordForTenant(uow, tenantID, &specializationID,
				repository.Select("`id`"),
				repository.Filter("`branch_name`=?", academicExcel.SpecializationName),
				repository.Filter("`degree_id`=?", degreeID.ID),
				repository.Filter("`tenant_id`=?", tenantID),
				repository.Table("specializations"))

			if err != nil {
				log.NewLogger().Error(err.Error())
				uow.Commit()
				if err == gorm.ErrRecordNotFound {
					return errors.NewValidationError("Invalid specialization name")
				} else {
					return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
				}
			}

			// Give specialization ID to academic.
			academic.SpecializationID = specializationID.ID

			// Give college name to academic.
			academic.College = academicExcel.CollegeName

			// Give percentage to academic.
			academic.Percentage = academicExcel.Percentage

			// Give year of passout to academic.
			academic.Passout = academicExcel.YearOfPassout

			// Push academic into academics.
			academics = append(academics, &academic)
		}

		// Give academics to talent.
		talent.Academics = academics
	}

	// Give first name to talent.
	talent.FirstName = talentExcel.FirstName

	// Give last name to talent.
	talent.LastName = talentExcel.LastName

	// Give email to talent.
	talent.Email = talentExcel.Email

	// Give contact to talent.
	talent.Contact = talentExcel.Contact

	// Give academic year to talent.
	talent.AcademicYear = talentExcel.AcademicYear

	// Give is active to talent.
	talent.IsActive = talentExcel.IsActive

	// Give is swabhav talent to talent.
	talent.IsSwabhavTalent = talentExcel.IsSwabhavTalent

	// Give address to talent.
	talent.Address.Address = talentExcel.Address

	// Give city to talent.
	talent.City = talentExcel.City

	// Give Pin code to talent.
	talent.PINCode = talentExcel.PINCode

	// Give tenant ID to talent.
	talent.TenantID = tenantID

	// Give created by to talent.
	talent.CreatedBy = credentialID

	uow.Commit()

	// Add talent individually.
	if _, err := service.AddTalent(&talent, false); err != nil {
		return err
	}

	return nil
}

// UpdateTalent updates one talent in database by specific talent id.
func (service *TalentService) UpdateTalent(talent *tal.Talent, uows ...*repository.UnitOfWork) error {
	// Give tenant id to academics and experiences.
	service.setTenantID(talent, talent.TenantID)

	// Get credential id from UpdatedBy field of talent(set in controller).
	credentialID := talent.UpdatedBy

	// Extract all foreign key IDs,assign to entityID field and make entity object nil.
	service.extractID(talent)

	// Validate tenant id.
	if err := service.doesTenantExist(talent.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, talent.TenantID); err != nil {
		return err
	}

	// Validate foreign keys.
	if err := service.doForeignKeysExist(talent); err != nil {
		return err
	}

	// Starting new transaction.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Check for same email conflict.
	exists, err := service.doesEmailExist(uow, talent)
	if err != nil {
		uow.RollBack()
		return err
	}
	if exists {
		uow.RollBack()
		return errors.NewValidationError("Email already exists")
	}

	// Set college id by college name.
	if err := service.setCollegeNameAndID(uow, talent); err != nil {
		uow.RollBack()
		return err
	}

	// Give 01 to all talent's experiences' fromDate and toDate field.
	if talent.Experiences != nil && len(talent.Experiences) != 0 {
		for i := 0; i < len(talent.Experiences); i++ {
			if len(talent.Experiences[i].FromDate) != 0 {
				util.AddDateToMonth(&(talent.Experiences[i].FromDate))
			}
			if talent.Experiences[i].ToDate != nil && len(*talent.Experiences[i].ToDate) != 0 {
				util.AddDateToMonth((talent.Experiences[i].ToDate))
			}
		}
	}

	// Create bucket for getting talent already present in database.
	tempTalent := tal.Talent{}

	// Get talent for getting created_by field of talent from database.
	if err := service.Repository.GetForTenant(uow, talent.TenantID, talent.ID, &tempTalent,
		repository.PreloadAssociations([]string{"Academics", "Experiences", "MastersAbroad", "MastersAbroad.Scores", "Technologies"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give password from temp talent to talent to be updated.
	talent.Password = tempTalent.Password

	// Give created_by id from temp talent to talent to be updated.
	talent.CreatedBy = tempTalent.CreatedBy

	// Give credential id to updated_by of talent.
	talent.UpdatedBy = credentialID

	// Give code to talent.
	talent.Code = tempTalent.Code

	// Give lifetime value to talent.
	talent.LifetimeValue = tempTalent.LifetimeValue

	// Update academics.
	err = service.updateAcademics(uow, talent.Academics, talent.TenantID, talent.UpdatedBy, talent.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make academics nil to avoid any inserts or updates in academics table.
	talent.Academics = nil

	// Update experiences.
	err = service.updateExperiences(uow, talent.Experiences, talent.TenantID, talent.UpdatedBy, talent.ID)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make experiences nil to avoid any inserts or updates in experiences table.
	talent.Experiences = nil

	// Update matsers abroad.
	err = service.updateMastersAbroad(uow, talent.MastersAbroad, tempTalent.MastersAbroad, talent.TenantID,
		talent.UpdatedBy, talent.ID, talent)
	if err != nil {
		uow.RollBack()
		return err
	}

	// Make matsers abroad nil to avoid any inserts or updates in matsers abroad table.
	talent.MastersAbroad = nil

	// Update talent associations.
	if err := service.updateTalentAssociation(uow, talent); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be updated", http.StatusInternalServerError)
	}

	// Make technologies nil so that it is not inserted again.
	talent.Technologies = nil

	// Update talent.
	if err := service.Repository.Save(uow, talent); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be updated", http.StatusInternalServerError)
	}

	// Update is_active field for talent credential.
	if err := service.updateTalentCredential(uow, talent); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// UpdateTalentsSalesperson updates multiple talents' salesperson id.
func (service *TalentService) UpdateTalentsSalesperson(talents *[]tal.TalentUpdate, salepersonID uuid.UUID,
	tenantID uuid.UUID, credentialID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, tenantID); err != nil {
		return err
	}

	// Validate salespersonid (login id).
	if err := service.doesSalespersonExist(salepersonID, tenantID); err != nil {
		return err
	}

	// Collect all talent ids in variable.
	var talentIDs []uuid.UUID
	for _, talent := range *talents {
		talentIDs = append(talentIDs, talent.TalentID)
	}

	// Validate all talents.
	for _, talentID := range talentIDs {
		if err := service.doesTalentExist(talentID, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Update sales_person_id field of all talents.
	if err := service.Repository.UpdateWithMap(uow, &tal.Talent{}, map[string]interface{}{
		"SalesPersonID": salepersonID,
		"UpdatedBy":     credentialID,
	},
		repository.Filter("id IN (?)", talentIDs)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Sales person could not be allocated to talents", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteTalent deletes one talent from database by specific id.
func (service *TalentService) DeleteTalent(talent *tal.Talent, uows ...*repository.UnitOfWork) error {
	// Get credential id from DeletedBy field of talent(set in controller).
	credentialID := talent.DeletedBy

	// Validate tenant id.
	if err := service.doesTenantExist(talent.TenantID); err != nil {
		return err
	}

	// Validate credential id.
	if err := service.doesCredentialExist(credentialID, talent.TenantID); err != nil {
		return err
	}

	// Starting new transaction.
	var uow *repository.UnitOfWork
	length := len(uows)
	if length == 0 {
		uow = repository.NewUnitOfWork(service.DB, false)
	} else {
		uow = uows[0]
	}

	// Create bucket for counting batches of talent.
	batchCount := 0

	// Get count of all batches for the talent.
	err := service.Repository.GetCount(uow, &tal.Talent{}, &batchCount,
		repository.Join("INNER JOIN batch_talents ON batch_talents.`talent_id` = talents.`id`"),
		repository.Filter("batch_talents.`talent_id`=?", talent.ID),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.Filter("batch_talents.`tenant_id`=? AND talents.`tenant_id`=?", talent.TenantID, talent.TenantID),
	)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Talent should not be deleted if talent is assigned to any batch.
	if batchCount > 0 {
		log.NewLogger().Error("Cannot delete talent as it is assigned to a batch")
		uow.RollBack()
		return errors.NewValidationError("Cannot delete talent as it is assigned to a batch")
	}

	// Get talent for updating deleted_by field of talent.
	if err := service.Repository.GetForTenant(uow, talent.TenantID, talent.ID, talent,
		repository.PreloadAssociations([]string{"Technologies", "MastersAbroad"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//  Delete talent association from database.
	if err := service.deleteTalentAssociation(uow, talent, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete talent", http.StatusInternalServerError)
	}

	// Make technologies and masters abroad field nil to avoid any updates or inserts.
	talent.Technologies = nil
	talent.MastersAbroad = nil

	// Update talent for updating deleted_by and deleted_at fields of talent.
	if err := service.Repository.UpdateWithMap(uow, talent, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	// Create bucket for credential.
	credential := general.Credential{}

	// Call delete credential service method.
	loginService := genService.NewCredentialService(service.DB, service.Repository)
	if err := loginService.DeleteCredential(&credential, talent.TenantID, credentialID, talent.ID, "`talent_id`=?", uow); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete talent", http.StatusInternalServerError)
	}

	// Commit only if no transaction has been passed.
	if length == 0 {
		uow.Commit()
	}

	return nil
}

// GetTalent gets one talent from database by specific id.
func (service *TalentService) GetTalent(talent *tal.DTO, tenantID uuid.UUID) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent.
	if err := service.Repository.GetForTenant(uow, tenantID, talent.ID, talent,
		repository.PreloadAssociations(TalentAssociationNames)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	talents := []tal.DTO{*talent}

	// Sort the child tables.
	service.sortTalentChildTables(&talents)

	uow.Commit()
	return nil
}

// GetEligibleTalent returns talent by specific talent id for getting fields for eligibility check.
func (service *TalentService) GetEligibleTalent(talent *tal.EligibleTalentDTO, tenantID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get talent.
	if err := service.Repository.GetForTenant(uow, tenantID, talent.ID, talent,
		repository.PreloadAssociations([]string{"Academics", "Experiences", "Technologies", "Experiences.Technologies", "Academics.Degree", "Experiences.Designation"})); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetTalents gets all talents from database.
func (service *TalentService) GetTalents(talents *[]tal.DTO, tenantID uuid.UUID, limit int, offset int, totalCount *int,
	totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Login filter***************************************************
	// Variables for role name and login id.

	// Get query params for role name and login id.
	// #niranjan changed
	roleName := queryParams.Get("roleName")
	isViewAllBatches := queryParams.Get("isViewAllBatches")
	loginID, err := uuid.FromString(queryParams.Get("loginID"))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
	}

	// Create query processors according to conditions.
	var queryProcesors []repository.QueryProcessor

	// If login is salesperson.
	if roleName == "SalesPerson" {
		// Validate salespersonid (login id).
		if err := service.doesSalespersonExist(loginID, tenantID); err != nil {
			return err
		}
		// Add salesperson filter.
		queryProcesors = append(queryProcesors, repository.Filter("talents.`sales_person_id`=?", loginID))
	}

	// If login is faculty.
	if roleName == "Faculty" {

		// Validate facultyid (login id).
		if err := service.doesFacultyExist(tenantID, loginID); err != nil {
			return err
		}

		// Add batch filter and joins.
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN batch_talents ON talents.`id` = batch_talents.`talent_id`"),
			repository.Join("INNER JOIN batches ON batch_talents.`batch_id` = batches.`id`"),
			// // #niranjan changed
			// repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id`"+
			// 	"AND `talents`.`tenant_id`= `batch_modules`.`tenant_id`"),
			// repository.Filter("batch_modules.`faculty_id` = ?", loginID),
			// // end #
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
			repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("talents.`id`"))

		// If isViewAllBatches then give faculty filter.
		if (isViewAllBatches == "0"){
			queryProcesors = append(queryProcesors,
				repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id`"+
				"AND `talents`.`tenant_id`= `batch_modules`.`tenant_id`"),
				repository.Filter("batch_modules.`faculty_id` = ?", loginID))
		}
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Add paginate and preload.
	queryProcesors = append(queryProcesors, repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount),
		repository.Filter("talents.`tenant_id`=?", tenantID))

	// Get talents from database.
	if err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcesors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// If talents is not present then do not get extra fields.
	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************
	// Get lifetime value only for admin and salesperson.
	if roleName == "Admin" || roleName == "SalesPerson" {
		var queryProcessorsForLifetiemValue []repository.QueryProcessor
		queryProcessorsForLifetiemValue = append(queryProcessorsForLifetiemValue,
			repository.Filter("talents.`deleted_at` IS NULL"),
			repository.Filter("talents.`tenant_id`=?", tenantID),
			repository.Table("talents"),
			repository.Select("sum(talents.lifetime_value) as total_lifetime_value"))

		// Add filter for salesperson login.
		if roleName == "SalesPerson" {
			queryProcessorsForLifetiemValue = append(queryProcessorsForLifetiemValue,
				repository.Filter("talents.`sales_person_id`=?", loginID))
		}

		// Get lifetime value from database.
		if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValue...); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetAllSearchTalents gets all searched talents from database.
func (service *TalentService) GetAllSearchTalents(talents *[]tal.DTO, talentSearch *tal.Search, tenantID uuid.UUID,
	limit int, offset int, totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {
	//********************************************Login filter***************************************************
	// Variables for role name and login id.
	roleName := queryParams.Get("roleName")
	isViewAllBatches := queryParams.Get("isViewAllBatches")
	loginID, err := uuid.FromString(queryParams.Get("loginID"))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If search all talents is true then skip filter by salesperson and faculty.
	if talentSearch.SearchAllTalents != nil && !(*talentSearch.SearchAllTalents) {

		// If login is salesperson.
		if roleName == "SalesPerson" {

			// Validate salespersonid (login id).
			if err := service.doesSalespersonExist(loginID, tenantID); err != nil {
				return err
			}

			// Add salesperson filter.
			queryProcessors = append(queryProcessors, repository.Filter("talents.`sales_person_id`=?", loginID))
		}

		// If login is faculty.
		if roleName == "faculty" {

			// Validate facultyid (login id).
			if err := service.doesFacultyExist(tenantID, loginID); err != nil {
				return err
			}

			// Add faculty filter and joins.
			queryProcessors = append(queryProcessors,
				repository.Join("INNER JOIN batch_talents ON talents.`id` = batch_talents.`talent_id`"),
				repository.Join("INNER JOIN batches ON batch_talents.`batch_id` = batches.`id`"),
				// // #niranjan changed
				// repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id`"+
				// 	"AND `talents`.`tenant_id`= `batch_modules`.`tenant_id`"),
				// repository.Filter("batch_modules.`faculty_id` = ?", loginID),
				// // end #
				repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
				repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
				repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
				repository.GroupBy("talents.`id`"))
		
			// If isViewAllBatches then give faculty filter.
			if (isViewAllBatches == "0"){
				queryProcessors = append(queryProcessors,
					repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id`"+
					"AND `talents`.`tenant_id`= `batch_modules`.`tenant_id`"),
					repository.Filter("batch_modules.`faculty_id` = ?", loginID))
			}
		}
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Add preload, pagination and search queries.
	queryProcessors = append(queryProcessors, service.talentAddSearchQueries(tenantID, talentSearch, roleName)...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations(TalentAssociationNames),
		repository.PaginateWithoutModel(limit, offset, totalCount),
		repository.Filter("talents.`tenant_id`=?", tenantID))

	// Get talents from database.
	if err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// If talents is not present then do not get courses.
	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************
	// Get lifetime value only for admin and salesperson.
	if roleName == "admin" || roleName == "SalesPerson" {
		// Query precessors for sub query.
		var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`deleted_at` IS NULL"),
			repository.Filter("talents.`tenant_id`=?", tenantID),
			repository.Table("talents"),
			repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"))
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery, service.talentAddSearchQueries(tenantID, talentSearch, roleName)...)

		// Add filter for salesperson login.
		if roleName == "SalesPerson" {
			queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
				repository.Filter("talents.`sales_person_id`=?", loginID))
		}

		// Create query expression for sub query.
		subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}

		// Query processors for query.
		var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
		queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
			repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

		// Get total lifetime value of all talents.
		if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return errors.NewValidationError("Record not found")
		}
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsByRequirementID gets all talent from database by company requirement.
func (service *TalentService) GetTalentsByRequirementID(talents *[]tal.DTO, requirementID uuid.UUID,
	tenantID uuid.UUID, limit int, offset int, totalCount *int) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate company requirement id.
	if err := service.doesCompanyRequirementExist(requirementID, tenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all talents by company requirement id
	if err := service.Repository.GetAllInOrder(uow, talents, "`first_name`",
		// repository.Filter("id IN(?)", talentIDs),
		repository.Join("INNER JOIN company_requirements_talents ON talents.`id` = company_requirements_talents.`talent_id`"),
		repository.Filter("company_requirements_talents.`requirement_id` = ?", requirementID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Range talents for getting courses, faculties and expected ctc.
	err := service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsForExcelDownload returns talents for excel download.
func (service *TalentService) GetTalentsForExcelDownload(talents *[]tal.ExcelDTO, tenantID uuid.UUID, limit int,
	offset int, totalCount *int) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create bucket for temo talents.
	temptalents := []tal.DTO{}

	// Get talents.
	if err := service.Repository.GetAllInOrder(uow, &temptalents, "`first_name`, `last_name`",
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.PreloadAssociations([]string{"Technologies", "State", "Country", "SalesPerson",
			"TalentSource", "MastersAbroad", "MastersAbroad.Scores", "MastersAbroad.Scores.Examination",
			"MastersAbroad.Degree", "MastersAbroad.Universities", "MastersAbroad.Countries"}),
		repository.Paginate(limit, offset, totalCount)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Create talents.
	if err := service.createTalentsForExcel(uow, talents, &temptalents, tenantID); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// GetSearchedTalentsForExcelDownload gets searched talents from database.
func (service *TalentService) GetSearchedTalentsForExcelDownload(talents *[]tal.ExcelDTO, talentSearch *tal.Search, tenantID uuid.UUID,
	limit int, offset int, totalCount *int) error {

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Create bucket for temo talents.
	temptalents := []tal.DTO{}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Add preload, pagination and search queries.
	queryProcessors = append(queryProcessors, service.talentAddSearchQueries(tenantID, talentSearch, "")...)
	queryProcessors = append(queryProcessors, repository.PreloadAssociations([]string{"Technologies", "State", "Country", "SalesPerson",
		"TalentSource", "MastersAbroad", "MastersAbroad.Scores", "MastersAbroad.Scores.Examination",
		"MastersAbroad.Degree", "MastersAbroad.Universities", "MastersAbroad.Countries"}),
		repository.PaginateWithoutModel(limit, offset, totalCount),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("talents.`deleted_at` IS NULL"))

	// Get talents from database.
	if err := service.Repository.GetAllInOrder(uow, &temptalents, "`first_name`, `last_name`", queryProcessors...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Create talents.
	if err := service.createTalentsForExcel(uow, talents, &temptalents, tenantID); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	uow.Commit()
	return nil
}

// createTalentsForExcel create talent with minimum details for excel download.
func (service *TalentService) createTalentsForExcel(uow *repository.UnitOfWork, talents *[]tal.ExcelDTO, tempTalents *[]tal.DTO, tenantID uuid.UUID) error {

	for i := range *tempTalents {

		tempTalent := tal.ExcelDTO{}

		// Give basic details to talents from temp talents.
		tempTalent.Code = (*tempTalents)[i].Code
		tempTalent.FirstName = (*tempTalents)[i].FirstName
		tempTalent.LastName = (*tempTalents)[i].LastName
		tempTalent.Email = (*tempTalents)[i].Email
		tempTalent.Mobile = (*tempTalents)[i].Contact
		tempTalent.AlternateEmail = (*tempTalents)[i].AlternateEmail
		tempTalent.AlternateMobile = (*tempTalents)[i].AlternateContact
		tempTalent.AcademicYear = (*tempTalents)[i].AcademicYear
		if (*tempTalents)[i].Country != nil {
			tempTalent.Country = &(*tempTalents)[i].Country.Name
		}
		if (*tempTalents)[i].State != nil {
			tempTalent.State = &(*tempTalents)[i].State.Name
		}
		tempTalent.City = (*tempTalents)[i].City
		tempTalent.Pincode = (*tempTalents)[i].PINCode
		tempTalent.TalentType = (*tempTalents)[i].TalentType
		tempTalent.PersonalityType = (*tempTalents)[i].PersonalityType
		if (*tempTalents)[i].SalesPerson != nil {
			salesPersonName := (*tempTalents)[i].SalesPerson.FirstName + " " + (*tempTalents)[i].SalesPerson.LastName
			tempTalent.SalesPerson = &(salesPersonName)
		}
		if (*tempTalents)[i].TalentSource != nil {
			tempTalent.Source = &(*tempTalents)[i].TalentSource.Description
		}
		tempTalent.IsSwabhavTalent = (*tempTalents)[i].IsSwabhavTalent
		tempTalent.IsActive = (*tempTalents)[i].IsActive
		tempTalent.FacebookURL = (*tempTalents)[i].FacebookURL
		tempTalent.InstagramURL = (*tempTalents)[i].InstagramURL
		tempTalent.GithubURL = (*tempTalents)[i].GithubURL
		tempTalent.LinkedInURL = (*tempTalents)[i].LinkedInURL
		if (*tempTalents)[i].Technologies != nil && len((*tempTalents)[i].Technologies) > 0 {
			tempTechnologyArray := []string{}
			for j := range (*tempTalents)[i].Technologies {
				tempTechnologyArray = append(tempTechnologyArray, (*tempTalents)[i].Technologies[j].Language)
				tempTechnologiesString := strings.Join(tempTechnologyArray[:], ", ")
				tempTalent.Technologies = &tempTechnologiesString
			}

		}
		tempTalent.LoyaltyPoints = (*tempTalents)[i].LoyaltyPoints
		tempTalent.Resume = (*tempTalents)[i].Resume
		tempTalent.LifetimeValue = (*tempTalents)[i].LifetimeValue
		tempTalent.TotalYearOfExp = (*tempTalents)[i].ExperienceInMonths

		// Get latest talent academic.
		talentAcademic := tal.AcademicDTO{}
		if err := service.Repository.GetRecordForTenant(uow, tenantID, &talentAcademic,
			repository.Filter("talent_academics.`tenant_id`=?", tenantID),
			repository.Filter("talent_academics.`deleted_at` IS NULL"),
			repository.Filter("talent_academics.`talent_id`=?", (*tempTalents)[i].ID),
			repository.OrderBy("`passout` DESC"),
			repository.PreloadAssociations([]string{"Degree", "Specialization"})); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewValidationError("Internal server error")
			}
		}

		// Give latest academic details to talent.
		tempTalent.Qualification = &(talentAcademic.Degree.Name)
		tempTalent.Specialization = &(talentAcademic.Specialization.BranchName)
		tempTalent.CollegeName = &(talentAcademic.College)
		tempTalent.YearOfPassout = &(talentAcademic.Passout)
		tempTalent.CGPA = &(talentAcademic.Percentage)

		// Get latest talent experience.
		talentExperience := tal.ExperienceDTO{}
		if err := service.Repository.GetRecordForTenant(uow, tenantID, &talentExperience,
			repository.Filter("talent_experiences.`tenant_id`=?", tenantID),
			repository.Filter("talent_experiences.`deleted_at` IS NULL"),
			repository.Filter("talent_experiences.`talent_id`=?", (*tempTalents)[i].ID),
			repository.OrderBy("`from_date` DESC"),
			repository.PreloadAssociations([]string{"Technologies", "Designation"})); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewValidationError("Internal server error")
			}
		}

		// Give latest experience details to talent.
		tempTalent.Designation = &(talentExperience.Designation.Position)
		tempTalent.CurrentCompany = &(talentExperience.Company)
		tempTalent.Package = talentExperience.Package
		tempTalent.FromYear = &(talentExperience.FromDate)
		tempTalent.ToYear = talentExperience.ToDate
		if talentExperience.Technologies != nil && len(talentExperience.Technologies) > 0 {
			tempTechnologyArray := []string{}
			for j := range talentExperience.Technologies {
				tempTechnologyArray = append(tempTechnologyArray, talentExperience.Technologies[j].Language)
				tempTechnologiesString := strings.Join(tempTechnologyArray[:], ", ")
				tempTalent.ExperienceTechnologies = &tempTechnologiesString
			}

		}

		// Get talent's courses.
		tempCourses := []crs.Course{}
		err := service.Repository.GetAll(uow, &tempCourses,
			repository.Join("LEFT JOIN batches ON courses.`id` = batches.`course_id`"),
			repository.Join("LEFT JOIN batch_talents ON batch_talents.`batch_id` = batches.`id`"),
			repository.Join("LEFT JOIN talents ON batch_talents.`talent_id` = talents.`id`"),
			repository.Filter("talents.`id`=?", (*tempTalents)[i].ID),
			repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("courses.`tenant_id`=? AND courses.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("courses.`id`"),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		if tempCourses != nil && len(tempCourses) > 0 {
			tempCourseArray := []string{}
			for j := range tempCourses {
				tempCourseArray = append(tempCourseArray, tempCourses[j].Name)
				tempCoursesString := strings.Join(tempCourseArray[:], ", ")
				tempTalent.Courses = &tempCoursesString
			}
		}

		// Get talent's faculties.
		tempFaculties := []list.Faculty{}
		err = service.Repository.GetAll(uow, &tempFaculties,
			repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`faculty_id` = `faculties`.`id`"),
			repository.Join("INNER JOIN `batches` ON `batches`.`id` = `batch_modules`.`batch_id`"),
			repository.Join("JOIN batch_talents on batch_talents.`batch_id` = batches.`id`"),
			repository.Filter("batch_talents.`talent_id`=?", (*tempTalents)[i].ID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_modules.`tenant_id`=? AND batch_modules.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("faculties.`id`"),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		if tempFaculties != nil && len(tempFaculties) > 0 {
			tempFacultyArray := []string{}
			for j := range tempFaculties {
				tempFacultyArray = append(tempFacultyArray, (tempFaculties[j].FirstName + " " + tempFaculties[j].LastName))
				tempFacultiesString := strings.Join(tempFacultyArray[:], ", ")
				tempTalent.Faculties = &tempFacultiesString
			}
		}

		// Get talent's expected ctc.
		exepectedCTC := &tal.ExpectedCTCLatest{}
		if err := service.Repository.Scan(uow, exepectedCTC,
			repository.Filter("talents.`deleted_at` IS NULL AND talent_call_records.`deleted_at` IS NULL AND talents.`tenant_id`=? AND talent_call_records.`tenant_id`=? AND talents.`id`=?",
				tenantID, tenantID, (*tempTalents)[i].ID),
			repository.Filter("talent_call_records.`expected_ctc` IS NOT NULL"),
			repository.Table("talents"),
			repository.Join("JOIN talent_call_records on talent_call_records.`talent_id` = talents.`id`"),
			repository.Select("expected_ctc"),
			repository.GroupBy("talents.`id`"),
			repository.OrderBy("talent_call_records.`date_time`")); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewValidationError("Internal server error")
			}
		}

		if exepectedCTC != nil && exepectedCTC.ExpectedCTC > 0 {
			tempTalent.ExpectedCTC = &exepectedCTC.ExpectedCTC
		}

		*talents = append(*talents, tempTalent)
	}

	return nil
}

// GetTalentsForBatch returns all batches from database.
func (service *TalentService) GetTalentsForBatch(talents *[]tal.DTO, tenantID, batchID uuid.UUID,
	limit, offset int, totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult) error {
	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate batchID.
	err := service.doesBatchExist(tenantID, batchID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllInOrder(uow, talents, "`first_name`",
		repository.Join("LEFT JOIN batch_talents ON talents.`id` = batch_talents.`talent_id`"),
		repository.Filter("batch_talents.`batch_id`=? AND batch_talents.`is_active` = ?", batchID, true),
		repository.Filter("batch_talents.`tenant_id`=? AND talents.`tenant_id`=?", tenantID, tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Join("LEFT JOIN batch_talents ON talents.`id` = batch_talents.`talent_id`"),
		repository.Filter("batch_talents.`batch_id`=?", batchID),
		repository.Filter("batch_talents.`tenant_id`=?", tenantID),
		repository.Filter("batch_talents.`deleted_at` IS NULL"))

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsByWaitingList returns all talents by waiting list.
func (service *TalentService) GetTalentsByWaitingList(talents *[]tal.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************ID filter***************************************************
	// Variables for all possible IDs.
	companyBranchID := uuid.Nil
	requirementID := uuid.Nil
	courseID := uuid.Nil
	batchID := uuid.Nil
	technologyID := uuid.Nil

	// Get query params for all IDs.
	// Company branch ID.
	IDArray, _ := queryParams["companyBranchID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		companyBranchID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Company requirement ID.
	IDArray, _ = queryParams["requirementID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		requirementID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Course ID.
	IDArray, _ = queryParams["courseID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		courseID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Batch ID.
	IDArray, _ = queryParams["batchID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		batchID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Technology ID.
	IDArray, _ = queryParams["technologyID"]
	if IDArray != nil && len(IDArray) != 0 {
		var err error
		technologyID, err = uuid.FromString(IDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Create query processors according to conditions.
	var queryProcesors []repository.QueryProcessor

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate company branch ID if present.
	if companyBranchID != uuid.Nil {
		if err := service.doesCompanyBranchExist(companyBranchID, tenantID); err != nil {
			return err
		}

		// Add company branch filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`company_branch_id`=?", companyBranchID))
	}

	// Validate company requirement ID if present.
	if requirementID != uuid.Nil {
		if err := service.doesCompanyRequirementExist(requirementID, tenantID); err != nil {
			return err
		}

		// Add company requirement filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`company_requirement_id`=?", requirementID))
	}

	// Validate course ID if present.
	if courseID != uuid.Nil {
		if err := service.doesCourseExist(courseID, tenantID); err != nil {
			return err
		}

		// Add course filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`course_id`=?", courseID))
	}

	// Validate batch ID if present.
	if batchID != uuid.Nil {
		if err := service.doesBatchExist(tenantID, batchID); err != nil {
			return err
		}

		// Add batch filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`batch_id`=?", batchID))
	}

	// Validate technology ID if present.
	if technologyID != uuid.Nil {
		if err := service.doesTechnologyExist(technologyID, tenantID); err != nil {
			return err
		}

		// Add technology filter.
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Join("LEFT JOIN company_requirements on waiting_list.`company_requirement_id` = company_requirements.`id`"),
			repository.Join("LEFT JOIN company_requirements_technologies on company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
			repository.Join("LEFT JOIN courses on waiting_list.`course_id` = courses.`id`"),
			repository.Join("LEFT JOIN courses_technologies on courses.`id` = courses_technologies.`course_id`"),
			repository.Filter("company_requirements_technologies.`technology_id`=? OR courses_technologies.`technology_id`=?", technologyID, technologyID))
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcesors = append(queryProcesors,
		repository.Filter("talents.`id` IS NOT NULL"),
		repository.Filter("waiting_list.`deleted_at` IS NULL"),
		repository.Filter("waiting_list.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.GroupBy("talents.`id`"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcesors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Table("talents"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Filter("waiting_list.`tenant_id`=?", tenantID),
		repository.Filter("waiting_list.`deleted_at` IS NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// Filter for company branch ID.
	if companyBranchID != uuid.Nil {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("LEFT JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`company_branch_id`=?", companyBranchID))
	}

	// Filter for company requirement ID.
	if requirementID != uuid.Nil {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("LEFT JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`company_requirement_id`=?", companyBranchID))
	}

	// Filter for course ID.
	if courseID != uuid.Nil {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("LEFT JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`course_id`=?", courseID))
	}

	// Filter for batch ID.
	if batchID != uuid.Nil {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("LEFT JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`batch_id`=?", batchID))
	}

	// Filter for technology ID.
	if technologyID != uuid.Nil {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Join("LEFT JOIN company_requirements on waiting_list.`company_requirement_id` = company_requirements.`id`"),
			repository.Join("LEFT JOIN company_requirements_technologies on company_requirements.`id` = company_requirements_technologies.`requirement_id`"),
			repository.Join("LEFT JOIN courses on waiting_list.`course_id` = courses.`id`"),
			repository.Join("LEFT JOIN courses_technologies on courses.`id` = courses_technologies.`course_id`"),
			repository.Filter("company_requirements_technologies.`technology_id`=? OR courses_technologies.`technology_id`=?", technologyID, technologyID))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsByCampusDrive returns all talents by campus drive id.
func (service *TalentService) GetTalentsByCampusDrive(talents *[]tal.DTO, tenantID, campusDriveID uuid.UUID,
	limit, offset int, totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Appeared filter***************************************************
	// Variables for role name and login id.
	hasAppeared := ""

	// Get query params for has appeared.
	hasAppearedArray, _ := queryParams["hasAppeared"]
	if hasAppearedArray != nil && len(hasAppearedArray) != 0 {
		hasAppeared = hasAppearedArray[0]
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If has appeared is present then add has appeared filter.
	if hasAppeared == "1" {
		queryProcessors = append(queryProcessors,
			repository.Filter("campus_talent_registrations.`has_attempted`=?", hasAppeared))
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate campus drive id.
	if err := service.doesCampusDriveExist(campusDriveID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessors = append(queryProcessors, repository.Join("JOIN campus_talent_registrations ON talents.`id` = campus_talent_registrations.`talent_id`"),
		repository.Filter("talents.`id` IS NOT NULL"),
		repository.Filter("campus_talent_registrations.`campus_drive_id`=?", campusDriveID),
		repository.Filter("campus_talent_registrations.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("campus_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.GroupBy("talents.`id`"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Join("LEFT JOIN campus_talent_registrations ON talents.`id` = campus_talent_registrations.`talent_id`"),
		repository.Filter("campus_talent_registrations.`campus_drive_id`=?", campusDriveID),
		repository.Filter("campus_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("campus_talent_registrations.`deleted_at` IS NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// If has appeared is present then add has appeared filter.
	if hasAppeared == "1" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("campus_talent_registrations.`has_attempted`=?", hasAppeared))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsBySeminar returns all talents by seminar id.
func (service *TalentService) GetTalentsBySeminar(talents *[]tal.DTO, tenantID, seminarID uuid.UUID,
	limit, offset int, totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Visited filter***************************************************
	// Variables for role name and login id.
	hasVisited := ""

	// Get query params for has visited.
	hasVisitedArray, _ := queryParams["hasVisited"]
	if hasVisitedArray != nil && len(hasVisitedArray) != 0 {
		hasVisited = hasVisitedArray[0]
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If has visited is present then add has visited filter.
	if hasVisited == "1" {
		queryProcessors = append(queryProcessors,
			repository.Filter("seminar_talent_registrations.`has_visited`=?", hasVisited))
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate seminar id.
	if err := service.doesSeminarExist(seminarID, tenantID); err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessors = append(queryProcessors, repository.Join("JOIN seminar_talent_registrations ON talents.`id` = seminar_talent_registrations.`talent_id`"),
		repository.Filter("talents.`id` IS NOT NULL"),
		repository.Filter("seminar_talent_registrations.`seminar_id`=?", seminarID),
		repository.Filter("seminar_talent_registrations.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("seminar_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.GroupBy("talents.`id`"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Join("LEFT JOIN seminar_talent_registrations ON talents.`id` = seminar_talent_registrations.`talent_id`"),
		repository.Filter("seminar_talent_registrations.`seminar_id`=?", seminarID),
		repository.Filter("seminar_talent_registrations.`tenant_id`=?", tenantID),
		repository.Filter("seminar_talent_registrations.`deleted_at` IS NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// If has visited is present then add has visited filter.
	if hasVisited == "1" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("seminar_talent_registrations.`has_visited`=?", hasVisited))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsForProSummaryReport returns all talents for professional summary report.
func (service *TalentService) GetTalentsForProSummaryReport(talents *[]tal.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Appeared filter***************************************************
	// Variables for company name and category.
	companyName := ""
	category := ""
	isCompany := ""

	// Get query params for has appeared.
	companyNameArray, _ := queryParams["companyName"]
	if companyNameArray != nil && len(companyNameArray) != 0 {
		companyName = companyNameArray[0]
	}
	categoryArray, _ := queryParams["category"]
	if categoryArray != nil && len(categoryArray) != 0 {
		category = categoryArray[0]
	}
	isCompanyArray, _ := queryParams["isCompany"]
	if isCompanyArray != nil && len(isCompanyArray) != 0 {
		isCompany = isCompanyArray[0]
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If category is first then find talents with 12-24 months of experience.
	if category == "first" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>=12 AND `experience_in_months`<=24"))
	}

	// If category is second then find talents with 24-60 months of experience.
	if category == "second" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>24 AND `experience_in_months`<=60"))
	}

	// If category is third then find talents with 60-84 months of experience.
	if category == "third" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>60 AND `experience_in_months`<=84"))
	}

	// If category is fourth then find talents with above 84 months of experience.
	if category == "fourth" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months` > 84"))
	}

	// If category is fifth then find talents with above 12 months of experience.
	if category == "fifth" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months` > 12"))
	}

	// If category is fifth then find talents with above 0 months of experience.
	if category == "fifth" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months` > 0"))
	}

	// If company name is present then give its filter.
	if isCompany == "true" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`company`=?", companyName))
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate company name.
	if companyName != "none" {
		if err := service.doesCompanyNameExist(companyName, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessors = append(queryProcessors,
		repository.Join("JOIN talent_experiences on talents.`id` = talent_experiences.`talent_id`"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talent_experiences.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`to_date` IS NULL"),
		repository.Filter("`from_date` IS NOT NULL"),
		repository.GroupBy("talents.`id`"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	// Get all talents form database.
	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	// If talents are not present then return.
	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Join("JOIN talent_experiences on talents.`id` = talent_experiences.`talent_id`"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("talent_experiences.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("to_date IS NULL"),
		repository.Filter("from_date IS NOT NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// If category is first then find talents with 12-24 months of experience.
	if category == "first" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>=12 AND `experience_in_months`<=24"))
	}

	// If category is second then find talents with 24-60 months of experience.
	if category == "second" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>24 AND `experience_in_months`<=60"))
	}

	// If category is third then find talents with 60-84 months of experience.
	if category == "third" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>60 AND `experience_in_months`<=84"))
	}

	// If category is fourth then find talents with above 84 months of experience.
	if category == "fourth" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months` > 84"))
	}

	// If category is fifth then find talents with above 12 months of experience.
	if category == "fifth" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months` > 12"))
	}

	// If company name is present then give its filter.
	if isCompany == "true" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`company`=?", companyName))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsForProSummaryReportTechnologyCount returns all talents for professional summary report by technology talent count.
func (service *TalentService) GetTalentsForProSummaryReportTechnologyCount(talents *[]tal.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Appeared filter***************************************************
	// Variables for company name and category.
	companyName := ""
	category := ""
	technologyID := uuid.Nil
	isCompany := ""

	// Get query params for has appeared.
	companyNameArray, _ := queryParams["companyName"]
	if companyNameArray != nil && len(companyNameArray) != 0 {
		companyName = companyNameArray[0]
	}
	categoryArray, _ := queryParams["category"]
	if categoryArray != nil && len(categoryArray) != 0 {
		category = categoryArray[0]
	}
	isCompanyArray, _ := queryParams["isCompany"]
	if isCompanyArray != nil && len(isCompanyArray) != 0 {
		isCompany = isCompanyArray[0]
	}
	technologyIDArray, _ := queryParams["technologyID"]
	if technologyIDArray != nil && len(technologyIDArray) != 0 {
		var err error
		technologyID, err = uuid.FromString(technologyIDArray[0])
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Unable to get talents", http.StatusInternalServerError)
		}
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	// If category is first then find talents with 12-24 months of experience.
	if category == "first" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>=12 AND `experience_in_months`<=24"))
	}

	// If category is second then find talents with 24-60 months of experience.
	if category == "second" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>24 AND `experience_in_months`<=60"))
	}

	// If category is third then find talents with 60-84 months of experience.
	if category == "third" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months`>60 AND `experience_in_months`<=84"))
	}

	// If category is fourth then find talents with above 84 months of experience.
	if category == "fourth" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months` > 84"))
	}

	// If category is fifth then find talents with above 12 months of experience.
	if category == "fifth" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`experience_in_months` > 12"))
	}

	// If company name is present then give its filter.
	if isCompany == "true" {
		queryProcessors = append(queryProcessors,
			repository.Filter("`company`=?", companyName))
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate technology id.
	if err := service.doesTechnologyExist(technologyID, tenantID); err != nil {
		return err
	}

	// Validate company name.
	if companyName != "none" {
		if err := service.doesCompanyNameExist(companyName, tenantID); err != nil {
			return err
		}
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	queryProcessors = append(queryProcessors,
		repository.Join("JOIN talent_experiences on talents.`id` = talent_experiences.`talent_id`"),
		repository.Join("JOIN talent_experiences_technologies on talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("`technology_id`=?", technologyID),
		repository.Filter("talent_experiences.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("`to_date` IS NULL"),
		repository.Filter("`from_date` IS NOT NULL"),
		repository.GroupBy("talents.`id`"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	// Get talents from database.
	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	// If talents is not present then retuen.
	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Join("JOIN talent_experiences on talents.`id` = talent_experiences.`talent_id`"),
		repository.Join("JOIN talent_experiences_technologies on talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.Filter("`technology_id`=?", technologyID),
		repository.Filter("talent_experiences.`tenant_id`=?", tenantID),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Filter("to_date IS NULL"),
		repository.Filter("from_date IS NOT NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// If company name is present then give its filter.
	if isCompany == "true" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`company`=?", companyName))
	}

	// If category is first then find talents with 12-24 months of experience.
	if category == "first" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>=12 AND `experience_in_months`<=24"))
	}

	// If category is second then find talents with 24-60 months of experience.
	if category == "second" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>24 AND `experience_in_months`<=60"))
	}

	// If category is third then find talents with 60-84 months of experience.
	if category == "third" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months`>60 AND `experience_in_months`<=84"))
	}

	// If category is fourth then find talents with above 84 months of experience.
	if category == "fourth" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months` > 84"))
	}

	// If category is fifth then find talents with above 12 months of experience.
	if category == "fifth" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("`experience_in_months` > 12"))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsForFresherSummaryReport returns all talents for fresher summary report.
func (service *TalentService) GetTalentsForFresherSummaryReport(talents *[]tal.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	//********************************************Appeared filter***************************************************
	// Variables for academicYear, talentType and technology.
	// var academicYear string
	var talentType string
	var isExperienced string
	var isLookingForJob string
	var technology string

	var techLanguage = []string{
		"Advance Java", "Dotnet", "Java", "Machine Learning", "Cloud", "Golang",
	}

	talentTypeArray := queryParams["talentType"]
	if len(talentTypeArray) != 0 {
		talentType = queryParams.Get("talentType")
	}

	isExperiencedArray := queryParams["isExperienced"]
	if len(isExperiencedArray) != 0 {
		isExperienced = queryParams.Get("isExperienced")
	}

	isLookingForJobArray := queryParams["isExperienced"]
	if len(isLookingForJobArray) != 0 {
		isLookingForJob = queryParams.Get("isLookingForJob")
	}

	technologyArray := queryParams["fresherTechnology"]
	if len(technologyArray) != 0 {
		technology = queryParams.Get("fresherTechnology")
		// var err error
		// technology, err = uuid.FromString(queryParams.Get("technology"))
		// if err != nil {
		//   return err
		// }
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor
	var subQueryProcessor []repository.QueryProcessor

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	//  && isExperienced != "all"
	if isExperiencedArray != nil {
		queryProcessors = append(queryProcessors, repository.Filter("talents.`is_experience` = ?", isExperienced))
	}

	if isLookingForJob == "1" {

		subQueryProcessor = append(subQueryProcessor, repository.Table("talent_call_records"),
			repository.Select("max(talent_call_records.`date_time`)"), repository.GroupBy("talent_call_records.`talent_id`"))

		// Create query expression for sub query one.
		subQueryOne, err := service.Repository.SubQuery(uow, talent.CallRecord{}, subQueryProcessor...)
		if err != nil {
			uow.RollBack()
			return err
		}

		queryProcessors = append(queryProcessors, repository.Table("talents"),
			repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
				"talents.`tenant_id` = talent_call_records.`tenant_id`"),
			repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id` = outcomes.`id` AND "+
				"outcomes.`tenant_id` = talent_call_records.`tenant_id`"), repository.Filter("talents.`tenant_id` = ?", tenantID),
			repository.Filter("`talents`.`is_experience` = '1'"), repository.Filter("talents.`deleted_at` IS NULL"),
			repository.Filter("talent_call_records.`deleted_at` IS NULL AND outcomes.`deleted_at` IS NULL"),
			repository.Filter("talent_call_records.`date_time` IN ? AND outcomes.`outcome` = 'Job switch'", subQueryOne))
	}

	// If talentType is Outstanding then find talents with talent_type >= 5.
	if talentType == "Outstanding" {
		queryProcessors = append(queryProcessors, repository.Filter("talents.`talent_type` >= 5"))
	}

	// If talentType is Excellent then find talents with talent_type BETWEEN 3 and 4.
	if talentType == "Excellent" {
		queryProcessors = append(queryProcessors, repository.Filter("talents.`talent_type` BETWEEN 3 AND 4"))
	}

	// If talentType is Average then find talents with talent_type <= 2.
	if talentType == "Average" {
		queryProcessors = append(queryProcessors, repository.Filter("talents.`talent_type` <= 2"))
	}

	// If talentType is Unranked then find talents with talent_type is null.
	if talentType == "Unranked" {
		queryProcessors = append(queryProcessors, repository.Filter("talents.`talent_type` IS NULL"))
	}

	// Get query params for has appeared.
	academicYear := queryParams["academicYear"]
	if len(academicYear) != 0 {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`academic_year` = ?", queryParams.Get("academicYear")))
	}

	if technologyArray != nil && technology != "Other" {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Filter("talents_technologies.`technology_id` = ? ", technology))
	}

	if technology == "Other" {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Join("INNER JOIN technologies ON technologies.`id` = talents_technologies.`technology_id`"),
			repository.Filter("technologies.`language` NOT IN (?)", techLanguage))
	}

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	queryProcessors = append(queryProcessors,
		// repository.GroupBy("talents.id"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.GroupBy("talents.`id`"),
	)

	// If talentType is Outstanding then find talents with talent_type >= 5.
	if talentType == "Outstanding" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`talent_type` >= 5"))
	}

	// If talentType is Excellent then find talents with talent_type BETWEEN 3 and 4.
	if talentType == "Excellent" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`talent_type` BETWEEN 3 AND 4"))
	}

	// If talentType is Average then find talents with talent_type <= 2.
	if talentType == "Average" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`talent_type` <= 2"))
	}

	// If talentType is Unranked then find talents with talent_type is null.
	if talentType == "Unranked" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`talent_type` IS NULL"))
	}

	// Get query params for has appeared.
	if len(academicYear) != 0 {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`academic_year` = ?", queryParams.Get("academicYear")))
	}

	if technologyArray != nil && technology != "Other" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Filter("talents_technologies.`technology_id` = ? ", technology))
	}

	if technology == "Other" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"),
			repository.Join("INNER JOIN technologies ON technologies.`id` = talents_technologies.`technology_id`"),
			repository.Filter("technologies.`language` NOT IN (?)", techLanguage))
	}

	if isLookingForJob == "1" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id` AND "+
				"talents.`tenant_id` = talent_call_records.`tenant_id`"),
			repository.Join("INNER JOIN outcomes ON talent_call_records.`outcome_id` = outcomes.`id` AND "+
				"outcomes.`tenant_id` = talent_call_records.`tenant_id`"),
			repository.Filter("outcomes.`outcome` = 'Job switch'"),
			repository.OrderBy("talent_call_records.`date_time` DESC"))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentsForPackageSummaryReport returns all talents for package summary report.
func (service *TalentService) GetTalentsForPackageSummaryReport(talents *[]tal.DTO, tenantID uuid.UUID, limit, offset int,
	totalCount *int, totalLifetimeValue *tal.TotalLifetimeValueResult, queryParams url.Values) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	//********************************************Appeared filter***************************************************
	// Variables for academicYear, talentType and technology.
	var packageType string
	var technology string
	var packageExperience string

	var techLanguage = []string{
		"Advance Java", "Dotnet", "Java", "Machine Learning", "Cloud", "Golang",
	}

	var experience = []string{
		"0-3", "3-6", "6-8", "8-10", "10+",
	}

	packageTypeArray := queryParams["packageType"]
	if len(packageTypeArray) != 0 {
		packageType = queryParams.Get("packageType")
	}

	technologyArray := queryParams["packageTechnology"]
	if len(technologyArray) != 0 {
		technology = queryParams.Get("packageTechnology")
	}

	packageExperienceArray := queryParams["packageExperience"]
	if len(packageExperienceArray) != 0 {
		packageExperience = queryParams.Get("packageExperience")
	}

	// Create query processors according to conditions.
	var queryProcessors []repository.QueryProcessor

	queryProcessors = append(queryProcessors,
		repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` "+
			" AND talents.`tenant_id` = talent_experiences.`tenant_id`"), repository.Filter("talent_experiences.`to_date` IS NULL"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Filter("talent_experiences.`package` IS NOT NULL AND talent_experiences.`from_date` IS NOT NULL"))

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// experience -> 0 to 3 years
	if packageExperience == experience[0] {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 1, 36))
	}

	// experience -> 3 to 6 years
	if packageExperience == experience[1] {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 37, 72))
	}

	// experience -> 6 to 8 years
	if packageExperience == experience[2] {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 73, 96))
	}

	// experience -> 8 to 10 years
	if packageExperience == experience[3] {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 97, 120))
	}

	// experience -> 10+ years
	if packageExperience == experience[4] {
		queryProcessors = append(queryProcessors,
			repository.Filter("talents.`experience_in_months` > ?", 120))
	}

	// If packageType is LessThanThree then find talents with package >= 300000.
	if packageType == "LessThanThree" {
		queryProcessors = append(queryProcessors,
			repository.Filter("talent_experiences.`package` <= ?", 300000))
	}

	// If packageType is ThreeToFive then find talents with talent_type BETWEEN 300001 and 500000.
	if packageType == "ThreeToFive" {
		queryProcessors = append(queryProcessors,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 300001, 500000))
	}

	// If packageType is FiveToTen then find talents with package BETWEEN 500001 AND 1000000.
	if packageType == "FiveToTen" {
		queryProcessors = append(queryProcessors,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 500001, 1000000))
	}

	// If packageType is TenToFifteen then find talents with package BETWEEN 1000001 AND 1500000.
	if packageType == "TenToFifteen" {
		queryProcessors = append(queryProcessors,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 1000001, 1500000))
	}

	// If packageType is GreaterThanFifteen then find talents with package > 1500000.
	if packageType == "GreaterThanFifteen" {
		queryProcessors = append(queryProcessors,
			repository.Filter("talent_experiences.`package` > ?", 1500000))
	}

	if technologyArray != nil && technology != "Other" {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN talent_experiences_technologies ON "+
				"talent_experiences.`id` = talent_experiences_technologies.`experience_id`"),
			repository.Filter("talent_experiences_technologies.`technology_id` = ? ", technology))
	}

	if technology == "Other" {
		queryProcessors = append(queryProcessors,
			repository.Join("INNER JOIN talent_experiences_technologies ON "+
				"talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
			repository.Join("INNER JOIN technologies ON talent_experiences_technologies.`technology_id` = technologies.`id`"),
			repository.Filter("technologies.`id` NOT IN (?)", techLanguage))
	}

	queryProcessors = append(queryProcessors,
		// repository.GroupBy("talents.id"),
		repository.PreloadAssociations(TalentAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	err := service.Repository.GetAllInOrder(uow, talents, "`first_name`, `last_name`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
	}

	if talents == nil || len(*talents) == 0 {
		return nil
	}

	// Range talents for getting courses, faculties and expected ctc.
	err = service.getValuesForTalent(uow, talents, tenantID)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	//**************************************************Lifetime value***********************************************************

	// Query precessors for sub query.
	var queryProcessorsForLifetiemValueSubQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
		repository.Table("talents"),
		repository.Select("sum(talents.lifetime_value) as total_lifetime_value_all_talents"),
		repository.Filter("talents.`tenant_id`=?", tenantID),
		repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` "+
			" AND talents.`tenant_id` = talent_experiences.`tenant_id`"), repository.Filter("talent_experiences.`to_date` IS NULL"),
		repository.Filter("talent_experiences.`deleted_at` IS NULL"), repository.Filter("talents.`tenant_id` = ?", tenantID),
		repository.Filter("talent_experiences.`package` IS NOT NULL AND talent_experiences.`from_date` IS NOT NULL"),
		repository.GroupBy("talents.`id`"),
	)

	// experience -> 0 to 3 years
	if packageExperience == experience[0] {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 1, 36))
	}

	// experience -> 3 to 6 years
	if packageExperience == experience[1] {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 37, 72))
	}

	// experience -> 6 to 8 years
	if packageExperience == experience[2] {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 73, 96))
	}

	// experience -> 8 to 10 years
	if packageExperience == experience[3] {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`experience_in_months` BETWEEN ? AND ?", 97, 120))
	}

	// experience -> 10+ years
	if packageExperience == experience[4] {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talents.`experience_in_months` > ?", 120))
	}

	// If packageType is LessThanThree then find talents with package >= 300000.
	if packageType == "LessThanThree" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talent_experiences.`package` <= ?", 300000))
	}

	// If packageType is ThreeToFive then find talents with talent_type BETWEEN 300001 and 500000.
	if packageType == "ThreeToFive" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 300001, 500000))
	}

	// If packageType is FiveToTen then find talents with package BETWEEN 500001 AND 1000000.
	if packageType == "FiveToTen" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 500001, 1000000))
	}

	// If packageType is TenToFifteen then find talents with package BETWEEN 1000001 AND 1500000.
	if packageType == "TenToFifteen" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talent_experiences.`package` BETWEEN ? AND ?", 1000001, 1500000))
	}

	// If packageType is GreaterThanFifteen then find talents with package > 1500000.
	if packageType == "GreaterThanFifteen" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Filter("talent_experiences.`package` > ?", 1500000))
	}

	if technologyArray != nil && technology != "Other" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("INNER JOIN talent_experiences_technologies ON "+
				"talent_experiences.`id` = talent_experiences_technologies.`experience_id`"),
			repository.Filter("talent_experiences_technologies.`technology_id` = ? ", technology))
	}

	if technology == "Other" {
		queryProcessorsForLifetiemValueSubQuery = append(queryProcessorsForLifetiemValueSubQuery,
			repository.Join("INNER JOIN talent_experiences_technologies ON "+
				"talent_experiences_technologies.`experience_id` = talent_experiences.`id`"),
			repository.Join("INNER JOIN technologies ON talent_experiences_technologies.`technology_id` = technologies.`id`"),
			repository.Filter("technologies.`id` NOT IN (?)", techLanguage))
	}

	// Create query expression for sub query.
	subQuery, err := service.Repository.SubQuery(uow, talents, queryProcessorsForLifetiemValueSubQuery...)
	if err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Query processors for query.
	var queryProcessorsForLifetiemValueQuery []repository.QueryProcessor
	queryProcessorsForLifetiemValueQuery = append(queryProcessorsForLifetiemValueQuery,
		repository.RawQuery("select sum(total_lifetime_value_all_talents) as total_lifetime_value from ? as sub_query", subQuery))

	// Get total lifetime value of all talents.
	if err := service.Repository.Scan(uow, totalLifetimeValue, queryProcessorsForLifetiemValueQuery...); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Sort the child tables.
	service.sortTalentChildTables(talents)

	uow.Commit()
	return nil
}

// GetTalentBatches gets all batches for a specific talent.
func (service *TalentService) GetTalentBatches(batches *[]bat.BatchDTO, tenantID, talentID uuid.UUID) error {
	// Check if tenant exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Check if talent exist.
	err = service.doesTalentExist(talentID, tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get all batches for the talent.
	err = service.Repository.GetAllInOrder(uow, batches, "status DESC",
		repository.Join("INNER JOIN batch_talents ON batch_talents.`batch_id` = batches.`id`"),
		repository.Filter("batch_talents.`talent_id`=?", talentID),
		repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
		repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
		repository.PreloadAssociations(BatchAssociations),
		repository.PreloadWithCustomCondition(repository.Preload{
			Schema: "Timing",
			Queryprocessors: []repository.QueryProcessor{
				repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
				repository.Filter("days.`deleted_at` IS NULL AND batch_timing.`tenant_id` = ?", tenantID),
				repository.OrderBy("days.`order`"),
			}}))
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}
	// err = service.Repository.GetAllInOrderForTenant(uow, tenantID, batches,
	// 	repository.Join("INNER JOIN batch_talents ON batch_talents.`batch_id` = batches.`id`"),
	// 	repository.Filter("batch_talents.`talent_id`=?", talentID),
	// 	repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
	// 	repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
	// 	repository.PreloadWithCustomCondition(repository.Preload{
	// 		Schema: "Timing",
	// 		Queryprocessors: []repository.QueryProcessor{
	// 			repository.Join("INNER JOIN days ON days.`id` = batch_timing.`day_id` AND days.`tenant_id` = batch_timing.`tenant_id`"),
	// 			repository.Filter("days.`deleted_at` IS NULL AND batch_timing.`tenant_id` = ?", tenantID),
	// 			repository.OrderBy("days.`order`"),
	// 		}}), repository.PreloadAssociations(BatchAssociations),
	// 	repository.OrderBy("batches.`created_at` DESC, batches.`is_active` DESC"))
	// if err != nil {
	// 	log.NewLogger().Error(err.Error())
	// 	uow.RollBack()
	// 	return err
	// }

	for index := range *batches {
		// err = service.Repository.GetAllInOrder(uow, &(*batches)[index].Timing, "days.order",
		// 	repository.PreloadAssociations([]string{"Day"}),
		// 	repository.Join("INNER JOIN days ON days.`id`=batch_timing.`day_id`"),
		// 	repository.Filter("batch_timing.`batch_id`=?", (*batches)[index].ID),
		// 	repository.Filter("batch_timing.`tenant_id`=? AND batch_timing.`deleted_at` IS NULL", tenantID),
		// 	repository.Filter("days.`tenant_id`=? AND days.`deleted_at` IS NULL", tenantID))
		// if err != nil {
		// 	uow.RollBack()
		// 	return err
		// }
		
		err = service.Repository.GetAll(uow, &(*batches)[index].Faculty,
		repository.Join("INNER JOIN batch_modules ON faculties.id = batch_modules.faculty_id"+
			" AND faculties.tenant_id = batch_modules.tenant_id"),
		repository.Join("INNER JOIN batches ON batches.id = batch_modules.batch_id"),
		repository.Filter("batch_modules.batch_id=?", &(*batches)[index].ID),
		repository.GroupBy("faculties.`id`"))
		if err != nil {
			uow.RollBack()
			return err
		}

		err = service.getBatchSessionsCount(uow, &(*batches)[index].TotalSessionCount,
			&(*batches)[index].CompletedSessionCount, tenantID, (*batches)[index].ID)
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// updateTalentAssociation updates talent's associations.
func (service *TalentService) updateTalentAssociation(uow *repository.UnitOfWork, talent *tal.Talent) error {
	// Replace technologies of talent.
	if err := service.Repository.ReplaceAssociations(uow, talent, "Technologies",
		talent.Technologies); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	return nil
}

// updateAcademics will update the academics for specified talent.
func (service *TalentService) updateAcademics(uow *repository.UnitOfWork, academics []*tal.Academic,
	tenantID, credentialID, talentID uuid.UUID) error {

	// If previous academics is present and current academics is not present.
	if academics == nil {
		err := service.Repository.UpdateWithMap(uow, tal.Academic{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`talent_id`=?", talentID))
		if err != nil {
			return err
		}
	}

	// Create temp academics to get presvious academics of talent.
	tempAcademics := []tal.Academic{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempAcademics,
		repository.Filter("`talent_id`=?", talentID))
	if err != nil {
		return err
	}

	// Make map to count occurences of academic id in previous and current academics.
	academicIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous academic ID.
	for _, tempAcademic := range tempAcademics {
		academicIDMap[tempAcademic.ID]++
	}

	// Count the number of occurrence of current academic ID to know total count of occurrenec of each ID.
	for _, academic := range academics {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(academic.ID) {
			academicIDMap[academic.ID]++
		} else {
			// If ID is nil create new academic entry i table.
			academic.CreatedBy = credentialID
			academic.TenantID = tenantID
			academic.TalentID = talentID
			err = service.Repository.Add(uow, &academic)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current academics) then update academic.
		if academicIDMap[academic.ID] > 1 {
			academic.UpdatedBy = credentialID
			err = service.Repository.Update(uow, &academic)
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after updating academic.
			academicIDMap[academic.ID] = 0
		}
	}

	// If the number of occurrence is one, the academic was presnt in previous academics only, delete it.
	for _, tempAcademic := range tempAcademics {
		if academicIDMap[tempAcademic.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, tal.Academic{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempAcademic.ID))
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after deleting academic.
			academicIDMap[tempAcademic.ID] = 0
		}
	}
	return nil
}

// updateExperiences will update the experiences for specified talent.
func (service *TalentService) updateExperiences(uow *repository.UnitOfWork, experiences []*tal.Experience,
	tenantID, credentialID, talentID uuid.UUID) error {

	// If previous experiences is present and current experiences is not present.
	if experiences == nil {
		err := service.Repository.UpdateWithMap(uow, tal.Experience{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		}, repository.Filter("`talent_id`=?", talentID))
		if err != nil {
			return err
		}
	}

	// Create temp experiences to get presvious experiences of talent.
	tempExperiences := []tal.Experience{}

	err := service.Repository.GetAllForTenant(uow, tenantID, &tempExperiences,
		repository.Filter("`talent_id`=?", talentID))
	if err != nil {
		return err
	}

	// Make map to count occurences of experience id in previous and current experiences.
	experienceIDMap := make(map[uuid.UUID]uint)

	// Count the number of occurence of previous experience ID.
	for _, tempExperience := range tempExperiences {
		experienceIDMap[tempExperience.ID]++
	}

	// Count the number of occurrence of current experience ID to know total count of occurrenec of each ID.
	for _, experience := range experiences {

		// If ID is valid then push its occurence in ID map.
		if util.IsUUIDValid(experience.ID) {
			experienceIDMap[experience.ID]++
		} else {
			// If ID is nil create new experience entry i table.
			experience.CreatedBy = credentialID
			experience.TenantID = tenantID
			experience.TalentID = talentID
			err = service.Repository.Add(uow, &experience)
			if err != nil {
				return err
			}
		}

		// If number of occurrence is more than one (present in previous and current experiences) then update experience.
		if experienceIDMap[experience.ID] > 1 {
			// Give created_by field of previous experiences to current experiences.
			for _, tempExperience := range tempExperiences {
				if tempExperience.ID == experience.ID {
					experience.CreatedBy = tempExperience.CreatedBy
				}
			}

			// Replace technologies of expereince.
			if err := service.Repository.ReplaceAssociations(uow, experience, "Technologies",
				experience.Technologies); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			// Make technologies nil to avoid insertion during update talent.
			experience.Technologies = nil

			// Gice updated_by field to current experiences.
			experience.UpdatedBy = credentialID
			experience.TalentID = talentID
			err = service.Repository.Save(uow, &experience)
			if err != nil {
				return err
			}

			// Make the number of occurrences 0 after updating experience.
			experienceIDMap[experience.ID] = 0
		}
	}

	// If the number of occurrence is one, the experience was presnt in previous experiences only, delete it.
	for _, tempExperience := range tempExperiences {
		if experienceIDMap[tempExperience.ID] == 1 {
			err = service.Repository.UpdateWithMap(uow, tal.Experience{}, map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			}, repository.Filter("`id` = ?", tempExperience.ID))
			if err != nil {
				return err
			}
			// Make the number of occurrences 0 after deleting experience.
			experienceIDMap[tempExperience.ID] = 0
		}
	}
	return nil
}

// updateMastersAbroad will update masters abroad for specified talent.
func (service *TalentService) updateMastersAbroad(uow *repository.UnitOfWork, mastersAbroad, tempMastersAbroad *general.MastersAbroad,
	tenantID, credentialID, talentID uuid.UUID, talent *tal.Talent) error {
	// If previous masters abroad does not exist and current masters abroad exists.
	if tempMastersAbroad == nil && mastersAbroad != nil {
		// Give masters abroad and its score arrays created by field.
		mastersAbroad.CreatedBy = credentialID
		for i := 0; i < len(mastersAbroad.Scores); i++ {
			mastersAbroad.Scores[i].CreatedBy = credentialID
		}
		// Add masters abroad.
		err := service.Repository.Add(uow, &mastersAbroad)
		if err != nil {
			return err
		}
	}

	// Do checks with previous data only if previous masters abroad and current masters abroad exists.
	if tempMastersAbroad != nil && mastersAbroad != nil {
		// Create temp scores to get presvious scores of masters abroad.
		tempScores := []general.Score{}

		err := service.Repository.GetAllForTenant(uow, tenantID, &tempScores,
			repository.Filter("`masters_abroad_id`=?", tempMastersAbroad.ID))
		if err != nil {
			return err
		}

		// Make map to count occurences of scores id in previous and current scores.
		scoreIDMap := make(map[uuid.UUID]uint)

		// Count the number of occurence of previous score ID.
		for _, tempScore := range tempScores {
			scoreIDMap[tempScore.ID]++
		}

		// Count the number of occurrence of current score ID to know total count of occurrence of each ID.
		for _, score := range mastersAbroad.Scores {

			// If ID is valid then push its occurence in ID map.
			if util.IsUUIDValid(score.ID) {
				scoreIDMap[score.ID]++
			} else {
				// If ID is nil create new score entry in table.
				score.CreatedBy = credentialID
				score.TenantID = tenantID
				err = service.Repository.Add(uow, &score)
				if err != nil {
					return err
				}
			}

			// If number of occurrence is more than one (present in previous and current scores) then update score.
			if scoreIDMap[score.ID] > 1 {

				// Gice updated_by field to current scores.
				score.UpdatedBy = credentialID
				err = service.Repository.Update(uow, &score)
				if err != nil {
					return err
				}

				// Make the number of occurrences 0 after updating experience.
				scoreIDMap[score.ID] = 0
			}
		}

		// If the number of occurrence is one, the score was presnt in previous scores only, delete it.
		for _, tempScore := range tempScores {
			if scoreIDMap[tempScore.ID] == 1 {
				err = service.Repository.UpdateWithMap(uow, general.Score{}, map[string]interface{}{
					"DeletedBy": credentialID,
					"DeletedAt": time.Now(),
				}, repository.Filter("`id` = ?", tempScore.ID))
				if err != nil {
					return err
				}
				// Make the number of occurrences 0 after deleting scores.
				scoreIDMap[tempScore.ID] = 0
			}
		}

		// Replace countries and universities of masters abroad of talent.
		// Replace countries.
		if err := service.Repository.ReplaceAssociations(uow, mastersAbroad, "Countries",
			mastersAbroad.Countries); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}
		// Replace universities.
		if err := service.Repository.ReplaceAssociations(uow, mastersAbroad, "Universities",
			mastersAbroad.Universities); err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Make universities empty to avoid unnecessary updates.
		mastersAbroad.Universities = nil

		// Make countries empty to avoid unnecessary updates.
		mastersAbroad.Countries = nil

		// Make scores empty to avoid unnecessary updates.
		mastersAbroad.Scores = nil

		// Give updated_by field to masters abroad.
		mastersAbroad.UpdatedBy = credentialID

		// Update masters abroad.
		err = service.Repository.Update(uow, &mastersAbroad)
		if err != nil {
			return err
		}
	}

	// Check if current talent's masters abroad exists or not, if doesnt exist then make it nil
	// and delete masters abroad from database.
	if mastersAbroad == nil && tempMastersAbroad != nil {
		// Make talent's master abroad id as nil.
		talent.IsMastersAbroad = false

		// Update score(s) for updating deleted_by field of score(s).
		err := service.Repository.UpdateWithMap(uow, &general.Score{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		},
			repository.Filter("`masters_abroad_id`=? AND `tenant_id`=?", tempMastersAbroad.ID, talent.TenantID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}

		// Delete masters abroad form database.
		// Update masters abroad for updating deleted_by field of masters abroad.
		if err := service.Repository.UpdateWithMap(uow, &general.MastersAbroad{}, map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
			repository.Filter("`talent_id`=? AND `tenant_id`=?", talentID, talent.TenantID)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Masters Abroad details could not be deleted", http.StatusInternalServerError)
		}
	}
	return nil
}

// deleteTalentAssociation deletes talent's associations.
func (service *TalentService) deleteTalentAssociation(uow *repository.UnitOfWork, talent *tal.Talent, credentialID uuid.UUID) error {

	//**********************************************EXPERIENCES***************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.Experience{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//**********************************************ACADEMICS******************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.Academic{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//**********************************************MASTERS ABROAD*******************************************************
	if talent.IsMastersAbroad {
		// Update score(s) for updating deleted_by field of score(s).
		err := service.Repository.UpdateWithMap(uow, &general.Score{}, map[string]interface{}{
			"DeletedAt": time.Now(),
			"DeletedBy": credentialID,
		},
			repository.Filter("`masters_abroad_id`=? AND `tenant_id`=?", talent.MastersAbroad.ID, talent.TenantID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}

		// Make countries and universities field nil to avoid any updates or inserts.
		talent.MastersAbroad.Countries = nil
		talent.MastersAbroad.Universities = nil

		// Update masters abroad for updating deleted_by and deleted_at fields of masters abroad.
		if err := service.Repository.UpdateWithMap(uow, &general.MastersAbroad{},
			map[string]interface{}{
				"DeletedBy": credentialID,
				"DeletedAt": time.Now(),
			},
			repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
		}
	}

	//***********************************************CALL RECORDS************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.CallRecord{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//*******************************************INTERVIEWS AND INTERVIEW SCHEDULES********************************
	interviewSchedules := []tal.InterviewSchedule{}
	// Get interview schedules from database.
	if err := service.Repository.GetAllForTenant(uow, talent.TenantID, &interviewSchedules,
		repository.Filter("`talent_id`=?", talent.ID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Collect all interview schedules ids in a variable.
	var interviewScheduleIDs []uuid.UUID
	for _, interviewSchedule := range interviewSchedules {
		interviewScheduleIDs = append(interviewScheduleIDs, interviewSchedule.ID)
	}

	// Deleting interviews.
	if err := service.Repository.UpdateWithMap(uow, &tal.Interview{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`schedule_id` IN (?) AND `tenant_id`=?", interviewScheduleIDs, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	// Deleting interview schedules.
	if err := service.Repository.UpdateWithMap(uow, &tal.InterviewSchedule{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//*********************************************LIFETIME VALUE*****************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.LifetimeValue{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//************************************************CAREER PLANS************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.CareerPlan{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//************************************************NEXT ACTION************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.NextAction{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//************************************************WAITING LIST************************************************
	if err := service.Repository.UpdateWithMap(uow, &tal.WaitingList{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//************************************************CAMPUS TALENT REGISTRATION************************************************
	if err := service.Repository.UpdateWithMap(uow, &college.CampusTalentRegistration{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	//************************************************SEMINAR TALENT REGISTRATION************************************************
	if err := service.Repository.UpdateWithMap(uow, &college.SeminarTalentRegistration{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`talent_id`=? AND `tenant_id`=?", talent.ID, talent.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Talent could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// talentAddSearchQueries adds all search queries by comparing with the talent data recieved from request body.
func (service *TalentService) talentAddSearchQueries(tenantID uuid.UUID, talent *tal.Search, roleName string) []repository.QueryProcessor {

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	fmt.Println("===================================SEARCH TALENT===================================================")
	fmt.Println("talent ->", talent)
	fmt.Println("=================================================================================================")

	if talent.FirstName != nil {
		util.AddToSlice("talents.`first_name`", "LIKE ?", "AND", "%"+*talent.FirstName+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if talent.LastName != nil {
		util.AddToSlice("talents.`last_name`", "LIKE ?", "AND", "%"+*talent.LastName+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if talent.Email != nil {
		util.AddToSlice("talents.`email`", "LIKE ?", "AND", "%"+*talent.Email+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if talent.IsExperience != nil {
		util.AddToSlice("talents.`is_experience`", "= ?", "AND", *talent.IsExperience,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.AcademicYears != nil && len(talent.AcademicYears) != 0 {
		util.AddToSlice("talents.`academic_year`", "IN(?)", "AND", talent.AcademicYears,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.IsActive != nil {
		util.AddToSlice("talents.`is_active`", "= ?", "AND", *talent.IsActive, &columnNames,
			&conditions, &operators, &values)
	}
	if talent.PersonalityType != nil {
		util.AddToSlice("talents.`personality_type`", "= ?", "AND", *talent.PersonalityType,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.TalentType != nil {
		util.AddToSlice("talents.`talent_type`", ">= ?", "AND", *talent.TalentType,
			&columnNames, &conditions, &operators, &values)
	}
	// remove uuid.nil element from the IDs & add an IS NULL query.
	if len(talent.SalesPersonIDs) > 0 {
		util.AddToSlice("talents.`sales_person_id`", "IN(?)", "AND", talent.SalesPersonIDs,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.LifetimeValue != nil {
		util.AddToSlice("talents.`lifetime_value`", ">=?", "AND", talent.LifetimeValue,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.IsMastersAbroad != nil {
		util.AddToSlice("talents.`is_masters_abroad`", "=?", "AND", talent.IsMastersAbroad,
			&columnNames, &conditions, &operators, &values)
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN masters_abroad ON talents.`id` = masters_abroad.`talent_id`"),
			repository.Filter("masters_abroad.`deleted_at` IS NULL"),
			repository.Filter("masters_abroad.`tenant_id`=?", tenantID))

		if talent.YearOfMS != nil {
			util.AddToSlice("masters_abroad.`year_of_ms`", "=?", "AND", talent.YearOfMS,
				&columnNames, &conditions, &operators, &values)
		} else {
			queryProcesors = append(queryProcesors, repository.Filter("masters_abroad.`year_of_ms` IS NULL"))
		}
	}
	if talent.IsSwabhavTalent != nil {
		util.AddToSlice("talents.`is_swabhav_talent`", "=?", "AND", talent.IsSwabhavTalent,
			&columnNames, &conditions, &operators, &values)
	}
	if talent.City != nil {
		util.AddToSlice("talents.`city`", "LIKE ?", "AND", "%"+*talent.City+"%",
			&columnNames, &conditions, &operators, &values)
	}
	if talent.CountryID != nil {
		util.AddToSlice("talents.`country_id`", "=?", "AND", talent.CountryID,
			&columnNames, &conditions, &operators, &values)
	}

	if talent.NextActionTypeID != nil {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_next_actions ON talents.`id` = talent_next_actions.`talent_id`"),
			repository.Filter("talent_next_actions.`deleted_at` IS NULL"),
			repository.Filter("talent_next_actions.`tenant_id`=?", tenantID))

		util.AddToSlice("talent_next_actions.`action_type_id`", "=?", "AND", talent.NextActionTypeID,
			&columnNames, &conditions, &operators, &values)
	}

	if talent.CourseID != nil || talent.BatchID != nil || talent.FacultyID != nil {

		if roleName != "Faculty" {
			queryProcesors = append(queryProcesors,
				repository.Join("INNER JOIN batch_talents ON talents.`id` = batch_talents.`talent_id`"),
				repository.Join("INNER JOIN batches ON batches.`id` = batch_talents.`batch_id`"),
				repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`batch_id` = `batches`.`id`"),
				repository.Filter("batch_modules.`deleted_at` IS NULL"),
				repository.Filter("batch_talents.`deleted_at` IS NULL"),
				repository.Filter("batch_modules.`tenant_id`=?", tenantID),
				repository.Filter("batch_talents.`tenant_id`=?", tenantID),
				repository.Filter("batches.`deleted_at` IS NULL"),
				repository.Filter("batches.`tenant_id`=?", tenantID))
			if talent.FacultyID != nil {
				fmt.Println("====================IN Faculty===============================")
				util.AddToSlice("batch_modules.`faculty_id`", "=?", "AND",
					talent.FacultyID, &columnNames, &conditions, &operators, &values)
			}
		}

		if talent.BatchID != nil {
			fmt.Println("====================IN Batch===============================")
			util.AddToSlice("batches.`id`", "=?", "AND",
				talent.BatchID, &columnNames, &conditions, &operators, &values)
		}
		if talent.CourseID != nil {
			fmt.Println("====================IN Course===============================")
			util.AddToSlice("batches.`course_id`", "=?", "AND",
				talent.CourseID, &columnNames, &conditions, &operators, &values)
		}
	}

	// If college or qualifications is present then join talent_academics table.
	if talent.College != nil || len(talent.Qualifications) != 0 {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_academics ON talents.`id` = talent_academics.`talent_id`"),
			repository.Filter("talent_academics.`deleted_at` IS NULL"),
			repository.Filter("talent_academics.`tenant_id`=?", tenantID))

		if talent.College != nil {
			util.AddToSlice("talent_academics.`college`", "LIKE ?", "AND", "%"+*talent.College+"%",
				&columnNames, &conditions, &operators, &values)
		}

		if len(talent.Qualifications) > 0 {
			util.AddToSlice("talent_academics.`degree_id`", "IN(?)", "AND", talent.Qualifications,
				&columnNames, &conditions, &operators, &values)
		}
	}

	if talent.CallRecordOutcomeID != nil || talent.CallRecordPurposeID != nil {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_call_records ON talents.`id` = talent_call_records.`talent_id`"),
			repository.Filter("talent_call_records.`deleted_at` IS NULL"),
			repository.Filter("talent_call_records.`tenant_id`=?", tenantID))

		if talent.CallRecordPurposeID != nil {
			util.AddToSlice("talent_call_records.`purpose_id`", "=?", "AND", talent.CallRecordPurposeID,
				&columnNames, &conditions, &operators, &values)
		}
		if talent.CallRecordOutcomeID != nil {
			util.AddToSlice("talent_call_records.`outcome_id`", "=?", "AND", talent.CallRecordOutcomeID,
				&columnNames, &conditions, &operators, &values)
		}
	}

	// If experince technologies is present then join talent_experiences table, talent_experiences_technologies table.
	if len(talent.ExperienceTechnologies) != 0 || len(talent.Designations) != 0 || talent.CompanyName != nil {

		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id`"),
			repository.Join("LEFT JOIN talent_experiences_technologies ON talent_experiences.`id` = talent_experiences_technologies.`experience_id`"),
			repository.Filter("talent_experiences.`deleted_at` IS NULL"),
			repository.Filter("talent_experiences.`tenant_id`=?", tenantID))
		if len(talent.ExperienceTechnologies) > 0 {
			fmt.Println("====================IN EXPERIENCE tech===============================")
			util.AddToSlice("talent_experiences_technologies.`technology_id`", "IN(?)", "AND",
				talent.ExperienceTechnologies, &columnNames, &conditions, &operators, &values)
		}
		if len(talent.Designations) > 0 {
			fmt.Println("====================IN EXPERIENCE Designation===============================")
			util.AddToSlice("talent_experiences.`designation_id`", "IN(?)", "AND",
				talent.Designations, &columnNames, &conditions, &operators, &values)
		}
		if talent.CompanyName != nil {
			fmt.Println("====================IN EXPERIENCE Company name===============================")
			util.AddToSlice("talent_experiences.`company`", "LIKE ?", "AND", "%"+*talent.CompanyName+"%",
				&columnNames, &conditions, &operators, &values)
		}
	}

	// If technologies is present then join talents_technologies and technolgies table.
	if len(talent.Technologies) != 0 {
		fmt.Println("====================IN REGULAR TECH===============================")
		queryProcesors = append(queryProcesors, repository.Join("INNER JOIN talents_technologies ON talents.`id` = talents_technologies.`talent_id`"))

		if len(talent.Technologies) > 0 {
			util.AddToSlice("talents_technologies.`technology_id`", "IN(?)", "AND", talent.Technologies, &columnNames, &conditions, &operators, &values)
		}
	}

	// If one or more waiting list fields are present then join waiting_list table.
	if talent.WaitingFor != nil || talent.WaitingForCompanyBranchID != nil || talent.WaitingForRequirementID != nil ||
		talent.WaitingForCourseID != nil || talent.WaitingForBatchID != nil || talent.WaitingForIsActive != nil ||
		talent.WaitingForFromDate != nil || talent.WaitingForToDate != nil {
		queryProcesors = append(queryProcesors,
			repository.Join("INNER JOIN waiting_list ON talents.`id` = waiting_list.`talent_id`"),
			repository.Filter("waiting_list.`deleted_at` IS NULL"),
			repository.Filter("waiting_list.`tenant_id`=?", tenantID))

		// Waiting for is 'Company' search.
		if talent.WaitingFor != nil && *talent.WaitingFor == "Company" {
			queryProcesors = append(queryProcesors,
				repository.Filter("waiting_list.`company_branch_id` IS NOT NULL OR waiting_list.`company_requirement_id` IS NOT NULL"))
		}

		// Waiting for is 'Course' search.
		if talent.WaitingFor != nil && *talent.WaitingFor == "Course" {
			queryProcesors = append(queryProcesors,
				repository.Filter("waiting_list.`course_id` IS NOT NULL OR waiting_list.`batch_id` IS NOT NULL"))
		}

		// Company branch id search.
		if talent.WaitingForCompanyBranchID != nil {
			util.AddToSlice("waiting_list.`company_branch_id`", "=?", "AND", talent.WaitingForCompanyBranchID, &columnNames, &conditions, &operators, &values)
		}

		// Company requirement id search.
		if talent.WaitingForRequirementID != nil {
			util.AddToSlice("waiting_list.`company_requirement_id`", "=?", "AND", talent.WaitingForRequirementID, &columnNames, &conditions, &operators, &values)
		}

		// Course id search.
		if talent.WaitingForCourseID != nil {
			util.AddToSlice("waiting_list.`course_id`", "=?", "AND", talent.WaitingForCourseID, &columnNames, &conditions, &operators, &values)
		}

		// Batch id search.
		if talent.WaitingForBatchID != nil {
			util.AddToSlice("waiting_list.`batch_id`", "=?", "AND", talent.WaitingForBatchID, &columnNames, &conditions, &operators, &values)
		}

		// Is active search.
		if talent.WaitingForIsActive != nil {
			util.AddToSlice("waiting_list.`is_active`", "=?", "AND", talent.WaitingForIsActive, &columnNames, &conditions, &operators, &values)
		}

		// Applied for from date.
		if talent.WaitingForFromDate != nil {
			util.AddToSlice("waiting_list.`created_at`", ">= ?", "AND", talent.WaitingForFromDate, &columnNames, &conditions, &operators, &values)
		}

		// Applied for to date.
		if talent.WaitingForToDate != nil {
			util.AddToSlice("waiting_list.`created_at`", "<= ?", "AND", talent.WaitingForToDate, &columnNames, &conditions, &operators, &values)
		}
	}

	var havingQuery, minExp, maxExp string

	// If total experience is present then give having clause.
	if talent.MinimumExperience != nil || talent.MaximumExperience != nil {
		fmt.Println("====================IN TOTAL EXP===========================")
		havingQuery = "SUM(TIMESTAMPDIFF(month, talent_experiences.from_date, " +
			"IFNULL(talent_experiences.to_date, CURDATE()))) "
		// If experineces technologies is not present then join talent experineces table.
		if len(talent.ExperienceTechnologies) == 0 {
			queryProcesors = append(queryProcesors,
				repository.Join("INNER JOIN talent_experiences ON talents.`id` = talent_experiences.`talent_id` AND "+
					"talents.`tenant_id`=talent_experiences.`tenant_id` AND talent_experiences.`deleted_at` IS NULL"))
		}
		if talent.MinimumExperience != nil && talent.MaximumExperience != nil {
			minExp = strconv.Itoa(int(*talent.MinimumExperience * 12))
			maxExp = strconv.Itoa(int(*talent.MaximumExperience * 12))
			havingQuery = havingQuery + "> " + minExp + " AND " + havingQuery + "< " + maxExp
		}
		if talent.MinimumExperience != nil && talent.MaximumExperience == nil {
			minExp = strconv.Itoa(int(*talent.MinimumExperience * 12))
			havingQuery = havingQuery + "> " + minExp
		}
		if talent.MaximumExperience != nil && talent.MinimumExperience == nil {
			maxExp = strconv.Itoa(int(*talent.MaximumExperience * 12))
			havingQuery = havingQuery + "< " + maxExp
		}
	}
	// // #need to modify when deleted records search is given!
	// util.AddToSlice("talents.`deleted_at`", "IS NULL", "AND", nil,
	// 	&columnNames, &conditions, &operators, &values)
	queryProcesors = append(queryProcesors,
		repository.Table("talents"),
		repository.Filter("talents.`deleted_at` IS NULL"),
		repository.FilterWithOperator(columnNames, conditions, operators, values),
		repository.GroupBy("talents.`id`"),
		repository.Having(havingQuery))

	return queryProcesors
}

// =================================================NIRANJAN TEST===============================================================

// doesCredentialExist validates if credental exists or not in database.
func (service *TalentService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *TalentService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doForeignKeysExist validates all foreign keys of talent.
func (service *TalentService) doForeignKeysExist(talent *tal.Talent) error {
	uow := repository.NewUnitOfWork(service.DB, true)

	// Check parent talent present or not if referral id is not nil.
	if talent.ReferralID != nil { //referral id is present
		exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, tal.Talent{},
			repository.Filter("`id` = ?", talent.ReferralID))
		if err := util.HandleError("Invalid referral talent ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check parent source exists or not.
	if talent.SourceID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Source{},
			repository.Filter("`id` = ?", talent.SourceID))
		if err := util.HandleError("Invalid source ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check parent salesperson exists or not.
	if talent.SalesPersonID != nil {
		exists, err := repository.DoesRecordExist(uow.DB, general.User{},
			repository.Join("left join roles ON users.`role_id` = roles.`id`"),
			repository.Filter("users.`id`=? AND roles.`role_name`=? AND users.`tenant_id`=? AND roles.`tenant_id`=?",
				talent.SalesPersonID, "salesperson", talent.TenantID, talent.TenantID))
		if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if country exists or not.
	if talent.CountryID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Country{},
			repository.Filter("`id` = ?", talent.CountryID))
		if err := util.HandleError("Invalid country ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if state exists or not.
	if talent.StateID != nil {
		exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.State{},
			repository.Filter("`id` = ?", talent.StateID))
		if err := util.HandleError("Invalid state ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}
	}

	// Check if technologies exist or not.
	if talent.Technologies != nil && len(talent.Technologies) != 0 {
		if err := service.doTechnologiesExist(talent.Technologies, talent.TenantID); err != nil {
			return err
		}
	}

	// Check if experiences' technologies and designation exists or not.
	if talent.Experiences != nil && len(talent.Experiences) != 0 {
		for _, experience := range talent.Experiences {
			// Check if designation exists or not.
			exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Designation{},
				repository.Filter("`id` = ?", experience.DesignationID))
			if err := util.HandleError("Invalid designation ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
			if experience.Technologies != nil && len(experience.Technologies) != 0 {
				if err := service.doTechnologiesExist(experience.Technologies, talent.TenantID); err != nil {
					return err
				}
			}
		}
	}

	// Check if academic's degree and specialization exists or not.
	if talent.Academics != nil && len(talent.Academics) != 0 {
		for _, academic := range talent.Academics {
			// Check if degree exists or not.
			exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Degree{},
				repository.Filter("`id` = ?", academic.DegreeID))
			if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}

			// Check if specialization exists or not.
			exists, err = repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Specialization{},
				repository.Filter("`id`=? AND `degree_id`=?", academic.SpecializationID, academic.DegreeID))
			if err := util.HandleError("Invalid specialization ID", exists, err); err != nil {
				log.NewLogger().Error(err.Error())
				return err
			}
		}
	}

	// Check if masters abroad degree, universities and countries exists or not.
	if talent.MastersAbroad != nil {
		// Check if degree exists or not.
		exists, err := repository.DoesRecordExistForTenant(uow.DB, talent.TenantID, general.Degree{},
			repository.Filter("`id` = ?", talent.MastersAbroad.DegreeID))
		if err := util.HandleError("Invalid degree ID", exists, err); err != nil {
			log.NewLogger().Error(err.Error())
			return err
		}

		// Check if universities exist or not.
		var universityIDs []uuid.UUID
		for _, university := range talent.MastersAbroad.Universities {
			universityIDs = append(universityIDs, university.ID)
		}
		var count int = 0
		err = service.Repository.GetCountForTenant(uow, talent.TenantID, general.University{}, &count,
			repository.Filter("`id` IN (?)", universityIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(talent.MastersAbroad.Universities) {
			log.NewLogger().Error("University ID is invalid")
			return errors.NewValidationError("University ID is invalid")
		}

		// Check if countries exist or not.
		var countryIDs []uuid.UUID
		for _, country := range talent.MastersAbroad.Countries {
			countryIDs = append(countryIDs, country.ID)
		}
		err = service.Repository.GetCountForTenant(uow, talent.TenantID, general.Country{}, &count,
			repository.Filter("`id` IN (?)", countryIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(talent.MastersAbroad.Countries) {
			log.NewLogger().Error("Country ID is invalid")
			return errors.NewValidationError("Country ID is invalid")
		}

		// Check if score's examination exist or not.
		var examinationIDs []uuid.UUID
		for _, score := range talent.MastersAbroad.Scores {
			examinationIDs = append(examinationIDs, score.ExaminationID)
		}

		err = service.Repository.GetCountForTenant(uow, talent.TenantID, general.Examination{}, &count,
			repository.Filter("`id` IN (?)", examinationIDs))
		if err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
		}
		if count != len(examinationIDs) {
			log.NewLogger().Error("Examination ID is invalid")
			return errors.NewValidationError("Examination ID is invalid")
		}
	}
	return nil
}

// doTechnologiesExist validates if technolgy exists or not in database.
func (service *TalentService) doTechnologiesExist(technologies []*general.Technology, tenantID uuid.UUID) error {
	uow := repository.NewUnitOfWork(service.DB, true)
	// Keep all technology ids in one variable.
	var technologyIds []uuid.UUID
	for _, technology := range technologies {
		technologyIds = append(technologyIds, technology.ID)
	}
	// Get count for technologyIDs.
	var count int = 0
	err := service.Repository.GetCountForTenant(uow, tenantID, general.Technology{}, &count,
		repository.Filter("`id` IN (?)", technologyIds))
	if err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError(err.Error(), http.StatusInternalServerError)
	}
	if count != len(technologies) {
		log.NewLogger().Error("Technology ID is invalid")
		return errors.NewValidationError("Technology ID is invalid")
	}
	return nil
}

// doesSalespersonExist validates if salesperson exists or not in database.
func (service *TalentService) doesSalespersonExist(salespersonID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.User{},
		repository.Filter("`id` = ?", salespersonID))
	if err := util.HandleError("Invalid salesperson ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCompanyNameExist validates if company name exists or not in database.
func (service *TalentService) doesCompanyNameExist(companyName string, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Experience{},
		repository.Filter("`company` = ? AND `deleted_at` IS NULL", companyName))
	if err := util.HandleError("Company name does not exist for any talent experiences", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTechnologyExist validates if technology exists or not in database.
func (service *TalentService) doesTechnologyExist(technologyID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Technology{},
		repository.Filter("`id` = ?", technologyID))
	if err := util.HandleError("Invalid technology ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBatchExist validates if batch exists or not in database.
func (service *TalentService) doesBatchExist(tenantID uuid.UUID, batchID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, bat.Batch{},
		repository.Filter("`id` = ?", batchID))
	if err := util.HandleError("Batch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCourseExist validates if course exists or not in database.
func (service *TalentService) doesCourseExist(courseID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, crs.Course{},
		repository.Filter("`id` = ?", courseID))
	if err := util.HandleError("Course not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyRequirementExist validates if company requirement exists or not in database.
func (service *TalentService) doesCompanyRequirementExist(requirementID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Requirement{},
		repository.Filter("`id` = ?", requirementID))
	if err := util.HandleError("Company requirement not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCompanyBranchExist validates if company branch exists or not in database.
func (service *TalentService) doesCompanyBranchExist(companyBranchID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, company.Branch{},
		repository.Filter("`id` = ?", companyBranchID))
	if err := util.HandleError("Company branch not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCampusDriveExist validates if campus drive exists or not in database.
func (service *TalentService) doesCampusDriveExist(campsuDriveID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.CampusDrive{},
		repository.Filter("`id` = ?", campsuDriveID))
	if err := util.HandleError("Campus drive not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesSeminarExist validates if seminar exists or not in database.
func (service *TalentService) doesSeminarExist(seminarID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, college.Seminar{},
		repository.Filter("`id` = ?", seminarID))
	if err := util.HandleError("Seminar not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesTalentExist validates if talent exists or not in database.
func (service *TalentService) doesTalentExist(talentID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, tal.Talent{}, repository.Filter("`id` = ?", talentID))
	if err := util.HandleError("Talent not found", exists, err); err != nil {
		return err
	}
	return nil
}

//doesEmailExist check for same email conflict, if email exists return true.
func (service *TalentService) doesEmailExist(uow *repository.UnitOfWork, talent *tal.Talent) (bool, error) {
	var count int
	if err := service.Repository.GetCountForTenant(uow, talent.TenantID, &tal.Talent{}, &count,
		repository.Filter("`email`=? AND `id`!=?", talent.Email, talent.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		return false, errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if count != 0 { //email already present
		log.NewLogger().Error("Validate email:Email already exists")
		return true, nil
	}

	return false, nil
}

// doesFacultyExist validates faculty exists in database or not.
func (service *TalentService) doesFacultyExist(tenantID, facultyID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, fclt.Faculty{}, repository.Filter("`id` = ?", facultyID))
	if err := util.HandleError("Faculty not found", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// AddCredential adds credential to credential table according to talent details and generates password.
func (service *TalentService) AddCredential(uow *repository.UnitOfWork, talent *tal.Talent, credentialID uuid.UUID) error {
	// Create bucket for credential.
	credential := general.Credential{}

	// Create bucket for role.
	role := general.Role{}

	// Get role id by role name as 'talent'.
	if err := service.Repository.GetAllForTenant(uow, talent.TenantID, &role, repository.Filter("`role_name`=?", "talent")); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Record not found", http.StatusInternalServerError)
	}

	// Give talent details to credential.
	credential.FirstName = talent.FirstName
	credential.LastName = &talent.LastName
	credential.Email = talent.Email
	credential.Password = talent.Password
	credential.Contact = talent.Contact
	credential.RoleID = role.ID
	credential.TenantID = talent.TenantID
	credential.TalentID = &talent.ID
	credential.CreatedBy = credentialID

	// Add credential to database.
	loginService := genService.NewCredentialService(service.DB, service.Repository)
	if err := loginService.AddCredential(&credential, uow); err != nil {
		if err.Error() == "Credential already exists" || err.Error() == "Credential with same email and role exist" {
			errorString := "Email: " + credential.Email + " already exists"
			return errors.NewValidationError(errorString)
		}
		return err
	}
	return nil
}

// setTenantID gives tenant id to all academics and experiences of talent.
func (service *TalentService) setTenantID(talent *tal.Talent, tenantID uuid.UUID) {
	// If academics is present then give all the academics tenant id.
	if talent.Academics != nil && len(talent.Academics) != 0 {
		for _, academic := range talent.Academics {
			academic.TenantID = tenantID
		}
	}

	// If experiences is present then give all the experiences tenant id.
	if talent.Experiences != nil && len(talent.Experiences) != 0 {
		for _, experience := range talent.Experiences {
			experience.TenantID = tenantID
		}
	}

	// If masters abroad is present then give masters abroad and scores tenant id.
	if talent.MastersAbroad != nil {
		talent.MastersAbroad.TenantID = tenantID
		for _, score := range talent.MastersAbroad.Scores {
			score.TenantID = tenantID
		}
	}
}

// sortTalentChildTables sorts talent's academics and experineces.
func (service *TalentService) sortTalentChildTables(talents *[]tal.DTO) {
	if talents != nil && len(*talents) != 0 {
		for i := 0; i < len(*talents); i++ {

			// Sort academics by order od passout year in ascending order.
			academics := &(*talents)[i].Academics
			for j := 0; j < len(*academics); j++ {
				if (*academics)[j].Passout == 0 {
					return
				}
			}
			for j := 0; j < len(*academics); j++ {
				sort.Slice(*academics, func(p, q int) bool {
					return (*academics)[p].Passout < (*academics)[q].Passout
				})
			}

			// Sort experiences by order od from date in ascending order.
			experiences := &(*talents)[i].Experiences
			for j := 0; j < len(*experiences); j++ {
				if len((*experiences)[j].FromDate) == 0 {
					return
				}
			}
			for j := 0; j < len(*experiences); j++ {
				sort.Slice(*experiences, func(p, q int) bool {
					FromDateSmallInStr := (*experiences)[p].FromDate[:4]
					FromDateLargeInStr := (*experiences)[q].FromDate[:4]
					FromDateSmallInInt, _ := strconv.Atoi(FromDateSmallInStr)
					FromDateLargeInInt, _ := strconv.Atoi(FromDateLargeInStr)
					return FromDateSmallInInt < FromDateLargeInInt
				})
			}
		}
	}
}

// Extracts ID from object and removes data from the object.
// this is done so that the foreign key entity records are not updated in their respective tables
// when the college branch entity is being added or updated.
func (service *TalentService) extractID(talent *tal.Talent) error {
	// State field.
	if talent.State != nil {
		talent.StateID = &talent.State.ID
	}

	// Country field.
	if talent.Country != nil {
		talent.CountryID = &talent.Country.ID
	}

	return nil
}

// AddTalentAndCredential checks if talent exists in db using it's email & if it does it will check for
// that talent in credential table using talent's ID.
// If talent doesn't exist in the first place, it will add a talent & a credential record accoridingly.
func (service *TalentService) AddTalentAndCredential(uow *repository.UnitOfWork, talent *tal.Talent) error {

	credentialID := talent.CreatedBy

	// Check if email already exists.
	exists, err := service.doesEmailExist(uow, talent)
	if err != nil {
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if exists { // Email exists.
		// Get talent that exists.
		// Create talent bucket to get talent already present in DB.
		talentInDB := tal.Talent{}
		if err := service.Repository.GetRecordForTenant(uow, talent.TenantID, &talentInDB,
			repository.Filter("`email`=?", talent.Email)); err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		// If password is blank in database then generate password and give it to talent's password.
		if talentInDB.Password == "" || talentInDB.Password == talentInDB.FirstName {
			password := util.GeneratePassword()
			// Update talent with new password.
			if err := service.Repository.UpdateWithMap(uow, &talentInDB, map[string]interface{}{"Password": password}); err != nil {
				log.NewLogger().Error(err.Error())
				return errors.NewHTTPError("Talent could not be added", http.StatusInternalServerError)
			}
		}

		// Create credential for talent.
		if err := service.AddCredential(uow, &talentInDB, talentInDB.CreatedBy); err != nil {
			return err
		}

		// Return error to user for same email.
		return errors.NewValidationError("Email already exists")

	}
	// Email does not exist.
	// Generate password and give it to talent's password.
	talent.Password = util.GeneratePassword()

	// Add talent to database.
	if err := service.Repository.Add(uow, talent); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Talent could not be added", http.StatusInternalServerError)
	}

	// Create credential for talent.
	if err := service.AddCredential(uow, talent, credentialID); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	return nil
}

// updateTalentCredential updates specified talent record in credentials table.
func (service *TalentService) updateTalentCredential(uow *repository.UnitOfWork, talent *tal.Talent) error {
	err := service.Repository.UpdateWithMap(uow, &general.Credential{}, map[string]interface{}{
		"UpdatedBy": talent.UpdatedBy,
		"IsActive":  talent.IsActive,
	}, repository.Filter("`talent_id` = ?", talent.ID))
	if err != nil {
		return err
	}
	return nil
}

// setCollegeNameAndID gets the college by college name and sets the college id.
func (service *TalentService) setCollegeNameAndID(uow *repository.UnitOfWork, talent *tal.Talent) error {
	academics := talent.Academics
	for i := 0; i < len(academics); i++ {
		// Get college branch from database.
		collegeBranch := list.Branch{}
		if err := service.Repository.GetRecordForTenant(uow, talent.TenantID, &collegeBranch,
			repository.Filter("`branch_name`=?", academics[i].College)); err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		// Give college id to academics collegeID field.
		academics[i].CollegeID = &collegeBranch.ID
	}
	return nil
}

// getValuesForTalent gets values for talent by firing individual query for each talent.
func (service *TalentService) getValuesForTalent(uow *repository.UnitOfWork, talents *[]tal.DTO, tenantID uuid.UUID) error {
	for index := range *talents {
		//********************************************Courses*********************************************************************

		// Create bucket for courses.
		tempCourses := []crs.Course{}

		// Get all courses form database.
		err := service.Repository.GetAll(uow, &tempCourses,
			repository.Join("LEFT JOIN batches ON courses.`id` = batches.`course_id`"),
			repository.Join("LEFT JOIN batch_talents ON batch_talents.`batch_id` = batches.`id`"),
			repository.Join("LEFT JOIN talents ON batch_talents.`talent_id` = talents.`id`"),
			repository.Filter("talents.`id`=?", (*talents)[index].ID),
			repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("talents.`tenant_id`=? AND talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("courses.`tenant_id`=? AND courses.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("courses.`id`"),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Give courses to talent.
		(*talents)[index].Courses = &tempCourses

		//***************************************************Faculties*****************************************************

		// Create bucket for faculties.
		tempFaculties := []list.Faculty{}

		// Get all faculties form database.
		err = service.Repository.GetAll(uow, &tempFaculties,
			repository.Join("INNER JOIN `batch_modules` ON `batch_modules`.`faculty_id` = `faculties`.`id`"),
			repository.Join("INNER JOIN `batches` ON `batches`.`id` = `batch_modules`.`batch_id`"),
			repository.Join("JOIN batch_talents on batch_talents.`batch_id` = batches.`id`"),
			repository.Filter("batch_talents.`talent_id`=?", (*talents)[index].ID),
			repository.Filter("batch_talents.`tenant_id`=? AND batch_talents.`deleted_at` IS NULL", tenantID),
			repository.Filter("batches.`tenant_id`=? AND batches.`deleted_at` IS NULL", tenantID),
			repository.Filter("faculties.`tenant_id`=? AND faculties.`deleted_at` IS NULL", tenantID),
			repository.Filter("batch_modules.`tenant_id`=? AND batch_modules.`deleted_at` IS NULL", tenantID),
			repository.GroupBy("faculties.`id`"),
		)
		if err != nil {
			uow.RollBack()
			log.NewLogger().Error(err.Error())
			return err
		}

		// Give faculties to talent.
		(*talents)[index].Faculties = &tempFaculties

		//**************************************************ExpectedCTC(call records)***********************************************************

		// Create bucket for expected ctc.
		exepectedCTC := &tal.ExpectedCTCLatest{}

		// Get expected ctc from database.
		if err := service.Repository.Scan(uow, exepectedCTC,
			repository.Filter("talents.`deleted_at` IS NULL AND talent_call_records.`deleted_at` IS NULL AND talents.`tenant_id`=? AND talent_call_records.`tenant_id`=? AND talents.`id`=?",
				tenantID, tenantID, (*talents)[index].ID),
			repository.Filter("talent_call_records.`expected_ctc` IS NOT NULL"),
			repository.Table("talents"),
			repository.Join("JOIN talent_call_records on talent_call_records.`talent_id` = talents.`id`"),
			repository.Select("expected_ctc"),
			repository.GroupBy("talents.`id`"),
			repository.OrderBy("talent_call_records.`date_time`")); err != nil {
			if err != gorm.ErrRecordNotFound {
				log.NewLogger().Error(err.Error())
				uow.RollBack()
				return errors.NewValidationError("Internal server error")
			}
		}

		// Give expected CTC to talent
		if exepectedCTC != nil && exepectedCTC.ExpectedCTC > 0 {
			(*talents)[index].ExpectedCTC = &exepectedCTC.ExpectedCTC
		}
	}
	return nil
}

//temporary function**********************************************************
func (service *TalentService) AddTalentAndCredentialForUpdatingPassword(uow *repository.UnitOfWork, talent *tal.Talent) error {
	credentialID := talent.CreatedBy
	talentExists := false

	var count int
	if err := service.Repository.GetCountForTenant(uow, talent.TenantID, &tal.Talent{}, &count,
		repository.Filter("`email`=? AND `id`=?", talent.Email, talent.ID)); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	if count != 0 { //talent exists
		talentExists = true
	}

	if talentExists { //email exists
		//get talent that exists
		//create talent bucket to get talent already present in DB
		talentInDB := tal.Talent{}
		if err := service.Repository.GetRecordForTenant(uow, talent.TenantID, &talentInDB,
			repository.Filter("`email`=?", talent.Email)); err != nil {
			log.NewLogger().Error(err.Error())
			return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
		}

		//if password is blank in database then generate password and give it to talent's password
		if talentInDB.Password == "" || talentInDB.Password == talentInDB.FirstName {
			password := util.GeneratePassword()
			//update talent with new password
			if err := service.Repository.UpdateWithMap(uow, &talentInDB, map[string]interface{}{"Password": password}); err != nil {
				log.NewLogger().Error(err.Error())
				return errors.NewHTTPError("Talent could not be added", http.StatusInternalServerError)
			}
		}

		//create credential for talent
		if err := service.AddCredential(uow, &talentInDB, talentInDB.CreatedBy); err != nil {
			return err
		}

		//return error to user for same email
		return errors.NewValidationError("Email already exists")

	}
	//email does not exist
	//generate password and give it to talent's password
	talent.Password = util.GeneratePassword()

	// Add talent to database
	if err := service.Repository.Add(uow, talent); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Talent could not be added", http.StatusInternalServerError)
	}

	// Create credential for talent
	if err := service.AddCredential(uow, talent, credentialID); err != nil {
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Internal server error", http.StatusInternalServerError)
	}
	return nil
}

//===================================================UPDATE EXPERIENCE IN MONTHS OF ALL TALENTS===========================================================================

// UpdateExperienceInMonthsOfTalent updates the total expereince of each talent.
func (service *TalentService) UpdateExperienceInMonthsOfTalent() error {

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update experience in months of all talents.
	err := service.Repository.Scan(uow, tal.Talent{},
		repository.RawQuery("UPDATE talents AS t1 INNER JOIN "+
			"(SELECT sum( TIMESTAMPDIFF(month, from_date, if(to_date is null, curdate(), to_date))) as months, "+
			"talent_experiences.* FROM talent_experiences where talent_experiences.deleted_at IS NULL GROUP BY talent_id) t2 "+
			"ON t1.id = t2.talent_id SET t1.experience_in_months = months"))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		log.NewLogger().Error(err.Error())
		return errors.NewHTTPError("Talents could not be updated", http.StatusInternalServerError)
	}
	return nil
}

// getBatchSessionsCount will return count of total sessions and completed sessions for the specified batch
func (service *TalentService) getBatchSessionsCount(uow *repository.UnitOfWork,
	totalSession, completedSession *uint, tenantID, batchID uuid.UUID) error {

	err := service.Repository.GetCountForTenant(uow, tenantID, bat.Session{}, totalSession,
		repository.Filter("batch_sessions.`batch_id`=?", batchID))
	if err != nil {
		return err
	}

	err = service.Repository.GetCountForTenant(uow, tenantID, bat.Session{}, completedSession,
		repository.Filter("batch_sessions.`batch_id` = ? AND batch_sessions.`is_completed` = ?", batchID, true))
	if err != nil {
		return err
	}

	return nil
}
