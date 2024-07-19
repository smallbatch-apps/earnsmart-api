package database

import (
	"os"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

func CreateTigerBeetleClient() (tb.Client, error) {
	tbAddress := os.Getenv("TB_ADDRESS")
	if len(tbAddress) == 0 {
		tbAddress = "3000"
	}
	client, err := tb.NewClient(tbt.ToUint128(0), []string{tbAddress}, 256)
	if err != nil {
		return nil, err
	}
	return client, nil
}
