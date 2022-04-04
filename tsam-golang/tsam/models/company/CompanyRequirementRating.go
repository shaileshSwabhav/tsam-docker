package company

// // CompanyRequirementRating will consist of table columns for compnay_requirement_rating.
// type CompanyRequirementRating struct {
// 	general.TenantBase
// 	RequirementID uuid.UUID `json:"requirementID" gorm:"type:varchar(36)"`
// 	QuestionID    uuid.UUID `json:"questionID" gorm:"type:varchar(36)"`
// 	OptionID      uuid.UUID `json:"optionID" gorm:"type:varchar(36)"`
// 	Answer        string    `json:"answer" gorm:"type:varchar(250)"`

// 	// models
// 	Requirement *Requirement              `json:"requirement" gorm:"foreignkey:RequirementID;association_autocreate:false;association_autoupdate:false"`
// 	Question    *general.FeedbackQuestion `json:"question" gorm:"foreignkey:QuestionID;association_autocreate:false;association_autoupdate:false"`
// 	Option      *general.FeedbackOption   `json:"option" gorm:"foreignkey:OptionID;association_autocreate:false;association_autoupdate:false"`
// }

// // ValidateRating will validate field of company_requirement_rating.
// func (rating *CompanyRequirementRating) ValidateRating() error {

// 	if rating.RequirementID == uuid.Nil {
// 		return errors.NewValidationError("Requirement must be specified")
// 	}

// 	if rating.QuestionID == uuid.Nil {
// 		return errors.NewValidationError("Question must be specified")
// 	}

// 	if rating.OptionID == uuid.Nil {
// 		return errors.NewValidationError("Option must be specified")
// 	}

// 	if util.IsEmpty(rating.Answer) {
// 		return errors.NewValidationError("Option must be specified")
// 	}

// 	return nil
// }
