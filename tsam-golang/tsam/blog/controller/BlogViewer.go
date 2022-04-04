package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/blog/service"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/blog"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// BlogViewController provide method to update, delete, add, get for blog view.
type BlogViewController struct {
	BlogViewService *service.BlogViewService
}

// NewBlogViewController creates new instance of BlogViewController.
func NewBlogViewController(BlogViewService *service.BlogViewService) *BlogViewController {
	return &BlogViewController{
		BlogViewService: BlogViewService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *BlogViewController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all blog views by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/blog-viewer/limit/{limit}/offset/{offset}",
		controller.GetBlogViews).Methods(http.MethodGet)

	// Get one blog view.
	router.HandleFunc("/tenant/{tenantID}/blog-viewer/{blogViewID}",
		controller.GetBlogView).Methods(http.MethodGet)

	// Add one blog view.
	router.HandleFunc("/tenant/{tenantID}/blog-viewer/credential/{credentialID}",
		controller.AddBlogView).Methods(http.MethodPost)

	// Update blog view.
	router.HandleFunc("/tenant/{tenantID}/blog-viewer/{blogViewID}/credential/{credentialID}",
		controller.UpdateBlogView).Methods(http.MethodPut)

	// Delete blog view.
	router.HandleFunc("/tenant/{tenantID}/blog-viewer/{blogViewID}/credential/{credentialID}",
		controller.DeleteBlogView).Methods(http.MethodDelete)

	log.NewLogger().Info("Blog view Route Registered")
}

// GetBlogViews returns all the blog views.
func (controller *BlogViewController) GetBlogViews(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBlogGetBlogViewsViewers called==============================")

	// Create bucket.
	blogViews := []blog.BlogView{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// For pagination.
	var totalCount int
	limit, offset := web.GetLimitAndOffset(r)

	// Call get all service method.
	err = controller.BlogViewService.GetBlogViews(&blogViews, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, blogViews)
}

// GetBlogView returns the specified blog view by id.
func (controller *BlogViewController) GetBlogView(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogView call**************************************")

	// Create bucket.
	blogView := blog.BlogView{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogView.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog view id from param and parsing it to uuid.
	blogView.ID, err = util.ParseUUID(mux.Vars(r)["blogViewID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog view id", http.StatusBadRequest))
		return
	}

	// Call get blog view by id service method.
	err = controller.BlogViewService.GetBlogView(&blogView)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogView)
}

// AddBlogView adds new blog view.
func (controller *BlogViewController) AddBlogView(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddBlogView call**************************************")

	// Create bucket.
	blogView := blog.BlogView{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &blogView)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = blogView.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	blogView.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogView.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.BlogViewService.AddBlogView(&blogView)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog view added successfully")
}

// UpdateBlogView updates the specified blog view by id.
func (controller *BlogViewController) UpdateBlogView(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateBlogView call**************************************")

	// Create bucket.
	blogView := blog.BlogView{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogView.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting langauge id from param and parsing it to uuid.
	blogView.ID, err = util.ParseUUID(mux.Vars(r)["blogViewID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog view id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogView.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &blogView)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate blog view fields.
	err = blogView.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.BlogViewService.UpdateBlogView(&blogView)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog view updated successfully")
}

// DeleteBlogView deletes specific blog view by id.
func (controller *BlogViewController) DeleteBlogView(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteBlogView call**************************************")

	// Create bucket.
	blogView := blog.BlogView{}

	var err error

	param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	blogView.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogView.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog view id from param and parsing it to uuid.
	blogView.ID, err = util.ParseUUID(param["blogViewID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog view id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.BlogViewService.DeleteBlogView(&blogView)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog view deleted successfully")
}
