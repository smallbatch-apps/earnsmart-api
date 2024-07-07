package routes

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/controllers"
)

func RegisterTransactionRoutes(mux *http.ServeMux, controller *controllers.TransactionController) {
	mux.HandleFunc("GET /transaction", controller.ListTransactions)
	mux.HandleFunc("GET /transaction/{id}", controller.GetTransaction)
	mux.HandleFunc("POST /transaction", controller.AddTransaction)
	mux.HandleFunc("PATCH /transaction/{id}", controller.EditTransaction)
}
