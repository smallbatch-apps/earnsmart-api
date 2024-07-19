package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type FundController struct {
	service *services.FundService
}

func NewFundController(service *services.FundService) *FundController {
	return &FundController{service: service}
}

func (c *FundController) ListFunds(w http.ResponseWriter, r *http.Request) {
	funds, err := c.service.ListFunds()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := schema.FundsResponse{Status: "success", Funds: funds}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *FundController) GetFund(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fund, err := c.service.GetFund(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := schema.FundResponse{Status: "success", Fund: fund}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
