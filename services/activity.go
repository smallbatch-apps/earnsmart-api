package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

type ActivityService struct {
	*BaseService
}

func NewActivityService(db *gorm.DB, tbClient tb.Client) *ActivityService {
	return &ActivityService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *ActivityService) GetAll(userID uint) ([]models.Activity, error) {
	var activities []models.Activity
	if err := s.db.Where("user_id = ?", userID).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}
