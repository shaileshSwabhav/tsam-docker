package blog

import (
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/models/list"
	"github.com/techlabs/swabhav/tsam/util"
)

//************************************* ADD/ UPDATE MODEL *********************************************************

// Blog contains add update fields required for blog.
type Blog struct {
	general.TenantBase

	// Related table IDs.
	AuthorID uuid.UUID `json:"authorID" gorm:"type:varchar(36)"`

	// Other fields.
	Content       string  `json:"content" gorm:"type:varchar(5000)"`
	Title         string  `json:"title" gorm:"type:varchar(200)"`
	Description   string  `json:"description" gorm:"type:varchar(300)"`
	TimeToRead    uint16  `json:"timeToRead" gorm:"type:smallint(4)"`
	PublishedDate *string `json:"publishedDate" gorm:"type:date"`
	BannerImage              *string  `json:"bannerImage" gorm:"type:varchar(200)"`

	// Flags.
	IsVerified  bool `json:"isVerified"`
	IsPublished bool `json:"isPublished"`

	// Maps.
	BlogTopics []*BlogTopic `json:"blogTopics" gorm:"many2many:blogs_blog_topics;ASSOCIATION_AUTOCREATE:false;association_autoupdate:false"`
}

// Validate validates compulsary fields of Blog.
func (blog *Blog) Validate() error {

	// Check if content is blank or not.
	if util.IsEmpty(blog.Content) {
		return errors.NewValidationError("Content must be specified")
	}

	// Content maximum characters.
	if len(blog.Content) > 5000 {
		return errors.NewValidationError("Content can have maximum 5000 characters")
	}

	// Check if title is blank or not.
	if util.IsEmpty(blog.Title) {
		return errors.NewValidationError("Title must be specified")
	}

	// Title maximum characters.
	if len(blog.Title) > 200 {
		return errors.NewValidationError("Title can have maximum 200 characters")
	}

	// Check if description is blank or not.
	if util.IsEmpty(blog.Description) {
		return errors.NewValidationError("Description must be specified")
	}

	// Description maximum characters.
	if len(blog.Description) > 300 {
		return errors.NewValidationError("Description can have maximum 300 characters")
	}

	// Time To Read minimum.
	if blog.TimeToRead <= 0 {
		return errors.NewValidationError("Time To Read can be minimum 1")
	}

	// Time To Read maximum.
	if blog.TimeToRead > 240 {
		return errors.NewValidationError("Time To Read can be maximum 240")
	}

	// Check if published date is empty or not.
	if blog.PublishedDate != nil && util.IsEmpty(*blog.PublishedDate) {
		return errors.NewValidationError("Published Date must be specified")
	}

	return nil
}

// ===========Defining many to many structs===========

// BlogBlogTopics is the map of blog and blog topic.
type BlogBlogTopics struct {
	BlogID      uuid.UUID `gorm:"type:varchar(36)"`
	BlogTopicID uuid.UUID `gorm:"type:varchar(36)"`
}

// TableName defines table name of the struct.
func (*BlogBlogTopics) TableName() string {
	return "blogs_blog_topics"
}

//************************************* DTO MODEL *************************************************************

// DTO contains all fields required for blog.
type DTO struct {
	general.TenantBase

	// Single model.
	Author   list.Credential `json:"author" gorm:"foreignkey:AuthorID"`
	AuthorID *uuid.UUID      `json:"-"`

	// Other fields.
	Content       string  `json:"content"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	TimeToRead    uint16  `json:"timeToRead"`
	PublishedDate *string `json:"publishedDate"`
	BannerImage              *string  `json:"bannerImage"`
	BlogViewCount uint16  `json:"blogViewCount"`

	// Flags.
	IsVerified  bool `json:"isVerified"`
	IsPublished bool `json:"isPublished"`

	// Maps.
	BlogTopics []*BlogTopic `json:"blogTopics" gorm:"many2many:blogs_blog_topics;association_jointable_foreignkey:blog_topic_id;jointable_foreignkey:blog_id"`

	// Mutiple field.
	Reactions []BlogReactionDTO  `json:"reactions" gorm:"foreignkey:BlogID"`
}

// TableName defines table name of the struct.
func (*DTO) TableName() string {
	return "blogs"
}

// SnippetDTO contains all fields required for blog snippet.
type SnippetDTO struct {
	general.TenantBase

	// Single model.
	Author   list.Credential `json:"author" gorm:"foreignkey:AuthorID"`
	AuthorID *uuid.UUID      `json:"-"`

	// Other fields.
	Title         string  `json:"title" gorm:"type:varchar(200)"`
	Description   string  `json:"description" gorm:"type:varchar(300)"`
	TimeToRead    uint16  `json:"timeToRead" gorm:"type:smallint(4)"`
	PublishedDate *string `json:"publishedDate" gorm:"type:date"`
	BannerImage              *string  `json:"bannerImage"`
	ClapCount uint16  `json:"clapCount"`
	SlapCount uint16  `json:"slapCount"`

	// Maps.
	BlogTopics []*BlogTopic `json:"blogTopics" gorm:"many2many:blogs_blog_topics;association_jointable_foreignkey:blog_topic_id;jointable_foreignkey:blog_id"`
}

// TableName defines table name of the struct.
func (*SnippetDTO) TableName() string {
	return "blogs"
}

// FlagsUpdateModel is used for updating the flags of blog.
type FlagsUpdateModel struct {
	ID uuid.UUID `json:"id"`

	// Flags.
	IsVerified  bool `json:"isVerified"`
	IsPublished bool `json:"isPublished"`
}

// ClapSlapCountModel is used for getting the clap and slap count of blogs.
type ClapSlapCountModel struct {
	ClapCount uint16 
	SlapCount uint16  
}

// TableName defines table name of the struct.
func (*ClapSlapCountModel) TableName() string {
	return "blogs"
}
