package general

import (
	"regexp"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Address Definition
type Address struct {
	Address   *string    `json:"address,omitempty"`
	City      *string    `json:"city,omitempty" gorm:"type:varchar(50)"`
	PINCode   *uint32    `json:"pinCode,omitempty" gorm:"type:int(10)"`
	Country   *Country   `json:"country,omitempty" gorm:"association_autocreate:false;association_autoupdate:false"`
	State     *State     `json:"state,omitempty" gorm:"association_autocreate:false;association_autoupdate:false"`
	StateID   *uuid.UUID `json:"-" gorm:"type:varchar(36)"`
	CountryID *uuid.UUID `json:"-" gorm:"type:varchar(36)"`
}

// ValidateAddress checks if address is valid, if not returns error
func (address *Address) ValidateAddress() error {
	pinCodePattern := regexp.MustCompile("^[1-9][0-9]{5}$")
	// if util.ValidateAddress(address.Address) {
	// 	return errors.NewValidationError("Address must be specified and should have only characters")
	// }
	if address.City != nil {
		if !util.ValidateStringWithSpace(*address.City) {
			return errors.NewValidationError("City must be specified and should have only characters and space")
		}
	}
	if address.PINCode != nil {
		if !pinCodePattern.MatchString(strconv.Itoa(int(*address.PINCode))) {
			return errors.NewValidationError("Invalid Pincode")
		}
	}
	// if util.IsUUIDValid(address.CountryID) {
	// 	return errors.NewValidationError("CountryID must be specified")
	// }
	// if util.IsUUIDValid(address.StateID) {
	// 	return errors.NewValidationError("StateID must be specified")
	// }
	// if address.State != nil {
	// 	if err := address.State.ValidateState(); err != nil {
	// 		return err
	// 	}
	// }
	// if address.Country != nil {
	// 	if err := address.Country.ValidateCountry(); err != nil {
	// 		return err
	// 	}
	// }
	// Add Validation in Main Model
	// if address.PINCode == 0 {
	// 	return errors.NewValidationError("PIN Code must be specified")
	// }
	// if util.IsEmpty(address.City) {
	// 	return errors.NewValidationError("City must be specified")
	// }
	return nil
}

// MandatoryValidation can be used where all fields of address is mandatory.
func (address *Address) MandatoryValidation() error {

	switch {
	case address.Address == nil:
		return errors.NewValidationError("Address must be specified")
	case address.City == nil:
		return errors.NewValidationError("City must be specified")
	case address.Country == nil:
		return errors.NewValidationError("Country must be specified")
	case address.State == nil:
		return errors.NewValidationError("State must be specified")
	case address.PINCode == nil:
		return errors.NewValidationError("PIN code must be specified")
	}

	pinCodePattern := regexp.MustCompile("^[1-9][0-9]{5}$")
	if util.IsEmpty(*address.Address) {
		return errors.NewValidationError("Address must not be empty")
	}
	if util.IsEmpty(*address.City) || !util.ValidateStringWithSpace(*address.City) {
		return errors.NewValidationError("City must be non-empty and characters only")
	}
	if !pinCodePattern.MatchString(strconv.Itoa(int(*address.PINCode))) {
		return errors.NewValidationError("Invalid Pincode")
	}
	if !util.IsUUIDValid(address.Country.ID) && util.IsEmpty(address.Country.Name) {
		return errors.NewValidationError("Country must be specified")
	}
	if !util.IsUUIDValid(address.State.ID) && util.IsEmpty(address.State.Name) {
		return errors.NewValidationError("State must be specified")
	}
	// if err := address.State.ValidateState(); err != nil {
	// 	return err
	// }
	// if err := address.Country.ValidateCountry(); err != nil {
	// 	return err
	// }
	return nil
}
