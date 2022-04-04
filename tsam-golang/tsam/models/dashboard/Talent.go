package dashboard

import (
	uuid "github.com/satori/go.uuid"
)

// TalentDashboard contains fields of talent to be displayed on the dashboard.
type TalentDashboard struct {
	TotalTalents             uint `json:"totalTalents"`
	TotalFreshers            uint `json:"totalFreshers"`
	TotalExperienced         uint `json:"totalExperienced"`
	B2CTalents               uint `json:"b2cTalents"`
	SwabhavTalents           uint `json:"swabhavTalents"`
	NonSwabhavTalents        uint `json:"nonSwabhavTalents"`
	TotalInterestedInForeign uint `json:"totalInterestedInForeign"`
	TotalLifetimeValue       uint `json:"totalLifetimeValue"`
}

// TalentSegregation classifies on the basis of talent's current academic year.
type TalentSegregation struct {
	FirstYear  uint `json:"firstYear"`
	SecondYear uint `json:"secondYear"`
	ThirdYear  uint `json:"thirdYear"`
	FourthYear uint `json:"fourthYear"`
}

// GetSumOfAllTalents adds 1,2,3 & 4th year talents and returns the sum.
func (ts *TalentSegregation) GetSumOfAllTalents() uint {
	return ts.FirstYear + ts.SecondYear + ts.ThirdYear + ts.FourthYear
}

// EnquiryDashboard contains fields of enquiry to be displayed on the dashboard.
// Data will only be of the last 30 days.
type EnquiryDashboard struct {
	TotalEnquiries       uint `json:"totalEnquiries"`
	NewEnquiries         uint `json:"newEnquiries"`
	EnquiriesAssigned    uint `json:"enquiriesAssigned"`
	EnquiriesNotHandled  uint `json:"enquiriesNotHandled"`
	EnquiriesNotAssigned uint `json:"enquiriesNotAssigned"`
	EnquiriesConverted   uint `json:"enquiriesConverted"`
}

// EnquirySource contains source of enquiry and count of that source
type EnquirySource struct {
	SourceID     uuid.UUID `json:"sourceID"`
	Description  string    `json:"description"`
	EnquiryCount uint      `json:"enquiryCount"`
}

// FacultyFeedbackToTalentDashboard contains minimum details displaying feedback on talent dasboard.
type FacultyFeedbackToTalentDashboard struct {
	Answer     float32 `json:"answer"`
	Question  string    `json:"question"`
}

// ThisAndPreviousWeekFacultyToTalentFeedackDashboard contains details for this week and previous week for
// displaying feedback on talent dasboard.
type ThisAndPreviousWeekFacultyToTalentFeedackDashboard struct {
	ThisWeekFeedbacks     []FacultyFeedbackToTalentDashboard `json:"thisWeekFeedbacks"`
	PreviousWeekFeedbacks     []FacultyFeedbackToTalentDashboard `json:"previousWeekFeedbacks"`
}
 
// FacultyFeedbackRatingLeaderBoard contains minimum details for displaying the ranks of the talents of a batch.
type FacultyFeedbackRatingLeaderBoard struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	TalentID string `json:"talentID"`
	Rating   float32 `json:"rating"`
	Rank     uint    `json:"rank"`
}
