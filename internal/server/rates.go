package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) getRates(c *gin.Context) {

	currencies := c.Query("currencies")
	if currencies == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	currList := strings.Split(currencies, ",")
	if len(currList) <= 1 {
		c.Status(http.StatusBadRequest)
		return
	}

	usdRate := 1.0
	rates, err := s.exchangeClient.GetLatestRatesForUSD()
	if err != nil {
		slog.Error("failed to get latest rates", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}

	currRates := make(map[string]float64)

	for _, currency := range currList {
		if currency == "USD" {
			currRates[currency] = usdRate
			continue
		}

		rate, ok := rates[currency]
		if !ok {
			slog.Warn("currency not found", "currency", currency)
		} else {
			currRates[currency] = rate
		}
	}

	res := make([]map[string]any, 0, len(currList)*(len(currList)-1))

	for i := 0; i < len(currList); i++ {
		for j := 0; j < len(currList); j++ {
			if i != j {
				res = append(res, map[string]any{
					"from": currList[i],
					"to":   currList[j],
					"rate": currRates[currList[j]] / currRates[currList[i]],
				})
			}
		}
	}

	c.JSON(http.StatusOK, res)
}
