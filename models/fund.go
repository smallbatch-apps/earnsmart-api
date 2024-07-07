package models

import (
	"encoding/json"
)

type FundPeriod int

const (
	FundPeriodImmediate FundPeriod = iota
	FundPeriodMonth
	FundPeriod3Months
	FundPeriod6Months
	FundPeriodYear
)

type Fund struct {
	CustomModel
	Name     string
	Currency string
	Period   uint
	Rate     float32
}

func (fund Fund) MarshalJSON() ([]byte, error) {
	type Alias Fund

	return json.Marshal(&struct {
		ID       uint    `json:"id"`
		Name     string  `json:"name"`
		Currency string  `json:"currency"`
		Period   uint    `json:"period"`
		Rate     float32 `json:"rate"`
		Alias
	}{
		ID:       fund.ID,
		Name:     fund.Name,
		Currency: fund.Currency,
		Period:   fund.Period,
		Rate:     fund.Rate,
		Alias:    (Alias)(fund),
	})
}
