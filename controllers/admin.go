package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type AdminController struct {
	service *services.AdminService
}

func NewAdminController(service *services.AdminService) *AdminController {
	return &AdminController{service: service}
}

func (c *AdminController) SeedData(w http.ResponseWriter, r *http.Request) {
	adminPassword := r.URL.Query().Get("admin_password")

	if adminPassword != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Invalid admin password", http.StatusUnauthorized)
		return
	}

	if c.service.SafeToMigrate() {
		log.Println("Database is in a safe state to migrate")
		c.service.SeedData()
		log.Println("Migration and data seeding complete")
	} else {
		log.Println("Database is not in a safe state to migrate")
		return
	}
}
