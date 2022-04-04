package test

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewTestModuleConfig Create New Test Module Config
func NewTestModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	// Table List
	var models []interface{} = []interface{}{
		Option{},
		Question{},
	}

	// Table Migrant
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	// Add Foreign Key
	config.DB.Model(Option{}).
		AddForeignKey("question_id", "questions(id)", "CASCADE", "CASCADE")

	log.NewLogger().Info("Test Module Configured.")
}
