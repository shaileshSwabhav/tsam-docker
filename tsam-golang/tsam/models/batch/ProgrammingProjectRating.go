package batch

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// ProgrammingProjectRatings will store map for project rating
type ProgrammingProjectRating struct {
	general.TenantBase
	TalentID                            uuid.UUID                          `json:"talentID" gorm:"type:varchar(36)"`
	BatchID                             uuid.UUID                          `json:"batchID" gorm:"type:varchar(36)"`
	TalentSubmissionID                  uuid.UUID                          `json:"talentSubmissionID" gorm:"type:varchar(36)"`
	ProgrammingProjectRatingParameterID uuid.UUID                          `json:"programmingProjectRatingParameterID" gorm:"type:varchar(36)"`
	ProgrammingProjectRatingParameter   *ProgrammingProjectRatingParameter `json:"programmingProjectRatingParameter"`
	Score                               uint                               `json:"score" gorm:"type:int"`
}

// TableName overrides name of the table
func (*ProgrammingProjectRating) TableName() string {
	return "programming_project_ratings"
}

// Validate will validate all the fields of programming_project_ratings.
func (programmingprojectRating *ProgrammingProjectRating) Validate() error {

	if !util.IsUUIDValid(programmingprojectRating.BatchID) {
		return errors.NewValidationError("Batch ID must be specified.")
	}

	if !util.IsUUIDValid(programmingprojectRating.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}

	if !util.IsUUIDValid(programmingprojectRating.ProgrammingProjectRatingParameterID) {
		return errors.NewValidationError("Programming Project Rating Parameter ID must be specified.")
	}

	if !util.IsUUIDValid(programmingprojectRating.TalentSubmissionID) {
		return errors.NewValidationError("Talent submission ID must be specified.")
	}

	// if util.IsEmpty(assignment.DueDate) {
	// 	return errors.NewValidationError("Due date must be specified.")
	// }

	// if assignment.BatchAssignmentTask != nil {
	// 	if err := assignment.BatchAssignmentTask.Validate(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type ProgrammingProjectRatingDTO struct {
	ID                                  uuid.UUID                          `json:"id"`
	DeletedAt                           *time.Time                         `json:"-"`
	TalentID                            uuid.UUID                          `json:"talentID"`
	BatchID                             uuid.UUID                          `json:"-"`
	TalentSubmissionID                  uuid.UUID                          `json:"-"`
	ProgrammingProjectRatingParameterID uuid.UUID                          `json:"-"`
	ProgrammingProjectRatingParameter   *ProgrammingProjectRatingParameter `json:"-"`
	Score                               uint                               `json:"score"`
}
