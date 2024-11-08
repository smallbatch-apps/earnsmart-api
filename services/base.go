package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

// BaseService provides common CRUD operations
type BaseService struct {
	db       *gorm.DB
	tbClient tb.Client
}

func NewBaseService(db *gorm.DB, tbClient tb.Client) *BaseService {
	return &BaseService{db: db, tbClient: tbClient}
}

func (s *BaseService) GetDB() *gorm.DB {
	return s.db
}

func (s *BaseService) GetTBClient() tb.Client {
	return s.tbClient
}

func (s *BaseService) LogActivity(activityType models.ActivityType, message string, userID uint64) (models.Activity, error) {
	activity := models.Activity{
		Type:         activityType,
		Message:      message,
		OwnableModel: models.OwnableModel{UserID: userID},
	}

	err := s.db.Create(&activity).Error
	if err != nil {
		return activity, err
	}

	return activity, nil
}
