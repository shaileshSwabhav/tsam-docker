package report

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
)

// NextActionReport will provide next action for every credential.
type NextActionReport struct {
	CredentialID         uuid.UUID `json:"credentialID"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	TotalNextActionCount uint      `json:"totalNextActionCount"`
	TotalTalentCount     uint      `json:"totalTalentCount"`
}

// TalentNextActionReportDTO will provide next action and talent data.
type TalentNextActionReportDTO struct {
	LoginID   string `json:"loginID"`
	LoginName string `json:"loginName"`
	// Talent details
	Talent `json:"talent"`
	// next action details
	ID            string               `json:"id"`
	TalentID      string               `json:"talentID"`
	Stipend       *uint                `json:"stipend,omitempty"`
	ReferralCount *uint                `json:"referralCount,omitempty"`
	FromDate      *string              `json:"fromDate,omitempty"`
	ToDate        *string              `json:"toDate,omitempty"`
	TargetDate    *string              `json:"targetDate,omitempty"`
	Comment       *string              `json:"comment,omitempty"`
	ActionType    string               `json:"actionType"`
	Courses       []list.Course        `json:"courses,omitempty"`
	Companies     []list.CompanyBranch `json:"companies,omitempty"`
	Technologies  []general.Technology `json:"technologies,omitempty"`
}

// Talent contains the fields needed for next action record.
type Talent struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Contact   string `json:"contact"`
	Email     string `json:"email"`
}

type NextActionSearch struct {
	LoginID       *string      `json:"loginID,omitempty"`
	Stipend       *uint        `json:"stipend,omitempty"`
	ReferralCount *uint        `json:"referralCount,omitempty"`
	FromDate      *string      `json:"fromDate,omitempty"`
	ToDate        *string      `json:"toDate,omitempty"`
	TargetDate    *string      `json:"targetDate,omitempty"`
	Comment       *string      `json:"comment,omitempty"`
	ActionType    *string      `json:"actionType"`
	Courses       *[]uuid.UUID `json:"courses,omitempty"`
	Companies     *[]uuid.UUID `json:"companies,omitempty"`
	Technologies  *[]uuid.UUID `json:"technologies,omitempty"`
}
