package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ExchangeClient struct {
	client  *http.Client
	baseURL string
	appID   string
}

func NewExchangeClient() *ExchangeClient {

	appID := os.Getenv("OPENEXCHANGERATES_APP_ID")
	if appID == "" {
		panic("OPENEXCHANGERATES_APP_ID environment variable is not set")
	}

	return &ExchangeClient{
		client:  &http.Client{},
		baseURL: "https://openexchangerates.org/api/",
		appID:   appID,
	}
}

func (ec ExchangeClient) GetLatestRatesForUSD() (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ec.baseURL+"latest.json?base=USD", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+ec.appID)

	resp, err := ec.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get latest rates: %s", resp.Status)
	}

	var ratesResponse struct {
		Disclaimer string             `json:"disclaimer"`
		License    string             `json:"license"`
		Timestamp  int64              `json:"timestamp"`
		Base       string             `json:"base"`
		Rates      map[string]float64 `json:"rates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ratesResponse); err != nil {
		return nil, err
	}

	if ratesResponse.Base != "USD" {
		return nil, fmt.Errorf("unexpected base currency: %s", ratesResponse.Base)
	}

	slog.Info("latest rates fetched", "timestamp", time.Unix(ratesResponse.Timestamp, 0).Format(time.RFC3339), "base", ratesResponse.Base)

	return ratesResponse.Rates, nil
}
