package controllers

import (
	"net/http"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type AdminController struct {
	service *services.FundService
}

func NewAdminController(service *services.FundService) *AdminController {
	return &AdminController{service: service}
}

func (c *AdminController) SeedData(w http.ResponseWriter, r *http.Request) {
	admin_password := r.URL.Query().Get("admin_password")
	if admin_password != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Invalid admin password", http.StatusUnauthorized)
		return
	}

	c.service.db.AutoMigrate(&models.Price{}, &models.Fund{}, &models.Setting{}, &models.User{})

}
