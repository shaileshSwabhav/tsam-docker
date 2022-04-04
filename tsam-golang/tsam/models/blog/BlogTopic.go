package blog

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// BlogTopic contains add update fields required for blog topic.
type BlogTopic struct {
	general.TenantBase
	Name string `json:"name" gorm:"type:varchar(100)"`
}

// Validate validates compulsary fields of BlogTopic.
func (topic *BlogTopic) Validate() error {

	// Check if name is blank or not.
	if util.IsEmpty(topic.Name) {
		return errors.NewValidationError("Name must be specified")
	}

	// Name maximum characters.
	if len(topic.Name) > 100 {
		return errors.NewValidationError("Name can have maximum 100 characters")
	}

	return nil
}
