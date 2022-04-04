package talent

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// TotalLifetimeValueResult contains information about total lifetime value of all talents in database.
type TotalLifetimeValueResult struct {
	TotalLifetimeValue int64
}

// LifetimeValue contains single talent's lifetime value details.
type LifetimeValue struct {
	general.TenantBase
	Upsell    *uint     `json:"upsell"`
	Placement *uint     `json:"placement"`
	Knowledge *uint     `json:"knowledge"`
	Teaching  *uint     `json:"teaching"`
	TalentID  uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
}

// LifetimeValueReport contains information to be displayed in lifetime value report.
type LifetimeValueReport struct {
	LifetimeValue
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// TableName will name the table of LifetimeValue model as "talent_lifetime_values".
func (*LifetimeValue) TableName() string {
	return "talent_lifetime_values"
}

// Validate Validates fields of talent lifetime value.
func (lifetimeValue *LifetimeValue) Validate() error {
	// Check if upsell value is below 0.
	if lifetimeValue.Upsell != nil && *lifetimeValue.Upsell < 0 {
		return errors.NewValidationError("Upsell value cannot be below 0")
	}

	// Check if upsell value is above 999999999999.
	if lifetimeValue.Upsell != nil && *lifetimeValue.Upsell > 999999999999 {
		return errors.NewValidationError("Upsell value cannot be above 999999999999")
	}

	// Check if placement value is below 0.
	if lifetimeValue.Placement != nil && *lifetimeValue.Placement < 0 {
		return errors.NewValidationError("Placement value cannot be below 0")
	}

	// Check if placement value is above 999999999999.
	if lifetimeValue.Placement != nil && *lifetimeValue.Placement > 999999999999 {
		return errors.NewValidationError("Placement value cannot be above 999999999999")
	}

	// Check if knowledge value is below 0.
	if lifetimeValue.Knowledge != nil && *lifetimeValue.Knowledge < 0 {
		return errors.NewValidationError("Knowledge value cannot be below 0")
	}

	// Check if knowledge value is above 999999999999.
	if lifetimeValue.Knowledge != nil && *lifetimeValue.Knowledge > 999999999999 {
		return errors.NewValidationError("Knowledge value cannot be above 999999999999")
	}

	// Check if teaching value is below 0.
	if lifetimeValue.Teaching != nil && *lifetimeValue.Teaching < 0 {
		return errors.NewValidationError("Teaching value cannot be below 0")
	}

	// Check if teaching value is above 999999999999.
	if lifetimeValue.Teaching != nil && *lifetimeValue.Teaching > 999999999999 {
		return errors.NewValidationError("Teaching value cannot be above 999999999999")
	}

	// Check if sum of all values is above 999999999999.
	if lifetimeValue.Upsell != nil && lifetimeValue.Placement != nil && lifetimeValue.Knowledge != nil && lifetimeValue.Teaching != nil &&
		(*lifetimeValue.Upsell+*lifetimeValue.Placement+*lifetimeValue.Knowledge+*lifetimeValue.Teaching) > 999999999999 {
		return errors.NewValidationError("Total lifetime value cannot be above 999999999999")
	}

	return nil
}
