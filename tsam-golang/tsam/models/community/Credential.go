package community

import common "github.com/techlabs/swabhav/tsam/models/general"

// Credential Contain id, First and Last Name
type Credential struct {
	common.BaseDTO
	FirstName string `json:"firstName" example:"John"`
	LastName  string `json:"lastName,omitempty" example:"Deo"`
}
