package company

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
)

// EnquirySearch contains all searchable fields of enquiry.
type EnquirySearch struct {
	general.Address
	CompanyName   *string               `json:"companyName,omitempty"`
	Outcome       *string               `json:"outcome,omitempty"`
	Domains       []*Domain             `json:"domains,omitempty"`
	Technologies  []*general.Technology `json:"technologies,omitempty"`
	EnquiryDate   *string               `json:"enquiryDate,omitempty"`
	EnquiryType   *string               `json:"enquiryType,omitempty"`
	EnquirySource *string               `json:"enquirySource,omitempty"`
	SalesPersonID *uuid.UUID            `json:"salesPersonID,omitempty"`
	FromDate      *string               `json:"fromDate,omitempty" `
	TillDate      *string               `json:"tillDate,omitempty" `
}
