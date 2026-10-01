package ecbank

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"moneyconverter/money"
	"net/http"
)

const baseCurrencyCode = "EUR"

const (
	clientErrorClass = 4
	serverErrorClass = 5
)

var (
	ErrServerSide         = errors.New("error calling server")
	ErrClientSide         = errors.New("HTTP client error")
	ErrUnknownStatusCode  = errors.New("unexpected HTTP status code")
	ErrUnexpectedFormat   = errors.New("not valid format for response")
	ErrChangeRateNotFound = errors.New("couldn't find the exchange rate")
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
		return money.ExchangeRate{}, fmt.Errorf("failed to find the source currency %s", source)
	}
	targetValue, targetFound := rates[target]
	if !targetFound {
		return money.ExchangeRate{}, fmt.Errorf("failed to find the source currency %s", target)
	}
	rate, err := money.ParseDecimal(fmt.Sprintf("%.10f", targetValue/sourceValue))
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("unable to parse exchange rate from %s to %s: %w", source, target, err)
	}
	return money.ExchangeRate(rate), nil
}

func readRateFromResponse(source, target string, respBody io.Reader) (money.ExchangeRate, error) {
	decoder := xml.NewDecoder(respBody)
	var response envelope
	err := decoder.Decode(&response)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrUnexpectedFormat, err)
	}
	rate, err := response.exchangeRate(source, target)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrChangeRateNotFound, err)
	}
	return rate, nil
}

func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {

	const path = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	resp, err := http.Get(path)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s",
			ErrServerSide, err.Error())
	}
	defer resp.Body.Close()
	if err = checkStatusCode(resp.StatusCode); err != nil {
		return money.ExchangeRate{}, err
	}
	rate, err := readRateFromResponse(source.Code(), target.Code(), resp.Body)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	return rate, nil
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
