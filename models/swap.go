package models

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

type SwapStatus uint16

const (
	_ SwapStatus = iota
	SwapStatusPending
	SwapStatusConfirmed
	SwapStatusCancelled
	SwapStatusFailed
)

type Swap struct {
	CustomModel
	OwnableModel

	FromAmount     decimal.Decimal `json:"from_amount" gorm:"type:decimal(32,16);not null"`
	FromCurrency   string          `json:"from_currency" gorm:"not null"`
	FromTransferID string          `json:"from_transfer_id" gorm:"column:from_transfer_id;type:char(32)"`

	ToAmount     decimal.Decimal `json:"to_amount" gorm:"type:decimal(32,16);not null"`
	ToCurrency   string          `json:"to_currency" gorm:"not null"`
	ToTransferID string          `json:"to_transfer_id" gorm:"column:to_transfer_id;type:char(32)"`

	Rate   float64    `json:"rate" gorm:"type:decimal(32,16);not null"`
	Status SwapStatus `json:"status" gorm:"not null;default:1"`
}

func (swap Swap) MarshalJSON() ([]byte, error) {
	type Alias Swap

	return json.Marshal(&struct {
		Status uint16 `json:"status"`
		Alias
	}{
		Status: uint16(swap.Status),
		Alias:  (Alias)(swap),
	})
}
