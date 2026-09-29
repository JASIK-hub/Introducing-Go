package ecbank

import (
	"errors"
	"fmt"
	"moneyconverter/money"
	"net/http"
)

type Client struct{}

var ErrServerSide = errors.New("error calling server")

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
