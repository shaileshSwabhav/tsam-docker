package talentenquiry

import (
	uuid "github.com/satori/go.uuid"
)

// User is the DTO for salesperson of the enquiry.
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName  string    `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
}

// TableName will get Salesperson model from "users" table.
func (*User) TableName() string {
	return "users"
}

// func (salesperson *User) Validate() error {
// 	//check if salesperson id exists or not
// 	if !util.IsUUIDValid(salesperson.ID) {
// 		return errors.NewValidationError("Salesperson id must be specified")
// 	}
// 	if util.IsEmpty(salesperson.FirstName) || !util.ValidateString(salesperson.FirstName) {
// 		return errors.NewValidationError("User FirstName must be specified and must have characters only")
// 	}
// 	if util.IsEmpty(salesperson.LastName) || !util.ValidateString(salesperson.LastName) {
// 		return errors.NewValidationError("User LastName must be specified and must have characters only")
// 	}
// 	return nil
// }
