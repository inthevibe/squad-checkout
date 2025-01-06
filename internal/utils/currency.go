package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Currency struct {
	CountryCurrencyDesc string `json:"country_currency_desc"`
}

func GetSupportedCurrencies() ([]string, error) {
	urlStr := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange?fields=country_currency_desc&page[size]=300"

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(urlStr)
	if err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) || strings.Contains(err.Error(), "timeout") {
			return nil, errors.New("API unresponsive: request timed out")
		}
		return nil, fmt.Errorf("failed to fetch supported currencies: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Data []Currency `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	var currencies []string
	for _, currency := range result.Data {
		if currency.CountryCurrencyDesc == strings.ToUpper(currency.CountryCurrencyDesc) {
			continue
		}
		currencies = append(currencies, currency.CountryCurrencyDesc)
	}

	return currencies, nil
}

type ExchangeRateResponse struct {
	Data []struct {
		CountryCurrencyDesc string `json:"country_currency_desc"`
		ExchangeRate        string `json:"exchange_rate"`
		RecordDate          string `json:"record_date"`
	} `json:"data"`
}

var ApiURL = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"

func GetExchangeRate(currency string, date time.Time) (float64, error) {
	encodedCurrency := url.QueryEscape(currency)
	startDate := date.AddDate(0, -6, 0).Format("2006-01-02")
	dateStr := date.Format("2006-01-02")

	urlStr := fmt.Sprintf(
		"%s?fields=country_currency_desc,exchange_rate,record_date&filter=country_currency_desc:eq:%s,record_date:gte:%s,record_date:lte:%s&sort=-record_date&page[size]=1",
		ApiURL, encodedCurrency, startDate, dateStr,
	)

	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(urlStr)
	if err != nil {
		if errors.Is(err, http.ErrHandlerTimeout) || strings.Contains(err.Error(), "timeout") {
			return 0, errors.New("API unresponsive: request timed out")
		}
		return 0, fmt.Errorf("failed to fetch exchange rate: %v", err)
	}
	defer resp.Body.Close()

	var result ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode API response: %v", err)
	}

	if len(result.Data) == 0 {
		return 0, fmt.Errorf("no exchange rate found for %s between %s and %s", encodedCurrency, startDate, dateStr)
	}

	exchangeRate, err := strconv.ParseFloat(result.Data[0].ExchangeRate, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse exchange rate: %v", err)
	}

	return exchangeRate, nil
}
