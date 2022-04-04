package admin

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// TargetCommunity consists of estimated targets on different communities.
type TargetCommunity struct {
	general.TenantBase

	// Maps.
	Colleges  []list.Branch        `json:"colleges" gorm:"many2many:target_communities_colleges;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Companies []list.CompanyBranch `json:"companies" gorm:"many2many:target_communities_companies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Courses   []list.Course        `json:"courses" gorm:"many2many:target_communities_courses;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	DepartmentID  uuid.UUID `json:"departmentID" gorm:"type:varchar(36)"`
	SalesPersonID uuid.UUID `json:"salesPersonID" gorm:"type:varchar(36)"`
	FunctionID    uuid.UUID `json:"functionID" gorm:"type:varchar(36)"`
	FacultyID     uuid.UUID `json:"facultyID" gorm:"type:varchar(36)"`

	// Other fields.
	TargetType           string   `json:"targetType" gorm:"type:varchar(50)"`
	StudentType          string   `json:"studentType" example:"Ravi" gorm:"type:varchar(50)"`
	NumberOfBatches      *uint8   `json:"numberOfBatches" gorm:"type:INT"`
	TargetStartDate      string   `json:"targetStartDate" gorm:"type:date"`
	TargetEndDate        string   `json:"targetEndDate" gorm:"type:date"`
	IsTargetAchieved     bool     `json:"isTargetAchieved" gorm:"type:tinyint(1)"`
	TargetStudentCount   *uint32  `json:"targetStudentCount" gorm:"type:INT"`
	Hours                *uint16  `json:"hours" gorm:"type:INT"`
	Fees                 *float32 `json:"fees" gorm:"type:decimal(9,2)"`
	Rating               *uint8   `json:"rating" gorm:"type:smallint(2)"`
	RequiredTalentRating uint8    `json:"requiredTalentRating" gorm:"type:smallint(2)"`
	TalentType           uint8    `json:"talentType" example:"4" gorm:"type:tinyint(2)"`
	MinExperienceYears   *uint8   `json:"minExperienceYears" example:"4" gorm:"type:tinyint(2)"`
	MaxExperienceYears   *uint8   `json:"maxExperienceYears" example:"4" gorm:"type:tinyint(2)"`
	Salary               *float32 `json:"salary" gorm:"type:decimal(12,2)"`
	UpSell               *uint16  `json:"upSell" gorm:"type:SMALLINT(4)"`
	CrossSell            *uint16  `json:"crossSell" gorm:"type:SMALLINT(4)"`
	Referral             *uint16  `json:"referral" gorm:"type:SMALLINT(4)"`
	Action               *string  `json:"action" gorm:"type:varchar(1000)"`
}

