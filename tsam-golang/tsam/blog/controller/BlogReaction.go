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

// BlogReactionController provide method to update, delete, add, get for blog reaction.
type BlogReactionController struct {
	BlogReactionService *service.BlogReactionService
}

// NewBlogReactionController creates new instance of BlogReactionController.
func NewBlogReactionController(blogReactionService *service.BlogReactionService) *BlogReactionController {
	return &BlogReactionController{
		BlogReactionService: blogReactionService,
	}
}

// RegisterRoutes registers all endpoints to router excluding list.
func (controller *BlogReactionController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// Get all blog reactions.
	router.HandleFunc("/tenant/{tenantID}/blog-reaction",
		controller.GetBlogReactions).Methods(http.MethodGet)

	// Get one blog reaction.
	router.HandleFunc("/tenant/{tenantID}/blog-reaction/{blogReactionID}",
		controller.GetBlogReaction).Methods(http.MethodGet)

	// Add one blog reaction.
	router.HandleFunc("/tenant/{tenantID}/blog-reaction/credential/{credentialID}",
		controller.AddBlogReaction).Methods(http.MethodPost)

	// Update blog reaction.
	router.HandleFunc("/tenant/{tenantID}/blog-reaction/{blogReactionID}/credential/{credentialID}",
		controller.UpdateBlogReaction).Methods(http.MethodPut)

	// Delete blog reaction.
	router.HandleFunc("/tenant/{tenantID}/blog-reaction/{blogReactionID}/credential/{credentialID}",
		controller.DeleteBlogReaction).Methods(http.MethodDelete)

	log.NewLogger().Info("Blog reaction Route Registered")
}

// GetBlogReactions returns all the blog reactions.
func (controller *BlogReactionController) GetBlogReactions(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("==============================GetBlogReactions called==============================")

	// Create bucket.
	blogReactions := []blog.BlogReactionDTO{}

	// Getting tenant id from param and parsing it to uuid.
	tenantID, err := util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Parsing for query params.
	r.ParseForm()

	// Call get all service method.
	err = controller.BlogReactionService.GetBlogReactions(&blogReactions, r.Form, tenantID)
	if err != nil {
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogReactions)
}

// GetBlogReaction returns the specified blog reaction by id.
func (controller *BlogReactionController) GetBlogReaction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************GetBlogReaction call**************************************")

	// Create bucket.
	blogReaction := blog.BlogReactionDTO{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogReaction.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog reaction id from param and parsing it to uuid.
	blogReaction.ID, err = util.ParseUUID(mux.Vars(r)["blogReactionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse blog reaction id", http.StatusBadRequest))
		return
	}

	// Call get blog reaction by id service method.
	err = controller.BlogReactionService.GetBlogReaction(&blogReaction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter
	web.RespondJSON(w, http.StatusOK, blogReaction)
}

// AddBlogReaction adds new blog reaction.
func (controller *BlogReactionController) AddBlogReaction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************AddBlogReaction call**************************************")

	// Create bucket.
	blogReaction := blog.BlogReaction{}

	// Unmarshal json.
	err := web.UnmarshalJSON(r, &blogReaction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Invalid Request", http.StatusBadRequest))
		return
	}

	// Validate compulsary fields.
	err = blogReaction.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Getting tenant id from param and parsing it to uuid.
	blogReaction.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReaction.CreatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Call add service method.
	err = controller.BlogReactionService.AddBlogReaction(&blogReaction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reaction added successfully")
}

// UpdateBlogReaction updates the specified blog reaction by id.
func (controller *BlogReactionController) UpdateBlogReaction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************UpdateBlogReaction call**************************************")

	// Create bucket.
	blogReaction := blog.BlogReaction{}

	var err error

	// Getting tenant id from param and parsing it to uuid.
	blogReaction.TenantID, err = util.ParseUUID(mux.Vars(r)["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting langauge id from param and parsing it to uuid.
	blogReaction.ID, err = util.ParseUUID(mux.Vars(r)["blogReactionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog reaction id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReaction.UpdatedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Unmarshal json.
	err = web.UnmarshalJSON(r, &blogReaction)
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse data", http.StatusBadRequest))
		return
	}

	// Validate blog reaction fields.
	err = blogReaction.Validate()
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, err)
		return
	}

	// Call update service method.
	err = controller.BlogReactionService.UpdateBlogReaction(&blogReaction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reaction updated successfully")
}

// DeleteBlogReaction deletes specific blog reaction by id.
func (controller *BlogReactionController) DeleteBlogReaction(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("********************************DeleteBlogReaction call**************************************")

	// Create bucket.
	blogReaction := blog.BlogReaction{}

	var err error

	param := mux.Vars(r)

	// Getting tenant id from param and parsing it to uuid.
	blogReaction.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting credential id from param and parsing it to uuid.
	blogReaction.DeletedBy, err = util.ParseUUID(mux.Vars(r)["credentialID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("unable to parse tenant id", http.StatusBadRequest))
		return
	}

	// Getting blog reaction id from param and parsing it to uuid.
	blogReaction.ID, err = util.ParseUUID(param["blogReactionID"])
	if err != nil {
		log.NewLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Unable to parse blog reaction id", http.StatusBadRequest))
		return
	}

	// Call delete service method.
	err = controller.BlogReactionService.DeleteBlogReaction(&blogReaction)
	if err != nil {
		web.RespondError(w, err)
		return
	}

	// Writing Response with OK Status to ResponseWriter.
	web.RespondJSON(w, http.StatusOK, "Blog reaction deleted successfully")
}
