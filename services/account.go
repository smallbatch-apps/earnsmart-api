package services

import (
	"encoding/json"
	"errors"
	"fmt"
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
	err := s.db.Where(searchParams).Find(&accounts).Error
	return accounts, err
}

func (s *AccountService) GetOrCreateAccount(accountSearch models.Account) (models.Account, error) {
	var account models.Account
	s.db.Where(accountSearch).First(&account)

	if account.ID == 0 {
		account, err := s.CreateFundAccount(accountSearch.UserID, accountSearch.Currency, accountSearch.AccountCode)
		return account, err
	}
	return account, nil
}

func (s *AccountService) CreateWalletAccount(userID uint64, currency string) (models.Account, error) {
	return s.CreateFundAccount(userID, currency, models.AccountCodeWallet)
}

func (s *AccountService) CreateFundAccount(userID uint64, currency string, code models.AccountCode) (models.Account, error) {
	var accounts = []tbt.Account{}
	localCurrency := models.AllCurrencies[currency]

	var account = models.Account{
		OwnableModel: models.OwnableModel{UserID: userID},
		Currency:     currency,
		AccountCode:  code,
	}
	err := s.db.Create(&account).Error

	if err != nil {
		log.Println("Error creating account: ", err.Error())
		return account, err
	} else {
		log.Println("Created db account: ", account.ID)
	}

	flags := tbt.AccountFlags{History: true}
	// Treasury accounts (userID 1) can go negative,
	// all others must maintain positive balance
	if userID != 1 {
		flags.DebitsMustNotExceedCredits = true
	}

	accounts = append(accounts, tbt.Account{
		ID:         account.TbID(),
		Code:       uint16(code),
		Ledger:     localCurrency.LedgerID,
		UserData64: userID,
		Flags:      flags.ToUint16(),
	})

	result, err := s.tbClient.CreateAccounts(accounts)

	if err != nil {
		log.Println("Error creating account: ", err.Error())
		return account, err
	} else {
		log.Println("Created tigerbeetle account: ", account.ID)
		log.Printf("Result: %+v\n", result)
	}

	accountType := "fund"
	if code == models.AccountCodeWallet {
		accountType = "wallet"
	}

	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Creating %s account for %s", accountType, currency), userID)
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

func (s *AccountService) GetAccountBalanceHistory(accountIds []tbt.Uint128) ([]AccountBalanceWithID, error) {

	var balances = []AccountBalanceWithID{}

	var filter = tbt.AccountFilter{
		AccountID: tbt.ToUint128(0),
		Limit:     10,
		Flags: tbt.AccountFilterFlags{
			Debits:  true,
			Credits: true,
		}.ToUint32(),
	}

	for _, accountId := range accountIds {
		filter.AccountID = accountId

		balance, err := s.tbClient.GetAccountBalances(filter)
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

type AccountWithBalance struct {
	tbt.Account
	Currency          string
	Balance           tbt.Uint128
	BalanceUSD        float64
	BalancePending    tbt.Uint128
	BalancePendingUSD float64
}

func (a AccountWithBalance) MarshalJSON() ([]byte, error) {
	type AccountResponse struct {
		ID                string  `json:"id"`
		Ledger            uint32  `json:"ledger"`
		Currency          string  `json:"currency"`
		Code              uint16  `json:"code"`
		Timestamp         uint64  `json:"timestamp"`
		Balance           string  `json:"balance"`
		BalanceUSD        float64 `json:"balance_usd"`
		BalancePending    string  `json:"balance_pending"`
		BalancePendingUSD float64 `json:"balance_pending_usd"`
	}
	balanceBig := a.Balance.BigInt()
	balance := balanceBig.String()

	balancePendingBig := a.BalancePending.BigInt()
	balancePending := balancePendingBig.String()

	response := AccountResponse{
		ID:                a.ID.String(),
		Ledger:            a.Ledger,
		Currency:          a.Currency,
		Code:              a.Code,
		Timestamp:         a.Timestamp,
		Balance:           balance,
		BalanceUSD:        a.BalanceUSD,
		BalancePending:    balancePending,
		BalancePendingUSD: a.BalancePendingUSD,
	}

	return json.Marshal(response)
}

func (s *AccountService) LookupAccounts(accountIds []tbt.Uint128) ([]AccountWithBalance, error) {
	priceService := NewPriceService(s.db, s.tbClient)

	var balances = []AccountWithBalance{}

	prices, err := priceService.ListPriceMap()
	if err != nil {
		return balances, err
	}

	accounts, err := s.tbClient.LookupAccounts(accountIds)
	if err != nil {
		return balances, err
	}

	for _, account := range accounts {
		currency := models.LedgerCurrency[account.Ledger]
		rate := prices[currency]
		accountWithBalance := s.GetBalancesForAccount(account, rate)
		balances = append(balances, accountWithBalance)
	}

	return balances, nil
}

func (s *AccountService) GetBalancesForAccount(account tbt.Account, rate float64) AccountWithBalance {
	priceService := NewPriceService(s.db, s.tbClient)
	currency := models.LedgerCurrency[account.Ledger]
	credits := account.CreditsPosted.BigInt()
	debits := account.DebitsPosted.BigInt()
	balance := new(big.Int).Sub(&credits, &debits)
	balanceUint128 := tbt.BigIntToUint128(*balance)
	balanceUSD := priceService.AmountToUSD(currency, rate, balanceUint128)

	creditsPending := account.CreditsPending.BigInt()
	debitsPending := account.DebitsPending.BigInt()
	balancePending := new(big.Int).Sub(&creditsPending, &debitsPending)
	balancePendingUint128 := tbt.BigIntToUint128(*balancePending)
	balancePendingUSD := priceService.AmountToUSD(currency, rate, balancePendingUint128)

	accountWithBalance := AccountWithBalance{
		Currency:          currency,
		Account:           account,
		Balance:           balanceUint128,
		BalanceUSD:        balanceUSD,
		BalancePending:    balancePendingUint128,
		BalancePendingUSD: balancePendingUSD,
	}

	return accountWithBalance
}

func (s *AccountService) CombineAccountBalances(accounts []AccountWithBalance) []AccountWithBalance {
	ledgerBalances := make(map[uint32]AccountWithBalance)

	for _, account := range accounts {
		if combined, exists := ledgerBalances[account.Ledger]; exists {
			credits := account.Balance.BigInt()
			existingCredits := combined.Balance.BigInt()
			totalBalance := new(big.Int).Add(&credits, &existingCredits)

			pendingCredits := account.BalancePending.BigInt()
			existingPending := combined.BalancePending.BigInt()
			totalPending := new(big.Int).Add(&pendingCredits, &existingPending)

			totalBalanceUSD := combined.BalanceUSD + account.BalanceUSD
			totalPendingUSD := combined.BalancePendingUSD + account.BalancePendingUSD

			combined.Balance = tbt.BigIntToUint128(*totalBalance)
			combined.BalancePending = tbt.BigIntToUint128(*totalPending)
			combined.BalanceUSD = totalBalanceUSD
			combined.BalancePendingUSD = totalPendingUSD

			ledgerBalances[account.Ledger] = combined
		} else {
			ledgerBalances[account.Ledger] = account
		}
	}

	var result []AccountWithBalance
	for _, account := range ledgerBalances {
		result = append(result, account)
	}

	return result
}

func (s *AccountService) AccountHasSufficientBalance(account models.Account, amount tbt.Uint128) bool {
	accountIds, err := s.ExtractIDs([]models.Account{account})
	if err != nil {
		return false
	}

	accountBalances, err := s.GetAccountBalanceHistory(accountIds)

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

func (s *AccountService) GetTreasuryWallet(currency string) (models.Account, error) {
	var account models.Account
	searchModel := models.Account{OwnableModel: models.OwnableModel{UserID: 1}, Currency: currency, AccountCode: models.AccountCodeWallet}

	s.db.Where(&searchModel).First(&account)

	if account.ID == 0 {
		return account, errors.New("treasury wallet not found")
	}
	return account, nil
}

func (s *AccountService) GetUserWallet(userID uint64, currency string) (models.Account, error) {
	var account models.Account

	s.db.Where(models.Account{OwnableModel: models.OwnableModel{UserID: userID}, Currency: currency, AccountCode: models.AccountCodeWallet}).First(&account)

	if account.ID == 0 {
		return account, errors.New("currency wallet not found")
	}
	return account, nil
}
