package services

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

type TransferService struct {
	*BaseService
}

// Convert float to fixed-point integer representation with 24 decimal places
func PriceToUint128(amount float64) tbt.Uint128 {
	log.Printf("PriceToUint128 input: %v", amount)

	precision := math.Pow10(6)
	precisionResult := amount * precision

	precisionResult = math.Round(precisionResult)
	log.Printf("After precision multiplication: %v", precisionResult)

	// Check for potential overflow
	if precisionResult >= float64(math.MaxUint64) {
		log.Printf("WARNING: Value will overflow uint64 max (%d)", uint64(math.MaxUint64))
	}

	intAmount := uint64(precisionResult)
	log.Printf("After uint64 conversion: %v", intAmount)

	result := tbt.ToUint128(intAmount)
	log.Printf("Final Uint128 bytes: %v", result)

	return result
}

func UserData128ToPrice(value tbt.Uint128) float64 {
	precision := math.Pow10(6)
	numericValue := float64(binary.LittleEndian.Uint64(value[:8]))
	return numericValue / precision
}

func stringToUint128(amount string) (tbt.Uint128, error) {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return tbt.Uint128{}, err
	}
	return tbt.BigIntToUint128(*amountDecimal.BigInt()), nil
}

func NewTransferService(db *gorm.DB, tbClient tb.Client) *TransferService {
	return &TransferService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *TransferService) CreateSubscribeTransfer(account models.Account, amount tbt.Uint128, currency string) (tbt.Uint128, error) {
	localCurrency := models.AllCurrencies[currency]
	accountService := NewAccountService(s.db, s.tbClient)
	priceService := NewPriceService(s.db, s.tbClient)

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	userWallet, err := accountService.GetUserWallet(account.UserID, currency)
	if err != nil {
		return tbt.Uint128{}, err
	}
	id := models.GenerateUUID()
	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  userWallet.TbID(),
			CreditAccountID: account.TbID(),
			Amount:          amount,
			Code:            uint16(models.TransferCodeSubscribe),
			Ledger:          localCurrency.LedgerID,
			UserData128:     PriceToUint128(amountPrice),
		},
	}

	_, err = s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating subscription", err)
		return tbt.Uint128{}, err
	}

	return id, nil
}

func (s *TransferService) CreateRedeemTransfer(account models.Account, amount tbt.Uint128, currency string) (tbt.Uint128, error) {
	localCurrency := models.AllCurrencies[currency]
	accountService := NewAccountService(s.db, s.tbClient)
	priceService := NewPriceService(s.db, s.tbClient)

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	userWallet, err := accountService.GetUserWallet(account.UserID, currency)
	if err != nil {
		return tbt.Uint128{}, err
	}
	id := models.GenerateUUID()
	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  account.TbID(),
			CreditAccountID: userWallet.TbID(),
			Amount:          amount,
			Code:            uint16(models.TransferCodeRedeem),
			Ledger:          localCurrency.LedgerID,
			UserData128:     PriceToUint128(amountPrice),
		},
	}

	_, err = s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating redemption", err)
		return tbt.Uint128{}, err
	}

	return id, nil
}

func (s *TransferService) CreateDepositTransfer(account models.Account, amount tbt.Uint128, currency string) (tbt.Uint128, error) {

	localCurrency := models.AllCurrencies[currency]
	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)
	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	adminWallet, err := accountService.GetTreasuryWallet(currency)

	if err != nil {
		return tbt.Uint128{}, errors.New("admin wallet not found")
	}

	id := models.GenerateUUID()

	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  adminWallet.TbID(),
			CreditAccountID: account.TbID(),
			Amount:          amount,
			Code:            uint16(models.TransferCodeDeposit),
			Ledger:          localCurrency.LedgerID,
			UserData128:     PriceToUint128(amountPrice),
		},
	}

	transferResults, err := s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating deposit", err)
		return tbt.Uint128{}, err
	}

	log.Println("deposit created", transferResults)

	amountFloat, err := priceService.AmountToFloat(currency, amount)
	if err != nil {
		return tbt.Uint128{}, err
	}

	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Confirm deposit: %f%s", amountFloat, currency), account.UserID)
	return id, nil
}

func (s *TransferService) ConfirmTransfer(transfer tbt.Transfer) (tbt.Uint128, error) {
	id := models.GenerateUUID()
	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  transfer.DebitAccountID,
			CreditAccountID: transfer.CreditAccountID,
			Amount:          transfer.Amount,
			Ledger:          transfer.Ledger,
			PendingID:       transfer.ID,
			UserData128:     transfer.UserData128,
			Code:            transfer.Code,
			Flags:           tbt.TransferFlags{PostPendingTransfer: true}.ToUint16(),
		},
	}

	_, err := s.tbClient.CreateTransfers(transfers)
	if err != nil {
		return tbt.Uint128{}, err
	}

	return id, nil
}

