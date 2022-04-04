package company

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Enquiry contains all fields of company enquiry.
type Enquiry struct {
	general.TenantBase

	// Address.
	general.Address

	// Maps.
	Domains      []Domain             `json:"domains" gorm:"many2many:company_enquiries_domains;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Technologies []general.Technology `json:"technologies" gorm:"many2many:company_enquiries_technologies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	CompanyBranchID *uuid.UUID `json:"companyBranchID" gorm:"type:varchar(36)"`
	SalesPersonID   *uuid.UUID `json:"salesPersonID" gorm:"type:varchar(36)"`

	// Other fields.
	CompanyName    string  `json:"companyName" gorm:"type:varchar(200)"`
	Website        *string `json:"website" gorm:"type:varchar(100)"`
	Email          *string `json:"email" gorm:"type:varchar(100)"`
	Code           string  `json:"code" gorm:"type:varchar(10);not null"`
	HRName         *string `json:"hrName" gorm:"type:varchar(100)"`
	HRContact      *string `json:"hrContact" gorm:"type:varchar(20)"`
	FounderName    *string `json:"founderName" gorm:"type:varchar(100)"`
	Vacancy        *uint16 `json:"vacancy" gorm:"type:smallint(5)"`
	EnquiryDate    *string `json:"enquiryDate" gorm:"type:date"`
	EnquiryType    *string `json:"enquiryType" gorm:"varchar(50)"`
	EnquirySource  *string `json:"enquirySource" gorm:"varchar(50)"`
	JobRole        *string `json:"jobRole" gorm:"varchar(50)"`
	PackageOffered *uint64 `json:"packageOffered" gorm:"type:int(10)"`
	Subject        string  `json:"subject" gorm:"type:varchar(500)"`
	Message        string  `json:"message" gorm:"type:varchar(4000)"`
}

// TableName will name the table of enquiry model as "company_enquiries".
func (*Enquiry) TableName() string {
	return "company_enquiries"
}

