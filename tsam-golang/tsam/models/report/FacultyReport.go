package report

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
)

type FacultyReport struct {
	Faculty            list.Faculty               `json:"faculty"`
	Monday             []Batch                    `json:"monday"`
	Tuesday            []Batch                    `json:"tuesday"`
	Wednesday          []Batch                    `json:"wednesday"`
	Thursday           []Batch                    `json:"thursday"`
	Friday             []Batch                    `json:"friday"`
	Saturday           []Batch                    `json:"saturday"`
	Sunday             []Batch                    `json:"sunday"`
	WorkingHours       map[uuid.UUID]WorkingHours `json:"workingHours"`
	TotalTrainingHours float64                    `json:"totalTrainingHours"`
}

// Batch is used for listing of batches
type Batch struct {
	general.BaseDTO
	BatchName       string              `json:"batchName"`
	Status          *string             `json:"batchStatus" `
	TotalDailyHours float64             `json:"totalDailyHours"`
	BatchTimings    []BatchModuleTiming `json:"batchTimings"`
}

// WorkingHours will store batch and and total hours and minutes for the batch in a week.
type WorkingHours struct {
	BatchName  string  `json:"batchName"`
	TotalHours float64 `json:"totalHours"`
}

// type BatchModule struct {
// 	general.BaseDTO
// 	BatchID            uuid.UUID           `json:"-"`
// 	ModuleID           uuid.UUID           `json:"-"`
// 	BatchModuleTimings []batchModuleTiming `json:"batchTimings"`
// }

// BatchModuleTiming will store the schedule for sessions of batch
type BatchModuleTiming struct {
	general.BaseDTO
	BatchID       uuid.UUID   `json:"-"`
	BatchModuleID uuid.UUID   `json:"batchModuleID"`
	Day           general.Day `json:"day" gorm:"type:foreignkey:DayID"`
	DayID         uuid.UUID   `json:"-"`
	FromTime      *string     `json:"fromTime"`
	ToTime        *string     `json:"toTime"`
}

// TableName defines table name of the struct.
func (*BatchModuleTiming) TableName() string {
	return "batch_module_timings"
}
