package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/services"
)

type PriceController struct {
	service *services.PriceService
}

func NewPriceController(service *services.PriceService) *PriceController {
	return &PriceController{service: service}
}

func (c *PriceController) ListPrices(w http.ResponseWriter, r *http.Request) {
	prices, err := c.service.GetLatestPrices()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
