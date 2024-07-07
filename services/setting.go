package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"
	"gorm.io/gorm"
)

// SettingService extends BaseService for settings
type SettingService struct {
	*BaseService
}

// NewSettingService creates a new SettingService
func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{
		BaseService: NewBaseService(db),
	}
}

func (s *SettingService) GetAll(userID uint) ([]models.Setting, error) {
	var settings []models.Setting
	if err := s.db.Where("user_id = ?", userID).Or("type = ?", "app").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *SettingService) SetSetting(setting *models.Setting) error {
	return s.db.Save(setting).Error
}
