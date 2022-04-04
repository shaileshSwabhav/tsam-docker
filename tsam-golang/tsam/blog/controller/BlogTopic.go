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

// BlogTopicController provide method to update, delete, add, get for blog topic.
type BlogTopicController struct {
	BlogTopicService *service.BlogTopicService
}

// NewBlogTopicController creates new instance of BlogTopicController.
func NewBlogTopicController(blogTopicService *service.BlogTopicService) *BlogTopicController {
	return &BlogTopicController{
		BlogTopicService: blogTopicService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *BlogTopicController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all blog topics by limit and offset.
	router.HandleFunc("/tenant/{tenantID}/blog-topic/limit/{limit}/offset/{offset}",
		controller.GetBlogTopics).Methods(http.MethodGet)

	// Get blog topic list.
	router.HandleFunc("/tenant/{tenantID}/blog-topic",
		controller.GetBlogTopicList).Methods(http.MethodGet)

	// Get one blog topic.
	router.HandleFunc("/tenant/{tenantID}/blog-topic/{blogTopicID}",
		controller.GetBlogTopic).Methods(http.MethodGet)

	// Add one blog topic.
	router.HandleFunc("/tenant/{tenantID}/blog-topic/credential/{credentialID}",
		controller.AddBlogTopic).Methods(http.MethodPost)

	// Update blog topic.
	router.HandleFunc("/tenant/{tenantID}/blog-topic/{blogTopicID}/credential/{credentialID}",
		controller.UpdateBlogTopic).Methods(http.MethodPut)

	// Delete blog topic.
	router.HandleFunc("/tenant/{tenantID}/blog-topic/{blogTopicID}/credential/{credentialID}",
		controller.DeleteBlogTopic).Methods(http.MethodDelete)

	log.NewLogger().Info("Blog topic Route Registered")
}

// GetBlogTopics returns all the blog topics.
func (controller *BlogTopicController) GetBlogTopics(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBlogTopics called==============================")

	// Create bucket.
	blogTopics := []blog.BlogTopic{}

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
	err = controller.BlogTopicService.GetBlogTopics(&blogTopics, r.Form, tenantID, limit, offset, &totalCount)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, blogTopics)
}

// GetBlogTopicList returns all the blog topics (without pagination).
func (controller *BlogTopicController) GetBlogTopicList(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogTopicList call**************************************")

	// Create bucket.
	blogTopics := []blog.BlogTopic{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call get all service method.
	err = controller.BlogTopicService.GetBlogTopicList(&blogTopics, tenantID)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, blogTopics)
}

// GetBlogTopic returns the specified blog topic by id.
func (controller *BlogTopicController) GetBlogTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogTopic call**************************************")

	// Create bucket.
	blogTopic := blog.BlogTopic{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogTopic.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog topic id from param and parsing it to uuid.
	blogTopic.ID, err = util.ParseUUID(mux.Vars(r)["blogTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog topic id", http.StatusBadRequest))
		return
	}

	// Call get blog topic by id service method.
	err = controller.BlogTopicService.GetBlogTopic(&blogTopic)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogTopic)
}

// AddBlogTopic adds new blog topic.
func (controller *BlogTopicController) AddBlogTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddBlogTopic call**************************************")

	// Create bucket.
	blogTopic := blog.BlogTopic{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = blogTopic.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	blogTopic.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogTopic.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.BlogTopicService.AddBlogTopic(&blogTopic)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog topic added successfully")
}

// UpdateBlogTopic updates the specified blog topic by id.
func (controller *BlogTopicController) UpdateBlogTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateBlogTopic call**************************************")

	// Create bucket.
	blogTopic := blog.BlogTopic{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogTopic.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting langauge id from param and parsing it to uuid.
	blogTopic.ID, err = util.ParseUUID(mux.Vars(r)["blogTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog topic id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogTopic.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &blogTopic)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate blog topic fields.
	err = blogTopic.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.BlogTopicService.UpdateBlogTopic(&blogTopic)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog topic updated successfully")
}

// DeleteBlogTopic deletes specific blog topic by id.
func (controller *BlogTopicController) DeleteBlogTopic(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteBlogTopic call**************************************")

	// Create bucket.
	blogTopic := blog.BlogTopic{}

	var err error

	param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	blogTopic.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogTopic.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog topic id from param and parsing it to uuid.
	blogTopic.ID, err = util.ParseUUID(param["blogTopicID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog topic id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.BlogTopicService.DeleteBlogTopic(&blogTopic)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog topic deleted successfully")
}
