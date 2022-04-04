package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//Batch contain information about single batch
type Batch struct {
	general.TenantBase
	BatchName      string            `json:"batchName" gorm:"type:varchar(50);not null"`
	Code           string            `json:"code" gorm:"type:varchar(10);not null"`
	StartDate      *string           `json:"startDate" gorm:"type:date"`
	TotalStudents  *uint8            `json:"totalStudents" gorm:"type:int(3)"`
	TotalIntake    *uint8            `json:"totalIntake" gorm:"type:int(3)"`
	Status         string            `json:"batchStatus" gorm:"type:varchar(20)"`
	IsActive       *bool             `json:"isActive" gorm:"DEFAULT:true"`
	Eligibility    *Eligibility      `json:"eligibility" gorm:"foreignkey:EligibilityID"`
	Course         *general.Course   `json:"course" gorm:"foreignkey:CourseID;association_autocreate:false;association_autoupdate:false"`
	SalesPerson    *list.User        `json:"salesPerson" gorm:"foreignkey:SalesPersonID;association_autocreate:false;association_autoupdate:false"`
	Requirement    *list.Requirement `json:"requirement" gorm:"foreignkey:RequirementID;association_autocreate:false;association_autoupdate:false"`
	SalesPersonID  uuid.UUID         `json:"-" gorm:"type:varchar(36)"`
	CourseID       uuid.UUID         `json:"-" gorm:"type:varchar(36)"`
	EligibilityID  *uuid.UUID        `json:"-" gorm:"type:varchar(36)"`
	RequirementID  *uuid.UUID        `json:"-" gorm:"type:varchar(36)"`
	IsB2B          *bool             `json:"isB2B" gorm:"type:tinyint(1);column:is_b2b"`
	BatchObjective string            `json:"batchObjective" gorm:"type:varchar(100)"`
	Brochure       *string           `json:"brochure" gorm:"type:varchar(200)"`
	MeetLink       *string           `json:"meetLink" gorm:"type:varchar(100)"`
	TelegramLink   *string           `json:"telegramLink" gorm:"type:varchar(100)"`

	EstimatedEndDate *string `json:"estimatedEndDate" gorm:"type:date"`
	FinalEndDate     *string `json:"finalEndDate" gorm:"type:date"`

	Timing []Timing `json:"batchTimings"`
	// Faculty        *list.Faculty     `json:"faculty" gorm:"foreignkey:FacultyID;association_autocreate:false;association_autoupdate:false"`
	// FacultyID      uuid.UUID         `json:"-" gorm:"type:varchar(36)"`
	// Logo           *string    `json:"logo" gorm:"type:varchar(200)"`
	// EndDate        *string           `json:"endDate" gorm:"type:date"`
}

func (batch *Batch) ValidateBatch() error {

	if util.IsEmpty(batch.BatchObjective) {
		return errors.NewValidationError("batch objective must be specified")
	}

	if util.IsEmpty(batch.BatchName) {
		return errors.NewValidationError("batch name must be specified")
	}

	if batch.Course == nil {
		return errors.NewValidationError("course must be specified")
	}

	if batch.TotalIntake != nil && *batch.TotalIntake <= 0 {
		return errors.NewValidationError("total intake must be greater than 0")
	}

	if batch.SalesPerson == nil {
		return errors.NewValidationError("salesperson must be specified")
	}

	// if batch.Timing != nil {
	// 	for _, batchTime := range batch.Timing {
	// 		err := batchTime.ValidateBatchTiming()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

// BatchDTO contains information about batch
type BatchDTO struct {
	ID                    uuid.UUID         `json:"id"`
	DeletedAt             *time.Time        `json:"-"`
	BatchName             string            `json:"batchName" `
	Code                  string            `json:"code" `
	StartDate             *string           `json:"startDate" `
	EndDate               *string           `json:"endDate" `
	TotalStudents         *uint8            `json:"totalStudents"`
	TotalIntake           *uint8            `json:"totalIntake" `
	Status                *string           `json:"batchStatus" `
	IsActive              *bool             `json:"isActive"`
	SalesPersonID         uuid.UUID         `json:"-"`
	CourseID              uuid.UUID         `json:"-"`
	EligibilityID         *uuid.UUID        `json:"-"`
	RequirementID         *uuid.UUID        `json:"-"`
	SalesPerson           *list.User        `json:"salesPerson"`
	Course                *general.Course   `json:"course"`
	Eligibility           *Eligibility      `json:"eligibility"`
	Requirement           *list.Requirement `json:"requirement"`
	IsB2B                 *bool             `json:"isB2B" gorm:"column:is_b2b"`
	Brochure              *string           `json:"brochure"`
	Logo                  *string           `json:"logo"`
	BatchObjective        string            `json:"batchObjective"`
	TotalSessionCount     uint              `json:"totalSessionCount"`
	CompletedSessionCount uint              `json:"completedSessionCount"`
	TotalApplicants       *uint16           `json:"totalApplicants"`
	Timing                []*Timing         `json:"batchTimings" gorm:"foreignkey:BatchID"`
	Faculty               []*list.Faculty   `json:"faculty"`
	// FacultyID             uuid.UUID         `json:"-"`

	EstimatedEndDate *string `json:"estimatedEndDate" gorm:"type:date"`
	FinalEndDate     *string `json:"finalEndDate" gorm:"type:date"`
}

// TableName defines table name of the struct.
func (*BatchDTO) TableName() string {
	return "batches"
}

//*************************************** UPCOMING BATCH ******************************************************

// UpcomingBatch contains few fields related to upcoming batches.
type UpcomingBatch struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Batch related fields.
	TotalEnrolled uint   `json:"totalEnrolled"`
	StartDate     string `json:"startDate"`
	TotalIntake   uint8  `json:"totalIntake"`
	BatchName     string `json:"batchName"`

	// Course realted fields.
	Course   MyUpcomingBatchCourse `json:"course" gorm:"foreignkey:CourseID"`
	CourseID uuid.UUID             `json:"-"`
}

