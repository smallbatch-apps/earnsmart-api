package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterSwapRoutes(mux *http.ServeMux, controller *controllers.SwapController) {
	mux.HandleFunc("GET /swap", controller.ListSwaps)
	mux.HandleFunc("POST /swap", controller.CreateSwap)
}
