package controllers

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type PriceController struct {
	services *services.Services
}

func NewPriceController(services *services.Services) *PriceController {
	return &PriceController{services}
}

func (c *PriceController) ListPrices(w http.ResponseWriter, r *http.Request) {
	prices, err := c.services.Price.GetPrices()
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "prices", prices)
}
