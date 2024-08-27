package services

import (
	"errors"
	"log"
	"math/big"

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

func (s *AccountService) CreateWalletAccount(userID uint, currency string) (models.Account, error) {
	return s.CreateFundAccount(userID, currency, models.AccountCodeWallet)
}

func (s *AccountService) CreateFundAccount(userID uint, currency string, code models.AccountCode) (models.Account, error) {
	var accounts = []tbt.Account{}
	localCurrency := models.AllCurrencies[currency]

	var account = models.Account{
		UserID:      userID,
		Currency:    currency,
		AccountCode: code,
	}
	err := s.db.Create(&account).Error

	if err != nil {
		log.Println("Error creating account: ", err.Error())
		return account, err
	} else {
		log.Println("Created db account: ", account.ID)
	}

	accounts = append(accounts, tbt.Account{
		ID:         tbt.ToUint128(uint64(account.ID)),
		Code:       uint16(code),
		Ledger:     uint32(localCurrency.LedgerID),
		UserData64: uint64(userID),
		Flags:      tbt.AccountFlags{History: true}.ToUint16(),
	})

	// utils.LogAccount(accounts[0])

	result, err := s.tbClient.CreateAccounts(accounts)

	if err != nil {
		log.Println("Error creating account: ", err.Error())
		return account, err
	} else {
		log.Println("Created tigerbeetle account: ", account.ID)
		log.Printf("Result: %+v\n", result)
	}
	s.LogActivity(models.ActivityTypeAdmin, "Creating fund account", userID)
	return account, nil
}

func (s *AccountService) ExtractIDs(accounts []models.Account) ([]tbt.Uint128, error) {
	var accountIds = []tbt.Uint128{}
	for _, account := range accounts {
		accountIds = append(accountIds, tbt.ToUint128(uint64(account.ID)))
	}
	return accountIds, nil
}

type AccountBalanceWithID struct {
	tbt.AccountBalance
	AccountID tbt.Uint128
}

func (s *AccountService) LookupAccountBalances(accountIds []tbt.Uint128) ([]AccountBalanceWithID, error) {

	var balances = []AccountBalanceWithID{}

	//balance, err := s.tbClient.GetAccountBalances()
	var filter = tbt.AccountFilter{
		AccountID: tbt.ToUint128(0),
		// TimestampMin: 0,
		// TimestampMax: 0,
		Limit: 10,
		Flags: tbt.AccountFilterFlags{
			Debits:  true,
			Credits: true,
		}.ToUint32(),
	}

	// utils.LogAccountIDs(accountIds)

	for _, accountId := range accountIds {
		filter.AccountID = accountId

		// utils.LogAccountFilter(filter)

		balance, err := s.tbClient.GetAccountBalances(filter)
		// utils.LogJson("Balance:", balance)

		// utils.LogAccountBalance(balance[0])

		if err != nil {
			log.Println("Error getting account balances: ", err)
			return balances, err
		}

		if len(balance) == 0 {
			balance = []tbt.AccountBalance{
				{
					DebitsPending:  tbt.ToUint128(0),
					DebitsPosted:   tbt.ToUint128(0),
					CreditsPending: tbt.ToUint128(0),
					CreditsPosted:  tbt.ToUint128(0),
					Timestamp:      0,
				},
			}
		}

		balanceWithID := AccountBalanceWithID{
			AccountBalance: balance[0],
			AccountID:      accountId,
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

func (s *AccountService) GetTreasuryWallet(currency string) (tbt.Uint128, error) {
	var account models.Account
	s.db.Where(models.Account{UserID: 1, Currency: currency, AccountCode: models.AccountCodeWallet}).First(&account)

	if account.ID == 0 {
		return tbt.ToUint128(0), errors.New("treasury wallet not found")
	}
	return tbt.ToUint128(uint64(account.ID)), nil

}
