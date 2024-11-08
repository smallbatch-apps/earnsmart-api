package controllers

import (
	"net/http"
	"strconv"

	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type FundController struct {
	services *services.Services
}

func NewFundController(services *services.Services) *FundController {
	return &FundController{services}
}

func (c *FundController) ListFunds(w http.ResponseWriter, r *http.Request) {
	funds, err := c.services.Fund.ListFunds()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "funds", funds)
}

func (c *FundController) GetFund(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fund, err := c.services.Fund.GetFund(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "fund", fund)
}
