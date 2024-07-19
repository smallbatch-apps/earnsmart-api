package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type SettingController struct {
	service *services.SettingService
}

func NewSettingController(service *services.SettingService) *SettingController {
	return &SettingController{service: service}
}

func (c *SettingController) ListSettings(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	settings, err := c.service.GetAll(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(settings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (c *SettingController) GetSetting(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	setting, err := c.service.Get(id)
// 	fmt.Fprintf(w, "get a setting: %v\n", id)
// }

// func (c *SettingController) AddSetting(w http.ResponseWriter, r *http.Request) {
//   userID, err := middleware.GetUserIDFromContext(r.Context())
//   //
// 	// fmt.Fprint(w, "adding a setting\n")
// }

func (c *SettingController) EditSetting(w http.ResponseWriter, r *http.Request) {
	// id := r.PathValue("id")
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

	setting, err := c.service.GetSetting(userID, payload.Setting)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if setting.Type != models.SettingTypeUser {
		http.Error(w, "Invalid setting type", http.StatusUnauthorized)
		return
	}

	err = c.service.SetSetting(userID, payload.Setting, payload.Value)

	if err == nil {
		setting.Value = payload.Value
	}

	response := schema.SettingResponse{Status: "ok", Setting: setting}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
