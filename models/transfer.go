package models

import (
	"encoding/json"

	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type TransactionType uint16

const (
	_ TransactionType = iota
	TransactionTypeDeposit
	TransactionTypeWithdraw
	TransactionTypeSwap
	TransactionTypeSubscribe
	TransactionTypeInterest
	TransactionTypeRedeem
)

type TransferCode uint16

const (
	_ TransactionType = iota
	TransferCodeDeposit
	TransferCodeWithdraw
	TransferCodeSwapFrom
	TransferCodeSwapTo
	TransferCodeSubscribe
	TransferCodeInterest
	TransferCodeRedeem
)

type TransferWrapper struct {
	tbt.Transfer
}

func (t TransferWrapper) MarshalJSON() ([]byte, error) {
	type Alias tbt.Transfer
	return json.Marshal(&struct {
		ID              uint64
		CreditAccountID uint64
		DebitAccountID  uint64
		Amount          uint64
		PendingID       uint64
		UserData128     uint64
		*Alias
	}{
		ID:              BigIntToUint64(t.ID.BigInt()),
		CreditAccountID: BigIntToUint64(t.CreditAccountID.BigInt()),
		DebitAccountID:  BigIntToUint64(t.DebitAccountID.BigInt()),
		Amount:          BigIntToUint64(t.Amount.BigInt()),
		PendingID:       BigIntToUint64(t.PendingID.BigInt()),
		UserData128:     BigIntToUint64(t.UserData128.BigInt()),
		Alias:           (*Alias)(&t.Transfer),
	})
}
