package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRatesClient struct {
	latestRates map[string]float64
	err         error
}

func (m MockRatesClient) GetLatestRatesForUSD() (map[string]float64, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.latestRates, nil
}

func TestGetRates(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	testClient := &MockRatesClient{
		latestRates: map[string]float64{
			"EUR": 1.2,
			"GBP": 1.5,
		},
	}

	s := NewServer(testClient)

	testCases := []struct {
		name           string
		query          string
		expectedStatus int
		expectedLen    int
		withApiError   bool
		validation     func(t *testing.T, res []map[string]any)
	}{
		{
			name:           "Valid Request",
			query:          "currencies=USD,EUR,GBP",
			expectedStatus: http.StatusOK,
			expectedLen:    6,
			validation: func(t *testing.T, res []map[string]any) {
				for _, rate := range res {
					if rate["from"] == "USD" && rate["to"] == "USD" {
						t.Errorf("Unexpected rate from USD to USD")
					}
					if rate["from"] == "EUR" && rate["to"] == "EUR" {
						t.Errorf("Unexpected rate from EUR to EUR")
					}
					if rate["from"] == "GBP" && rate["to"] == "GBP" {
						t.Errorf("Unexpected rate from GBP to GBP")
					}
					if rate["from"] == "USD" && rate["to"] == "EUR" {
						if math.Abs(rate["rate"].(float64)-1.2) > 0.0001 {
							t.Errorf("Expected rate from USD to EUR to be 1.2, got %v", rate["rate"])
						}
					}
					if rate["from"] == "USD" && rate["to"] == "GBP" {
						if math.Abs(rate["rate"].(float64)-1.5) > 0.0001 {
							t.Errorf("Expected rate from USD to GBP to be 1.5, got %v", rate["rate"])
						}
					}
					if rate["from"] == "EUR" && rate["to"] == "GBP" {
						if math.Abs(rate["rate"].(float64)-1.5/1.2) > 0.0001 {
							t.Errorf("Expected rate from EUR to GBP to be 1.5/1.2, got %v", rate["rate"])
						}
					}
					if rate["from"] == "GBP" && rate["to"] == "EUR" {
						if math.Abs(rate["rate"].(float64)-1.2/1.5) > 0.0001 {
							t.Errorf("Expected rate from GBP to EUR to be 1.2/1.5, got %v", rate["rate"])
						}
					}
				}
			},
		},
		{
			name:           "Invalid Request - No Currencies",
			query:          "",
			expectedStatus: http.StatusBadRequest,
			expectedLen:    0,
		},
		{
			name:           "Invalid Request - Single Currency",
			query:          "currencies=USD",
			expectedStatus: http.StatusBadRequest,
			expectedLen:    0,
		},
		{
			name:           "Invalid Request - API Error",
			query:          "currencies=USD,EUR,GBP",
			expectedStatus: http.StatusBadRequest,
			expectedLen:    0,
			withApiError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.withApiError {
				testClient.err = fmt.Errorf("API error")
			} else {
				testClient.err = nil
			}

			req, _ := http.NewRequest(http.MethodGet, "/rates?"+tc.query, nil)
			rec := httptest.NewRecorder()

			s.router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedStatus == http.StatusOK {
				var res []map[string]any

				err := json.Unmarshal(rec.Body.Bytes(), &res)
				if err != nil {
					t.Errorf("Failed to parse response body: %v", err)
					return
				}

				if len(res) != tc.expectedLen {
					t.Errorf("Expected length %d, got %d", tc.expectedLen, len(res))
				}

				if tc.validation != nil {
					tc.validation(t, res)
				}
			}
		})
	}
}
