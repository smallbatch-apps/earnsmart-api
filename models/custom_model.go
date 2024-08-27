package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomModel struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type CustomMarshalType struct {
	ID        uint   `json:"id" `
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
