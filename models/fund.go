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
	Code     AccountCode
	Rate     float32
}

func (fund Fund) MarshalJSON() ([]byte, error) {
	type Alias Fund

	return json.Marshal(&struct {
		ID       uint64  `json:"id"`
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

var AllFunds = []Fund{
	{Name: "Bitcoin Immediate Return", Currency: "BTC", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Bitcoin Short Term", Currency: "BTC", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Bitcoin Medium - 3 months", Currency: "BTC", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Bitcoin Medium - 6 months", Currency: "BTC", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Bitcoin Long Term", Currency: "BTC", Period: uint(FundPeriodYear), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Ether Immediate", Currency: "ETH", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Ether Short", Currency: "ETH", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Ether Medium - 3 months", Currency: "ETH", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Ether Medium - 6 months", Currency: "ETH", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Ether Long Term", Currency: "ETH", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Tether Immediate", Currency: "USDT", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Tether Short", Currency: "USDT", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Tether Medium - 3 months", Currency: "USDT", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Tether Medium - 6 months", Currency: "USDT", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Tether Long Term", Currency: "USDT", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "USDC Immediate", Currency: "USDC", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "USDC Short", Currency: "USDC", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "USDC Medium - 3 months", Currency: "USDC", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "USDC Medium - 6 months", Currency: "USDC", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "USDC Long Term", Currency: "USDC", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "BAT Immediate", Currency: "BAT", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "BAT Short", Currency: "BAT", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "BAT Medium - 3 months", Currency: "BAT", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "BAT Medium - 6 months", Currency: "BAT", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "BAT Long Term", Currency: "BAT", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "ADA Immediate", Currency: "ADA", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "ADA Short", Currency: "ADA", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "ADA Medium - 3 months", Currency: "ADA", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "ADA Medium - 6 months", Currency: "ADA", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "ADA Long Term", Currency: "ADA", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "AVAX Immediate", Currency: "AVAX", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "AVAX Short", Currency: "AVAX", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "AVAX Medium - 3 months", Currency: "AVAX", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "AVAX Medium - 6 months", Currency: "AVAX", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "AVAX Long Term", Currency: "AVAX", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "DOT Immediate", Currency: "DOT", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "DOT Short", Currency: "DOT", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "DOT Medium - 3 months", Currency: "DOT", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "DOT Medium - 6 months", Currency: "DOT", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "DOT Long Term", Currency: "DOT", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "BNB Immediate", Currency: "BNB", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "BNB Short", Currency: "BNB", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "BNB Medium - 3 months", Currency: "BNB", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "BNB Medium - 6 months", Currency: "BNB", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "BNB Long Term", Currency: "BNB", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Dai Immediate", Currency: "DAI", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Dai Short", Currency: "DAI", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Dai Medium - 3 months", Currency: "DAI", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Dai Medium - 6 months", Currency: "DAI", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Dai Long Term", Currency: "DAI", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Matic Immediate", Currency: "AVAX", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Matic Short", Currency: "AVAX", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Matic Medium - 3 months", Currency: "AVAX", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Matic Medium - 6 months", Currency: "AVAX", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Matic Long Term", Currency: "AVAX", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Solana Immediate", Currency: "SOL", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Solana Short", Currency: "SOL", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Solana Medium - 3 months", Currency: "SOL", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Solana Medium - 6 months", Currency: "SOL", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Solana Long Term", Currency: "SOL", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "HBAR Immediate", Currency: "HBAR", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "HBAR Short", Currency: "HBAR", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "HBAR Medium - 3 months", Currency: "HBAR", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "HBAR Medium - 6 months", Currency: "HBAR", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "HBAR Long Term", Currency: "HBAR", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "BNB Immediate", Currency: "BNB", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "BNB Short", Currency: "BNB", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "BNB Medium - 3 months", Currency: "BNB", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "BNB Medium - 6 months", Currency: "BNB", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "BNB Long Term", Currency: "BNB", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Ripple Immediate", Currency: "XRP", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Ripple Short", Currency: "XRP", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Ripple Medium - 3 months", Currency: "XRP", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Ripple Medium - 6 months", Currency: "XRP", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Ripple Long Term", Currency: "XRP", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Uniswap Immediate", Currency: "UNI", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Uniswap Short", Currency: "UNI", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Uniswap Medium - 3 months", Currency: "UNI", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Uniswap Medium - 6 months", Currency: "UNI", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Uniswap Long Term", Currency: "UNI", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Link Immediate", Currency: "LINK", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Link Short", Currency: "LINK", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Link Medium - 3 months", Currency: "LINK", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Link Medium - 6 months", Currency: "LINK", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Link Long Term", Currency: "LINK", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},

	{Name: "Tron Immediate", Currency: "TRX", Period: uint(FundPeriodImmediate), Code: AccountCodeFundImmediate, Rate: 0.04},
	{Name: "Tron Short", Currency: "TRX", Period: uint(FundPeriodMonth), Code: AccountCodeFundMonth, Rate: 0.06},
	{Name: "Tron Medium - 3 months", Currency: "TRX", Period: uint(FundPeriod3Months), Code: AccountCodeFund3Months, Rate: 0.07},
	{Name: "Tron Medium - 6 months", Currency: "TRX", Period: uint(FundPeriod6Months), Code: AccountCodeFund6Months, Rate: 0.08},
	{Name: "Tron Long Term", Currency: "TRX", Period: uint(FundPeriod6Months), Code: AccountCodeFundOneYear, Rate: 0.095},
}
