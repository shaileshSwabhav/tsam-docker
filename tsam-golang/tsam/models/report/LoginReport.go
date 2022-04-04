package report

import uuid "github.com/satori/go.uuid"

// CredentialLoginReport will contain details of login for specified credential.
type CredentialLoginReport struct {
	LoginName string `json:"loginName"`
	// RoleName   string  `json:"roleName"`
	LoginTime  string  `json:"loginTime"`
	LogoutTime *string `json:"logoutTime"`
	TotalHours string  `json:"totalHours"`
}

type LoginReport struct {
	LoginSessionID uuid.UUID `json:"loginSessionID"`
	CredentialID   uuid.UUID `json:"credentialID"`
	LoginName      string    `json:"loginName"`
	RoleName       string    `json:"roleName"`
	LastLoginTime  string    `json:"lastLoginTime"`
	LastLogoutTime *string   `json:"lastLogoutTime"`
	LoginTime      *string   `json:"loginTime"`
	LogoutTime     *string   `json:"logoutTime"`
	LoginCount     int       `json:"loginCount"`
	TotalHours     string    `json:"totalHours"`
}
