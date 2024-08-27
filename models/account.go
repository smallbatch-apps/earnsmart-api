package models

import (
	"encoding/json"
	"errors"
	"math/big"

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
	UserID      uint        `json:"user_id" gorm:"type:uuid"`
	Currency    string      `json:"currency"`
	AccountCode AccountCode `json:"account_code" gorm:"type:int"`
}

func (account Account) MarshalJSON() ([]byte, error) {
	type Alias Account

	return json.Marshal(&struct {
		ID          uint   `json:"id"`
		Currency    string `json:"currency"`
		AccountCode uint   `json:"account_code"`
		Alias
	}{
		ID:          account.ID,
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
	_ AccountCode = iota
	AccountCodeWallet
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

type AccountWrapper struct {
	tbt.Account
}

func (a AccountWrapper) MarshalJSON() ([]byte, error) {
	type Alias tbt.Account
	return json.Marshal(&struct {
		ID             uint64
		DebitsPending  uint64
		DebitsPosted   uint64
		CreditsPending uint64
		CreditsPosted  uint64
		UserData128    uint64
		*Alias
	}{
		ID:             BigIntToUint64(a.ID.BigInt()),
		DebitsPending:  BigIntToUint64(a.DebitsPending.BigInt()),
		DebitsPosted:   BigIntToUint64(a.DebitsPosted.BigInt()),
		CreditsPending: BigIntToUint64(a.CreditsPending.BigInt()),
		CreditsPosted:  BigIntToUint64(a.CreditsPosted.BigInt()),
		UserData128:    BigIntToUint64(a.UserData128.BigInt()),
		Alias:          (*Alias)(&a.Account),
	})
}

func BigIntToUint64(bi big.Int) uint64 {
	return bi.Uint64()
}

type AccountFilterWrapper struct {
	tbt.AccountFilter
}

func (a AccountFilterWrapper) MarshalJSON() ([]byte, error) {
	type Alias tbt.AccountFilter
	return json.Marshal(&struct {
		AccountID uint64
		Reserved  uint64 `json:"-"`
		*Alias
	}{
		AccountID: BigIntToUint64(a.AccountID.BigInt()),
		Alias:     (*Alias)(&a.AccountFilter),
	})
}

type AccountBalanceWrapper struct {
	tbt.AccountBalance
}

func (a AccountBalanceWrapper) MarshalJSON() ([]byte, error) {
	type Alias tbt.AccountBalance
	return json.Marshal(&struct {
		DebitsPending  uint64
		DebitsPosted   uint64
		CreditsPending uint64
		CreditsPosted  uint64
		Reserved       uint64 `json:"-"`
		*Alias
	}{
		DebitsPending:  BigIntToUint64(a.DebitsPending.BigInt()),
		DebitsPosted:   BigIntToUint64(a.DebitsPosted.BigInt()),
		CreditsPending: BigIntToUint64(a.CreditsPending.BigInt()),
		CreditsPosted:  BigIntToUint64(a.CreditsPosted.BigInt()),
		Alias:          (*Alias)(&a.AccountBalance),
	})
}