// Validate will check if all fields are valid in Target Community.
func (community *TargetCommunity) Validate() error {

	// Department id.
	if !util.IsUUIDValid(community.DepartmentID) {
		return errors.NewValidationError("Department ID must be specified")
	}

	// Function id.
	if !util.IsUUIDValid(community.FunctionID) {
		return errors.NewValidationError("Function ID must be specified")
	}

	// SalesPserson id.
	if !util.IsUUIDValid(community.SalesPersonID) {
		return errors.NewValidationError("SalesPerson ID must be specified")
	}

	// Faculty id.
	if !util.IsUUIDValid(community.FacultyID) {
		return errors.NewValidationError("Faculty ID must be specified")
	}

	// Validate target type.
	if util.IsEmpty(community.TargetType) {
		return errors.NewValidationError("Target type must be specified")
	}

	// Validate student type.
	if util.IsEmpty(community.StudentType) {
		return errors.NewValidationError("Student type must be specified")
	}

	// Validate number of batches.
	if community.NumberOfBatches != nil {
		if *community.NumberOfBatches < 0 {
			return errors.NewValidationError("Number of batches cannot be less than 0")
		}
		if *community.NumberOfBatches > 255 {
			return errors.NewValidationError("Number of batches cannot be more than 255")
		}
	}

	// Validate target student count.
	if community.TargetStudentCount != nil {
		if *community.TargetStudentCount < 0 {
			return errors.NewValidationError("Target Student Count cannot be less than 0")
		}
		if *community.TargetStudentCount > 2147483647 {
			return errors.NewValidationError("Target Student Count cannot be more than 2147483647")
		}
	}

	// Validate hours.
	if community.Hours != nil {
		if *community.Hours < 0 {
			return errors.NewValidationError("Hours cannot be less than 0")
		}
		if *community.Hours > 65535 {
			return errors.NewValidationError("Hours cannot be more than 65535")
		}
	}

	// Validate fees.
	if community.Fees != nil {
		if *community.Fees < 0 {
			return errors.NewValidationError("Fees cannot be less than 0")
		}
		if *community.Fees > 1000000 {
			return errors.NewValidationError("Fees cannot be more than 1000000")
		}
	}

	// Validate target start date.
	if util.IsEmpty(community.TargetStartDate) {
		return errors.NewValidationError("Target start date must be specified")
	}

	// Validate target end date.
	if util.IsEmpty(community.TargetEndDate) {
		return errors.NewValidationError("Target end date must be specified")
	}

	// Validate talent type.
	if community.TalentType > 99 {
		return errors.NewValidationError("Talent type cannot be more than 99")
	}

	if community.TalentType < 0 {
		return errors.NewValidationError("Talent type cannot be less than 0")
	}

	// Validate min experience years.
	if community.MinExperienceYears != nil && *community.MinExperienceYears > 99 {
		return errors.NewValidationError("Minimum experience years cannot be more than 99")
	}

	if community.MinExperienceYears != nil && *community.MinExperienceYears < 0 {
		return errors.NewValidationError("Minimum experience years cannot be less than 0")
	}

	// Validate max experience years.
	if community.MaxExperienceYears != nil && *community.MaxExperienceYears > 99 {
		return errors.NewValidationError("Maximum experience years cannot be more than 99")
	}

	if community.MaxExperienceYears != nil && *community.MaxExperienceYears < 0 {
		return errors.NewValidationError("Maximum experience years cannot be less than 0")
	}

	// Validate salary.
	if community.Salary != nil && *community.Salary > 9999999999 {
		return errors.NewValidationError("Salary cannot be more than 9999999999")
	}

	if community.Salary != nil && *community.Salary < 0 {
		return errors.NewValidationError("Salary cannot be less than 0")
	}

	// Validate upsell.
	if community.UpSell != nil && *community.UpSell > 1000 {
		return errors.NewValidationError("UpSell cannot be more than 1000")
	}

	if community.UpSell != nil && *community.UpSell < 0 {
		return errors.NewValidationError("UpSell cannot be less than 0")
	}

	// Validate crossSell.
	if community.CrossSell != nil && *community.CrossSell > 1000 {
		return errors.NewValidationError("CrossSell cannot be more than 1000")
	}

	if community.CrossSell != nil && *community.CrossSell < 0 {
		return errors.NewValidationError("CrossSell cannot be less than 0")
	}

	// Validate referral.
	if community.Referral != nil && *community.Referral > 1000 {
		return errors.NewValidationError("Referral cannot be more than 1000")
	}

	if community.Referral != nil && *community.Referral < 0 {
		return errors.NewValidationError("Referral cannot be less than 0")
	}

	// Validate action.
	if community.Action != nil && len(*community.Action) > 1000 {
		return errors.NewValidationError("Action can have maximum 1000 characters")
	}

	// Validate required talent rating.
	if community.RequiredTalentRating == 0 {
		return errors.NewValidationError("Required Talent Rating must be specified")
	}

	// Check if rating is present or not if target type is college or company.
	if (community.TargetType == "Company" || community.TargetType == "College") && community.Rating == nil {
		return errors.NewValidationError("Rating must be specified")
	}

	// Check if rating is below 0 or not if target type is college or company.
	if (community.TargetType == "Company" || community.TargetType == "College") && community.Rating == nil && *community.Rating < 0 {
		return errors.NewValidationError("Rating cannot be below 0")
	}

	// Check if rating is above 10 if target type is college or company.
	if (community.TargetType == "Company" || community.TargetType == "College") && community.Rating == nil && *community.Rating < 0 {
		return errors.NewValidationError("Rating cannot be above 10")
	}

	return nil
}

