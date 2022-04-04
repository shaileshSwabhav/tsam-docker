package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// WaitingList is used for adding and updating waiting list, consists of IDs of related information.
type WaitingList struct {
	general.TenantBase

	// Talent/Enquiry related fields.
	TalentID  *uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	EnquiryID *uuid.UUID `json:"enquiryID" gorm:"type:varchar(36)"`
	IsActive  bool       `json:"isActive"`
	Email     *string    `json:"email" gorm:"type:varchar(100)"`

	// Company related IDs.
	CompanyBranchID      *uuid.UUID `json:"companyBranchID" gorm:"type:varchar(36)"`
	CompanyRequirementID *uuid.UUID `json:"companyRequirementID" gorm:"type:varchar(36)"`

	// Course related IDs.
	CourseID *uuid.UUID `json:"courseID" gorm:"type:varchar(36)"`
	BatchID  *uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`

	// Other fields.
	SourceID uuid.UUID `json:"sourceID" gorm:"type:varchar(36)"`
}

// TwoWaitingLists is for getting ywo waiting lists related to talent and enquiry.
type TwoWaitingLists struct {
	TalentWaitingList  *[]WaitingList `json:"talentWaitingList"`
	EnquiryWaitingList *[]WaitingList `json:"enquiryWaitingList"`
}

// UpdateWaitingList is for updating a batch of waiting list.
type UpdateWaitingList struct {
	CompanyBranchID *uuid.UUID    `json:"companyBranchID"`
	RequirementID   *uuid.UUID    `json:"requirementID"`
	CourseID        *uuid.UUID    `json:"courseID"`
	BatchID         *uuid.UUID    `json:"batchID"`
	WaitingLists    []WaitingList `json:"waitingLists"`
}

// TableName defines table name of the struct.
func (*WaitingList) TableName() string {
	return "waiting_list"
}

// Validate Validates compulsary fields of waiting list.
func (waitingList *WaitingList) Validate() error {

	// Email should be in proper format.
	if waitingList.Email != nil && !util.ValidateEmail(*waitingList.Email) {
		return errors.NewValidationError("Email must in the format : email@example.com")
	}

	// Email maximum characters.
	if waitingList.Email != nil && len(*waitingList.Email) > 100 {
		return errors.NewValidationError("Email can have maximum 100 characters")
	}

	// Check if talent or enquiry id is given or not.
	if waitingList.TalentID == nil && waitingList.EnquiryID == nil {
		return errors.NewValidationError("Talent ID or Enquiry must be mentioned")
	}

	// Check if source id is present or not.
	if !util.IsUUIDValid(waitingList.SourceID) {
		return errors.NewValidationError("Source ID must be mentioned")
	}

	return nil
}

//************************************* DTO MODEL *************************************************************

// WaitingListDTO is used for adding and updating waiting list, consists of IDs of related information.
type WaitingListDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Talent/Enquiry related fields.
	TalentID  *uuid.UUID `json:"talentID"`
	EnquiryID *uuid.UUID `json:"enquiryID"`
	IsActive  bool       `json:"isActive"`
	Email     *string    `json:"email"`

	// Company related IDs.
	CompanyBranchID      *uuid.UUID          `json:"-"`
	CompanyBranch        *list.CompanyBranch `json:"companyBranch" gorm:"foreignkey:CompanyBranchID"`
	CompanyRequirementID *uuid.UUID          `json:"-"`
	CompanyRequirement   *list.Requirement   `json:"companyRequirement" gorm:"foreignkey:CompanyRequirementID"`

	// Course related IDs.
	CourseID *uuid.UUID   `json:"-"`
	Course   *list.Course `json:"course" gorm:"foreignkey:CourseID"`
	BatchID  *uuid.UUID   `json:"-"`
	Batch    *list.Batch  `json:"batch" gorm:"foreignkey:BatchID"`

	// Count related.
	TotalCount *uint16 `json:"totalCount" gorm:"type:SMALLINT(4)"`

	// Other fields.
	SourceID uuid.UUID      `json:"-"`
	Source   general.Source `json:"source" gorm:"foreignkey:SourceID"`
}

// TableName defines table name of the struct.
func (*WaitingListDTO) TableName() string {
	return "waiting_list"
}
