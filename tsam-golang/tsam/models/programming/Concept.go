package programming

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ProgrammingConcept contains add update fields required for programming concept.
type ProgrammingConcept struct {
	general.TenantBase
	Name                string                       `json:"name" gorm:"type:varchar(100)"`
	IsModuleIndependent bool                         `json:"isModuleIndependent"`
	Complexity          uint8                        `json:"complexity" gorm:"type:tinyint(2)"`
	ConceptModules      []*ModuleProgrammingConcepts `json:"modules,omitempty"`
	// ProgrammingQuestions   []list.ProgrammingQuestion `json:"programmingQuestions" gorm:"many2many:programming_concepts_programming_questions;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	Description            *string                     `json:"description" gorm:"type:varchar(2000)"`
	// Logo                   *string                     `json:"logo" gorm:"type:varchar(200)"`
	// ProgrammingConceptID   *uuid.UUID                  `json:"programmingConceptID" gorm:"type:varchar(36)"`
	// SubProgrammingConcepts []*ProgrammingConcept       `json:"subProgrammingConcepts" gorm:"foreignkey:ProgrammingConceptID"`
}

// TableName defines table name of the struct.
func (*ProgrammingConcept) TableName() string {
	return "programming_concepts"
}

// Validate validates compulsary fields of ProgrammingConcept.
func (concept *ProgrammingConcept) Validate() error {

	// Check if name is blank or not.
	if util.IsEmpty(concept.Name) {
		return errors.NewValidationError("Name must be specified")
	}

	// Name maximum characters.
	if len(concept.Name) > 100 {
		return errors.NewValidationError("Name can have maximum 100 characters")
	}

	// Level must be specified.
	if concept.Complexity == 0 {
		return errors.NewValidationError("Level must be specified")
	}

	// Level maximum.
	if concept.Complexity > 99 {
		return errors.NewValidationError("Level cannot be above 99")
	}

	// // Is Module independent.
	// if !concept.IsModuleIndependent && (concept.ConceptModules == nil || (concept.ConceptModules != nil && len(concept.ConceptModules) == 0)) {
	// 	return errors.NewValidationError("If concept is not module independent then spcify the modules")
	// }

	// Check if all concpets are having unique module ids.
	moduleMap := make(map[uuid.UUID]uint)
	for _, course := range concept.ConceptModules {
		moduleMap[course.ModuleID]++
		if moduleMap[course.ModuleID] > 1 {
			return errors.NewValidationError("Concept cannot have same modules")
		}
	}

	// Description maximum characters.
	if concept.Description != nil && len(*concept.Description) > 2000 {
		return errors.NewValidationError("Description can have maximum 2000 characters")
	}

	// // Programming Questions.
	// if concept.ProgrammingQuestions == nil  {
	// 	return errors.NewValidationError("Programming Questions must be specified")
	// }

	// // Programming Questions.
	// if concept.ProgrammingQuestions != nil && len(concept.ProgrammingQuestions) <= 0 {
	// 	return errors.NewValidationError("Programming Questions must be specified")
	// }

	// // If programming concept id is null then description must be present.
	// if concept.ProgrammingConceptID == nil && concept.Description == nil {
	// 	return errors.NewValidationError("Description must be present if it is a parent programming concept")
	// }

	return nil
}

//************************************* DTO MODEL *********************************************************

// ProgrammingConceptDTO contains add update fields required for programming concept.
type ProgrammingConceptDTO struct {
	ID                  uuid.UUID                    `json:"id"`
	DeletedAt           *time.Time                   `json:"-"`
	Name                string                       `json:"name"`
	IsModuleIndependent bool                         `json:"isModuleIndependent"`
	Complexity          uint8                        `json:"complexity"`
	ConceptModules      []*ModuleProgrammingConcepts `json:"modules" gorm:"foreignkey:ProgrammingConceptID"`
	Description            *string                     `json:"description"`
	// Logo                   *string                     `json:"logo" gorm:"type:varchar(200)"`
	// ProgrammingConceptID   *uuid.UUID                  `json:"programmingConceptID" gorm:"type:varchar(36)"`
	// SubProgrammingConcepts []*ProgrammingConcept       `json:"subProgrammingConcepts" gorm:"foreignkey:ProgrammingConceptID"`
	// ProgrammingQuestions   []list.ProgrammingQuestion `json:"programmingQuestions" gorm:"many2many:programming_concepts_programming_questions;association_jointable_foreignkey:programming_question_id;jointable_foreignkey:programming_concept_id"`
}

// TableName defines table name of the struct.
func (*ProgrammingConceptDTO) TableName() string {
	return "programming_concepts"
}

// ===========Defining many to many structs===========

// ProgrammingConceptsProgrammingQuestions is the map of programming concept and programming question.
type ProgrammingConceptsProgrammingQuestions struct {
	ProgrammingConceptID  uuid.UUID `gorm:"type:varchar(36)"`
	ProgrammingQuestionID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ProgrammingConceptsProgrammingQuestions) TableName() string {
	return "programming_concepts_programming_questions"
}
