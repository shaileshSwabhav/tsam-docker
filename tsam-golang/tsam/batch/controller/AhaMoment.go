package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/batch/service"
	"github.com/techlabs/swabhav/tsam/log"
	"github.com/techlabs/swabhav/tsam/models/batch"
	"github.com/techlabs/swabhav/tsam/util"
	"github.com/techlabs/swabhav/tsam/web"
)

// AhaMomentController provides method to update, delete, add, get method for aha_moment.
type AhaMomentController struct {
	AhaMomentService *service.AhaMomentService
}

// NewAhaMomentController creates new instance of AhaMomentController.
func NewAhaMomentController(ahaMomentService *service.AhaMomentService) *AhaMomentController {
	return &AhaMomentController{
		AhaMomentService: ahaMomentService,
	}
}

// RegisterRoutes registers all endpoints to router.
func (con *AhaMomentController) RegisterRoutes(router *mux.Router, exclude *[]*mux.Route) {

	// add
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/aha-moment/credential/{credentialID}",
		con.AddAhaMoments).Methods(http.MethodPost)

	// update
	// router.HandleFunc("/tenant/{tenantID}/aha-moment/{ahaMomentID}/credential/{credentialID}", con.UpdateAhaMoment).Methods(http.MethodPut)

	// delete
	router.HandleFunc("/tenant/{tenantID}/batch/aha-moment/{ahaMomentID}/credential/{credentialID}",
		con.DeleteAhaMoment).Methods(http.MethodDelete)

	// get
	router.HandleFunc("/tenant/{tenantID}/batch/{batchID}/session/{sessionID}/aha-moment", con.GetAllAhaMoments).Methods(http.MethodGet)

	log.NewLogger().Info("Aha Moment Routes Registered")
}

// AddAhaMoments will add multiple ahaMoments and its response to the table.
func (con *AhaMomentController) AddAhaMoments(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===========================AddAhaMoments called===========================")
	param := mux.Vars(r)
	ahaMoments := []batch.AhaMoment{}

	err := web.UnmarshalJSON(r, &ahaMoments)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	credentialID, err := util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	sessionID, err := util.ParseUUID(param["sessionID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	for _, ahaMoment := range ahaMoments {
		err = ahaMoment.ValidateAhaMoment()
		if err != nil {
			log.NewLogger().Error(err)
			web.RespondError(w, err)
			return
		}
	}

	err = con.AhaMomentService.AddAhaMoments(&ahaMoments, tenantID, batchID, sessionID, credentialID)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Aha Moments successfully added")
}

// DeleteAhaMoment will delete aha moment and its responses.
func (con *AhaMomentController) DeleteAhaMoment(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===========================DeleteAhaMoment called===========================")
	param := mux.Vars(r)
	var err error
	ahaMoment := batch.AhaMoment{}

	ahaMoment.DeletedBy, err = util.ParseUUID(param["credentialID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	ahaMoment.TenantID, err = util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	ahaMoment.ID, err = util.ParseUUID(param["ahaMomentID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	// ahaMoment.TalentID, err = util.ParseUUID(param["talentID"])
	// if err != nil {
	// 	log.NewLogger().Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	// ahaMoment.FacultyID, err = util.ParseUUID(param["facultyID"])
	// if err != nil {
	// 	log.NewLogger().Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	// ahaMoment.BatchID, err = util.ParseUUID(param["batchID"])
	// if err != nil {
	// 	log.NewLogger().Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	// ahaMoment.BatchSessionID, err = util.ParseUUID(param["sessionID"])
	// if err != nil {
	// 	log.NewLogger().Error(err)
	// 	web.RespondError(w, err)
	// 	return
	// }

	err = con.AhaMomentService.DeleteAhaMoment(&ahaMoment)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Aha Moments successfully deleted")
}

// GetAllAhaMoments will return all the aha moments and its response for specified batch and session
func (con *AhaMomentController) GetAllAhaMoments(w http.ResponseWriter, r *http.Request) {
	log.NewLogger().Info("===========================GetAllAhaMoments called===========================")
	param := mux.Vars(r)
	var err error
	ahaMoments := []batch.AhaMomentDTO{}

	tenantID, err := util.ParseUUID(param["tenantID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	batchID, err := util.ParseUUID(param["batchID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	sessionID, err := util.ParseUUID(param["sessionID"])
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}

	err = con.AhaMomentService.GetAllAhaMoments(&ahaMoments, tenantID, batchID, sessionID)
	if err != nil {
		log.NewLogger().Error(err)
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, ahaMoments)
}

// // UpdateAhaMoment will update ahaMoment and its response.
// func (con *AhaMomentController) UpdateAhaMoment(w http.ResponseWriter, r *http.Request) {
// 	log.NewLogger().Info("===========================UpdateAhaMoment called===========================")
// 	param := mux.Vars(r)
// 	ahaMoment := batch.AhaMoment{}

// 	err := web.UnmarshalJSON(r, &ahaMoment)
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	ahaMoment.UpdatedBy, err = util.ParseUUID(param["credentialID"])
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	ahaMoment.TenantID, err = util.ParseUUID(param["tenantID"])
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	ahaMoment.ID, err = util.ParseUUID(param["ahaMomentID"])
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = ahaMoment.ValidateAhaMoment()
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}

// 	err = con.AhaMomentService.UpdateAhaMoment(&ahaMoment)
// 	if err != nil {
// 		log.NewLogger().Error(err)
// 		web.RespondError(w, err)
// 		return
// 	}
// 	web.RespondJSON(w, http.StatusOK, "Aha Moments successfully updated")
// }
