package controllers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

type QuoteController struct {
	service *services.PriceService
}

func NewQuoteController(service *services.PriceService) *PriceController {
	return &PriceController{service: service}
}

func (c *QuoteController) GetQuote(w http.ResponseWriter, r *http.Request) {
	payload := schema.QuoteRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errs.InternalError(w, err.Error(), err.Error())
		return
	}

	priceFrom, err := c.service.GetPriceForCurrency(payload.FromCurrency)
	if err != nil {
		errs.InternalError(w, "Unable to get price for currency:"+err.Error(), "Internal server error")
	}

	priceTo, err := c.service.GetPriceForCurrency(payload.ToCurrency)
	if err != nil {
		errs.InternalError(w, "Unable to get price for currency:"+err.Error(), "Internal server error")
	}

	decimals := models.AllCurrencies[payload.ToCurrency].Decimals
	totalUsd := priceFrom.Rate * float32(payload.Amount)
	totalTo := uint((totalUsd / priceTo.Rate) * float32(math.Pow10(int(decimals))))

	quote := schema.QuoteResponseData{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
		AmountFrom:   payload.Amount,
		AmountTo:     totalTo,
		Rate:         priceTo.Rate,
	}

	if err := json.NewEncoder(w).Encode(quote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
