package talent

import (
	"time"

	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"

	uuid "github.com/satori/go.uuid"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// NextAction contains the details of the next action planned to be taken on talent.
type NextAction struct {
	general.TenantBase

	// Maps.
	Courses      []list.Course        `json:"courses" gorm:"many2many:talent_next_actions_courses;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Companies    []list.CompanyBranch `json:"companies" gorm:"many2many:talent_next_actions_company_branches;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Technologies []general.Technology `json:"technologies" gorm:"many2many:talent_next_actions_technologies;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`

	// Related table IDs.
	TalentID     uuid.UUID `json:"talentID" gorm:"type:varchar(36)"`
	ActionTypeID uuid.UUID `json:"actionTypeID" gorm:"type:varchar(36)"`

	// Other fields.
	Stipend       *uint   `json:"stipend,omitempty" gorm:"type:int(6)"`
	ReferralCount *uint   `json:"referralCount,omitempty" gorm:"type:SMALLINT(4)"`
	FromDate      *string `json:"fromDate,omitempty" gorm:"type:date"`
	ToDate        *string `json:"toDate,omitempty" gorm:"type:date"`
	TargetDate    *string `json:"targetDate,omitempty" gorm:"type:date"`
	Comment       *string `json:"comment,omitempty" gorm:"type:varchar(1000)"`
}

// TableName will create the table for model NextAction with name talent_next_actions.
func (*NextAction) TableName() string {
	return "talent_next_actions"
}

// Validate will validate the fields of next action.
func (action *NextAction) Validate() error {

	// Next Action Type ID.
	if !util.IsUUIDValid(action.ActionTypeID) {
		return errors.NewValidationError("Next Action Type ID must be specified.")
	}

	// Stipend minimum.
	if action.Stipend != nil && *action.Stipend < 0 {
		return errors.NewValidationError("Stipend cannot be lesser than 0")
	}

	// Stipend maximum.
	if action.Stipend != nil && *action.Stipend > 999999 {
		return errors.NewValidationError("Stipend cannot be more than 999999")
	}

	// Referral count minimum.
	if action.ReferralCount != nil && *action.ReferralCount < 0 {
		return errors.NewValidationError("Referral Count cannot be lesser than 0")
	}

	// Referral Count maximum.
	if action.ReferralCount != nil && *action.ReferralCount > 9999 {
		return errors.NewValidationError("Referral Count cannot be more than 9999")
	}

	// Comment minimum.
	if action.Comment != nil && len(*action.Comment) < 0 {
		return errors.NewValidationError("Comment cannot be lesser than 0")
	}

	// Comment maximum.
	if action.Comment != nil && len(*action.Comment) > 1000 {
		return errors.NewValidationError("Comment cannot be more than 1000")
	}

	return nil
}

// ===========Defining many to many structs===========

// NextActionCourse is the map of next action and course.
type NextActionCourse struct {
	NextActionID uuid.UUID
	CourseID     uuid.UUID
}

// TableName defines table name of the struct.
func (*NextActionCourse) TableName() string {
	return "talent_next_actions_courses"
}

// NextActionCompanyBranch is the map of next action and company branch.
type NextActionCompanyBranch struct {
	NextActionID    uuid.UUID
	CompanyBranchID uuid.UUID
}

// TableName defines table name of the struct.
func (*NextActionCompanyBranch) TableName() string {
	return "talent_next_actions_company_branches"
}

// NextActionTechnology is the map of next action and technology.
type NextActionTechnology struct {
	NextActionID uuid.UUID
	TechnologyID uuid.UUID
}

// TableName defines table name of the struct.
func (*NextActionTechnology) TableName() string {
	return "talent_next_actions_technologies"
}

//************************************* DTO MODEL *************************************************************

// NextAction contains the details of the next action planned to be taken on talent.
type NextActionDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`

	// Maps.
	Courses      []list.Course        `json:"courses" gorm:"many2many:talent_next_actions_courses;association_jointable_foreignkey:course_id;jointable_foreignkey:next_action_id"`
	Companies    []list.CompanyBranch `json:"companies" gorm:"many2many:talent_next_actions_company_branches;association_jointable_foreignkey:company_branch_id;jointable_foreignkey:next_action_id"`
	Technologies []general.Technology `json:"technologies" gorm:"many2many:talent_next_actions_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:next_action_id"`

	// Related tables.
	ActionType   *NextActionType `json:"actionType" gorm:"foreignkey:ActionTypeID"`
	ActionTypeID uuid.UUID       `json:"-"`

	// Other fields.
	TalentID      uuid.UUID `json:"talentID"`
	Stipend       *uint     `json:"stipend"`
	ReferralCount *uint     `json:"referralCount"`
	FromDate      *string   `json:"fromDate"`
	ToDate        *string   `json:"toDate"`
	TargetDate    *string   `json:"targetDate"`
	Comment       *string   `json:"comment"`
}

// TableName will create the table for model NextAction with name talent_next_actions.
func (*NextActionDTO) TableName() string {
	return "talent_next_actions"
}
