package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterPriceRoutes(mux *http.ServeMux, controller *controllers.PriceController) {
	mux.HandleFunc("GET /prices", controller.ListPrices)
}
