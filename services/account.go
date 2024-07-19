package services

import (
	"math/big"

	"github.com/google/uuid"
	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

type AccountService struct {
	*BaseService
}

func NewAccountService(db *gorm.DB, tbClient tb.Client) *AccountService {
	return &AccountService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *AccountService) GetAccounts(searchParams models.Account) ([]models.Account, error) {
	var accounts = []models.Account{}
	s.db.Where(&searchParams).Find(&accounts)
	return accounts, nil
}

func (s *AccountService) GetOrCreateAccount(accountSearch models.Account) (models.Account, error) {
	var account models.Account
	s.db.Where(accountSearch).First(&account)

	if account.ID == 0 {
		account_id, err := s.CreateFundAccount(accountSearch.UserID, accountSearch.Currency, accountSearch.AccountCode)
		if err != nil {
			return account, err
		}
		s.db.Where("account_id = ?", account_id).First(&account)
		return account, nil
	}
	return account, nil
}

func (s *AccountService) CreateWalletAccount(userID uint, currency string) (uuid.UUID, error) {
	return s.CreateFundAccount(userID, currency, models.AccountCodeWallet)
}

func (s *AccountService) CreateFundAccount(userID uint, currency string, code models.AccountCode) (uuid.UUID, error) {
	var accounts = []tbt.Account{}

	local_currency := models.AllCurrencies[currency]

	account_id, _ := uuid.NewV7()

	accounts = append(accounts, tbt.Account{
		ID:         models.ConvertUUIDToUint128(account_id),
		Code:       uint16(code),
		Ledger:     uint32(local_currency.LedgerID),
		UserData64: uint64(userID),
		Flags:      tbt.AccountFlags{History: true}.ToUint16(),
	})

	_, err := s.tbClient.CreateAccounts(accounts)

	if err != nil {
		return account_id, err
	}

	var account = models.Account{
		AccountID:   account_id,
		UserID:      userID,
		Currency:    currency,
		AccountCode: code,
	}
	err = s.db.Create(&account).Error

	if err != nil {
		return account_id, err
	}

	return account_id, nil
}

func (s *AccountService) ExtractIDs(accounts []models.Account) ([]tbt.Uint128, error) {
	var account_ids = []tbt.Uint128{}
	for _, account := range accounts {
		account_ids = append(account_ids, models.ConvertUUIDToUint128(account.AccountID))
	}
	return account_ids, nil
}

type AccountBalanceWithID struct {
	tbt.AccountBalance
	AccountID uuid.UUID
}

func (s *AccountService) LookupAccountBalances(accountIds []tbt.Uint128) ([]AccountBalanceWithID, error) {

	var balances = []AccountBalanceWithID{}

	//balance, err := s.tbClient.GetAccountBalances()
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

	for _, accountId := range accountIds {
		filter.AccountID = accountId
		balance, err := s.tbClient.GetAccountBalances(filter)

		if err != nil {
			return balances, err
		}

		balanceWithID := AccountBalanceWithID{
			AccountBalance: balance[0],
			AccountID:      models.ConvertUint128ToUUID(accountId),
		}
		balances = append(balances, balanceWithID)
	}
	return balances, nil
}

func (s *AccountService) AccountHasSufficientBalance(account models.Account, amount tbt.Uint128) bool {
	accountIds, err := s.ExtractIDs([]models.Account{account})
	if err != nil {
		return false
	}

	accountBalances, err := s.LookupAccountBalances(accountIds)

	if err != nil {
		return false
	}

	if len(accountBalances) == 0 {
		return false
	}

	credits := accountBalances[0].CreditsPosted.BigInt()
	debits := accountBalances[0].DebitsPosted.BigInt()
	balance := new(big.Int).Sub(&credits, &debits)

	amountBigInt := amount.BigInt()

	result := new(big.Int).Sub(balance, &amountBigInt)
	return result.Sign() >= 0
}
