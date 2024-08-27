package models

import (
	"encoding/json"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	CustomModel
	Name     string `json:"name"`
	Email    string `gorm:"index:idx_email_unique,unique" json:"email"`
	Password string `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Password, err = HashPassword(u.Password)
	return
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (user User) MarshalJSON() ([]byte, error) {
	type Alias User

	return json.Marshal(&struct {
		ID        uint   `json:"id" `
		Name      string `json:"name"`
		Email     string `json:"email"`
		Password  string `json:"-"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Alias
	}{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Alias:     (Alias)(user),
	})
}
