package services

import (
	"gorm.io/gorm"
)

type TransactionService struct {
	*BaseService
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		BaseService: NewBaseService(db),
	}
}
