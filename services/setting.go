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

func (s *SettingService) ListSettings(userID uint64) ([]models.Setting, error) {
	var settings []models.Setting
	if err := s.db.Where("user_id = ?", userID).Or("type = ?", "app").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) UpdateSetting(userID uint64, setting string, value string) (models.Setting, error) {
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("update setting %s to %s", setting, value), userID)

	dbSetting := models.Setting{
		OwnableModel: models.OwnableModel{UserID: userID},
		Name:         setting,
		Value:        value,
		Type:         models.SettingTypeUser,
	}

	result := s.db.Where("user_id = ? AND name = ?", userID, setting).
		Assign(models.Setting{Value: value}).
		FirstOrCreate(&dbSetting)

	if result.Error != nil {
		return models.Setting{}, result.Error
	}

	return dbSetting, nil
}

func (s *SettingService) GetSetting(userID uint64, setting string) (models.Setting, error) {
	var dbSetting models.Setting
	err := s.db.Where("user_id = ? AND name = ?", userID, setting).First(&dbSetting).Error
	return dbSetting, err
}
