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

// BlogReplyController provide method to update, delete, add, get for blog reply.
type BlogReplyController struct {
	BlogReplyService *service.BlogReplyService
}

// NewBlogReplyController creates new instance of BlogReplyController.
func NewBlogReplyController(blogReplyService *service.BlogReplyService) *BlogReplyController {
	return &BlogReplyController{
		BlogReplyService: blogReplyService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *BlogReplyController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all blog replies by blog id.
	router.HandleFunc("/tenant/{tenantID}/blog-reply/blog/{blogID}",
		controller.GetBlogReplies).Methods(http.MethodGet)

	// Get one blog reply.
	router.HandleFunc("/tenant/{tenantID}/blog-reply/{blogReplyID}",
		controller.GetBlogReply).Methods(http.MethodGet)

	// Add one blog reply.
	router.HandleFunc("/tenant/{tenantID}/blog-reply/credential/{credentialID}",
		controller.AddBlogReply).Methods(http.MethodPost)

	// Update blog reply.
	router.HandleFunc("/tenant/{tenantID}/blog-reply/{blogReplyID}/credential/{credentialID}",
		controller.UpdateBlogReply).Methods(http.MethodPut)

	// Delete blog reply.
	router.HandleFunc("/tenant/{tenantID}/blog-reply/{blogReplyID}/credential/{credentialID}",
		controller.DeleteBlogReply).Methods(http.MethodDelete)

	log.NewLogger().Info("Blog reply Route Registered")
}

// GetBlogReplies returns all the blog replies.
func (controller *BlogReplyController) GetBlogReplies(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBlogReplies called==============================")

	// Create bucket.
	blogReplies := []blog.BlogReplyDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog id from param and parsing it to uuid.
	blogID, err := util.ParseUUID(mux.Vars(r)["blogID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// Call get all service method.
	err = controller.BlogReplyService.GetBlogReplies(&blogReplies, r.Form, tenantID, blogID)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogReplies)
}

// GetBlogReply returns the specified blog reply by id.
func (controller *BlogReplyController) GetBlogReply(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogReply call**************************************")

	// Create bucket.
	blogReply := blog.BlogReplyDTO{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogReply.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog reply id from param and parsing it to uuid.
	blogReply.ID, err = util.ParseUUID(mux.Vars(r)["blogReplyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog reply id", http.StatusBadRequest))
		return
	}

	// Call get blog reply by id service method.
	err = controller.BlogReplyService.GetBlogReply(&blogReply)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogReply)
}

// AddBlogReply adds new blog reply.
func (controller *BlogReplyController) AddBlogReply(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddBlogReply call**************************************")

	// Create bucket.
	blogReply := blog.BlogReply{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = blogReply.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	blogReply.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReply.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse credential id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.BlogReplyService.AddBlogReply(&blogReply)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reply added successfully")
}

// UpdateBlogReply updates the specified blog reply by id.
func (controller *BlogReplyController) UpdateBlogReply(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateBlogReply call**************************************")

	// Create bucket.
	blogReply := blog.BlogReply{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogReply.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting langauge id from param and parsing it to uuid.
	blogReply.ID, err = util.ParseUUID(mux.Vars(r)["blogReplyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog reply id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReply.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &blogReply)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate blog reply fields.
	err = blogReply.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.BlogReplyService.UpdateBlogReply(&blogReply)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reply updated successfully")
}

// DeleteBlogReply deletes specific blog reply by id.
func (controller *BlogReplyController) DeleteBlogReply(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteBlogReply call**************************************")

	// Create bucket.
	blogReply := blog.BlogReply{}

	var err error

	param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	blogReply.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReply.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog reply id from param and parsing it to uuid.
	blogReply.ID, err = util.ParseUUID(param["blogReplyID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog reply id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.BlogReplyService.DeleteBlogReply(&blogReply)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reply deleted successfully")
}
