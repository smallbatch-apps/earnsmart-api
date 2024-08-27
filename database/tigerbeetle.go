package database

import (
	"log"
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
		log.Println("Failed to create TigerBeetle client! \n", err.Error())
		return nil, err
	} else {
		log.Println("Connected to TigerBeetle successfully")
	}
	return client, nil
}
