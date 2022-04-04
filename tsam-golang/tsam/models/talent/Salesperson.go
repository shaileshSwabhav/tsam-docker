package talent

import (
	uuid "github.com/satori/go.uuid"
)

//User is the DTO for salesperson of the talent
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName  string    `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
}

// TableName will get Salesperson model from "users" table
func (*User) TableName() string {
	return "users"
}
