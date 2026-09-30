package ecbank

import (
	"errors"
	"fmt"
	"moneyconverter/money"
	"net/http"
)

const baseCurrencyCode = "EUR"

const (
	clientErrorClass = 4
	serverErrorClass = 5
)

var (
	ErrServerSide        = errors.New("error calling server")
	ErrClientSide        = errors.New("HTTP client error")
	ErrUnknownStatusCode = errors.New("unexpected HTTP status code")
)

type Client struct{}
type envelope struct {
	Rates []currencyRate `xml:"Cube>Cube>Cube"`
}

type currencyRate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

func (e envelope) mappedExchangeRates() map[string]float64 {
	rates := make(map[string]float64, len(e.Rates))

	for _, c := range e.Rates {
		rates[c.Currency] = c.Rate
	}

	rates[baseCurrencyCode] = 1
	return rates
}

func (e envelope) exchangeRate(source, target string) (money.ExchangeRate, error) {
	if source == target {
		one, err := money.ParseDecimal("1")
		if err != nil {
			return money.ExchangeRate{}, err
		}
		return money.ExchangeRate(one), nil
	}
	rates := e.mappedExchangeRates()

	sourceValue, sourceFound := rates[source]
	if !sourceFound {
		return money.ExchangeRate{}, fmt.Errorf("failed to find the source currency %s", sourceValue)
	}
	targetValue, targetFound := rates[target]
	if !targetFound {
		return money.ExchangeRate{}, fmt.Errorf("failed to find the source currency %s", targetValue)
	}
	rate, err := money.ParseDecimal(fmt.Sprintf("%.10f", targetValue/sourceValue))
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("unable to parse exchange rate from %s to %s: %w", source, target, err)
	}
	return money.ExchangeRate(rate), nil
}

func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {

	const path = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	resp, err := http.Get(path)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s",
			ErrServerSide, err.Error())
	}
	defer resp.Body.Close()

	return money.ExchangeRate{}, nil
}

func checkStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil
	case httpStatusCode(statusCode) == clientErrorClass:
		return fmt.Errorf("%w: %d", ErrClientSide, statusCode)
	case httpStatusCode(statusCode) == serverErrorClass:
		return fmt.Errorf("%w: %d", ErrServerSide, statusCode)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownStatusCode, statusCode)

	}
}

func httpStatusCode(statusCode int) int {
	const httpErrorClassSize = 100
	return statusCode / httpErrorClassSize
}
