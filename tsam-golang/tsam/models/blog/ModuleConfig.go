package blog

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/techlabs/swabhav/tsam/log"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	DB *gorm.DB
}

// NewBlogModuleConfig Return New blog Module Config.
func NewBlogModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (module *ModuleConfig) TableMigration(wg *sync.WaitGroup) {

	var models []interface{} = []interface{}{
		&BlogTopic{},
		&Blog{},
		&BlogReply{},
		&BlogReaction{},
		&BlogView{},
		&BlogNotification{},
	}

	for _, model := range models {
		if err := module.DB.Debug().AutoMigrate(model).Error; err != nil {
			log.NewLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	//********************************** BLOG TOPIC FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&BlogTopic{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BLOG FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&Blog{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Author.
	if err := module.DB.Model(&Blog{}).
		AddForeignKey("author_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//**********************************Blog's maps' foreign keys*******************************

	// Blog topic map.
	if err := module.DB.Model(&BlogBlogTopics{}).
		AddForeignKey("blog_id", "blogs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	if err := module.DB.Model(&BlogBlogTopics{}).
		AddForeignKey("blog_topic_id", "blog_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BLOG REPLY FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&BlogReply{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Replier.
	if err := module.DB.Model(&BlogReply{}).
		AddForeignKey("replier_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Reply.
	if err := module.DB.Model(&BlogReply{}).
		AddForeignKey("reply_id", "blog_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Blog.
	if err := module.DB.Model(&BlogReply{}).
		AddForeignKey("blog_id", "blogs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BLOG REACTION FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&BlogReaction{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Reactor.
	if err := module.DB.Model(&BlogReaction{}).
		AddForeignKey("reactor_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Reply.
	if err := module.DB.Model(&BlogReaction{}).
		AddForeignKey("reply_id", "blog_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Blog.
	if err := module.DB.Model(&BlogReaction{}).
		AddForeignKey("blog_id", "blogs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//********************************** BLOG VIEWER FOREIGN KEYS *******************************

	// Tenant.
	if err := module.DB.Model(&BlogView{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Reactor.
	if err := module.DB.Model(&BlogView{}).
		AddForeignKey("viewer_id", "credentials(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	// Blog.
	if err := module.DB.Model(&BlogView{}).
		AddForeignKey("blog_id", "blogs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	log.NewLogger().Info("Blog Module Configured.")

	//********************************** BLOG Notification FOREIGN KEYS *******************************
	// Tenant.
	if err := module.DB.Model(&BlogNotification{}).
		AddForeignKey("tenant_id", "tenants(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//BlogTopic
	if err := module.DB.Model(&BlogNotification{}).
		AddForeignKey("blog_topic_id", "blog_topics(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}

	//BlogReaction
	if err := module.DB.Model(&BlogNotification{}).
		AddForeignKey("blog_reaction_id", "blog_reactions(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())

	}
	// Reply.
	if err := module.DB.Model(&BlogNotification{}).
		AddForeignKey("blog_reply_id", "blog_replies(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
	// blog ID.
	if err := module.DB.Model(&BlogNotification{}).
		AddForeignKey("blog_id", "blogs(id)", "CASCADE", "CASCADE").Error; err != nil {
		log.NewLogger().Error(err.Error())
	}
}