// TargetCommunityUpdate is the model for updating some fields of target community.
type TargetCommunityUpdate struct {
	TargetCommunityID uuid.UUID `gorm:"type:varchar(36)"`
}

// ===========Defining many to many structs===========

// TargetCommunityColleges is the map of target community and college.
type TargetCommunityColleges struct {
	TargetCommunityID uuid.UUID `gorm:"type:varchar(36)"`
	BranchID          uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*TargetCommunityColleges) TableName() string {
	return "target_communities_colleges"
}

// TargetCommunityCompanies is the map of target community and companies.
type TargetCommunityCompanies struct {
	TargetCommunityID uuid.UUID `gorm:"type:varchar(36)"`
	CompanyBranchID   uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*TargetCommunityCompanies) TableName() string {
	return "target_communities_companies"
}

// TargetCommunityCourses is the map of target community and courses.
type TargetCommunityCourses struct {
	TargetCommunityID uuid.UUID `gorm:"type:varchar(36)"`
	CourseID          uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*TargetCommunityCourses) TableName() string {
	return "target_communities_courses"
}

//************************************* DTO MODEL *************************************************************

// TargetCommunityDTO consists of estimated targets on different communities.
type TargetCommunityDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	Colleges  []list.Branch        `json:"colleges" gorm:"many2many:target_communities_colleges;association_jointable_foreignkey:branch_id;jointable_foreignkey:target_community_id"`
	Companies []list.CompanyBranch `json:"companies" gorm:"many2many:target_communities_companies;association_jointable_foreignkey:company_branch_id;jointable_foreignkey:target_community_id"`
	Courses   []list.Course        `json:"courses" gorm:"many2many:target_communities_courses;association_jointable_foreignkey:course_id;jointable_foreignkey:target_community_id"`

	// Related tables.
	Department    general.DepartmentDTO           `json:"department" gorm:"foreignkey:DepartmentID"`
	DepartmentID  uuid.UUID                       `json:"-"`
	SalesPerson   general.User                    `json:"salesPerson" gorm:"foreignkey:SalesPersonID"`
	SalesPersonID uuid.UUID                       `json:"-"`
	Function      general.TargetCommunityFunction `json:"function" gorm:"foreignkey:FunctionID"`
	FunctionID    uuid.UUID                       `json:"-"`
	Faculty       list.Faculty                    `json:"faculty" gorm:"foreignkey:FacultyID"`
	FacultyID     uuid.UUID                       `json:"-"`

	// Other fields.
	TargetType           string   `json:"targetType"`
	StudentType          string   `json:"studentType"`
	NumberOfBatches      *uint8   `json:"numberOfBatches"`
	TargetStartDate      string   `json:"targetStartDate"`
	TargetEndDate        string   `json:"targetEndDate"`
	IsTargetAchieved     bool     `json:"isTargetAchieved"`
	TargetStudentCount   *uint32  `json:"targetStudentCount"`
	Hours                *uint16  `json:"hours"`
	Fees                 *float32 `json:"fees"`
	Rating               *uint8   `json:"rating"`
	RequiredTalentRating uint8    `json:"requiredTalentRating"`
	TalentType           uint8    `json:"talentType"`
	MinExperienceYears   *uint8   `json:"minExperienceYears"`
	MaxExperienceYears   *uint8   `json:"maxExperienceYears"`
	Salary               *float32 `json:"salary"`
	UpSell               *uint16  `json:"upSell"`
	CrossSell            *uint16  `json:"crossSell"`
	Referral             *uint16  `json:"referral"`
	Action               *string  `json:"action"`
}

// TableName will name the table of Experience model as "target_communities".
func (*TargetCommunityDTO) TableName() string {
	return "target_communities"
}
