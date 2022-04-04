package programming

import (

	uuid "github.com/satori/go.uuid"
)

//************************************* DTO MODEL *********************************************************

// ComplexConcept contains fields for getting complex concepts.
type ComplexConcept struct {
	ConceptName         string                       `json:"conceptName"`
	TalentID             uuid.UUID                   `json:"talentID"`
	Score                *float32                          `json:"score"`
	Complexity          	uint8                        `json:"complexity"`
	Level                uint8                    `json:"level"`
	Description            *string                     `json:"description"`
}


