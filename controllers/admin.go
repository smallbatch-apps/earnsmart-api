package controllers

import (
	"net/http"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type AdminController struct {
	service *services.AdminService
}

func NewAdminController(service *services.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (c *AdminController) SeedData(w http.ResponseWriter, r *http.Request) {
	admin_password := r.URL.Query().Get("admin_password")
	if admin_password != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Invalid admin password", http.StatusUnauthorized)
		return
	}
	serviceDb := c.service.GetDB()
	serviceDb.AutoMigrate(&models.Account{}, &models.Fund{}, &models.Price{}, &models.Setting{}, &models.User{})
	c.service.SeedData()
}
