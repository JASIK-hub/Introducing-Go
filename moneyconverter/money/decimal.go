package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidDecimal = errors.New("unable to convert the decimal")
	ErrTooLarge       = errors.New("quantity over 10^12 is too large")
)

const maxDecimal = 1e12

type Decimal struct {
	subunits  int64
	precision byte
}

func ParseDecimal(value string) (Decimal, error) {
	intPart, fracPart, hasDot := strings.Cut(value, ".")

	if hasDot && fracPart == "" {
		return Decimal{}, ErrInvalidDecimal
	}
	subunits, err := strconv.ParseInt(intPart+fracPart, 10, 64)

	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	if subunits > maxDecimal {
		return Decimal{}, ErrTooLarge
	}
	precision := byte(len(fracPart))

	return Decimal{subunits: subunits, precision: precision}, nil
}

func (d *Decimal) String() string {
	if d.precision == 0 {
		return fmt.Sprintf("%d", d.subunits)
	}
	centsPerUnit := pow10(d.precision)
	frac := d.subunits % centsPerUnit
	integer := d.subunits / centsPerUnit
	format := "%d.%0" + strconv.Itoa(int(d.precision)) + "d"
	return fmt.Sprintf(format, integer, frac)
}
