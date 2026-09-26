package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Decimal struct {
	subunits  int64
	precision byte
}

var (
	ErrInvalidDecimal = errors.New("unable to convert the decimal")
	ErrTooLarge       = errors.New("quantity over 10^12 is too large")
)

const maxDecimal = 1e12

func ParseDecimal(value string) (Decimal, error) {
	intPart, fracPart, _ := strings.Cut(value, ".")

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
