package company

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// RequirementSearch contains all searchable fields of requirement.
type RequirementSearch struct {
	general.Address
	Qualifications      []*general.Degree     `json:"qualifications,omitempty"`
	PersonalityType     *string               `json:"personalityType,omitempty"`
	MinimumExperience   *uint8                `json:"minimumExperience,omitempty"`
	MaximumExperience   *uint8                `json:"maximumExperience,omitempty"`
	JobRole             *string               `json:"jobRole,omitempty"`
	Technologies        []*general.Technology `json:"technologies,omitempty"`
	SalesPersonID       *uuid.UUID            `json:"salesPersonID,omitempty"`
	PackageOffered      *uint64               `json:"packageOffered,omitempty"`
	RequirementFromDate *string               `json:"requirementFromDate,omitempty"`
	RequirementTillDate *string               `json:"requirementTillDate,omitempty"`
	IsActive            *bool                 `json:"isActive,omitempty"`
}
