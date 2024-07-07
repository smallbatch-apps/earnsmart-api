package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterAdminRoutes(mux *http.ServeMux, controller *controllers.AdminController) {
	mux.HandleFunc("GET /seed", controller.SeedData)
}
