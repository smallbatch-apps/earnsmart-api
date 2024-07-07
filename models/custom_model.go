package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomModel struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CustomMarshalType struct {
	ID        uint   `json:"id" `
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
