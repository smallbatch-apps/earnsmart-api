package controllers

import (
	"fmt"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type FundController struct {
	service *services.FundService
}

func NewFundController(service *services.FundService) *FundController {
	return &FundController{service: service}
}

func (c *FundController) ListFunds(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "listing all funds\n")
}

func (c *FundController) GetFund(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get a transactions\n")
}
