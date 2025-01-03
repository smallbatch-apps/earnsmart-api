package services

import (
	"fmt"

	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/utils"
	tb "github.com/tigerbeetle/tigerbeetle-go"
	"gorm.io/gorm"
)

type AllocationService struct {
	*BaseService
}

func NewAllocationService(db *gorm.DB, tbClient tb.Client) *AllocationService {
	return &AllocationService{
		BaseService: NewBaseService(db, tbClient),
	}
}

func (s *AllocationService) GetAllocations(userID uint64) ([]models.AllocationPlan, error) {
	var allocations = []models.AllocationPlan{}
	s.db.Where("user_id = ?", userID).Find(&allocations)
	return allocations, nil
}

func (s *AllocationService) GetAllocation(id uint64) (models.AllocationPlan, error) {
	var allocation = models.AllocationPlan{}
	err := s.db.First(&allocation, id).Error
	return allocation, err
}

func (s *AllocationService) CreateAllocation(allocation *models.AllocationPlan) (*models.AllocationPlan, error) {
	err := s.db.Create(allocation).Error
	formattedAmount := utils.FormatCurrencyAmount(allocation.Amount, models.AllCurrencies[allocation.FromCurrency].Decimals)
	s.LogActivity(models.ActivityTypeUser, fmt.Sprintf("Create new asset allocation: %s %s in %s", formattedAmount, allocation.FromCurrency, allocation.ToCurrency), allocation.UserID)
	return allocation, err
}

func (s *AllocationService) UpdateAllocation(allocation *models.AllocationPlan, payload models.AllocationPlan) (*models.AllocationPlan, error) {
	err := s.db.Model(allocation).Updates(payload).Error
	return allocation, err
}

func (s *AllocationService) DeleteAllocation(allocation models.AllocationPlan) error {
	return s.db.Delete(allocation).Error
}