func (s *TransferService) VoidTransfer(transfer tbt.Transfer) (tbt.Uint128, error) {
	id := models.GenerateUUID()
	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  transfer.DebitAccountID,
			CreditAccountID: transfer.CreditAccountID,
			Amount:          transfer.Amount,
			Ledger:          transfer.Ledger,
			PendingID:       transfer.ID,
			UserData128:     transfer.UserData128,
			Code:            transfer.Code,
			Flags:           tbt.TransferFlags{VoidPendingTransfer: true}.ToUint16(),
		},
	}

	_, err := s.tbClient.CreateTransfers(transfers)
	if err != nil {
		return tbt.Uint128{}, err
	}

	return id, nil
}

func (s *TransferService) CreateWithdrawalTransfer(account models.Account, amount tbt.Uint128, currency string) (tbt.Uint128, error) {

	localCurrency := models.AllCurrencies[currency]
	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	if !accountService.AccountHasSufficientBalance(account, amount) {
		return tbt.Uint128{}, errors.New("insufficient balance")
	}

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	adminWallet, err := accountService.GetTreasuryWallet(currency)
	if err != nil {
		return tbt.Uint128{}, errors.New("admin wallet not found")
	}

	id := models.GenerateUUID()

	transfers := []tbt.Transfer{
		{
			ID:              id,
			DebitAccountID:  account.TbID(),
			CreditAccountID: adminWallet.TbID(),
			Amount:          amount,
			Ledger:          localCurrency.LedgerID,
			UserData128:     PriceToUint128(amountPrice),
			Code:            uint16(models.TransferCodeWithdraw),
		},
	}

	_, err = s.tbClient.CreateTransfers(transfers)
	if err != nil {
		return tbt.Uint128{}, err
	}

	amountFloat, err := priceService.AmountToFloat(currency, amount)
	if err != nil {
		return tbt.Uint128{}, err
	}

	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Confirm withdrawal: %f%s", amountFloat, currency), account.UserID)
	return id, nil
}

func (s *TransferService) CreateTransfer(creditAccount models.Account, debitAccount models.Account, amount tbt.Uint128, currency string, code uint16) (tbt.Uint128, error) {

	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	if !accountService.AccountHasSufficientBalance(debitAccount, amount) {
		return tbt.Uint128{}, errors.New("insufficient balance")
	}

	localCurrency := models.AllCurrencies[currency]

	amountPrice, _ := priceService.GetAmountPrice(currency, amount)
	id := models.GenerateUUID()
	transfers := []tbt.Transfer{
		{
			ID:              id,
			CreditAccountID: tbt.ToUint128(uint64(creditAccount.ID)),
			DebitAccountID:  tbt.ToUint128(uint64(debitAccount.ID)),
			Amount:          amount,
			Ledger:          uint32(localCurrency.LedgerID),
			UserData128:     PriceToUint128(amountPrice),
			UserData32:      uint32(models.TransactionTypeWithdraw),
			Code:            code,
			Flags:           tbt.TransferFlags{Pending: true}.ToUint16(),
		},
	}

	_, err := s.tbClient.CreateTransfers(transfers)

	if err != nil {
		log.Println("error creating transfer", err)
		return tbt.Uint128{}, err
	}

	return id, nil
}

type AccountTransferWithID struct {
	tbt.Transfer
	AmountUSD float64
	Currency  string
}

func (t AccountTransferWithID) MarshalJSON() ([]byte, error) {
	type TransferResponse struct {
		ID            string  `json:"id"`
		Amount        string  `json:"amount"`
		AmountUSD     float64 `json:"amount_usd"`
		Currency      string  `json:"currency"`
		Code          uint16  `json:"code"`
		Timestamp     int64   `json:"timestamp"`
		DebitAccount  uint64  `json:"debit_account"`
		CreditAccount uint64  `json:"credit_account"`
	}

	id := uuid.UUID(t.ID)

	amountBig := t.Amount.BigInt()
	amount := amountBig.String()
	debitAccount := binary.LittleEndian.Uint64(t.DebitAccountID[:8])
	creditAccount := binary.LittleEndian.Uint64(t.CreditAccountID[:8])

	response := TransferResponse{
		ID:            id.String(),
		Amount:        amount,
		AmountUSD:     t.AmountUSD,
		Currency:      t.Currency,
		Code:          uint16(t.Code),
		Timestamp:     int64(t.Timestamp),
		DebitAccount:  debitAccount,
		CreditAccount: creditAccount,
	}

	return json.Marshal(response)
}

func (s *TransferService) ListTransfers(accountIds []tbt.Uint128) ([]AccountTransferWithID, error) {
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
				AmountUSD: UserData128ToPrice(transfer.UserData128),
				Currency:  models.LedgerCurrency[transfer.Ledger],
			}
			allTransfers = append(allTransfers, transferWithID)
		}

	}
	return allTransfers, nil
}

func (s *TransferService) GetTransferByID(transferID tbt.Uint128) (tbt.Transfer, error) {
	transfers, err := s.tbClient.LookupTransfers([]tbt.Uint128{transferID})
	if err != nil {
		return tbt.Transfer{}, err
	}
	return transfers[0], nil
}
