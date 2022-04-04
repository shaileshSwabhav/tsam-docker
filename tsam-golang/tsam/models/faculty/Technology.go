package faculty

// // Technology contain single FacultyTechnology detail
// type Technology struct {
// 	general.TenantBase
// 	Language  string    `json:"language" gorm:"type:varchar(50)"`
// 	FacultyID uuid.UUID `json:"facultyID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
// }

// // ValidateFacultyTechnologies validates fields of the facultyTechnologies
// func (technologies *Technology) ValidateFacultyTechnologies() error {
// 	if util.IsEmpty(technologies.Language) || !util.ValidateStringWithSpace(technologies.Language) {
// 		return errors.NewValidationError("Language is compulsory and should contain only characters")
// 	}
// 	return nil
// }

// // TableName overrides name of the table
// func (*Technology) TableName() string {
// 	return "faculties_technologies"
// }
