package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"

	"github.com/smallbatch-apps/earnsmart-api/database"
	"github.com/smallbatch-apps/earnsmart-api/models"
	"github.com/smallbatch-apps/earnsmart-api/services"
)

func main() {
	godotenv.Load()
	c := cron.New()
	database.Connect()
	tbClient, _ := database.CreateTigerBeetleClient()
	priceService := services.NewPriceService(database.DB, tbClient)

	c.AddFunc("10 * * * * *", func() {
		if err := priceService.UpdatePrices(models.CurrencyPeriod1m); err != nil {
			log.Println("Error updating prices:", err)
		}
	})

	c.AddFunc("50 20 * * * *", func() {
		if err := priceService.UpdatePrices(models.CurrencyPeriod1h); err != nil {
			log.Println("Error updating 1h prices:", err)
		}
	})

	// Every day at 03:15:25
	c.AddFunc("25 15 3 * * *", func() {
		if err := priceService.UpdatePrices(models.CurrencyPeriodDay); err != nil {
			log.Println("Error updating 1d prices:", err)
		}
	})

	c.Start()

	select {}
}
