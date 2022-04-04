package batch

import "github.com/techlabs/swabhav/tsam/models/general"

// ProgrammingProjectRatingParameter will store map for project rating
type ProgrammingProjectRatingParameter struct {
	general.TenantBase
	Label       string `json:"label" gorm:"type:varchar(25)"`
	Description string `json:"description" gorm:"type:varchar(250)"`
	TotalScore  uint   `json:"total_score" gorm:"type:int"`
	IsActive    *bool  `json:"isActive" gorm:"default:true"`
}

// TableName overrides name of the table
func (*ProgrammingProjectRatingParameter) TableName() string {
	return "programming_project_rating_parameter"
}
