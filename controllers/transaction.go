package controllers

import (
	"fmt"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type TransactionController struct {
	service *services.TransactionService
}

func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (c *TransactionController) ListTransactions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "listing all transactions\n")
}

func (c *TransactionController) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "get a transactions: id=%v\n", id)
}

func (c *TransactionController) AddTransaction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "adding a transaction\n")
}

func (c *TransactionController) EditTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "editing a transaction: %v\n", id)
}

func (c *TransactionController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "deleting a transaction: %v\n", id)
}
