package services

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

type TransactionService struct {
	*BaseService
}

func NewTransactionService(db *gorm.DB, tbClient tb.Client) *TransactionService {
	return &TransactionService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *TransactionService) CreateSubscription(userID uint64, amount string, currency string, accountCode models.AccountCode) (models.Transaction, error) {
	accountService := NewAccountService(s.db, s.tbClient)
	transferService := NewTransferService(s.db, s.tbClient)

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return models.Transaction{}, err
	}
	account, err := accountService.GetOrCreateAccount(models.Account{
		AccountCode:  accountCode,
		Currency:     currency,
		OwnableModel: models.OwnableModel{UserID: userID},
	})
	if err != nil {
		return models.Transaction{}, err
	}

	transferID, err := transferService.CreateSubscribeTransfer(account, tbt.BigIntToUint128(*amountDecimal.BigInt()), currency)
	if err != nil {
		return models.Transaction{}, err
	}

	transaction := models.Transaction{
		Type:         models.TransactionTypeRedeem,
		Status:       models.TransactionStatusConfirmed,
		Amount:       amountDecimal,
		Currency:     currency,
		TransferID:   transferID.String(),
		OwnableModel: models.OwnableModel{UserID: userID},
	}

	err = s.db.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, nil
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Create new fund subscription for %s%s", amountDecimal, currency), userID)
	return transaction, err
}

func (s *TransactionService) CreateRedemption(userID uint64, amount string, currency string, accountCode models.AccountCode) (models.Transaction, error) {
	accountService := NewAccountService(s.db, s.tbClient)
	transferService := NewTransferService(s.db, s.tbClient)

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return models.Transaction{}, err
	}
	account, err := accountService.GetOrCreateAccount(models.Account{
		AccountCode:  accountCode,
		Currency:     currency,
		OwnableModel: models.OwnableModel{UserID: userID},
	})
	if err != nil {
		return models.Transaction{}, err
	}

	transferID, err := transferService.CreateRedeemTransfer(account, tbt.BigIntToUint128(*amountDecimal.BigInt()), currency)
	if err != nil {
		return models.Transaction{}, err
	}

	transaction := models.Transaction{
		OwnableModel: models.OwnableModel{UserID: userID},
		Type:         models.TransactionTypeRedeem,
		Status:       models.TransactionStatusConfirmed,
		Amount:       amountDecimal,
		Currency:     currency,
		TransferID:   transferID.String(),
	}

	err = s.db.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, nil
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Redeem %s%s from fund", amountDecimal, currency), userID)
	return transaction, nil
}

func (s *TransactionService) CreatePendingDeposit(userID uint64, amount string, depositAddress string, currency string) (models.Transaction, error) {
	priceService := NewPriceService(s.db, s.tbClient)
	transferService := NewTransferService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return models.Transaction{}, err
	}

	account, err := accountService.GetOrCreateAccount(models.Account{
		AccountCode:  models.AccountCodeWallet,
		Currency:     currency,
		OwnableModel: models.OwnableModel{UserID: userID},
	})
	if err != nil {
		return models.Transaction{}, err
	}

	transferID, err := transferService.CreateDepositTransfer(account, tbt.BigIntToUint128(*amountDecimal.BigInt()), currency)
	if err != nil {
		return models.Transaction{}, err
	}

	transaction := models.Transaction{
		Type:         models.TransactionTypeWithdraw,
		Status:       models.TransactionStatusPending,
		Amount:       amountDecimal,
		Currency:     currency,
		TransferID:   transferID.String(),
		Address:      depositAddress,
		OwnableModel: models.OwnableModel{UserID: userID},
	}

	err = s.db.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, err
	}
	amountFloat, err := priceService.AmountToFloat(currency, transaction.AmountAsUint128())
	if err != nil {
		return models.Transaction{}, err
	}

	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Create pending deposit: %f%s", amountFloat, currency), userID)
	return transaction, nil
}

func (s *TransactionService) CreatePendingWithdrawal(userID uint64, amount string, withdrawalAddress string, currency string) (models.Transaction, error) {
	priceService := NewPriceService(s.db, s.tbClient)
	transferService := NewTransferService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return models.Transaction{}, err
	}

	account, err := accountService.GetOrCreateAccount(models.Account{
		AccountCode:  models.AccountCodeWallet,
		Currency:     currency,
		OwnableModel: models.OwnableModel{UserID: userID},
	})
	if err != nil {
		return models.Transaction{}, err
	}

	transferID, _ := transferService.CreateDepositTransfer(account, tbt.BigIntToUint128(*amountDecimal.BigInt()), currency)

	transaction := models.Transaction{
		Type:       models.TransactionTypeWithdraw,
		Status:     models.TransactionStatusPending,
		Amount:     amountDecimal,
		Currency:   currency,
		TransferID: transferID.String(),
		Address:    withdrawalAddress,
	}

	err = s.db.Create(&transaction).Error
	if err != nil {
		return models.Transaction{}, err
	}
	amountFloat, err := priceService.AmountToFloat(currency, transaction.AmountAsUint128())
	if err != nil {
		return models.Transaction{}, err
	}

	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Creating pending withdrawal: %f%s", amountFloat, currency), userID)
	return transaction, nil
}

func (s *TransactionService) ApproveTransaction(transaction models.Transaction) (models.Transaction, error) {
	transferService := NewTransferService(s.db, s.tbClient)

	err := s.db.Model(&transaction).Updates(models.Transaction{Status: models.TransactionStatusConfirmed}).Error
	if err != nil {
		return transaction, err
	}
	transfer, err := transferService.GetTransferByID(transaction.TransferIDAsUint128())
	if err != nil {
		return transaction, err
	}

	_, err = transferService.ConfirmTransfer(transfer)
	if err != nil {
		return transaction, err
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Approved transaction: %s%s", transaction.Amount, transaction.Currency), transaction.UserID)
	return transaction, err
}

func (s *TransactionService) DeclineTransaction(transaction models.Transaction, status models.TransactionStatus) (models.Transaction, error) {
	transferService := NewTransferService(s.db, s.tbClient)

	err := s.db.Model(&transaction).Updates(models.Transaction{Status: status}).Error
	if err != nil {
		return transaction, err
	}

	transfer, err := transferService.GetTransferByID(transaction.TransferIDAsUint128())
	if err != nil {
		return transaction, err
	}

	_, err = transferService.VoidTransfer(transfer)
	if err != nil {
		return transaction, err
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Declined transaction: %s%s", transaction.Amount, transaction.Currency), transaction.UserID)
	return transaction, err
}
