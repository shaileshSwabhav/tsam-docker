package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Feeling will consist of different feeling names.
type Feeling struct {
	TenantBase
	FeelingName   string         `json:"feelingName" gorm:"type:varchar(50)"`
	FeelingLevels []FeelingLevel `json:"feelingLevels"`
}

// ValidateFeeling will check if valid name for feeling is provided.
func (feeling *Feeling) ValidateFeeling() error {

	// Feeling level must be specified.
	if util.IsEmpty(feeling.FeelingName) {
		return errors.NewValidationError("Feeling name must be specified")
	}

	// Feeling level maximum characters.
	if len(feeling.FeelingName) > 50 {
		return errors.NewValidationError("Feeling name can have maximum 50 characters")
	}

	// Check if all feeling levels are having unique feeling number.
	feelingLevelMap := make(map[uint]uint)
	for _, feelingLevel := range feeling.FeelingLevels {
		feelingLevelMap[feelingLevel.LevelNumber]++
		if feelingLevelMap[feelingLevel.LevelNumber] > 1 {
			return errors.NewValidationError("Same feeling level number cannot be assigned to multiple levels")
		}
	}

	// Validate feeling levels.
	if feeling.FeelingLevels != nil {
		for _, feelingLevel := range feeling.FeelingLevels {
			err := feelingLevel.Validate()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
