package tests

import (
	"errors"
	"squad-checkout/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetExchangeRate_RealAPI(t *testing.T) {
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
			date:          time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedRate:  0,
			expectedError: errors.New("no exchange rate found"),
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
