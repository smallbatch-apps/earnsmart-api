package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
	"gorm.io/gorm"
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

	settings, err := c.services.Setting.ListSettings(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "settings", settings)
}

func (c *SettingController) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	var payload schema.SettingPayload
	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	setting, err := c.services.Setting.GetSetting(userID, payload.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	if err == nil && setting.Type != models.SettingTypeUser {
		http.Error(w, "user cannot modify admin setting", http.StatusForbidden)
		return
	}

	setting, err = c.services.Setting.UpdateSetting(userID, payload.Name, payload.Value)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "setting", setting)
}
