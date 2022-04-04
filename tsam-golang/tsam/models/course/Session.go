package course

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/util"
)

// CourseSession Contain Sessions Details.
type CourseSession struct {
	general.TenantBase
	Name           string               `json:"name"`
	Hours          *float32             `json:"hours" gorm:"type:decimal(6,2)"`
	Order          uint                 `json:"order" gorm:"type:TINYINT"`
	StudentOutput  string               `json:"studentOutput" gorm:"type:varchar(100)"`
	SubSessions    []CourseSession      `json:"subSessions" gorm:"foreignkey:SessionID;association_autoupdate:false"`
	SessionID      *uuid.UUID           `json:"sessionID" gorm:"type:varchar(36)"`
	CourseID       uuid.UUID            `json:"-" gorm:"type:varchar(36)"`
	CourseModuleID uuid.UUID            `json:"courseModuleID" gorm:"type:varchar(36)"`
	Resources      []*resource.Resource `json:"resource" gorm:"many2many:course_sessions_resources;association_jointable_foreignkey:resource_id;jointable_foreignkey:course_session_id;association_autocreate:false;association_autoupdate:false;"`
}

// TableName overrides name of the table
func (*CourseSession) TableName() string {
	return "course_sessions"
}

// ValidateSession validates fields of session
func (session *CourseSession) ValidateSession() error {

	if util.IsEmpty(session.Name) {
		return errors.NewValidationError("Session name must be specified")
	}

	if session.Hours != nil {
		if *session.Hours <= 0 {
			return errors.NewValidationError("Session length should not be zero")
		}
	}

	if session.Order <= 0 {
		return errors.NewValidationError("Session order should not be zero")
	}

	if util.IsEmpty(session.StudentOutput) {
		return errors.NewValidationError("Student output must be specified")
	}

	if session.SubSessions != nil {
		for _, subSession := range session.SubSessions {
			if err := subSession.ValidateSession(); err != nil {
				return err
			}
		}
	}

	// for _, resource := range session.Resources {
	// 	if err := resource.ValidateResource(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// CourseSessionDTO Contain Sessions Details.
type CourseSessionDTO struct {
	ID                           uuid.UUID                      `json:"id" gorm:"primaryKey"`
	DeletedAt                    *time.Time                     `json:"-"`
	Name                         string                         `json:"name"`
	Hours                        *float32                       `json:"hours"`
	Order                        uint                           `json:"order"`
	StudentOutput                string                         `json:"studentOutput"`
	SubSessions                  []CourseSessionDTO             `json:"subSessions" gorm:"foreignkey:SessionID"`
	SessionID                    *uuid.UUID                     `json:"sessionID"`
	CourseID                     uuid.UUID                      `json:"-"`
	CourseModuleID               uuid.UUID                      `json:"-"`
	Resources                    []*resource.ResourceDTO        `json:"resource" gorm:"many2many:course_sessions_resources;association_jointable_foreignkey:resource_id;jointable_foreignkey:course_session_id;"`
	CourseProgrammingAssignments []*TopicProgrammingQuestionDTO `json:"courseProgrammingAssignment"`
}

// TableName overrides name of the table
func (*CourseSessionDTO) TableName() string {
	return "course_sessions"
}
