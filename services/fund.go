package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

type FundService struct {
	*BaseService
}

func NewFundService(db *gorm.DB, tbClient tb.Client) *FundService {
	return &FundService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *FundService) GetFund(id uint) (models.Fund, error) {
	fund := models.Fund{}
	if err := s.db.First(&fund, id).Error; err != nil {
		return fund, err
	}
	return fund, nil
}

func (s *FundService) ListFunds() ([]models.Fund, error) {
	var funds []models.Fund
	if err := s.db.Order("currency, period DESC").Find(&funds).Error; err != nil {
		return nil, err
	}
	return funds, nil
}
