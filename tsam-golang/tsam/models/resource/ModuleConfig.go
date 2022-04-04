package resource

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewResourceModuleConfig Return New Module Config.
func NewResourceModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&Resource{},
		&Download{},
		&Like{},
		// &Book{},
	}

	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	// Add Foreign Key

	//************************************* RESOURCE FOREIGN KEY ******************************

	if err := config.DB.Model(&Resource{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Resource{}).
		AddForeignKey("technology_id", "technologies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//************************************* RESOURCE DOWNLOAD FOREIGN KEY ******************************

	if err := config.DB.Model(&Download{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Download{}).
		AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Download{}).
		AddForeignKey("credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//************************************* RESOURCE LIKE FOREIGN KEY ******************************

	if err := config.DB.Model(&Like{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Like{}).
		AddForeignKey("resource_id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := config.DB.Model(&Like{}).
		AddForeignKey("credential_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//************************************* BOOK FOREIGN KEY ******************************

	// if err := config.DB.Model(&Book{}).
	// 	AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

	// if err := config.DB.Model(&Book{}).
	// 	AddForeignKey("id", "resources(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Error(err.Error())
	// }

}
