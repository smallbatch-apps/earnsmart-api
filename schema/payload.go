package schema

import (
	"encoding/json"
	"math/big"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type AccountWithBalance struct {
	models.Account
	Balance    tbt.Uint128 `json:"balance"`
	BalanceUSD float64     `json:"balance_usd"`
}

func (a AccountWithBalance) MarshalJSON() ([]byte, error) {
	// Marshal the embedded Account struct
	accountJSON, err := json.Marshal(a.Account)
	if err != nil {
		return nil, err
	}

	// Create a map to hold the combined fields
	var accountMap map[string]interface{}
	if err := json.Unmarshal(accountJSON, &accountMap); err != nil {
		return nil, err
	}

	// Add the additional fields
	accountMap["balance"] = bigIntToUint64(a.Balance.BigInt())
	accountMap["balance_usd"] = a.BalanceUSD

	return json.Marshal(accountMap)
}

// bigIntToUint64 converts a big.Int to uint64
func bigIntToUint64(bi big.Int) uint64 {
	return (&bi).Uint64()
}

type SettingPayload struct {
	Setting string `json:"setting"`
	Value   string `json:"value"`
}

type TransactionPayload struct {
	TransactionType models.TransactionType `json:"transaction_type"`
	Amount          uint64                 `json:"amount"`
	Currency        string                 `json:"currency"`
}

type FundTransactionPayload struct {
	TransactionPayload
	FundID uint `json:"amount"`
}

type NewUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserResponse struct {
	User   models.User `json:"user"`
	Status string      `json:"status"`
}

type TransactionsResponse struct {
	Transactions []services.AccountTransferWithID `json:"transactions"`
	Status       string                           `json:"status"`
}

type SettingsResponse struct {
	Settings []models.Setting `json:"settings"`
	Status   string           `json:"status"`
}

type SettingResponse struct {
	Setting models.Setting `json:"setting"`
	Status  string         `json:"status"`
}

type QuoteRequest struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Amount       uint   `json:"amount"`
}

type QuoteResponseData struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	AmountFrom   uint    `json:"amount_from"`
	AmountTo     uint    `json:"amount_to"`
	Rate         float32 `json:"rate"`
}

type QuoteResponse struct {
	Status string            `json:"status"`
	Quote  QuoteResponseData `json:"quote"`
}

type AccountResponse struct {
	Status   string               `json:"status"`
	Accounts []AccountWithBalance `json:"accounts"`
}

type FundsResponse struct {
	Status string        `json:"status"`
	Funds  []models.Fund `json:"funds"`
}

type FundResponse struct {
	Status string      `json:"status"`
	Fund   models.Fund `json:"fund"`
}

type AccountSerializer struct {
	tbt.Account
}

type ActivitiesResponse struct {
	Activities []models.Activity `json:"activities"`
	Status     string            `json:"status"`
}
