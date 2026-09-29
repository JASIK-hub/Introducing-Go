package money

import "math"

type ExchangeRate Decimal

func Convert(amount Amount, to Currency) (Amount, error) {
	convertedValue := applyExchangeRate(amount, to, ExchangeRate{subunits: 2, precision: 0})
	if err := convertedValue.validate(); err != nil {
		return Amount{}, err
	}
	return convertedValue, nil
}

func pow10(power byte) int64 {
	switch power {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	default:
		return int64(math.Pow(10, (float64(power))))
	}
}

func (d Decimal) simplify() Decimal {
	for d.precision > 0 && d.subunits%10 == 0 {
		d.subunits /= 10
		d.precision--
	}
	return d
}

func multiply(d Decimal, r ExchangeRate) Decimal {
	dec := Decimal{
		subunits:  d.subunits * r.subunits,
		precision: d.precision + r.precision,
	}
	dec.simplify()
	return dec
}

func applyExchangeRate(a Amount, target Currency, r ExchangeRate) Amount {
	converted := multiply(a.quantity, r)
	if converted.precision > target.precision {
		converted.subunits /= pow10(converted.precision - target.precision)
	} else if converted.precision < target.precision {
		converted.subunits *= pow10(target.precision - converted.precision)
	}
	converted.precision = target.precision
	return Amount{
		currency: target,
		quantity: converted,
	}
}
