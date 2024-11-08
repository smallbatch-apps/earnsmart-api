package services

import (
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

type Services struct {
	Db          *gorm.DB
	TbClient    tb.Client
	Account     *AccountService
	Activity    *ActivityService
	Allocation  *AllocationService
	Fund        *FundService
	Price       *PriceService
	Setting     *SettingService
	Swap        *SwapService
	Transaction *TransactionService
	Transfer    *TransferService
	User        *UserService
}

func NewServices(db *gorm.DB, tbClient tb.Client) *Services {
	return &Services{
		Db:          db,
		TbClient:    tbClient,
		Account:     NewAccountService(db, tbClient),
		Activity:    NewActivityService(db, tbClient),
		Allocation:  NewAllocationService(db, tbClient),
		Fund:        NewFundService(db, tbClient),
		Price:       NewPriceService(db, tbClient),
		Setting:     NewSettingService(db, tbClient),
		Swap:        NewSwapService(db, tbClient),
		Transaction: NewTransactionService(db, tbClient),
		Transfer:    NewTransferService(db, tbClient),
		User:        NewUserService(db, tbClient),
	}
}
