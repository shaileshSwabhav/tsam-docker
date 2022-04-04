package util

import "time"

// GetISTINRFC3339 Return IST Date time in RFC3339 Format
func GetISTINRFC3339() time.Time {
	dateTime := time.Now().Add(5 * time.Hour).Add(30 * time.Minute)
	dateTime.Format(time.RFC3339)
	return dateTime
}
