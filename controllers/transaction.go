package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type TransactionController struct {
	services *services.Services
}

func NewTransactionController(services *services.Services) *TransactionController {
	return &TransactionController{services}
}

func (c *TransactionController) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	accountService := c.services.Account
	transferService := c.services.Transfer

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accounts, err := accountService.GetAccounts(models.Account{UserID: userID})
	if err != nil {
		errs.InternalError(w, "Unable to get accounts", "Internal server error")
		return
	}

	accountIds, err := accountService.ExtractIDs(accounts)
	if err != nil {
		errs.InternalError(w, "Unable to extract ids", "Internal server error")
		return
	}

	transfers, err := transferService.GetAllTransfers(accountIds)
	if err != nil {
		errs.InternalError(w, "Failed to successfully get transactions", "Internal server error")
		return
	}

	utils.RespondOk(w, "transactions", transfers)
}

func (c *TransactionController) AddFundTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.RespondError(w, err, http.StatusUnauthorized)
		return
	}

	var payload schema.FundTransactionPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.RespondError(w, err, http.StatusBadRequest)
		return
	}

	tType := payload.TransactionType

	var transaction models.Transaction

	if tType == models.TransactionTypeSubscribe {
		transaction, err = c.services.Transaction.CreateSubscription(userID, payload.Amount, payload.Currency, payload.AccountCode)
		if err != nil {
			utils.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	} else if tType == models.TransactionTypeRedeem {
		transaction, err = c.services.Transaction.CreateRedemption(userID, payload.Amount, payload.Currency, payload.AccountCode)
		if err != nil {
			utils.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	utils.RespondOk(w, "transaction", transaction)
}

func (c *TransactionController) AddTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		errs.UserTokenNotValidError(w)
		return
	}

	var payload schema.TransactionPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		errs.InvalidPayloadError(w)
		return
	}

	tType := payload.TransactionType
	var transaction models.Transaction

	if tType == models.TransactionTypeDeposit {
		transaction, err = c.services.Transaction.CreatePendingDeposit(userID, payload.Amount, payload.Address, payload.Currency)
		if err != nil {
			utils.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	} else if tType == models.TransactionTypeWithdraw {
		transaction, err = c.services.Transaction.CreatePendingWithdrawal(userID, payload.Amount, payload.Address, payload.Currency)
		if err != nil {
			utils.RespondError(w, err, http.StatusInternalServerError)
			return
		}
	}

	watcher := &services.TransactionWatcher{
		TransactionService: c.services.Transaction,
	}
	watcher.WatchTransaction(transaction)

	w.WriteHeader(http.StatusCreated)
	utils.RespondOk(w, "transaction", transaction)
}
