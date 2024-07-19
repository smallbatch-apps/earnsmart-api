package models

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type CurrencyOption int

const (
	ETH CurrencyOption = iota
	BTC
	USDT
	USDC
	DAI
	BNB
	MATIC
	AVAX
	SOL
	BAT
	LINK
	UNI
	XRP
	ADA
	HBAR
	DOT
	TRX
)

type Account struct {
	CustomModel
	AccountID   uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID      uint      `gorm:"type:uuid"`
	Currency    string
	AccountCode AccountCode `gorm:"type:int"`
}

func (account Account) MarshalJSON() ([]byte, error) {
	type Alias Account

	return json.Marshal(&struct {
		ID          uint      `json:"id"`
		AccountID   uuid.UUID `json:"account_id"`
		Currency    string    `json:"currency"`
		AccountCode uint      `json:"account_code"`
		Alias
	}{
		ID:          account.ID,
		AccountID:   account.AccountID,
		Currency:    account.Currency,
		AccountCode: uint(account.AccountCode),
		Alias:       (Alias)(account),
	})
}

var currencies = [...]string{
	"ETH", "BTC", "USDT", "USDC", "DAI", "BNB", "MATIC", "AVAX",
	"SOL", "BAT", "LINK", "UNI", "XRP", "ADA", "HBAR", "DOT", "TRX",
}

func (c CurrencyOption) String() string {
	if c < ETH || c > TRX {
		return "Unknown"
	}
	return currencies[c]
}

func CurrencyFromInt(value int) (CurrencyOption, error) {
	if value < int(ETH) || value > int(TRX) {
		return 0, errors.New("invalid currency code")
	}
	return CurrencyOption(value), nil
}

type AccountCode uint64

const (
	AccountCodeWallet AccountCode = iota
	AccountCodeFundImmediate
	AccountCodeFundMonth
	AccountCodeFund3Months
	AccountCodeFund6Months
	AccountCodeFundOneYear
)

func GenerateUUID() tbt.Uint128 {
	var u128 tbt.Uint128
	id, _ := uuid.NewV7()
	copy(u128[:], id[:]) // Copy the UUID bytes into the Uint128 array
	return u128
}

func ConvertUUIDToUint128(id uuid.UUID) tbt.Uint128 {
	var u128 tbt.Uint128
	copy(u128[:], id[:]) // Copy the UUID bytes into the Uint128 array
	return u128
}

func ConvertUint128ToUUID(u128 tbt.Uint128) uuid.UUID {
	var id uuid.UUID
	copy(id[:], u128[:]) // Copy the Uint128 bytes into the UUID array
	return id
}
