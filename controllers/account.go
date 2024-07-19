package controllers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"slices"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"

	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type AccountController struct {
	service *services.AccountService
}

func NewAccountController(service *services.AccountService) *AccountController {
	return &AccountController{service: service}
}

func (c *AccountController) ListAccounts(w http.ResponseWriter, r *http.Request) {

	priceService := services.NewPriceService(c.service.GetDB(), c.service.GetTBClient())
	userID, err := middleware.GetUserIDFromContext(r.Context())

	if err != nil {
		errs.UnauthorisedError(w, "Unable to get user context")
		return
	}

	accounts, err := c.service.GetAccounts(models.Account{UserID: userID})
	if err != nil {
		errs.InternalError(w, "Unable to get accounts", "Unable to get accounts")
		return
	}

	accountIds, err := c.service.ExtractIDs(accounts)
	if err != nil {
		errs.InternalError(w, "Unable to extract account ids", "Internal server error")
		return
	}

	balances, err := c.service.LookupAccountBalances(accountIds)
	if err != nil {
		errs.InternalError(w, "Unable to get account balances", "Internal server error")
		return
	}

	var accountsWithBalances []schema.AccountWithBalance

	for _, account := range accounts {
		idx := slices.IndexFunc(balances, func(b services.AccountBalanceWithID) bool { return b.AccountID == account.AccountID })
		foundBalance := balances[idx]
		credits := foundBalance.CreditsPosted.BigInt()
		debits := foundBalance.DebitsPosted.BigInt()
		balance := new(big.Int).Sub(&credits, &debits)
		balanceBigInt := tbt.BigIntToUint128(*balance)
		balanceUsd, _ := priceService.GetAmountPrice(account.Currency, balanceBigInt)

		accountsWithBalances = append(accountsWithBalances, schema.AccountWithBalance{
			Account:    account,
			Balance:    balanceBigInt,
			BalanceUSD: balanceUsd,
		})
	}

	response := schema.AccountResponse{
		Status:   "success",
		Accounts: accountsWithBalances,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		errs.InternalError(w, err.Error(), err.Error())
		return
	}
}
