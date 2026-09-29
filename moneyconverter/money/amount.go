package money

import (
	"errors"
)

type Amount struct {
	quantity Decimal
	currency Currency
}

var (
	ErrTooPrecise = errors.New("quantity is too precise")
)

func NewAmount(quantity Decimal, currency Currency) (Amount, error) {
	if quantity.precision > currency.precision {
		return Amount{}, ErrTooPrecise
	}
	quantity.subunits *= pow10(currency.precision - quantity.precision)
	quantity.precision = currency.precision
	return Amount{quantity: quantity, currency: currency}, nil
}

func (a Amount) validate() error {
	switch {
	case a.quantity.subunits > maxDecimal:
		return ErrTooLarge
	case a.quantity.precision > a.currency.precision:
		return ErrTooPrecise
	}
	return nil
}

func (a Amount) String() string {
	return a.quantity.String() + " " + a.currency.code
}
