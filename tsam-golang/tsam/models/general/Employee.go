package general

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

//Employee will contain employees information
type Employee struct {
	TenantBase
	Address
	Code          string        `json:"code" gorm:"type:varchar(10);not null"`
	FirstName     string        `json:"firstName" example:"Ravi" gorm:"type:varchar(50)"`
	LastName      string        `json:"lastName" example:"Sharma" gorm:"type:varchar(50)"`
	Email         string        `json:"email" example:"abc@gmail.com" gorm:"type:varchar(50)"`
	Contact       string        `json:"contact" example:"9700795509" gorm:"type:varchar(15)"`
	DateOfBirth   *string       `json:"dateOfBirth" example:"2000-02-30" gorm:"type:date"`
	DateOfJoining *string       `json:"dateOfJoining" example:"2019-02-23" gorm:"type:date"`
	Technologies  []*Technology `json:"technologies" gorm:"many2many:employees_technologies;association_autoupdate:false;"`
	Resume        *string       `json:"resume" gorm:"type:varchar(200)"`
	IsActive      bool          `json:"isActive" gorm:"DEFAULT:true"`
	Type          string        `json:"type" gorm:"varchar(20)"`
	RoleID        uuid.UUID     `json:"-" gorm:"type:varchar(36)"`
	Role          *Role         `json:"role" gorm:"foreignkey:RoleID;association_autoupdate:false;"`
	// Academics     []*Academic   `json:"academics"`
	// Experiences   []*Experience         `json:"experiences"`
	// temp
	Password string `json:"-" gorm:"type:varchar(255)"`
}

//association_foreignkey:technology_id;foreignkey:employee_id

// ===========Defining many to many structs===========

// // EmployeeTechnologies is the map of employee and technologies.
// type EmployeeTechnologies struct {
// 	EmployeeID   uuid.UUID
// 	TechnologyID uuid.UUID
// }

// // TableName will create table with name employees_technologies
// func (*EmployeeTechnologies) TableName() string {
// 	return "employees_technologies"
// }

// ValidateEmployee validate employee details
func (employee *Employee) ValidateEmployee() error {
	if util.IsEmpty(employee.FirstName) || !util.ValidateString(employee.FirstName) {
		return errors.NewValidationError("Employee FirstName must be specified and must have characters only")
	}

	if util.IsEmpty(employee.LastName) || !util.ValidateString(employee.LastName) {
		return errors.NewValidationError("Employee LastName must be specified and must have characters only")
	}

	if util.IsEmpty(employee.Email) || !util.ValidateEmail(employee.Email) {
		return errors.NewValidationError("Employee Email must be specified and should be of the type abc@domain.com")
	}

	if util.IsEmpty(employee.Contact) || !util.ValidateContact(employee.Contact) {
		return errors.NewValidationError("Employee Contact must be specified and have 10 digits")
	}

	if !util.IsUUIDValid(employee.TenantID) {
		return errors.NewValidationError("Invalid tenant ID")
	}

	if employee.Role == nil {
		return errors.NewValidationError("Role must be specified")
	}

	err := employee.ValidateAddress()
	if err != nil {
		return err
	}
	// if util.IsEmpty(employee.Type) {
	// 	return errors.NewValidationError("Employee type must be specified")
	// }
	// if employee.Academics == nil {
	// 	return errors.NewValidationError("Academics must be specified")
	// }
	// for _, academic := range employee.Academics {
	// 	if err := academic.ValidateEmployeeAcademic(); err != nil {
	// 		return errors.NewValidationError(err.Error())
	// 	}
	// }
	// if employee.Experiences != nil {
	// 	for _, experience := range employee.Experiences {
	// 		if err := experience.ValidateEmployeeExperiences(); err != nil {
	// 			return errors.NewValidationError(err.Error())
	// 		}
	// 	}
	// }
	// if employee.Technologies != nil {
	// 	for _, technology := range employee.Technologies {
	// 		if err := technology.Validate(); err != nil {
	// 			return errors.NewValidationError(err.Error())
	// 		}
	// 	}
	// }

	return nil
}

//Employee will contain employees information
type EmployeeDTO struct {
	ID        uuid.UUID  `json:"id"`
	DeletedAt *time.Time `json:"-"`
	Address
	Code          string        `json:"code"`
	FirstName     string        `json:"firstName"`
	LastName      string        `json:"lastName"`
	Email         string        `json:"email"`
	Contact       string        `json:"contact"`
	DateOfBirth   *string       `json:"dateOfBirth"`
	DateOfJoining *string       `json:"dateOfJoining"`
	Technologies  []*Technology `json:"technologies" gorm:"many2many:employees_technologies;association_jointable_foreignkey:technology_id;jointable_foreignkey:employee_id;"`
	Resume        *string       `json:"resume"`
	IsActive      bool          `json:"isActive"`
	Type          string        `json:"type"`
	RoleID        uuid.UUID     `json:"-"`
	Role          *Role         `json:"role" gorm:"foreignkey:RoleID"`
	// Academics     []*Academic   `json:"academics"`
	// Experiences   []*Experience         `json:"experiences"`
}

// TableName defines table name of the struct.
func (*EmployeeDTO) TableName() string {
	return "employees"
}
