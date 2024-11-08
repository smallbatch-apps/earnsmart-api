package models

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
)

type AllocationPlanStatus uint16

const (
	AllocationPlanStatusActive AllocationPlanStatus = iota
	AllocationPlanStatusPaused
	AllocationPlanStatusCompleted
	AllocationPlanStatusCancelled
)

type AllocationPlan struct {
	CustomModel
	OwnableModel
	FromCurrency string               `json:"from_currency"`
	ToCurrency   string               `json:"to_currency"`
	Amount       decimal.Decimal      `json:"amount"`
	Frequency    string               `json:"frequency"`
	Status       AllocationPlanStatus `json:"status"`
	NextRun      time.Time            `json:"next_run"`
}

func (ap *AllocationPlan) GetNextRun() (time.Time, error) {
	schedule, err := cron.ParseStandard(ap.Frequency)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid frequency format: %w", err)
	}
	return schedule.Next(time.Now()), nil
}

// ValidateFrequency ensures the frequency string is a valid cron expression
func (ap *AllocationPlan) ValidateFrequency() error {
	_, err := cron.ParseStandard(ap.Frequency)
	return err
}
