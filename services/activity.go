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

func (s *ActivityService) ListActivities(userID uint64) ([]models.Activity, error) {
	var activities []models.Activity
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&activities).Error
	return activities, err
}

func (s *ActivityService) CreateActivity(activity models.Activity) error {
	return s.db.Create(&activity).Error
}

func (s *ActivityService) CreateUserActivity(userID uint64, message string) error {
	activity := models.Activity{
		Type:         models.ActivityTypeUser,
		Message:      message,
		OwnableModel: models.OwnableModel{UserID: userID},
	}
	return s.CreateActivity(activity)
}

func (s *ActivityService) CreateAdminActivity(message string) error {
	activity := models.Activity{
		Type:    models.ActivityTypeAdmin,
		Message: message,
	}
	return s.CreateActivity(activity)
}
