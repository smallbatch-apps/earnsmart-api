package controllers

import (
	"fmt"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type SettingController struct {
	service *services.SettingService
}

func NewSettingController(service *services.SettingService) *SettingController {
	return &SettingController{service: service}
}

func (c *SettingController) ListSettings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "listing all settings\n")
}

func (c *SettingController) GetSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "get a setting: %v\n", id)
}

func (c *SettingController) AddSetting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "adding a setting\n")
}

func (c *SettingController) EditSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "editing a setting: %v\n", id)
}

func (c *SettingController) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "deleting a setting: %v\n", id)
}
