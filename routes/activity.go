package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterActivityRoutes(mux *http.ServeMux, controller *controllers.ActivityController) {
	mux.HandleFunc("GET /activity", controller.ListActivities)
}