// Validate validates CompanyEnquiry.
func (enquiry *Enquiry) Validate() error {
	// Validate address.
	if err := enquiry.Address.MandatoryValidation(); err != nil {
		return err
	}

	// Check if company name is blank or not.
	if util.IsEmpty(enquiry.CompanyName) {
		return errors.NewValidationError("Company name must be specified")
	}

	// Company name maximum characters.
	if len(enquiry.CompanyName) > 200 {
		return errors.NewValidationError("Company name can have maximum 200 characters")
	}

	// Website maximum characters.
	if enquiry.Website != nil && len(*enquiry.Website) > 100 {
		return errors.NewValidationError("Website can have maximum 100 characters")
	}

	// Email maximum characters.
	if enquiry.Email != nil && len(*enquiry.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// HR Name maximum characters.
	if enquiry.HRName != nil && len(*enquiry.HRName) > 100 {
		return errors.NewValidationError("HR Name can have maximum 100 characters")
	}

	// HR Contact should consist of only 10 numbers.
	if enquiry.HRContact != nil && !util.ValidateContact(*enquiry.HRContact) {
		return errors.NewValidationError("HR Contact should consist of only 10 numbers")
	}

	// Founder Name maximum characters.
	if enquiry.FounderName != nil && len(*enquiry.FounderName) > 100 {
		return errors.NewValidationError("Founder Name can have maximum 100 characters")
	}

	// Vacancy less than 0.
	minValue := 0
	if enquiry.Vacancy != nil && *enquiry.Vacancy < uint16(minValue) {
		return errors.NewValidationError("Vacancy cannot be less than 0")
	}

	// Vacancy more than 99999.
	maxValue := 10000
	if enquiry.Vacancy != nil && *enquiry.Vacancy > uint16(maxValue) {
		return errors.NewValidationError("Vacancy cannot be more than 10000")
	}

	// Job Role characters.
	if enquiry.JobRole != nil && len(*enquiry.JobRole) > 50 {
		return errors.NewValidationError("Job Role can have maximum 50 characters")
	}

	// Package Offered less than 0.
	minValue = 0
	if enquiry.PackageOffered != nil && *enquiry.PackageOffered < uint64(minValue) {
		return errors.NewValidationError("Package Offered cannot be less than 0")
	}

	// Package Offered more than 9999999999.
	maxValue = 9999999999
	if enquiry.PackageOffered != nil && *enquiry.PackageOffered > uint64(maxValue) {
		return errors.NewValidationError("Package Offered cannot be more than 9999999999")
	}

	// Check if subject is blank or not.
	if util.IsEmpty(enquiry.Subject) {
		return errors.NewValidationError("Subject must be specified")
	}

	// Check if subject is blank or not.
	if util.IsEmpty(enquiry.Subject) {
		return errors.NewValidationError("Subject must be specified")
	}

	// Subject maximum characters.
	if len(enquiry.Subject) > 500 {
		return errors.NewValidationError("Subject can have maximum 500 characters")
	}

	// Check if message is blank or not.
	if util.IsEmpty(enquiry.Message) {
		return errors.NewValidationError("Message must be specified")
	}

	// Message maximum characters.
	if len(enquiry.Message) > 4000 {
		return errors.NewValidationError("Message can have maximum 4000 characters")
	}

	// Salesperson ID.
	if enquiry.SalesPersonID != nil && !util.IsUUIDValid(*enquiry.SalesPersonID) {
		return errors.NewValidationError("Salesperson ID must be a proper uuid")
	}

	// Technologies.
	if enquiry.Technologies == nil || (enquiry.Technologies != nil && len(enquiry.Technologies) == 0) {
		return errors.NewValidationError("Technology(s) must be specified")
	}

	// Domains.
	if enquiry.Domains == nil || (enquiry.Domains != nil && len(enquiry.Domains) == 0) {
		return errors.NewValidationError("Domain(s) must be specified")
	}

	return nil
}

// ===========Defining many to many structs===========

// CompanyEnquiryTechnologies is the map of enquiry and technology.
type CompanyEnquiryTechnologies struct {
	EnquiryID    uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyEnquiryTechnologies) TableName() string {
	return "company_enquiries_technologies"
}

// CompanyEnquiryDomains is the map of enquiry and domain.
type CompanyEnquiryDomains struct {
	EnquiryID uuid.UUID `gorm:"type:varchar(36)"`
	DomainID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*CompanyEnquiryDomains) TableName() string {
	return "company_enquiries_domains"
}

// EnquiryUpdate is the struct for updating enquiry' salesperson.
type EnquiryUpdate struct {
	EnquiryID uuid.UUID `json:"enquiryID"`
}

// Enquiry contains all fields of company enquiry.
type EnquiryDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Address.
	general.Address

	// Maps.
	Domains      []Domain             `json:"domains" gorm:"many2many:company_enquiries_domains;association_jointable_foreignkey:domain_id;jointable_foreignkey:enquiry_id"`
	Technologies []general.Technology `json:"technologies" gorm:"many2many:company_enquiries_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:enquiry_id"`

	// Single model.
	SalesPerson   *SalesPerson `json:"salesPerson" gorm:"foreignkey:SalesPersonID"`
	SalesPersonID *uuid.UUID   `json:"-"`
	// CompanyBranchID *uuid.UUID   `json:"companyBranchID" gorm:"type:varchar(36);FOREIGNKEY:CompanyBranchID"`

	// Other fields.
	CompanyName    string  `json:"companyName"`
	Website        *string `json:"website"`
	Email          *string `json:"email"`
	Code           string  `json:"code"`
	HRName         *string `json:"hrName"`
	HRContact      *string `json:"hrContact"`
	FounderName    *string `json:"founderName"`
	Vacancy        *uint16 `json:"vacancy"`
	EnquiryDate    *string `json:"enquiryDate"`
	EnquiryType    *string `json:"enquiryType"`
	EnquirySource  *string `json:"enquirySource"`
	JobRole        *string `json:"jobRole,"`
	PackageOffered *uint64 `json:"packageOffered"`
	Subject        string  `json:"subject"`
	Message        string  `json:"message"`
}

// TableName will name the table of enquiry model as "company_enquiries".
func (*EnquiryDTO) TableName() string {
	return "company_enquiries"
}
