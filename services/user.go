package services

import (
	"fmt"

	"github.com/smallbatch-apps/earnsmart-api/models"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

type UserService struct {
	*BaseService
}

func NewUserService(db *gorm.DB, tbClient tb.Client) *UserService {
	return &UserService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *UserService) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (s *UserService) GetUser(id uint64) (models.User, error) {
	user := models.User{}
	err := s.db.First(&user, id).Error
	return user, err
}

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	err := s.db.Create(user).Error
	if err != nil {
		return &models.User{}, err
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Created user: %s", user.Name), user.ID)
	return user, nil
}