// TableName defines table name of the struct.
func (*UpcomingBatch) TableName() string {
	return "batches"
}

// MyUpcomingBatchCourse contains course details for upcoming batches.
type MyUpcomingBatchCourse struct {
	ID          uuid.UUID  `json:"id"`
	DeletedAt   *time.Time `json:"-"`
	Name        string     `json:"name"`
	Logo        *string    `json:"logo"`
	Description string     `json:"description"`
}

// TableName defines table name of the struct.
func (*MyUpcomingBatchCourse) TableName() string {
	return "courses"
}

//*************************************** BATCH-DETAILS ******************************************************

// BatchDetails contains few fields related to one batch.
type BatchDetails struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Batch related fields.
	BatchName string  `json:"batchName"`
	StartDate string  `json:"startDate"`
	EndDate   *string `json:"endDate"`
	Location  *string `json:"location"`
	Brochure  *string `json:"brochure"`
	// Faculty            *FacultyDetails `json:"faculty" gorm:"foreignkey:FacultyID"`
	Faculty            []FacultyDetails `json:"faculty"`
	FacultyID          *uuid.UUID       `json:"-"`
	TotalIntake        uint8            `json:"totalIntake"`
	TotalSessionsCount uint             `json:"totalSessionsCount"`
	BatchMeetLink      *string          `json:"batchMeetLink"`
	BatchTelegramLink  *string          `json:"batchTelegramLink"`
	TotalHours         float32          `json:"totalHours"`
	BatchTimings       []Timing         `json:"batchTimings" gorm:"foreignkey:BatchID"`
	BatchStatus        *string          `json:"batchStatus"`
	TotalSessionsCompleted uint        `json:"totalSessionsCompleted"`
	TotalCompletedHours    float32     `json:"totalCompletedHours"`
	TotalStudents         *uint8            `json:"totalStudents"`

	EstimatedEndDate *string `json:"estimatedEndDate" gorm:"type:date"`
	FinalEndDate     *string `json:"finalEndDate" gorm:"type:date"`

	// Course realted fields.
	Course   MyUpcomingBatchCourse `json:"course" gorm:"foreignkey:CourseID"`
	CourseID uuid.UUID             `json:"-"`
}

// TableName defines table name of the struct.
func (*BatchDetails) TableName() string {
	return "batches"
}

// Faculty is DTO which is used for getting minimum details of faculty.
type FacultyDetails struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	TelegramID *string   `json:"telegramID"`
}

// TableName defines table name of the struct.
func (*FacultyDetails) TableName() string {
	return "faculties"
}

// BatchDetailsSessionDTO will get batch session realeted details for batch details.
type BatchDetailsSessionDTO struct {
	general.BaseDTO
	FacultyID *uuid.UUID      `json:"-"`
	Faculty   *FacultyDetails `json:"faculty"`
}

// TableName overrides name of the table
func (*BatchDetailsSessionDTO) TableName() string {
	return "batch_sessions"
}
