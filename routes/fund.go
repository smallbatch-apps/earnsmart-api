package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterFundRoutes(mux *http.ServeMux, controller *controllers.FundController) {
	mux.HandleFunc("GET /funds", controller.ListFunds)
	mux.HandleFunc("GET /funds/{id}", controller.GetFund)
}
