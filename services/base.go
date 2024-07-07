package services

import (
	"gorm.io/gorm"
)

// BaseService provides common CRUD operations
type BaseService struct {
	db *gorm.DB
}

func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{db: db}
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
