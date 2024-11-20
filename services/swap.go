package services

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/utils"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

// SettingService extends BaseService for settings
type SwapService struct {
	*BaseService
}

func NewSwapService(db *gorm.DB, tbClient tb.Client) *SwapService {
	return &SwapService{BaseService: NewBaseService(db, tbClient)}
}

func (s *SwapService) CreateSwap(userID uint64, amountFrom string, currencyFrom string, amountTo string, currencyTo string, rate float64) (models.Swap, error) {
	priceService := NewPriceService(s.db, s.tbClient)
	fromTreasury, toTreasury, fromUser, toUser, err := s.getSwapWallets(userID, currencyFrom, currencyTo)
	if err != nil {
		return models.Swap{}, err
	}

	currencyFromLedgerID := models.AllCurrencies[currencyFrom].LedgerID
	currencyToLedgerID := models.AllCurrencies[currencyTo].LedgerID

	fromTransferID := models.GenerateUUID()
	toTransferID := models.GenerateUUID()

	amountFrom128, err := stringToUint128(amountFrom)
	if err != nil {
		return models.Swap{}, err
	}
	amountTo128, err := stringToUint128(amountTo)
	if err != nil {
		return models.Swap{}, err
	}
	amountFromDecimal, err := decimal.NewFromString(amountFrom)
	if err != nil {
		return models.Swap{}, err
	}
	amountToDecimal, err := decimal.NewFromString(amountTo)
	if err != nil {
		return models.Swap{}, err
	}

	amountFromPrice, _ := priceService.GetAmountPrice(currencyFrom, amountFrom128)
	amountToPrice, _ := priceService.GetAmountPrice(currencyTo, amountTo128)

	accountService := NewAccountService(s.db, s.tbClient)
	if !accountService.AccountHasSufficientBalance(fromUser, amountFrom128) {
		return models.Swap{}, errors.New("insufficient balance")
	}

	transfers := []tbt.Transfer{
		{
			ID:              fromTransferID,
			Ledger:          currencyFromLedgerID,
			DebitAccountID:  fromUser.TbID(),
			Amount:          amountFrom128,
			CreditAccountID: fromTreasury.TbID(),
			Flags:           tbt.TransferFlags{Linked: true}.ToUint16(),
			UserData128:     PriceToUint128(amountFromPrice),
		},
		{
			ID:              toTransferID,
			Ledger:          currencyToLedgerID,
			Amount:          amountTo128,
			DebitAccountID:  toTreasury.TbID(),
			CreditAccountID: toUser.TbID(),
			Flags:           tbt.TransferFlags{Linked: true}.ToUint16(),
			UserData128:     PriceToUint128(amountToPrice),
		},
	}

	_, err = s.tbClient.CreateTransfers(transfers)
	if err != nil {
		return models.Swap{}, err
	}

	swap := models.Swap{
		OwnableModel:   models.OwnableModel{UserID: userID},
		FromAmount:     amountFromDecimal,
		FromCurrency:   currencyFrom,
		ToAmount:       amountToDecimal,
		ToCurrency:     currencyTo,
		Rate:           rate,
		FromTransferID: fromTransferID.String(),
		ToTransferID:   toTransferID.String(),
	}

	err = s.db.Create(&swap).Error

	formattedFromAmount := utils.FormatCurrencyAmount(amountFromDecimal, models.AllCurrencies[currencyFrom].Decimals)
	formattedToAmount := utils.FormatCurrencyAmount(amountToDecimal, models.AllCurrencies[currencyTo].Decimals)
	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Created swap: %s %s to %s %s", formattedFromAmount, currencyFrom, formattedToAmount, currencyTo), userID)

	return swap, err
}

func (s *SwapService) ListSwaps(userID uint64) ([]models.Swap, error) {
	var swaps []models.Swap
	err := s.db.Where("user_id = ?", userID).Find(&swaps).Order("created_at DESC").Error
	return swaps, err
}

func (s *SwapService) GetSwap(id uint64) (models.Swap, error) {
	var swap models.Swap
	err := s.db.First(&swap, id).Error
	return swap, err
}

func (s *SwapService) getSwapWallets(userID uint64, currencyFrom, currencyTo string) (fromTreasury, toTreasury, fromUser, toUser models.Account, err error) {
	accountService := NewAccountService(s.db, s.tbClient)
	if fromTreasury, err = accountService.GetTreasuryWallet(currencyFrom); err != nil {
		return
	}
	if toTreasury, err = accountService.GetTreasuryWallet(currencyTo); err != nil {
		return
	}
	if fromUser, err = accountService.GetUserWallet(userID, currencyFrom); err != nil {
		return
	}
	if toUser, err = accountService.GetUserWallet(userID, currencyTo); err != nil {
		return
	}
	return
}
