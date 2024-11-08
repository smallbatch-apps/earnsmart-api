package services

import (
	"fmt"

	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

// SettingService extends BaseService for settings
type SettingService struct {
	*BaseService
}

func NewSettingService(db *gorm.DB, tbClient tb.Client) *SettingService {
	return &SettingService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *SettingService) GetAll(userID uint64) ([]models.Setting, error) {
	var settings []models.Setting
	if err := s.db.Where("user_id = ?", userID).Or("type = ?", "app").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) SetSetting(userID uint64, setting string, value string) error {
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("update setting %s to %s", setting, value), userID)
	return s.db.Exec("INSERT INTO settings (user_id, setting, value) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE value = ?", setting, value, value).Error
}

func (s *SettingService) GetSetting(userID uint64, setting string) (models.Setting, error) {
	var dbSetting models.Setting
	err := s.db.Where("user_id", userID).Where("setting = ?", setting).First(dbSetting).Error

	if err != nil {
		return dbSetting, err
	}

	return dbSetting, nil
}
