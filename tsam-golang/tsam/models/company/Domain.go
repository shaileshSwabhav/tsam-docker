package company

import (
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"

	"github.com/techlabs/swabhav/tsam/models/general"
)

// Domain struct contain domain information of company
type Domain struct {
	general.TenantBase
	DomainName string `json:"domainName" gorm:"varchar(100)"`
}

// ValidateDomain Return Error If Not Valid
func (domain *Domain) ValidateDomain() error {
	if util.IsEmpty(domain.DomainName) {
		return errors.NewValidationError("Domaine Name must be specified")
	}
	return nil
}
