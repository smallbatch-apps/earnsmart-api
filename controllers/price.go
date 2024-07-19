package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type PriceController struct {
	service *services.PriceService
}

func NewPriceController(service *services.PriceService) *PriceController {
	return &PriceController{service: service}
}

func (c *PriceController) ListPrices(w http.ResponseWriter, r *http.Request) {
	var prices []models.Price
	var err error

	if prices, err = c.service.GetPrices(); err != nil {
		errs.InternalError(w, err.Error(), err.Error())
	}

	if err := json.NewEncoder(w).Encode(prices); err != nil {
		errs.InternalError(w, "Encoding error in JSON response", err.Error())
	}
}
