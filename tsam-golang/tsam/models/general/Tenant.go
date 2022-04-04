package general

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Tenant contains details about the tenants using the tsam app.
type Tenant struct {
	Base
	TenantName    string `json:"tenantName" gorm:"type:varchar(100)"`
	TenantEmail   string `json:"tenantEmail" example:"abc@gmail.com" gorm:"type:varchar(50)"`
	TenantContact string `json:"tenantContact" example:"9700795509" gorm:"type:varchar(15)"`
}

// ValidateTenant validates all the fields of tenant.
func (tenant *Tenant) ValidateTenant() error {
	if util.IsEmpty(tenant.TenantName) {
		return errors.NewValidationError("Tenant name must be specified")
	}
	if util.IsEmpty(tenant.TenantEmail) || !util.ValidateEmail(tenant.TenantEmail) {
		return errors.NewValidationError("Tenant Email must be specified and should be of the type abc@domain.com")
	}
	if util.IsEmpty(tenant.TenantContact) || !util.ValidateContact(tenant.TenantContact) {
		return errors.NewValidationError("Tenant Contact must be specified and have 10 digits")
	}
	return nil
}
