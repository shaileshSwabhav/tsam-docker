package general

import uuid "github.com/satori/go.uuid"

// TokenDTO have Token, Name and Email
type TokenDTO struct {
	CredentialID   uuid.UUID  `json:"credentialID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	TenantID       uuid.UUID  `json:"tenantID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	FirstName      string     `json:"firstName" example:"John"`
	LastName       string     `json:"lastName" example:"Doe"`
	Email          string     `json:"email" example:"example@gmail.com"`
	RoleID         uuid.UUID  `json:"roleID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	LoginID        uuid.UUID  `json:"loginID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	DepartmentID   *uuid.UUID `json:"departmentID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	LoginSessionID *uuid.UUID `json:"loginSessionID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	Token          string     `json:"token" example:"sdhlasjdls09d789qowekjleuoqwueiqw#*Jl&*(hklh(UH&^()"`
}
