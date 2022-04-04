package general

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"time"
)

// Base contains master fields for all entities.
type Base struct {
	ID        uuid.UUID  `gorm:"type:varchar(36);primary_key" json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	CreatedBy uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedBy uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedBy uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// BeforeCreate will be called before the entity is added to db.
func (*Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid.String())
}
