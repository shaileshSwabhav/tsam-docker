package talent

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/programming"
	"github.com/techlabs/swabhav/tsam/util"
)

// TalentConceptRating will store score for programming concept.
type TalentConceptRating struct {
	general.TenantBase
	ModuleProgrammingConceptID uuid.UUID `json:"programmingConceptModuleID" gorm:"type:varchar(36)"`
	TalentSubmissionID         uuid.UUID `json:"talentSubmissionID" gorm:"type:varchar(36)"`
	TalentID                   uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	Score                      uint      `json:"score" gorm:"type:INT"`
}

type TalentConceptRatingDTO struct {
	ModuleProgrammingConcept   programming.ModuleProgrammingConcepts `json:"programmingConceptModule" gorm:"foreignkey:ModuleProgrammingConceptID"`
	ID                         uuid.UUID                             `json:"id"`
	DeletedAt                  *time.Time                            `json:"-"`
	ModuleProgrammingConceptID uuid.UUID                             `json:"-"`
	TalentID                   uuid.UUID                             `json:"talentID"`
	TalentSubmissionID         uuid.UUID                             `json:"-"`
	Score                      uint                                  `json:"score"`
}

func (tcr *TalentConceptRating) Validate() error {

	if !util.IsUUIDValid(tcr.TalentID) {
		return errors.NewValidationError("Talent ID must be specified.")
	}
	if !util.IsUUIDValid(tcr.ModuleProgrammingConceptID) {
		return errors.NewValidationError("Concept moudle ID must be specified.")
	}
	if !util.IsUUIDValid(tcr.TalentSubmissionID) {
		return errors.NewValidationError("Talent submission ID must be specified.")
	}
	if tcr.Score <= 0 {
		return errors.NewValidationError("Score must be specified.")
	}
	return nil
}

// TableName defines table name of the struct.
func (*TalentConceptRatingDTO) TableName() string {
	return "talent_concept_ratings"
}
