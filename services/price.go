package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"

	"gorm.io/gorm"
)

type PriceService struct {
	*BaseService
}

func NewPriceService(db *gorm.DB) *PriceService {
	return &PriceService{
		BaseService: NewBaseService(db),
	}
}

func (s *PriceService) GetPrices() ([]models.Price, error) {
	var prices []models.Price
	err := s.db.Where("is_current = ?", true).Order("currency DESC").Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func (s *PriceService) GetPricesForPeriod(currency string, period uint) ([]models.Price, error) {
	var prices []models.Price
	if err := s.db.Where("currency = ? AND period = ?", currency, period).Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}

func (s *PriceService) RequestForQuote(fromCurrency string, toCurrency string, amount uint) {
}

func (s *PriceService) UpdatePrices() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("CMC_HOST")+"/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
	}

	q := url.Values{}
	q.Add("symbol", "ETH,BTC,USDT,USDC,DAI,BNB,MATIC,AVAX,SOL,BAT,LINK,UNI,XRP,ADA,HBAR,DOT,TRX")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", os.Getenv("COINMARKETCAP_API_KEY"))
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		os.Exit(1)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal([]byte(respBody), &data)

	var prices []models.Price

	for _, value := range data {
		currencyData := value.(map[string]interface{})
		symbol := currencyData["symbol"].(string)
		quoteData := currencyData["quote"].(map[string]interface{})
		usdData := quoteData["USD"].(map[string]interface{})
		price := usdData["price"].(float64)

		prices = append(prices, models.Price{
			Currency:  symbol,
			IsCurrent: true,
			Period:    uint(models.CurrencyPeriod1h),
			Rate:      float32(price),
		})
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Model(&models.Price{}).Where("is_current = ?", true).Update("is_current", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&prices).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// InsertPricesBatch inserts multiple price records in one batch.
func (s *PriceService) InsertPricesBatch(prices []models.Price) error {
	return s.db.Create(&prices).Error
}
