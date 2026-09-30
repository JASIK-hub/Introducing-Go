package ecbank

import (
	"errors"
	"fmt"
	"moneyconverter/money"
	"net/http"
)

type Client struct{}

const (
	clientErrorClass = 4
	serverErrorClass = 5
)

var (
	ErrServerSide        = errors.New("error calling server")
	ErrClientSide        = errors.New("HTTP client error")
	ErrUnknownStatusCode = errors.New("unexpected HTTP status code")
)

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
