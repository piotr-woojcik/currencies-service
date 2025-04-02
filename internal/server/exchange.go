package server

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CryptoExchangeRate struct {
	DecimalPlaces int     `json:"decimal_places"`
	USDRate       float64 `json:"exchange_rate"`
}

var cryptoExchangeRates = map[string]CryptoExchangeRate{
	"BEER": {
		DecimalPlaces: 18,
		USDRate:       0.00002461,
	},
	"FLOKI": {
		DecimalPlaces: 18,
		USDRate:       0.0001428,
	},
	"GATE": {
		DecimalPlaces: 18,
		USDRate:       6.87,
	},
	"USDT": {
		DecimalPlaces: 6,
		USDRate:       0.999,
	},
	"WBTC": {
		DecimalPlaces: 8,
		USDRate:       57037.22,
	},
}

func (s *Server) getCryptoExchange(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")
	if from == "" || to == "" || amount == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	fromRate, ok := cryptoExchangeRates[from]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	toRate, ok := cryptoExchangeRates[to]
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	rate := fromRate.USDRate / toRate.USDRate
	convertedAmount := amountFloat * rate

	res := map[string]any{
		"from":   from,
		"to":     to,
		"amount": roundToDecimalPlaces(convertedAmount, toRate.DecimalPlaces),
	}

	c.JSON(http.StatusOK, res)
}

func roundToDecimalPlaces(value float64, decimalPlaces int) float64 {
	pow := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*pow) / pow
}
