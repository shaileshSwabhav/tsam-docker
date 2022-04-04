package controller

// Contains all the API links & param information.
const (
	// Parent routes
	TenantIDLink      = "/tenant/{tenantID}"
	TalentEnquiryLink = "/talent-enquiry"
	PaginationLink    = "/limit/{limit}/offset/{offset}"

	// parameter names for their respective fields
	// make them private if you decide not to move them in a separate pkg.
	talentID      = "talentID"
	enquiryID     = "enquriyID"
	salesPersonID = "salespersonID"
	tenantID      = "tenantID"
)
