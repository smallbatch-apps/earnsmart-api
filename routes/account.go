package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterAccountRoutes(mux *http.ServeMux, controller *controllers.AccountController) {
	mux.HandleFunc("GET /accounts", controller.ListAccounts)
}
