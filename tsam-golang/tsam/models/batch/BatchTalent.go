package batch

import (
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
)

// MappedTalent stores mapping of batches and talents along with date of joining.
type MappedTalent struct {
	general.TenantBase
	BatchID        uuid.UUID `json:"batchID" gorm:"type:varchar(36)"`
	TalentID       uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	DateOfJoining  *string   `json:"dateOfJoining" gorm:"type:date"`
	IsActive       bool      `json:"isActive" example:"true"`
	SuspensionDate *string   `json:"suspensionDate" gorm:"type:date"`
}

// ValidateMappedTalent validates fields of MappedTalent struct
func (talent *MappedTalent) ValidateMappedTalent() error {
	// datePattern := regexp.MustCompile("^\\d{4}\\/(0[1-9]|1[012])\\/(0[1-9]|[12][0-9]|3[01])$")
	datePattern := regexp.MustCompile(`^(19|20)\d\d[-/.](0[1-9]|1[012])[-/.](0[1-9]|[12][0-9]|3[01])$`)
	// accpets date in yyyy-mm-dd pattern
	if talent.DateOfJoining != nil && !datePattern.MatchString(*talent.DateOfJoining) {
		return errors.NewValidationError("Date should be of the form yyyy/mm/dd")
	}
	return nil
}

// TableName overrides name of the table
func (*MappedTalent) TableName() string {
	return "batch_talents"
}

// BatchTalentDTO stores mapping of batches and talents along with date of joining.
type BatchTalentDTO struct {
	ID                     uuid.UUID   `json:"id"`
	DeletedAt              *time.Time  `json:"-"`
	BatchID                uuid.UUID   `json:"-"`
	TalentID               uuid.UUID   `json:"-"`
	Batch                  list.Batch  `json:"batch" gorm:"foreignkey:BatchID"`
	Talent                 list.Talent `json:"talent" gorm:"foreignkey:TalentID"`
	DateOfJoining          *string     `json:"dateOfJoining" gorm:"type:date"`
	SessionsAttended       uint        `json:"sessionsAttendedCount"`
	TotalSessionsCount     uint        `json:"totalSessionsCount"`
	TotalSessionsCompleted uint        `json:"totalSessionsCompleted"`
	TotalHours             float32     `json:"totalHours"`
	TotalCompletedHours    float32     `json:"totalCompletedHours"`
	AttendedHours          float32     `json:"attendedHours"`
	AverageRating          float32     `json:"averageRating"`
	TotalFeedbacksGiven    uint        `json:"totalFeedbacksGiven"`
	IsActive               *bool       `json:"isActive"`
	SuspensionDate         *string     `json:"suspensionDate"`
	// BatchStartTime         *string     `json:"batchStartTime"`
	// BatchEndTime           *string     `json:"batchEndTime"`
	BatchMeetLink          *string     `json:"batchMeetLink"`
	BatchTelegramLink          *string     `json:"batchTelegramLink"`
	// FacultyTelegram        *string     `json:"facultyTelegram"`
	BatchTimings           []Timing    `json:"batchTimings"`
	StartDate              string  `json:"startDate"`
	EstimatedEndDate       *string `json:"estimatedEndDate" gorm:"type:date"`
}

// TableName overrides name of the table
func (*BatchTalentDTO) TableName() string {
	return "batch_talents"
}

// UpdateBatchTalentSuspension is used for updating suspension date of batch talent.
type UpdateBatchTalentSuspension struct {
	ID        uuid.UUID `json:"id"`
	IsSuspend bool      `json:"isSuspend"`
}

// UpdateBatchTalentIsActive is used for updating is active of batch talent.
type UpdateBatchTalentIsActive struct {
	ID       uuid.UUID `json:"id"`
	IsActive bool      `json:"isActive"`
	BatchID  uuid.UUID `json:"batchID"`
}

//************************************* MINIMUM DTO MODEL *********************************************************

// MinimumBatchTalentForTalent is used for getting minimum details of batch talent.
type MinimumBatchTalentForTalent struct {
	ID          uuid.UUID     `json:"id"`
	DeletedAt   *time.Time    `json:"-"`
	BatchID     uuid.UUID     `json:"batchID"`
	CourseName  string        `json:"courseName"`
	BatchStatus string        `json:"batchStatus"`
	Faculty     *list.Faculty `json:"faculty" gorm:"foreignkey:FacultyID"`
	FacultyID   *uuid.UUID    `json:"-"`
}

// TableName overrides name of the table
func (*MinimumBatchTalentForTalent) TableName() string {
	return "batch_talents"
}
