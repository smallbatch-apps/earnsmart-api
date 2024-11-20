package models

type LocalCurrency struct {
	Name     string
	Symbol   string
	Decimals uint32
	LedgerID uint32
}

var AllCurrencies = map[string]LocalCurrency{
	"ADA":   {Name: "Cardano", Symbol: "ADA", Decimals: 18, LedgerID: 1},
	"AVAX":  {Name: "AVAX", Symbol: "AVAX", Decimals: 18, LedgerID: 2},
	"BAT":   {Name: "Basic Attention Token", Symbol: "BAT", Decimals: 18, LedgerID: 3},
	"BNB":   {Name: "BNB", Symbol: "BNB", Decimals: 18, LedgerID: 4},
	"BTC":   {Name: "Bitcoin", Symbol: "BTC", Decimals: 8, LedgerID: 5},
	"DAI":   {Name: "DAI", Symbol: "DAI", Decimals: 18, LedgerID: 6},
	"DOT":   {Name: "Polkadot", Symbol: "DOT", Decimals: 10, LedgerID: 7},
	"ETH":   {Name: "Ether", Symbol: "ETH", Decimals: 18, LedgerID: 8},
	"HBAR":  {Name: "Hedera HBAR", Symbol: "HBAR", Decimals: 8, LedgerID: 9},
	"LINK":  {Name: "LINK", Symbol: "LINK", Decimals: 24, LedgerID: 10},
	"MATIC": {Name: "MATIC", Symbol: "MATIC", Decimals: 18, LedgerID: 11},
	"SOL":   {Name: "Solana", Symbol: "SOL", Decimals: 9, LedgerID: 12},
	"TRX":   {Name: "TRX", Symbol: "TRX", Decimals: 6, LedgerID: 13},
	"UNI":   {Name: "Uniswap", Symbol: "UNI", Decimals: 18, LedgerID: 14},
	"USDT":  {Name: "USDT", Symbol: "USDT", Decimals: 6, LedgerID: 15},
	"USDC":  {Name: "USDC", Symbol: "USDC", Decimals: 6, LedgerID: 16},
	"XRP":   {Name: "Ripple", Symbol: "XRP", Decimals: 6, LedgerID: 17},
}

var LedgerCurrency = map[uint32]string{
	1:  "ADA",
	2:  "AVAX",
	3:  "BAT",
	4:  "BNB",
	5:  "BTC",
	6:  "DAI",
	7:  "DOT",
	8:  "ETH",
	9:  "HBAR",
	10: "LINK",
	11: "MATIC",
	12: "SOL",
	13: "TRX",
	14: "UNI",
	15: "USDT",
	16: "USDC",
	17: "XRP",
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
// 		ID       uint64   `json:"id"`
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
