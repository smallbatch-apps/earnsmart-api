package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"

	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TransactionController struct {
	service *services.TransactionService
}

func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (c *TransactionController) ListTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())

	accountService := services.NewAccountService(c.service.GetDB(), c.service.GetTBClient())
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

	transactions, err := c.service.GetAllTransactions(accountIds)
	if err != nil {
		errs.InternalError(w, "Failed to successfully get transactions", "Internal server error")
		return
	}

	transactionsResponse := schema.TransactionsResponse{
		Status:       "success",
		Transactions: transactions,
	}

	if err := json.NewEncoder(w).Encode(transactionsResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *TransactionController) AddTransaction(w http.ResponseWriter, r *http.Request) {
	fundService := services.NewFundService(c.service.GetDB(), c.service.GetTBClient())
	accountService := services.NewAccountService(c.service.GetDB(), c.service.GetTBClient())
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

	accountSearch := models.Account{
		UserID:   userID,
		Currency: payload.Currency,
	}

	tType := payload.TransactionType
	isWalletType := tType == models.TransactionTypeDeposit || tType == models.TransactionTypeWithdraw
	isFundType := tType == models.TransactionTypeDeploy || tType == models.TransactionTypeRedeem

	if isWalletType {
		accountSearch.AccountCode = models.AccountCodeWallet

		if payload.TransactionType == models.TransactionTypeDeposit {
			account, err := accountService.GetOrCreateAccount(accountSearch)
			if err != nil {
				errs.InternalError(w, err.Error(), "Internal server error")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			c.service.CreateDeposit(account.AccountID, tbt.ToUint128(payload.Amount), payload.Currency)
		} else if payload.TransactionType == models.TransactionTypeWithdraw {
			accounts, err := accountService.GetAccounts(accountSearch)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(accounts) == 0 {
				http.Error(w, "Account not found", http.StatusBadRequest)
				return
			}
			c.service.CreateWithdrawal(accounts[0], tbt.ToUint128(payload.Amount), payload.Currency)
		}
	} else if isFundType {
		var payload schema.FundTransactionPayload
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			errs.InvalidPayloadError(w)
			return
		}

		fund, err := fundService.GetFund(payload.FundID)
		if err != nil {
			errs.InternalError(w, err.Error(), "Internal server error")
			return
		}

		walletSearchParams := models.Account{UserID: userID, Currency: payload.Currency, AccountCode: models.AccountCodeWallet}
		walletAccounts, err := accountService.GetAccounts(walletSearchParams)

		if len(walletAccounts) == 0 || err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		walletAccount := walletAccounts[0]
		accountSearch.AccountCode = fund.Code

		accounts, err := accountService.GetAccounts(accountSearch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var account models.Account
		if len(accounts) == 0 {
			if payload.TransactionType == models.TransactionTypeDeploy {
				account, err = accountService.GetOrCreateAccount(accountSearch)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				http.Error(w, "Account not found", http.StatusBadRequest)
				return
			}
		} else {
			account = accounts[0]
		}

		var creditAccount models.Account
		var debitAccount models.Account

		if tType == models.TransactionTypeDeploy {
			creditAccount = account
			debitAccount = walletAccount
		} else {
			creditAccount = walletAccount
			debitAccount = account
		}

		_, err = c.service.CreateTransfer(creditAccount, debitAccount, tbt.ToUint128(payload.Amount), payload.Currency)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Transaction created successfully"))
}
