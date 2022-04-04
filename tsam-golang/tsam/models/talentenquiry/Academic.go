package talentenquiry

import (
	"time"

	uuid "github.com/satori/go.uuid"
	// "github.com/techlabs/swabhav/tsam/errors"
	// "github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Academic contains the academic details of enquiry which are enough for adding and updating enquiry academic.
type Academic struct {
	general.TenantBase
	// Related table IDs.
	CollegeID        *uuid.UUID `json:"-" gorm:"type:varchar(36);column:college_branch_id"`
	DegreeID         uuid.UUID  `json:"degreeID" gorm:"type:varchar(36)"`
	SpecializationID uuid.UUID  `json:"specializationID" gorm:"type:varchar(36)"`
	EnquiryID        uuid.UUID  `json:"enquiryID" gorm:"type:varchar(36)"`

	// Other fields.
	College    string  `json:"college" gorm:"type:varchar(200)"`
	Percentage float32 `json:"percentage" gorm:"type:decimal(4,2)"`
	Passout    uint16  `json:"passout" gorm:"type:SMALLINT(4)"`
}

// TableName will name the table of Academic model as "talent_enquiry_academics".
func (*Academic) TableName() string {
	return "talent_enquiry_academics"
}

// Validate Validates fields of talent academic.
func (academic *Academic) Validate() error {

	// College ID.
	if academic.CollegeID != nil && !util.IsUUIDValid(*academic.CollegeID) {
		return errors.NewValidationError("College ID must ne a proper uuid")
	}

	// Degree ID.
	if !util.IsUUIDValid(academic.DegreeID) {
		return errors.NewValidationError("Degree ID must ne a proper uuid")
	}

	// Check if specialization id is specified or not.
	if !util.IsUUIDValid(academic.SpecializationID) {
		return errors.NewValidationError("Specialization ID must be specified")
	}

	// Check if college name is specified or not.
	if util.IsEmpty(academic.College) {
		return errors.NewValidationError("College Name must be specified")
	}

	// College name maximum length.
	if len(academic.College) > 200 {
		return errors.NewValidationError("College Name can have maximum 200 characters")
	}

	// Check if percentage is specified or not.
	if academic.Percentage == 0 {
		return errors.NewValidationError("Percentage must be specified")
	}

	// Check if passout year is specified or not.
	if academic.Passout == 0 {
		return errors.NewValidationError("Passout must be specified")
	}

	//*********************************CHECKS FOR NON COMPULSARY FIELDS(OLD CODE)****************************

	// // Check if specialization id is specified or not.
	// if !isEnquiryForm && (academic.SpecializationID == nil || (academic.SpecializationID != nil && !util.IsUUIDValid(*academic.SpecializationID))) {
	// 	return errors.NewValidationError("Specialization ID must be specified")
	// }

	// // Check if college name is specified or not.
	// if !isEnquiryForm && (academic.College == nil || (academic.College != nil && util.IsEmpty(*academic.College))) {
	// 	return errors.NewValidationError("College Name must be specified")
	// }

	// // College name maximum length.
	// if !isEnquiryForm && academic.College != nil && (len(*academic.College) > 200) {
	// 	return errors.NewValidationError("College Name can have maximum 200 characters")
	// }

	// // Check if percentage is specified or not.
	// if !isEnquiryForm && (academic.Percentage == nil || (academic.Percentage != nil && *academic.Percentage == 0)) {
	// 	return errors.NewValidationError("Percentage must be specified")
	// }

	return nil
}

//************************************* DTO MODEL *************************************************************

// AcademicDTO contains the complete information of enquiry academic which is needed to display.
type AcademicDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Realted tables.
	Degree           general.Degree          `json:"degree" gorm:"foreignkey:DegreeID"`
	DegreeID         uuid.UUID               `json:"-"`
	Specialization   *general.Specialization `json:"specialization" gorm:"foreignkey:SpecializationID"`
	SpecializationID *uuid.UUID              `json:"-"`
	College          *string                 `json:"college"`
	CollegeID        *uuid.UUID              `json:"-" gorm:"column:college_branch_id"`

	// Other fields.
	Percentage *float32  `json:"percentage"`
	Passout    uint16    `json:"passout"`
	EnquiryID  uuid.UUID `json:"enquiryID"`
}

// TableName will name the table of Academic model as "talent_academics".
func (*AcademicDTO) TableName() string {
	return "talent_enquiry_academics"
}

//************************************* EXCEL MODEL *************************************************************

// Academic contains the academic details of talent which are enough for adding and updating talent academic.
type AcademicExcel struct {
	CollegeName        string  `json:"collegeName"`
	DegreeName         string  `json:"degreeName"`
	SpecializationName string  `json:"specializationName"`
	Percentage         float32 `json:"percentage"`
	YearOfPassout      uint16  `json:"yearOfPassout"`
}

// Validate validates fields of talent excel's academic excel.
func (academic *AcademicExcel) Validate() error {

	// College name.
	if util.IsEmpty(academic.CollegeName) {
		return errors.NewValidationError("College name must be specified")
	}

	// Degree name.
	if util.IsEmpty(academic.CollegeName) {
		return errors.NewValidationError("Degree name must be specified")
	}

	// Specialization name.
	if util.IsEmpty(academic.SpecializationName) {
		return errors.NewValidationError("Specialization name must be specified")
	}

	// College name maximum length.
	if len(academic.CollegeName) > 200 {
		return errors.NewValidationError("College Name can have maximum 200 characters")
	}

	// Check if percentage is specified or not.
	if academic.Percentage == 0 {
		return errors.NewValidationError("Percentage must be specified")
	}

	// Check if passout year is specified or not.
	if academic.YearOfPassout == 0 {
		return errors.NewValidationError("Passout must be specified")
	}

	return nil
}
