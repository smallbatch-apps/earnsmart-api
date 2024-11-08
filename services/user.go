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

func (s *UserService) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// func (s *UserService) FindUserById(db *gorm.DB, id string) error {
// 	user := models.User{}
// 	return s.db.Where("id = ?", id).First(&user).Error
// }

func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	err := s.db.Create(user).Error
	if err != nil {
		return &models.User{}, err
	}
	s.LogActivity(models.ActivityTypeAdmin, fmt.Sprintf("Created user: %s", user.Name), user.ID)
	return user, nil
}
