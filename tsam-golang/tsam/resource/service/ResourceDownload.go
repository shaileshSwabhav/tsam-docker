package service

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/models/general"
	res "github.com/techlabs/swabhav/tsam/models/resource"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// ResourceDownloadService Provide method to Update, Delete, Add, Get Method For ResourceDownload.
type ResourceDownloadService struct {
	DB          *gorm.DB
	Repository  repository.Repository
	association []string
}

// NewResourceDownloadService returns new instance of ResourceDownloadService.
func NewResourceDownloadService(db *gorm.DB, repository repository.Repository) *ResourceDownloadService {
	return &ResourceDownloadService{
		DB:         db,
		Repository: repository,
		association: []string{
			"Resource", "Credential",
		},
	}
}

// AddResourceDownload will add download details.
func (service *ResourceDownloadService) AddResourceDownload(resourceDownload *res.Download) error {

	// Check if all foreign keys are valid
	err := service.doesForeignKeyExist(resourceDownload, resourceDownload.CreatedBy)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, false)

	err = service.Repository.Add(uow, resourceDownload)
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetResourceDownload will return all the downloaded resources and its credential
func (service *ResourceDownloadService) GetResourceDownload(tenantID, resourceID uuid.UUID, resourceDownloads *[]res.DownloadDTO) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if unique resource is added.
	err = service.doesResourceExist(tenantID, resourceID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetAllForTenant(uow, tenantID, resourceDownloads,
		repository.Filter("`resource_id` = ?", resourceID), repository.GroupBy("`credential_id`"),
		repository.PreloadAssociations(service.association))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetResourceCount will return count of download of specified resource
func (service *ResourceDownloadService) GetResourceCount(tenantID, resourceID uuid.UUID, downloadCount *res.DownloadCount) error {

	// Check if tenant exist.
	err := service.doesTenantExists(tenantID)
	if err != nil {
		return err
	}

	// Check if unique resource is added.
	err = service.doesResourceExist(tenantID, resourceID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.DB, true)

	err = service.Repository.GetCountForTenant(uow, tenantID, &res.Download{}, &downloadCount.TotalCount,
		repository.Filter("`resource_id` = ?", resourceID), repository.GroupBy("`credential_id`"))
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

// doesForeignKeyExist will check if all foreign keys exist.
func (service *ResourceDownloadService) doesForeignKeyExist(download *res.Download, credentialID uuid.UUID) error {

	// Check if tenant exist.
	err := service.doesTenantExists(download.TenantID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(download.TenantID, credentialID)
	if err != nil {
		return err
	}

	// Check if credential exist.
	err = service.doesCredentialExists(download.TenantID, download.CredentialID)
	if err != nil {
		return err
	}

	// Check if unique resource is added.
	err = service.doesResourceExist(download.TenantID, download.ResourceID)
	if err != nil {
		return err
	}

	return nil
}

// doesTenantExists validates tenantID.
func (service *ResourceDownloadService) doesTenantExists(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{},
		repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Tenant not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesCredentialExists validates credentialID.
func (service *ResourceDownloadService) doesCredentialExists(tenantID, credentialID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Credential not found", exists, err); err != nil {
		return err
	}
	return nil
}

// doesResourceExist checks if resource already exist.
func (service *ResourceDownloadService) doesResourceExist(tenantID, resourceID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, res.Resource{},
		repository.Filter("`id`=?", resourceID))
	if err := util.HandleError("Resource not found", exists, err); err != nil {
		return err
	}
	return nil
}
