package services

import (
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"

	"github.com/smallbatch-apps/earnsmart-api/models"

	tb "github.com/tigerbeetle/tigerbeetle-go"
	tbt "github.com/tigerbeetle/tigerbeetle-go/pkg/types"
	"gorm.io/gorm"
)

type PriceService struct {
	*BaseService
}

func NewPriceService(db *gorm.DB, tbClient tb.Client) *PriceService {
	return &PriceService{
		BaseService: NewBaseService(db, tbClient),
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

func (s *PriceService) GetPriceForCurrency(currency string) (models.Price, error) {
	var price models.Price
	query := models.Price{Currency: currency, IsCurrent: true}
	err := s.db.Where(query).Order("currency DESC").First(&price).Error
	if err != nil {
		return models.Price{}, err
	}
	return price, nil
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
	req, err := http.NewRequest("GET", os.Getenv("COINMARKETCAP_HOST")+"/v1/cryptocurrency/quotes/latest", nil)
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
		log.Println(err.Error())
		log.Println("Error in Coin Market Cap API key", os.Getenv("COINMARKETCAP_API_KEY"))
		return err
		//os.Exit(1)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal([]byte(respBody), &response)
	dataField, ok := response["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("Error: 'data' field is missing or not an array")
	}
	var prices []models.Price

	for _, value := range dataField {

		currencyData, ok := value.(map[string]interface{})
		if !ok {
			log.Println("Invalid currency data format")
			continue
		}

		currencyDataJson, _ := json.MarshalIndent(currencyData, "", "  ")

		symbol, ok := currencyData["symbol"].(string)
		if !ok {
			log.Println("Missing or invalid 'symbol' field", symbol)
			log.Printf("currencyData: %+v\n\n", string(currencyDataJson))
			continue
		}

		quoteData, ok := currencyData["quote"].(map[string]interface{})
		if !ok {
			log.Println("Missing or invalid 'quote' field")
			log.Printf("currencyData: %+v\n\n", string(currencyDataJson))
			continue
		}

		usdData, ok := quoteData["USD"].(map[string]interface{})
		if !ok {
			log.Println("Missing or invalid 'USD' field")
			log.Printf("currencyData: %+v\n\n", string(currencyDataJson))
			continue
		}

		price, ok := usdData["price"].(float64)
		if !ok {
			log.Println("Missing or invalid 'price' field")
			log.Printf("currencyData: %+v\n\n", string(currencyDataJson))
			continue
		}

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

func (s *PriceService) GetAmountPrice(currency string, amount tbt.Uint128) (float64, error) {
	var price models.Price
	if err := s.db.Where("currency = ?", currency).First(&price).Error; err != nil {
		return 0, err
	}

	decimals := uint64(models.AllCurrencies[currency].Decimals)

	decimalsBigInt := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountBigInt := amount.BigInt()
	shiftedAmount := new(big.Float).Quo(new(big.Float).SetInt(&amountBigInt), new(big.Float).SetInt(decimalsBigInt))

	shiftedAmountFloat, _ := shiftedAmount.Float64()

	amountUSD := float64(price.Rate) * shiftedAmountFloat

	return amountUSD, nil
}
