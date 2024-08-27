package routes

import (
	"log"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterAccountRoutes(mux *http.ServeMux, controller *controllers.AccountController) {
	log.Println("REGISTERING ACCOUNT ROUTES")
	mux.HandleFunc("GET /accounts", controller.ListAccounts)
}
