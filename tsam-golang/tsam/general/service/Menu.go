package service

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// MenuService provide methods to update, delete, add, get and get by role id for menu in database
type MenuService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewMenuService returns new instance of menu.
func NewMenuService(db *gorm.DB, repository repository.Repository) *MenuService {
	return &MenuService{
		DB:         db,
		Repository: repository,
	}
}

// AddMenu adds new menu in database.
func (service *MenuService) AddMenu(menu *general.Menu) error {

	// Validate tenant id.
	err := service.doesTenantExist(menu.TenantID)
	if err != nil {
		return err
	}

	// Validate role id.
	if err := service.doesRoleExist(menu.TenantID, menu.RoleID); err != nil {
		return err
	}

	// Validate parent menu id if present.
	if menu.MenuID != nil {
		if err := service.doesMenuExist(*menu.MenuID, menu.TenantID); err != nil {
			return err
		}
		err = service.doesOrderExistForSubMenu(menu.TenantID, *menu.MenuID, menu.RoleID, menu.ID, menu.Order)
		if err != nil {
			return err
		}
	} else {
		err = service.doesOrderExistForMenu(menu.TenantID, menu.ID, menu.RoleID, menu.Order)
		if err != nil {
			return err
		}
	}

	// Check if menu name exists or not.
	err = service.doesMenuNameExist(menu)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Add menu to database.
	if err := service.Repository.Add(uow, menu); err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateMenu updates one menu by specific menu id in database.
func (service *MenuService) UpdateMenu(menu *general.Menu) error {

	// Validate tenant id.
	if err := service.doesTenantExist(menu.TenantID); err != nil {
		return err
	}

	// Validate role id.
	if err := service.doesRoleExist(menu.TenantID, menu.RoleID); err != nil {
		return err
	}

	// Validate parent menu id if present.
	if menu.MenuID != nil {
		if err := service.doesMenuExist(*menu.MenuID, menu.TenantID); err != nil {
			return err
		}
		err := service.doesOrderExistForSubMenu(menu.TenantID, *menu.MenuID, menu.RoleID, menu.ID, menu.Order)
		if err != nil {
			return err
		}
	} else {
		err := service.doesOrderExistForMenu(menu.TenantID, menu.ID, menu.RoleID, menu.Order)
		if err != nil {
			return err
		}
	}

	// Check if menu name exists or not.
	err := service.doesMenuNameExist(menu)
	if err != nil {
		return err
	}

	// Validate menu id.
	if err := service.doesMenuExist(menu.ID, menu.TenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting menu already present in database.
	tempMenu := general.Menu{}

	// Get menu for getting created_by field of menu from database.
	if err := service.Repository.GetForTenant(uow, menu.TenantID, menu.ID, &tempMenu); err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp menu to menu to be updated.
	menu.CreatedBy = tempMenu.CreatedBy

	// Update menu.
	if err := service.Repository.Save(uow, menu); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Menu could not be updated", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// DeleteMenu deletes one menu by specific menu id form database.
func (service *MenuService) DeleteMenu(menu *general.Menu) error {

	// Validate tenant id.
	if err := service.doesTenantExist(menu.TenantID); err != nil {
		return err
	}

	// Validate menu id.
	if err := service.doesMenuExist(menu.ID, menu.TenantID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Update menu for updating deleted_by and deleted_at fields of menu
	if err := service.Repository.UpdateWithMap(uow, menu, map[string]interface{}{
		"DeletedBy": menu.DeletedBy,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", menu.TenantID)); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Menu could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetMenusByRole gets menus by role id form database.
func (service *MenuService) GetMenusByRole(menus *[]general.MenuDTO, roleID uuid.UUID, tenantID uuid.UUID) error {

	// Validate tenant id.
	if err := service.doesTenantExist(tenantID); err != nil {
		return err
	}

	// Validate role id.
	if err := service.doesRoleExist(tenantID, roleID); err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	//get all top level menus by role id
	if err := service.Repository.GetAllInOrderForTenant(uow, tenantID, menus, "`order`",
		repository.Filter("`role_id`=? AND `menu_id` IS NULL", roleID)); err != nil {
		uow.RollBack()
		return err
	}

	// Get sub menus recursively.
	for index := range *menus {
		err := service.getSubMenuByRole(uow, tenantID, (*menus)[index].ID, roleID, &((*menus)[index].SubMenus))
		if err != nil {
			uow.RollBack()
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetMenu returns one menu by specific id.
func (service *MenuService) GetMenu(menu *general.MenuDTO) error {
	// Validate tenant id.
	err := service.doesTenantExist(menu.TenantID)
	if err != nil {
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get one menu by menu id from database.
	err = service.Repository.GetForTenant(uow, menu.TenantID, menu.ID, menu)
	if err != nil {
		uow.RollBack()
		return errors.NewValidationError("Record not found")
	}

	// Get sub menus recursively.
	err = service.getSubMenu(uow, menu.TenantID, menu.ID, &menu.SubMenus)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// getSubMenuByRole will return all the sub-menus for the parent menu by role id.
func (service *MenuService) getSubMenuByRole(uow *repository.UnitOfWork, tenantID, menuID, roleID uuid.UUID,
	subMenus *[]general.MenuDTO) error {

	// Check if menu exists or not in database.
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Menu{},
		repository.Filter("`menu_id`=? AND `role_id`=?", menuID, roleID))
	if err != nil {
		return err
	}
	if exist { // if exists then get its sub menus from database.
		err = service.Repository.GetAllInOrderForTenant(uow, tenantID, subMenus, "`order`",
			repository.Filter("`menu_id`=? AND `role_id`=?", menuID, roleID),
			repository.PreloadAssociations([]string{"ParentMenu"}))
		if err != nil {
			return err
		}

		// Recursice call for all the sub menus.
		for index := range *subMenus {
			err = service.getSubMenuByRole(uow, tenantID, (*subMenus)[index].ID, roleID, &((*subMenus)[index].SubMenus))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getSubMenu will return all the sub-menus for the parent menu.
func (service *MenuService) getSubMenu(uow *repository.UnitOfWork, tenantID, menuID uuid.UUID, subMenus *[]general.MenuDTO) error {

	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Menu{},
		repository.Filter("`menu_id`=?", menuID))
	if err != nil {
		return err
	}
	if exist {

		err = service.Repository.GetAllInOrderForTenant(uow, tenantID, subMenus, "`order`",
			repository.Filter("`menu_id`=?", menuID))
		if err != nil {
			return err
		}

		// rescursive call
		for index := range *subMenus {
			err = service.getSubMenu(uow, tenantID, (*subMenus)[index].ID, &((*subMenus)[index].SubMenus))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// doesRoleExist validates if role exists or not in database.
func (service *MenuService) doesRoleExist(tenantID uuid.UUID, roleID uuid.UUID) error {
	//check if role exists or not
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Role{},
		repository.Filter("`id` = ?", roleID))
	if err := util.HandleError("Invalid role ID", exists, err); err != nil {
		return err
	}

	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *MenuService) doesTenantExist(tenantID uuid.UUID) error {
	//check if tenant(parent tenant) exists or not
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		return err
	}
	return nil
}

// doesMenuNameExist checks if menus's menu name already exists (if parent menu exists then checks only in
//that parent menus' child menus) in database
func (service *MenuService) doesMenuNameExist(menu *general.Menu) error {
	// Create query processor according to parent menu exists or not.
	var queryProcessor repository.QueryProcessor
	queryProcessor = repository.Filter("`menu_name`=? AND `menu_id` IS NULL AND `role_id`=? AND `id`!=?", menu.MenuName, menu.RoleID, menu.ID)
	if menu.MenuID != nil { //check for same parent id also
		queryProcessor = repository.Filter("`menu_name`=? AND `menu_id`=? AND `role_id`=? AND `id`!=?", menu.MenuName, menu.MenuID, menu.RoleID, menu.ID)
	}

	// Check for same menu name conflict.
	exists, err := repository.DoesRecordExistForTenant(service.DB, menu.TenantID, &general.Menu{}, queryProcessor)
	if err := util.HandleIfExistsError("Menu name exists", exists, err); err != nil {
		return err
	}
	return nil
}

// doesMenuExist validates if menu exists or not in database.
func (service *MenuService) doesMenuExist(menuID uuid.UUID, tenantID uuid.UUID) error {
	// Check menu exists or not.
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Menu{},
		repository.Filter("`id` = ?", menuID))
	if err := util.HandleError("Invalid menu ID", exists, err); err != nil {
		return err
	}
	return nil
}

// sortMenusByOrder sorts all the menus by order recursively.
// func (service *MenuService) sortMenusByOrder(menus *[]general.Menu) {
// 	sort.SliceStable(*menus, func(p, q int) bool {
// 		return (*menus)[p].Order < (*menus)[q].Order
// 	})
// 	for _, menu := range *menus {
// 		if len(menu.SubMenus) != 0 {
// 			service.sortMenusByOrder(&menu.SubMenus)
// 		}
// 	}
// }

// Check if menu order exists for menu.
func (service *MenuService) doesOrderExistForMenu(tenantID, id, roleID uuid.UUID, order int) error {
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Menu{},
		repository.Filter("`role_id`=? AND `order`=? AND `menu_id` IS NULL AND `id` !=?", roleID, order, id))
	if err := util.HandleIfExistsError("Order for menu already exist", exist, err); err != nil {
		return err
	}

	return nil
}

// Check if menu order exists for submenu.
func (service *MenuService) doesOrderExistForSubMenu(tenantID, menuID, roleID, id uuid.UUID, order int) error {
	exist, err := repository.DoesRecordExistForTenant(service.DB, tenantID, &general.Menu{},
		repository.Filter("`order`=? AND `menu_id`=? And `role_id`=? AND `id`!=?", order, menuID, roleID, id))
	if err := util.HandleIfExistsError("Order for menu already exists", exist, err); err != nil {
		return err
	}

	return nil
}
