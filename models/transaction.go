package models

import (
	"math/big"

	"github.com/shopspring/decimal"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TransactionStatus uint16

const (
	_ TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusConfirmed
	TransactionStatusCancelled
	TransactionStatusFailed
)

type Transaction struct {
	CustomModel
	OwnableModel
	Amount     decimal.Decimal   `json:"amount"`
	Currency   string            `json:"currency"`
	Address    string            `json:"address"`
	Status     TransactionStatus `json:"status" gorm:"default:1"`
	Type       TransactionType   `json:"type"`
	TransferID string            `json:"transfer_id"`
}

func (t *Transaction) AmountAsUint128() tbt.Uint128 {
	amountStr := t.Amount.String()
	bigInt := new(big.Int)
	bigInt.SetString(amountStr, 10)
	return tbt.BigIntToUint128(*bigInt)
}

func (t *Transaction) TransferIDAsUint128() tbt.Uint128 {
	bigInt := new(big.Int)
	bigInt.SetString(t.TransferID, 10)
	return tbt.BigIntToUint128(*bigInt)
}
