package services

import (
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"gorm.io/gorm"
)

type AdminService struct {
	*BaseService
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		BaseService: NewBaseService(db),
	}
}

func (s *AdminService) SeedData() error {

	s.db.AutoMigrate(&models.Price{}, &models.Fund{}, &models.Setting{}, &models.User{})

	price_service := NewPriceService(s.db)

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

	funds := []models.Fund{
		{Name: "Bitcoin Immediate Return", Currency: "BTC", Period: uint(models.FundPeriodImmediate), Rate: 0.04},
		{Name: "Bitcoin Short Term", Currency: "BTC", Period: uint(models.FundPeriodMonth), Rate: 0.06},
		{Name: "Bitcoin Medium - 3 months", Currency: "BTC", Period: uint(models.FundPeriod3Months), Rate: 0.07},
		{Name: "Bitcoin Medium - 6 months", Currency: "BTC", Period: uint(models.FundPeriod6Months), Rate: 0.08},
		{Name: "Bitcoin ", Currency: "BTC", Period: uint(models.FundPeriodYear), Rate: 0.095},

		{Currency: "ETH", Period: uint(models.FundPeriodImmediate), Rate: 0.04},
		{Currency: "ETH", Period: uint(models.FundPeriodMonth), Rate: 0.06},
		{Currency: "ETH", Period: uint(models.FundPeriod3Months), Rate: 0.07},
		{Currency: "ETH", Period: uint(models.FundPeriod6Months), Rate: 0.08},
		{Currency: "ETH", Period: uint(models.FundPeriod6Months), Rate: 0.095},
	}

	return nil
}
