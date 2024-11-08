package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type AdminController struct {
	services *services.Services
}

func NewAdminController(services *services.Services) *AdminController {
	return &AdminController{services: services}
}

func (c *AdminController) SeedData(w http.ResponseWriter, r *http.Request) {
	adminPassword := r.URL.Query().Get("admin_password")
	adminService := services.NewAdminService(c.services.Db, c.services.TbClient)

	if adminPassword != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Invalid admin password", http.StatusUnauthorized)
		return
	}

	if adminService.SafeToMigrate() {
		log.Println("Database is in a safe state to migrate")
		adminService.SeedData()
		log.Println("Migration and data seeding complete")
	} else {
		log.Println("Database is not in a safe state to migrate")
		return
	}
}
