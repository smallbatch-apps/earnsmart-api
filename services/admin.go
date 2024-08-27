package services

import (
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

	s.db.AutoMigrate(&models.Account{}, &models.Fund{}, &models.Price{}, &models.Setting{}, &models.User{}, &models.Activity{})

	priceService := NewPriceService(s.db, s.tbClient)
	accountService := NewAccountService(s.db, s.tbClient)
	transactionService := NewTransactionService(s.db, s.tbClient)

	log.Println("Updating prices")

	priceService.UpdatePrices()

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
			log.Println("Error creating eth wallet account: ", err.Error())
			return err
		}
	}

	// _, err = accountService.CreateWalletAccount(adminUser.ID, "ETH")

	// create demo user
	testUser := models.User{
		Name:     "Test User",
		Email:    "test@earnsmart.com",
		Password: "123456",
	}

	err = s.db.Create(&testUser).Error
	if err != nil {
		log.Println("Error creating test user: ", err.Error())
		return err
	}

	log.Println("Creating test accounts and deposits")

	var testAccount models.Account
	_, err = accountService.CreateWalletAccount(testUser.ID, "ETH")
	if err != nil {
		log.Println("Error creating eth wallet account: ", err.Error())
		return err
	}

	// if _, err = transactionService.CreateDeposit(testAccountId, tbt.ToUint128(2997924580000000000), "ETH"); err != nil {
	// 	log.Println("Error creating eth deposit: ", err.Error())
	// 	return err
	// }

	if testAccount, err = accountService.CreateWalletAccount(testUser.ID, "USDT"); err != nil {
		log.Println("Error creating usdt wallet account: ", err.Error())
		return err
	}

	if _, err = transactionService.CreateDeposit(testAccount, tbt.ToUint128(2997924580000000000), "USDT"); err != nil {
		log.Println("Error creating usdt deposit: ", err.Error())
		return err
	}
	// log.Println("Populating all funds")
	// s.db.Create(&models.AllFunds)

	return nil
}
