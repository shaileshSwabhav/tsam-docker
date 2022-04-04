package college

import (
	uuid "github.com/satori/go.uuid"
)

// DriveCandidates defines the association between campus drive & talent.
type DriveCandidates struct {
	CampusDriveID  uuid.UUID `gorm:"type:varchar(36)"`
	TalentID       uuid.UUID `gorm:"type:varchar(36)"`
	IsTestLinkSent bool      `gorm:"default:false;not null"`
}

// SeminarCandidates defines the association between Seminar & talent.
type SeminarCandidates struct {
	SeminarID uuid.UUID `gorm:"type:varchar(36)"`
	TalentID  uuid.UUID `gorm:"type:varchar(36)"`
}

// TalentTechnologies defines the association between talent and technologies.
type TalentTechnologies struct {
	TalentID     uuid.UUID `gorm:"type:varchar(36)"`
	TechnologyID uuid.UUID `gorm:"type:varchar(36)"`
}
