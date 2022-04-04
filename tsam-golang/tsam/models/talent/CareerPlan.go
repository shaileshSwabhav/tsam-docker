package talent

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// CareerPlan consists of fields of talent's career plan.
type CareerPlan struct {
	general.TenantBase
	CareerObjectiveID         uuid.UUID `json:"careerObjectiveID" gorm:"type:varchar(36)"`
	CareerObjectivesCoursesID uuid.UUID `json:"careerObjectivesCoursesID" gorm:"type:varchar(36)"`
	FacultyID                 uuid.UUID `json:"facultyID" gorm:"type:varchar(36)"`
	TalentID                  uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	CurrentRating             uint8     `json:"currentRating" gorm:"type:TINYINT(2)"`
}

// TableName will name the table of CareerPlan model as "talent_career_plans".
func (*CareerPlan) TableName() string {
	return "talent_career_plans"
}

// Validate Validates fields of career plan.
func (careerPlan *CareerPlan) Validate() error {
	// Check if current rating is specified or not.
	if careerPlan.CurrentRating == 0 {
		return errors.NewValidationError("Current Rating must be specified")
	}

	// Check if current rating is below 1 or not.
	if careerPlan.CurrentRating < 1 {
		return errors.NewValidationError("Current Rating must be above 0")
	}

	// Check if current rating is above 10 or not.
	if careerPlan.CurrentRating > 10 {
		return errors.NewValidationError("Current Rating must be below 10")
	}

	// Check if CareerObjectiveID is specified or not.
	if !util.IsUUIDValid(careerPlan.CareerObjectiveID) {
		return errors.NewValidationError("Career Objective ID must be specified")
	}

	// Check if CareerObjectivesCoursesID is specified or not.
	if !util.IsUUIDValid(careerPlan.CareerObjectivesCoursesID) {
		return errors.NewValidationError("Career Objectives Courses ID must be specified")
	}

	// Check if FacultyID is specified or not.
	if !util.IsUUIDValid(careerPlan.FacultyID) {
		return errors.NewValidationError("Faculty ID must be specified")
	}

	return nil
}
