package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/smallbatch-apps/earnsmart-api/models"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func Debug(message string) {
	fmt.Println("DEBUG:", message)
}

func FormatString(str string) string {
	return fmt.Sprintf("Formatted: %s", str)
}

func LogJson(message string, data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	log.Println(message, string(jsonData))
}

func LogAccount(account tbt.Account) {
	LogJson("Account:", models.AccountWrapper{Account: account})
}

func LogTransfer(transfer tbt.Transfer) {
	LogJson("Transfer:", models.TransferWrapper{Transfer: transfer})
}

func LogAccountFilter(accountFilter tbt.AccountFilter) {
	LogJson("AccountFilter:", models.AccountFilterWrapper{AccountFilter: accountFilter})
}

func LogAccountBalance(accountBalance tbt.AccountBalance) {
	LogJson("AccountBalance:", models.AccountBalanceWrapper{AccountBalance: accountBalance})
}

func LogAccountIDs(ids []tbt.Uint128) {
	var results []string
	for _, id := range ids {
		uint64Value := FromUint128(id)
		results = append(results, fmt.Sprintf("%d", uint64Value))
	}
	log.Println("Account IDs:", strings.Join(results, ","))
}

func ToUint128(id uint) tbt.Uint128 {
	return tbt.ToUint128(uint64(id))
}

func FromUint128(u tbt.Uint128) uint64 {
	bigIntValue := u.BigInt()
	return bigIntValue.Uint64()
}
