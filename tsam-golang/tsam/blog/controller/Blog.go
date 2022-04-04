package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/blog/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	blg "github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BlogController provides method to update, delete, add, get For blog.
type BlogController struct {
	BlogService *service.BlogService
}

// NewBlogController creates new instance of BlogController.
func NewBlogController(blogService *service.BlogService) *BlogController {
	return &BlogController{
		BlogService: blogService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *BlogController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all blogs with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/blog/limit/{limit}/offset/{offset}",
		controller.GetBlogs).Methods(http.MethodGet)

	// Get all latest blog snippets with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/blog-snippet-latest/limit/{limit}/offset/{offset}",
		controller.GetLatestBlogSnippets).Methods(http.MethodGet)

	// Get all trendong blog snippets with limit and offset.
	router.HandleFunc("/tenant/{tenantID}/blog-snippet-trending/limit/{limit}/offset/{offset}",
		controller.GetTrendingBlogSnippets).Methods(http.MethodGet)

	// Get one blog by id.
	router.HandleFunc("/tenant/{tenantID}/blog/{blogID}",
		controller.GetBlog).Methods(http.MethodGet)

	// Add one blog.
	router.HandleFunc("/tenant/{tenantID}/blog/credential/{credentialID}",
		controller.AddBlog).Methods(http.MethodPost)

	// Update one blog.
	router.HandleFunc("/tenant/{tenantID}/blog/{blogID}/credential/{credentialID}",
		controller.UpdateBlog).Methods(http.MethodPut)

	// Update is flags of blog.
	router.HandleFunc("/tenant/{tenantID}/blog-flags-update/credential/{credentialID}",
		controller.UpdateBlogFlags).Methods(http.MethodPut)

	// Delete one blog.
	router.HandleFunc("/tenant/{tenantID}/blog/{blogID}/credential/{credentialID}",
		controller.DeleteBlog).Methods(http.MethodDelete)

	log.NewLogger().Info("Blog Route Registered")
}

// GetBlogs returns all blogs.
func (controller *BlogController) GetBlogs(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogs call**************************************")

	// Create bucket.
	blogs := []blg.DTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all blogs service method.
	err = controller.BlogService.GetBlogs(&blogs, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, blogs)
}

// GetLatestBlogSnippets returns all blog snippets.
func (controller *BlogController) GetLatestBlogSnippets(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetLatestBlogSnippets call**************************************")

	// Create bucket.
	blogs := []blg.SnippetDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all blogs service method.
	err = controller.BlogService.GetLatestBlogSnippets(&blogs, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, blogs)
}

// GetTrendingBlogSnippets returns all trending blog snippets.
func (controller *BlogController) GetTrendingBlogSnippets(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetTrendingBlogSnippets call**************************************")

	// Create bucket.
	blogs := []blg.SnippetDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parse for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all blogs service method.
	err = controller.BlogService.GetTrendingBlogSnippets(&blogs, tenantID, r.Form, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, blogs)
}

// GetBlog return specific blog by id.
func (controller *BlogController) GetBlog(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlog call**************************************")

	// Create bucket.
	blog := blg.DTO{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blog.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting blog id from param and parsing it to uuid.
	blog.ID, err = util.ParseUUID(mux.Vars(r)["blogID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog ID", http.StatusBadRequest))
		return
	}

	// Call get blog by id service method.
	err = controller.BlogService.GetBlog(&blog)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, blog)
}

// AddBlog adds new blog.
func (controller *BlogController) AddBlog(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddBlog call**************************************")

	// Create bucket.
	blog := blg.Blog{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	if err := blog.Validate(); err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	blog.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blog.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blog.AuthorID, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.BlogService.AddBlog(&blog)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog added successfully")
}

// UpdateBlog updates the specified blog by id.
func (controller *BlogController) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateBlog call**************************************")

	// Create bucket.
	blog := blg.Blog{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blog.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting blog id from param and parsing it to uuid.
	blog.ID, err = util.ParseUUID(mux.Vars(r)["blogID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blog.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blog.AuthorID, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse requested data", http.StatusBadRequest))
		return
	}

	// Validate country fields
	err = blog.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.BlogService.UpdateBlog(&blog)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, "Blog updated successfully")
}

// DeleteBlog deletes the specified blog by id.
func (controller *BlogController) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteBlog call**************************************")

	// Create bcuket.
	blog := blg.Blog{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blog.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant ID", http.StatusBadRequest))
		return
	}

	// Getting blog id from param and parsing it to uuid.
	blog.ID, err = util.ParseUUID(mux.Vars(r)["blogID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog ID", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blog.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.BlogService.DeleteBlog(&blog)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog deleted successfully")
}

// UpdateBlogFlags will update the flags of specified blog record in the table.
func (controller *BlogController) UpdateBlogFlags(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================UpdateBlogFlags called==============================")

	// Params.
	param := mux.Vars(r)

	// Create bucket for flags update model.
	blog := blg.FlagsUpdateModel{}

	err := web.UnmarshalJSON(r, &blog)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	err = controller.BlogService.UpdateBlogFlags(&blog, tenantID, credentialID)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, "Blog updated successfully")
}
