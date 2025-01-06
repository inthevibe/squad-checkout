package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"squad-checkout/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetExchangeRate_Unit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "Brazil-Real") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":[{"country_currency_desc":"Brazil-Real","exchange_rate":"5.434","record_date":"2024-10-01"}]}`))
		} else if strings.Contains(r.URL.String(), "Invalid-Currency") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":[]}`))
		} else if strings.Contains(r.URL.String(), "Timeout-Currency") {
			time.Sleep(3 * time.Second)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`Internal Server Error`))
		}
	}))
	defer server.Close()

	originalApiURL := utils.ApiURL
	utils.ApiURL = server.URL
	defer func() { utils.ApiURL = originalApiURL }()

	tests := []struct {
		name          string
		currency      string
		date          time.Time
		expectedRate  float64
		expectedError error
	}{
		{
			name:          "Valid exchange rate",
			currency:      "Brazil-Real",
			date:          time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedRate:  5.434,
			expectedError: nil,
		},
		{
			name:          "No exchange rate found",
			currency:      "Invalid-Currency",
			date:          time.Date(2030, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedRate:  0,
			expectedError: errors.New("no exchange rate found"),
		},
		{
			name:          "API timeout",
			currency:      "Timeout-Currency",
			date:          time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedRate:  0,
			expectedError: errors.New("Timeout"),
		},
		{
			name:          "API server error",
			currency:      "Error-Currency",
			date:          time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedRate:  0,
			expectedError: errors.New("failed to decode API response"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, err := utils.GetExchangeRate(tt.currency, tt.date)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRate, rate)
			}
		})
	}
}
