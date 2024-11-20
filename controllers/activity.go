package controllers

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type ActivityController struct {
	services *services.Services
}

func NewActivityController(services *services.Services) *ActivityController {
	return &ActivityController{services}
}

func (c *ActivityController) ListActivities(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activities, err := c.services.Activity.ListActivities(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "activities", activities)
}
