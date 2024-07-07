package models

import "errors"

type CurrencyOption int

const (
	ETH CurrencyOption = iota
	BTC
	USDT
	USDC
	DAI
	BNB
	MATIC
	AVAX
	SOL
	BAT
	LINK
	UNI
	XRP
	ADA
	HBAR
	DOT
	TRX
)

var currencies = [...]string{
	"ETH", "BTC", "USDT", "USDC", "DAI", "BNB", "MATIC", "AVAX",
	"SOL", "BAT", "LINK", "UNI", "XRP", "ADA", "HBAR", "DOT", "TRX",
}

func (c CurrencyOption) String() string {
	if c < ETH || c > TRX {
		return "Unknown"
	}
	return currencies[c]
}

func CurrencyFromInt(value int) (CurrencyOption, error) {
	if value < int(ETH) || value > int(TRX) {
		return 0, errors.New("invalid currency code")
	}
	return CurrencyOption(value), nil
}
