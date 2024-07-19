package models

type TransactionType uint

const (
	TransactionTypeDeposit TransactionType = iota
	TransactionTypeWithdraw
	TransactionTypeSwap
	TransactionTypeDeploy
	TransactionTypeInterest
	TransactionTypeRedeem
)
