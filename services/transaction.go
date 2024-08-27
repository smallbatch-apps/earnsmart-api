package services

import (
	"errors"
	"log"
	"math"
	"math/big"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/utils"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

type TransactionService struct {
	*BaseService
}

// Convert float to fixed-point integer representation with 24 decimal places
func floatToUint128(amount float64) tbt.Uint128 {
	precision := math.Pow10(24) // Fixed precision factor (24 decimal places)
	intAmount := uint64(amount * precision)
	return tbt.ToUint128(intAmount)
}

func uint128ToFloat(amount tbt.Uint128) float64 {
	precision := new(big.Float).SetFloat64(math.Pow10(24)) // Fixed precision factor (24 decimal places)
	intAmount := new(big.Int).SetBytes(amount[:])

	amountBig := new(big.Float).SetInt(intAmount)
	amountBig.Quo(amountBig, precision)

	result, _ := amountBig.Float64()
	return result
}

func NewTransactionService(db *gorm.DB, tbClient tb.Client) *TransactionService {
	return &TransactionService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *TransactionService) CreateDeposit(account models.Account, amount tbt.Uint128, currency string) ([]tbt.TransferEventResult, error) {

	localCurrency := models.AllCurrencies[currency]
	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)
	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	adminWallet, err := accountService.GetTreasuryWallet(currency)

	if err != nil {
		return nil, errors.New("admin wallet not found")
	}

	transfers := []tbt.Transfer{
		{
			ID:              models.GenerateUUID(),
			DebitAccountID:  adminWallet,
			CreditAccountID: utils.ToUint128(account.ID),
			Amount:          amount,
			Code:            models.TransferCodeDeposit,
			Ledger:          uint32(localCurrency.LedgerID),
			UserData128:     floatToUint128(amountPrice),
		},
	}

	transferResults, err := s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating deposit", err)
		return nil, err
	}

	log.Println("deposit created", transferResults)
	s.LogActivity(models.ActivityTypeUser, "Creating fund account", account.UserID)
	return transferResults, nil
}

func (s *TransactionService) CreateWithdrawal(account models.Account, amount tbt.Uint128, currency string) ([]tbt.TransferEventResult, error) {

	localCurrency := models.AllCurrencies[currency]
	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	if !accountService.AccountHasSufficientBalance(account, amount) {
		return nil, errors.New("insufficient balance")
	}

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	adminWallet, err := accountService.GetTreasuryWallet(currency)
	if err != nil {
		return nil, errors.New("admin wallet not found")
	}

	transfers := []tbt.Transfer{
		{
			ID:              models.GenerateUUID(),
			DebitAccountID:  utils.ToUint128(account.ID),
			CreditAccountID: adminWallet,
			Amount:          amount,
			Ledger:          uint32(localCurrency.LedgerID),
			UserData128:     floatToUint128(float64(amountPrice)),
			Code:            models.TransferCodeWithdraw,
		},
	}

	transferResults, err := s.tbClient.CreateTransfers(transfers)

	if err != nil {
		return nil, err
	}

	return transferResults, nil
}

func (s *TransactionService) CreateTransfer(creditAccount models.Account, debitAccount models.Account, amount tbt.Uint128, currency string, code uint16) ([]tbt.TransferEventResult, error) {

	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	if !accountService.AccountHasSufficientBalance(debitAccount, amount) {
		return nil, errors.New("insufficient balance")
	}

	localCurrency := models.AllCurrencies[currency]

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)

	transfers := []tbt.Transfer{
		{
			ID:              models.GenerateUUID(),
			CreditAccountID: tbt.ToUint128(uint64(creditAccount.ID)),
			DebitAccountID:  tbt.ToUint128(uint64(debitAccount.ID)),
			Amount:          amount,
			Ledger:          uint32(localCurrency.LedgerID),
			UserData128:     floatToUint128(float64(amountPrice)),
			UserData32:      uint32(models.TransactionTypeWithdraw),
			Code:            code,
		},
	}

	transferResults, err := s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating transfer", err)
		return nil, err
	}

	return transferResults, nil
}

type AccountTransferWithID struct {
	tbt.Transfer
	AmountUSD float64
	Currency  string
}

func (s *TransactionService) GetAllTransactions(accountIds []tbt.Uint128) ([]AccountTransferWithID, error) {
	var filter = tbt.AccountFilter{
		AccountID:    tbt.ToUint128(0),
		TimestampMin: 0,
		TimestampMax: 0,
		Limit:        10,
		Flags: tbt.AccountFilterFlags{
			Debits:  true,
			Credits: true,
		}.ToUint32(),
	}

	var allTransfers []AccountTransferWithID

	for _, accountId := range accountIds {
		filter.AccountID = accountId
		transfers, _ := s.tbClient.GetAccountTransfers(filter)

		for _, transfer := range transfers {
			transferWithID := AccountTransferWithID{
				Transfer:  transfer,
				AmountUSD: uint128ToFloat(transfer.UserData128),
				Currency:  models.LedgerCurrency[transfer.Ledger],
			}
			allTransfers = append(allTransfers, transferWithID)
		}

	}
	return allTransfers, nil
}
