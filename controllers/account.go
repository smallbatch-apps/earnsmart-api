package controllers

import (
	"log"
	"math/big"
	"net/http"
	"slices"

	"github.com/smallbatch-apps/earnsmart-api/errs"
	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/schema"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"

	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type AccountController struct {
	services *services.Services
}

func NewAccountController(services *services.Services) *AccountController {
	return &AccountController{services}
}

func (c *AccountController) ListAccounts(w http.ResponseWriter, r *http.Request) {

	s := c.services

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		errs.UnauthorisedError(w, "Unable to get user context")
		return
	}

	accounts, err := s.Account.GetAccounts(models.Account{OwnableModel: models.OwnableModel{UserID: userID}})
	if err != nil {
		errs.InternalError(w, "Unable to get accounts", "Unable to get accounts")
		return
	}

	accountIds, err := c.services.Account.ExtractIDs(accounts)
	if err != nil {
		errs.InternalError(w, "Unable to extract account ids", "Internal server error")
		return
	}

	balances, err := s.Account.LookupAccountBalances(accountIds)
	if err != nil {
		errs.InternalError(w, "Unable to get account balances", "Internal server error")
		return
	}

	var accountsWithBalances []schema.AccountWithBalance

	for _, account := range accounts {
		idx := slices.IndexFunc(balances, func(b services.AccountBalanceWithID) bool {
			return b.AccountID == tbt.ToUint128(uint64(account.ID))
		})

		foundBalance := balances[idx]
		credits := foundBalance.CreditsPosted.BigInt()
		debits := foundBalance.DebitsPosted.BigInt()
		balance := new(big.Int).Sub(&credits, &debits)
		balanceBigInt := tbt.BigIntToUint128(*balance)
		balanceUsd, _ := s.Price.GetAmountPrice(account.Currency, balanceBigInt)

		accountsWithBalances = append(accountsWithBalances, schema.AccountWithBalance{
			Account:    account,
			Balance:    balanceBigInt,
			BalanceUSD: balanceUsd,
		})
	}
	log.Printf("AccountsWithBalances: %+v", accountsWithBalances)

	utils.RespondOk(w, "accounts", accountsWithBalances)
}
