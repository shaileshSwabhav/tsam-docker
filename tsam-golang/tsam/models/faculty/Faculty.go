package faculty

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//Faculty will contain all the details regarding a faculty..
type Faculty struct {
	general.TenantBase
	general.Address
	Code          string                `json:"code" gorm:"type:varchar(10);not null"`
	FirstName     string                `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName      string                `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email         string                `json:"email" example:"abc@gmail.com" gorm:"type:varchar(50)"`
	Contact       string                `json:"contact" example:"9700795509" gorm:"type:varchar(15)"`
	DateOfBirth   *string               `json:"dateOfBirth" example:"2000-02-30" gorm:"type:date"`
	DateOfJoining *string               `json:"dateOfJoining" example:"2019-02-23" gorm:"type:date"`
	Technologies  []*general.Technology `json:"technologies" gorm:"many2many:faculties_technologies;association_autoupdate:false"`
	Resume        *string               `json:"resume" gorm:"type:varchar(200)"`
	Academics     []*Academic           `json:"academics"`
	Experiences   []*Experience         `json:"experiences"`
	IsActive      *bool                 `json:"isActive" gorm:"DEFAULT:true"`
	IsFullTime    *bool                 `json:"isFullTime" gorm:"type:tinyint(1)"`
	TelegramID    *string               `json:"telegramID" gorm:"type:varchar(50)"`

	// temp
	Password               string  `json:"-" gorm:"type:varchar(255)"`
	AverageAssessmentScore float64 `json:"averageAssessmentScore" gorm:"-"`
	AverageTechnicalScore  float64 `json:"averageTechnicalScore" gorm:"-"`
}

// Validate validate faculty details.
func (faculty *Faculty) Validate() error {

	if util.IsEmpty(faculty.FirstName) || !util.ValidateString(faculty.FirstName) {
		return errors.NewValidationError("Faculty FirstName must be specified and must have characters only")
	}

	if util.IsEmpty(faculty.LastName) || !util.ValidateString(faculty.LastName) {
		return errors.NewValidationError("Faculty LastName must be specified and must have characters only")
	}

	if util.IsEmpty(faculty.Email) || !util.ValidateEmail(faculty.Email) {
		return errors.NewValidationError("Faculty Email must be specified and should be of the type abc@domain.com")
	}

	if util.IsEmpty(faculty.Contact) || !util.ValidateContact(faculty.Contact) {
		return errors.NewValidationError("Faculty Contact must be specified and have 10 digits")
	}

	if !util.IsUUIDValid(faculty.TenantID) {
		return errors.NewValidationError("Invalid tenantID")
	}

	if faculty.Academics == nil {
		return errors.NewValidationError("Academics must be specified")
	}

	for _, academic := range faculty.Academics {
		if err := academic.ValidateFacultyAcademic(); err != nil {
			return errors.NewValidationError(err.Error())
		}
	}

	if faculty.Experiences != nil {
		for _, experience := range faculty.Experiences {
			if err := experience.ValidateFacultyExperiences(); err != nil {
				return errors.NewValidationError(err.Error())
			}
		}
	}

	// if faculty.Technologies == nil {
	// 	return errors.NewValidationError("Technologies must be specified")
	// }
	// if faculty.Technologies != nil {
	// 	for _, technology := range faculty.Technologies {
	// 		if err := technology.Validate(); err != nil {
	// 			return errors.NewValidationError(err.Error())
	// 		}
	// 	}
	// }

	return nil
}

// FacultyDTO will provide faculty details.
type FacultyDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	general.Address

	Code string `json:"code"`

	FirstName     string  `json:"firstName" example:"Ravi"`
	LastName      string  `json:"lastName" example:"Sharma"`
	Email         string  `json:"email" example:"abc@gmail.com"`
	Contact       string  `json:"contact" example:"9700795509"`
	DateOfBirth   *string `json:"dateOfBirth" example:"2000-02-30"`
	DateOfJoining *string `json:"dateOfJoining" example:"2019-02-23"`
	Resume        *string `json:"resume"`

	Technologies []*general.Technology `json:"technologies" gorm:"many2many:faculties_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:faculty_id"`
	Academics    []*AcademicDTO        `json:"academics" gorm:"foreignkey:FacultyID"`
	Experiences  []*ExperienceDTO      `json:"experiences" gorm:"foreignkey:FacultyID"`

	IsActive   *bool `json:"isActive"`
	IsFullTime *bool `json:"isFullTime"`

	Password               string  `json:"-"`
	AverageAssessmentScore float64 `json:"averageAssessmentScore"`
	AverageTechnicalScore  float64 `json:"averageTechnicalScore"`
}

// TableName defines table name of the struct.
func (*FacultyDTO) TableName() string {
	return "faculties"
}
