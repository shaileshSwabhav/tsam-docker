package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Score contains details about the score the talent or enquiry has obtained for the examination.
type Score struct {
	TenantBase

	// Related table IDs.
	ExaminationID   uuid.UUID `json:"examinationID" gorm:"type:varchar(36)"`
	MastersAbroadID uuid.UUID `json:"mastersAbroadID" gorm:"type:varchar(36)"`

	// Other fields.
	MarksObtained float64 `json:"marksObtained" gorm:"type:decimal(10, 2)"`
}

// TableName will name the table of Score model as "scores".
func (*Score) TableName() string {
	return "scores"
}

// Validate score fields.
func (score *Score) Validate() error {

	// Examination ID must be secified.
	if !util.IsUUIDValid(score.ExaminationID) {
		return errors.NewValidationError("Examination ID must be secified")
	}

	// if score.MarksObtained < 0 {
	// 	errorString := score.Examination.Name + " score cannot be below 0"
	// 	return errors.NewValidationError(errorString)
	// }
	// if score.MarksObtained > score.Examination.TotalMarks {
	// 	errorString := score.Examination.Name + " score cannot be above " + strconv.FormatFloat(score.Examination.TotalMarks, 'f', 6, 64)
	// 	return errors.NewValidationError(errorString)
	// }
	return nil
}

//************************************* DTO MODEL *************************************************************

// Score contains details about the score the talent or enquiry has obtained for the examination.
type ScoreDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Related tables and IDs.
	Examination     Examination `json:"examination" gorm:"foreignkey:ExaminationID"`
	ExaminationID   uuid.UUID   `json:"-"`
	MastersAbroadID uuid.UUID   `json:"mastersAbroadID"`

	// Other fields.
	MarksObtained float64 `json:"marksObtained"`
	// MarksObtained   uint16       `json:"marksObtained" gorm:"type:smallint"` -> ielts score has decimal
}

// TableName will name the table of Score model as "scores".
func (*ScoreDTO) TableName() string {
	return "scores"
}
