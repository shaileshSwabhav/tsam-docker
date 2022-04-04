package general

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/util"
)

// Menu contains all information to create a dynamic navigation bar with links to the respective pages.
// MenuID is a self referencing foreign key to determine the position of the menu item.
type Menu struct {
	TenantBase
	Permission `json:"permission"`
	Order      int        `json:"order" gorm:"type:tinyint(2)"`
	MenuName   string     `json:"menuName" example:"Company" gorm:"varchar(30)"`
	MenuID     *uuid.UUID `json:"menuID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"FOREIGNKEY:ID;type:varchar(36)"`
	RoleID     uuid.UUID  `json:"roleID" example:"cfe25758-f5fe-48f0-874d-e72cd4edd9b9" gorm:"type:varchar(36)"`
	URL        *string    `json:"url" example:"./company/" gorm:"varchar(100)"`
	SubMenus   []Menu     `json:"menus"`
	IsVisible  bool       `json:"isVisible"`
}

// MenuDTO contains all information to create a dynamic navigation bar with links to the respective pages.
// MenuID is a self referencing foreign key to determine the position of the menu item.
type MenuDTO struct {
	TenantBase
	Permission `json:"permission"`
	Order      int        `json:"order"`
	MenuName   string     `json:"menuName"`
	MenuID     *uuid.UUID `json:"menuID"`
	RoleID     uuid.UUID  `json:"roleID"`
	URL        *string    `json:"url"`
	SubMenus   []MenuDTO     `json:"menus" gorm:"foreignkey:MenuID"`
	IsVisible  bool       `json:"isVisible"`
	ParentMenu *ParentMenu  `json:"parentMenu" gorm:"foreignkey:MenuID"`
}

// TableName will name the table of ParentMenu model as "menus".
func (*MenuDTO) TableName() string {
	return "menus"
}

// Permission contains different permissions(access rights) like add,update & delete.
type Permission struct {
	Add    bool `json:"add" example:"true"`
	Update bool `json:"update" example:"false"`
	Delete bool `json:"delete" example:"true"`
}

// Validate Validates menu fields.
func (menu *Menu) Validate() error {
	// Check if menu name is blank ot not.
	if util.IsEmpty(menu.MenuName) {
		return errors.NewValidationError("Menu name must be specified")
	}

	// Check if order is lesser than or equal to 0.
	if menu.Order <= 0 {
		return errors.NewValidationError("Menu order cannot be 0 or less")
	}

	// Check if role id exists or not.
	if !util.IsUUIDValid(menu.RoleID) {
		return errors.NewValidationError("Role id must be specified")
	}

	// Check if parent menu id is valid or not.
	if menu.MenuID != nil && !util.IsUUIDValid(*menu.MenuID) {
		return errors.NewValidationError("Parent Menu ID is not valid")
	}

	// Check if all submenus have unique order.
	menuMap := make(map[int]uint)
	for _, subMenu := range menu.SubMenus {
		menuMap[subMenu.Order]++

		if menuMap[subMenu.Order] > 1 {
			return errors.NewValidationError("Sub-menus with same order not allowed")
		}
	}

	// If menu contains menus then validate all child menus
	if len(menu.SubMenus) > 0 {
		for _, m := range menu.SubMenus {
			if err := m.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

// ParentMenu is used to get basic info of parent menu.
type ParentMenu struct {
	ID         uuid.UUID  `json:"id"`
	MenuName   string     `json:"menuName"`
}

// TableName will name the table of ParentMenu model as "menus".
func (*ParentMenu) TableName() string {
	return "menus"
}
