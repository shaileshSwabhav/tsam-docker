package talent

import (
	"regexp"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	crs "github.com/techlabs/swabhav/tsam/models/course"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Talent is used for adding and updating talent, consists of IDs of related information.
type Talent struct {
	general.TenantBase

	// Personal Details.
	general.Address
	Code               string   `json:"code" gorm:"type:varchar(10);not null"`
	FirstName          string   `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName           string   `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email              string   `json:"email" example:"john@doe.com" gorm:"type:varchar(100)"`
	Contact            string   `json:"contact" example:"1234567890" gorm:"type:varchar(15)"`
	AcademicYear       uint8    `json:"academicYear" example:"1" gorm:"type:tinyint(2)"`
	TalentType         *uint8   `json:"talentType" example:"4" gorm:"type:tinyint(2)"`
	PersonalityType    *string  `json:"personalityType" example:"Introvert" gorm:"type:varchar(50)"`
	Resume             *string  `json:"resume" gorm:"type:varchar(200)"`
	AlternateContact   *string  `json:"alternateContact" example:"9730220343" gorm:"type:varchar(15)"`
	AlternateEmail     *string  `json:"alternateEmail" example:"9730220343" gorm:"type:varchar(100)"`
	LoyaltyPoints      *float32 `json:"loyaltyPoints" example:"20.00"`
	Password           string   `json:"-" gorm:"type:varchar(100)"`
	LifetimeValue      *uint    `json:"lifetimeValue"`
	ExperienceInMonths *uint    `json:"experienceInMonths"`
	Image              *string  `json:"image" gorm:"type:varchar(200)"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talents_technologies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	SalesPersonID *uuid.UUID `json:"salesPersonID" gorm:"type:varchar(36)"`
	SourceID      *uuid.UUID `json:"sourceID" gorm:"type:varchar(36)"`
	ReferralID    *uuid.UUID `json:"referralID" gorm:"type:varchar(36)"`

	// Flags.
	IsActive        bool `json:"isActive" example:"true" gorm:"type:tinyint(1)"`
	IsMastersAbroad bool `json:"isMastersAbroad"`
	IsSwabhavTalent bool `json:"isSwabhavTalent"`
	IsExperience    bool `json:"isExperience" example:"true" gorm:"type:tinyint(1)"`

	// Child tables.
	Academics     []*Academic            `json:"academics"`
	Experiences   []*Experience          `json:"experiences"`
	MastersAbroad *general.MastersAbroad `json:"mastersAbroad"`

	// Social media.
	FacebookURL  *string `json:"facebookUrl" gorm:"type:varchar(200)"`
	InstagramURL *string `json:"instagramUrl" gorm:"type:varchar(200)"`
	GithubURL    *string `json:"githubUrl" gorm:"type:varchar(200)"`
	LinkedInURL  *string `json:"linkedInUrl" gorm:"type:varchar(200)"`
}

// ***********Removed Fields************ you must delete all references
// TestHistories   []*TestHistory        `json:"testHistories" gorm:"many2many:talent_test_histories;"`

// Validate validates compulsary fields of talent.
func (talent *Talent) Validate() error {

	// City maximum characters.
	if talent.City != nil && len(*talent.City) > 50 {
		return errors.NewValidationError("City can have maximum 50 characters")
	}

	// Check if first name is blank or not.
	if util.IsEmpty(talent.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(talent.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// First name maximum characters.
	if len(talent.FirstName) > 50 {
		return errors.NewValidationError("First name can have maximum 50 characters")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(talent.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(talent.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Last name maximum characters.
	if len(talent.LastName) > 50 {
		return errors.NewValidationError("Last name can have maximum 50 characters")
	}

	// Check if email is blank or not.
	if util.IsEmpty(talent.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(talent.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Email maximum characters.
	if len(talent.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// Alternate Email should be in proper format.
	if talent.AlternateEmail != nil && !util.ValidateEmail(*talent.AlternateEmail) {
		return errors.NewValidationError("Alternate Email must in the format : email@example.com")
	}

	// Alternate Email maximum characters.
	if talent.AlternateEmail != nil && len(*talent.AlternateEmail) > 100 {
		return errors.NewValidationError("Alternate Email can have maximum 100 characters")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(talent.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(talent.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Alternate Contact should consist of only 10 numbers.
	if talent.AlternateContact != nil && !util.ValidateContact(*talent.AlternateContact) {
		return errors.NewValidationError("Alternate Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if talent.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// Salesperson ID.
	if talent.SalesPersonID != nil && !util.IsUUIDValid(*talent.SalesPersonID) {
		return errors.NewValidationError("Salesperson ID must ne a proper uuid")
	}

	// Source ID.
	if talent.SourceID != nil && !util.IsUUIDValid(*talent.SourceID) {
		return errors.NewValidationError("Source ID must be a proper uuid")
	}

	// Referral ID.
	if talent.ReferralID != nil && !util.IsUUIDValid(*talent.ReferralID) {
		return errors.NewValidationError("Referral ID must be a proper uuid")
	}

	// Validate address.
	if err := talent.Address.ValidateAddress(); err != nil {
		return err
	}

	// Check if academics exist, if true then validate.
	if talent.Academics != nil {
		for _, academic := range talent.Academics {
			if err := academic.Validate(); err != nil {
				return err
			}
		}
	}

	// Check if experiences exist, if true then validate.
	if talent.Experiences != nil {
		isWorking := false
		for _, experience := range talent.Experiences {
			if experience.ToDate == nil {
				if isWorking {
					return errors.NewValidationError("More than one experince cannot have To Date field as null")
				}
				isWorking = true
			}
			if err := experience.Validate(); err != nil {
				return err
			}
		}
	}

	// Validate master abroad fields.
	if talent.MastersAbroad != nil {
		err := talent.MastersAbroad.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

// ===========Defining many to many structs===========

// TalentTechnologies is the map of talent and technology.
type TalentTechnologies struct {
	TalentID     uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*TalentTechnologies) TableName() string {
	return "talents_technologies"
}

// TalentUpdate is the struct for updating talents' salesperson.
type TalentUpdate struct {
	TalentID uuid.UUID `json:"talentID"`
}

// EligibleTalent is the struct for getting qualification. technologies, experience technologies and designations of talent.
type EligibleTalentDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Multiple Model.
	Academics   []*AcademicDTO   `json:"academics" gorm:"foreignkey:TalentID"`
	Experiences []*ExperienceDTO `json:"experiences" gorm:"foreignkey:TalentID"`

	// Flags.
	IsExperience *bool `json:"isExperience"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talents_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:talent_id"`
}

// TableName defines table name of the struct.
func (*EligibleTalentDTO) TableName() string {
	return "talents"
}

//************************************* DTO MODEL *************************************************************

// DTO is used for getting talent from database including all the talent related information also.
type DTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Personal Details.
	general.Address
	Code             string     `json:"code"`
	FirstName        string     `json:"firstName" example:"Ravi"`
	LastName         string     `json:"lastName" example:"Sharma"`
	Email            string     `json:"email" example:"john@doe.com"`
	Contact          string     `json:"contact" example:"1234567890"`
	AlternateContact *string    `json:"alternateContact" example:"9730220343"`
	AlternateEmail   *string    `json:"alternateEmail" example:"9730220343"`
	AcademicYear     uint8      `json:"academicYear" example:"1"`
	TalentType       *uint8     `json:"talentType" example:"4"`
	PersonalityType  *string    `json:"personalityType" example:"Introvert"`
	ExperienceInMonths *uint    `json:"experienceInMonths"`
	Resume           *string    `json:"resume"`
	ReferralID       *uuid.UUID `json:"referralID"`
	LoyaltyPoints    *float32   `json:"loyaltyPoints" example:"20.00"`
	Password         string     `json:"-"`
	LifetimeValue    *uint      `json:"lifetimeValue"`
	ExpectedCTC      *int64     `json:"expectedCTC"`
	Image            *string    `json:"image"`

	// Single model.
	SalesPerson   *User                     `json:"salesPerson" gorm:"foreignkey:SalesPersonID"`
	SalesPersonID *uuid.UUID                `json:"-"`
	TalentSource  *general.Source           `json:"talentSource" gorm:"foreignkey:SourceID"`
	SourceID      *uuid.UUID                `json:"-"`
	MastersAbroad *general.MastersAbroadDTO `json:"mastersAbroad" gorm:"foreignkey:TalentID"`

	// Multiple Model.
	Academics   []*AcademicDTO   `json:"academics" gorm:"foreignkey:TalentID"`
	Experiences []*ExperienceDTO `json:"experiences" gorm:"foreignkey:TalentID"`
	Faculties   *[]list.Faculty  `json:"faculties"`
	Courses     *[]crs.Course    `json:"courses"`

	// Maps.
	Technologies []*general.Technology `json:"technologies" gorm:"many2many:talents_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:talent_id"`

	// Flags.
	IsSwabhavTalent *bool `json:"isSwabhavTalent"`
	IsMastersAbroad bool  `json:"isMastersAbroad"`
	IsExperience    *bool `json:"isExperience" example:"true"`
	IsActive        *bool `json:"isActive" example:"true"`

	// Social Media.
	FacebookURL  *string `json:"facebookUrl"`
	InstagramURL *string `json:"instagramUrl"`
	GithubURL    *string `json:"githubUrl"`
	LinkedInURL  *string `json:"linkedInUrl"`
}

