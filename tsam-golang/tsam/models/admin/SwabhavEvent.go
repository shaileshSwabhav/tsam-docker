package admin

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/util"
)

// SwabhavEvent struct will contain details about event.
type SwabhavEvent struct {
	general.TenantBase
	general.Address
	Title                string  `json:"title" gorm:"type:varchar(50)"`
	Description          string  `json:"description" gorm:"type:varchar(500)"`
	EntryFee             uint    `json:"entryFee" gorm:"type:int"`
	IsOnline             *bool   `json:"isOnline"`
	FromDate             string  `json:"fromDate" gorm:"type:date"`
	ToDate               string  `json:"toDate" gorm:"type:date"`
	FromTime             string  `json:"fromTime" gorm:"type:time"`
	ToTime               string  `json:"toTime" gorm:"type:time"`
	TotalHours           uint    `json:"totalHours" gorm:"type:int"`
	LastRegistrationDate string  `json:"lastRegistrationDate" gorm:"type:date"`
	EventStatus          string  `json:"eventStatus" gorm:"type:varchar(30)"`
	EventMeetingLink     *string `json:"eventMeetingLink" gorm:"type:varchar(500)"`
	IsActive             *bool   `json:"isActive"`
	EventImage           *string `json:"eventImage" gorm:"type:varchar(200)"`
}

// ValidateEvent will validate all the fields of event struct.
func (event *SwabhavEvent) ValidateEvent() error {

	err := event.ValidateAddress()
	if err != nil {
		return err
	}

	// Dates validation
	fromDate, err := time.Parse(time.RFC3339, event.FromDate)
	if err != nil {
		return err
	}

	toDate, err := time.Parse(time.RFC3339, event.ToDate)
	if err != nil {
		return err
	}

	lastRegistrationDate, err := time.Parse(time.RFC3339, event.LastRegistrationDate)
	if err != nil {
		return err
	}

	// After reports whether the time instant from date is after to date.
	if fromDate.After(toDate) {
		return errors.NewValidationError("To date must be greater than or equal to from date.")
	}

	if lastRegistrationDate.After(fromDate) {
		return errors.NewValidationError("Last registration date must be less than from date.")
	}

	event.FromDate = fromDate.Format("2006-01-02")
	event.ToDate = toDate.Format("2006-01-02")
	event.LastRegistrationDate = lastRegistrationDate.Format("2006-01-02")

	if util.IsEmpty(event.Title) {
		return errors.NewValidationError("Title must be specified.")
	}

	if util.IsEmpty(event.Description) {
		return errors.NewValidationError("Description must be specified.")
	}

	if event.IsOnline == nil {
		return errors.NewValidationError("Whether event is online or offline must be specified.")
	}

	if event.IsActive == nil {
		return errors.NewValidationError("Event's active status must be specified.")
	}

	if util.IsEmpty(event.FromDate) {
		return errors.NewValidationError("From date must be specified.")
	}

	if util.IsEmpty(event.ToDate) {
		return errors.NewValidationError("To date must be specified.")
	}

	if util.IsEmpty(event.FromTime) {
		return errors.NewValidationError("From time must be specified.")
	}

	if util.IsEmpty(event.ToTime) {
		return errors.NewValidationError("To time must be specified.")
	}

	if event.TotalHours == 0 {
		return errors.NewValidationError("Total Hours must be specified and should be greater than zero.")
	}

	if util.IsEmpty(event.LastRegistrationDate) {
		return errors.NewValidationError("Last registration date must be specified.")
	}

	if util.IsEmpty(event.EventStatus) {
		return errors.NewValidationError("Event status must be specified.")
	}

	return nil
}

// SwabhavEventDTO will contain details for get request.
type SwabhavEventDTO struct {
	ID                   uuid.UUID  `json:"id"`
	DeletedAt            *time.Time `json:"-"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	EntryFee             uint       `json:"entryFee"`
	IsOnline             *bool      `json:"isOnline"`
	FromDate             string     `json:"fromDate"`
	ToDate               string     `json:"toDate"`
	FromTime             string     `json:"fromTime"`
	ToTime               string     `json:"toTime"`
	TotalHours           uint       `json:"totalHours"`
	LastRegistrationDate string     `json:"lastRegistrationDate"`
	EventStatus          string     `json:"eventStatus"`
	EventMeetingLink     *string    `json:"eventMeetingLink"`
	IsActive             *bool      `json:"isActive"`
	EventImage           *string    `json:"eventImage"`
	TotalRegistrations   uint       `json:"totalRegistrations"`
	general.Address
}

// TableName will refer it to events.
func (*SwabhavEventDTO) TableName() string {
	return "swabhav_events"
}
