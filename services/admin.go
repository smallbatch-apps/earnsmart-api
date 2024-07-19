package services

import (
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"gorm.io/gorm"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

type AdminService struct {
	*BaseService
}

func NewAdminService(db *gorm.DB, tbClient tb.Client) *AdminService {
	return &AdminService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *AdminService) SeedData() error {

	s.db.AutoMigrate(&models.Price{}, &models.Fund{}, &models.Setting{}, &models.User{})

	price_service := NewPriceService(s.db, s.tbClient)
	account_service := NewAccountService(s.db, s.tbClient)
	transaction_service := NewTransactionService(s.db, s.tbClient)

	// populate prices
	price_service.UpdatePrices()

	// create admin user
	admin_user := models.User{
		Name:     "Michael Smith",
		Email:    "m.smith@earnsmart.com",
		Password: os.Getenv("ADMIN_USER_PASSWORD"),
	}
	err := s.db.Create(admin_user).Error
	if err != nil {
		return err
	}

	// create demo user
	test_user := models.User{
		Name:     "Test User",
		Email:    "test@earnsmart.com",
		Password: "123456",
	}

	err = s.db.Create(test_user).Error
	if err != nil {
		return err
	}

	test_account_id, err := account_service.CreateWalletAccount(test_user.ID, "ETH")
	if err != nil {
		return err
	}
	transaction_service.CreateDeposit(test_account_id, tbt.ToUint128(2997924580000000000), "ETH")

	test_account_id, err = account_service.CreateWalletAccount(test_user.ID, "USDT")
	if err != nil {
		return err
	}
	transaction_service.CreateDeposit(test_account_id, tbt.ToUint128(2997924580000000000), "USDT")

	s.db.Create(&models.AllFunds)

	return nil
}
