package community

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

//ModuleConfig Automigrate Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

//NewCommunityModuleConfig Return New Module Config.
func NewCommunityModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db.Set("gorm:table_options", "ENGINE=InnoDB"),
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&Channel{},
		&Discussion{},
		&Reply{},
		&Reaction{},
		&Notification{},
		&NotificationType{},
	}
	for _, community := range models {
		if err := config.DB.AutoMigrate(community).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration Failed ==> %s", err.Error())
		}
	}

	if err := config.DB.Model(Channel{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in channel table, Error %s", err.Error())
	}

	if err := config.DB.Model(Discussion{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in channel table, Error %s", err.Error())
	}

	if err := config.DB.Model(Discussion{}).
		AddForeignKey("channel_id", "community_channels(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add channel foreign key in Discussion Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Discussion{}).
		AddForeignKey("author_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add talent foreign key in Discussion Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Discussion{}).
		AddForeignKey("best_reply_id", "community_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add reply foreign key in Discussion Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reply{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in channel table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reply{}).
		AddForeignKey("discussion_id", "community_discussions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reply Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reply{}).
		AddForeignKey("reply_id", "community_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reply Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reply{}).
		AddForeignKey("replier_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reply Table, Error %s", err.Error())
	}

	// if err := config.DB.Model(Notification{}).
	// 	AddForeignKey("notifier_talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Errorf("failed to add foreign key in Discussion Table, Error %s", err.Error())
	// }

	// if err := config.DB.Model(Notification{}).
	// 	AddForeignKey("notified_talent_id", "talents(id)", "CASCADE", "CASCADE").Error; err != nil {
	// 	log.NewLogger().Errorf("failed to add foreign key in Discussion Table, Error %s", err.Error())
	// }

	if err := config.DB.Model(Reaction{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in channel table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reaction{}).
		AddForeignKey("reactor_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reaction Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reaction{}).
		AddForeignKey("discussion_id", "community_discussions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reaction Table, Error %s", err.Error())
	}

	if err := config.DB.Model(Reaction{}).
		AddForeignKey("reply_id", "community_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Errorf("failed to add foreign key in Reaction Table, Error %s", err.Error())
	}
	log.NewLogger().Info("Community Module Configured.")
}
