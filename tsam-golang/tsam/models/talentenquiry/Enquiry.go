package talentenquiry

import (
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	tal "github.com/techlabs/swabhav/tsam/models/talent"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Enquiry is used for adding and updating enquiry, consists of IDs of related information.
type Enquiry struct {
	general.TenantBase

	// Personal Details.
	general.Address
	Code              string     `json:"code" gorm:"type:varchar(10);not null"`
	FirstName         string     `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName          string     `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email             string     `json:"email" gorm:"type:varchar(100)"`
	Contact           string     `json:"contact" gorm:"type:varchar(15)"`
	AlternateContact  *string    `json:"alternateContact" example:"9730220343" gorm:"type:varchar(15)"`
	AlternateEmail    *string    `json:"alternateEmail" example:"9730220343" gorm:"type:varchar(100)"`
	EnquiryDate       string     `json:"enquiryDate" gorm:"type:date"`
	EnquiryType       string     `json:"enquiryType" gorm:"type:varchar(50)"`
	TalentID          *uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	AdditionalDetails *string    `json:"additionalDetails" gorm:"type:varchar(1000)"`
	AcademicYear      uint8      `json:"academicYear" example:"1" gorm:"type:tinyint(2)"`
	Resume            *string    `json:"resume" gorm:"type:varchar(200)"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talent_enquiries_technologies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Courses      []*list.Course        `json:"courses" gorm:"many2many:talent_enquiries_courses;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	SalesPersonID *uuid.UUID `json:"salesPersonID" gorm:"type:varchar(36)"`
	SourceID      *uuid.UUID `json:"sourceID" gorm:"type:varchar(36)"`

	// Flags.
	IsMastersAbroad bool `json:"isMastersAbroad"`
	IsExperience    bool `json:"isExperience" example:"true" gorm:"type:tinyint(1)"`

	// Child tables.
	Academics     []*Academic            `json:"academics"`
	Experiences   []*Experience          `json:"experiences"`
	MastersAbroad *general.MastersAbroad `json:"mastersAbroad"`

	// Social media.
	FacebookUrl  *string `json:"facebookUrl" gorm:"type:varchar(200)"`
	InstagramUrl *string `json:"instagramUrl" gorm:"type:varchar(200)"`
	GithubUrl    *string `json:"githubUrl" gorm:"type:varchar(200)"`
	LinkedInUrl  *string `json:"linkedInUrl" gorm:"type:varchar(200)"`
}

// TableName will name the table of Enquiry model as "talent_enquiries".
func (*Enquiry) TableName() string {
	return "talent_enquiries"
}

// Validate validates the fields of enquiry.
func (enquiry *Enquiry) Validate(isEnquiryForm bool) error {

	// City maximum characters.
	if enquiry.City != nil && len(*enquiry.City) > 50 {
		return errors.NewValidationError("City can have maximum 50 characters")
	}

	// Check if first name is blank or not.
	if util.IsEmpty(enquiry.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(enquiry.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(enquiry.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(enquiry.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Check if email is blank or not.
	if util.IsEmpty(enquiry.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(enquiry.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Alternate Email should be in proper format.
	if enquiry.AlternateEmail != nil && !util.ValidateEmail(*enquiry.AlternateEmail) {
		return errors.NewValidationError("Alternate Email must in the format : email@example.com")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(enquiry.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(enquiry.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Alternate Contact should consist of only 10 numbers.
	if enquiry.AlternateContact != nil && !util.ValidateContact(*enquiry.AlternateContact) {
		return errors.NewValidationError("Alternate Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if enquiry.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// Check if enquiry type is empty or not.
	if util.IsEmpty(enquiry.EnquiryType) {
		return errors.NewValidationError("Enquiry Type must be specified")
	}

	// Check if enquiry date is empty or not.
	if util.IsEmpty(enquiry.EnquiryDate) {
		return errors.NewValidationError("Enquiry Date must be specified")
	}

	// Validate address.
	if err := enquiry.Address.ValidateAddress(); err != nil {
		return err
	}

	// Check if additional details has more than maximum characters or not.
	if enquiry.AdditionalDetails != nil && len(*enquiry.AdditionalDetails) > 1000 {
		return errors.NewValidationError("Additional lDetails cannot have more than 1000 characters")
	}

	// Check if academics exist, if true then validate.
	if enquiry.Academics != nil {
		for _, academic := range enquiry.Academics {
			if err := academic.Validate(); err != nil {
				return err
			}
		}
	}

	// Check if experiences exist, if true then validate.
	if enquiry.Experiences != nil {
		isWorking := false
		for _, experience := range enquiry.Experiences {
			if experience.ToDate == nil {
				if isWorking {
					return errors.NewValidationError("More than one experince cannot be currently working")
				}
				isWorking = true
			}
			if err := experience.Validate(); err != nil {
				return err
			}
		}
	}

	// Validate master abroad fields.
	if enquiry.MastersAbroad != nil {
		err := enquiry.MastersAbroad.Validate()
		if err != nil {
			return err
		}
	}

	//*********************************CHECKS FOR NON COMPULSARY FIELDS(OLD CODE)****************************
	// // Check if academic year is present or not.
	// if !isEnquiryForm && (enquiry.AcademicYear == nil || (enquiry.AcademicYear != nil && *enquiry.AcademicYear == 0)) {
	// 	return errors.NewValidationError("Academic Year must be specified")
	// }

	// // Check if enquiry type is empty or not.
	// if !isEnquiryForm && (enquiry.EnquiryType == nil || (enquiry.EnquiryType != nil && util.IsEmpty(*enquiry.EnquiryType))) {
	// 	return errors.NewValidationError("Enquiry Type must be specified")
	// }

	// // Check if academics exist, if true then validate.
	// if enquiry.Academics != nil {
	// 	for _, academic := range enquiry.Academics {
	// 		if err := academic.Validate(isEnquiryForm); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

// ===========Defining many to many structs===========

// EnquiryTechnologies is the map of enquiry and technology.
type EnquiryTechnologies struct {
	EnquiryID    uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*EnquiryTechnologies) TableName() string {
	return "talent_enquiries_technologies"
}

// EnquiryCourses is the map of enquiry and course.
type EnquiryCourses struct {
	EnquiryID uuid.UUID `gorm:"type:varchar(36)"`
	CourseID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*EnquiryCourses) TableName() string {
	return "talent_enquiries_courses"
}

// Struct for update enquiries' salesperson.
type EnquiryUpdate struct {
	EnquiryID uuid.UUID `json:"enquiryID"`
}

// ApplicationForm is the model for adding enquiry and waiting list.
type ApplicationForm struct {
	Enquiry     Enquiry         `json:"enquiry"`
	WaitingList tal.WaitingList `json:"waitingList"`
}

//************************************* DTO MODEL *************************************************************

// DTO is used for getting enquiry from database including all the enquiry related information also.
type DTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Personal Details.
	general.Address
	Code              string     `json:"code"`
	FirstName         string     `json:"firstName" example:"Ravi"`
	LastName          string     `json:"lastName" example:"Sharma"`
	Email             string     `json:"email"`
	Contact           string     `json:"contact"`
	AlternateContact  *string    `json:"alternateContact" example:"9730220343"`
	AlternateEmail    *string    `json:"alternateEmail" example:"9730220343"`
	EnquiryDate       string     `json:"enquiryDate"`
	EnquiryType       *string    `json:"enquiryType"`
	TalentID          *uuid.UUID `json:"talentID"`
	AdditionalDetails *string    `json:"additionalDetails"`
	AcademicYear      *uint8     `json:"academicYear" example:"1"`
	Resume            *string    `json:"resume"`
	ExpectedCTC       *int64     `json:"expectedCTC"`

	// Single model.
	SalesPerson   *User                     `json:"salesPerson" gorm:"foreignkey:SalesPersonID"`
	SalesPersonID *uuid.UUID                `json:"-"`
	EnquirySource *general.Source           `json:"enquirySource" gorm:"foreignkey:SourceID"`
	SourceID      *uuid.UUID                `json:"-"`
	MastersAbroad *general.MastersAbroadDTO `json:"mastersAbroad" gorm:"foreignkey:EnquiryID"`

	// Multiple Model.
	Academics   []*AcademicDTO   `json:"academics" gorm:"foreignkey:EnquiryID"`
	Experiences []*ExperienceDTO `json:"experiences" gorm:"foreignkey:EnquiryID"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talent_enquiries_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:enquiry_id"`
	Courses      []*list.Course        `json:"courses" gorm:"many2many:talent_enquiries_courses;association_jointable_foreignkey:course_id;jointable_foreignkey:enquiry_id"`

	// Flags.
	IsMastersAbroad bool  `json:"isMastersAbroad"`
	IsExperience    *bool `json:"isExperience" example:"true"`

	// Social Media.
	FacebookUrl  *string `json:"facebookUrl"`
	InstagramUrl *string `json:"instagramUrl"`
	GithubUrl    *string `json:"githubUrl"`
	LinkedInUrl  *string `json:"linkedInUrl"`
}

// TableName will name the table of Enquiry model as "talent_enquiries".
func (*DTO) TableName() string {
	return "talent_enquiries"
}

//************************************* EXCEL MODEL *************************************************************

// EnquiryExcel is used for adding data coming from excel file.
type EnquiryExcel struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Contact      string `json:"contact"`
	AcademicYear uint8  `json:"academicYear"`

	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	PINCode     *uint32 `json:"pinCode,omitempty"`
	CountryName *string `json:"countryName,omitempty"`
	StateName   *string `json:"stateName,omitempty"`

	Academics []*AcademicExcel `json:"academics"`
}

// Validate validates compulsary fields of enquiry excel.
func (enquiry *EnquiryExcel) Validate() error {
	pinCodePattern := regexp.MustCompile("^[1-9][0-9]{5}$")

	// Check if first name is blank or not.
	if util.IsEmpty(enquiry.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(enquiry.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// First name maximum characters.
	if len(enquiry.FirstName) > 50 {
		return errors.NewValidationError("First name can have maximum 50 characters")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(enquiry.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(enquiry.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Last name maximum characters.
	if len(enquiry.LastName) > 50 {
		return errors.NewValidationError("Last name can have maximum 50 characters")
	}

	// Check if email is blank or not.
	if util.IsEmpty(enquiry.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(enquiry.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Email maximum characters.
	if len(enquiry.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(enquiry.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(enquiry.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if enquiry.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// City should consist of only alphabets.
	if enquiry.City != nil {
		if !util.ValidateStringWithSpace(*enquiry.City) {
			return errors.NewValidationError("City should have only characters and space")
		}
	}

	// PIN Code should be in proper format.
	if enquiry.PINCode != nil {
		if !pinCodePattern.MatchString(strconv.Itoa(int(*enquiry.PINCode))) {
			return errors.NewValidationError("Invalid Pincode")
		}
	}

	// Check if academics exist, if true then validate.
	if enquiry.Academics != nil {
		for _, academic := range enquiry.Academics {
			if err := academic.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

// IDModel is used for gettig id from database.
type IDModel struct {
	ID uuid.UUID
}
