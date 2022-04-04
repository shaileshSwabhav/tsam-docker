package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ResourceLikeService Provide method to Update, Delete, Add, Get Method For ResourceLike.
type ResourceLikeService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewResourceLikeService returns new instance of ResourceLikeService.
func NewResourceLikeService(db *gorm.DB, repository repository.Repository) *ResourceLikeService {
	return &ResourceLikeService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Resource", "Credential",
		},
	}
}

// AddLike will add entry in the resource-like table.
func (service *ResourceLikeService) AddLike(like *res.Like) error {

	// Check if all foreignkey exists.
	err := service.doesForeignKeyExist(like, like.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, like)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateLike will update entry in the resource-like table.
func (service *ResourceLikeService) UpdateLike(like *res.Like) error {

	// Check if all foreignkey exists.
	err := service.doesForeignKeyExist(like, like.UpdatedBy)
	if err != nil {
		return err
	}

	// check if like exist.
	err = service.doesResourceLikeExist(like.TenantID, like.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	// Get createdby
	tempLike := res.Like{}
	err = service.Repository.GetRecordForTenant(uow, like.TenantID, &tempLike,
		repository.Filter("`id` = ?", like.ID), repository.Select("`created_by`"))
	if err != nil {
		uow.RollBack()
		return err
	}

	like.CreatedBy = tempLike.CreatedBy

	err = service.Repository.Save(uow, like)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetResourceLike will return count of like of specified resource.
func (service *ResourceLikeService) GetResourceLike(tenantID, resourceID, credentialID uuid.UUID,
	like *res.LikeDTO) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if resource is exist.
	err = service.doesResourceExist(tenantID, resourceID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(tenantID, credentialID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	// repository.PreloadAssociations(ser.association)
	err = service.Repository.GetRecordForTenant(uow, tenantID, &like,
		repository.Filter("`resource_id` = ? AND `credential_id` = ?", resourceID, credentialID))
	if err != nil {
		uow.RollBack()
		return err
	}

	err = service.Repository.GetCountForTenant(uow, tenantID, &res.Like{}, &like.TotalCount,
		repository.Filter("`resource_id` = ? AND `is_liked` = '1'", resourceID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetAllResourceLikes will return all the like details for specified resource.
func (service *ResourceLikeService) GetAllResourceLikes(tenantID, resourceID uuid.UUID,
	like *[]res.LikeDTO) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if resource is exist.
	err = service.doesResourceExist(tenantID, resourceID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, &like,
		repository.Filter("`resource_id` = ? AND `is_liked` = '1'", resourceID),
		repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	return nil
}

// ==========================================================================================================
// private methods below!
// ==========================================================================================================

// doesForeignKeyExist will check if all foreign keys exist.
func (service *ResourceLikeService) doesForeignKeyExist(like *res.Like, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(like.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(like.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(like.TenantID, like.CredentialID)
	if err != nil {
		return err
	}

	// Check if resource is exist.
	err = service.doesResourceExist(like.TenantID, like.ResourceID)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (service *ResourceLikeService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExists validates credentialID.
func (service *ResourceLikeService) doesCredentialExists(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceExist checks if resource already exist.
func (service *ResourceLikeService) doesResourceExist(tenantID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, res.Resource{},
		repository.Filter("`id`=?", resourceID))
	if err := util.HandleError("Resource not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceLikeExist checks if like already exist.
func (service *ResourceLikeService) doesResourceLikeExist(tenantID, likeID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, res.Like{},
		repository.Filter("`id`=?", likeID))
	if err := util.HandleError("Like not found", exists, err); err != nil {
		return err
	}
	return nil
}
