package controllers

import (
	"encoding/json"
	"math"
	"math/big"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type QuoteController struct {
	services *services.Services
}

func NewQuoteController(services *services.Services) *QuoteController {
	return &QuoteController{services}
}

func (c *QuoteController) GetQuote(w http.ResponseWriter, r *http.Request) {
	payload := schema.QuoteRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errs.InternalError(w, err.Error(), err.Error())
		return
	}

	priceFrom, err := c.services.Price.GetPriceForCurrency(payload.FromCurrency)
	if err != nil {
		errs.InternalError(w, "Unable to get price for currency:"+err.Error(), "Internal server error")
	}

	priceTo, err := c.services.Price.GetPriceForCurrency(payload.ToCurrency)
	if err != nil {
		errs.InternalError(w, "Unable to get price for currency:"+err.Error(), "Internal server error")
	}

	decimals := models.AllCurrencies[payload.ToCurrency].Decimals

	// Convert amount to big.Float
	amountBig := new(big.Float).SetUint64(uint64(payload.Amount))

	// Convert rates to big.Float
	rateFromBig := new(big.Float).SetFloat64(priceFrom.Rate)
	rateToBig := new(big.Float).SetFloat64(priceTo.Rate)

	// Calculate: amount * fromRate / toRate * 10^decimals
	result := new(big.Float).Mul(amountBig, rateFromBig)
	result = new(big.Float).Quo(result, rateToBig)
	result = result.Mul(result, new(big.Float).SetFloat64(math.Pow10(int(decimals))))

	// Convert result to uint
	resultUint, _ := result.Uint64()
	totalTo := uint(resultUint)

	quote := schema.QuoteResponseData{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
		AmountFrom:   payload.Amount,
		AmountTo:     totalTo,
		Rate:         priceTo.Rate,
	}

	utils.RespondOk(w, "quote", quote)
}
