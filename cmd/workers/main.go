package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"github.com/smallbatch-apps/earnsmart-api/database"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

func main() {
	godotenv.Load()
	c := cron.New()
	database.Connect()
	tbClient, _ := database.CreateTigerBeetleClient()
	priceService := services.NewPriceService(database.DB, tbClient)

	c.AddFunc("10 * * * * *", func() {
		if err := priceService.UpdatePrices(); err != nil {
			log.Println("Error updating prices:", err)
		}
	})

	c.Start()

	select {}
}
