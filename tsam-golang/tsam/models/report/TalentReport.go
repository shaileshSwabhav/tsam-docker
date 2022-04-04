package report

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/list"
)

type TalentReport struct {
	Talent             list.Talent                `json:"talent"`
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
