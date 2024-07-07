package models

type LocalCurrency struct {
	Name     string
	Symbol   string
	Decimals uint
	LedgerID uint
}

var AllCurrencies = map[string]LocalCurrency{
	"ADA":   {Name: "Cardano", Symbol: "ADA", Decimals: 18, LedgerID: 1},
	"AVAX":  {Name: "AVAX", Symbol: "AVAX", Decimals: 18, LedgerID: 2},
	"BAT":   {Name: "Basic Attention Token", Symbol: "BAT", Decimals: 18, LedgerID: 3},
	"BNB":   {Name: "BNB", Symbol: "BNB", Decimals: 18, LedgerID: 4},
	"BTC":   {Name: "Bitcoin", Symbol: "BTC", Decimals: 18, LedgerID: 5},
	"DAI":   {Name: "DAI", Symbol: "DAI", Decimals: 18, LedgerID: 6},
	"DOT":   {Name: "Polkadot", Symbol: "DOT", Decimals: 18, LedgerID: 7},
	"ETH":   {Name: "Ether", Symbol: "ETH", Decimals: 18, LedgerID: 8},
	"HBAR":  {Name: "Hedera HBAR", Symbol: "HBAR", Decimals: 18, LedgerID: 9},
	"LINK":  {Name: "LINK", Symbol: "LINK", Decimals: 24, LedgerID: 10},
	"MATIC": {Name: "MATIC", Symbol: "MATIC", Decimals: 18, LedgerID: 11},
	"SOL":   {Name: "Solana", Symbol: "SOL", Decimals: 18, LedgerID: 12},
	"TRX":   {Name: "TRX", Symbol: "TRX", Decimals: 18, LedgerID: 13},
	"UNI":   {Name: "Uniswap", Symbol: "UNI", Decimals: 18, LedgerID: 14},
	"USDT":  {Name: "USDT", Symbol: "USDT", Decimals: 18, LedgerID: 15},
	"USDC":  {Name: "USDC", Symbol: "USDC", Decimals: 18, LedgerID: 16},
	"XRP":   {Name: "Ripple", Symbol: "XRP", Decimals: 18, LedgerID: 17},
}

// type Currency struct {
// 	CustomModel
// 	Name     string
// 	Symbol   string
// 	Decimals uint
// }

// func (currency Currency) MarshalJSON() ([]byte, error) {
// 	type Alias Currency

// 	return json.Marshal(&struct {
// 		ID       uint   `json:"id"`
// 		Name     string `json:"name"`
// 		Symbol   string `json:"symbol"`
// 		Decimals uint   `json:"decimals"`
// 		Alias
// 	}{
// 		ID:       currency.ID,
// 		Name:     currency.Name,
// 		Symbol:   currency.Symbol,
// 		Decimals: currency.Decimals,
// 		Alias:    (Alias)(currency),
// 	})
// }
