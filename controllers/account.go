package controllers

import (
	"net/http"

	"github.com/smallbatch-apps/earnsmart-api/middleware"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
	"github.com/smallbatch-apps/earnsmart-api/utils"
)

type AccountController struct {
	services *services.Services
}

func NewAccountController(services *services.Services) *AccountController {
	return &AccountController{services}
}

func (c *AccountController) ListWalletAccounts(w http.ResponseWriter, r *http.Request) {
	accountService := c.services.Account

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.RespondError(w, err, http.StatusUnauthorized)
		return
	}

	accounts, err := accountService.GetAccounts(models.Account{OwnableModel: models.OwnableModel{UserID: userID}, AccountCode: models.AccountCodeWallet})
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	accountIds, err := accountService.ExtractIDs(accounts)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	balances, err := accountService.LookupAccounts(accountIds)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "accounts", balances)
}

func (c *AccountController) ListFundingAccounts(w http.ResponseWriter, r *http.Request) {
	accountService := c.services.Account

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.RespondError(w, err, http.StatusUnauthorized)
		return
	}

	accounts, err := accountService.GetAccounts(models.Account{OwnableModel: models.OwnableModel{UserID: userID}})
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	fundingAccounts := []models.Account{}
	for _, account := range accounts {
		if account.AccountCode != models.AccountCodeWallet {
			fundingAccounts = append(fundingAccounts, account)
		}
	}
	if len(fundingAccounts) == 0 {
		utils.RespondOk(w, "funding_accounts", fundingAccounts)
		return
	}

	accountIds, err := accountService.ExtractIDs(fundingAccounts)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	balances, err := accountService.LookupAccounts(accountIds)
	if err != nil {
		utils.RespondError(w, err, http.StatusInternalServerError)
		return
	}

	utils.RespondOk(w, "funding_accounts", balances)
}
