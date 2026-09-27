package money

import (
	"errors"
)

type Amount struct {
	quantity Decimal
	currency Currency
}

var (
	ErrTooPresice = errors.New("quantity is too precise")
)

func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
	if quantity.precision > currency.precision {
		return Amount{}, ErrTooPresice
	}
	quantity.precision = currency.precision
	return Amount{quantity: quantity, currency: currency}, nil
}
