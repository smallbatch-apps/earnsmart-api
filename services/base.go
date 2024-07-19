package services

import (
	"gorm.io/gorm"

	tb "github.com/tigerbeetle/tigerbeetle-go"
)

// BaseService provides common CRUD operations
type BaseService struct {
	db       *gorm.DB
	tbClient tb.Client
}

func NewBaseService(db *gorm.DB, tbClient tb.Client) *BaseService {
	return &BaseService{db: db, tbClient: tbClient}
}

func (s *BaseService) GetDB() *gorm.DB {
	return s.db
}

func (s *BaseService) GetTBClient() tb.Client {
	return s.tbClient
}

// // Get retrieves an entity by ID
// func (s *BaseService[T]) Get(id uint) (*T, error) {
// 	var entity T
// 	if err := s.db.First(&entity, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &entity, nil
// }

// // GetAll retrieves all entities
// func (s *BaseService[T]) GetAll() ([]T, error) {
// 	var entities []T
// 	if err := s.db.Find(&entities).Error; err != nil {
// 		return nil, err
// 	}
// 	return entities, nil
// }

// // Create creates a new entity
// func (s *BaseService[T]) Create(entity *T) error {
// 	return s.db.Create(entity).Error
// }

// // Update updates an existing entity
// func (s *BaseService[T]) Update(entity *T) error {
// 	return s.db.Save(entity).Error
// }

// // Delete deletes an entity by ID
// func (s *BaseService[T]) Delete(id uint) error {
// 	return s.db.Delete(new(T), id).Error
// }
