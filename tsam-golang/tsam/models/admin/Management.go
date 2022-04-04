package admin

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// Employee contains the credential details of employees with supervisors.
type Employee struct {
	general.TenantBase
	FirstName   string       `json:"firstName" example:"Ravi" `
	LastName    string       `json:"lastName" example:"Sharma" `
	Email       string       `json:"email" example:"abc@gmail.com" `
	Contact     string       `json:"contact" example:"9700795509"`
	RoleID      uuid.UUID    `json:"-" gorm:"type:varchar(36)"`
	Role        general.Role `json:"role"`
	Supervisors []Employee   `json:"supervisors" gorm:"many2many:employee_supervisors;association_jointable_foreignkey:supervisor_credential_id;jointable_foreignkey:employee_credential_id"`
}

// TableName sets table name as credentials.
func (*Employee) TableName() string {
	return "credentials"
}

// EmployeeSupervisor contains supervisor details.
type EmployeeSupervisor struct {
	EmployeeCredentialID   uuid.UUID `json:"employeeCredentialID" gorm:"type:varchar(36)"`
	SupervisorCredentialID uuid.UUID `json:"supervisorCredentialID" gorm:"type:varchar(36)"`
}

// TableName sets table name as credentials.
func (*Employee) EmployeeSupervisor() string {
	return "employee_supervisors"
}

// Validate will validate the fields in supervisor.
func (supervisor *EmployeeSupervisor) Validate() error {
	if !util.IsUUIDValid(supervisor.EmployeeCredentialID) {
		return errors.NewValidationError("Invalid employee credential id")
	}
	if !util.IsUUIDValid(supervisor.SupervisorCredentialID) {
		return errors.NewValidationError("Invalid supervisor credential id")
	}
	if supervisor.EmployeeCredentialID == supervisor.SupervisorCredentialID {
		return errors.NewValidationError("Employee and supervisor can't be the same.")
	}
	return nil
}

// CountModel is used for gettig count from database.
type CountModel struct {
	TotalCount int `json:"totalCount"`
}
