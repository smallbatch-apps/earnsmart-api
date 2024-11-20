package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type SwapController struct {
	services *services.Services
}

func NewSwapController(services *services.Services) *SwapController {
	return &SwapController{services}
}

func (c *SwapController) ListSwaps(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	settings, err := c.services.Swap.ListSwaps(userID)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "swaps", settings)
}

func (c *SwapController) CreateSwap(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	var payload schema.SwapPayload
	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	swap, err := c.services.Swap.CreateSwap(userID, payload.FromAmount, payload.FromCurrency, payload.ToAmount, payload.ToCurrency, payload.Rate)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "swap", swap)
}
