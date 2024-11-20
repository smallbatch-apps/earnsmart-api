package models

type CurrencyPeriod int

const (
	CurrencyPeriod1m CurrencyPeriod = iota
	CurrencyPeriod1h
	CurrencyPeriodDay
)

type Price struct {
	CustomModel
	Currency  string  `gorm:"not null;index" json:"currency"`
	IsCurrent bool    `gorm:"not null;default:false;index" json:"is_current"`
	Period    uint    `gorm:"not null;index" json:"period"`
	Rate      float64 `gorm:"not null;index" json:"rate"`
	Change24h float64 `json:"change_24h"`
}
