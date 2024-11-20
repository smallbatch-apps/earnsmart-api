package services

import (
	"fmt"
	"log"
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

func (s *AdminService) SafeToMigrate() bool {
	var users models.User
	err := s.db.Find(&users).Error
	return err != nil
}

func (s *AdminService) SeedData() error {

	s.db.AutoMigrate(&models.Account{}, &models.AllocationPlan{}, &models.Activity{}, &models.Fund{}, &models.Price{}, &models.Setting{}, &models.Swap{}, &models.Transaction{}, &models.User{})

	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)
	transferService := NewTransferService(s.db, s.tbClient)

	log.Println("Updating prices")

	priceService.UpdatePrices(models.CurrencyPeriod1m)

	log.Println("Creating test users")

	// create admin user
	adminUser := models.User{
		Name:     "Michael Smith",
		Email:    "m.v.smith@earnsmart.com",
		Password: os.Getenv("ADMIN_USER_PASSWORD"),
	}
	err := s.db.Create(&adminUser).Error
	if err != nil {
		log.Println("Error creating admin user: ", err.Error())
		return err
	}

	log.Println("Creating admin wallet accounts")

	for currency := range models.AllCurrencies {
		_, err = accountService.CreateWalletAccount(adminUser.ID, currency)
		if err != nil {
			log.Println(fmt.Sprintf("Error creating %s wallet account", currency), err.Error())
			return err
		}
	}

	// _, err = accountService.CreateWalletAccount(adminUser.ID, "ETH")

	// create demo user
	testUser := models.User{
		Name:     "Test User",
		Email:    "test@earnsmart.com",
		Password: "admin123456",
	}

	err = s.db.Create(&testUser).Error
	if err != nil {
		log.Println("Error creating test user: ", err.Error())
		return err
	}

	log.Println("Creating test accounts and deposits")

	// var testAccount models.Account
	testAccount, err := accountService.CreateWalletAccount(testUser.ID, "ETH")
	if err != nil {
		log.Println("Error creating eth wallet account: ", err.Error())
		return err
	}

	if _, err = transferService.CreateDepositTransfer(testAccount, tbt.ToUint128(3290915267000000000), "ETH"); err != nil {
		log.Println("Error creating eth deposit: ", err.Error())
		return err
	}

	if testAccount, err = accountService.CreateWalletAccount(testUser.ID, "USDT"); err != nil {
		log.Println("Error creating usdt wallet account: ", err.Error())
		return err
	}

	if _, err = transferService.CreateDepositTransfer(testAccount, tbt.ToUint128(1200000000), "USDT"); err != nil {
		log.Println("Error creating usdt deposit: ", err.Error())
		return err
	}
	log.Println("Populating all funds")
	s.db.Create(&models.AllFunds)

	return nil
}
