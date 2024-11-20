package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterAccountRoutes(mux *http.ServeMux, controller *controllers.AccountController) {
	mux.HandleFunc("GET /accounts", controller.ListWalletAccounts)
	mux.HandleFunc("GET /fund-accounts", controller.ListFundingAccounts)
}
