package callreport

import (
	uuid "github.com/satori/go.uuid"
)

// LoginwiseCallingReport will provide calling records for every credential.
type LoginwiseCallingReport struct {
	CredentialID      uuid.UUID `json:"credentialID"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	TotalCallingCount uint      `json:"totalCallingCount"`
	TotalTalentCount  uint      `json:"totalTalentCount"`
}

// TalentCallingReportDTO will provide call records and talent data.
type TalentCallingReportDTO struct {
	// LoginID   string `json:"loginID"`
	LoginName string `json:"loginName"`
	// Talent details
	Talent `json:"talent"`
	// call record details
	DateTime     string  `json:"dateTime"`
	Purpose      string  `json:"purpose"`
	Outcome      string  `json:"outcome"`
	Comment      *string `json:"comment,omitempty"`
	ExpectedCTC  *uint   `json:"expectedCTC,omitempty"`
	NoticePeriod *uint8  `json:"noticePeriod,omitempty"`
	TargetDate   *string `json:"targetDate,omitempty"`
}

// Talent contains the fields needed for call record.
type Talent struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Contact   string `json:"contact"`
	Email     string `json:"email"`
}

// DaywiseCallingReport will provide calling records of a particular day.
type DaywiseCallingReport struct {
	Date              string `json:"date"`
	TotalCallingCount uint   `json:"totalCallingCount"`
	TotalTalentCount  uint   `json:"totalTalentCount"`
}

// =======================================================TALENT-ENQUIRY=======================================================

// LoginwiseTalentEnquiryCallingReport will provide calling records for every credential.
type LoginwiseTalentEnquiryCallingReport struct {
	CredentialID      uuid.UUID `json:"credentialID"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	TotalCallingCount uint      `json:"totalCallingCount"`
	TotalEnquiryCount uint      `json:"totalEnquiryCount"`
}

// DaywiseTalentEnquiryCallingReport will provide calling records of a particular day for talent-enquiry.
type DaywiseTalentEnquiryCallingReport struct {
	Date              string `json:"date"`
	TotalCallingCount uint   `json:"totalCallingCount"`
	TotalEnquiryCount uint   `json:"totalEnquiryCount"`
}

// TalentEnquiryCallingReportDTO will provide call records and talent-enquiry data.
type TalentEnquiryCallingReportDTO struct {
	// LoginID   string `json:"loginID"`
	LoginName string `json:"loginName"`
	// Talent details
	TalentEnquiry `json:"talentEnquiry"`
	// call record details
	DateTime     string  `json:"dateTime"`
	Purpose      string  `json:"purpose"`
	Outcome      string  `json:"outcome"`
	Comment      *string `json:"comment,omitempty"`
	ExpectedCTC  *uint   `json:"expectedCTC,omitempty"`
	NoticePeriod *uint8  `json:"noticePeriod,omitempty"`
	TargetDate   *string `json:"targetDate,omitempty"`
}

// TalentEnquiry contains the fields needed for call record.
type TalentEnquiry struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Contact   string `json:"contact"`
	Email     string `json:"email"`
}
