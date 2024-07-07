package services

import (
	"github.com/smallbatch-apps/earnsmart-api/models"

	"gorm.io/gorm"
)

type UserService struct {
	*BaseService
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		BaseService: NewBaseService(db),
	}
}

func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) FindUserById(db *gorm.DB, id string) error {
	user := models.User{}
	return s.db.Where("id = ?", id).First(&user).Error
}
