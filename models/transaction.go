package models

import (
	"encoding/json"
	"math/big"
	"time"

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

func (activity Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        uint64 `json:"id"`
		Type      uint16 `json:"type"`
		Currency  string `json:"currency"`
		Address   string `json:"address"`
		Status    uint16 `json:"status"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		ID:        activity.ID,
		Type:      uint16(activity.Type),
		Currency:  activity.Currency,
		Address:   activity.Address,
		Status:    uint16(activity.Status),
		CreatedAt: activity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: activity.UpdatedAt.Format(time.RFC3339),
	})
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
