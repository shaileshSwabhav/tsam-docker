package programming

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// ModuleProgrammingConcepts contains add update fields required for concept module.
type ModuleProgrammingConcepts struct {
	general.TenantBase
	ProgrammingConceptID   uuid.UUID                  `json:"programmingConceptID" gorm:"type:varchar(36)"`
	ModuleID   				uuid.UUID                  `json:"moduleID" gorm:"type:varchar(36)"`
	Level                     uint8                    `json:"level" gorm:"type:tinyint(2)"`
	ParentModuleProgrammingConcepts []ModuleProgrammingConcepts   `json:"parentModuleProgrammingConcepts" gorm:"many2many:module_programming_concepts_parent_module_programming_concepts;association_jointable_foreignkey:parent_module_programming_concept_id;jointable_foreignkey:module_programming_concept_id;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
}

// Validate validates compulsary fields of ConceptModule.
func (conceptModule *ModuleProgrammingConcepts) Validate() error {

	// Programming Concept ID.
	if !util.IsUUIDValid(conceptModule.ProgrammingConceptID) {
		return errors.NewValidationError("Programming Question ID must be a proper uuid")
	}
	
	// Module ID.
	if !util.IsUUIDValid(conceptModule.ModuleID) {
		return errors.NewValidationError("Module ID must be a proper uuid")
	}

	// Level must be specified.
	if conceptModule.Level < 0 {
		return errors.NewValidationError("Level cannot be lesser than 0")
	}

	// Level maximum.
	if conceptModule.Level > 99 {
		return errors.NewValidationError("Level cannot be above 99")
	}

	return nil
}

// TableName defines table name of the struct.
func (*ModuleProgrammingConcepts) TableName() string {
	return "modules_programming_concepts"
}

// //************************************* CONCEPT MODULE DTO *********************************************************

// ModuleProgrammingConceptsDTO all fields for comcept modules.
type ModuleProgrammingConceptsDTO struct {
	general.TenantBase
	ProgrammingConcept   	ProgrammingConcept        `json:"programmingConcept" gorm:"foreignkey:ProgrammingConceptID"`
	ProgrammingConceptID   uuid.UUID                  `json:"programmingConceptID"`
	ModuleID   				uuid.UUID                  `json:"moduleID"`
	Level                     uint8                    `json:"level"`
}

// TableName defines table name of the struct.
func (*ModuleProgrammingConceptsDTO) TableName() string {
	return "modules_programming_concepts"
}

// ModuleProgrammingConceptsForTalentScore contains fields for getting concept modules for talent score.
type ModuleProgrammingConceptsForTalentScore struct {
	general.TenantBase
	ProgrammingConceptID   uuid.UUID                  `json:"-"`
	ProgrammingConcept   	ProgrammingConcept        `json:"programmingConcept" gorm:"foreignkey:ProgrammingConceptID"`
	ModuleID   				uuid.UUID                  `json:"moduleID"`
	Level                     uint8                    `json:"level"`
	ParentModuleProgrammingConcepts []ModuleProgrammingConcepts   `json:"parentModuleProgrammingConcepts" gorm:"many2many:module_programming_concepts_parent_module_programming_concepts;association_jointable_foreignkey:parent_module_programming_concept_id;jointable_foreignkey:module_programming_concept_id;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
	AverageScore                    *float32                    `json:"averageScore"`
}

// TableName defines table name of the struct.
func (*ModuleProgrammingConceptsForTalentScore) TableName() string {
	return "modules_programming_concepts"
}

//************************************* MANY TO MANY MODEL *********************************************************

// ModuleProgrammingConceptsParentModuleProgrammingConcepts contains parent programming concept module details.
type ModuleProgrammingConceptsParentModuleProgrammingConcepts struct {
	ParentModuleProgrammingConceptID   uuid.UUID `json:"parentModuleProgrammingConceptID" gorm:"type:varchar(36)"`
	ModuleProgrammingConceptID uuid.UUID `json:"moduleProgrammingConceptID" gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*ModuleProgrammingConceptsParentModuleProgrammingConcepts) TableName() string {
	return "module_programming_concepts_parent_module_programming_concepts"
}

// Validate will validate the fields in parent programming concept module.
func (parentConcepModule *ModuleProgrammingConceptsParentModuleProgrammingConcepts) Validate() error {
	if !util.IsUUIDValid(parentConcepModule.ParentModuleProgrammingConceptID) {
		return errors.NewValidationError("Invalid concept module id")
	}
	if !util.IsUUIDValid(parentConcepModule.ModuleProgrammingConceptID) {
		return errors.NewValidationError("Invalid parent concept module id")
	}
	if parentConcepModule.ParentModuleProgrammingConceptID == parentConcepModule.ModuleProgrammingConceptID {
		return errors.NewValidationError("Concept module and parent concept module cannot be same")
	}
	return nil
}
