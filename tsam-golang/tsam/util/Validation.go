package util

import (
	"regexp"
)

// ValidateString validates if string is valid and there is no space in it
func ValidateString(name string) bool {
	stringPattern := regexp.MustCompile("^[a-zA-Z]*$")
	return stringPattern.MatchString(name)
}

// ValidateStringWithSpace validates string that contains space
func ValidateStringWithSpace(str string) bool {
	stringPattern := regexp.MustCompile("^[a-zA-Z ]*$")
	return stringPattern.MatchString(str)
}

// ValidateContact validates 10 digit contact number statring from 0/+91
// Allowed contact numbers -> 9883443344, 09883443344, 0919883443344, +919883443344.....
func ValidateContact(contact string) bool {
	contactPattern := regexp.MustCompile(`^(?:(?:\+|0{0,2})91(\s*[\-]\s*)?|[0]?)?[6789]\d{9}$`)
	return contactPattern.MatchString(contact)
}

// ValidateEmail validates email which should be of the type example@domain.com
func ValidateEmail(email string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,4}`)
	return emailPattern.MatchString(email)
}

// ValidateAddress validates if address is valid or not. (regex to be added, currently checks if address is empty)
func ValidateAddress(address string) bool {
	return IsEmpty(address)
}

// ValidateTenantID validates tenantID in the tenant table
// func ValidateTenantID(db *gorm.DB, id uuid.UUID, out interface{}) (bool, error) {
// 	exists, err := repository.DoesRecordExist(db, out, repository.Filter("`id` = ?", id))
// 	if err != nil {
// 		return false, err
// 	}
// 	return exists, nil
// }
