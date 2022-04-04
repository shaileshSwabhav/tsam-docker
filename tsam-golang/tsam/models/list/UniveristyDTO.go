package list

import "github.com/techlabs/swabhav/tsam/models/general"

// use ID instead of general. Avoid unnecessary imports

// University contains info about university.
type University struct {
	general.TenantBase
	UniversityName string `json:"universityName"`
}
