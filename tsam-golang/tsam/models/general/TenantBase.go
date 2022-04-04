package general

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// TenantBase struct contains the common fields and an additional TenantID field for specific entities.
type TenantBase struct {
	ID         uuid.UUID  `gorm:"type:varchar(36);primary_key" json:"id" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" `
	CreatedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	UpdatedAt  *time.Time `json:"-"`
	DeletedBy  uuid.UUID  `gorm:"type:varchar(36)" json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	TenantID   uuid.UUID  `gorm:"type:varchar(36)" json:"tenantID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9"`
	IgnoreHook bool       `gorm:"-" json:"-"`
}

// BeforeCreate will be called before the entity is added to db.
// func (tb *TenantBase) BeforeCreate() {
// 	if tb.IgnoreHook {
// 		return
// 	}
// 	tb.ID = uuid.NewV4()
// 	return
// }

func (tenantBase *TenantBase) BeforeCreate(scope *gorm.Scope) error {
	if tenantBase.IgnoreHook {
		return nil
	}
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid.String())
}
