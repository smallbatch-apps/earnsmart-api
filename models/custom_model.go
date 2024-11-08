package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomModel struct {
	ID        uint64         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type CustomMarshalType struct {
	ID        uint64 `json:"id" `
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OwnableModel struct {
	UserID uint64 `json:"-" gorm:"index"`
}

func (om *OwnableModel) IsOwner(userID uint64) bool {
	return om.UserID == userID
}

func (om *OwnableModel) SetUserID(userID uint64) {
	om.UserID = userID
}
