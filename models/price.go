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
	Rate      float32 `gorm:"not null;index" json:"rate"`
}

// func (price Price) MarshalJSON() ([]byte, error) {
// 	type Alias Price

// 	return json.Marshal(&struct {
// 		ID        uint    `json:"id"`
// 		Currency  string  `json:"currency"`
// 		IsCurrent bool    `json:"is_current"`
// 		Period    uint    `json:"period"`
// 		Rate      float32 `json:"rate"`
// 		Alias
// 	}{
// 		ID:        price.ID,
// 		Currency:  price.Currency,
// 		IsCurrent: price.IsCurrent,
// 		Period:    price.Period,
// 		Rate:      price.Rate,
// 		Alias:     (Alias)(price),
// 	})
// }
