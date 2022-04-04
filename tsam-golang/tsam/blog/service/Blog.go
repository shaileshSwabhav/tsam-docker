package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	blg "github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/models/general"
	"github.com/techlabs/swabhav/tsam/repository"
	"github.com/techlabs/swabhav/tsam/util"
)

// BlogService provide method to update, delete, add, get method for blog.
type BlogService struct {
	DB         *gorm.DB
	Repository repository.Repository
}

// NewBlogService creates new instance of BlogService.
func NewBlogService(db *gorm.DB, repository repository.Repository) *BlogService {
	return &BlogService{
		DB:         db,
		Repository: repository,
	}
}

// multipleBlogAssociationNames provides preload associations array for blog.
var multipleBlogAssociationNames []string = []string{
	"BlogTopics", "Author", "Author.Role",
}

// BlogAssociationNames provides preload associations array for blog.
var BlogAssociationNames []string = []string{
	"BlogTopics", "Author", "Author.Role", "Reactions", "Reactions.Reactor", "Reactions.Reactor.Role",
}

// AddBlog adds new blog to database.
func (service *BlogService) AddBlog(blog *blg.Blog) error {

	// Validate tenant id.
	err := service.doesTenantExist(blog.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same blog title exists.
	err = service.doesBlogTitleExist(blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blog.CreatedBy, blog.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate author ID.
	err = service.doesCredentialExist(blog.AuthorID, blog.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	//  Creating unit of work.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Add blog.
	err = service.Repository.Add(uow, blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// UpdateBlog updates blog in database.
func (service *BlogService) UpdateBlog(blog *blg.Blog) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blog.TenantID)
	if err != nil {
		return err
	}

	// Validate blog ID.
	err = service.doesBlogExist(blog.ID, blog.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(blog.UpdatedBy, blog.TenantID)
	if err != nil {
		return err
	}

	// Validate author ID.
	err = service.doesCredentialExist(blog.AuthorID, blog.TenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Validate if same blog title exists.
	err = service.doesBlogTitleExist(blog)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create bucket for getting blog already present in database.
	tempBlog := blg.Blog{}

	// Get blog for getting created_by field of blog from database.
	if err := service.Repository.GetForTenant(uow, blog.TenantID, blog.ID, &tempBlog); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Give created_by id from temp blog to blog to be updated.
	blog.CreatedBy = tempBlog.CreatedBy

	// Update blog associations.
	err = service.updateBlogAssociation(uow, blog)
	if err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Blog could not be updated", http.StatusInternalServerError)
	}

	// Make blog topics nil so that it is not inserted again.
	blog.BlogTopics = nil

	// Update blog.
	err = service.Repository.Save(uow, blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// DeleteBlog delete blog from database.
func (service *BlogService) DeleteBlog(blog *blg.Blog) error {
	credentialID := blog.DeletedBy

	// Validate tenant ID.
	err := service.doesTenantExist(blog.TenantID)
	if err != nil {
		return err
	}

	// Validate blog ID.
	err = service.doesBlogExist(blog.ID, blog.TenantID)
	if err != nil {
		return err
	}

	// Validate credential ID.
	err = service.doesCredentialExist(credentialID, blog.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	//  Delete blog association from database.
	if err := service.deleteBlogAssociation(uow, blog, credentialID); err != nil {
		uow.RollBack()
		return errors.NewHTTPError("Failed to delete blog", http.StatusInternalServerError)
	}

	// Update blog for updating deleted_by and deleted_at fields of blog.
	if err := service.Repository.UpdateWithMap(uow, blog, map[string]interface{}{
		"DeletedBy": credentialID,
		"DeletedAt": time.Now(),
	},
		repository.Filter("`tenant_id`=?", blog.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog could not be deleted", http.StatusInternalServerError)
	}

	uow.Commit()
	return nil
}

// GetBlog returns one blog by id.
func (service *BlogService) GetBlog(blog *blg.DTO) error {

	// Validate tenant ID.
	err := service.doesTenantExist(blog.TenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Get blog by id from database.
	err = service.Repository.GetForTenant(uow, blog.TenantID, blog.ID, blog,
		repository.PreloadAssociations(BlogAssociationNames),
	)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Create bucket to get count og blog viewers.
	var totalCount int

	// Count blog viewers.
	err = service.Repository.GetCount(uow, &blg.BlogView{}, &totalCount, 
		repository.Filter("`blog_id`=?", blog.ID),
		repository.Filter("`tenant_id`=? AND `deleted_at` IS NULL", blog.TenantID))
	if err != nil {
		uow.RollBack()
		return err
	}
	blog.BlogViewCount = uint16(totalCount)

	uow.Commit()
	return nil
}

// GetBlogs returns all blogs.
func (service *BlogService) GetBlogs(blogs *[]blg.DTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create quesry precessors.
	queryProcessors := service.addSearchQueries(form, tenantID)
	queryProcessors = append(queryProcessors, 
		repository.PreloadAssociations(multipleBlogAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	// Get all blogs.
	err = service.Repository.GetAllInOrderForTenant(uow, tenantID, blogs, "`published_date`", queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetLatestBlogSnippets returns all latest blog snippets.
func (service *BlogService) GetLatestBlogSnippets(blogs *[]blg.SnippetDTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create quesry precessors.
	queryProcessors := service.addSearchQueries(form, tenantID)
	queryProcessors = append(queryProcessors, 
		repository.Select("SUM(IF(blog_reactions.`is_clap` = 1, 1, 0)) as clap_count, SUM(IF(blog_reactions.`is_clap` = 0, 1, 0)) as slap_count, blogs.*"),
		repository.Join("LEFT JOIN blog_reactions on blogs.`id` = blog_reactions.`blog_id` AND blogs.`tenant_id` = blog_reactions.`tenant_id`"),
		repository.Filter("blog_reactions.`deleted_at` IS NULL"),
		repository.Filter("blogs.`is_verified`=1 AND blogs.`is_published`=1"),
		repository.Filter("blogs.`tenant_id`=?", tenantID),
		repository.GroupBy("blogs.`id`"),
		repository.OrderBy("blogs.`published_date` DESC"),
		repository.PreloadAssociations(multipleBlogAssociationNames),
		repository.Paginate(limit, offset, totalCount),
	)


	// Get all blogs.
	err = service.Repository.GetAll(uow, blogs, queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// GetTrendingBlogSnippets returns all trending blog snippets.
func (service *BlogService) GetTrendingBlogSnippets(blogs *[]blg.SnippetDTO, tenantID uuid.UUID,
	form url.Values, limit, offset int, totalCount *int) error {

	// Validate tenant id.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}

	// Start new transcation.
	uow := repository.NewUnitOfWork(service.DB, true)

	// Create quesry precessors.
	queryProcessors := service.addSearchQueries(form, tenantID)
	queryProcessors = append(queryProcessors, 
		repository.Select("SUM(IF(blog_reactions.`is_clap` = 1, 1, 0)) as clap_count, blogs.*"),
		repository.Join("JOIN blog_reactions on blogs.`id` = blog_reactions.`blog_id`"),
		repository.Filter("blog_reactions.`created_at` between DATE_SUB(NOW(),INTERVAL 2 WEEK) and NOW()"),
		repository.Filter("blog_reactions.`deleted_at` IS NULL"),
		repository.Filter("blog_reactions.`is_clap`=1"),
		repository.Filter("blogs.`is_verified`=1 AND blogs.`is_published`=1"),
		repository.Filter("blogs.`tenant_id`=? AND blog_reactions.`tenant_id`=?", tenantID, tenantID),
		repository.GroupBy("blogs.`id`"),
		repository.OrderBy("clap_count DESC"),
		repository.PreloadAssociations(multipleBlogAssociationNames),
		repository.Paginate(limit, offset, totalCount))

	// Get all blogs.
	err = service.Repository.GetAll(uow, blogs, queryProcessors...)
	if err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return err
	}

	// Count total claps and slaps of all blog snippets.
	for index := range *blogs {

		blogsSlapsClaps := blg.ClapSlapCountModel{}

		err = service.Repository.GetAll(uow, &blogsSlapsClaps,
			repository.Select("SUM(IF(blog_reactions.`is_clap` = 1, 1, 0)) as clap_count, SUM(IF(blog_reactions.`is_clap` = 0, 1, 0)) as slap_count"),
			repository.Join("JOIN blog_reactions on blogs.`id` = blog_reactions.`blog_id`"),
			repository.Filter("blog_reactions.`deleted_at` IS NULL AND blogs.`deleted_at` IS NULL"),
			repository.Filter("blogs.`id`=?", (*blogs)[index].ID),
			repository.Filter("blogs.`tenant_id`=? AND blog_reactions.`tenant_id`=?", tenantID, tenantID))
		if err != nil {
			log.NewLogger().Error(err.Error())
			uow.RollBack()
			return err
		}
		(*blogs)[index].ClapCount = uint16(blogsSlapsClaps.ClapCount)
		(*blogs)[index].SlapCount = uint16(blogsSlapsClaps.SlapCount)
	}

	uow.Commit()
	return nil
}

// UpdateBlogFlags will update the flags of specified blog record in the table.
func (service *BlogService) UpdateBlogFlags(blog *blg.FlagsUpdateModel, tenantID, credentialID uuid.UUID) error {

	// Check if tenant record exist.
	err := service.doesTenantExist(tenantID)
	if err != nil {
		return err
	}

	// Check if credential record exist.
	err = service.doesCredentialExist(credentialID, tenantID)
	if err != nil {
		return err
	}

	// Check if blog record exist.
	err = service.doesBlogExist(blog.ID, tenantID)
	if err != nil {
		return err
	}

	// Start new transaction.
	uow := repository.NewUnitOfWork(service.DB, false)

	// Create variable for pulished date as current date.
	currentTime := time.Now()
	publishedDateInString := currentTime.Format("2006-01-02")
	publishedDate := &publishedDateInString

	// If blog is being verified then dont give published date.
	if !blog.IsVerified || !blog.IsPublished{
		publishedDate = nil 
	}

	// Update flags.
	err = service.Repository.UpdateWithMap(uow, &blg.Blog{}, map[interface{}]interface{}{
		"IsVerified": blog.IsVerified,
		"IsPublished": blog.IsPublished,
		"UpdatedBy":   credentialID,
		"PublishedDate": publishedDate,
	}, repository.Filter("`id`=?", blog.ID))
	if err != nil {
		uow.RollBack()
		return err
	}

	uow.Commit()
	return nil
}

// deleteBlogAssociation deletes blog's associations.
func (service *BlogService) deleteBlogAssociation(uow *repository.UnitOfWork, blog *blg.Blog, credentialID uuid.UUID) error {

	//********************************************** Blog Reply ***************************************************

	blogReplies := []blg.BlogReply{}

	// Get blog replies from database.
	if err := service.Repository.GetAllForTenant(uow, blog.TenantID, &blogReplies,
		repository.Filter("`blog_id`=?", blog.ID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Delete replies' reactions, sub replies and its reactions from database.
	for _, blogReply := range blogReplies {
		if err := service.DeleteReplyAssociation(uow, &blogReply, credentialID); err != nil {
			uow.RollBack()
			return errors.NewHTTPError("Failed to delete blog", http.StatusInternalServerError)
		}
	}

	// Deleting replies.
	if err := service.Repository.UpdateWithMap(uow, &blg.BlogReply{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`blog_id`=? AND `tenant_id`=?", blog.ID, blog.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog could not be deleted", http.StatusInternalServerError)
	}

	//********************************************** Blog Reaction ***************************************************
	
	if err := service.Repository.UpdateWithMap(uow, &blg.BlogReaction{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`blog_id`=? AND `tenant_id`=?", blog.ID, blog.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// DeleteReplyAssociation deletes reply's associations.
func (service *BlogService) DeleteReplyAssociation(uow *repository.UnitOfWork, blogReply *blg.BlogReply, credentialID uuid.UUID) error {

	//********************************************** Sub Reply Reactions ***************************************************

	subReplies := []blg.BlogReply{}

	// Get sub replies.
	if err := service.Repository.GetAllForTenant(uow, blogReply.TenantID, &subReplies,
		repository.Filter("`reply_id`=?", blogReply.ID)); err != nil {
		uow.RollBack()
		log.NewLogger().Error(err.Error())
		return errors.NewValidationError("Record not found")
	}

	// Collect all sub replies' ids in a variable.
	var subReplyIDs []uuid.UUID
	for _, subReply := range subReplies {
		subReplyIDs = append(subReplyIDs, subReply.ID)
	}

	// Deleting sub replies' reactions.
	if err := service.Repository.UpdateWithMap(uow, &blg.BlogReaction{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`reply_id` IN (?) AND `tenant_id`=?", subReplyIDs, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	// Deleting sub replies.
	if err := service.Repository.UpdateWithMap(uow, &blg.BlogReply{},
		map[string]interface{}{
			"DeletedBy": credentialID,
			"DeletedAt": time.Now(),
		},
		repository.Filter("`reply_id`=? AND `tenant_id`=?", blogReply.ID, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}
	
	//********************************************** Blog Reply Reaction ***************************************************
	if err := service.Repository.UpdateWithMap(uow, &blg.BlogReaction{},
		map[string]interface{}{"DeletedBy": credentialID, "DeletedAt": time.Now()},
		repository.Filter("`reply_id`=? AND `tenant_id`=?", blogReply.ID, blogReply.TenantID)); err != nil {
		log.NewLogger().Error(err.Error())
		uow.RollBack()
		return errors.NewHTTPError("Blog reply could not be deleted", http.StatusInternalServerError)
	}

	return nil
}

// addSearchQueries adds all search queries by comparing with the blog data.
func (service *BlogService) addSearchQueries(requestForm url.Values, tenantID uuid.UUID) []repository.QueryProcessor {
	if len(requestForm) == 0 {
		return nil
	}

	var columnNames []string
	var conditions []string
	var operators []string
	var values []interface{}

	var queryProcesors []repository.QueryProcessor

	// Title.
	if _, ok := requestForm["title"]; ok {
		util.AddToSlice("`title`", "LIKE ?", "AND", "%"+requestForm.Get("title")+"%", &columnNames, &conditions, &operators, &values)
	}

	// Time to read(greater than or equal to) search.
	if timeToRead, ok := requestForm["timeToRead"]; ok {
		util.AddToSlice("`time_to_read`", ">=?", "AND", timeToRead, &columnNames, &conditions, &operators, &values)
	}

	// Is verified search.
	if isVerified, ok := requestForm["isVerified"]; ok {
		util.AddToSlice("`is_verified`", "=?", "AND", isVerified, &columnNames, &conditions, &operators, &values)
	}

	// Is published search.
	if isPublished, ok := requestForm["isPublished"]; ok {
		util.AddToSlice("`is_published`", "=?", "AND", isPublished, &columnNames, &conditions, &operators, &values)
	}

	// Published date from date.
	if publishedFromDate, ok := requestForm["publishedFromDate"]; ok {
		util.AddToSlice("`published_date`", ">= ?", "AND", publishedFromDate, &columnNames, &conditions, &operators, &values)
	}

	// Published date to date.
	if publishedEndDate, ok := requestForm["publishedEndDate"]; ok {
		util.AddToSlice("`published_date`", "<= ?", "AND", publishedEndDate, &columnNames, &conditions, &operators, &values)
	}

	// Author id.
	if _, ok := requestForm["authorID"]; ok {
		util.AddToSlice("author_id", "= ?", "AND", requestForm.Get("authorID"), &columnNames, &conditions, &operators, &values)
	}

	// Blog topic search.
	if blogTopicIDs, ok := requestForm["blogTopicIDs"]; ok {
		fmt.Println(";aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		fmt.Println(blogTopicIDs)
		queryProcesors = append(queryProcesors,
			repository.Join("JOIN blogs_blog_topics ON blogs.`id` = blogs_blog_topics.`blog_id`"))

		util.AddToSlice("blogs_blog_topics.`blog_topic_id`", "IN(?)", "AND", blogTopicIDs,
			&columnNames, &conditions, &operators, &values)
	}

	// Group by blog id and add all filters
	queryProcesors = append(queryProcesors,
		repository.FilterWithOperator(columnNames, conditions, operators, values))
	return queryProcesors
}

// updateBlogAssociation updates blog's associations.
func (service *BlogService) updateBlogAssociation(uow *repository.UnitOfWork, enquiry *blg.Blog) error {
	
	// Replace technologies of blog.
	if err := service.Repository.ReplaceAssociations(uow, enquiry, "BlogTopics",
		enquiry.BlogTopics); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogTitleExist returns true if the bog title already exists for the blog in database.
func (service *BlogService) doesBlogTitleExist(blog *blg.Blog) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, blog.TenantID, &blg.Blog{},
		repository.Filter("`title`=? AND `id`!=?", blog.Title, blog.ID))
	if err := util.HandleIfExistsError("Title : "+blog.Title+" exists", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesTenantExist validates if tenant exists or not in database.
func (service *BlogService) doesTenantExist(tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExist(service.DB, general.Tenant{}, repository.Filter("`id` = ?", tenantID))
	if err := util.HandleError("Invalid tenant ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesBlogExist validates if blog exists or not in database.
func (service *BlogService) doesBlogExist(blogID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, blg.Blog{},
		repository.Filter("`id` = ?", blogID))
	if err := util.HandleError("Invalid blog ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}

// doesCredentialExist validates if credential exists or not in database.
func (service *BlogService) doesCredentialExist(credentialID uuid.UUID, tenantID uuid.UUID) error {
	exists, err := repository.DoesRecordExistForTenant(service.DB, tenantID, general.Credential{},
		repository.Filter("`id` = ?", credentialID))
	if err := util.HandleError("Invalid credential ID", exists, err); err != nil {
		log.NewLogger().Error(err.Error())
		return err
	}
	return nil
}
