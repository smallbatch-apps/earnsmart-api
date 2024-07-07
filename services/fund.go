package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"

	"gorm.io/gorm"
)

type FundService struct {
	*BaseService
}

func NewFundService(db *gorm.DB) *FundService {
	return &FundService{
		BaseService: NewBaseService(db),
	}
}

func (s *FundService) GetFund(id string) (*models.Fund, error) {
	fund := models.Fund{}
	if err := s.db.Where("id = ?", id).First(&fund).Error; err != nil {
		return nil, err
	}
	return &fund, nil
}

func (s *FundService) ListFunds() ([]models.Fund, error) {
	var funds []models.Fund
	if err := s.db.Order("currency, period DESC").Find(&funds).Error; err != nil {
		return nil, err
	}
	return funds, nil
}