// TableName defines table name of the struct.
func (*DTO) TableName() string {
	return "talents"
}

//************************************* EXCEL MODEL *************************************************************

// TalentExcel is used for adding data coming from excel file.
type TalentExcel struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Contact         string `json:"contact"`
	AcademicYear    uint8  `json:"academicYear"`
	IsActive        bool   `json:"isActive"`
	IsSwabhavTalent bool   `json:"isSwabhavTalent"`

	Address     *string `json:"address,omitempty"`
	City        *string `json:"city,omitempty"`
	PINCode     *uint32 `json:"pinCode,omitempty"`
	CountryName *string `json:"countryName,omitempty"`
	StateName   *string `json:"stateName,omitempty"`

	Academics []*AcademicExcel `json:"academics"`
}

// Validate validates compulsary fields of talent excel.
func (talent *TalentExcel) Validate() error {
	pinCodePattern := regexp.MustCompile("^[1-9][0-9]{5}$")

	// Check if first name is blank or not.
	if util.IsEmpty(talent.FirstName) {
		return errors.NewValidationError("First name must be specified")
	}

	// First name should consist of only alphabets.
	if !util.ValidateString(talent.FirstName) {
		return errors.NewValidationError("First name should consist of only alphabets")
	}

	// First name maximum characters.
	if len(talent.FirstName) > 50 {
		return errors.NewValidationError("First name can have maximum 50 characters")
	}

	// Check if last name is blank or not.
	if util.IsEmpty(talent.LastName) {
		return errors.NewValidationError("Last name must be specified")
	}

	// Last name should consist of only alphabets.
	if !util.ValidateString(talent.LastName) {
		return errors.NewValidationError("Last name should consist of only alphabets")
	}

	// Last name maximum characters.
	if len(talent.LastName) > 50 {
		return errors.NewValidationError("Last name can have maximum 50 characters")
	}

	// Check if email is blank or not.
	if util.IsEmpty(talent.Email) {
		return errors.NewValidationError("Email must be specified")
	}

	// Email should be in proper format.
	if !util.ValidateEmail(talent.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Email maximum characters.
	if len(talent.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// Check if contact is blank or not.
	if util.IsEmpty(talent.Contact) {
		return errors.NewValidationError("Contact must be specified")
	}

	// Contact should consist of only 10 numbers.
	if !util.ValidateContact(talent.Contact) {
		return errors.NewValidationError("Contact should consist of only 10 numbers")
	}

	// Check if academic year is present or not.
	if talent.AcademicYear == 0 {
		return errors.NewValidationError("Academic Year must be specified")
	}

	// City should consist of only alphabets.
	if talent.City != nil {
		if !util.ValidateStringWithSpace(*talent.City) {
			return errors.NewValidationError("City should have only characters and space")
		}
	}

	// PIN Code should be in proper format.
	if talent.PINCode != nil {
		if !pinCodePattern.MatchString(strconv.Itoa(int(*talent.PINCode))) {
			return errors.NewValidationError("Invalid Pincode")
		}
	}

	// Check if academics exist, if true then validate.
	if talent.Academics != nil {
		for _, academic := range talent.Academics {
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

// ExcelDTO is used for getting talent from database including all the talent related information also for excel.
type ExcelDTO struct {

	// Personal Details.
	Code             string     `json:"code"`
	FirstName        string     `json:"firstName"`
	LastName         string     `json:"lastName"`
	Email            string     `json:"email"`
	Mobile           string     `json:"mobile"`
	AlternateEmail   *string    `json:"alternateEmail"`
	AlternateMobile *string    `json:"alternateMobile"`
	AcademicYear     uint8      `json:"academicYear"`
	Country *string    			`json:"country"`
	State *string    			`json:"state"`
	City *string    			`json:"city"`
	Pincode   *uint32    		`json:"pincode"`
	TalentType       *uint8     `json:"talentType"`
	PersonalityType  *string    `json:"personalityType"`
	SalesPerson   *string        `json:"salesPerson"`
	Source  *string      `json:"source"`
	IsSwabhavTalent *bool `json:"isSwabhavTalent"`
	IsActive        *bool `json:"isActive" example:"true"`
	FacebookURL  *string `json:"facebookUrl"`
	InstagramURL *string `json:"instagramUrl"`
	GithubURL    *string `json:"githubUrl"`
	LinkedInURL  *string `json:"linkedInUrl"`
	Technologies *string `json:"technologies"`
	LoyaltyPoints    *float32   `json:"loyaltyPoints"`
	Resume           *string    `json:"resume"`
	TotalYearOfExp  *uint    `json:"totalYearOfExp"`
	LifetimeValue    *uint      `json:"lifetimeValue"`
	Courses     *string    `json:"courses"`
	Faculties    *string  `json:"faculties"`
	ExpectedCTC      *int64     `json:"expectedCTC"`

	// Academic Details.
	Qualification          	*string         `json:"qualification"`
	Specialization   	*string  `json:"specialization"`
	CollegeName          *string                 `json:"collegeName"`
	CGPA        		*float32             `json:"CGPA"`
	YearOfPassout    *uint16    `json:"yearOfPassout"`

	// Experience Details.
	ExperienceTechnologies *string `json:"experienceTechnologies"`
	Designation    *string `json:"designation"`
	CurrentCompany  *string    `json:"currentCompany"`
	FromYear *string    `json:"fromYear"`
	ToYear   *string   `json:"toYear"`
	Package  *uint     `json:"package"`
}