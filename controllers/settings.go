package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type SettingController struct {
	services *services.Services
}

func NewSettingController(services *services.Services) *SettingController {
	return &SettingController{services}
}

func (c *SettingController) ListSettings(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	settings, err := c.services.Setting.GetAll(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "settings", settings)
}

func (c *SettingController) EditSetting(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payload schema.SettingPayload
	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setting, err := c.services.Setting.GetSetting(userID, payload.Setting)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if setting.Type != models.SettingTypeUser {
		http.Error(w, "Invalid setting type", http.StatusUnauthorized)
		return
	}

	err = c.services.Setting.SetSetting(userID, payload.Setting, payload.Value)

	if err == nil {
		setting.Value = payload.Value
	}

	utils.RespondOk(w, "setting", setting)
}
